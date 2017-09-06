package mazenv

import (
	"errors"
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

	ends := shuffledSpaces(maze, maze.Start)
	if len(ends) == 0 {
		return nil, errors.New("no options for end")
	}
	maze.End = ends[0]

	return maze, nil
}

// IslandGenerator generates mazes by creating islands of
// walls such that all spaces are still connected.
//
// See the Python examples in:
// https://en.wikipedia.org/wiki/Maze_generation_algorithm.
type IslandGenerator struct {
	// Density controls the number of islands.
	// The value may range from 0 to 1.
	//
	// If 0, a default of 0.75 is used.
	Density float64

	// Complexity controls how large each island gets.
	// The value may range from 0 to 1.
	//
	// If 0, a default of 0.75 is used.
	Complexity float64
}

// Description returns a short description of what the
// algorithm does.
func (i *IslandGenerator) Description() string {
	return "create islands of walls in odd positions"
}

// AddFlags adds the generator's options as flags.
func (i *IslandGenerator) AddFlags(fs *flag.FlagSet) {
	fs.Float64Var(&i.Density, "density", 0.75, "the number of islands")
	fs.Float64Var(&i.Complexity, "complexity", 0.75, "the size of islands")
}

// Generate generates a random maze.
//
// Both dimensions must be odd.
func (i *IslandGenerator) Generate(rows, cols int) (*Maze, error) {
	if rows%2 == 0 || cols%2 == 0 {
		return nil, errors.New("maze dimensions must be odd")
	}

	maze := &Maze{
		Rows:  rows,
		Cols:  cols,
		Walls: make([]bool, rows*cols),
	}

	numIslands, islandSize := i.adjustedParams(rows, cols)

	for i := 0; i < numIslands; i++ {
		// Select an island start, which can possibly be on
		// the border around the grid.
		curPos := Position{
			Row: rand.Intn(rows/2+2)*2 - 1,
			Col: rand.Intn(cols/2+2)*2 - 1,
		}

		if maze.InBounds(curPos) {
			maze.Walls[maze.CellIndex(curPos)] = true
		}

		for j := 0; j < islandSize; j++ {
			destinations := spacesTwoCellsAway(maze, curPos)
			if len(destinations) == 0 {
				break
			}
			destination := destinations[rand.Intn(len(destinations))]
			midpoint := Position{
				Row: curPos.Row + (destination.Row-curPos.Row)/2,
				Col: curPos.Col + (destination.Col-curPos.Col)/2,
			}
			for _, setMe := range []Position{midpoint, destination} {
				maze.Walls[maze.CellIndex(setMe)] = true
			}
			curPos = destination
		}
	}

	spaces := shuffledSpaces(maze)
	if len(spaces) < 2 {
		return nil, errors.New("not enough spaces")
	}
	maze.Start, maze.End = spaces[0], spaces[1]

	return maze, nil
}

func (i *IslandGenerator) adjustedParams(rows, cols int) (numIslands, islandSize int) {
	density := i.Density
	if density == 0 {
		density = 0.75
	}
	numIslands = int(density * float64(rows/2+cols/2))

	complexity := i.Complexity
	if complexity == 0 {
		complexity = 0.75
	}
	islandSize = int(complexity * 5 * float64(rows+cols))

	return
}

func spacesTwoCellsAway(m *Maze, pos Position) []Position {
	var res []Position
	for _, delta := range []int{-2, 2} {
		new1 := pos
		new2 := pos
		new1.Row += delta
		new2.Col += delta
		for _, p := range []Position{new1, new2} {
			if !m.Wall(p) {
				res = append(res, p)
			}
		}
	}
	return res
}
