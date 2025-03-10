package model_http

import (
	"github.com/Mopsgamer/space-soup/server/soup"
)

type OrbitInput struct {
	Dist int     `form:"dist"`
	Tau1 float64 `form:"tau1"`
	Tau2 float64 `form:"tau2"`
	V1   float64 `form:"v1"`
	V2   float64 `form:"v2"`
	V3   float64 `form:"v3"`
	Date string  `form:"date"`
}

func (p *OrbitInput) Movement() (*soup.Movement, error) {
	date, err := soup.ParseDate(p.Date)
	if err != nil {
		return nil, err
	}

	return soup.NewMovement(soup.Input{
		Dist:  p.Dist,
		Tau1:  p.Tau1,
		Tau2:  p.Tau2,
		V_avg: soup.Average([]float64{p.V1, p.V2, p.V3}),
		Date:  date,
	}), nil
}
