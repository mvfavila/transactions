package util

import "math"

// roundToCents rounds a float64 value to the nearest cent.
func RoundToCents(value float64) float64 {
	return math.Round(value*100) / 100
}
