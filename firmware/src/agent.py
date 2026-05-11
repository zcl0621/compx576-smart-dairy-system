"""Dairy cow collar agent.

Each boot:
  1. POST previous sample's report from RTC memory (fresh lwIP, no I2C yet)
  2. Sample DS18B20 + NEO-6M + MAX30102 for SAMPLE_MS
  3. Stash report in RTC memory, machine.reset() to fully clear lwIP/TLS state

This avoids EBUSY on HTTPS POSTs after MAX30102 SoftI2C polling, which
corrupts lwIP state in a way that disconnect/reconnect cannot recover.

Run: import agent
"""
import sys
import time
import json
import gc
import machine

sys.path.append("/lib")

import onewire
import ds18x20
from machine import UART, SoftI2C, Pin
from max30102 import MAX30102

import config
import wifi
import transport
import led


# ---- config ----
TEMP_PIN = 14
GPS_RX, GPS_TX = 16, 17
I2C_SDA, I2C_SCL = 21, 22
SAMPLE_MS = 20_000
HR_WINDOW_SAMPLES = 100  # 4 s at 25 Hz effective
HR_WINDOW_MS = 4000
MIN_GAP_SAMPLES = 9  # ~170 bpm cap at 25 Hz
FINGER_IR_MIN = 5000

# valid ranges — sensor sanity check, not health alerts
# (cow body temp 38-39; HR 60-80; SpO2 >95)
TEMP_MIN, TEMP_MAX = 0.0, 50.0
BPM_MIN, BPM_MAX = 40, 180
SPO2_MIN, SPO2_MAX = 70, 100


# ---- Phase 1: POST cached report from RTC memory BEFORE any I2C ----
rtc = machine.RTC()
try:
    stashed = rtc.memory()
except Exception:
    stashed = b""

if stashed:
    try:
        cached = json.loads(stashed.decode("utf-8"))
        print("found cached report, connecting wifi to post...")
        led.set_pattern("wifi")
        wifi.connect(config.WIFI_SSID, config.WIFI_PASSWORD)
        led.set_pattern("ok")
        print(json.dumps(cached))
        sent, failed = transport.send_report(
            config.AGENT_BASE_URL, config.COW_ID, config.AGENT_TOKEN, cached
        )
        print("posted: {} ok, {} failed".format(sent, failed))
        if sent > 0 and failed == 0:
            led.flash(100)
    except Exception as e:
        print("post cached report err:", e)
    finally:
        rtc.memory(b"")  # clear regardless so we don't retry forever


