package main

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/sendelivery/go-trace-rays/color"
	"github.com/sendelivery/go-trace-rays/hittable"
	"github.com/sendelivery/go-trace-rays/hittable/sphere"
	"github.com/sendelivery/go-trace-rays/intervals"
	"github.com/sendelivery/go-trace-rays/ray"
	"github.com/sendelivery/go-trace-rays/vec3"
)

const dampen = 0.5

func rayColor(r ray.Ray, world hittable.Hittabler) color.Color {
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

func main() {
	// Image setup
	aspectRatio := 16.0 / 9.0
	imageWidth := 800

	// Calculate image height, ensuring it's at least 1
	imageHeight := int(float64(imageWidth) / aspectRatio)
	if imageHeight < 1 {
		imageHeight = 1
	}

	// World
	var world hittable.HittableList
	s1 := sphere.New(vec3.New(0, 0, -1), 0.5)
	s2 := sphere.New(vec3.New(0, -100.5, -1), 100)
	world.Add(&s1, &s2)

	// Camera setup
	focalLength := 1.0 // Distance between camera centre and viewport
	viewportHeight := 2.0
	viewportWidth := viewportHeight * float64(imageWidth) / float64(imageHeight)
	cameraCentre := vec3.New(0, 0, 0)

	// Calculate the vectors across the horizontal and down the vertical viewport edges
	viewportU := vec3.New(viewportWidth, 0, 0)
	viewportV := vec3.New(0, -viewportHeight, 0)

	// Calculate the horizontal and vertical delta vectors from pixel to pixel
	pixelDeltaU := vec3.Div(viewportU, float64(imageWidth))
	pixelDeltaV := vec3.Div(viewportV, float64(imageHeight))

	// Calculate the location of the upper left pixel
	viewportUpperLeft := vec3.Duplicate(cameraCentre)
	viewportUpperLeft.Sub(
		vec3.New(0, 0, focalLength),
	).Sub(
		vec3.Div(viewportU, 2),
	).Sub(
		vec3.Div(viewportV, 2),
	)
	pixel00Loc := vec3.Add(viewportUpperLeft, vec3.Mulf(vec3.Add(pixelDeltaU, pixelDeltaV), 0.5)) // Inset pixels

	// Render
	fmt.Printf("P3\n%d %d\n255\n", imageWidth, imageHeight)

	start := time.Now()

	for j := range imageHeight {
		fmt.Fprintf(os.Stderr, "\rScanlines remaining: %d ", imageHeight-j)
		for i := range imageWidth {
			pixelCentre := vec3.Add(
				pixel00Loc,
				vec3.Add(
					vec3.Mulf(pixelDeltaU, float64(i)),
					vec3.Mulf(pixelDeltaV, float64(j)),
				),
			)

			rayDirection := vec3.Sub(pixelCentre, cameraCentre)
			r := ray.New(cameraCentre, rayDirection)

			col := rayColor(r, &world)
			color.WriteColor(os.Stdout, col)
		}
	}

	elapsed := time.Since(start).Milliseconds()

	fmt.Fprintf(os.Stderr, "\rDone in %dms.         \n", elapsed)
}
