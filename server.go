package main

import (
	"bufio"
	"strconv"
	"strings"
	"time"

	"github.com/kr/pretty"
	"go.bug.st/serial"
)

func main() {
	// prepare pipe
	pipe := make(chan []float64, 100)

	// read values
	go read(pipe)

	// process values
	go process(pipe)

	// block
	select {}
}

func read(pipe chan<- []float64) {
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

			// decode values
			values := make([]float64, 0, len(parts))
			for _, seg := range parts {
				value, _ := strconv.ParseFloat(seg, 64)
				values = append(values, value)
			}

			// send or drop
			select {
			case pipe <- values:
			default:
			}
		}
	}
}

func process(pipe <-chan []float64) {
	for values := range pipe {
		pretty.Println(values)
	}
}
