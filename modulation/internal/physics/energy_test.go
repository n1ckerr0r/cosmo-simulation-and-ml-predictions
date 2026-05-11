package physics

import (
	"math"
	"testing"
)

func TestEnergyCalculations(t *testing.T) {
	bodies := []Body{
		{
			Mass:     2,
			Position: Vector{X: 0, Y: 0, Z: 0},
			Velocity: Vector{X: 3, Y: 0, Z: 0},
		},
		{
			Mass:     4,
			Position: Vector{X: 2, Y: 0, Z: 0},
			Velocity: Vector{X: 0, Y: 2, Z: 0},
		},
	}

	if got := KineticEnergy(bodies); got != 17 {
		t.Fatalf("KineticEnergy() = %v", got)
	}

	wantPotential := -G * bodies[0].Mass * bodies[1].Mass / (2 + 1e-9)
	if got := PotentialEnergy(bodies); math.Abs(got-wantPotential) > 1e-20 {
		t.Fatalf("PotentialEnergy() = %v, want %v", got, wantPotential)
	}

	wantTotal := 17 + wantPotential
	if got := TotalEnergy(bodies); math.Abs(got-wantTotal) > 1e-20 {
		t.Fatalf("TotalEnergy() = %v, want %v", got, wantTotal)
	}
}
