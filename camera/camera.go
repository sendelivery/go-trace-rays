package camera

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/sendelivery/go-trace-rays/color"
	"github.com/sendelivery/go-trace-rays/hittable"
	"github.com/sendelivery/go-trace-rays/intervals"
	"github.com/sendelivery/go-trace-rays/ray"
	"github.com/sendelivery/go-trace-rays/vec3"
)

type Camera struct {
	AspectRatio float64      // Ratio of image width over height
	ImageWidth  int          // Rendered image width in pixel count
	imageHeight int          // Rendered image height
	centre      vec3.Vector3 // Camera center
	pixel00Loc  vec3.Vector3 // Location of pixel 0, 0
	pixelDeltaU vec3.Vector3 // Offset to pixel to the right
	pixelDeltaV vec3.Vector3 // Offset to pixel below
}

func New() *Camera {
	c := Camera{
		AspectRatio: 1.0,
		ImageWidth:  100,
	}
	return &c
}

func (c *Camera) Render(world hittable.Hittabler) {
	c.initialise()

	// Render
	fmt.Printf("P3\n%d %d\n255\n", c.ImageWidth, c.imageHeight)

	start := time.Now()

	for j := range c.imageHeight {
		fmt.Fprintf(os.Stderr, "\rScanlines remaining: %d ", c.imageHeight-j)
		for i := range c.ImageWidth {
			pixelCentre := vec3.Add(
				c.pixel00Loc,
				vec3.Add(
					vec3.Mulf(c.pixelDeltaU, float64(i)),
					vec3.Mulf(c.pixelDeltaV, float64(j)),
				),
			)

			rayDirection := vec3.Sub(pixelCentre, c.centre)
			r := ray.New(c.centre, rayDirection)

			col := c.rayColor(r, world)
			color.WriteColor(os.Stdout, col)
		}
	}

	elapsed := time.Since(start).Milliseconds()

	fmt.Fprintf(os.Stderr, "\rDone in %dms.         \n", elapsed)
}

func (c *Camera) initialise() {
	// Calculate image height, ensuring it's at least 1
	c.imageHeight = max(int(float64(c.ImageWidth)/c.AspectRatio), 1)

	c.centre = vec3.New(0, 0, 0)

	// Determine viewport dimensions
	focalLength := 1.0 // Distance between camera centre and viewport
	viewportHeight := 2.0
	viewportWidth := viewportHeight * float64(c.ImageWidth) / float64(c.imageHeight)

	// Calculate the vectors across the horizontal and down the vertical viewport edges
	viewportU := vec3.New(viewportWidth, 0, 0)
	viewportV := vec3.New(0, -viewportHeight, 0)

	// Calculate the horizontal and vertical delta vectors from pixel to pixel
	c.pixelDeltaU = vec3.Div(viewportU, float64(c.ImageWidth))
	c.pixelDeltaV = vec3.Div(viewportV, float64(c.imageHeight))

	// Calculate the location of the upper left pixel
	viewportUpperLeft := vec3.Duplicate(c.centre)
	viewportUpperLeft.Sub(
		vec3.New(0, 0, focalLength),
	).Sub(
		vec3.Div(viewportU, 2),
	).Sub(
		vec3.Div(viewportV, 2),
	)
	c.pixel00Loc = vec3.Add(viewportUpperLeft, vec3.Mulf(vec3.Add(c.pixelDeltaU, c.pixelDeltaV), 0.5)) // Inset pixels
}

const dampen = 0.5

func (c *Camera) rayColor(r ray.Ray, world hittable.Hittabler) color.Color {
	if hr, ok := world.Hit(r, intervals.New(0, math.Inf(1))); ok {
		return vec3.Mulf(vec3.Add(hr.Normal(), color.New(1, 1, 1)), 0.5)
	}

	unitDirection := vec3.UnitVector(r.Direction())
	a := dampen * (unitDirection.Y() + 1)
	white := color.New(1, 1, 1)
	blue := color.New(0.5, 0.7, 1)

	return vec3.Add(
		vec3.Mulf(white, (1.0-a)),
		vec3.Mulf(blue, a),
	)
}
