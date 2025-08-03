package hittable

import (
	"github.com/sendelivery/go-trace-rays/interval"
	"github.com/sendelivery/go-trace-rays/object/material"
	"github.com/sendelivery/go-trace-rays/ray"
	"github.com/sendelivery/go-trace-rays/vec3"
)

type Hittabler interface {
	// Hit takes in a Ray and a range, returning a pointer to a HitRecord and a bool
	// indicating if the object was hit by the ray. The pointer will be nil and the bool
	// false if it was not hit.
	Hit(r ray.Ray, rt interval.Interval) (HitRecord, bool)
}

type HitRecord struct {
	point, normal vec3.Vector3
	t             float64
	frontFace     bool
	mat           material.Scatterer
}

func NewHitRecord(r ray.Ray, t float64, outwardNormal vec3.Vector3, mat material.Scatterer) HitRecord {
	hr := HitRecord{
		t:     t,
		point: r.At(t),
		mat:   mat,
	}
	hr.setFaceNormal(r, outwardNormal)
	return hr
}

// setFaceNormal sets the HitRecord's normal vector.
// The outwardNormal argument is assumed to have unit length.
func (hr *HitRecord) setFaceNormal(r ray.Ray, outwardNormal vec3.Vector3) {
	hr.frontFace = vec3.Dot(r.Direction(), outwardNormal) < 0
	if hr.frontFace {
		hr.normal = vec3.Duplicate(outwardNormal)
	} else {
		hr.normal = vec3.Mulf(outwardNormal, -1)
	}
}

func (hr *HitRecord) Point() vec3.Vector3          { return hr.point }
func (hr *HitRecord) Normal() vec3.Vector3         { return hr.normal }
func (hr *HitRecord) T() float64                   { return hr.t }
func (hr *HitRecord) FrontFace() bool              { return hr.frontFace }
func (hr *HitRecord) Material() material.Scatterer { return hr.mat }

type HittableList struct {
	objects []Hittabler
}

func (hl *HittableList) Clear() {
	hl.objects = make([]Hittabler, 0)
}

func (hl *HittableList) Add(o ...Hittabler) {
	hl.objects = append(hl.objects, o...)
}

func (hl HittableList) Hit(r ray.Ray, rt interval.Interval) (HitRecord, bool) {
	var result HitRecord
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
