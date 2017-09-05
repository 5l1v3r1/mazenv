package mazenv

func neighbors(m *Maze, p Position) []Position {
	var res []Position
	for off := -1; off <= 1; off += 2 {
		newPos := Position{p.Row + off, p.Col}
		if m.InBounds(newPos) {
			res = append(res, newPos)
		}
	}
	for off := -1; off <= 1; off += 2 {
		newPos := Position{p.Row, p.Col + off}
		if m.InBounds(newPos) {
			res = append(res, newPos)
		}
	}
	return res
}

func neighboringSpaces(m *Maze, p Position) []Position {
	var res []Position
	for _, neighbor := range neighbors(m, p) {
		if !m.Wall(neighbor) {
			res = append(res, neighbor)
		}
	}
	return res
}
