package main

import (
	"fmt"

	"github.com/n1ckerr0r/cosmo-simulation-and-ml-predictions/modulation/internal/output"
	"github.com/n1ckerr0r/cosmo-simulation-and-ml-predictions/modulation/internal/physics"
	"github.com/n1ckerr0r/cosmo-simulation-and-ml-predictions/modulation/internal/simulation"
)

func main() {
    bodies := []physics.Body{
        {
            Mass:     1.989e30,
            Position: physics.Vector{X: 0, Y: 0, Z: 0},
            Velocity: physics.Vector{X: 0, Y: 0, Z: 0},
        },
        {
            Mass:     5.972e24,
            Position: physics.Vector{X: 1.496e11, Y: 0, Z: 0},
            Velocity: physics.Vector{X: 0, Y: 29780, Z: 0},
        },
    }

    acc := physics.CalculateAccelerations(bodies)
    for i := range bodies {
        bodies[i].Acceleration = acc[i]
    }

    engine := simulation.NewEngine(bodies, 60*60)
    writer := output.NewWriter("data/output.csv")
    defer writer.Close()

    steps := 24 * 365 * 3

    engine.Run(steps, physics.Step, func(step int, bodies []physics.Body) {
        if step%100 == 0 {
            energy := physics.TotalEnergy(bodies)
            fmt.Println("step:", step, "energy:", energy)
        }
        writer.Write(step, bodies)
    })

    fmt.Println("Simulation finished")
}
