package main

import (
	"bufio"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.bug.st/serial"
)

var pipeLen = flag.Int("pipe-len", 100, "Length of sample buffer") // ~0.5s

var min = flag.Float64("min", 1000, "Sample range minimum")

var maxWindow = flag.Int("max-window", 1000, "Maximum window size") // ~10s
var maxMin = flag.Float64("max-min", 2000, "Maximum range minimum")
var maxMax = flag.Float64("max-max", 20_000, "Maximum range maximum")

var debug = flag.Bool("debug", false, "Debug mode")

func main() {
	// parse flags
	flag.Parse()

	// prepare pipe
	pipe := make(chan sample, *pipeLen)

	// read values
	go read(pipe)

	// process values
	go process(pipe)

	// block
	select {}
}

func read(pipe chan<- sample) {
	for {
		// get list
		list, err := serial.GetPortsList()
		if err != nil {
			println(err.Error())
			time.Sleep(time.Second)
			continue
		}

		// check port
		var port string
		for _, name := range list {
			if strings.Contains(name, "usbmodem") {
				port = name
			}
		}
		if port == "" {
			println("no device")
			time.Sleep(time.Second)
			continue
		}

		// open device
		device, err := serial.Open(port, &serial.Mode{
			BaudRate: 115200,
			DataBits: 7,
			Parity:   serial.NoParity,
			StopBits: serial.OneStopBit,
		})
		if err != nil {
			println(err.Error())
			time.Sleep(time.Second)
			continue
		}

		// prepare reader
		reader := bufio.NewReader(device)

		// read data
		for {
			// read line
			line, err := reader.ReadString('\n')
			if err != nil {
				println(err.Error())
				_ = device.Close()
				continue
			}

			// split
			parts := strings.Split(line, ",")

			// decode sample
			sample := make(sample, 0, len(parts))
			for _, seg := range parts {
				value, _ := strconv.ParseFloat(seg, 64)
				sample = append(sample, value)
			}

			// send or drop sample
			select {
			case pipe <- sample:
			default:
			}
		}
	}
}

func process(pipe <-chan sample) {
	// prepare window
	wMax := newWindow(*maxWindow)

	// process values
	for values := range pipe {
		// add min and max
		_, vMax := values.minMax()
		wMax.add(vMax)

		// get max
		_, max := wMax.minMax()

		// adjust max
		max = clamp(max/2, *maxMin, *maxMax)

		// scale
		scaled := make(sample, len(values))
		for i, v := range values {
			scaled[i] = clamp(scale(v, *min, max, 0, 1), 0, 1)
		}

		// debug
		if *debug {
			fmt.Printf("Values: %s | Max: %.0f\n", scaled.String(), max)
		}
	}
}
