package goschedule

import (
	"sync"

	"github.com/streadway/amqp"
)

type rabbitMQEvent struct {
	ch                  *amqp.Channel
	exchange            string
	routingKey          string
	mandatory           bool
	immediate           bool
	contentType         string
	body                []byte
	identifier          string
	latestPublishFailed bool
	failedReason        string
}

type rabbitMQEventInfo struct {
	LatestPublishFailed string
	FailedReason        string
}

func newRabbitMQEvent(ch *amqp.Channel, exchange string, routingKey string, mandatory bool, immediate bool, contentType string, body []byte, identifier string) *rabbitMQEvent {
	return &rabbitMQEvent{
		ch:          ch,
		exchange:    exchange,
		routingKey:  routingKey,
		mandatory:   mandatory,
		immediate:   immediate,
		contentType: contentType,
		body:        body,
		identifier:  identifier,
	}
}

func (r *rabbitMQEvent) publishEvent(wg *sync.WaitGroup) {
	err := r.ch.Publish(
		r.exchange,
		r.routingKey,
		r.mandatory,
		r.immediate,
		amqp.Publishing{
			ContentType: r.contentType,
			Body:        r.body,
		},
	)
	if err != nil {
		r.latestPublishFailed = true
		r.failedReason = err.Error()
	}
}

func (r *rabbitMQEvent) GetRabbitEventInfo() *rabbitMQEventInfo {
	status := "success"
	if r.latestPublishFailed {
		status = "fail"
	}
	return &rabbitMQEventInfo{
		LatestPublishFailed: status,
		FailedReason:        r.failedReason,
	}
}
