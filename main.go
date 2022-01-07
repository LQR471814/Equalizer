package main

import (
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func main() {
	var targetFreq float64 = 2

	samples := InitializeProceduralSamples(
		PlaybackConfig{
			rate:     100,
			duration: 1,
			dt:       0.01,
		},
		func(x float64) float64 {
			return math.Cos(2*targetFreq*math.Pi*x) + math.Cos(2*5*math.Pi*x)
		},
	)

	plot := plot.New()

	plot.X.Min = 0
	plot.X.Max = 100
	plot.Y.Min = -0.2
	plot.Y.Max = 1

	exp := plotter.NewFunction(func(x float64) float64 {
		return transform(float64(x), samples)
	})
	exp.Samples = 1000

	plot.Add(exp)

	if err := plot.Save(8*vg.Inch, 8*vg.Inch, "points.png"); err != nil {
		panic(err)
	}
}
