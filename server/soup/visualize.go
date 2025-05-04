package soup

import (
	"fmt"
	"image/color"
	"io"
	"math"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

type AngleTicker struct {
	Step float64
}

func (t AngleTicker) Ticks(min, max float64) []plot.Tick {
	var ticks []plot.Tick
	for val := min; val <= max; val += t.Step {
		ticks = append(ticks, plot.Tick{Value: val, Label: fmt.Sprintf("%.1f°", val)})
	}
	return ticks
}

type VisualizeConfig struct {
	Dark           bool
	Tests          []MovementTest
	Description    string
	XLabel, YLabel string
	GetXY          func(m Movement) (x float64, y float64)
}

func Visualize(config VisualizeConfig) (io.WriterTo, error) {
	p := plot.New()
	p.Title.Text = fmt.Sprintf(
		"Generation Date: %s | Orbit Count: %d",
		time.Now().Format("2006-01-02 15:04:05"), len(config.Tests),
	)
	p.Title.Padding = 10

	if config.Description != "" {
		p.Title.Text += " | " + config.Description
	}

	p.X.Label.Text, p.Y.Label.Text = config.XLabel, config.YLabel
	p.X.Tick.Length, p.Y.Tick.Length = 5, 5
	p.X.Tick.Marker, p.Y.Tick.Marker = AngleTicker{Step: 45}, AngleTicker{Step: 10}
	p.X.Min, p.X.Max = 0, 360
	p.Y.Min, p.Y.Max = -90, 90

	pointsActual := plotter.XYs{}
	for _, m := range config.Tests {
		if m.Actual.Fail != nil {
			continue
		}
		xActual, yActual := config.GetXY(m.Actual)
		pointsActual = append(pointsActual, plotter.XY{X: xActual, Y: yActual})
	}
	scatter, err := plotter.NewScatter(pointsActual)
	if err != nil {
		return nil, err
	}
	scatter.Shape = draw.CircleGlyph{}
	scatter.GlyphStyle.Radius = vg.Points(.4)

	ecliptic := plotter.NewFunction(func(x float64) float64 {
		amplitude := 23.5
		frequency := 2 * math.Pi / 360
		phase := 0.0

		return amplitude * math.Sin(frequency*x+phase)
	})
	ecliptic.Width = vg.Points(1)
	p.Legend.TextStyle.Color = color.White
	p.X.Padding = vg.Centimeter

	grid := plotter.NewGrid()

	if config.Dark {
		p.BackgroundColor = color.RGBA{R: 0, G: 0, B: 0, A: 255}
		p.Title.TextStyle.Color = color.White

		p.X.Label.TextStyle.Color = color.White
		p.Y.Label.TextStyle.Color = color.White
		p.X.Tick.Label.Color = color.White
		p.Y.Tick.Label.Color = color.White
		p.X.Tick.Color = color.White
		p.Y.Tick.Color = color.White
		p.X.Color = color.White
		p.Y.Color = color.White

		grid.Horizontal.Color = color.RGBA{0, 0, 255, 50}
		grid.Vertical.Color = grid.Horizontal.Color

		scatter.Color = color.White

		ecliptic.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	}

	p.Add(ecliptic)
	p.Add(grid)
	p.Legend.Add("Ecliptic", ecliptic)
	p.Add(scatter)

	w, h := 16*45, 9*45
	return p.WriterTo(vg.Length(w), vg.Length(h), "png")
}
