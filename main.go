package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/256dpi/god"
)

var addr = flag.String("addr", "0.0.0.0:8080", "WebSocket server address")

var min = flag.Float64("min", 1000, "Sample range minimum")

var maxWindow = flag.Int("max-window", 1000, "Maximum window size") // ~10s
var maxMin = flag.Float64("max-min", 2000, "Maximum range minimum")
var maxMax = flag.Float64("max-max", 20000, "Maximum range maximum")

var freq = flag.Int("freq", 60, "Sample publish frequency")

var debug = flag.Bool("debug", false, "Debug mode")

func main() {
	// parse flags
	flag.Parse()

	// run debug
	if *debug {
		god.Init(god.Options{})
	}

	// open stream
	go stream(*addr)

	// manage devices
	go manage()

	// process samples
	process()
}

func process() {
	// timeout
	timeout := time.Second / time.Duration(*freq)

	// prepare windows
	lwMax := newWindow(*maxWindow)
	rwMax := newWindow(*maxWindow)

	for {
		// get samples
		l := left.Load().(sample)
		r := right.Load().(sample)

		// add min and max
		_, lvMax := l.minMax()
		_, rvMax := r.minMax()
		lwMax.add(lvMax)
		rwMax.add(rvMax)

		// get max
		_, lMax := lwMax.minMax()
		_, rMax := rwMax.minMax()

		// adjust max
		lMax = clamp(lMax/2, *maxMin, *maxMax)
		rMax = clamp(rMax/2, *maxMin, *maxMax)

		// get result
		result := make(sample, 12)
		for i, v := range l {
			result[i] = clamp(scale(v, *min, lMax, 0, 1), 0, 1)
		}
		for i, v := range r {
			result[i+6] = clamp(scale(v, *min, rMax, 0, 1), 0, 1)
		}

		// emit
		emit(result)

		// debug
		if *debug {
			fmt.Printf("Values: %s | Max: %.0f %.0f | Devices: %d | Clients: %d\n", result.String(), lMax, rMax, numDevices, numClients)
		}

		// sleep
		time.Sleep(timeout)
	}
}
