package camera

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/sendelivery/go-trace-rays/color"
	"github.com/sendelivery/go-trace-rays/interval"
	"github.com/sendelivery/go-trace-rays/object/hittable"
	"github.com/sendelivery/go-trace-rays/ray"
	"github.com/sendelivery/go-trace-rays/utility"
	"github.com/sendelivery/go-trace-rays/vec3"
)

type Camera struct {
	SamplesPerPixel int     // Count of random samples for each pixel
	AspectRatio     float64 // Ratio of image width over height
	ImageWidth      int     // Rendered image width in pixel count
	MaxDepth        int     // Maximum number of ray bounces into the scene

	imageHeight      int          // Rendered image height
	centre           vec3.Vector3 // Camera center
	pixel00Loc       vec3.Vector3 // Location of pixel 0, 0
	pixelDeltaU      vec3.Vector3 // Offset to pixel to the right
	pixelDeltaV      vec3.Vector3 // Offset to pixel below
	pixelSampleScale float64      // Color scale factor for a sum of pixel samples
}

func New() *Camera {
	c := Camera{
		AspectRatio:     1.0,
		ImageWidth:      100,
		SamplesPerPixel: 1,
		MaxDepth:        1,
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
			col := color.New(0, 0, 0)
			for range c.SamplesPerPixel {
				r := c.getRay(i, j)
				col.Add(c.rayColor(r, c.MaxDepth, world))
			}
			col.Mulf(c.pixelSampleScale)
			color.WriteColor(os.Stdout, col)
		}
	}

	elapsed := time.Since(start).Milliseconds()

	fmt.Fprintf(os.Stderr, "\rDone in %dms.         \n", elapsed)
}

func (c *Camera) initialise() {
	// Calculate image height, ensuring it's at least 1
	c.imageHeight = max(int(float64(c.ImageWidth)/c.AspectRatio), 1)

	c.pixelSampleScale = 1.0 / float64(c.SamplesPerPixel)

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

// getRay construct a camera ray originating from the origin and directed at a randomly sampled
// point around the pixel location i, j
func (c *Camera) getRay(i, j int) ray.Ray {
	offset := c.sampleSquare()
	pixelSample := vec3.Add(
		c.pixel00Loc,
		vec3.Add(
			vec3.Mulf(c.pixelDeltaU, float64(i)+offset.X()),
			vec3.Mulf(c.pixelDeltaV, float64(j)+offset.Y()),
		),
	)
	rayDirection := vec3.Sub(pixelSample, c.centre)
	return ray.New(c.centre, rayDirection)
}

// sampleSquare returns the vector to a random point in the [-.5,-.5]-[+.5,+.5] unit square
func (c *Camera) sampleSquare() vec3.Vector3 {
	return vec3.New(utility.Random()-0.5, utility.Random()-0.5, 0)
}

const dampen = 0.5

func (c *Camera) rayColor(r ray.Ray, depth int, world hittable.Hittabler) color.Color {
	if depth <= 0 {
		return color.Black
	}

	if hr, ok := world.Hit(r, interval.New(1e-3, math.Inf(1))); ok {
		if attenuation, scattered, ok := hr.Material().Scatter(r, hr); ok {
			return vec3.Mulv(attenuation, c.rayColor(scattered, depth-1, world))
		}
		return color.Black
	}

	unitDirection := vec3.UnitVector(r.Direction())
	a := dampen * (unitDirection.Y() + 1)
	blue := color.New(0.5, 0.7, 1)

	return vec3.Add(
		vec3.Mulf(color.White, (1.0-a)),
		vec3.Mulf(blue, a),
	)
}
