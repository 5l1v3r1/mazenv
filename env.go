package mazenv

import (
	"errors"

	"github.com/unixpickle/anyrl"
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

// Env is a generic maze environment.
//
// Actions are one-hot vectors with five possibilities.
// See ActionNop, ActionUp, etc.
//
// Observations and rewards are dependent on context.
// See, for example, NewEnv.
type Env interface {
	anyrl.Env

	// Creator is what the environment uses to create
	// observations.
	Creator() anyvec.Creator

	// Maze returns the environment's map.
	Maze() *Maze

	// Position returns the player's current position.
	Position() Position
}

// rawEnv is a barebones environment for a maze.
type rawEnv struct {
	creator  anyvec.Creator
	maze     *Maze
	position Position
}

// NewEnv creates an Env for the maze.
//
// Observations are row-major representations of the
// maze grid.
// Each cell is represented as a boolean (is current
// position) followed by a one-hot vector of four
// components: space, wall, start, end.
//
// Rewards are -1 until the maze is solved, at which point
// the episode ends and the reward is 0.
// This way, shorter solutions are preferred.
func NewEnv(cr anyvec.Creator, maze *Maze) Env {
	return &rawEnv{creator: cr, maze: maze}
}

// Creator returns the creator for observations.
func (r *rawEnv) Creator() anyvec.Creator {
	return r.creator
}

// Maze returns the maze.
func (r *rawEnv) Maze() *Maze {
	return r.maze
}

// Position returns the current position.
func (r *rawEnv) Position() Position {
	return r.position
}

// Reset resets the player's position to the start.
func (r *rawEnv) Reset() (obs anyvec.Vector, err error) {
	r.position = r.maze.Start
	return r.observation(), nil
}

// Step takes a step in the environment.
func (r *rawEnv) Step(action anyvec.Vector) (obs anyvec.Vector, reward float64,
	done bool, err error) {
	if r.position == r.maze.End {
		err = errors.New("step: maze is already solved")
		return
	}
	newPos := r.position
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
	if !r.maze.Wall(newPos) {
		r.position = newPos
	}
	if r.position == r.maze.End {
		reward = 0
		done = true
	} else {
		reward = -1
	}
	obs = r.observation()
	return
}

func (r *rawEnv) observation() anyvec.Vector {
	grid := oneHotGrid(r.maze, r.position, 0, 0, r.maze.Rows, r.maze.Cols)
	vecData := r.creator.MakeNumericList(grid)
	return r.creator.MakeVectorData(vecData)
}
