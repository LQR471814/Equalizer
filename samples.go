package main

import (
	"fmt"
	"io"
	"os"

	"github.com/youpy/go-wav"
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
	rate     uint32
	duration float64
}

func (c PlaybackConfig) dt() float64 {
	return c.duration / float64(c.rate)
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

func InitializeSamplesFromWav() RawSamples {
	file, err := os.Open("Jangle.wav")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	reader := wav.NewReader(file)

	format, err := reader.Format()
	if err != nil {
		panic(err)
	}

	duration, err := reader.Duration()
	if err != nil {
		panic(err)
	}

	result := []float64{}

	for {
		samples, err := reader.ReadSamples()
		if err == io.EOF {
			break
		}

		s := make([]float64, len(samples))
		for i, sample := range samples {
			s[i] = (reader.FloatValue(sample, 0) + reader.FloatValue(sample, 1)) / 2
		}

		result = append(result, s...)
	}

	return RawSamples{
		config: PlaybackConfig{
			rate:     format.SampleRate,
			duration: duration.Seconds(),
		},
		data: result,
	}
}

func InitializeRawSamples(rate uint32, data []float64) RawSamples {
	return RawSamples{
		config: PlaybackConfig{
			rate:     rate,
			duration: float64(len(data)) / float64(rate),
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
