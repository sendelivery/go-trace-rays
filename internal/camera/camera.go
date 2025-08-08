package camera

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sendelivery/go-trace-rays/internal/color"
	"github.com/sendelivery/go-trace-rays/internal/image"
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

	// Below fields are used by the parallel workflow
	parallel  bool        // Whether to render the image using the parallel workflow
	workers   int         // Number of goroutines rendering the image
	img       image.Image // An Image struct used by the parallel workflow
	numChunks int         // Number of chunks to break the image up in
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

	// Timer
	start := time.Now()
	defer func() {
		elapsed := time.Since(start).Seconds()
		c.statusf("\rDone in %.2fs.         \n", elapsed)
	}()

	fmt.Printf("P3\n%d %d\n255\n", c.ImageWidth, c.imageHeight)

	for j := range c.imageHeight {
		c.statusf("\rScanlines remaining: %d ", c.imageHeight-j)
		for i := range c.ImageWidth {
			col := c.calculatePixel(i, j, world)
			color.WriteColor(os.Stdout, col)
		}
	}
}

func (c *Camera) RenderParallel(world hittable.Hittabler) {
	c.parallel = true
	c.initialise()

	c.statusf("%d workers\n", c.workers)
	c.statusf("%d chunks\n", c.numChunks)

	// Timer
	start := time.Now()
	defer func() {
		elapsed := time.Since(start).Seconds()
		c.statusf("\rDone in %.2fs.         \n", elapsed)
	}()

	chunks := make(chan image.Chunk, c.numChunks)
	c.queueChunks(chunks)
	close(chunks) // All the data that needs to be sent has been sent

	var wg sync.WaitGroup
	wg.Add(c.numChunks)

	chunksLeft := atomic.Int32{}
	chunksLeft.Add(int32(c.numChunks))
	c.statusf("\rChunks remaining: %d ", c.numChunks)

	// Queue up the workers
	for range c.workers {
		go func() {
			for ch := range chunks {
				c.processChunk(ch, world)
				chunksLeft.Add(-1)
				c.statusf("\rChunks remaining: %d ", chunksLeft.Load())
				wg.Done()
			}
		}()
	}

	wg.Wait()

	// Draw the image
	fmt.Printf("P3\n%d %d\n255\n", c.ImageWidth, c.imageHeight)
	for y := range c.imageHeight {
		for x := range c.ImageWidth {
			col, ok := c.img.Get(image.NewPixelCoord(x, y))
			if !ok {
				panic(fmt.Sprintf("issue when drawing image at coord %d, %d", x, y))
			}
			color.WriteColor(os.Stdout, col)
		}
	}
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

	if !c.parallel {
		return
	}

	// If left unspecified, set concurrency to max CPU - 2 to leave headroom for the system
	if c.workers == 0 {
		c.workers = max(runtime.NumCPU()-2, 1)
	}

	pixelCount := c.ImageWidth * c.imageHeight
	c.numChunks = c.workers * c.workers

	for c.workers > 1 {
		if pixelCount/c.numChunks > 0 {
			break
		}

		// Special case, our image is so small that evenly dividing it up across the number
		// of workers would result in chunks less than 0 pixels in area. We'll reduce the
		// workers by 20% until we reach an acceptable number of workers or 1.
		c.workers = int(float64(c.workers) * 0.8)
		c.numChunks = c.workers * c.workers
	}

	c.img = image.New(c.ImageWidth, c.imageHeight)
}

func (c *Camera) calculatePixel(x, y int, world hittable.Hittabler) color.Color {
	col := color.New(0, 0, 0)
	for range c.SamplesPerPixel {
		r := c.getRay(x, y)
		col.Add(c.rayColor(r, c.MaxDepth, world))
	}
	return *col.Mulf(c.pixelSampleScale)
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

// queueChunks sends all the chunks to be computed to the ch channel
func (c *Camera) queueChunks(ch chan<- image.Chunk) {
	chunkDeltaU := c.ImageWidth / c.workers
	chunkDeltaV := c.imageHeight / c.workers

	startY := 0

	for j := range c.workers {
		endY := startY + chunkDeltaV
		if j == c.workers-1 {
			endY = c.imageHeight
		}

		startX := 0

		for i := range c.workers {
			endX := startX + chunkDeltaU
			if i == c.workers-1 {
				endX = c.ImageWidth
			}

			ch <- image.NewChunk(
				image.NewPixelCoord(startX, startY),
				image.NewPixelCoord(endX, endY),
			)

			startX = endX
		}
		startY = endY
	}
}

// processChunk calculates all the pixel colours for the given chunk and writes them to our
// camera's img
func (c *Camera) processChunk(chunk image.Chunk, world hittable.Hittabler) {
	for x := chunk.Start().X(); x < chunk.End().X(); x++ {
		for y := chunk.Start().Y(); y < chunk.End().Y(); y++ {
			col := c.calculatePixel(x, y, world)
			if err := c.img.Add(image.NewPixelCoord(x, y), col); err != nil {
				panic(err)
			}
		}
	}
}

func (c *Camera) statusf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
}
