package main

import (
	"fmt"
	"os"

	"github.com/sendelivery/go-trace-rays/color"
	"github.com/sendelivery/go-trace-rays/ray"
	"github.com/sendelivery/go-trace-rays/vec3"
)

func rayColor(r ray.Ray) color.Color {
	dampen := 0.5

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
	imageWidth := 400

	// Calculate image height, ensuring it's at least 1
	imageHeight := int(float64(imageWidth) / aspectRatio)
	if imageHeight < 1 {
		imageHeight = 1
	}

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

			col := rayColor(r)
			color.WriteColor(os.Stdout, col)
		}
	}

	fmt.Fprintf(os.Stderr, "\rDone.                 \n")
}
