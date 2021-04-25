package goschedule

import (
	"testing"
	"time"
)

func TestWorker(t *testing.T) {
	tests := []struct {
		implmenetation interface{}
		params         []interface{}
		result         interface{}
		identifier     string
		nextRun        time.Time
		every          time.Duration
	}{
		{
			implmenetation: func(l int, b int) int { return l - b },
			params:         []interface{}{1, 2},
			result:         -1,
			identifier:     "test3",
			nextRun:        time.Now().Add(5 * time.Second),
			every:          5 * time.Second,
		}, {
			implmenetation: func(l string, b string) string { return l + b },
			params:         []interface{}{"Mah", "moud"},
			result:         "Mahmoud",
			identifier:     "test4",
			nextRun:        time.Now().Add(5 * time.Second),
			every:          5 * time.Second,
		},
	}

	worker := newWorker()
	go worker.work()

	t1 := newFunctionJob("example1", tests[0].implmenetation, tests[0].params).At(tests[0].nextRun).Every(tests[0].every)
	t2 := newFunctionJob("example2", tests[1].implmenetation, tests[1].params).At(tests[1].nextRun).Every(tests[1].every)

	worker.jobs <- t1
	worker.jobs <- t2

	close(worker.jobs)

	time.Sleep(5 * time.Second)

	if len(t1.f.res) > 0 {
		res1 := t1.f.res[0].Interface().(int)
		if res1 != tests[0].result {
			t.Error("Got: ", res1, " Wanted: ", tests[0].result)
		}
	} else {
		t.Error("function results are empty")
	}

	if len(t2.f.res) > 0 {
		res2 := t2.f.res[0].Interface().(string)
		if res2 != tests[1].result {
			t.Error("Got: ", res2, " Wanted: ", tests[0].result)
		}
	} else {
		t.Error("function results are empty")
	}

}
