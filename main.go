package main

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func main() {
	samples := InitializeSamplesFromWav()

	p := plot.New()

	p.X.Min = MIN_EQ_FREQ
	p.X.Max = MAX_EQ_FREQ
	p.Y.Min = 0
	p.Y.Max = 1

	p.X.Scale = plot.LogScale{}

	exp := plotter.NewFunction(func(x float64) float64 {
		return clamp(transform(x, samples)*1000, 0, 1)
	})

	exp.Samples = 500
	p.Add(exp)

	if err := p.Save(19*vg.Inch, 10*vg.Inch, "points.png"); err != nil {
		panic(err)
	}
}
