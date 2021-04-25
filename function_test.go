package goschedule

import (
	"sync"
	"testing"
)

func TestFunction(t *testing.T) {
	tests := []struct {
		implmenetation interface{}
		params         []interface{}
		result         interface{}
		identifier     string
	}{
		{
			implmenetation: func(l int, b int) int { return l + b },
			params:         []interface{}{1, 2},
			result:         3,
			identifier:     "test1",
		},
		{
			implmenetation: func(l int, b int) int { return l * b },
			params:         []interface{}{1, 2},
			result:         2,
			identifier:     "test2",
		}, {
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

	var wg sync.WaitGroup

	for idx, test := range tests {
		f := newFunction(test.implmenetation, test.params, test.identifier)
		wg.Add(1)
		go f.runFunc(&wg)
		wg.Wait()
		var res int
		var res3 string
		if idx != 3 {
			res = f.res[0].Interface().(int)
			if res != test.result {
				t.Error("Got: ", res, " Wanted: ", test.result)
			}
		} else {
			res3 = f.res[0].Interface().(string)
			if res3 != test.result {
				t.Error("Got: ", res3, " Wanted: ", test.result)
			}
		}
	}
}
