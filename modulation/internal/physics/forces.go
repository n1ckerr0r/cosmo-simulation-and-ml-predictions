package physics

const G = 6.67430e-11

// CalculateAccelerations подсчитывает как все тела действуют друг на друга
// Возвращает вектор новых ускорейний тел
func CalculateAccelerations(bodies []Body) []Vector {

	n := len(bodies)
	newAccelerations := make([]Vector, n)

	for i := range n {
		for j := i + 1; j < n; j++ {
			// Расстояние между планетами
			posDiff := bodies[j].Position.Sub(bodies[i].Position)

			// Прибавляем маленько число для предовращения деления на ноль
			dist := posDiff.Length() + 1e-9

			// Подсчет силы по закону всемирного тяготения
			forceMag := G * bodies[i].Mass * bodies[j].Mass / (dist * dist)

			// Нормализованный вектор направления
			direction := posDiff.Mul(1 / dist)

			// Нормализованная сила
			force := direction.Mul(forceMag)

			newAccelerations[i] = newAccelerations[i].Add(force.Mul(1 / bodies[i].Mass))
			newAccelerations[j] = newAccelerations[j].Add(force.Mul(-1 / bodies[j].Mass))
		}
	}

	return newAccelerations
}
