package physics

import (
	"math"
	"testing"
)

func TestStepMovesBodyWithCurrentAcceleration(t *testing.T) {
	bodies := []Body{
		{
			Mass:         1,
			Position:     Vector{X: 0, Y: 0, Z: 0},
			Velocity:     Vector{X: 2, Y: 0, Z: 0},
			Acceleration: Vector{X: 1, Y: 0, Z: 0},
		},
	}

	next := Step(bodies, 2)

	if math.Abs(next[0].Position.X-6) > 1e-12 {
		t.Fatalf("position x = %v, want 6", next[0].Position.X)
	}
	if math.Abs(next[0].Velocity.X-2) > 1e-12 {
		t.Fatalf("velocity x = %v, want unchanged velocity without external force", next[0].Velocity.X)
	}
	if next[0].Acceleration != (Vector{}) {
		t.Fatalf("acceleration = %+v, want zero after recalculation", next[0].Acceleration)
	}
}
