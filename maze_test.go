package mazenv

import "testing"

func TestMazeString(t *testing.T) {
	maze := testingMaze()
	actual := maze.String()
	expected := ".ww.\nw.wx\n.Aww\nw..."
	if actual != expected {
		t.Errorf("expected %#v but got %#v", expected, actual)
	}

	shouldFail := []string{"xw\nA", "xw\nxA", "xw\n..", "Aw\n..", "Aw\nxA"}
	for _, s := range shouldFail {
		if _, err := ParseMaze(s); err == nil {
			t.Errorf("expected failure for %#v", s)
		}
	}
}

func TestMazeParse(t *testing.T) {
	expected := testingMaze()
	actual, err := ParseMaze(expected.String())
	if err != nil {
		t.Error(err)
	} else {
		if actual.String() != expected.String() {
			t.Errorf("expected %#v but got %#v", expected, actual)
		}
	}
}

func testingMaze() *Maze {
	return &Maze{
		Rows:  4,
		Cols:  4,
		Start: Position{2, 1},
		End:   Position{1, 3},
		Walls: []bool{
			false, true, true, false,
			true, false, true, false,
			false, false, true, true,
			true, false, false, false,
		},
	}
}
