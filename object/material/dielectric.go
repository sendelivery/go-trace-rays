package material

import (
	"math"

	"github.com/sendelivery/go-trace-rays/color"
	"github.com/sendelivery/go-trace-rays/object/hitrecord"
	"github.com/sendelivery/go-trace-rays/ray"
	"github.com/sendelivery/go-trace-rays/utility"
	"github.com/sendelivery/go-trace-rays/vec3"
)

type Dielectric struct {
	refractionIndex float64
}

func NewDielectric(refractionIndex float64) Dielectric {
	return Dielectric{
		refractionIndex: refractionIndex,
	}
}

func (d *Dielectric) Scatter(in ray.Ray, hr hitrecord.HitRecord) (color.Color, ray.Ray, bool) {
	ri := d.refractionIndex
	if hr.FrontFace() {
		ri = 1.0 / d.refractionIndex
	}

	unitDir := vec3.UnitVector(in.Direction())

	cosTheta := min(vec3.Dot(vec3.Mulf(unitDir, -1), hr.Normal()), 1)
	sinTheta := math.Sqrt(1 - cosTheta*cosTheta)

	cannotRefract := ri*sinTheta > 1

	var direction vec3.Vector3
	if cannotRefract || reflectance(cosTheta, ri) > utility.Random() {
		direction = vec3.Reflect(unitDir, hr.Normal())
	} else {
		direction = vec3.Refract(unitDir, hr.Normal(), ri)
	}

	scattered := ray.New(hr.Point(), direction)
	return color.White, scattered, true
}

// reflectance uses Schlick's approximation of reflectance.
func reflectance(cosine, refractionIndex float64) float64 {
	r0 := (1 - refractionIndex) / (1 + refractionIndex)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow((1-cosine), 5)
}
