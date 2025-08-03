package sphere

import (
	"math"

	"github.com/sendelivery/go-trace-rays/interval"
	"github.com/sendelivery/go-trace-rays/object/hittable"
	"github.com/sendelivery/go-trace-rays/object/material"
	"github.com/sendelivery/go-trace-rays/ray"
	"github.com/sendelivery/go-trace-rays/vec3"
)

type Sphere struct {
	centre vec3.Vector3
	radius float64
	mat    material.Scatterer
}

func New(centre vec3.Vector3, radius float64, mat material.Scatterer) Sphere {
	return Sphere{centre, math.Max(0, radius), mat}
}

func (s Sphere) Hit(r ray.Ray, rt interval.Interval) (hittable.HitRecord, bool) {
	oc := vec3.Sub(s.centre, r.Origin())
	a := r.Direction().LengthSquared()
	h := vec3.Dot(r.Direction(), oc)
	c := oc.LengthSquared() - s.radius*s.radius

	discriminant := h*h - a*c
	if discriminant < 0 {
		return hittable.HitRecord{}, false
	}

	sqrtd := math.Sqrt(discriminant)

	// Find the nearest root that lies in the acceptable range
	root := (h - sqrtd) / a
	if !rt.Surrounds(root) {
		root = (h + sqrtd) / a
		if !rt.Surrounds(root) {
			return hittable.HitRecord{}, false
		}
	}
	outwardNormal := vec3.Div(vec3.Sub(r.At(root), s.centre), s.radius)
	hr := hittable.NewHitRecord(r, root, outwardNormal, s.mat)

	return hr, true
}
