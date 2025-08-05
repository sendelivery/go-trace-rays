package material

import (
	"github.com/sendelivery/go-trace-rays/color"
	"github.com/sendelivery/go-trace-rays/object/hitrecord"
	"github.com/sendelivery/go-trace-rays/ray"
	"github.com/sendelivery/go-trace-rays/vec3"
)

type Metal struct {
	albedo color.Color
	fuzz   float64
}

func NewMetal(albedo color.Color, fuzz float64) Metal {
	return Metal{
		albedo: albedo,
		fuzz:   min(fuzz, 1),
	}
}

func (m *Metal) Scatter(in ray.Ray, hr hitrecord.HitRecord) (color.Color, ray.Ray, bool) {
	reflected := vec3.Reflect(in.Direction(), hr.Normal())
	reflected = vec3.UnitVector(reflected)
	reflected.Add(vec3.Mulf(vec3.NewRandomUnitVector(), m.fuzz))

	scattered := ray.New(hr.Point(), reflected)

	// scatter being false signals that we should absorb the ray,
	// meaning black should be used for this ray
	scatter := !(vec3.Dot(scattered.Direction(), hr.Normal()) <= 0)

	return m.albedo, scattered, scatter
}
