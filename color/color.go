package color

import (
	"fmt"
	"io"
	"math"

	"github.com/sendelivery/go-trace-rays/interval"
	"github.com/sendelivery/go-trace-rays/vec3"
)

type Color = vec3.Vector3

func New(r, g, b float64) Color {
	return vec3.New(r, g, b)
}

func WriteColor(w io.Writer, pixelColor Color) {
	r := linearToGamma(pixelColor.X())
	g := linearToGamma(pixelColor.Y())
	b := linearToGamma(pixelColor.Z())

	// Translate [0,1] component values to the byte range [0,255]
	intensity := interval.New(0, 0.999)
	rByte := int(256 * intensity.Clamp(r))
	gByte := int(256 * intensity.Clamp(g))
	bByte := int(256 * intensity.Clamp(b))

	fmt.Fprintf(w, "%d %d %d\n", rByte, gByte, bByte)
}

func linearToGamma(lc float64) float64 {
	if lc <= 0 {
		return 0
	}
	return math.Sqrt(lc)
}

var Black = vec3.New(0, 0, 0)
var White = vec3.New(1, 1, 1)
