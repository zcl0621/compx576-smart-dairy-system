"""WiFi helper — connect with timeout, return IP or raise."""
import network
import time


def connect(ssid, password, timeout_s=20):
    wlan = network.WLAN(network.STA_IF)
    wlan.active(True)
    if wlan.isconnected():
        return wlan.ifconfig()[0]

    print("wifi: connecting to {!r}...".format(ssid))
    wlan.connect(ssid, password)

    t0 = time.ticks_ms()
    while not wlan.isconnected():
        if time.ticks_diff(time.ticks_ms(), t0) > timeout_s * 1000:
            raise RuntimeError("wifi: timed out connecting to {}".format(ssid))
        time.sleep_ms(200)

    ip = wlan.ifconfig()[0]
    print("wifi: ok, ip={}".format(ip))
    return ip
