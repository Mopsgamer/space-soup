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
	Tests       []MovementTest
	Description string
}

func VisualizeDeclRasc(config VisualizeConfig) (io.WriterTo, error) {
	p := plot.New()
	p.Title.Text = fmt.Sprintf(
		"Generation Date: %s | Orbit Count: %d",
		time.Now().Format("2006-01-02 15:04:05"), len(config.Tests),
	)
	if config.Description != "" {
		p.Title.Text += " | " + config.Description
	}
	p.X.Label.Text, p.Y.Label.Text = "Right Ascension", "Declination"
	p.X.Tick.Length, p.Y.Tick.Length = 5, 5
	p.X.Tick.Marker, p.Y.Tick.Marker = AngleTicker{Step: 45}, AngleTicker{Step: 10}
	p.X.Min, p.X.Max = 0, 360
	p.Y.Min, p.Y.Max = -90, 90

	pointsActual := plotter.XYs{}
	// pointsExpected := plotter.XYs{}
	for _, m := range config.Tests {
		if m.Actual.Fail != nil {
			continue
		}
		pointsActual = append(pointsActual, plotter.XY{X: m.Actual.Alpha, Y: m.Actual.Delta})
		// pointsExpected = append(pointsActual, plotter.XY{X: m.Expected.Alpha, Y: m.Expected.Delta})
	}
	scatter, err := plotter.NewScatter(pointsActual)
	if err != nil {
		return nil, err
	}
	scatter.Shape = draw.CircleGlyph{}
	scatter.GlyphStyle.Radius = vg.Points(.4)
	p.Add(scatter)

	ecliptic := plotter.NewFunction(func(x float64) float64 {
		amplitude := 23.5
		frequency := 2 * math.Pi / 360
		phase := 0.0

		return amplitude * math.Sin(frequency*x+phase)
	})
	ecliptic.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	ecliptic.Width = vg.Points(1)
	p.Add(ecliptic)
	p.Title.Padding = 10
	p.Legend.Add("Ecliptic", ecliptic)
	// scatter2, err := plotter.NewScatter(pointsExpected)
	// if err != nil {
	// 	return err
	// }
	// scatter2.Shape = draw.CircleGlyph{}
	// scatter2.GlyphStyle.Radius = vg.Points(.2)
	// scatter2.GlyphStyle.Color = color.RGBA{255, 0, 0, 255}
	// p.Add(scatter2)

	// Add a grid
	grid := plotter.NewGrid()
	p.Add(grid)

	scale := 1
	w, h := 4*80*scale, 3*80*scale
	return p.WriterTo(vg.Length(w*scale), vg.Length(h*scale), "png")
}
