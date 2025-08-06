package camera

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/sendelivery/go-trace-rays/internal/color"
	"github.com/sendelivery/go-trace-rays/internal/interval"
	"github.com/sendelivery/go-trace-rays/internal/object/hittable"
	"github.com/sendelivery/go-trace-rays/internal/ray"
	"github.com/sendelivery/go-trace-rays/internal/utility"
	"github.com/sendelivery/go-trace-rays/internal/vec3"
)

type Camera struct {
	SamplesPerPixel int     // Count of random samples for each pixel
	AspectRatio     float64 // Ratio of image width over height
	ImageWidth      int     // Rendered image width in pixel count
	MaxDepth        int     // Maximum number of ray bounces into the scene

	VerticalFov float64      // Vertical view angle (field of view)
	LookFrom    vec3.Vector3 // The point the camera is looking from
	LookAt      vec3.Vector3 // The point the camera is looking at
	VUp         vec3.Vector3 // Camera-relative up direction

	DefocusAngle  float64 // Variation angle of rays through each pixel
	FocusDistance float64 // Distance from the camera look from point to the plane of perfect focus

	imageHeight      int          // Rendered image height
	centre           vec3.Vector3 // Camera center
	pixel00Loc       vec3.Vector3 // Location of pixel 0, 0
	pixelDeltaU      vec3.Vector3 // Offset to pixel to the right
	pixelDeltaV      vec3.Vector3 // Offset to pixel below
	pixelSampleScale float64      // Color scale factor for a sum of pixel samples
	u, v, w          vec3.Vector3 // Camera frame basis vectors
	defocusDiskU     vec3.Vector3 // Defocus disk horizontal radius
	defocusDiskV     vec3.Vector3 // Defocus disk vertical radius
}

func New() *Camera {
	c := Camera{
		AspectRatio:     1.0,
		ImageWidth:      100,
		SamplesPerPixel: 10,
		MaxDepth:        10,
		VerticalFov:     90,
		LookFrom:        vec3.New(0, 0, 0),
		LookAt:          vec3.New(0, 0, -1),
		VUp:             vec3.New(0, 1, 0),
		DefocusAngle:    0,
		FocusDistance:   10,
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

	c.centre = c.LookFrom

	// Determine viewport dimensions
	theta := utility.Deg2Rad(c.VerticalFov)
	h := math.Tan(theta / 2)
	viewportHeight := 2.0 * h * c.FocusDistance
	viewportWidth := viewportHeight * float64(c.ImageWidth) / float64(c.imageHeight)

	// Calculate u, v, w unit basis vectors for the camera coordinate frame.
	c.w = vec3.UnitVector(vec3.Sub(c.LookFrom, c.LookAt))
	c.u = vec3.UnitVector(vec3.Cross(c.VUp, c.w))
	c.v = vec3.Cross(c.w, c.u)

	// Calculate the vectors across the horizontal and down the vertical viewport edges
	viewportU := vec3.Mulf(c.u, viewportWidth)                 // Vector across the viewport horizontal edge
	viewportV := vec3.Mulf(vec3.Mulf(c.v, -1), viewportHeight) // Vector down viewport vertical edge

	// Calculate the horizontal and vertical delta vectors from pixel to pixel
	c.pixelDeltaU = vec3.Div(viewportU, float64(c.ImageWidth))
	c.pixelDeltaV = vec3.Div(viewportV, float64(c.imageHeight))

	// Calculate the location of the upper left pixel
	viewportUpperLeft := vec3.Duplicate(c.centre)
	viewportUpperLeft.Sub(
		vec3.Mulf(c.w, c.FocusDistance),
	).Sub(
		vec3.Div(viewportU, 2),
	).Sub(
		vec3.Div(viewportV, 2),
	)
	c.pixel00Loc = vec3.Add(
		viewportUpperLeft,
		vec3.Mulf(vec3.Add(c.pixelDeltaU, c.pixelDeltaV), 0.5),
	) // Inset pixels

	// Calculate the camera defocus disk basis vectors
	defocusRadius := c.FocusDistance * math.Tan(utility.Deg2Rad(c.DefocusAngle/2))
	c.defocusDiskU = vec3.Mulf(c.u, defocusRadius)
	c.defocusDiskV = vec3.Mulf(c.v, defocusRadius)
}

// getRay construct a camera ray originating from the defocus disk and directed at a randomly
// sampled point around the pixel location i, j
func (c *Camera) getRay(i, j int) ray.Ray {
	offset := c.sampleSquare()
	pixelSample := vec3.Add(
		c.pixel00Loc,
		vec3.Add(
			vec3.Mulf(c.pixelDeltaU, float64(i)+offset.X()),
			vec3.Mulf(c.pixelDeltaV, float64(j)+offset.Y()),
		),
	)
	rayOrigin := c.centre
	if c.DefocusAngle > 0 {
		rayOrigin = c.defocusDiskSample()
	}
	rayDirection := vec3.Sub(pixelSample, rayOrigin)
	return ray.New(rayOrigin, rayDirection)
}

// sampleSquare returns the vector to a random point in the [-.5,-.5]-[+.5,+.5] unit square
func (c *Camera) sampleSquare() vec3.Vector3 {
	return vec3.New(utility.Random()-0.5, utility.Random()-0.5, 0)
}

// defocusDiskSample returns a random point in the camera defocus disk
func (c *Camera) defocusDiskSample() vec3.Vector3 {
	p := vec3.RandomInUnitDisk()
	x := vec3.Duplicate(c.centre)
	x.Add(vec3.Mulf(c.defocusDiskU, p.X()))
	x.Add(vec3.Mulf(c.defocusDiskV, p.Y()))
	return x
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
