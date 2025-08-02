package vec3_test

import (
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sendelivery/go-trace-rays/vec3"
)

// Vector3 struct methods that mutate the struct
func TestVector3Methods(t *testing.T) {
	tests := []struct {
		method    string
		v1, v2    vec3.Vector3
		f         float64
		expected  vec3.Vector3
		expectedf float64
	}{
		{
			method:   "add",
			v1:       vec3.New(1, 0, 3),
			v2:       vec3.New(-1, 4, 2),
			expected: vec3.New(0, 4, 5),
		},
		{
			method:   "sub",
			v1:       vec3.New(1, 0, 3),
			v2:       vec3.New(-1, 4, 2),
			expected: vec3.New(2, -4, 1),
		},
		{
			method:   "mulv",
			v1:       vec3.New(1, 2, 3),
			v2:       vec3.New(1, 5, 7),
			expected: vec3.New(1, 10, 21),
		},
		{
			method:   "mulf",
			v1:       vec3.New(1, 2, 3),
			f:        5,
			expected: vec3.New(5, 10, 15),
		},
		{
			method:   "div",
			v1:       vec3.New(5, 10, 15),
			f:        5,
			expected: vec3.New(1, 2, 3),
		},
		{
			method:   "cross",
			v1:       vec3.New(1, 2, 3),
			v2:       vec3.New(1, 5, 7),
			expected: vec3.New(-1, -4, 3),
		},
	}

	for _, tc := range tests {
		t.Run(tc.method, func(t *testing.T) {
			t.Parallel()

			switch tc.method {
			case "add":
				tc.v1.Add(tc.v2)
			case "sub":
				tc.v1.Sub(tc.v2)
			case "mulv":
				tc.v1.Mulv(tc.v2)
			case "mulf":
				tc.v1.Mulf(tc.f)
			case "div":
				tc.v1.Div(tc.f)
			case "cross":
				tc.v1.Cross(tc.v2)
			default:
				t.Errorf("unsupported method, got=%s.", tc.method)
			}

			if !vec3.Equal(tc.v1, tc.expected) {
				t.Errorf("unexpected result, got=%q. want=%q.", &tc.v1, &tc.expected)
			}
		})
	}
}

// Vector3 struct methods that return a float
func TestVector3MethodsFloat(t *testing.T) {
	tests := []struct {
		method   string
		v1, v2   vec3.Vector3
		f        float64
		expected float64
	}{
		{
			method:   "dot",
			v1:       vec3.New(1, 0, 3),
			v2:       vec3.New(-1, 4, 2),
			expected: 5,
		},
		{
			method:   "length",
			v1:       vec3.New(1, 0, 3),
			expected: math.Sqrt(10),
		},
		{
			method:   "lengthSquared",
			v1:       vec3.New(1, 0, 3),
			expected: 10,
		},
	}

	for _, tc := range tests {
		t.Run(tc.method, func(t *testing.T) {
			t.Parallel()

			var result float64

			switch tc.method {
			case "dot":
				result = tc.v1.Dot(tc.v2)
			case "length":
				result = tc.v1.Length()
			case "lengthSquared":
				result = tc.v1.LengthSquared()
			default:
				t.Errorf("unsupported method, got=%s.", tc.method)
			}

			if result != tc.expected {
				t.Errorf("unexpected result, got=%f. want=%f.", result, tc.expected)
			}
		})
	}
}

func TestString(t *testing.T) {
	myVec := vec3.New(1, 2, 3)

	expected := "1.000000 2.000000 3.000000"

	if diff := cmp.Diff(myVec.String(), expected); diff != "" {
		t.Errorf("%s", diff)
	}
}

func TestEqual(t *testing.T) {
	myVec1 := vec3.New(1, 2, 3)
	myVec2 := vec3.New(1, 2, 3)

	if !vec3.Equal(myVec1, myVec2) {
		t.Error("expected vectors to be equal")
	}
}

// Package functions

