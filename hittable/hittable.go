package hittable

import (
	"github.com/sendelivery/go-trace-rays/intervals"
	"github.com/sendelivery/go-trace-rays/ray"
	"github.com/sendelivery/go-trace-rays/vec3"
)

type Hittabler interface {
	// Hit takes in a Ray and a range, returning a pointer to a HitRecord and a bool
	// indicating if the object was hit by the ray. The pointer will be nil and the bool
	// false if it was not hit.
	Hit(r ray.Ray, rt intervals.Interval) (HitRecord, bool)
}

type HitRecord struct {
	Point, normal vec3.Vector3
	T             float64
	FrontFace     bool
}

// SetFaceNormal sets the HitRecord's normal vector.
// The outwardNormal argument is assumed to have unit length.
func (hr *HitRecord) SetFaceNormal(r ray.Ray, outwardNormal vec3.Vector3) {
	hr.FrontFace = vec3.Dot(r.Direction(), outwardNormal) < 0
	if hr.FrontFace {
		hr.normal = vec3.Duplicate(outwardNormal)
	} else {
		hr.normal = vec3.Mulf(outwardNormal, -1)
	}
}

func (hr *HitRecord) Normal() vec3.Vector3 {
	return hr.normal
}

type HittableList struct {
	objects []Hittabler
}

func (hl *HittableList) Clear() {
	hl.objects = make([]Hittabler, 0)
}

func (hl *HittableList) Add(o ...Hittabler) {
	hl.objects = append(hl.objects, o...)
}

func (hl HittableList) Hit(r ray.Ray, rt intervals.Interval) (HitRecord, bool) {
	var result HitRecord
	var hitAnything bool
	closest := rt.Max

	for _, o := range hl.objects {
		if hr, ok := o.Hit(r, intervals.New(rt.Min, closest)); ok {
			hitAnything = true
			closest = hr.T
			result = hr
		}
	}

	return result, hitAnything
}
