package vec3

import (
	"fmt"
	"math"

	"github.com/sendelivery/go-trace-rays/utility"
)

type Vector3 struct {
	x, y, z float64
}

func New(x, y, z float64) Vector3 {
	return Vector3{x, y, z}
}

func NewRandom() Vector3 {
	return Vector3{
		utility.Random(),
		utility.Random(),
		utility.Random(),
	}
}

func NewRandomN(min, max float64) Vector3 {
	return Vector3{
		utility.RandomN(min, max),
		utility.RandomN(min, max),
		utility.RandomN(min, max),
	}
}

func NewRandomUnitVector() Vector3 {
	for {
		p := NewRandomN(-1, 1)
		lensq := p.LengthSquared()
		if 1e-160 < lensq && lensq <= 1 {
			return Div(p, math.Sqrt(lensq))
		}
	}
}

func NewRandomOnHemisphere(normal Vector3) Vector3 {
	u := NewRandomUnitVector()
	if Dot(u, normal) > 0.0 {
		// In the same hemisphere as the normal
		return u
	}
	return Mulf(u, -1) // Invert the vector
}

func Duplicate(v Vector3) Vector3 {
	return v
}

// Vector3 Struct Methods

func (v *Vector3) X() float64 { return v.x }
func (v *Vector3) Y() float64 { return v.y }
func (v *Vector3) Z() float64 { return v.z }

// String returns a string representation of a Vector3
func (v *Vector3) String() string {
	return fmt.Sprintf("%f %f %f", v.x, v.y, v.z)
}

// Add mutates the source Vector3 by adding v to it.
// A reference to the source Vector3 is returned for chaining.
func (vec *Vector3) Add(v Vector3) *Vector3 {
	add(vec, &v)
	return vec
}

// Sub mutates the source Vector3 by subtracting v from it.
// A reference to the source Vector3 is returned for chaining.
func (vec *Vector3) Sub(v Vector3) *Vector3 {
	sub(vec, &v)
	return vec
}

// Mulv mutates the source Vector3 by multiplying it with v.
// A reference to the source Vector3 is returned for chaining.
func (vec *Vector3) Mulv(v Vector3) *Vector3 {
	mulv(vec, &v)
	return vec
}

// Mulf mutates the source Vector3 by multiplying it by f.
// A reference to the source Vector3 is returned for chaining.
func (vec *Vector3) Mulf(f float64) *Vector3 {
	mulf(vec, f)
	return vec
}

// Div mutates the source Vector3 by dividing it by f.
// A reference to the source Vector3 is returned for chaining.
func (vec *Vector3) Div(f float64) *Vector3 {
	div(vec, f)
	return vec
}

// Dot returns the dot product of the source Vector3 and v.
func (vec *Vector3) Dot(v Vector3) float64 {
	return dot(vec, &v)
}

// Cross mutate the source Vector3 with the cross product of itself and v.
func (vec *Vector3) Cross(v Vector3) {
	cross(vec, &v)
}

// Equal
func Equal(v1, v2 Vector3) bool {
	return equal(v1, v2)
}

// Length returns the length of the Vector3
func (vec Vector3) Length() float64 {
	return length(&vec)
}

// LengthSquared returns the squared length of the Vector3
func (vec Vector3) LengthSquared() float64 {
	return lengthSquared(&vec)
}

// Package Functions

// Add returns v1 plus v2 as a new Vector3.
func Add(v1, v2 Vector3) Vector3 {
	add(&v1, &v2)
	return v1
}

// Sub returns v1 minus v2 as a new Vector3.
func Sub(v1, v2 Vector3) Vector3 {
	sub(&v1, &v2)
	return v1
}

// Mulv returns v1 times v2 as a new Vector3.
func Mulv(v1, v2 Vector3) Vector3 {
	mulv(&v1, &v2)
	return v1
}

// Mulf returns v times f as a new Vector3.
func Mulf(v Vector3, f float64) Vector3 {
	mulf(&v, f)
	return v
}

// Div returns v divided by f as a new Vector3.
func Div(v Vector3, f float64) Vector3 {
	div(&v, f)
	return v
}

// Dot returns the dot product of v1 and v2.
func Dot(v1, v2 Vector3) float64 {
	return dot(&v1, &v2)
}

// Cross returns the cross product of v1 and v2 as a new Vector3.
func Cross(v1, v2 Vector3) Vector3 {
	cross(&v1, &v2)
	return v1
}

// UnitVector returns the unit vector of v as a new Vector3.
func UnitVector(v Vector3) Vector3 {
	unitVector(&v)
	return v
}

// IsNearZero returns true if the Vector3 v is near zero in all dimensions.
func IsNearZero(v Vector3) bool {
	s := 1e-8
	return math.Abs(v.x) < s && math.Abs(v.y) < s && math.Abs(v.z) < s
}

func Reflect(v, normal Vector3) Vector3 {
	d := Dot(v, normal) * 2
	return Sub(v, Mulf(normal, d))
}

// Vector utility functions

// add takes in two Vector3 structs for addition.
// a is mutated by adding b, nothing is returned.
func add(a, b *Vector3) {
	a.x += b.x
	a.y += b.y
	a.z += b.z
}

// sub takes in two Vector3 structs for subtraction.
// a is mutated by subtracting b, nothing is returned.
func sub(a, b *Vector3) {
	a.x -= b.x
	a.y -= b.y
	a.z -= b.z
}

// mulv takes in two Vector3 structs for multiplication.
// a is mutated by multiplying it with b, nothing is returned.
func mulv(a, b *Vector3) {
	a.x *= b.x
	a.y *= b.y
	a.z *= b.z
}

// mulf takes in a Vector3 and a float64 for multiplication.
// a is mutated by multiplying it with f, nothing is returned.
func mulf(a *Vector3, f float64) {
	a.x *= f
	a.y *= f
	a.z *= f
}

// div takes in a Vector3 and a float64 for division.
// a is mutated by dividing it by f, nothing is returned.
func div(a *Vector3, f float64) {
	a.x /= f
	a.y /= f
	a.z /= f
}

// dot takes in two Vector3 pointers and returns their dot product.
func dot(a, b *Vector3) float64 {
	return a.x*b.x +
		a.y*b.y +
		a.z*b.z
}

// cross takes in two Vector3 structs.
// a is mutated by calculating the cross product with b, nothing is returned.
func cross(a, b *Vector3) {
	v := New(a.x, a.y, a.z)
	a.x = v.y*b.z - v.z*b.y
	a.y = v.z*b.x - v.x*b.z
	a.z = v.x*b.y - v.y*b.x
}

// unitVector takes in a Vector3 and normalises it.
// The Vector3 is mutated, nothing is returned.
func unitVector(a *Vector3) {
	div(a, length(a))
}

// equal returns true if the x, y, and z values of two Vector3s are equal.
func equal(a, b Vector3) bool {
	if a.x != b.x {
		return false
	}
	if a.y != b.y {
		return false
	}
	if a.z != b.z {
		return false
	}
	return true
}

// lengthSquared returns the squared length of a Vector3
func lengthSquared(a *Vector3) float64 {
	return a.x*a.x + a.y*a.y + a.z*a.z
}

// length returns the length of a Vector3
func length(a *Vector3) float64 {
	return math.Sqrt(lengthSquared(a))
}