func TestAdd(t *testing.T) {
	myVec1 := vec3.New(1, 2, 3)
	myVec2 := vec3.New(1, 2, 3)

	expected := vec3.New(2, 4, 6)

	newVec := vec3.Add(myVec1, myVec2)

	if !vec3.Equal(newVec, expected) {
		t.Errorf("unexpected result, want=%q. got=%q.", &newVec, &expected)
	}

	if !vec3.Equal(myVec1, vec3.New(1, 2, 3)) {
		t.Error("myVec1 was mutated")
	}

	if !vec3.Equal(myVec2, vec3.New(1, 2, 3)) {
		t.Error("myVec2 was mutated")
	}
}

func TestSub(t *testing.T) {
	myVec1 := vec3.New(3, 6, 9)
	myVec2 := vec3.New(1, 2, 3)

	expected := vec3.New(2, 4, 6)

	newVec := vec3.Sub(myVec1, myVec2)

	if !vec3.Equal(newVec, expected) {
		t.Errorf("unexpected result, want=%q. got=%q.", &newVec, &expected)
	}

	if !vec3.Equal(myVec1, vec3.New(3, 6, 9)) {
		t.Error("myVec1 was mutated")
	}

	if !vec3.Equal(myVec2, vec3.New(1, 2, 3)) {
		t.Error("myVec2 was mutated")
	}
}

func TestMulv(t *testing.T) {
	myVec1 := vec3.New(1, 2, 3)
	myVec2 := vec3.New(1, 5, 7)

	expected := vec3.New(1, 10, 21)

	newVec := vec3.Mulv(myVec1, myVec2)

	if !vec3.Equal(newVec, expected) {
		t.Errorf("unexpected result, want=%q. got=%q.", &newVec, &expected)
	}

	if !vec3.Equal(myVec1, vec3.New(1, 2, 3)) {
		t.Error("myVec1 was mutated")
	}

	if !vec3.Equal(myVec2, vec3.New(1, 5, 7)) {
		t.Error("myVec2 was mutated")
	}
}

func TestMulf(t *testing.T) {
	myVec1 := vec3.New(1, 2, 3)
	f := float64(5)

	expected := vec3.New(5, 10, 15)

	newVec := vec3.Mulf(myVec1, f)

	if !vec3.Equal(newVec, expected) {
		t.Errorf("unexpected result, want=%q. got=%q.", &newVec, &expected)
	}

	if !vec3.Equal(myVec1, vec3.New(1, 2, 3)) {
		t.Error("myVec1 was mutated")
	}
}

func TestDiv(t *testing.T) {
	myVec1 := vec3.New(5, 10, 15)
	f := float64(5)

	expected := vec3.New(1, 2, 3)

	newVec := vec3.Div(myVec1, f)

	if !vec3.Equal(newVec, expected) {
		t.Errorf("unexpected result, want=%q. got=%q.", &newVec, &expected)
	}

	if !vec3.Equal(myVec1, vec3.New(5, 10, 15)) {
		t.Error("myVec1 was mutated")
	}
}

func TestDot(t *testing.T) {
	myVec1 := vec3.New(1, 0, 3)
	myVec2 := vec3.New(-1, 4, 2)
	expected := float64(5)

	result := vec3.Dot(myVec1, myVec2)

	if result != expected {
		t.Errorf("unexpected result, want=%f. got=%f.", result, expected)
	}

	if !vec3.Equal(myVec1, vec3.New(1, 0, 3)) {
		t.Error("myVec1 was mutated")
	}

	if !vec3.Equal(myVec2, vec3.New(-1, 4, 2)) {
		t.Error("myVec2 was mutated")
	}
}

func TestCross(t *testing.T) {
	myVec1 := vec3.New(1, 2, 3)
	myVec2 := vec3.New(1, 5, 7)

	expected := vec3.New(-1, -4, 3)

	newVec := vec3.Cross(myVec1, myVec2)

	if !vec3.Equal(newVec, expected) {
		t.Errorf("unexpected result, want=%q. got=%q.", &newVec, &expected)
	}

	if !vec3.Equal(myVec1, vec3.New(1, 2, 3)) {
		t.Error("myVec1 was mutated")
	}

	if !vec3.Equal(myVec2, vec3.New(1, 5, 7)) {
		t.Error("myVec2 was mutated")
	}
}
