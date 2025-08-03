package material

import (
	"github.com/sendelivery/go-trace-rays/color"
	"github.com/sendelivery/go-trace-rays/ray"
	"github.com/sendelivery/go-trace-rays/vec3"
)

type Lambertian Material

func NewLambertian(albedo color.Color) Lambertian {
	return Lambertian{albedo: albedo}
}

func (l *Lambertian) Scatter(
	in ray.Ray,
	hitPoint vec3.Vector3,
	hitNormal vec3.Vector3,
) (color.Color, ray.Ray, bool) {
	scatterDir := vec3.Add(hitNormal, vec3.NewRandomUnitVector())

	// Catch a bad scatter direction (near zero)
	if vec3.IsNearZero(scatterDir) {
		scatterDir = hitNormal
	}

	scattered := ray.New(hitPoint, scatterDir)
	return l.albedo, scattered, true
}
