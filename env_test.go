package mazenv

import (
	"math"
	"testing"
)

func TestEnv(t *testing.T) {
	maze, err := ParseMaze("...w\n" + ".wxw\n" + "Awww\n" + "...w")
	if err != nil {
		t.Fatal(err)
	}

	env := NewEnv(maze)

	obs, err := env.Reset()
	if err != nil {
		t.Fatal(err)
	}

	expectedInitial := []float64{
		0, 1, 0, 0, 0,
		0, 1, 0, 0, 0,
		0, 1, 0, 0, 0,
		0, 0, 1, 0, 0,

		0, 1, 0, 0, 0,
		0, 0, 1, 0, 0,
		0, 0, 0, 0, 1,
		0, 0, 1, 0, 0,

		1, 0, 0, 1, 0,
		0, 0, 1, 0, 0,
		0, 0, 1, 0, 0,
		0, 0, 1, 0, 0,

		0, 1, 0, 0, 0,
		0, 1, 0, 0, 0,
		0, 1, 0, 0, 0,
		0, 0, 1, 0, 0,
	}
	testObsEqual(t, obs, expectedInitial)

	// All these directions should have no effect.
	for _, act := range []int{ActionNop, ActionRight, ActionLeft} {
		obs, reward, done, err := env.Step(oneHotAction(act))
		testNotDoneStepResult(t, reward, done, err)
		testObsEqual(t, obs, expectedInitial)
	}

	obs, reward, done, err := env.Step(oneHotAction(ActionDown))
	testNotDoneStepResult(t, reward, done, err)
	downRes := append([]float64{}, expectedInitial...)
	downRes[8*5] = 0
	downRes[12*5] = 1
	testObsEqual(t, obs, downRes)

	obs, reward, done, err = env.Step(oneHotAction(ActionUp))
	testNotDoneStepResult(t, reward, done, err)
	testObsEqual(t, obs, expectedInitial)

	lastObs := obs
	for _, act := range []int{ActionUp, ActionUp, ActionRight, ActionRight} {
		obs, reward, done, err := env.Step(oneHotAction(act))
		testNotDoneStepResult(t, reward, done, err)
		if obsEqual(lastObs, obs) {
			t.Errorf("observation didn't change after %v", act)
		}
		lastObs = obs
	}

	obs, reward, done, err = env.Step(oneHotAction(ActionDown))
	if err != nil {
		t.Error(err)
	}
	if !done {
		t.Error("expected done signal")
	}
	if reward != 0 {
		t.Error("expected reward of 0")
	}
	doneRes := append([]float64{}, expectedInitial...)
	doneRes[8*5] = 0
	doneRes[6*5] = 1
	testObsEqual(t, obs, doneRes)

	_, _, _, err = env.Step(oneHotAction(ActionUp))
	if err == nil {
		t.Error("expected error from step after end of episode")
	}
}

func testObsEqual(t *testing.T, actual, expected []float64) {
	if !obsEqual(actual, expected) {
		t.Errorf("expected %v but got %v", expected, actual)
	}
}

func testNotDoneStepResult(t *testing.T, reward float64, done bool, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if reward != -1 {
		t.Error("unexpected reward")
	}
	if done {
		t.Fatal("unexpected done")
	}
}

func oneHotAction(idx int) []float64 {
	data := make([]float64, 5)
	data[idx] = 1
	return data
}

func obsEqual(ob1 []float64, ob2 []float64) bool {
	if len(ob1) != len(ob2) {
		return false
	}
	for i, x1 := range ob1 {
		x2 := ob2[i]
		if math.Abs(x1-x2) > 1e-3 {
			return false
		}
	}
	return true
}
