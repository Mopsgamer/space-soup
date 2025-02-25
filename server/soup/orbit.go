package soup

import (
	"math"
	"time"
)

var (
	Pi0                = 1.7864122
	_0_4703            = 0.47033133
	_0_65              = 0.65
	_1_0398            = 1.0398 // FIXME: 0.0398 or 1.0398?
	_123_2             = 123.2
	_0_01672           = 0.01672
	e                  = 0.40918274
	sin_e, cos_e       = math.Sincos(e)
	e0                 = 0.01675
	T                  = 2 * 1e-3 // FIXME: T = 2 *10^-3 c. what is c
	l1                 = 4324.
	l2                 = 8422.
	m                  = l1 / l2
	phi                = RadiansFromRich(49, 24, 50)
	sin_phi, cos_phi   = math.Sincos(phi)
	phi1               = RadiansFromRich(34, 10, 16) // 34°10'16''
	sin_phi1, cos_phi1 = math.Sincos(phi1)
	phi2               = RadiansFromRich(110, 16, 22) // 110°16'22''
	sin_phi2, cos_phi2 = math.Sincos(phi2)
)

type MeteoroidMovement struct {
	K    float64
	A_gl float64
	// Азимут
	A float64
	// Зенитный угол радианта
	Z_avg float64
	// Зенитное расстояние радианта
	Z_fix float64
	// Склонение радианта
	Delta float64
	Sin_t float64
	Cos_t float64
	// Часовой угол
	T float64
	// Звездное время в момент наблюдения
	S float64
	// Исправленные экваториальные координаты радианта
	Alpha_fix float64
	// Исправленные экваториальные координаты радианта
	Delta_fix float64
	// Поправки за суточную аберрацию в экваториальных координатах
	Delta_alpha float64
	// Поправки за суточную аберрацию в экваториальных координатах
	Delta_delta float64
	// Геоцентрическая скорость
	V_g float64
	// Внеатмосферная скорость
	V_inf float64
}

type MeteoroidMovementInput struct {
	Dist int
	// Временная задержка
	Tau1 float64
	// Временная задержка
	Tau2  float64
	V_avg float64
	// Время и дата появления метеороида
	Date time.Time
}

