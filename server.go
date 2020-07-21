package main

import (
	"bufio"
	"strconv"
	"strings"
	"time"

	"github.com/kr/pretty"
	"go.bug.st/serial"
)

type sample []float64

func (s sample) minMax() (float64, float64) {
	// find min and max
	var max = s[0]
	var min = s[0]
	for _, value := range s {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}

	return min, max
}

func main() {
	// prepare pipe
	pipe := make(chan sample, 100)

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
	t := time.Now()
	for values := range pipe {
		pretty.Println(values, time.Since(t).String())
		t = time.Now()
	}
}
