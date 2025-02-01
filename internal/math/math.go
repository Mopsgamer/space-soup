package math

import (
	"math"
)

var (
	_1_94787   float64 = 1.94787
	_0_9252    float64 = 0.9252
	_0_4749    float64 = 0.4749
	_10_pow_n3         = 1e-3 // 10^-3
	_10_pow_n5         = 1e-5 // 10^-5
	_0_65              = 0.65
	_0_0398            = 0.0398
	_123_2             = 123.2

	_4_00_00 = RadiansFromDegrees(4)

	_phi       = RadiansFromPrecDegrees(49, 24, 50)
	_phi1      = RadiansFromPrecDegrees(34, 10, 16)  // 34°10'16''
	_phi2      = RadiansFromPrecDegrees(110, 16, 22) // 110°16'22''
	_20_16_22  = RadiansFromPrecDegrees(20, 16, 22)  // 20°16'22''
	_124_10_16 = RadiansFromPrecDegrees(124, 10, 16) // 124°10'16''
	_200_16_22 = RadiansFromPrecDegrees(200, 16, 22) // 200°16'22''
	_304_10_16 = RadiansFromPrecDegrees(304, 10, 16) // 304°10'16''
)

func azimuthInRange(A, more_than, less_or_eq float64) bool {
	return A > more_than && A <= less_or_eq
}

// step 1
func Azimuth(t1, t2 float64) (A float64, err error) {
	K := _1_94787*(t1/t2) + _10_pow_n5

	A_gl := math.Atan((math.Cos(_phi1) - K*math.Cos(_phi2)) / (K*math.Sin(_phi1) - math.Sin(_phi2)))

	if t1 > 0 && t2 >= 0 {
		var add float64 = math.Pi
		if A < 0 {
			add += math.Pi
		}

		A = A_gl + add

		if azimuthInRange(A, _200_16_22, _304_10_16) {
			return A, nil
		}
	} else if t1 <= 0 && t2 < 0 {
		var add float64 = 0
		if A < 0 {
			add += math.Pi
		}

		A = A_gl + add

		if azimuthInRange(A, _20_16_22, _124_10_16) {
			return A, nil
		}
	} else if t1 < 0 && t2 >= 0 { // FIXME: unsure
		A = A_gl
		if azimuthInRange(A_gl, 0, _20_16_22) {
			return A, nil
		}

		A = A_gl + math.Pi
		if A >= _304_10_16 {
			return A, nil
		}

	} else if t1 > 0 && t2 <= 0 {
		A = A_gl + math.Pi
		if azimuthInRange(A, _124_10_16, _200_16_22) {
			return A, nil
		}
	}

	return 0, ErrorSign1
}

// step 2
func ZenithAngle(t1, t2, A, v_avg float64) (z_fix float64, err error) {
	sin_z1 := -((_0_9252 * _10_pow_n3 * v_avg * t1) / math.Cos(A-_phi1))
	sin_z2 := -((_0_4749 * _10_pow_n3 * v_avg * t2) / math.Cos(A-_phi2))
	z1 := math.Asin(sin_z1)
	z2 := math.Asin(sin_z2)
	if z1 < 0 || z2 < 0 {
		return 0, ErrorSign2
	}

	if z1-z2 >= _4_00_00 {
		return 0, ErrorSign3
	}

	// step 3
	delta_v := _0_65 + _0_0398*v_avg // FIXME: 0.0398 or 1.0398?
	v0 := v_avg + delta_v
	z_avg := (z1 + z2) / 2
	if t1 == 0 {
		z_avg = z2
	} else if t2 == 0 {
		z_avg = z1
	}

	// step 4
	v_deriv := math.Sqrt(math.Pow(v0, 2) - _123_2)
	tan_deltaz_divide_2 := math.Abs((v_deriv-v0)/(v_deriv+v0)) * math.Tan(z_avg/2)
	delta_z := math.Atan(tan_deltaz_divide_2) * 2
	z_fix = z_avg + delta_z
	return z_fix, nil
}

// step 5
func RadiantDeclination(z_fix, A float64) (delta float64) {
	sin_delta := math.Sin(_phi)*math.Cos(z_fix) - math.Cos(_phi)*math.Sin(z_fix)*math.Cos(A)
	delta = math.Asin(sin_delta)
	return delta
}

// step 6
func RadiantClockAngle(A, z_fix, delta float64) (sin_t, cos_t, t float64) {
	sin_t = (math.Sin(z_fix) * math.Sin(A)) / math.Cos(delta)
	t = math.Asin(sin_t)
	cos_t = (math.Cos(_phi)*math.Cos(z_fix) + math.Sin(_phi)*math.Sin(z_fix)*math.Cos(A)) / t
	return sin_t, cos_t, t
}

// step 7
func StellarTime(c2, d, h, m float64) (S float64) {
	S = c2 + 0.98565*d + 15.0411*h + 0.25068*m
	return S
}

// step 8
func RightAscension(S, t float64) (alpha float64) {
	alpha = S - t
	return alpha
}
