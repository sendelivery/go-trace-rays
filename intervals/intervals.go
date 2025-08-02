package intervals

import "math"

type Interval struct {
	Min, Max float64
}

func Default() Interval {
	return Interval{Min: math.Inf(1), Max: math.Inf(-1)}
}

func New(min, max float64) Interval {
	return Interval{Min: min, Max: max}
}

func (i Interval) Size() float64 {
	return i.Max - i.Min
}

func (i Interval) Contains(x float64) bool {
	return i.Min <= x && x <= i.Max
}

func (i Interval) Surrounds(x float64) bool {
	return i.Min < x && x < i.Max
}

// EmptyInterval is an interval within which nothing lies
var EmptyInterval = Default()

// UniverseInterval encompasses everything
var UniverseInterval = New(math.Inf(-1), math.Inf(1))
