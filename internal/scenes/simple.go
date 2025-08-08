package scenes

import (
	"github.com/sendelivery/go-trace-rays/internal/color"
	"github.com/sendelivery/go-trace-rays/internal/object/hittable"
	"github.com/sendelivery/go-trace-rays/internal/object/material"
	"github.com/sendelivery/go-trace-rays/internal/object/sphere"
	"github.com/sendelivery/go-trace-rays/internal/vec3"
)

func NewSimple() hittable.Hittabler {
	var world hittable.HittableList

	materialGround := material.NewLambertian(color.New(0.8, 0.8, 0))
	materialCentre := material.NewLambertian(color.New(0.1, 0.2, 0.5))
	materialLeft := material.NewDielectric(1.5)
	materialBubble := material.NewDielectric(1 / 1.5)
	materialRight := material.NewMetal(color.New(0.8, 0.6, 0.2), 0.1)

	ground := sphere.New(vec3.New(0, -100.5, -1), 100, materialGround)
	centre := sphere.New(vec3.New(4, 0, 1), 0.5, materialCentre)
	left := sphere.New(vec3.New(3, 0, 2), 0.5, materialLeft)
	bubble := sphere.New(vec3.New(3, 0, 2), 0.4, materialBubble)
	right := sphere.New(vec3.New(3, 0, -0.5), 0.5, materialRight)

	world.Add(ground, centre, left, bubble, right)

	return world
}
