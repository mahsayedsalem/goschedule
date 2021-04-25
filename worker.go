package goschedule

import "sync"

type worker struct {
	jobs chan *Job
}

func newWorker() *worker {
	return &worker{
		jobs: make(chan *Job, 1),
	}
}

func (w *worker) work() {
	var wg sync.WaitGroup

	for j := range w.jobs {
		wg.Add(1)
		if j.isRabbitEvent {
			go j.rabbitEvent.publishEvent(&wg)
		} else {
			go j.f.runFunc(&wg)
		}
	}
	wg.Wait()
}
