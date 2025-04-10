package soup

import (
	"fmt"
	"math"
	"time"
)

var (
	allowedDeltaDegrees = 4.
	allowedDeltaRadians = RadiansFromDegrees(4) // 0.06981317007977318
	allowedDeltaSpeed   = 2.
	allowedDeltaAxis    = 0.1
	allowedDeltaExc     = 0.1
)

var (
	c2 = 1.61667
	c3 = 1.40042

	EarthRadius = 6371.12

	_2Pi    = 2 * math.Pi
	_halfPi = math.Pi / 2
	Pi0     = RadiansFromRich(102, 21, 14) // 1.7864123633315516

	_29_76     = 29.76                      // V_t?
	_26_948deg = RadiansFromDegrees(26.948) // 0.47033132682743195
	_0_9252    = 9252e-7
	_0_4749    = 4749e-7
	_0_65      = 0.65
	_1_0398    = 1.0398

	_123_2 = 123.2

	e    = RadiansFromRich(23, 26, 40) // 0.4091827468564484
	e0   = 0.01675
	m    = 1.94787
	phi  = RadiansFromRich(49, 24, 50)  // Широта місця спостереження, 0.8624350573257534
	phi1 = RadiansFromRich(34, 10, 16)  // 34°10'16'', 0.5963983979537067
	phi2 = RadiansFromRich(110, 16, 22) // 110°16'22'', 1.924623047542258

	// оптимізації
	sin_phi, cos_phi   = math.Sincos(phi)
	sin_phi1, cos_phi1 = math.Sincos(phi1)
	sin_phi2, cos_phi2 = math.Sincos(phi2)
	sin_e, cos_e       = math.Sincos(e)
)

type MovementGeneric[T any] struct {
	// Азимут
	A T
	// Зенітний кут радіанта
	Z_avg T
	// Зенітна відстань радіанта
	Z_fix T
	// Швидкість V0 з урахуванням поправки на гальмування
	V0 T
	// Схилення радіанта
	Delta T
	// Пряме сходження радіанта. 0 < mov.Alpha < 2*pi
	Alpha T
	// Екліптична широта радіанта
	Beta T
	// Екліптична довгота радіанта
	Lambda T
	// Довгота Сонця
	Lambda_theta T
	// Геліоцентрична швидкість
	V_h T
	// Широта істинного радіанта
	Beta_deriv T
	// Годинний кут
	T T
	// Висота
	H T
	// Довгота апекса
	Lambda_apex T
	// Довгота радіанта відносно апекса
	Diff_lambda T
	// Довгота істинного радіанта
	Lambda_deriv T
	// Зоряний час у момент спостереження
	S T
	// Виправлені екваторіальні координати радіанта
	Alpha_fix T
	// Виправлені екваторіальні координати радіанта
	Delta_fix T
	// Поправки на добову аберацію в екваторіальних координатах
	Delta_alpha T
	// Поправки на добову аберацію в екваторіальних координатах
	Delta_delta T
	// Елонгація
	Psi_E T
	// Геоцентрична швидкість
	V_g T
	// Орбітальна швидкість Землі для даного дня
	V_t T
	// Афелій – точка орбіти максимально віддалена від Сонця
	Q T
	// Велика піввісь
	Axis T
	// Кут, утворений радіус-вектором метеорного тіла з вектором його швидкості
	Psi T
	// Елонгація радіанта від Сонця
	E_theta_deriv T
	// Позаземна швидкість
	V_inf T
	// Нахил орбіти частинки до площини екліптики
	Inc T
	// Елонгація істинного радіанта від апекса
	E_deriv T
	// Перигелійна відстань
	DistPer T
	// Довгота висхідного вузла
	Omega T
	// Аргумент перигелію
	Wmega T
	// Істинна аномалія
	Nu T
	// Ексцентриситет
	Exc T
	// Годинний кут радіанта
	Angle T
	// Радіус-вектор орбіти Землі
	R T
	// Кут елонгації видимого радіанта від апекса руху Землі
	E_apex T
}

type Movement = MovementGeneric[float64]

