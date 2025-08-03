package material

import (
	"github.com/sendelivery/go-trace-rays/color"
	"github.com/sendelivery/go-trace-rays/ray"
	"github.com/sendelivery/go-trace-rays/vec3"
)

type Scatterer interface {
	Scatter(
		in ray.Ray,
		hitPoint vec3.Vector3,
		hitNormal vec3.Vector3,
	) (color.Color, ray.Ray, bool)
}

type Material struct {
	albedo color.Color
	Scatterer
}
