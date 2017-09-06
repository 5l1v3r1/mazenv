package mazenv

import "math/rand"

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

func oneHotGrid(m *Maze, curPos Position, startRow, startCol, rows, cols int) []float64 {
	var res []float64
	for row := startRow; row < startRow+rows; row++ {
		for col := startCol; col < startCol+cols; col++ {
			pos := Position{row, col}
			cellType := CellEmpty
			if m.Start == pos {
				cellType = CellStart
			} else if m.End == pos {
				cellType = CellEnd
			} else if m.Wall(pos) {
				cellType = CellWall
			}
			if pos == curPos {
				res = append(res, 1)
			} else {
				res = append(res, 0)
			}
			res = append(res, oneHot(4, cellType)...)
		}
	}
	return res
}

func oneHot(num, val int) []float64 {
	res := make([]float64, num)
	res[val] = 1
	return res
}

func shuffledSpaces(m *Maze, exclude ...Position) []Position {
	var options []Position

PosLoop:
	for _, p := range m.Positions() {
		if !m.Wall(p) {
			for _, x := range exclude {
				if x == p {
					continue PosLoop
				}
			}
			options = append(options, p)
		}
	}

	for i := 0; i < len(options); i++ {
		j := i + rand.Intn(len(options)-i)
		options[i], options[j] = options[j], options[i]
	}

	return options
}
