package main

import "github.com/unixpickle/mazenv"

func neighbors(m *mazenv.Maze, p mazenv.Position) []mazenv.Position {
	var res []mazenv.Position
	for off := -1; off <= 1; off += 2 {
		newPos := mazenv.Position{Row: p.Row + off, Col: p.Col}
		if m.InBounds(newPos) {
			res = append(res, newPos)
		}
	}
	for off := -1; off <= 1; off += 2 {
		newPos := mazenv.Position{Row: p.Row, Col: p.Col + off}
		if m.InBounds(newPos) {
			res = append(res, newPos)
		}
	}
	return res
}

func neighboringSpaces(m *mazenv.Maze, p mazenv.Position) []mazenv.Position {
	var res []mazenv.Position
	for _, neighbor := range neighbors(m, p) {
		if !m.Wall(neighbor) {
			res = append(res, neighbor)
		}
	}
	return res
}
