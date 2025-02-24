package soup

import (
	"math"
	"time"
)

var (
	Pi0      = 1.7864122
	_0_4703  = 0.47033133
	_0_65    = 0.65
	_1_0398  = 1.0398 // FIXME: 0.0398 or 1.0398?
	_123_2   = 123.2
	_0_01672 = 0.01672
	e        = 0.40918274
	e0       = 0.01675
	T        = 1e-3
	l1       = 4324.
	l2       = 8422.
	m        = l1 / l2
	phi      = RadiansFromRich(49, 24, 50)
	phi1     = RadiansFromRich(34, 10, 16)  // 34°10'16''
	phi2     = RadiansFromRich(110, 16, 22) // 110°16'22''
)

func StellarTime(c2, d, h, m int) (S float64) {
	S = float64(c2) + 0.98565*float64(d) + 15.0411*float64(h) + 0.25068*float64(m)
	return S
}

func SolarLongitude(c3, d, h, m int) (lambda_theta float64) {
	lambda_theta = -float64(c3) + 0.0000097*float64(m) + 0.000717*float64(h) + 0.017203*float64(d) + 0.034435*math.Sin(0.017203*float64(d-2))
	return lambda_theta
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
	Alpha_fix   float64
	Delta_fix   float64
	Delta_alpha float64
	Delta_delta float64
	V_geoc      float64
	V_vacuum    float64
}

func ParseDate(date string) (time.Time, error) {
	return time.Parse("2006-01-02T03:04", date)
}

type MeteoroidMovementInput struct {
	Dist  int
	Tau1  float64
	Tau2  float64
	V_avg float64
	Date  time.Time
}

