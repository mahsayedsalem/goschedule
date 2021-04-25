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
	}{
		{
			implmenetation: func(l int, b int) int { return l - b },
			params:         []interface{}{1, 2},
			result:         -1,
			identifier:     "test3",
		}, {
			implmenetation: func(l string, b string) string { return l + b },
			params:         []interface{}{"Mah", "moud"},
			result:         "Mahmoud",
			identifier:     "test4",
		},
	}

	worker := newWorker()
	go worker.work()

	t1 := newFunction(tests[0].implmenetation, tests[0].params, tests[0].identifier)
	t2 := newFunction(tests[1].implmenetation, tests[1].params, tests[1].identifier)

	worker.functions <- t1
	worker.functions <- t2

	close(worker.functions)

	time.Sleep(5 * time.Second)

	if len(t1.res) > 0 {
		res1 := t1.res[0].Interface().(int)
		if res1 != tests[0].result {
			t.Error("Got: ", res1, " Wanted: ", tests[0].result)
		}
	} else {
		t.Error("function results are empty")
	}

	if len(t2.res) > 0 {
		res2 := t2.res[0].Interface().(string)
		if res2 != tests[1].result {
			t.Error("Got: ", res2, " Wanted: ", tests[0].result)
		}
	} else {
		t.Error("function results are empty")
	}

}
