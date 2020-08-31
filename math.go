package main

import (
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
