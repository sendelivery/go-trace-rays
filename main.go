package main

import (
	"github.com/sendelivery/go-trace-rays/camera"
	"github.com/sendelivery/go-trace-rays/hittable"
	"github.com/sendelivery/go-trace-rays/hittable/sphere"
	"github.com/sendelivery/go-trace-rays/vec3"
)

func main() {
	var world hittable.HittableList
	s1 := sphere.New(vec3.New(0, 0, -1), 0.5)
	s2 := sphere.New(vec3.New(0, -100.5, -1), 100)
	world.Add(s1, s2)

	cam := camera.New()
	cam.AspectRatio = 16.0 / 9.0
	cam.ImageWidth = 400
	cam.SamplesPerPixel = 100
	cam.MaxDepth = 50

	cam.Render(world)
}
