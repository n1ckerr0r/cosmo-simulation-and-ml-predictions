package physics

// Step симулирет шаг во времени
func Step(bodies []Body, dt float64) []Body {
	n := len(bodies)

	for i := 0; i < n; i++ {
		b := &bodies[i]

		// Используем формулу равноускоренного движения
		b.Position = b.Position.Add(b.Velocity.Mul(dt)).Add(b.Acceleration.Mul(0.5 * dt * dt))
	}

	newAcceleration := CalculateAccelerations(bodies)

	// Считаем новые ускорения
	for i := 0; i < n; i++ {
		b := &bodies[i]

		b.Velocity = b.Velocity.Add(newAcceleration[i].Mul(dt))
		b.Acceleration = newAcceleration[i]
	}

	return bodies
} 