func NewMeteoroidMovement(inp MeteoroidMovementInput) (*MeteoroidMovement, error) {
	var temp float64
	// step 1

	k := m * (inp.Tau1 / inp.Tau2)

	// Главное значение азимута
	A_gl := math.Atan((math.Cos(phi1) - k*math.Cos(phi2)) / (k*math.Sin(phi1) - math.Sin(phi2)))

	// Азимут
	var A float64

	if inp.Tau1 <= 0 {
		if inp.Tau2 < 0 {
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
	} else { // inp.Tau1 > 0
		if inp.Tau2 < 0 {
			A = A_gl + math.Pi
		} else {
			if A_gl >= 0 {
				A = A_gl + math.Pi
			} else {
				A = A_gl + 2*math.Pi
			}
		}
	}

	// step 2

	sin_z1 := (2 * inp.V_avg * T * inp.Tau1) / (l1 * math.Cos(A-phi1))
	sin_z2 := (2 * inp.V_avg * T * inp.Tau2) / (l2 * math.Cos(A-phi2))

	W1 := math.Abs(math.Cos(A - phi1))
	W2 := math.Abs(math.Cos(A - phi2))
	// Зенитный угол радианта
	Z_avg := math.Asin((W1*sin_z1 + W2*sin_z2) / (W1 + W2))

	// step 3

	// Скорость V0 с учетом поправки за торможение
	V0 := _1_0398*inp.V_avg + _0_65

	// step 4

	V_deriv := math.Sqrt(math.Pow(V0, 2) - _123_2)
	delta_Z := 2 * math.Atan(math.Abs((V_deriv-V0))/(V_deriv+V0)*math.Tan(Z_avg/inp.V_avg))
	Z_fix := Z_avg + delta_Z

	// step 5

	delta := math.Asin(math.Sin(phi)*math.Cos(Z_fix) - math.Cos(phi)*math.Sin(Z_fix)*math.Cos(A))

	// step 6

	// Часовой угол радианта
	t_gl := math.Atan((math.Sin(Z_fix) * math.Sin(A)) / (math.Cos(phi)*math.Cos(Z_fix) + math.Sin(phi)*math.Sin(Z_fix)*math.Cos(A)))
	t := 0.
	temp = math.Sin(Z_fix) * math.Sin(A)
	if t_gl >= 0 {
		if temp >= 0 {
			t = t_gl
		} else {
			t = t_gl + math.Pi
		}
	} else { // t_gl < 0
		if temp >= 0 {
			t = t_gl + math.Pi
		} else {
			t = t_gl + 2*math.Pi
		}
	}

	// step 7

	c2 := inp.Dist // FIXME: c2
	S := StellarTime(c2, inp.Date.YearDay()-1, inp.Date.Hour(), inp.Date.Minute())

	// step 8

	alpha := S - t // 0 < alpha < 2 pi (6.28)

	// step 9

	delta_alpha := -((_0_4703 * math.Cos(t) * math.Cos(phi)) / (V_deriv * math.Cos(delta)))
	delta_delta := -((_0_4703 * math.Sin(t) * math.Sin(delta) * math.Cos(phi)) / (V_deriv))

	// step 10

	alpha_fix := alpha + delta_alpha // while alpha_fix > 0
	delta_fix := delta + delta_delta

	// step 11

	psi_E_gl := math.Acos(-math.Sin(t) * math.Cos(delta_fix))
	psi_E := psi_E_gl
	if psi_E_gl < 0 {
		psi_E += math.Pi
	}

	// step 12

	delta_s := math.Sqrt(math.Pow(delta_alpha*math.Cos(delta_fix), 2) + math.Pow(delta_delta, 2))
	V_g := V_deriv * (math.Sin(psi_E-delta_s) / math.Sin(psi_E))

	// step 13

	V_inf := math.Sqrt(math.Pow(V_g, 2) + _123_2)

	// step 14

	// Эклиптическая широта радианта
	beta := math.Asin(-math.Sin(e)*math.Sin(alpha_fix)*math.Cos(delta_fix) + math.Cos(e) + math.Sin(delta_fix))

	// step 15

	// Эклиптическая долгота радианта
	cos_lambda := ((math.Cos(delta_fix) * math.Cos(alpha_fix)) / math.Cos(beta))
	sin_lambda := (1 / math.Cos(beta)) * (math.Cos(delta_fix)*math.Sin(alpha_fix)*math.Cos(e) + math.Sin(delta_fix)*math.Sin(e))
	lambda := math.Atan2(sin_lambda, cos_lambda)

	// step 16

	c3 := 0 // FIXME: c3
	lambda_theta := SolarLongitude(c3, inp.Date.YearDay()-1, inp.Date.Hour(), inp.Date.Minute())

	// step 17

	delta_theta := _0_01672 * math.Sin(lambda_theta-Pi0)
	lambda_apex := lambda_theta + delta_theta - (math.Pi / 2)

	// step 18

	E_apex := math.Acos(math.Cos(beta) * math.Cos(lambda-lambda_apex))

	// step 19

	R := (1 - math.Pow(e0, 2)) / (1 - e0*math.Cos(lambda_theta-Pi0))

	// step 20

	// Орбитальная скорость Земли для данного дня
	V_t := 29.76 * math.Sqrt((2/R)-1)

	// step 21

	temp = lambda_theta + delta_theta - lambda
	cos_temp := math.Cos(temp)
	temp_deriv_gl := math.Atan2(math.Sin(temp)-(V_t/(V_g*math.Cos(beta))), cos_temp)
	temp_deriv := temp_deriv_gl
	if cos_temp < 0 {
		temp_deriv += math.Pi
	}

	lambda_deriv := lambda_theta + delta_theta - temp_deriv // 0 <= lambda_deriv <= 2*pi

	// step 22

	V_h := math.Sqrt(math.Pow(V_g, 2) + math.Pow(V_t, 2) - 2*V_g*V_t*math.Cos(E_apex))

	// step 23

	beta_deriv := math.Asin((V_g / V_h) * math.Sin(beta))

	// step 24

	E_deriv := math.Acos(math.Cos(beta_deriv) * math.Cos(lambda_deriv-lambda_apex))
	_ = E_deriv

	// step 25

	i_gl := math.Atan(-(math.Abs(math.Tan(beta_deriv)) / math.Sin(lambda_theta-lambda_deriv)))
	i := i_gl
	if i_gl > 0 {
		i += math.Pi
	}

	// step 26

	Q := math.Pow(V_h/V_t, 2)

	// step 27

	// Величина, обратная большой полуоси
	a := 1 / ((2 - Q) / R)

	// step 28

	psi := math.Acos(-math.Cos(beta_deriv) * math.Cos(lambda_theta-lambda_deriv))
	E_theta_deriv := math.Pi - psi

	// step 29

	p := math.Pow(R, 2) * Q * math.Pow(math.Sin(E_theta_deriv), 2)
	b := math.Sqrt(p * math.Abs(a))

	// step 30

	var c float64
	if a > 0 {
		c = math.Sqrt(math.Pow(a, 2) - math.Pow(b, 2))
	} else {
		c = math.Sqrt(math.Pow(a, 2) + math.Pow(b, 2))
	}
	e = math.Abs(c / a)

	// step 31

	q := a - c
	if a < 0 {
		q = c - math.Abs(a)
	}
	_ = q

	// step 32

	// Долгота восходящего узла
	omega := lambda_theta
	if beta_deriv < 0 {
		omega += math.Pi
	}
	_ = omega

	// step 33

	cos_v := (p - R) / (R * e)
	sin_v := p / (R * e) * math.Cos(i) * Ctg(lambda_deriv-lambda_theta)
	v := math.Atan2(sin_v, cos_v)

	// step 34

	var wmega float64
	if beta_deriv > 0 {
		wmega = math.Pi - v
	} else {
		wmega = -v
	}
	_ = wmega

	return &MeteoroidMovement{
		K:           k,
		A_gl:        A_gl,
		A:           A,
		Z_avg:       Z_avg,
		Z_fix:       Z_fix,
		Delta:       delta,
		T:           t,
		Alpha_fix:   alpha_fix,
		Delta_fix:   delta_fix,
		Delta_alpha: delta_alpha,
		Delta_delta: delta_delta,
		S:           S,
		V_geoc:      V_g,
		V_vacuum:    V_inf,
	}, nil
}
