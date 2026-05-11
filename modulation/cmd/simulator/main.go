package main

import (
	"fmt"

	"github.com/n1ckerr0r/cosmo-simulation-and-ml-predictions/modulation/internal/output"
	"github.com/n1ckerr0r/cosmo-simulation-and-ml-predictions/modulation/internal/physics"
	"github.com/n1ckerr0r/cosmo-simulation-and-ml-predictions/modulation/internal/simulation"
)

func main() {
	for _, scenario := range simulation.PlanetOrbitScenarios() {
		runScenario(scenario)
	}

	fmt.Println("Simulation finished")
}

func runScenario(scenario simulation.Scenario) {
	bodies := simulation.PrepareBodies(scenario.Bodies)
	engine := simulation.NewEngine(bodies, scenario.Dt)
	writer := output.NewWriter(scenario.OutputPath)
	defer writer.Close()

	engine.Run(scenario.Steps, physics.Step, func(step int, bodies []physics.Body) {
		if step%scenario.SampleEvery == 0 {
			writer.Write(step, bodies)
		}
		if step%(scenario.SampleEvery*100) == 0 {
			energy := physics.TotalEnergy(bodies)
			fmt.Println("scenario:", scenario.Name, "step:", step, "energy:", energy)
		}
	})

	fmt.Println("saved:", scenario.OutputPath)
}
