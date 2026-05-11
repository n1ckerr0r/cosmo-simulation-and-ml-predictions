package simulation

import (
	"math"
	"testing"

	"github.com/n1ckerr0r/cosmo-simulation-and-ml-predictions/modulation/internal/physics"
)

func TestPlanetOrbitScenariosCompleteAtLeastOneTurn(t *testing.T) {
	for _, scenario := range PlanetOrbitScenarios() {
		bodies := PrepareBodies(scenario.Bodies)
		engine := NewEngine(bodies, scenario.Dt)
		previousAngle := angleBetween(bodies[0], bodies[1])
		var totalAngle float64

		engine.Run(scenario.Steps, physics.Step, func(_ int, bodies []physics.Body) {
			currentAngle := angleBetween(bodies[0], bodies[1])
			delta := currentAngle - previousAngle
			if delta < -math.Pi {
				delta += 2 * math.Pi
			}
			if delta > math.Pi {
				delta -= 2 * math.Pi
			}
			totalAngle += delta
			previousAngle = currentAngle
		})

		if totalAngle < 2*math.Pi {
			t.Fatalf("%s completed %v radians, want at least %v", scenario.Name, totalAngle, 2*math.Pi)
		}
	}
}

func TestPrepareBodiesCopiesInputAndSetsAccelerations(t *testing.T) {
	input := PlanetOrbitScenarios()[0].Bodies

	prepared := PrepareBodies(input)

	if len(prepared) != len(input) {
		t.Fatalf("len = %v, want %v", len(prepared), len(input))
	}
	if prepared[1].Acceleration == (physics.Vector{}) {
		t.Fatal("planet acceleration was not initialized")
	}
	prepared[0].Mass = 1
	if input[0].Mass == 1 {
		t.Fatal("PrepareBodies changed input slice")
	}
}

func angleBetween(center, planet physics.Body) float64 {
	relative := planet.Position.Sub(center.Position)
	return math.Atan2(relative.Y, relative.X)
}
