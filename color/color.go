package color

import (
	"fmt"
	"io"

	"github.com/sendelivery/go-trace-rays/intervals"
	"github.com/sendelivery/go-trace-rays/vec3"
)

type Color = vec3.Vector3

func New(r, g, b float64) Color {
	return vec3.New(r, g, b)
}

func WriteColor(w io.Writer, pixelColor Color) {
	r := pixelColor.X()
	g := pixelColor.Y()
	b := pixelColor.Z()

	// Translate [0,1] component values to the byte range [0,255]
	intensity := intervals.New(0, 0.999)
	rByte := int(256 * intensity.Clamp(r))
	gByte := int(256 * intensity.Clamp(g))
	bByte := int(256 * intensity.Clamp(b))

	fmt.Fprintf(w, "%d %d %d\n", rByte, gByte, bByte)
}
