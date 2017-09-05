package mazenv

import (
	"flag"
	"math/rand"

	"github.com/unixpickle/essentials"
)

// A Generator generates mazes.
type Generator interface {
	Generate(rows, cols int) (*Maze, error)
}

// PrimGenerator is a Generator that uses a randomized
// variant of Prim's algorithm.
type PrimGenerator struct{}

// Description returns a short description of what the
// algorithm does.
func (p *PrimGenerator) Description() string {
	return "randomized variant of Prim's algorithm"
}

// AddFlags adds the generator's options as flags.
//
// Currently, this is a no-op.
func (p *PrimGenerator) AddFlags(f *flag.FlagSet) {
}

// Generate generates a random maze.
func (p *PrimGenerator) Generate(rows, cols int) (*Maze, error) {
	maze := &Maze{
		Rows: rows,
		Cols: cols,
		Start: Position{
			Row: rand.Intn(rows),
			Col: rand.Intn(cols),
		},
		Walls: make([]bool, rows*cols),
	}

	for i := range maze.Walls {
		maze.Walls[i] = true
	}
	maze.Walls[maze.CellIndex(maze.Start)] = false

	edges := neighbors(maze, maze.Start)
	visited := map[Position]bool{maze.Start: true}
	for _, p := range edges {
		visited[p] = true
	}
	for len(edges) > 0 {
		idx := rand.Intn(len(edges))
		pos := edges[idx]
		essentials.UnorderedDelete(&edges, idx)
		if len(neighboringSpaces(maze, pos)) > 1 {
			continue
		}
		maze.Walls[maze.CellIndex(pos)] = false
		for _, neighbor := range neighbors(maze, pos) {
			if !visited[neighbor] {
				edges = append(edges, neighbor)
				visited[neighbor] = true
			}
		}
	}

	var possibleEnds []Position
	for _, pos := range maze.Positions() {
		if !maze.Wall(pos) && pos != maze.Start {
			possibleEnds = append(possibleEnds, pos)
		}
	}
	maze.End = possibleEnds[rand.Intn(len(possibleEnds))]

	return maze, nil
}
