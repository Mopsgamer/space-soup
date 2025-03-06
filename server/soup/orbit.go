package soup

import (
	"math"
	"time"
)

var (
	c2 = 1.61222
	c3 = 1.4481

	Pi0 = 1.7864122

	_0_4703  = 0.47033133
	_0_65    = 0.65
	_1_0398  = 1.0398
	_123_2   = 123.2
	_0_01672 = 0.01672

	e    = 0.40918274
	e0   = 0.01675
	T    = 2 * 1e-3
	l1   = 8422.
	l2   = 4324.
	m    = l1 / l2
	phi  = 0.86244  // Широта места наблюдения
	phi1 = 0.596398 // 34°10'16''
	phi2 = 1.92462  // 110°16'22''

	// optimizations
	sin_phi, cos_phi   = math.Sincos(phi)
	sin_phi1, cos_phi1 = math.Sincos(phi1)
	sin_phi2, cos_phi2 = math.Sincos(phi2)
	sin_e, cos_e       = math.Sincos(e)
)

type Movement struct {
	// Азимут
	A float64
	// Зенитный угол радианта
	Z_avg float64
	// Зенитное расстояние радианта
	Z_fix float64
	// Скорость V0 с учетом поправки за торможение
	V0 float64
	// Склонение радианта
	Delta float64
	// Прямое восхождение радианта. 0 < mov.Alpha < 2*pi
	Alpha float64
	// Эклиптическая широта радианта
	Beta float64
	// Эклиптическая долгота радианта
	Lambda float64
	// Долгота Солнца
	Lambda_theta float64
	// Гелиоцентрическая скорость
	V_h float64
	// Широта истинного радианта
	Beta_deriv float64
	// Часовой угол
	T float64
	// Долгота апекса
	Lambda_apex float64
	// Долгота радианта относительно апекса
	Diff_lambda float64
	// Долгота истинного радианта
	Lambda_deriv float64
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
	// Элонгация
	Psi_E float64
	// Геоцентрическая скорость
	V_g float64
	// Орбитальная скорость Земли для данного дня
	V_t float64
	// Афелий – точка орбиты максимально удаленная от Солнца
	Q float64
	// Большая полуось
	Axis float64
	// Угол, образуемый радиус-вектором метеорного тела с вектором его скорости
	Psi float64
	// Элонгация радианта от Солнца
	E_theta_deriv float64
	// Внеатмосферная скорость
	V_inf float64
	// Наклонение орбиты частицы к плоскости эклиптики
	Inc float64
	// Элонгация истинного радианта от апекса
	E_deriv float64
	// Перигелийное расстояние
	DistPer float64
	// Долгота восходящего узла
	Omega float64
	// Аргумент перигелия
	Wmega float64
	// Истинная аномалия
	Nu float64
	// Эксцентриситет
	Exc float64
	// Часовой угол радианта
	Angle float64
	// Радиус вектор орбиты Земли
	R float64
	// Угол элонгации видимого радианта от апекса движения Земли
	E_apex float64
}

type Input struct {
	Dist int
	// Временная задержка
	Tau1 float64
	// Временная задержка
	Tau2  float64
	V_avg float64
	// Время и дата появления метеороида
	Date time.Time
}

