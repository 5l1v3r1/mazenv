package mazenv

import (
	"reflect"
	"testing"
)

func TestSolve(t *testing.T) {
	unsolvable, err := ParseMaze("ww..ww\nxw....\nwAwwww")
	if err != nil {
		t.Fatal(err)
	}
	if Solve(unsolvable) != nil {
		t.Error("found false solution")
	}

	solvable, err := ParseMaze("ww...w\nww.w.w\nwwAwx.\nww..w.\nww....")
	if err != nil {
		t.Fatal(err)
	}
	actual := Solve(solvable)
	expected := []Position{
		{2, 2},
		{1, 2},
		{0, 2},
		{0, 3},
		{0, 4},
		{1, 4},
		{2, 4},
	}
	compActual := make([]Position, len(actual))
	copy(compActual, actual)
	if !reflect.DeepEqual(compActual, expected) {
		t.Errorf("expected %v but got %v", expected, actual)
	}
}
