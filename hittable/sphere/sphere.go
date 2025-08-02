package sphere

import (
	"math"

	"github.com/sendelivery/go-trace-rays/hittable"
	"github.com/sendelivery/go-trace-rays/intervals"
	"github.com/sendelivery/go-trace-rays/ray"
	"github.com/sendelivery/go-trace-rays/vec3"
)

type Sphere struct {
	centre vec3.Vector3
	radius float64
}

func New(centre vec3.Vector3, radius float64) Sphere {
	return Sphere{centre, math.Max(0, radius)}
}

func (s *Sphere) Hit(r ray.Ray, rt intervals.Interval) (hr *hittable.HitRecord, ok bool) {
	oc := vec3.Sub(s.centre, r.Origin())
	a := r.Direction().LengthSquared()
	h := vec3.Dot(r.Direction(), oc)
	c := oc.LengthSquared() - s.radius*s.radius

	discriminant := h*h - a*c
	if discriminant < 0 {
		return nil, false
	}

	sqrtd := math.Sqrt(discriminant)

	// Find the nearest root that lies in the acceptable range
	root := (h - sqrtd) / a
	if !rt.Surrounds(root) {
		root = (h + sqrtd) / a
		if !rt.Surrounds(root) {
			return nil, false
		}
	}

	hr = &hittable.HitRecord{
		T:     root,
		Point: r.At(root),
	}
	outwardNormal := vec3.Div(vec3.Sub(hr.Point, s.centre), s.radius)
	hr.SetFaceNormal(r, outwardNormal)

	return hr, true
}