func NewMeteoroidMovement(inp MeteoroidMovementInput) *MeteoroidMovement {
	var temp float64
	// step 1

	k := m * (inp.Tau1 / inp.Tau2)

	// Главное значение азимута
	A_gl := math.Atan((cos_phi1 - k*cos_phi2) / (k*sin_phi1 - sin_phi2))

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
	sin_A, cos_A := math.Sincos(A)

	// step 2

	temp = 2 * inp.V_avg * T
	cos_A_minus_phi1 := math.Cos(A - phi1)
	cos_A_minus_phi2 := math.Cos(A - phi2)
	sin_z1 := (temp * inp.Tau1) / (l1 * cos_A_minus_phi1)
	sin_z2 := (temp * inp.Tau2) / (l2 * cos_A_minus_phi2)

	W1 := math.Abs(cos_A_minus_phi1)
	W2 := math.Abs(cos_A_minus_phi2)
	// Зенитный угол радианта
	Z_avg := math.Asin((W1*sin_z1 + W2*sin_z2) / (W1 + W2))

	// step 3

	// Скорость V0 с учетом поправки за торможение
	V0 := _1_0398*inp.V_avg + _0_65

	// step 4

	// Скорость
	V_deriv := math.Sqrt(math.Pow(V0, 2) - _123_2)
	delta_Z := 2 * math.Atan(math.Abs((V_deriv-V0))/(V_deriv+V0)*math.Tan(Z_avg/inp.V_avg))
	// Зенитное расстояние радианта
	Z_fix := Z_avg + delta_Z
	sin_Z_fix, cos_Z_fix := math.Sincos(Z_fix)
	sin_Z_fix_cos_A := sin_Z_fix * cos_A
	sin_Z_fix_sin_A := sin_Z_fix * sin_A

	// step 5

	// Склонение радианта
	delta := math.Asin(sin_phi*cos_Z_fix - cos_phi*sin_Z_fix_cos_A)
	sin_delta, cos_delta := math.Sincos(delta)

	// step 6

	t_gl := math.Atan((sin_Z_fix_sin_A) / (cos_phi*cos_Z_fix + sin_phi*sin_Z_fix_cos_A))
	// Часовой угол радианта
	t := 0.
	temp = sin_Z_fix_sin_A
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
	sin_t, cos_t := math.Sincos(t)

	// step 7

	// Константа, определяемая для каждого года
	c2 := inp.Dist // FIXME: c2

	// Звездное время в момент наблюдения
	S := StellarTime(c2, inp.Date.YearDay()-1, inp.Date.Hour(), inp.Date.Minute())

	// step 8

	// Прямое восхождение радианта
	alpha := S - t // 0 < alpha < 2 pi (6.28)

	// step 9

	// Поправки за суточную аберрацию в экваториальных координатах
	delta_alpha := -((_0_4703 * cos_t * cos_phi) / (V_deriv * cos_delta))
	// Поправки за суточную аберрацию в экваториальных координатах
	delta_delta := -((_0_4703 * sin_t * sin_delta * cos_phi) / (V_deriv))

	// step 10

	// Исправленные экваториальные координаты радианта
	alpha_fix := alpha + delta_alpha // while alpha_fix > 0
	sin_alpha_fix, cos_alpha_fix := math.Sincos(alpha_fix)
	// Исправленные экваториальные координаты радианта
	delta_fix := delta + delta_delta
	sin_delta_fix, cos_delta_fix := math.Sincos(delta_fix)

	// step 11

	psi_E_gl := math.Acos(-sin_t * cos_delta_fix)
	// Элонгация
	psi_E := psi_E_gl
	if psi_E_gl < 0 {
		psi_E += math.Pi
	}

	// step 12

	delta_s := math.Sqrt(math.Pow(delta_alpha*cos_delta_fix, 2) + math.Pow(delta_delta, 2))
	// Геоцентрическая скорость
	V_g := V_deriv * (math.Sin(psi_E-delta_s) / math.Sin(psi_E))
	pow_2_V_g := math.Pow(V_g, 2)

	// step 13

	// Внеатмосферная скорость
	V_inf := math.Sqrt(pow_2_V_g + _123_2)

	// step 14

	// Эклиптическая широта радианта
	beta := math.Asin(-sin_e*sin_alpha_fix*cos_delta_fix + cos_e + sin_delta_fix)
	sin_beta, cos_beta := math.Sincos(beta)

	// step 15

	cos_lambda := ((cos_delta_fix * cos_alpha_fix) / cos_beta)
	sin_lambda := (1 / cos_beta) * (cos_delta_fix*sin_alpha_fix*cos_e + sin_delta_fix*sin_e)
	// Эклиптическая долгота радианта
	lambda := math.Atan2(sin_lambda, cos_lambda)

	// step 16

	// Константа, определяемая для каждого года
	c3 := 0 // FIXME: c3
	// Долгота Солнца
	lambda_theta := SolarLongitude(c3, inp.Date.YearDay()-1, inp.Date.Hour(), inp.Date.Minute())

	// step 17

	delta_theta := _0_01672 * math.Sin(lambda_theta-Pi0)
	lambda_apex := lambda_theta + delta_theta - (math.Pi / 2)
	// Долгота радианта относительно апекса
	diff_lambda := lambda - lambda_apex
	_ = diff_lambda

	// step 18

	// Угол элонгации видимого радианта от апекса движения Земли
	E_apex := math.Acos(cos_beta * math.Cos(lambda-lambda_apex))

	// step 19

	// Радиус вектор орбиты Земли
	R := (1 - math.Pow(e0, 2)) / (1 - e0*math.Cos(lambda_theta-Pi0))

	// step 20

	// Орбитальная скорость Земли для данного дня
	V_t := 29.76 * math.Sqrt((2/R)-1)

	// step 21

	temp = lambda_theta + delta_theta - lambda
	sin_temp, cos_temp := math.Sincos(temp)
	temp_deriv_gl := math.Atan2(sin_temp-(V_t/(V_g*cos_beta)), cos_temp)
	temp_deriv := temp_deriv_gl
	if cos_temp < 0 {
		temp_deriv += math.Pi
	}
	// Долгота истинного радианта
	lambda_deriv := lambda_theta + delta_theta - temp_deriv // 0 <= lambda_deriv <= 2*pi
	sin_lambda_diff, cos_lambda_diff := math.Sincos(lambda_theta - lambda_deriv)

	// step 22

	// Гелиоцентрическая скорость
	V_h := math.Sqrt(pow_2_V_g + math.Pow(V_t, 2) - 2*V_g*V_t*math.Cos(E_apex))

	// step 23

	// Широта истинного радианта
	beta_deriv := math.Asin((V_g / V_h) * sin_beta)
	cos_beta_deriv := math.Cos(beta_deriv)

	// step 24

	// Элонгация истинного радианта от апекса
	E_deriv := math.Acos(cos_beta_deriv * math.Cos(lambda_deriv-lambda_apex))
	_ = E_deriv

	// step 25

	i_gl := math.Atan(-(math.Abs(math.Tan(beta_deriv)) / sin_lambda_diff))
	// Наклонение орбиты частицы к плоскости эклиптики
	i := i_gl
	if i_gl > 0 {
		i += math.Pi
	}

	// step 26

	// Афелий – точка орбиты максимально удаленная от Солнца
	Q := math.Pow(V_h/V_t, 2)

	// step 27

	// Большая полуось
	a := 1 / ((2 - Q) / R)

	// step 28

	// Угол, образуемый радиус-вектором метеорного тела с вектором его скорости
	psi := math.Acos(-cos_beta_deriv * cos_lambda_diff)
	// Элонгация радианта от Солнца
	E_theta_deriv := math.Pi - psi

	// step 29

	// Параметр орбиты
	p := math.Pow(R, 2) * Q * math.Pow(math.Sin(E_theta_deriv), 2)
	b := math.Sqrt(p * math.Abs(a))

	// step 30

	var c float64
	if a > 0 {
		c = math.Sqrt(math.Pow(a, 2) - math.Pow(b, 2))
	} else {
		c = math.Sqrt(math.Pow(a, 2) + math.Pow(b, 2))
	}
	// Эксцентриситет
	e = math.Abs(c / a)

	// step 31

	// Перигелийное расстояние
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
	// Истинная аномалия
	v := math.Atan2(sin_v, cos_v)

	// step 34

	// Аргумент перигелия
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
		V_g:         V_g,
		V_inf:       V_inf,
	}
}

// Format: 2006-01-02T03:04
func ParseDate(date string) (time.Time, error) {
	return time.Parse("2006-01-02T03:04", date)
}

func StellarTime(c2, d, h, m int) (S float64) {
	S = float64(c2) + 0.98565*float64(d) + 15.0411*float64(h) + 0.25068*float64(m)
	return
}

func SolarLongitude(c3, d, h, m int) (lambda_theta float64) {
	lambda_theta = -float64(c3) + 0.0000097*float64(m) + 0.000717*float64(h) + 0.017203*float64(d) + 0.034435*math.Sin(0.017203*float64(d-2))
	return
}
