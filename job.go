package goschedule

import (
	"time"

	"github.com/streadway/amqp"
)

type Job struct {
	identifier    string
	isPeriodic    bool
	latestRunAt   time.Time
	at            time.Time
	every         time.Duration
	firstRun      bool
	f             *function
	isRabbitEvent bool
	rabbitEvent   *rabbitMQEvent
}

type JobInfo struct {
	Identifier        string
	LatestRunAt       time.Time
	NextRunAt         time.Time
	Every             time.Duration
	RanBefore         bool
	FunctionInfo      *functionInfo
	IsRabbitEvent     bool
	RabbitMQEventInfo *rabbitMQEventInfo
}

func newFunctionJob(identifier string, f interface{}, params []interface{}) *Job {
	return &Job{
		identifier:    identifier,
		isPeriodic:    false,
		firstRun:      true,
		f:             newFunction(f, params, identifier),
		isRabbitEvent: false,
	}
}

func newEventJob(identifier string, ch *amqp.Channel, exchange string, routingKey string, mandatory bool, immediate bool, contentType string, body []byte) *Job {
	e := newRabbitMQEvent(ch, exchange, routingKey, mandatory, immediate, contentType, body, identifier)
	return &Job{
		identifier:    identifier,
		isPeriodic:    false,
		firstRun:      true,
		rabbitEvent:   e,
		isRabbitEvent: true,
	}
}

func (j *Job) At(at time.Time) *Job {
	j.at = at
	return j
}

func (j *Job) Every(every time.Duration) *Job {
	j.isPeriodic = true
	j.every = every
	return j
}

func (j *Job) GetIdentifier() string {
	return j.identifier
}

func (j *Job) GetNextRunTime() time.Time {
	return j.at
}

func (j *Job) GetIntervalsBetweenRuns() time.Duration {

	return j.every
}

func (j *Job) GetNextRunUnixNanoTime() int64 {
	return j.at.UnixNano()
}

func (j *Job) IsPeriodic() bool {
	return j.isPeriodic
}

func (j *Job) shouldRun() bool {
	return time.Now().After(j.at)
}

func (j *Job) updateForNextRun() {
	j.at = time.Now().Add(j.every)
}

func (j *Job) JobInfo() *JobInfo {
	if j.isRabbitEvent {
		return &JobInfo{
			Identifier:        j.identifier,
			LatestRunAt:       j.latestRunAt,
			NextRunAt:         j.at,
			Every:             j.every,
			RanBefore:         !j.firstRun,
			IsRabbitEvent:     j.isRabbitEvent,
			RabbitMQEventInfo: j.rabbitEvent.GetRabbitEventInfo(),
		}
	}
	return &JobInfo{
		Identifier:    j.identifier,
		LatestRunAt:   j.latestRunAt,
		NextRunAt:     j.at,
		Every:         j.every,
		RanBefore:     !j.firstRun,
		FunctionInfo:  j.f.GetFuncInfo(),
		IsRabbitEvent: j.isRabbitEvent,
	}
}
