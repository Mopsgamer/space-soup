package soup

import (
	"math"
	"time"
)

var (
	_1_94787   float64 = 1.94787
	_0_9252            = 0.9252
	_0_4749            = 0.4749
	_10_pow_n3         = 1e-3 // 10^-3
	_10_pow_n5         = 1e-5 // 10^-5
	_0_65              = 0.65
	_0_0398            = 0.0398
	_123_2             = 123.2
	_4_00_00           = RadiansFromDegrees(4)
	_phi               = RadiansFromRich(49, 24, 50)
	_phi1              = RadiansFromRich(34, 10, 16)  // 34°10'16''
	_phi2              = RadiansFromRich(110, 16, 22) // 110°16'22''
	_26_948            = RadiansFromDegrees(26.948)   // 26,948°

	// TODO: remove
	// _20_16_22          = RadiansFromRich(20, 16, 22)  // 20°16'22''
	// _124_10_16         = RadiansFromRich(124, 10, 16) // 124°10'16''
	// _200_16_22         = RadiansFromRich(200, 16, 22) // 200°16'22''
	// _304_10_16         = RadiansFromRich(304, 10, 16) // 304°10'16''
)

// TODO: remove
// func azimuthInRange(A, more_than, less_or_eq float64) bool {
// 	return A > more_than && A <= less_or_eq
// }

// step 1
func Azimuth(tau1, tau2 float64) (K, A_gl, A float64) {
	K = _1_94787*(tau1/tau2) + _10_pow_n5

	A_gl = math.Atan((math.Cos(_phi1) - K*math.Cos(_phi2)) / (K*math.Sin(_phi1) - math.Sin(_phi2)))

	if tau1 <= 0 {
		if tau2 < 0 {
			if A_gl >= 0 {
				A = A_gl
			} else {
				A = A_gl + math.Pi
			}
		} else {
			if A_gl >= 0 {
				A = A_gl
			} else {
				A = A_gl + 2*math.Pi
			}
		}
	} else { // tau1 > 0
		if tau2 < 0 {
			A = A_gl + math.Pi
		} else {
			if A_gl >= 0 {
				A = A_gl + math.Pi
			} else {
				A = A_gl + 2*math.Pi
			}
		}
	}

	return K, A_gl, A
}

// step 2
func ZenithAngle(t1, t2, A, v_avg float64) (z_avg, z_fix, v_deriv float64, err error) {
	sin_z1 := -((_0_9252 * _10_pow_n3 * v_avg * t1) / math.Cos(A-_phi1))
	sin_z2 := -((_0_4749 * _10_pow_n3 * v_avg * t2) / math.Cos(A-_phi2))
	z1 := math.Asin(sin_z1)
	z2 := math.Asin(sin_z2)
	if z1 < 0 || z2 < 0 {
		return 0, 0, 0, ErrorSign2
	}

	if z1-z2 >= _4_00_00 {
		return 0, 0, 0, ErrorSign3
	}

	// step 3
	delta_v := _0_65 + _0_0398*v_avg // FIXME: 0.0398 or 1.0398?
	v0 := v_avg + delta_v
	z_avg = (z1 + z2) / 2
	if t1 == 0 {
		z_avg = z2
	} else if t2 == 0 {
		z_avg = z1
	}

	// step 4
	v_deriv = math.Sqrt(math.Pow(v0, 2) - _123_2)
	tan_deltaz_divide_2 := math.Abs((v_deriv-v0)/(v_deriv+v0)) * math.Tan(z_avg/2)
	delta_z := math.Atan(tan_deltaz_divide_2) * 2
	z_fix = z_avg + delta_z
	return z_avg, z_fix, v_deriv, nil
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
	// TODO: maybe cos_t replacable with math.Cos(t)
	cos_t = (math.Cos(_phi)*math.Cos(z_fix) + math.Sin(_phi)*math.Sin(z_fix)*math.Cos(A)) / t
	return sin_t, cos_t, t
}

