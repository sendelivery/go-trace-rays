package material

import (
	"github.com/sendelivery/go-trace-rays/color"
	"github.com/sendelivery/go-trace-rays/object/hitrecord"
	"github.com/sendelivery/go-trace-rays/ray"
	"github.com/sendelivery/go-trace-rays/vec3"
)

type Lambertian struct {
	albedo color.Color
}

func NewLambertian(albedo color.Color) Lambertian {
	return Lambertian{
		albedo: albedo,
	}
}

func (l *Lambertian) Scatter(in ray.Ray, hr hitrecord.HitRecord) (color.Color, ray.Ray, bool) {
	scatterDir := vec3.Add(hr.Normal(), vec3.NewRandomUnitVector())

	// Catch a bad scatter direction (near zero)
	if vec3.IsNearZero(scatterDir) {
		scatterDir = hr.Normal()
	}

	scattered := ray.New(hr.Point(), scatterDir)
	return l.albedo, scattered, true
}
