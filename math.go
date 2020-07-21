package main

import (
	"math"
	"strconv"
)

func scale(v, iMin, iMax, oMin, oMax float64) float64 {
	return (v-iMin)*(oMax-oMin)/(iMax-iMin) + oMin
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	} else if v > max {
		return max
	}

	return v
}

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

func (s sample) String() string {
	b := make([]byte, 0, 128)
	for i, v := range s {
		if i > 0 {
			b = append(b, ' ')
		}

		b = strconv.AppendFloat(b, v, 'f', 2, 64)
	}

	return string(b)
}

type window struct {
	size int
	list []float64
	pos  int
	len  int
}

func newWindow(size int) *window {
	return &window{
		size: size,
		list: make([]float64, size),
	}
}

func (w *window) add(v float64) {
	// set value
	w.list[w.pos] = v

	// increment position
	w.pos++
	if w.pos > len(w.list)-1 {
		w.pos = 0
	}

	// increment length
	if w.len < w.size {
		w.len++
	}
}

func (w *window) minMax() (float64, float64) {
	// check length
	if w.len == 0 {
		return 0, 0
	}

	// calculate
	min := w.list[0]
	max := w.list[0]
	for i := 1; i < w.len; i++ {
		min = math.Min(min, w.list[i])
		max = math.Max(max, w.list[i])
	}

	return min, max
}
