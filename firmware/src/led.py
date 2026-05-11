"""Status LED on GPIO2 (onboard blue). Non-blocking patterns driven by tick()."""
from machine import Pin
import time

_pin = Pin(2, Pin.OUT, value=0)
_pattern = "off"
_last_tick = 0
_state = 0  # phase within the pattern

# patterns: list of (on/off, ms) cycles
_PATTERNS = {
    "off":      [(0, 1000)],
    "on":       [(1, 1000)],
    "boot":     [(1, 100), (0, 100)],                          # fast blink
    "wifi":     [(1, 500), (0, 500)],                          # slow blink
    "ok":       [(1, 100), (0, 4900)],                         # heartbeat every 5 s
    "post_ok":  [(1, 100), (0, 100)] * 2 + [(0, 1000)],        # 2 short flashes
    "error":    [(1, 1000), (0, 1000)],                        # 1 Hz heavy blink
}


def set_pattern(name):
    global _pattern, _state, _last_tick
    if name == _pattern:
        return
    _pattern = name
    _state = 0
    _last_tick = time.ticks_ms()
    seq = _PATTERNS.get(name) or _PATTERNS["off"]
    _pin.value(seq[0][0])


def tick():
    """Call frequently from main loop to advance the pattern."""
    global _state, _last_tick
    seq = _PATTERNS.get(_pattern) or _PATTERNS["off"]
    _, dur = seq[_state]
    if time.ticks_diff(time.ticks_ms(), _last_tick) >= dur:
        _state = (_state + 1) % len(seq)
        _last_tick = time.ticks_ms()
        _pin.value(seq[_state][0])


def flash(times_on_ms=100):
    """Blocking quick flash — for one-shot signals like POST success."""
    _pin.value(1)
    time.sleep_ms(times_on_ms)
    _pin.value(0)
