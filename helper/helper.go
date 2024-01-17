package helper

import "math"

func AlmostEqual(a, b float64, epsilon float64) bool {
    if a == b {
        return true
    }
    diff := math.Abs(a - b)
    return diff < epsilon
}
