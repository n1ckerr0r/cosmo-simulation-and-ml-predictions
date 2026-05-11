package simulation

import (
	"math"

	"github.com/n1ckerr0r/cosmo-simulation-and-ml-predictions/modulation/internal/physics"
)

const (
	sunMass = 1.989e30
	day     = 24 * 60 * 60
)

type Scenario struct {
	Name        string
	OutputPath  string
	Dt          float64
	Steps       int
	SampleEvery int
	Bodies      []physics.Body
}

func PlanetOrbitScenarios() []Scenario {
	return []Scenario{
		newPlanetOrbitScenario(
			"mercury_one_orbit",
			"data/output_mercury.csv",
			5.7909e10,
			3.301e23,
			7.0,
			100,
			30*60,
			4,
		),
		newPlanetOrbitScenario(
			"earth_one_orbit",
			"data/output_earth.csv",
			1.496e11,
			5.972e24,
			23.4,
			370,
			60*60,
			6,
		),
		newPlanetOrbitScenario(
			"mars_one_orbit",
			"data/output_mars.csv",
			2.279e11,
			6.39e23,
			25.2,
			700,
			2*60*60,
			6,
		),
	}
}

func PrepareBodies(bodies []physics.Body) []physics.Body {
	prepared := make([]physics.Body, len(bodies))
	copy(prepared, bodies)

	accelerations := physics.CalculateAccelerations(prepared)
	for i := range prepared {
		prepared[i].Acceleration = accelerations[i]
	}

	return prepared
}

func newPlanetOrbitScenario(name, outputPath string, radius, planetMass, inclinationDegrees float64, days int, dt float64, sampleEvery int) Scenario {
	speed := math.Sqrt(physics.G * sunMass / radius)
	inclination := inclinationDegrees * math.Pi / 180

	return Scenario{
		Name:        name,
		OutputPath:  outputPath,
		Dt:          dt,
		Steps:       int(float64(days*day) / dt),
		SampleEvery: sampleEvery,
		Bodies: []physics.Body{
			{
				Mass:     sunMass,
				Position: physics.Vector{X: 0, Y: 0, Z: 0},
				Velocity: physics.Vector{X: 0, Y: 0, Z: 0},
			},
			{
				Mass:     planetMass,
				Position: physics.Vector{X: radius, Y: 0, Z: 0},
				Velocity: physics.Vector{X: 0, Y: speed * math.Cos(inclination), Z: speed * math.Sin(inclination)},
			},
		},
	}
}
