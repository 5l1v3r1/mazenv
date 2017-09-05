package mazenv

import (
	"testing"

	"github.com/unixpickle/anyvec/anyvec64"
)

func TestSurroundingsEnv(t *testing.T) {
	maze, err := ParseMaze("...w\n" + ".wxw\n" + "Awww\n" + "...w")
	if err != nil {
		t.Fatal(err)
	}

	cr := anyvec64.DefaultCreator{}
	env := &SurroundingsEnv{Env: NewEnv(cr, maze), Horizon: 1}

	obs, err := env.Reset()
	if err != nil {
		t.Fatal(err)
	}

	testObsEqual(t, obs, []float64{
		0, 0, 1, 0, 0,
		0, 1, 0, 0, 0,
		0, 0, 1, 0, 0,

		0, 0, 1, 0, 0,
		1, 0, 0, 1, 0,
		0, 0, 1, 0, 0,

		0, 0, 1, 0, 0,
		0, 1, 0, 0, 0,
		0, 1, 0, 0, 0,
	})

	lastObs := obs
	for _, act := range []int{ActionUp, ActionUp, ActionRight, ActionRight} {
		obs, reward, done, err := env.Step(oneHotAction(cr, act))
		testNotDoneStepResult(t, reward, done, err)
		if vecsClose(lastObs, obs) {
			t.Errorf("observation didn't change after %v", act)
		}
		lastObs = obs
	}

	obs, reward, done, err := env.Step(oneHotAction(cr, ActionDown))
	if err != nil {
		t.Error(err)
	}
	if !done {
		t.Error("expected done signal")
	}
	if reward != 0 {
		t.Error("expected reward of 0")
	}
	testObsEqual(t, obs, []float64{
		0, 1, 0, 0, 0,
		0, 1, 0, 0, 0,
		0, 0, 1, 0, 0,

		0, 0, 1, 0, 0,
		1, 0, 0, 0, 1,
		0, 0, 1, 0, 0,

		0, 0, 1, 0, 0,
		0, 0, 1, 0, 0,
		0, 0, 1, 0, 0,
	})

	_, _, _, err = env.Step(oneHotAction(cr, ActionUp))
	if err == nil {
		t.Error("expected error from step after end of episode")
	}
}
