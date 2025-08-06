package ray

import "github.com/sendelivery/go-trace-rays/internal/vec3"

type Ray struct {
	origin    vec3.Vector3
	direction vec3.Vector3
}

func New(origin, direction vec3.Vector3) Ray {
	return Ray{origin, direction}
}

func (r *Ray) At(t float64) vec3.Vector3 {
	return vec3.Add(r.origin, vec3.Mulf(r.direction, t))
}

func (r *Ray) Origin() vec3.Vector3 {
	return r.origin
}

func (r *Ray) Direction() vec3.Vector3 {
	return r.direction
}
