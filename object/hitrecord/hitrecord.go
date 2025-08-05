package hitrecord

import (
	"github.com/sendelivery/go-trace-rays/color"
	"github.com/sendelivery/go-trace-rays/ray"
	"github.com/sendelivery/go-trace-rays/vec3"
)

type Scatterer interface {
	Scatter(in ray.Ray, hr HitRecord) (color.Color, ray.Ray, bool)
}

type HitRecord struct {
	point, normal vec3.Vector3
	t             float64
	frontFace     bool
	mat           Scatterer
}

func New(r ray.Ray, t float64, outwardNormal vec3.Vector3, mat Scatterer) HitRecord {
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

func (hr *HitRecord) Point() vec3.Vector3  { return hr.point }
func (hr *HitRecord) Normal() vec3.Vector3 { return hr.normal }
func (hr *HitRecord) T() float64           { return hr.t }
func (hr *HitRecord) FrontFace() bool      { return hr.frontFace }
func (hr *HitRecord) Material() Scatterer  { return hr.mat }
