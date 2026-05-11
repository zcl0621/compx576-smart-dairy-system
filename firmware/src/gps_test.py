"""NEO-6M GPS smoke test on UART2 (D16=RX, D17=TX).

Run on device:
    import gps_test

Prints raw NMEA lines first. Cold start outdoor fix takes 30s-2min.
Indoor near window may take 5+ min or never fix. Look for $GPRMC or $GPGGA
with non-empty lat/lon fields.
"""
from machine import UART
import time

uart = UART(2, baudrate=9600, tx=17, rx=16, timeout=200)
print("listening on UART2 (rx=GPIO16, tx=GPIO17) @ 9600 baud")
print("waiting for NMEA sentences...")

buf = b""
last_print = time.ticks_ms()
sentence_count = 0

while True:
    data = uart.read()
    if data:
        buf += data
        while b"\n" in buf:
            line, buf = buf.split(b"\n", 1)
            line = line.strip()
            if not line:
                continue
            try:
                s = line.decode("ascii")
            except Exception:
                continue
            print(s)
            sentence_count += 1
            # quick lat/lon extract from $GPRMC
            if s.startswith("$GPRMC") or s.startswith("$GNRMC"):
                f = s.split(",")
                if len(f) > 6 and f[2] == "A":
                    print(">>> FIX: lat={}{} lon={}{}".format(f[3], f[4], f[5], f[6]))
                elif len(f) > 2:
                    print(">>> no fix yet (status={})".format(f[2] or "empty"))
    if time.ticks_diff(time.ticks_ms(), last_print) > 5000:
        if sentence_count == 0:
            print("(no data yet — check TX/RX wiring, GPS power LED should blink)")
        last_print = time.ticks_ms()
        sentence_count = 0
    time.sleep_ms(50)
