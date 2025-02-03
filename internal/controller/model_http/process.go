package model_http

import (
	"github.com/Mopsgamer/space-soup/internal/math"
)

type Process struct {
	Tau1     float64 `form:"t1"`
	Tau2     float64 `form:"t2"`
	SpeedAvg float64 `form:"v_avg"`
}

func (p *Process) NewOrbit() (orbit *math.Orbit, err error) {
	return math.NewOrbit(p.Tau1, p.Tau2, p.SpeedAvg)
}
