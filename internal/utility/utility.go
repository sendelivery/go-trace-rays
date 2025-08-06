package utility

import "math/rand"

const (
	PI = 3.1415926535897932385
)

func Deg2Rad(d float64) float64 {
	return d * PI / 180
}

func Random() float64 {
	return rand.Float64()
}

func RandomN(min, max float64) float64 {
	return min + (max-min)*Random()
}
