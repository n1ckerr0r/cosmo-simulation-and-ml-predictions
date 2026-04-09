package physics

func KineticEnergy(bodies []Body) float64 {
	var energy float64

	for _, body := range bodies {
		velocity := body.Velocity.Length()
		energy += 0.5 * body.Mass * velocity * velocity
	}
	return energy
}

func PotentialEnergy(bodies []Body) float64 {
	var energy float64

	for i := 0; i < len(bodies); i++ {
		for j := i + 1; j < len(bodies); j++ {
			// Прибавляем константу, чтобы избежать деления на 0
			r := bodies[j].Position.Sub(bodies[i].Position).Length() + 1e-9
			energy -= G * bodies[i].Mass * bodies[j].Mass / r
		}
	}

	return energy
}

func TotalEnergy(bodies []Body) float64 {
	return KineticEnergy(bodies) + PotentialEnergy(bodies)
}

