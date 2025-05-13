package soup

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"

	expandrange "github.com/n0madic/expand-range"
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

// 50 items per cluster.
func Paginate[S any](slice []S) (paginated [][]S) {
	paginated = [][]S{}
	for v := range slices.Chunk(slice, 50) {
		paginated = append(paginated, v)
	}
	return
}

// Create identificator for file content.
func HashString(data []byte) string {
	hashBytes := sha256.Sum256(data)
	return hex.EncodeToString(hashBytes[:])
}

func Range(tests []MovementTest, bounds string) (testsRanged []MovementTest, err error) {
	testsRanged = []MovementTest{}
	err = nil
	if bounds == "" {
		testsRanged = tests
		return
	}

	rangeList, err := expandrange.Parse(bounds)
	if err != nil {
		return
	}

	for _, i := range rangeList {
		testsRanged = append(testsRanged, tests[i])
	}

	return
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
	return time.Parse("2006-01-02T15:04", date)
}

func Float64(str string) float64 {
	str = strings.ReplaceAll(str, ",", ".")
	result, _ := strconv.ParseFloat(str, 64)
	return result
}

func Float64Err(str string, error *error) float64 {
	str = strings.ReplaceAll(str, ",", ".")
	result, err := strconv.ParseFloat(str, 64)
	if err != nil {
		*error = errors.Join(*error, err)
	}
	return result
}

func Int(str string) int {
	result, _ := strconv.ParseInt(str, 0, strconv.IntSize)
	return int(result)
}

type TestResult uint

const (
	TestResultAbsoluteSuccess TestResult = iota
	TestResultAcceptable
	TestResultFailed
)

func InDelta(expected, actual, delta float64) TestResult {
	if math.IsNaN(expected) || math.IsNaN(actual) {
		return TestResultFailed
	}
	if expected == actual {
		return TestResultAcceptable
	}
	dt := expected - actual
	if dt >= -delta && dt <= delta {
		return TestResultAcceptable
	}

	return TestResultFailed
}
