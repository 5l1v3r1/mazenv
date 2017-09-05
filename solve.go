package mazenv

// Solve finds an optimal solution to the maze.
//
// The solution is represented as a list of positions that
// comprise the solution, including the start and end.
//
// If no solution is found, nil is returned.
func Solve(m *Maze) []Position {
	queue := []searchNode{{Path: []Position{m.Start}}}
	visited := map[Position]bool{m.Start: true}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		for _, neighbor := range neighboringSpaces(m, node.Pos()) {
			if neighbor == m.End {
				return node.Add(m.End).Path
			}
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, node.Add(neighbor))
			}
		}
	}
	return nil
}

type searchNode struct {
	Path []Position
}

func (s searchNode) Add(p Position) searchNode {
	return searchNode{Path: append(append([]Position{}, s.Path...), p)}
}

func (s searchNode) Pos() Position {
	return s.Path[len(s.Path)-1]
}
