package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/VividCortex/gohistogram"
)

var addr = flag.String("addr", "0.0.0.0:8080", "WebSocket server address")

var minRange = flag.Float64("min-range", 1000, "The minimum range")

var threshold = flag.Float64("threshold", 0, "The threshold for on/off values")

var freq = flag.Int("freq", 60, "Sample publish frequency")

var debug = flag.String("debug", "", `Debug server address e.g. ":1234"`)

var usbPrefix = flag.String("usb-prefix", "cu.usbmodem", "The prefix of usb devices to consider")

func main() {
	// parse flags
	flag.Parse()

	// run debug
	if *debug != "" {
		go func() {
			panic(http.ListenAndServe(*debug, uiHandler()))
		}()
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
		minimums := make(sample, 12)
		maximums := make(sample, 12)

		// calculate result
		for i, v := range s {
			// add to histogram
			histogram[i].Add(v)
			histogram[i].Add(v)

			// get quantiles
			min := histogram[i].Quantile(0.15)
			max := histogram[i].Quantile(0.95)

			// adjust range
			if max-min < +*minRange {
				max = min + *minRange
			}

			// get value
			v := clamp(scale(v, min, max, 0, 1), 0, 1)

			// apply threshold if set
			if *threshold > 0 {
				if v > *threshold {
					v = 1
				} else {
					v = 0
				}
			}

			// set result
			result[i] = v
			minimums[i] = min
			maximums[i] = max
		}

		// emit
		emit(result)

		// get min and max
		min, _ := minimums.minMax()
		_, max := maximums.minMax()

		// debug
		if *debug != "" {
			fmt.Printf("Values: %s | Range %.2f - %.2f |  Devices: %d | Clients: %d\n", result.String(), min, max, numDevices, numClients)
		}

		// sleep
		time.Sleep(timeout)
	}
}