// 0 - Успіх, 1 - Прийнятно, 2 - Неприйнятно
type MovementAssertion = MovementGeneric[uint]

type InputGeneric[T any] struct {
	Id *int
	// Нахилена дальність
	Dist T
	// Часова затримка
	Tau1 T
	// Часова затримка
	Tau2  T
	V_avg T
	// Час і дата появи метеороїда
	Date time.Time
}

type Input = InputGeneric[float64]

func NewMovement(inp Input) (*Movement, error) {
	mov := Movement{}

	// крок 1

	k := m * (inp.Tau1 / (inp.Tau2 + 1e-5))

	// Головне значення азимута
	A_gl := math.Atan2((cos_phi1 - k*cos_phi2), (k*sin_phi2 - sin_phi1))
	A_gl = LoopNumber(A_gl, -_halfPi, _halfPi)

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
			mov.A = A_gl + _2Pi
		}
	} else if inp.Tau1 > 0 && inp.Tau2 <= 0 {
		mov.A = A_gl + math.Pi
	} else if inp.Tau1 > 0 && inp.Tau2 >= 0 {
		if A_gl > 0 {
			mov.A = A_gl + math.Pi
		} else if A_gl < 0 {
			mov.A = A_gl + _2Pi
		}
	}
	sin_A, cos_A := math.Sincos(mov.A)

	// крок 2

	cos_A_substract_phi1 := math.Cos(mov.A - phi1)
	cos_A_substract_phi2 := math.Cos(mov.A - phi2)
	sin_z1 := -(_0_9252 * inp.V_avg * inp.Tau1) / cos_A_substract_phi1
	sin_z2 := -(_0_4749 * inp.V_avg * inp.Tau2) / cos_A_substract_phi2

	z1 := math.Asin(sin_z1)
	z2 := math.Asin(sin_z2)
	if z1 < 0 || z2 < 0 {
		return &mov, fmt.Errorf("z1 (%v) or z2 (%v) less than 0", z1, z2)
	}
	if delta := z1 - z2; DegreesFromRadians(delta) >= 4 {
		return &mov, fmt.Errorf("z1 (%v) and z2 (%v) delta (%v) greater than 4 deg (%v)", z1, z2, RadiansFromDegrees(delta), RadiansFromDegrees(4))
	}

	if inp.Tau1 == 0 {
		mov.Z_avg = z2
	} else if inp.Tau2 == 0 {
		mov.Z_avg = z2
	} else {
		mov.Z_avg = (z1 + z2) / 2
	}

	// крок 3

	mov.V0 = _1_0398*inp.V_avg + _0_65

	// крок 4

	// Швидкість
	V_deriv := math.Sqrt(math.Pow(mov.V0, 2) - _123_2)
	abs_VV := math.Abs((V_deriv - mov.V0) / (V_deriv + mov.V0))
	tan_delta_Z := abs_VV * math.Tan(mov.Z_avg/2)
	delta_Z := 2 * math.Atan(tan_delta_Z)
	mov.Z_fix = mov.Z_avg + delta_Z
	sin_Z_fix, cos_Z_fix := math.Sincos(mov.Z_fix)
	sin_Z_fix_cos_A := sin_Z_fix * cos_A
	sin_Z_fix_sin_A := sin_Z_fix * sin_A

	// крок 5

	mov.Delta = math.Asin(sin_phi*cos_Z_fix - cos_phi*sin_Z_fix_cos_A)
	mov.Delta = LoopNumber(mov.Delta, -_halfPi, _halfPi)
	sin_delta, cos_delta := math.Sincos(mov.Delta)

	// крок 6

	sin_t := (sin_Z_fix_sin_A) / cos_delta
	mov.Angle = math.Asin(sin_t)
	mov.Angle = LoopNumber(mov.Angle, 0, _2Pi)
	cos_t := (math.Cos(phi)*math.Cos(mov.Z_fix) + math.Sin(phi)*math.Sin(mov.Z_fix)*math.Cos(mov.A)) / mov.Angle

	// крок 7

	d := float64(inp.Date.YearDay() - 1)
	h := float64(inp.Date.Hour())
	m := float64(inp.Date.Minute())
	dMod := 0.98565
	hMod := 15.0411
	mMod := 0.25068
	mov.S = c2 + RadiansFromDegrees(dMod)*d + RadiansFromDegrees(hMod)*h + RadiansFromDegrees(mMod)*m

	// крок 8

	mov.Alpha = mov.S - mov.Angle
	mov.Alpha = LoopNumber(mov.Alpha, 0, _2Pi)

	// крок 9

	mov.Delta_alpha = -(_26_948deg / V_deriv) * (cos_t / cos_delta) * cos_phi
	mov.Delta_delta = (_26_948deg / V_deriv) * sin_t * sin_delta * cos_phi

	// крок 10

	mov.Alpha_fix = mov.Alpha + mov.Delta_alpha
	if mov.Alpha_fix <= 0 {
		return &mov, fmt.Errorf("alpha_fix (%v) should be greater than 0", mov.Alpha_fix)
	}
	sin_alpha_fix, cos_alpha_fix := math.Sincos(mov.Alpha_fix)

	mov.Delta_fix = mov.Delta + mov.Delta_delta
	sin_delta_fix, cos_delta_fix := math.Sincos(mov.Delta_fix)

	// крок 11

	psi_E_gl := math.Acos(-sin_t * cos_delta_fix)
	if psi_E_gl > 0 {
		mov.Psi_E = psi_E_gl
	} else if psi_E_gl < 0 {
		mov.Psi_E = psi_E_gl + math.Pi
	}

	// крок 12

	delta_s := math.Sqrt(math.Pow(mov.Delta_alpha*cos_delta_fix, 2) + math.Pow(mov.Delta_delta, 2))
	mov.V_g = V_deriv * (math.Sin(mov.Psi_E-delta_s) / math.Sin(mov.Psi_E))
	pow_2_V_g := math.Pow(mov.V_g, 2)

	// крок 13

	mov.V_inf = math.Sqrt(pow_2_V_g + _123_2)

	// крок 14

	sin_beta := -sin_e*sin_alpha_fix*cos_delta_fix + cos_e*sin_delta_fix
	mov.Beta = math.Asin(sin_beta)
	mov.Beta = LoopNumber(mov.Beta, -_halfPi, _halfPi)
	cos_beta := math.Cos(mov.Beta)

	// крок 15

	cos_lambda := cos_delta_fix * cos_alpha_fix / cos_beta
	sin_lambda := (cos_delta_fix*sin_alpha_fix*cos_e + sin_delta_fix*sin_e) / cos_beta
	mov.Lambda = math.Atan2(sin_lambda, cos_lambda)
	mov.Lambda = LoopNumber(mov.Lambda, 0, _2Pi)

	// крок 16

	dMod = 0.017202
	hMod = 0.000717
	mMod = 0.0000097
	lambdaMod := 0.034435
	mov.Lambda_theta = dMod*d + hMod*h + mMod*m + lambdaMod*math.Sin(dMod*(d-2)) - c3

	// крок 17

	delta_theta := 0.01677 * math.Sin(mov.Lambda_theta-Pi0)
	mov.Lambda_apex = mov.Lambda_theta + delta_theta - _halfPi
	mov.Lambda_apex = LoopNumber(mov.Lambda_apex, 0, _2Pi)
	mov.Diff_lambda = mov.Lambda - mov.Lambda_apex
	mov.Diff_lambda = LoopNumber(mov.Diff_lambda, 0, _2Pi)

	// крок 18

	mov.E_apex = math.Acos(cos_beta * math.Cos(mov.Diff_lambda))
	mov.E_apex = LoopNumber(mov.E_apex, 0, math.Pi)

	// крок 19

	mov.R = (1 - math.Pow(e0, 2)) / (1 - e0*math.Cos(mov.Lambda_theta-Pi0))

	// крок 20

	mov.V_t = _29_76 * math.Sqrt((2/mov.R)-1)

	// крок 21

	temp := mov.Lambda_theta + delta_theta - mov.Lambda
	sin_temp, cos_temp := math.Sincos(temp)
	temp_deriv_gl := math.Atan2(sin_temp-(mov.V_t/(mov.V_g*cos_beta)), cos_temp)
	temp_deriv := LoopNumber(temp_deriv_gl, 0, _2Pi)
	mov.Lambda_deriv = mov.Lambda_theta + delta_theta - temp_deriv
	mov.Lambda_deriv = LoopNumber(mov.Lambda_deriv, 0, _2Pi)

	sin_lambda_diff, cos_lambda_diff := math.Sincos(mov.Lambda_theta - mov.Lambda_deriv)

	// крок 22

	mov.V_h = math.Sqrt(pow_2_V_g + math.Pow(mov.V_t, 2) - 2*mov.V_g*mov.V_t*math.Cos(mov.E_apex))

	// крок 23

	sin_beta_deriv := (mov.V_g / mov.V_h) * sin_beta
	mov.Beta_deriv = math.Asin(sin_beta_deriv)
	mov.Beta_deriv = LoopNumber(mov.Beta_deriv, -_halfPi, _halfPi)
	cos_beta_deriv := math.Cos(mov.Beta_deriv)

	// крок 24

	mov.E_deriv = math.Acos(cos_beta_deriv * math.Cos(mov.Lambda_deriv-mov.Lambda_apex))

	// крок 25

	i_gl := math.Atan2(-math.Abs(math.Tan(mov.Beta_deriv)), sin_lambda_diff)
	if i_gl > 0 {
		mov.Inc = i_gl
	} else if i_gl < 0 {
		mov.Inc = i_gl + math.Pi
	}
	mov.Inc = LoopNumber(mov.Inc, 0, math.Pi)

	// крок 26

	mov.Q = math.Pow(mov.V_h/_29_76, 2)

	// крок 27

	// mov.Axis = 1 / ((2 - mov.Q) / mov.R)
	mov.Axis = mov.R / (2 - mov.R*mov.Q)

	// крок 28

	mov.Psi = math.Acos(-cos_beta_deriv * cos_lambda_diff)
	mov.Psi = LoopNumber(mov.Psi, 0, math.Pi)
	mov.E_theta_deriv = math.Pi - mov.Psi

	// крок 29

	// Параметр орбіти
	p := math.Pow(mov.R, 2) * mov.Q * math.Pow(math.Sin(mov.E_theta_deriv), 2)
	b := math.Sqrt(p * math.Abs(mov.Axis))

	// крок 30

	var c float64
	if mov.Axis > 0 {
		c = math.Sqrt(math.Pow(mov.Axis, 2) - math.Pow(b, 2))
	} else if mov.Axis < 0 {
		c = math.Sqrt(math.Pow(mov.Axis, 2) + math.Pow(b, 2))
	}
	mov.Exc = math.Abs(c / mov.Axis)

	// крок 31

	if mov.Axis > 0 {
		mov.DistPer = mov.Axis - c
	} else if mov.Axis < 0 {
		mov.DistPer = c - math.Abs(mov.Axis)
	}

	// крок 32

	if mov.Beta_deriv > 0 {
		mov.Omega = mov.Lambda_theta
	} else if mov.Beta_deriv < 0 {
		mov.Omega = mov.Lambda_theta + math.Pi
	}
	mov.Omega = LoopNumber(mov.Omega, 0, _2Pi)

	// крок 33

	cos_v := (p - mov.R) / (mov.R * e)
	sin_v := p / (mov.R * e) * math.Cos(mov.Inc) * Ctg(mov.Lambda_deriv-mov.Lambda_theta)
	mov.Nu = math.Atan2(sin_v, cos_v)
	mov.Nu = LoopNumber(mov.Nu, 0, _2Pi)

	// крок 34

	mov.Wmega = LoopNumber(-mov.Nu+math.Pi, 0, _2Pi)

	// eps := math.Atan(math.Tan(mov.Z_fix) * math.Cos(mov.A-mov.Axis))
	// sin_eps, cos_eps := math.Sincos(eps)
	// mov.H = inp.Dist*sin_eps + math.Pow(inp.Dist, 2)*math.Pow(cos_eps, 2)/(2*EarthRadius)
	// mov.H = DegreesFromRadians(_90deg - mov.Axis)
	return &mov, nil
}
