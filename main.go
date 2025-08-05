package main

import (
	"github.com/sendelivery/go-trace-rays/camera"
	"github.com/sendelivery/go-trace-rays/color"
	"github.com/sendelivery/go-trace-rays/object/hittable"
	"github.com/sendelivery/go-trace-rays/object/material"
	"github.com/sendelivery/go-trace-rays/object/sphere"
	"github.com/sendelivery/go-trace-rays/vec3"
)

func main() {
	var world hittable.HittableList

	materialGround := material.NewLambertian(color.New(0.8, 0.8, 0))
	materialCentre := material.NewLambertian(color.New(0.1, 0.2, 0.5))
	materialLeft := material.NewDielectric(1.5)
	materialBubble := material.NewDielectric(1 / 1.5)
	materialRight := material.NewMetal(color.New(0.8, 0.6, 0.2), 1)

	ground := sphere.New(vec3.New(0, -100.5, -1), 100, &materialGround)
	centre := sphere.New(vec3.New(0, 0, -1.2), 0.5, &materialCentre)
	left := sphere.New(vec3.New(-1, 0, -1), 0.5, &materialLeft)
	bubble := sphere.New(vec3.New(-1, 0, -1), 0.4, &materialBubble)
	right := sphere.New(vec3.New(1, 0, -1), 0.5, &materialRight)

	world.Add(ground, centre, left, bubble, right)

	cam := camera.New()
	cam.AspectRatio = 16.0 / 9.0
	cam.ImageWidth = 400
	cam.SamplesPerPixel = 100
	cam.MaxDepth = 50

	cam.VerticalFov = 20
	cam.LookFrom = vec3.New(-2, 2, 1)
	cam.LookAt = vec3.New(0, 0, -1)
	cam.VUp = vec3.New(0, 1, 0)

	cam.DefocusAngle = 10
	cam.FocusDistance = 3.4

	cam.Render(world)
}
