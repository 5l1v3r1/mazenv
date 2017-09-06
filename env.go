package mazenv

import (
	"errors"

	"github.com/unixpickle/anyrl"
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

	// Maze returns the environment's map.
	Maze() *Maze

	// Position returns the player's current position.
	Position() Position
}

// rawEnv is a barebones environment for a maze.
type rawEnv struct {
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
func NewEnv(maze *Maze) Env {
	return &rawEnv{maze: maze}
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
func (r *rawEnv) Reset() (obs []float64, err error) {
	r.position = r.maze.Start
	return r.observation(), nil
}

// Step takes a step in the environment.
func (r *rawEnv) Step(action []float64) (obs []float64, reward float64,
	done bool, err error) {
	if r.position == r.maze.End {
		err = errors.New("step: maze is already solved")
		return
	}
	newPos := r.position
	var actionIdx int
	for i, x := range action {
		if x != 0 {
			actionIdx = i
		}
	}
	switch actionIdx {
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

func (r *rawEnv) observation() []float64 {
	return oneHotGrid(r.maze, r.position, 0, 0, r.maze.Rows, r.maze.Cols)
}
