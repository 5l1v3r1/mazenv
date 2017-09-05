package mazenv

import "github.com/unixpickle/anyvec"

// SurroundingsEnv restricts the observations of an Env.
// In particular, it shows the agent an NxN rectangle with
// the agent at the center, where N is 2*Horizon+1.
//
// The agent will see walls in parts of its vision that go
// beyond the grid's bounds.
type SurroundingsEnv struct {
	Env

	// Horizon is the distance the agent can see.
	//
	// For example, if the Horizon is 1, then the agent
	// only sees a 3x3 grid with the agent at the center.
	Horizon int
}

// Reset resets the environment.
func (s *SurroundingsEnv) Reset() (obs anyvec.Vector, err error) {
	_, err = s.Env.Reset()
	if err != nil {
		return
	}
	obs = s.observe()
	return
}

// Step takes a step in the environment.
func (s *SurroundingsEnv) Step(act anyvec.Vector) (obs anyvec.Vector, rew float64,
	done bool, err error) {
	_, rew, done, err = s.Env.Step(act)
	if err != nil {
		return
	}
	obs = s.observe()
	return
}

func (s *SurroundingsEnv) observe() anyvec.Vector {
	p := s.Position()
	startRow := p.Row - s.Horizon
	startCol := p.Col - s.Horizon
	size := 2*s.Horizon + 1
	grid := oneHotGrid(s.Maze(), p, startRow, startCol, size, size)
	vecData := s.Creator().MakeNumericList(grid)
	return s.Creator().MakeVectorData(vecData)
}
