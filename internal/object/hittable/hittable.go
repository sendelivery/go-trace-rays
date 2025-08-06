package hittable

import (
	"github.com/sendelivery/go-trace-rays/internal/interval"
	"github.com/sendelivery/go-trace-rays/internal/object/hitrecord"
	"github.com/sendelivery/go-trace-rays/internal/ray"
)

type Hittabler interface {
	// Hit takes in a Ray and a range, returning a pointer to a HitRecord and a bool
	// indicating if the object was hit by the ray. The pointer will be nil and the bool
	// false if it was not hit.
	Hit(r ray.Ray, rt interval.Interval) (hitrecord.HitRecord, bool)
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

func (hl HittableList) Hit(r ray.Ray, rt interval.Interval) (hitrecord.HitRecord, bool) {
	var result hitrecord.HitRecord
	var hitAnything bool
	closest := rt.Max

	for _, o := range hl.objects {
		if hr, ok := o.Hit(r, interval.New(rt.Min, closest)); ok {
			hitAnything = true
			closest = hr.T()
			result = hr
		}
	}

	return result, hitAnything
}