# ---- helpers ----
def median(xs):
    if not xs:
        return None
    s = sorted(xs)
    n = len(s)
    return s[n // 2] if n % 2 else (s[n // 2 - 1] + s[n // 2]) / 2


def smooth(xs, k=5):
    out = []
    half = k // 2
    for i in range(len(xs)):
        s = 0
        c = 0
        for j in range(max(0, i - half), min(len(xs), i + half + 1)):
            s += xs[j]
            c += 1
        out.append(s // c)
    return out


def detect_peaks(sig, min_gap):
    if len(sig) < 5:
        return []
    hi = max(sig)
    lo = min(sig)
    if hi - lo < 50:
        return []
    threshold = lo + (hi - lo) * 4 // 10
    peaks = []
    for i in range(2, len(sig) - 2):
        if (sig[i] > threshold
            and sig[i] > sig[i - 1] and sig[i] > sig[i - 2]
            and sig[i] >= sig[i + 1] and sig[i] >= sig[i + 2]):
            if not peaks or (i - peaks[-1]) >= min_gap:
                peaks.append(i)
    return peaks


# ---- sensor init ----
ow = onewire.OneWire(Pin(TEMP_PIN))
ds = ds18x20.DS18X20(ow)
roms = ds.scan()
if not roms:
    raise RuntimeError("DS18B20 not found")
print("ds18b20 ok, rom=", roms[0].hex())

gps_uart = UART(2, baudrate=9600, tx=GPS_TX, rx=GPS_RX, timeout=0)
print("gps uart2 ready")

i2c = SoftI2C(scl=Pin(I2C_SCL), sda=Pin(I2C_SDA), freq=100000)
max_sensor = MAX30102(i2c=i2c)
for _ in range(5):
    try:
        max_sensor.setup_sensor()
        break
    except OSError:
        time.sleep_ms(200)
max_sensor.set_sample_rate(100)
max_sensor.set_fifo_average(4)
max_sensor.set_active_leds_amplitude(0xFF)
print("max30102 ok")


# ---- sample buffers ----
cycle = {"temps": [], "lats": [], "lons": [], "alts": [], "sats": [], "bpms": [], "spo2s": []}
last_temp_request = 0
temp_convert_started = False
gps_buf = b""
hr_buf_ir = []
hr_buf_red = []
last_hr_window = time.ticks_ms()


def parse_nmea_rmc(s):
    """Return (lat_deg, lon_deg) if RMC has fix, else None."""
    if not (s.startswith("$GPRMC") or s.startswith("$GNRMC")):
        return None
    f = s.split(",")
    if len(f) < 7 or f[2] != "A":
        return None
    try:
        lat = nmea_to_deg(f[3], f[4])
        lon = nmea_to_deg(f[5], f[6])
        return (lat, lon)
    except Exception:
        return None


def parse_nmea_gga(s):
    """Return (alt_m, num_sats) if GGA has fix."""
    if not (s.startswith("$GPGGA") or s.startswith("$GNGGA")):
        return None
    f = s.split(",")
    if len(f) < 10 or not f[6] or f[6] == "0":
        return None
    try:
        sats = int(f[7])
        alt = float(f[9])
        return (alt, sats)
    except Exception:
        return None


def nmea_to_deg(value, hemi):
    if not value:
        return None
    # ddmm.mmmm or dddmm.mmmm
    dot = value.find(".")
    deg = int(value[: dot - 2])
    minute = float(value[dot - 2 :])
    d = deg + minute / 60.0
    if hemi in ("S", "W"):
        d = -d
    return round(d, 6)


# ---- WiFi ----
led.set_pattern("wifi")
try:
    wifi.connect(config.WIFI_SSID, config.WIFI_PASSWORD)
except Exception as e:
    print("wifi failed:", e, "— continuing offline, reports will print only")
    led.set_pattern("error")

led.set_pattern("ok")
print("agent started, sampling for {}s then reboot".format(SAMPLE_MS // 1000))


# ---- one-shot sample window ----
sample_start = time.ticks_ms()
while time.ticks_diff(time.ticks_ms(), sample_start) < SAMPLE_MS:
    now = time.ticks_ms()

    # DS18B20 (1 Hz, non-blocking 2-step)
    if not temp_convert_started and time.ticks_diff(now, last_temp_request) >= 1000:
        try:
            ds.convert_temp()
            temp_convert_started = True
            last_temp_request = now
        except Exception as e:
            print("temp convert err:", e)
    elif temp_convert_started and time.ticks_diff(now, last_temp_request) >= 800:
        try:
            t = ds.read_temp(roms[0])
            if t is not None and TEMP_MIN <= t <= TEMP_MAX:
                cycle["temps"].append(t)
        except Exception:
            pass
        temp_convert_started = False

    # GPS (drain UART, defer parse)
    n_avail = gps_uart.any()
    d = gps_uart.read(n_avail) if n_avail else None
    if d:
        gps_buf = (gps_buf + d)[-1024:]

    # MAX30102 sampling + FIFO drain
    max_sensor.check()
    while max_sensor.available():
        red = max_sensor.pop_red_from_storage()
        ir = max_sensor.pop_ir_from_storage()
        hr_buf_ir.append(ir)
        hr_buf_red.append(red)
        if len(hr_buf_ir) > HR_WINDOW_SAMPLES:
            hr_buf_ir.pop(0)
            hr_buf_red.pop(0)

    # HR / SpO2 compute on rolling 4s window
    if (len(hr_buf_ir) >= HR_WINDOW_SAMPLES
        and time.ticks_diff(now, last_hr_window) >= HR_WINDOW_MS):
        last_hr_window = now
        sm_ir = smooth(hr_buf_ir, 5)
        sm_red = smooth(hr_buf_red, 5)
        dc_ir = sum(sm_ir) // len(sm_ir)
        dc_red = sum(sm_red) // len(sm_red)
        ac_ir = max(sm_ir) - min(sm_ir)
        ac_red = max(sm_red) - min(sm_red)
        if dc_ir >= FINGER_IR_MIN and ac_ir >= 100:
            if dc_red > 0 and ac_ir > 0:
                r = (ac_red * dc_ir * 1000) // (dc_red * ac_ir)
                sp = 110 - (25 * r) // 1000
                if SPO2_MIN <= sp <= SPO2_MAX:
                    cycle["spo2s"].append(sp)
            hp = [x - dc_ir for x in sm_ir]
            peaks = detect_peaks(hp, MIN_GAP_SAMPLES)
            if len(peaks) >= 2:
                gaps = [peaks[j + 1] - peaks[j] for j in range(len(peaks) - 1)]
                m_gap = median(gaps)
                if m_gap:
                    bpm = int(60 * 25 // m_gap)
                    if BPM_MIN <= bpm <= BPM_MAX:
                        cycle["bpms"].append(bpm)

    led.tick()
    time.sleep_ms(20)


# ---- sampling done, parse NMEA + build report ----
for raw in gps_buf.split(b"\n"):
    if not (raw.startswith(b"$GPRMC") or raw.startswith(b"$GNRMC")
            or raw.startswith(b"$GPGGA") or raw.startswith(b"$GNGGA")):
        continue
    try:
        s = raw.strip().decode("ascii")
    except Exception:
        continue
    ll = parse_nmea_rmc(s)
    if ll:
        cycle["lats"].append(ll[0])
        cycle["lons"].append(ll[1])
        continue
    gga = parse_nmea_gga(s)
    if gga:
        cycle["alts"].append(gga[0])
        cycle["sats"].append(gga[1])

report = {
    "cow_id": config.COW_ID,
    "ts": time.time(),
    "temp_c": round(median(cycle["temps"]), 2) if cycle["temps"] else None,
    "bpm": int(median(cycle["bpms"])) if cycle["bpms"] else None,
    "spo2": int(median(cycle["spo2s"])) if cycle["spo2s"] else None,
    "lat": median(cycle["lats"]),
    "lon": median(cycle["lons"]),
    "alt_m": round(median(cycle["alts"]), 1) if cycle["alts"] else None,
    "sats": int(median(cycle["sats"])) if cycle["sats"] else 0,
    "samples": {
        "temp": len(cycle["temps"]),
        "gps": len(cycle["lats"]),
        "hr": len(cycle["bpms"]),
        "spo2": len(cycle["spo2s"]),
    },
}
print(json.dumps(report))
gc.collect()

# Stash report in RTC memory; next boot POSTs it before any I2C runs.
try:
    rtc.memory(json.dumps(report).encode("utf-8"))
    print("stashed report in RTC memory, rebooting...")
except Exception as e:
    print("rtc stash err:", e)

time.sleep_ms(500)
machine.reset()
