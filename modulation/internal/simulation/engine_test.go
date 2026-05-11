package simulation

import (
	"testing"

	"github.com/n1ckerr0r/cosmo-simulation-and-ml-predictions/modulation/internal/physics"
)

func TestEngineRunCallsStepAndCallback(t *testing.T) {
	engine := NewEngine([]physics.Body{{Mass: 1}}, 0.5)
	var callbacks int

	engine.Run(3, func(bodies []physics.Body, dt float64) []physics.Body {
		if dt != 0.5 {
			t.Fatalf("dt = %v", dt)
		}
		bodies[0].Position.X++
		return bodies
	}, func(step int, bodies []physics.Body) {
		callbacks++
		if step != callbacks-1 {
			t.Fatalf("step = %v, want %v", step, callbacks-1)
		}
	})

	if callbacks != 3 {
		t.Fatalf("callbacks = %v, want 3", callbacks)
	}
	if engine.Bodies[0].Position.X != 3 {
		t.Fatalf("final x = %v, want 3", engine.Bodies[0].Position.X)
	}
}
