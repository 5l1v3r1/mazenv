package mazenv

import (
	"errors"

	"github.com/unixpickle/anyvec"
)

// Indices in one-hot action vectors.
const (
	ActionNop = iota
	ActionUp
	ActionRight
	ActionDown
	ActionLeft
)

// Indices in one-hot cell observations.
const (
	CellEmpty = iota
	CellWall
	CellStart
	CellEnd
)

// Env is a barebones environment for a single maze.
//
// Actions are one-hot vectors with five possibilities.
// See ActionNop, ActionUp, etc.
//
// Observations are row-major representations of the
// maze grid.
// Each cell is represented as a boolean (is current
// position) followed by a one-hot vector of four
// components: space, wall, start, end.
//
// Rewards are 0 until the maze is solved, at which point
// the episode ends and the reward is 1.
type Env struct {
	// Creator is used to create observation vectors.
	Creator anyvec.Creator

	// Maze is the map to use for the environment.
	Maze *Maze

	// Position is the current player position.
	Position Position
}

// Reset resets the player's position to the start.
func (e *Env) Reset() (obs anyvec.Vector, err error) {
	e.Position = e.Maze.Start
	return e.observation(), nil
}

// Step takes a step in the environment.
func (e *Env) Step(action anyvec.Vector) (obs anyvec.Vector, reward float64,
	done bool, err error) {
	if e.Position == e.Maze.End {
		err = errors.New("step: maze is already solved")
		return
	}
	newPos := e.Position
	switch anyvec.MaxIndex(action) {
	case ActionUp:
		newPos.Row--
	case ActionRight:
		newPos.Col++
	case ActionDown:
		newPos.Row++
	case ActionLeft:
		newPos.Col--
	}
	if !e.Maze.Wall(newPos) {
		e.Position = newPos
	}
	if e.Position == e.Maze.End {
		reward = 1
		done = true
	}
	obs = e.observation()
	return
}

func (e *Env) observation() anyvec.Vector {
	grid := oneHotGrid(e.Maze, e.Position, 0, 0, e.Maze.Rows, e.Maze.Cols)
	vecData := e.Creator.MakeNumericList(grid)
	return e.Creator.MakeVectorData(vecData)
}
