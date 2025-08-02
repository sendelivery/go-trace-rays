package color

import (
	"fmt"
	"io"

	"github.com/sendelivery/go-trace-rays/vec3"
)

type Color = vec3.Vector3

func New(r, g, b float64) Color {
	return vec3.New(r, g, b)
}

func WriteColor(w io.Writer, pixelColor Color) {
	rByte := int(255.999 * pixelColor.X())
	gByte := int(255.999 * pixelColor.Y())
	bByte := int(255.999 * pixelColor.Z())

	fmt.Fprintf(w, "%d %d %d\n", rByte, gByte, bByte)
}
