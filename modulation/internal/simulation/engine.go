package simulation

import "github.com/n1ckerr0r/cosmo-simulation-and-ml-predictions/modulation/internal/physics"

type Engine struct {
	Bodies []physics.Body
	Dt     float64
}

func NewEngine(bodies []physics.Body, dt float64) *Engine {
	return &Engine{
		Bodies: bodies,
		Dt:     dt,
	}
}

// Run запускает симуляцию
// Принимает функцию для симуляции шага
// callback здесь используется для последующей работы с измерениями, например для вывода
func (e *Engine) Run(
	steps int,
	stepFunc func([]physics.Body, float64) []physics.Body,
	callback func(step int, bodies []physics.Body),
) {
	for i := 0; i < steps; i++ {
		e.Bodies = stepFunc(e.Bodies, e.Dt)

		if callback != nil {
			callback(i, e.Bodies)
		}
	}
}

