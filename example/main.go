package main

import (
	"goschedule"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func failOnError(err error, baseLogger *logrus.Logger) {
	if err != nil {
		baseLogger.Fatalf(err.Error())
	}
}

func example1(a, b int, c string) {
	log.Println("Example 1: ", c, "is sum", a+b)
}

func example2(a, b int, c string) int {
	log.Println("Example 2: ", c, "is sum", a+b)
	return a + b
}

func main() {
	baseLogger := logrus.New()
	numberWorkers, err := strconv.Atoi(os.Getenv("WORKERSNUMBER"))
	failOnError(err, baseLogger)
	maxJobs, err := strconv.Atoi(os.Getenv("MAXJOBS"))
	failOnError(err, baseLogger)

	scheduler := goschedule.NewScheduler(maxJobs, numberWorkers)

	scheduler.Start()
	j, err := scheduler.FuncJob("exampleFunc1", example1, 1, 2, "Sum")
	failOnError(err, baseLogger)
	j.At(time.Now().Add(5 * time.Second)).Every(5 * time.Second)
	scheduler.Schedule(j)

	j1, err := scheduler.FuncJob("exampleFunc2", example2, 5, 2, "Sum")
	failOnError(err, baseLogger)
	j1.At(time.Now().Add(4 * time.Second)).Every(6 * time.Second)
	scheduler.Schedule(j1)

	j2, err := scheduler.FuncJob("example3", example2, 5, 2, "Sum")
	failOnError(err, baseLogger)
	j2.At(time.Now().Add(4 * time.Second))
	scheduler.Schedule(j2)

	// schedule event job
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	failOnError(err, baseLogger)
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, baseLogger)
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, baseLogger)
	body := []byte("Hello World 1!")
	exchange := ""
	routingKey := q.Name
	mandatory := false
	immediate := false
	contentType := "text/plain"

	j3, err := scheduler.EventJob("exampleEvent1", ch, exchange, routingKey, mandatory, immediate, contentType, body)
	failOnError(err, baseLogger)
	j3.At(time.Now().Add(10 * time.Second)).Every(5 * time.Second)
	scheduler.Schedule(j3)
	body = []byte("Hello World 2!")
	j4, err := scheduler.EventJob("exampleEvent2", ch, exchange, routingKey, mandatory, immediate, contentType, body)
	failOnError(err, baseLogger)
	j4.At(time.Now().Add(8 * time.Second)).Every(4 * time.Second)
	scheduler.Schedule(j4)

	time.Sleep(time.Second * 25)
	scheduler.Stop()

	baseLogger.Info(j.JobInfo().RanBefore, j.JobInfo().FunctionInfo.LatestExecutionTime, j.JobInfo().FunctionInfo.Results)
	baseLogger.Info(j1.JobInfo().RanBefore, j1.JobInfo().FunctionInfo.LatestExecutionTime, j1.JobInfo().FunctionInfo.Results)
}
