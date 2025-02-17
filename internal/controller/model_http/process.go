package model_http

import (
	"time"

	"github.com/Mopsgamer/space-soup/internal/math"
)

type Process struct {
	Tau1  float64 `form:"tau1"`
	Tau2  float64 `form:"tau2"`
	V_avg float64 `form:"v_avg"`
	Date  string  `form:"date"`
}

func (p *Process) NewMeteor() (meteor *math.Meteor, err error) {
	t, err := time.Parse("2006-01-02T03:04", p.Date)
	if err != nil {
		return nil, err
	}

	return math.NewMeteor(p.Tau1, p.Tau2, p.V_avg, t)
}
