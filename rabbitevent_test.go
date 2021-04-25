package goschedule

import (
	"log"
	"os"
	"sync"
	"testing"

	"github.com/streadway/amqp"
)

func TestRabbitEvent(t *testing.T) {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to rabbitmq: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}
	body := []byte("Hello World!")
	exchange := ""
	routingKey := q.Name
	mandatory := false
	immediate := false
	contentType := "text/plain"
	var wg sync.WaitGroup
	rabbitEvent := newRabbitMQEvent(ch, exchange, routingKey, mandatory, immediate, contentType, body, "example1")
	wg.Add(1)
	go rabbitEvent.publishEvent(&wg)
	wg.Wait()
	eventInfo := rabbitEvent.GetRabbitEventInfo()
	if eventInfo.LatestPublishStatus != "success" {
		t.Error("Got: ", eventInfo.LatestPublishStatus, " Wanted: success")
	}

}
