package sphere

import (
	"math"

	"github.com/sendelivery/go-trace-rays/internal/interval"
	"github.com/sendelivery/go-trace-rays/internal/object/hitrecord"
	"github.com/sendelivery/go-trace-rays/internal/ray"
	"github.com/sendelivery/go-trace-rays/internal/vec3"
)

type Sphere struct {
	centre vec3.Vector3
	radius float64
	mat    hitrecord.Scatterer
}

func New(centre vec3.Vector3, radius float64, mat hitrecord.Scatterer) Sphere {
	return Sphere{centre, math.Max(0, radius), mat}
}

func (s Sphere) Hit(r ray.Ray, rt interval.Interval) (hitrecord.HitRecord, bool) {
	oc := vec3.Sub(s.centre, r.Origin())
	a := r.Direction().LengthSquared()
	h := vec3.Dot(r.Direction(), oc)
	c := oc.LengthSquared() - s.radius*s.radius

	discriminant := h*h - a*c
	if discriminant < 0 {
		return hitrecord.HitRecord{}, false
	}

	sqrtd := math.Sqrt(discriminant)

	// Find the nearest root that lies in the acceptable range
	root := (h - sqrtd) / a
	if !rt.Surrounds(root) {
		root = (h + sqrtd) / a
		if !rt.Surrounds(root) {
			return hitrecord.HitRecord{}, false
		}
	}
	outwardNormal := vec3.Div(vec3.Sub(r.At(root), s.centre), s.radius)
	hr := hitrecord.New(r, root, outwardNormal, s.mat)

	return hr, true
}
