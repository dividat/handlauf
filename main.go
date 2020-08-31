package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/256dpi/god"
	"github.com/VividCortex/gohistogram"
)

var addr = flag.String("addr", "0.0.0.0:8080", "WebSocket server address")

var minRange = flag.Float64("min-range", 1000, "The minimum range")

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

	// prepare histograms
	var histogram []gohistogram.Histogram
	for i := 0; i < 12; i++ {
		histogram = append(histogram, gohistogram.NewHistogram(80))
	}

	for {
		// get samples
		l := left.Load().(sample)
		r := right.Load().(sample)

		// merge samples
		s := make(sample, 12)
		copy(s, l)
		copy(s[6:], r)

		// prepare result
		result := make(sample, 12)

		// calculate result
		for i, v := range s {
			// add to histogram
			histogram[i].Add(v)
			histogram[i].Add(v)

			// get quantiles
			min := histogram[i].Quantile(0.2)
			max := histogram[i].Quantile(0.8)

			// adjust max
			if max < min+*minRange {
				max = min + *minRange
			}

			// get value
			v := clamp(scale(v, min, max, 0, 1), 0, 1)

			// set result
			result[i] = v
		}

		// emit
		emit(result)

		// debug
		if *debug {
			fmt.Printf("Values: %s | Devices: %d | Clients: %d\n", result.String(), numDevices, numClients)
		}

		// sleep
		time.Sleep(timeout)
	}
}
