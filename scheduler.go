package goschedule

import (
	"errors"
	"os"
	"strconv"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Scheduler struct {
	queue         *PrioirityQueue
	currentJobs   map[string]*Job
	jobs          chan *Job
	quit          chan bool
	workers       []*worker
	numberWorkers int
	maxJobs       int
	isRunning     bool
	mu            sync.Mutex
}

func NewScheduler() (*Scheduler, error) {
	numberWorkers, err := strconv.Atoi(os.Getenv("WORKERSNUMBER"))
	if err != nil {
		return nil, errors.New("Workers aren't set or not int")
	}
	maxJobs, err := strconv.Atoi(os.Getenv("MAXJOBS"))
	if err != nil {
		return nil, errors.New("Maximum number of jobs aren't set or not int")
	}

	return &Scheduler{
		queue:         newQueue(maxJobs),
		currentJobs:   make(map[string]*Job),
		numberWorkers: numberWorkers,
		maxJobs:       maxJobs,
	}, nil
}

func (s *Scheduler) FuncJob(identifier string, f interface{}, params ...interface{}) (*Job, error) {
	if _, ok := s.currentJobs[identifier]; ok {
		return nil, errors.New("Identifiers must be unique for each job")
	}
	j := newFunctionJob(identifier, f, params)
	s.currentJobs[identifier] = j
	return j, nil
}

func (s *Scheduler) Schedule(job *Job) error {
	if !s.IsRunning() {
		log.Error("The scheduler isn't running")
		return errors.New("The scheduler isn't running")
	}
	if job.at.IsZero() {
		log.Error("No time scheduled for this job")
		return errors.New("No time scheduled for this job")
	}
	if job.IsPeriodic() && job.GetIntervalsBetweenRuns() == 0 {
		log.Error("Duration must be set in case of periodic job")
		return errors.New("Duration must be set")
	}
	now := time.Now()
	if !job.IsPeriodic() && job.GetNextRunTime().Before(now) {
		log.Error("Single run can't be in the past")
		return errors.New("Single run can't be in the past")
	}

	for job.at.Before(now) {
		job.at = job.at.Add(job.every)
	}
	s.mu.Lock()
	s.queue.insert(job)
	s.mu.Unlock()
	return nil
}

func (s *Scheduler) asyncAllocate() {
	for j := range s.jobs {
		log.Info("Allocating job ", j.GetIdentifier(), " to a worker")
		if len(s.workers) > 1 {
			w := s.workers[0]
			w.functions <- j.f
			s.workers = s.workers[1:]
			s.workers = append(s.workers, w)
		}
	}
}

func (s *Scheduler) asyncCheckJob() {
	for {
		select {
		case <-s.quit:
			return
		default:
			if s.queue.size > 0 && s.IsRunning() {
				s.mu.Lock()
				earilestJob := s.queue.remove()
				if earilestJob.shouldRun() {
					log.Info("Job ", earilestJob.GetIdentifier(), " is due")
					s.jobs <- earilestJob
					earilestJob.latestRunAt = time.Now()
					if earilestJob.firstRun {
						earilestJob.firstRun = false
					}
					if earilestJob.isPeriodic {
						earilestJob.updateForNextRun()
					}
				}
				if earilestJob.isPeriodic {
					s.queue.insert(earilestJob)
				}
				s.mu.Unlock()
			}
		}
	}
}

func (s *Scheduler) IsRunning() bool {
	return s.isRunning
}

func (s *Scheduler) Stop() {
	log.Info("Scheduler is stopping")
	if !s.IsRunning() {
		return
	}
	s.isRunning = false
	for _, w := range s.workers {
		close(w.functions)
	}
	s.quit <- true
	close(s.quit)
	close(s.jobs)
	log.Info("Scheduler Stopped")
}

func (s *Scheduler) Start() {
	log.Info("Scheduler is starting")
	var workers []*worker
	for i := 0; i < s.numberWorkers; i++ {
		workers = append(workers, newWorker())
	}
	s.workers = workers
	for _, worker := range s.workers {
		go worker.work()
	}
	s.jobs = make(chan *Job)
	s.quit = make(chan bool)
	go s.asyncAllocate()
	go s.asyncCheckJob()
	s.queue.buildMinHeap()
	s.isRunning = true
	log.Info("Scheduler started")
	time.Sleep(time.Second * 1)
}

func (s *Scheduler) GetJobInfo(identifier string) (*functionInfo, error) {
	if _, ok := s.currentJobs[identifier]; !ok {
		log.Error("No job with this identifier")
		return nil, errors.New("No job with this identifier")
	}
	j := s.currentJobs[identifier]
	return j.f.GetFuncInfo(), nil
}
