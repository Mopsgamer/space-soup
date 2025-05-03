package model_http

import (
	"github.com/Mopsgamer/space-soup/server/soup"
)

type OrbitInput struct {
	Tau1 float64 `form:"tau1"`
	Tau2 float64 `form:"tau2"`
	V1   float64 `form:"v1"`
	V2   float64 `form:"v2"`
	V3   float64 `form:"v3"`
	Date string  `form:"date"`
}

func (p *OrbitInput) Input() (soup.Input, error) {
	date, err := soup.ParseDateJSON(p.Date)
	if err != nil {
		return soup.Input{}, err
	}

	speedList := []float64{}
	for _, v := range []float64{p.V1, p.V2, p.V3} {
		if v > 999. {
			continue
		}
		speedList = append(speedList, v)
	}

	return soup.Input{
		Tau1:  p.Tau1,
		Tau2:  p.Tau2,
		V_avg: soup.Average(speedList),
		Date:  date,
	}, nil
}
