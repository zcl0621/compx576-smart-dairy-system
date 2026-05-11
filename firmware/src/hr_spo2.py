"""MAX30102 BPM + SpO2 with HP filter and peak detection.

Approach:
- 25 Hz effective sample rate (100 Hz / 4-sample FIFO average)
- 100-sample sliding window (4 s) for both red and IR
- DC = rolling mean, AC = window peak-to-peak
- HR: detect peaks on (IR - DC) with adaptive threshold, derive BPM from intervals
- SpO2: R = (AC_red/DC_red)/(AC_ir/DC_ir),  SpO2 = 110 - 25*R (approx)

Run: import hr_spo2
"""
import sys
import time

sys.path.append("/lib")
from machine import SoftI2C, Pin
from max30102 import MAX30102

FS = 25  # effective Hz
WIN = 100  # 4 s
FINGER_IR = 5000
MIN_BEAT_MS = 350  # ~170 bpm cap

i2c = SoftI2C(scl=Pin(22), sda=Pin(21), freq=100000)
sensor = MAX30102(i2c=i2c)
for _ in range(5):
    try:
        sensor.setup_sensor()
        break
    except OSError:
        time.sleep_ms(200)
sensor.set_sample_rate(100)
sensor.set_fifo_average(4)
sensor.set_active_leds_amplitude(0xFF)  # ~50mA, max signal
print("ready. place finger lightly on the window.")


def mean(xs):
    return sum(xs) // len(xs)


def smooth(xs, k=3):
    """k-point moving average (low-pass)."""
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
    """Return indices of local maxima above adaptive threshold."""
    if len(sig) < 5:
        return []
    hi = max(sig)
    lo = min(sig)
    if hi - lo < 50:  # too flat — no real pulse
        return []
    threshold = lo + (hi - lo) * 7 // 10  # 70% of peak-to-peak
    peaks = []
    for i in range(2, len(sig) - 2):
        if (sig[i] > threshold
            and sig[i] > sig[i - 1] and sig[i] > sig[i - 2]
            and sig[i] >= sig[i + 1] and sig[i] >= sig[i + 2]):
            if not peaks or (i - peaks[-1]) >= min_gap:
                peaks.append(i)
    return peaks


def median(xs):
    s = sorted(xs)
    n = len(s)
    return s[n // 2] if n % 2 else (s[n // 2 - 1] + s[n // 2]) // 2


buf_ir = []
buf_red = []
last_report = time.ticks_ms()
min_gap_samples = MIN_BEAT_MS * FS // 1000  # ~9

while True:
    sensor.check()
    while sensor.available():
        red = sensor.pop_red_from_storage()
        ir = sensor.pop_ir_from_storage()

        if ir < FINGER_IR:
            buf_ir = []
            buf_red = []
            if time.ticks_diff(time.ticks_ms(), last_report) > 2000:
                print("no finger (ir={})".format(ir))
                last_report = time.ticks_ms()
            continue

        buf_ir.append(ir)
        buf_red.append(red)
        if len(buf_ir) > WIN:
            buf_ir.pop(0)
            buf_red.pop(0)

        if len(buf_ir) < WIN:
            continue

        if time.ticks_diff(time.ticks_ms(), last_report) < 1000:
            continue
        last_report = time.ticks_ms()

        # smooth raw signals to kill high-freq noise
        sm_ir = smooth(buf_ir, 5)
        sm_red = smooth(buf_red, 5)

        dc_ir = mean(sm_ir)
        dc_red = mean(sm_red)
        ac_ir = max(sm_ir) - min(sm_ir)
        ac_red = max(sm_red) - min(sm_red)

        # signal quality gate
        if ac_ir < 100:
            print("weak signal (ir_ac={}) — press finger more firmly".format(ac_ir))
            continue

        # SpO2 — R ratio method
        spo2 = None
        if ac_ir > 0 and dc_ir > 0 and dc_red > 0:
            r = (ac_red * dc_ir * 1000) // (dc_red * ac_ir)
            spo2 = 110 - (25 * r) // 1000
            if spo2 > 100:
                spo2 = 100
            if spo2 < 70:
                spo2 = None

        # HR — peaks on HP-filtered smoothed IR, use median of intervals
        hp = [x - dc_ir for x in sm_ir]
        peaks = detect_peaks(hp, min_gap_samples)
        bpm = None
        if len(peaks) >= 3:
            gaps = [peaks[j + 1] - peaks[j] for j in range(len(peaks) - 1)]
            m_gap = median(gaps)
            bpm = 60 * FS // m_gap
            if bpm < 40 or bpm > 180:
                bpm = None

        print(
            "ir_ac={:4d} red_ac={:4d} peaks={:2d}  bpm={}  spo2={}".format(
                ac_ir, ac_red, len(peaks),
                bpm if bpm else "--",
                "{}%".format(spo2) if spo2 else "--",
            )
        )
