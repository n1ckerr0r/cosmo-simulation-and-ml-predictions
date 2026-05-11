package physics

import (
	"math"
	"testing"
)

func TestCalculateAccelerationsTwoBodies(t *testing.T) {
	bodies := []Body{
		{Mass: 10, Position: Vector{X: 0, Y: 0, Z: 0}},
		{Mass: 20, Position: Vector{X: 2, Y: 0, Z: 0}},
	}

	acc := CalculateAccelerations(bodies)

	dist := 2.0 + 1e-9
	wantFirst := G * bodies[1].Mass * 2 / (dist * dist * dist)
	wantSecond := -G * bodies[0].Mass * 2 / (dist * dist * dist)

	if math.Abs(acc[0].X-wantFirst) > 1e-20 || acc[0].Y != 0 || acc[0].Z != 0 {
		t.Fatalf("first acceleration = %+v, want x=%v", acc[0], wantFirst)
	}
	if math.Abs(acc[1].X-wantSecond) > 1e-20 || acc[1].Y != 0 || acc[1].Z != 0 {
		t.Fatalf("second acceleration = %+v, want x=%v", acc[1], wantSecond)
	}
}
