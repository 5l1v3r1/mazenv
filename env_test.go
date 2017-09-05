package mazenv

import (
	"testing"

	"github.com/unixpickle/anyvec"
	"github.com/unixpickle/anyvec/anyvec64"
)

func TestEnv(t *testing.T) {
	maze, err := ParseMaze("...w\n" + ".wxw\n" + "Awww\n" + "...w")
	if err != nil {
		t.Fatal(err)
	}

	cr := anyvec64.DefaultCreator{}
	env := &Env{Creator: cr, Maze: maze}

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
		obs, reward, done, err := env.Step(oneHotAction(cr, act))
		testNotDoneStepResult(t, reward, done, err)
		testObsEqual(t, obs, expectedInitial)
	}

	obs, reward, done, err := env.Step(oneHotAction(cr, ActionDown))
	testNotDoneStepResult(t, reward, done, err)
	downRes := append([]float64{}, expectedInitial...)
	downRes[8*5] = 0
	downRes[12*5] = 1
	testObsEqual(t, obs, downRes)

	obs, reward, done, err = env.Step(oneHotAction(cr, ActionUp))
	testNotDoneStepResult(t, reward, done, err)
	testObsEqual(t, obs, expectedInitial)

	lastObs := obs
	for _, act := range []int{ActionUp, ActionUp, ActionRight, ActionRight} {
		obs, reward, done, err := env.Step(oneHotAction(cr, act))
		testNotDoneStepResult(t, reward, done, err)
		if vecsClose(lastObs, obs) {
			t.Errorf("observation didn't change after %v", act)
		}
		lastObs = obs
	}

	obs, reward, done, err = env.Step(oneHotAction(cr, ActionDown))
	if err != nil {
		t.Error(err)
	}
	if !done {
		t.Error("expected done signal")
	}
	if reward != 1 {
		t.Error("expected reward of 1")
	}
	doneRes := append([]float64{}, expectedInitial...)
	doneRes[8*5] = 0
	doneRes[6*5] = 1
	testObsEqual(t, obs, doneRes)

	_, _, _, err = env.Step(oneHotAction(cr, ActionUp))
	if err == nil {
		t.Error("expected error from step after end of episode")
	}
}

func testObsEqual(t *testing.T, actual anyvec.Vector, expected []float64) {
	if !obsEqual(actual, expected) {
		t.Errorf("expected %v but got %v", expected, actual.Data())
	}
}

func testNotDoneStepResult(t *testing.T, reward float64, done bool, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if reward != 0 {
		t.Error("unexpected reward")
	}
	if done {
		t.Fatal("unexpected done")
	}
}

func oneHotAction(cr anyvec.Creator, idx int) anyvec.Vector {
	data := make([]float64, 5)
	data[idx] = 1
	return cr.MakeVectorData(cr.MakeNumericList(data))
}

func obsEqual(ob1 anyvec.Vector, ob2 []float64) bool {
	cr := ob1.Creator()
	return vecsClose(ob1, cr.MakeVectorData(cr.MakeNumericList(ob2)))
}

func vecsClose(v1, v2 anyvec.Vector) bool {
	cr := v1.Creator()
	diff := v2.Copy()
	diff.Sub(v1)
	diffNorm := anyvec.Norm(diff)
	return cr.NumOps().Less(diffNorm, cr.MakeNumeric(1e-3))
}
