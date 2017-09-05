package mazenv

import (
	"errors"

	"github.com/unixpickle/anyvec"
)

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

// MetaEnv creates a meta-learning environment around Env.
//
// Each meta-episode consists of NumRuns episodes of Env.
// The action space is the same, but the observations are
// augmented (at the end) with the previous action and
// reward (in that order).
//
// For the first observation, the action and reward values
// are both 0.
type MetaEnv struct {
	Env
	NumRuns int

	runsRemaining int
}

// Reset resets the environment.
func (m *MetaEnv) Reset() (obs anyvec.Vector, err error) {
	m.runsRemaining = m.NumRuns
	obs, err = m.Env.Reset()
	if err != nil {
		return
	}
	zeroVec := obs.Creator().MakeVector(6)
	obs = obs.Creator().Concat(obs, zeroVec)
	return
}

// Step takes a step in the environment.
func (m *MetaEnv) Step(act anyvec.Vector) (obs anyvec.Vector, rew float64,
	done bool, err error) {
	if m.runsRemaining <= 0 {
		err = errors.New("step: done sub-episodes in meta-environment")
		return
	}
	obs, rew, done, err = m.Env.Step(act)
	if err != nil {
		return
	}
	if done {
		m.runsRemaining--
		done = m.runsRemaining == 0
		if !done {
			obs, err = m.Env.Reset()
			if err != nil {
				return
			}
		}
	}
	c := obs.Creator()
	obs = c.Concat(obs, act, c.MakeVectorData(c.MakeNumericList([]float64{rew})))
	return
}
