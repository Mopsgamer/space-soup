package math

import "math"

func RadiansFromDegrees(deg float64) float64 {
	return deg * (math.Pi / 180.0)
}

func DegreesFromRadians(rad float64) float64 {
	return rad * (180.0 / math.Pi)
}

func DegreesRich(deg, min, sec float64) float64 {
	return deg + min/60.0 + sec/3600.0
}

func RadiansFromRich(deg, min, sec float64) float64 {
	return RadiansFromDegrees(DegreesRich(deg, min, sec))
}
