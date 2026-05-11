#!/usr/bin/env bash
# Flash MicroPython firmware and upload src/ to ESP32.
# Usage: ./flash.sh [erase|firmware|upload|reset|repl|dev|all]
#   dev = upload + reset + repl   (default; common loop, keeps firmware)
#   all = erase + firmware + upload + repl  (full bring-up from blank chip)
set -euo pipefail

PORT="${PORT:-$(ls /dev/cu.usbserial-* /dev/cu.SLAB_USBtoUART /dev/cu.wchusbserial* 2>/dev/null | head -n1)}"
BAUD="${BAUD:-460800}"
FW="${FW:-firmware-bin/ESP32_GENERIC.bin}"

if [ -z "${PORT:-}" ]; then
  echo "no serial port found, plug in ESP32 or set PORT=..."
  exit 1
fi

echo "port: $PORT"

erase()    { uv run esptool.py --chip esp32 --port "$PORT" erase_flash; }
firmware() {
  [ -f "$FW" ] || { echo "missing $FW, download from https://micropython.org/download/ESP32_GENERIC/"; exit 1; }
  uv run esptool.py --chip esp32 --port "$PORT" --baud "$BAUD" write_flash -z 0x1000 "$FW"
}
upload() {
  for f in src/*.py; do
    uv run mpremote connect "$PORT" cp "$f" :
  done
  if [ -d src/lib ] && [ -n "$(ls -A src/lib 2>/dev/null)" ]; then
    uv run mpremote connect "$PORT" cp -r src/lib :
  fi
}
reset()   { uv run mpremote connect "$PORT" reset; sleep 2; }
repl()    { uv run mpremote connect "$PORT" resume repl; }

case "${1:-dev}" in
  erase)    erase ;;
  firmware) firmware ;;
  upload)   upload ;;
  reset)    reset ;;
  repl)     repl ;;
  dev)      upload; reset; repl ;;
  all)      erase; firmware; upload; reset; repl ;;
  *)        echo "usage: $0 [erase|firmware|upload|reset|repl|dev|all]"; exit 1 ;;
esac
