package scenes

import (
	"github.com/sendelivery/go-trace-rays/internal/color"
	"github.com/sendelivery/go-trace-rays/internal/object/hitrecord"
	"github.com/sendelivery/go-trace-rays/internal/object/hittable"
	"github.com/sendelivery/go-trace-rays/internal/object/material"
	"github.com/sendelivery/go-trace-rays/internal/object/sphere"
	"github.com/sendelivery/go-trace-rays/internal/utility"
	"github.com/sendelivery/go-trace-rays/internal/vec3"
)

func NewComplex() hittable.Hittabler {
	var world hittable.HittableList

	groundMaterial := material.NewLambertian(color.New(0.5, 0.5, 0.5))
	groundSphere := sphere.New(vec3.New(0, -1000, 0), 1000, groundMaterial)

	world.Add(groundSphere)

	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			chooseMat := utility.Random()
			centre := vec3.New(float64(a)+0.9*utility.Random(), 0.2, float64(b)+0.9*utility.Random())

			if vec3.Sub(centre, vec3.New(4, 0.2, 0)).Length() > 0.9 {
				var sphereMat hitrecord.Scatterer

				if chooseMat < 0.8 {
					// diffuse
					albedo := color.NewRandom(0, 1)
					sphereMat = material.NewLambertian(albedo)
				} else if chooseMat < 0.95 {
					// metal
					albedo := color.NewRandom(0.5, 1)
					fuzz := utility.RandomN(0, 0.5)
					sphereMat = material.NewMetal(albedo, fuzz)
				} else {
					// glass
					sphereMat = material.NewDielectric(1.5)
				}

				world.Add(sphere.New(centre, 0.2, sphereMat))
			}
		}
	}

	mat1 := material.NewDielectric(1.5)
	mat2 := material.NewLambertian(color.New(0.4, 0.2, 0.1))
	mat3 := material.NewMetal(color.New(0.7, 0.6, 0.5), 0.0)

	world.Add(
		sphere.New(vec3.New(0, 1, 0), 1, mat1),
		sphere.New(vec3.New(-4, 1, 0), 1, mat2),
		sphere.New(vec3.New(4, 1, 0), 1, mat3),
	)

	return world
}
