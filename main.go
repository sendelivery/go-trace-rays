package main

import (
	"flag"

	"github.com/sendelivery/go-trace-rays/internal/camera"
	"github.com/sendelivery/go-trace-rays/internal/object/hittable"
	"github.com/sendelivery/go-trace-rays/internal/scenes"
	"github.com/sendelivery/go-trace-rays/internal/vec3"
)

func main() {
	complex := flag.Bool("complex", false, "whether to render a complex world")
	parallel := flag.Bool("parallel", false, "whether to use the parllelised rendering workflow")
	flag.Parse()

	var world hittable.Hittabler
	if *complex {
		world = scenes.NewComplex()
	} else {
		world = scenes.NewSimple()
	}

	cam := camera.New()
	cam.AspectRatio = 16.0 / 9.0
	cam.ImageWidth = 1200
	cam.SamplesPerPixel = 500
	cam.MaxDepth = 50

	cam.VerticalFov = 20
	cam.LookFrom = vec3.New(13, 2, 3)
	cam.LookAt = vec3.New(0, 0, 0)
	cam.VUp = vec3.New(0, 1, 0)

	cam.DefocusAngle = 0.6
	cam.FocusDistance = 10.0

	if *parallel {
		cam.RenderParallel(world)
	} else {
		cam.Render(world)
	}
}
