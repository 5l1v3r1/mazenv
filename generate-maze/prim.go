package main

import (
	"flag"
	"math/rand"

	"github.com/unixpickle/essentials"
	"github.com/unixpickle/mazenv"
)

type PrimAlgorithm struct{}

func (p *PrimAlgorithm) Description() string {
	return "randomized variant of Prim's algorithm"
}

func (p *PrimAlgorithm) AddFlags(f *flag.FlagSet) {
}

func (p *PrimAlgorithm) Generate(rows, cols int) (*mazenv.Maze, error) {
	maze := &mazenv.Maze{
		Rows: rows,
		Cols: cols,
		Start: mazenv.Position{
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
	visited := map[mazenv.Position]bool{maze.Start: true}
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

	var possibleEnds []mazenv.Position
	for _, pos := range maze.Positions() {
		if !maze.Wall(pos) && pos != maze.Start {
			possibleEnds = append(possibleEnds, pos)
		}
	}
	maze.End = possibleEnds[rand.Intn(len(possibleEnds))]

	return maze, nil
}