func NewMovement(inp Input) (mv *Movement) {
	mov := Movement{}
	var temp float64
	// step 1

	k := m * (inp.Tau1 / (inp.Tau2 + 1e-5))

	// Главное значение азимута
	A_gl := math.Atan2((cos_phi1 - k*cos_phi2), (k*sin_phi2 - sin_phi1))

	if inp.Tau1 <= 0 && inp.Tau2 < 0 {
		if A_gl > 0 {
			mov.A = A_gl
		} else if A_gl < 0 {
			mov.A = A_gl + math.Pi
		}
	} else if inp.Tau1 < 0 && inp.Tau2 >= 0 {
		if A_gl > 0 {
			mov.A = A_gl
		} else if A_gl < 0 {
			mov.A = A_gl + 2*math.Pi
		}
	} else if inp.Tau1 > 0 && inp.Tau2 <= 0 {
		mov.A = A_gl + math.Pi
	} else if inp.Tau1 > 0 && inp.Tau2 >= 0 {
		if A_gl > 0 {
			mov.A = A_gl + math.Pi
		} else if A_gl < 0 {
			mov.A = A_gl + 2*math.Pi
		}
	}
	sin_A, cos_A := math.Sincos(mov.A)

	// step 2

	temp = 2 * inp.V_avg * T
	cos_A_minus_phi1 := math.Cos(mov.A - phi1)
	cos_A_minus_phi2 := math.Cos(mov.A - phi2)
	sin_z1 := (temp * inp.Tau1) / (l1 * cos_A_minus_phi1)
	sin_z2 := (temp * inp.Tau2) / (l2 * cos_A_minus_phi2)

	W1 := math.Abs(cos_A_minus_phi1)
	W2 := math.Abs(cos_A_minus_phi2)
	mov.Z_avg = math.Asin((W1*sin_z1 + W2*sin_z2) / (W1 + W2))

	// step 3

	mov.V0 = _1_0398*inp.V_avg + _0_65

	// step 4

	// Скорость
	V_deriv := math.Sqrt(math.Pow(mov.V0, 2) - _123_2)
	delta_Z := 2 * math.Atan(math.Abs((V_deriv-mov.V0))/(V_deriv+mov.V0)*math.Tan(mov.Z_avg/inp.V_avg))
	mov.Z_fix = mov.Z_avg + delta_Z
	sin_Z_fix, cos_Z_fix := math.Sincos(mov.Z_fix)
	sin_Z_fix_cos_A := sin_Z_fix * cos_A
	sin_Z_fix_sin_A := sin_Z_fix * sin_A

	// step 5

	mov.Delta = math.Asin(sin_phi*cos_Z_fix - cos_phi*sin_Z_fix_cos_A)
	sin_delta, cos_delta := math.Sincos(mov.Delta)

	// step 6

	t_gl := math.Atan((sin_Z_fix_sin_A) / (cos_phi*cos_Z_fix + sin_phi*sin_Z_fix_cos_A))

	if t_gl >= 0 {
		if sin_Z_fix_sin_A >= 0 {
			mov.Angle = t_gl
		} else {
			mov.Angle = t_gl + math.Pi
		}
	} else { // t_gl < 0
		if sin_Z_fix_sin_A >= 0 {
			mov.Angle = t_gl + math.Pi
		} else {
			mov.Angle = t_gl + 2*math.Pi
		}
	}
	sin_t, cos_t := math.Sincos(mov.Angle)

	// step 7

	mov.S = StellarTime(c2, inp.Date.YearDay()-1, inp.Date.Hour(), inp.Date.Minute())

	// step 8

	mov.Alpha = mov.S - mov.Angle

	// step 9

	mov.Delta_alpha = -((_0_4703 * cos_t * cos_phi) / (V_deriv * cos_delta))
	mov.Delta_delta = -((_0_4703 * sin_t * sin_delta * cos_phi) / (V_deriv))

	// step 10

	mov.Alpha_fix = mov.Alpha + mov.Delta_alpha // while alpha_fix > 0
	sin_alpha_fix, cos_alpha_fix := math.Sincos(mov.Alpha_fix)
	mov.Delta_fix = mov.Delta + mov.Delta_delta
	sin_delta_fix, cos_delta_fix := math.Sincos(mov.Delta_fix)

	// step 11

	psi_E_gl := math.Acos(-sin_t * cos_delta_fix)
	mov.Psi_E = psi_E_gl
	if psi_E_gl < 0 {
		mov.Psi_E += math.Pi
	}

	// step 12

	delta_s := math.Sqrt(math.Pow(mov.Delta_alpha*cos_delta_fix, 2) + math.Pow(mov.Delta_delta, 2))
	mov.V_g = V_deriv * (math.Sin(mov.Psi_E-delta_s) / math.Sin(mov.Psi_E))
	pow_2_V_g := math.Pow(mov.V_g, 2)

	// step 13

	mov.V_inf = math.Sqrt(pow_2_V_g + _123_2)

	// step 14

	mov.Beta = math.Asin(-sin_e*sin_alpha_fix*cos_delta_fix + cos_e + sin_delta_fix)
	sin_beta, cos_beta := math.Sincos(mov.Beta)

	// step 15

	cos_lambda := ((cos_delta_fix * cos_alpha_fix) / cos_beta)
	sin_lambda := (1 / cos_beta) * (cos_delta_fix*sin_alpha_fix*cos_e + sin_delta_fix*sin_e)
	mov.Lambda = math.Atan2(sin_lambda, cos_lambda)

	// step 16

	mov.Lambda_theta = SolarLongitude(c3, inp.Date.YearDay()-1, inp.Date.Hour(), inp.Date.Minute())

	// step 17

	delta_theta := _0_01672 * math.Sin(mov.Lambda_theta-Pi0)
	mov.Lambda_apex = mov.Lambda_theta + delta_theta - (math.Pi / 2)
	mov.Diff_lambda = mov.Lambda - mov.Lambda_apex

	// step 18

	mov.E_apex = math.Acos(cos_beta * math.Cos(mov.Lambda-mov.Lambda_apex))

	// step 19

	mov.R = (1 - math.Pow(e0, 2)) / (1 - e0*math.Cos(mov.Lambda_theta-Pi0))

	// step 20

	mov.V_t = 29.76 * math.Sqrt((2/mov.R)-1)

	// step 21

	temp = mov.Lambda_theta + delta_theta - mov.Lambda
	sin_temp, cos_temp := math.Sincos(temp)
	temp_deriv_gl := math.Atan2(sin_temp-(mov.V_t/(mov.V_g*cos_beta)), cos_temp)
	temp_deriv := temp_deriv_gl
	if cos_temp < 0 {
		temp_deriv += math.Pi
	}
	mov.Lambda_deriv = mov.Lambda_theta + delta_theta - temp_deriv // 0 <= lambda_deriv <= 2*pi
	sin_lambda_diff, cos_lambda_diff := math.Sincos(mov.Lambda_theta - mov.Lambda_deriv)

	// step 22

	mov.V_h = math.Sqrt(pow_2_V_g + math.Pow(mov.V_t, 2) - 2*mov.V_g*mov.V_t*math.Cos(mov.E_apex))

	// step 23

	mov.Beta_deriv = math.Asin((mov.V_g / mov.V_h) * sin_beta)
	cos_beta_deriv := math.Cos(mov.Beta_deriv)

	// step 24

	mov.E_deriv = math.Acos(cos_beta_deriv * math.Cos(mov.Lambda_deriv-mov.Lambda_apex))

	// step 25

	i_gl := math.Atan(-(math.Abs(math.Tan(mov.Beta_deriv)) / sin_lambda_diff))
	mov.Inc = i_gl
	if i_gl > 0 {
		mov.Inc += math.Pi
	}

	// step 26

	mov.Q = math.Pow(mov.V_h/mov.V_t, 2)

	// step 27

	mov.Axis = 1 / ((2 - mov.Q) / mov.R)

	// step 28

	mov.Psi = math.Acos(-cos_beta_deriv * cos_lambda_diff)
	mov.E_theta_deriv = math.Pi - mov.Psi

	// step 29

	// Параметр орбиты
	p := math.Pow(mov.R, 2) * mov.Q * math.Pow(math.Sin(mov.E_theta_deriv), 2)
	b := math.Sqrt(p * math.Abs(mov.Axis))

	// step 30

	var c float64
	if mov.Axis > 0 {
		c = math.Sqrt(math.Pow(mov.Axis, 2) - math.Pow(b, 2))
	} else {
		c = math.Sqrt(math.Pow(mov.Axis, 2) + math.Pow(b, 2))
	}
	mov.Exc = math.Abs(c / mov.Axis)

	// step 31

	mov.DistPer = mov.Axis - c
	if mov.Axis < 0 {
		mov.DistPer = c - math.Abs(mov.Axis)
	}

	// step 32

	mov.Omega = mov.Lambda_theta
	if mov.Beta_deriv < 0 {
		mov.Omega += math.Pi
	}

	// step 33

	cos_v := (p - mov.R) / (mov.R * e)
	sin_v := p / (mov.R * e) * math.Cos(mov.Inc) * Ctg(mov.Lambda_deriv-mov.Lambda_theta)
	mov.Nu = math.Atan2(sin_v, cos_v)

	// step 34

	if mov.Beta_deriv > 0 {
		mov.Wmega = math.Pi - mov.Nu
	} else {
		mov.Wmega = -mov.Nu
	}

	mv = &mov
	return
}

// Format: 2006-01-02T03:04
func ParseDate(date string) (time.Time, error) {
	return time.Parse("2006-01-02T03:04", date)
}

func StellarTime(c2 float64, d, h, m int) (S float64) {
	S = c2 + 0.98565*float64(d) + 15.0411*float64(h) + 0.25068*float64(m)
	return
}

func SolarLongitude(c3 float64, d, h, m int) (lambda_theta float64) {
	lambda_theta = -c3 + 0.0000097*float64(m) + 0.000717*float64(h) + 0.017203*float64(d) + 0.034435*math.Sin(0.017203*float64(d-2))
	return
}
