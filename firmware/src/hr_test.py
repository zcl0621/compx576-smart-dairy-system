"""MAX30102 heart rate smoke test (I2C: SDA=GPIO21, SCL=GPIO22).

Run on device:
    import hr_test

Put a fingertip lightly on the sensor (red LED should be on).
IR value should jump from ~1000 (no finger) to >50000 (finger on).
Beats show as a slow up-down ripple in IR around its DC level.
"""
from machine import SoftI2C, Pin
import time
from max30102 import MAX30102

i2c = SoftI2C(scl=Pin(22), sda=Pin(21), freq=100000)
sensor = MAX30102(i2c=i2c)

if sensor.i2c_address not in i2c.scan():
    raise RuntimeError("MAX30102 not on I2C bus")
if not sensor.check_part_id():
    raise RuntimeError("wrong part ID — not a MAX30102?")

for attempt in range(5):
    try:
        sensor.setup_sensor()
        break
    except OSError:
        time.sleep_ms(200)
sensor.set_sample_rate(100)
sensor.set_fifo_average(4)
sensor.set_active_leds_amplitude(0x7F)
FINGER_THRESHOLD = 5000

print("sensor ready. put finger on now.")

# simple beat detector: trigger when IR rises through (avg + delta)
last_beat = 0
ir_avg = 0
beats = []
finger_on = False

while True:
    sensor.check()
    while sensor.available():
        red = sensor.pop_red_from_storage()
        ir = sensor.pop_ir_from_storage()

        # finger presence
        on_now = ir > FINGER_THRESHOLD
        if on_now != finger_on:
            finger_on = on_now
            print("finger:", "ON" if on_now else "OFF")
            ir_avg = ir
            beats = []
            continue
        if not finger_on:
            continue

        # exponential moving average as DC baseline
        ir_avg = (ir_avg * 15 + ir) // 16
        # beat when AC component crosses threshold rising
        now = time.ticks_ms()
        if ir > ir_avg + 50 and time.ticks_diff(now, last_beat) > 300:
            interval = time.ticks_diff(now, last_beat)
            last_beat = now
            if 300 < interval < 2000:  # 30-200 bpm
                beats.append(interval)
                if len(beats) > 5:
                    beats.pop(0)
                avg_ms = sum(beats) // len(beats)
                bpm = 60000 // avg_ms
                print("beat! interval={}ms  bpm~{}".format(interval, bpm))
