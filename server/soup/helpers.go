package soup

import (
	"math"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/constraints"
)

func LoopNumber[T constraints.Float | constraints.Integer](value, min, max T) T {
	rangeSize := max - min
	if rangeSize == 0 {
		return min
	}

	modVal := T(math.Mod(float64(value-min), float64(rangeSize)))

	// Ensure the result stays within the range
	if modVal < 0 {
		modVal += rangeSize
	}
	return modVal + min
}

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

func RichFromRadians(radians float64) (degrees, minutes, seconds float64) {
	degreesFloat := radians * 180 / math.Pi
	degrees = math.Floor(degreesFloat)

	minutesFloat := (degreesFloat - degrees) * 60
	minutes = math.Floor(minutesFloat)

	secondsFloat := (minutesFloat - minutes) * 60
	seconds = secondsFloat

	return degrees, minutes, seconds
}

func Ctg(x float64) float64 {
	sin_x, cos_x := math.Sincos(x)
	return cos_x / sin_x
}

func ArcCtg(x float64) float64 {
	return math.Pi/2 - math.Atan(x)
}

func Average(x []float64) (avg float64) {
	for _, v := range x {
		avg += v
	}
	avg /= float64(len(x))
	return
}

// Format: 2006-01-02T03:04
func ParseDateJSON(date string) (time.Time, error) {
	return time.Parse("2006-01-02T03:04", date)
}

func Float64(str string) float64 {
	str = strings.Replace(str, ",", ".", -1)
	result, _ := strconv.ParseFloat(str, 64)
	return result
}

func Int(str string) int {
	result, _ := strconv.ParseInt(str, 0, strconv.IntSize)
	return int(result)
}
