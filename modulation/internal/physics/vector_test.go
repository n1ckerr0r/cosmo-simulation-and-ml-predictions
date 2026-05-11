package physics

import (
	"math"
	"testing"
)

func TestVectorOperations(t *testing.T) {
	v := Vector{X: 3, Y: 4, Z: 12}
	o := Vector{X: 1, Y: -2, Z: 3}

	if got := v.Add(o); got != (Vector{X: 4, Y: 2, Z: 15}) {
		t.Fatalf("Add() = %+v", got)
	}
	if got := v.Sub(o); got != (Vector{X: 2, Y: 6, Z: 9}) {
		t.Fatalf("Sub() = %+v", got)
	}
	if got := o.Mul(2); got != (Vector{X: 2, Y: -4, Z: 6}) {
		t.Fatalf("Mul() = %+v", got)
	}
	if got := v.Length(); math.Abs(got-13) > 1e-12 {
		t.Fatalf("Length() = %v", got)
	}
}
