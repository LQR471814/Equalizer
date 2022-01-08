package main

import (
	"math"
)

type EvaluationSettings struct {
	function func(x float64) float64
	dx       float64
	min      float64
	max      float64
}

func (s EvaluationSettings) evaluate(callback func(v float64)) {
	x := s.min
	for x < s.max {
		callback(s.dx * s.function(x))
		x += s.dx
	}
}

func integrate(s EvaluationSettings) float64 {
	var sum float64 = 0
	s.evaluate(func(v float64) { sum += v })
	return sum
}

func transform(frequency float64, samples Samples) float64 {
	return integrate(
		EvaluationSettings{
			min: 0,
			max: samples.Config().duration,
			dx:  samples.Config().dt(),
			function: func(x float64) float64 {
				return samples.Fetch(x) * math.Cos(2*frequency*math.Pi*x)
			},
		},
	)
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	} else if v > max {
		return max
	}

	return v
}