// step 7
func StellarTime(c2 float64, d, h, m int) (S float64) {
	S = c2 + 0.98565*float64(d) + 15.0411*float64(h) + 0.25068*float64(m)
	return S
}

// step 8
func RightAscension(S, t float64) (alpha float64) {
	alpha = S - t // 0 < alpha < 2 pi (6.28)
	return alpha
}

// step 9
func FixDiaurnalAberration(sin_t, cos_t, delta, v_deriv float64) (delta_alpha, delta_delta float64) {
	delta_alpha = -(_26_948 / v_deriv) * (cos_t / math.Cos(delta)) * math.Cos(_phi)
	delta_delta = -(_26_948 / v_deriv) * sin_t * math.Sin(delta) * math.Cos(_phi)
	return delta_alpha, delta_delta
}

// step 10
func PrepareSpeed(alpha, delta_alpha, delta, delta_delta, t, v_deriv float64) (v_geoc, v_vacuum float64) {
	alpha_fix := alpha + delta_alpha // while alpha_fix > 0
	delta_fix := delta + delta_delta

	_ = alpha_fix // TODO: remove

	// step 11
	cos_psi_gl := -math.Sin(t) * math.Cos(delta_fix)
	psi_gl := math.Acos(cos_psi_gl)
	add := 0.
	if psi_gl > 0 {
		add = math.Pi
	}
	psi := psi_gl + add

	// step 12
	delta_S := math.Sqrt(math.Pow(delta_alpha*math.Cos(delta_fix), 2) + math.Pow(delta_delta, 2))
	v_geoc = (v_deriv * math.Sin(psi-delta_S)) / math.Sin(psi)
	// step 13
	v_vacuum = math.Sqrt(math.Pow(v_geoc, 2) + _123_2)
	return v_geoc, v_vacuum
}

type MeteoroidMovement struct {
	K           float64
	A_gl        float64
	A           float64
	Z_avg       float64
	Z_fix       float64
	Delta       float64
	Sin_t       float64
	Cos_t       float64
	T           float64
	S           float64
	Delta_alpha float64
	Delta_delta float64
	V_geoc      float64
	V_vacuum    float64
}

func ParseDate(date string) (time.Time, error) {
	return time.Parse("2006-01-02T03:04", date)
}

type MeteoroidMovementInput struct {
	Dist  int32
	Tau1  float64
	Tau2  float64
	V_avg float64
	Date  time.Time
}

func NewMeteoroidMovement(inp MeteoroidMovementInput) (*MeteoroidMovement, error) {
	K, A_gl, A := Azimuth(inp.Tau1, inp.Tau2)

	z_avg, z_fix, v_deriv, err := ZenithAngle(inp.Tau1, inp.Tau2, A, inp.V_avg)
	if err != nil {
		return nil, err
	}

	delta := RadiantDeclination(z_fix, A)
	sin_t, cos_t, t := RadiantClockAngle(A, z_fix, delta)

	c2 := 1. // FIXME: unsure C2
	S := StellarTime(c2, inp.Date.YearDay()-1, inp.Date.Hour(), inp.Date.Minute())
	alpha := RightAscension(S, t)
	delta_alpha, delta_delta := FixDiaurnalAberration(sin_t, cos_t, delta, v_deriv)
	v_geoc, v_vacuum := PrepareSpeed(alpha, delta_alpha, delta, delta_delta, t, v_deriv)

	return &MeteoroidMovement{
		K:           K,
		A_gl:        A_gl,
		A:           A,
		Z_avg:       z_avg,
		Z_fix:       z_fix,
		Delta:       delta,
		Sin_t:       sin_t,
		Cos_t:       cos_t,
		T:           t,
		Delta_alpha: delta_alpha,
		Delta_delta: delta_delta,
		S:           S,
		V_geoc:      v_geoc,
		V_vacuum:    v_vacuum,
	}, nil
}
