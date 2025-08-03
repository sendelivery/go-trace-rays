package material

import (
	"github.com/sendelivery/go-trace-rays/color"
	"github.com/sendelivery/go-trace-rays/ray"
	"github.com/sendelivery/go-trace-rays/vec3"
)

type Metal struct {
	Material
	fuzz float64
}

func NewMetal(albedo color.Color, fuzz float64) Metal {
	return Metal{
		Material: Material{albedo: albedo},
		fuzz:     min(fuzz, 1),
	}
}

func (m *Metal) Scatter(
	in ray.Ray,
	hitPoint vec3.Vector3,
	hitNormal vec3.Vector3,
) (color.Color, ray.Ray, bool) {
	reflected := vec3.Reflect(in.Direction(), hitNormal)
	reflected = vec3.UnitVector(reflected)
	reflected.Add(vec3.Mulf(vec3.NewRandomUnitVector(), m.fuzz))

	scattered := ray.New(hitPoint, reflected)

	// scatter being false signals that we should absorb the ray,
	// meaning black should be used for this ray
	scatter := !(vec3.Dot(scattered.Direction(), hitNormal) <= 0)

	return m.albedo, scattered, scatter
}
