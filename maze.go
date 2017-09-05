package mazenv

import (
	"strings"
)

// Position represents a place on a Maze.
type Position struct {
	Row int
	Col int
}

func (p Position) addOne() Position {
	return Position{Row: p.Row + 1, Col: p.Col + 1}
}

// Maze defines a matrix of cells which comprise a maze.
type Maze struct {
	Rows int
	Cols int

	Start Position
	End   Position

	// Walls is a row-major list specifying which cells
	// in the maze are walls.
	//
	// There should be no wall where the start and end
	// positions are.
	Walls []bool
}

// InBounds checks if the position is within the grid.
func (m *Maze) InBounds(pos Position) bool {
	return pos.Row >= 0 && pos.Row < m.Rows &&
		pos.Col >= 0 && pos.Col < m.Cols
}

// Positions returns all valid positions in the grid in
// the same order as m.Walls.
func (m *Maze) Positions() []Position {
	var res []Position
	for row := 0; row < m.Rows; row++ {
		for col := 0; col < m.Cols; col++ {
			res = append(res, Position{row, col})
		}
	}
	return res
}

// Wall checks if the grid entry is a wall.
//
// If the cell is out of bounds, true is returned.
// This way, mazes are implicitly surrounded by walls.
func (m *Maze) Wall(pos Position) bool {
	if !m.InBounds(pos) {
		return true
	}
	return m.Walls[m.CellIndex(pos)]
}

// CellIndex gets the index for the cell.
//
// The position must be within bounds.
func (m *Maze) CellIndex(pos Position) int {
	if !m.InBounds(pos) {
		panic("out of bounds")
	}
	return pos.Row*m.Cols + pos.Col
}

// Bordered creates a new maze by adding a border of wall
// cells to the grid.
func (m *Maze) Bordered() *Maze {
	res := &Maze{
		Rows:  m.Rows + 2,
		Cols:  m.Cols + 2,
		Start: m.Start.addOne(),
		End:   m.End.addOne(),
		Walls: make([]bool, (m.Rows+2)*(m.Cols+2)),
	}
	for i := 0; i < res.Rows; i++ {
		res.Walls[res.CellIndex(Position{i, 0})] = true
		res.Walls[res.CellIndex(Position{i, res.Cols - 1})] = true
	}
	for i := 0; i < res.Cols; i++ {
		res.Walls[res.CellIndex(Position{0, i})] = true
		res.Walls[res.CellIndex(Position{res.Rows - 1, i})] = true
	}
	for _, pos := range m.Positions() {
		res.Walls[res.CellIndex(pos.addOne())] = m.Wall(pos)
	}
	return res
}

// String produces an ASCII representation of the grid.
// Every wall is represented as a 'w', every space as a
// '.', the start as 'A', and the end as 'x'.
func (m *Maze) String() string {
	rows := make([]string, m.Rows)
	for row := 0; row < m.Rows; row++ {
		for col := 0; col < m.Cols; col++ {
			pos := Position{Row: row, Col: col}
			ch := '.'
			if pos == m.Start {
				ch = 'A'
			} else if pos == m.End {
				ch = 'x'
			} else if m.Wall(Position{row, col}) {
				ch = 'w'
			}
			rows[row] += string(ch)
		}
	}
	return strings.Join(rows, "\n")
}
