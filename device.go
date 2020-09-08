package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go.bug.st/serial"
)

var devices sync.Map
var numDevices int64

var left atomic.Value
var right atomic.Value

func init() {
	// initialize
	left.Store(make(sample, 6))
	right.Store(make(sample, 6))
}

func manage() {
	for {
		// get list
		list, err := serial.GetPortsList()
		if err != nil {
			println("manage:", err.Error())
			time.Sleep(time.Second)
			continue
		}

		// check devices
		for _, name := range list {
			if strings.Contains(name, *usbPrefix) {
				if _, ok := devices.Load(name); !ok {
					// add device
					devices.Store(name, true)
					atomic.AddInt64(&numDevices, 1)
					fmt.Printf("manage: added device %s\n", name)

					// read
					go read(name)
				}
			}
		}

		// sleep
		time.Sleep(time.Second)
	}
}

func read(name string) {
	// ensure remove
	defer func() {
		devices.Delete(name)
		atomic.AddInt64(&numDevices, -1)
		fmt.Printf("read: removed device %s\n", name)
	}()

	// open device
	device, err := serial.Open(name, &serial.Mode{
		BaudRate: 115200,
		DataBits: 7,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	})
	if err != nil {
		fmt.Printf("read: %s\n", err.Error())
		return
	}

	// ensure close
	defer device.Close()

	// prepare reader
	reader := bufio.NewReader(device)

	// read data
	for {
		// read line
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("read: %s\n", err.Error())
			return
		}

		// remove space
		line = strings.TrimSpace(line)

		// split
		parts := strings.Split(line, ",")

		// get orientation
		orientation := parts[0]
		parts = parts[1:]

		// decode sample
		values := make(sample, 6)
		for i := 0; i < len(values) && i < len(parts); i++ {
			values[i], _ = strconv.ParseFloat(parts[i], 64)
		}

		// set sample
		switch orientation {
		case "L":
			left.Store(values)
		case "R":
			right.Store(values)
		}
	}
}
