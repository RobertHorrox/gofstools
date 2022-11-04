package utils

import "math"

func IntMax(elem ...int) int {
	max := math.MinInt
	for _, i := range elem {
		if i > max {
			max = i
		}
	}

	return max
}
