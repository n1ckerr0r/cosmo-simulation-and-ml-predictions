package main

import (
	"fmt"

	"github.com/n1ckerr0r/cosmo-simulation-and-ml-predictions/modulation/internal/output"
	"github.com/n1ckerr0r/cosmo-simulation-and-ml-predictions/modulation/internal/physics"
	"github.com/n1ckerr0r/cosmo-simulation-and-ml-predictions/modulation/internal/simulation"
)

func main() {
	bodies := []physics.Body{
		// Солнце
		{
			Mass:     1.989e30,
			Position: physics.Vector{0, 0, 0},
			Velocity: physics.Vector{0, 0, 0},
		},

		// Земля
		{
			Mass:     5.972e24,
			Position: physics.Vector{1.496e11, 0, 0},
			Velocity: physics.Vector{0, 29780, 0},
		},

		// Луна
		{
			Mass:     7.35e22,
			Position: physics.Vector{1.496e11 + 3.84e8, 0, 0},
			Velocity: physics.Vector{0, 29780 + 1022, 0},
		},

		// Марс
		{
			Mass:     6.39e23,
			Position: physics.Vector{2.279e11, 0, 0},
			Velocity: physics.Vector{0, 24070, 0},
		},

		// Юпитер
		{
			Mass:     1.898e27,
			Position: physics.Vector{7.785e11, 0, 0},
			Velocity: physics.Vector{0, 13070, 0},
		},
	}

	acc := physics.CalculateAccelerations(bodies)
	for i := range bodies {
		bodies[i].Acceleration = acc[i]
	}

	engine := simulation.NewEngine(bodies, 600)

	writer := output.NewWriter("data/output_case5.csv")

	steps := 24 * 365 * 5

	engine.Run(
		steps,
		physics.Step,
		func(step int, bodies []physics.Body) {
			if step%100 == 0 {
				E := physics.TotalEnergy(bodies)
				fmt.Println("step:", step, "energy:", E)
			}
			writer.Write(step, bodies)
		},
	)

	writer.Close()
	fmt.Println("Simulation finished")
}
