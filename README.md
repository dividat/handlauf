# Handlauf

**PCB Design, Firmware and Driver vor the Dividat Handlauf Extension**

## PCB Design

See the `Handlauf.sch` and `Handlauf.brd` EAGLE files. The CAM output for
production can be found in the `cam.zip` file.

## Firmware

The `Handlauf.ino` is the Arduino based firmware for the built-in Arduino.
Download the "capacitive sensor" and "ewma" library from the Arduino library
manager.

## Driver

The `*.go` files make up the Go based driver for the sensor. It can be built
using the standard Go mechanisms. The binary provides the following options:

```
Usage of ./handlauf:
  -addr string
        WebSocket server address (default "0.0.0.0:8080")
  -debug string
        Debug server address e.g. ":1234"
  -freq int
        Sample publish frequency (default 60)
  -min-range float
        The minimum range (default 1000)
  -threshold float
        The threshold for on/off values
```

The driver contains a Web based debug interface to inspect the sensor values.
If the driver is started with the `-debug :1234` flag, the interface is
available at `http://0.0.0.0:1234`.
