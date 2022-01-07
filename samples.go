package main

import (
	"fmt"
)

type SampleRequestOutOfBounds struct {
	requested float64
	duration  float64
}

func (s *SampleRequestOutOfBounds) Error() string {
	return fmt.Sprintf(
		"The time requested is longer than the duration of the samples (%f > %f)",
		s.requested,
		s.duration,
	)
}

type PlaybackConfig struct {
	rate     int
	duration float64
	dt       float64
}

type Samples interface {
	Fetch(t float64) float64
	Config() PlaybackConfig
}

type ProceduralSamples struct {
	config   PlaybackConfig
	function func(t float64) float64
}

func InitializeProceduralSamples(
	config PlaybackConfig,
	function func(t float64) float64,
) ProceduralSamples {
	return ProceduralSamples{
		config:   config,
		function: function,
	}
}

func (s ProceduralSamples) Fetch(t float64) float64 {
	if t > s.config.duration {
		panic(&SampleRequestOutOfBounds{t, s.config.duration})
	}
	return s.function(t)
}

func (s ProceduralSamples) Config() PlaybackConfig {
	return s.config
}

type RawSamples struct {
	config PlaybackConfig
	data   []float64
}

func InitializeRawSamples(rate int, data []float64) RawSamples {
	if rate < 0 {
		rate = len(data)
	}

	duration := float64(len(data)) / float64(rate)
	return RawSamples{
		config: PlaybackConfig{
			rate:     rate,
			duration: duration,
			dt:       duration / float64(rate),
		},
		data: data,
	}
}

func (s RawSamples) Fetch(t float64) float64 {
	if t > s.config.duration {
		panic(&SampleRequestOutOfBounds{t, s.config.duration})
	}
	return s.data[int(float64(len(s.data))*t)]
}

func (s RawSamples) Config() PlaybackConfig {
	return s.config
}
