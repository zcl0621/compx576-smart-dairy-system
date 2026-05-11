# Firmware

ESP32 + MicroPython for the dairy cow collar demo.

## Hardware

- ESP32 dev board + sensor shield
- DS18B20 temp sensor — plugged into shield slot **D14** (S/V/G)
  - S → GPIO14, V → 3.3V, G → GND
  - Shield handles the 4.7k pull-up. If reading 85.0°C, add an external 4.7k between S and V.

## Layout

```
firmware/
├── pyproject.toml        # host-side tools (esptool, mpremote) via uv
├── flash.sh              # one-shot flash + upload + repl
├── firmware-bin/         # drop MicroPython .bin here
└── src/                  # files uploaded to ESP32
    ├── boot.py
    ├── main.py
    └── lib/              # third-party micropython libs
```

## First flash

1. Download MicroPython firmware for ESP32 from <https://micropython.org/download/ESP32_GENERIC/> and save it to `firmware-bin/ESP32_GENERIC.bin`.
2. Plug in the ESP32 over USB.
3. Run everything:

```bash
./flash.sh all
```

Or step by step:

```bash
./flash.sh erase       # wipe flash
./flash.sh firmware    # burn MicroPython
./flash.sh upload      # copy src/ to device
./flash.sh repl        # open serial REPL
```

Override port if auto-detect fails:

```bash
PORT=/dev/cu.usbserial-0001 ./flash.sh upload
```

## Expected output

```
found devices: ['28ff...']
rom=28ff... temp=24.31C
rom=28ff... temp=24.33C
```

## Troubleshooting

| Symptom | Cause |
|---|---|
| `no DS18B20 found` | wiring wrong, S/V swapped, or DQ not on GPIO14 |
| reads always `85.00` | missing pull-up, or sensor not powered |
| reads `-127.00` | bad contact on data line |
| port not found | install CP210x or CH340 USB driver, replug cable |

## Exit REPL

`Ctrl-X` to exit `mpremote`. `Ctrl-C` inside REPL stops `main.py`.
