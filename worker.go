package goschedule

import "sync"

type worker struct {
	functions chan *function
}

func newWorker() *worker {
	return &worker{
		functions: make(chan *function, 1),
	}
}

func (w *worker) work() {
	var wg sync.WaitGroup

	for f := range w.functions {
		wg.Add(1)
		go f.runFunc(&wg)
	}
	wg.Wait()
}
