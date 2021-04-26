<h1 align="center">
  <br>
  <img src="assets/logo.jpg" alt="Markdownify" width="300">
  <br>
  goschedule
  <br>
</h1>

<h4 align="center">An easy-to-use in-process scheduler to schedule functions and events.</h4>

<p align="center">
  <a href="#key-features">Key Features</a> •
  <a href="#how-to-use">How To Use</a> •
  <a href="#how-it-works">How It Works</a> •
  <a href="#usage">Usage</a> •
  <a href="#contribution-priority">Contributions Priority</a> •
  <a href="#inspired-from">Inspired From</a> •
  <a href="#license">License</a>
</p>

## Key Features

* Schedule functions to run once or periodically at a time and interval of your choice.
* Schedule RabbitMQ events to be published once or periodically at a time and interval of your choice.
* Logs for every job run.  
* Ability to view information about scheduled jobs (Functions and Events).
* Even if a function/event breaks during execution, the scheduler will remain running while updating the user that something broke.
* Events can be scheduled to be published and can be consumed by any service subscribing to the queue; hence, you're not stuck with Go in case you needed any other language.

## How To Use

1. Define number of Jobs and Workers.
2. Create a Scheduler and start it. 
3. Choose what kind of job do you want; Function or Event.
4. Schedule your Job. 
5. Monitor your Jobs while they are running.

## How It Works

1. Jobs are queued into a priorityqueue based on their next planning running time. 
2. When a job is due it's sent to a worker.
3. The worker executes the job.
4. Even if a scheduler stops. The worker will wait until all the jobs he executed finish running before it terminate itself.
5. Jobs are distributed among workers fairly. 
6. Everything works asynchronously.
7. Jobs are shared between priority-queue and workers / sheduler and priority-queue through go channels.
8. Special care was taken to avoid goroutine leakages.
## Usage

### Install

```sh
$ go get https://github.com/mahsayedsalem/goschedule
```

```
import "github.com/mahsayedsalem/goschedule"
```

### Create a Scheduler

```
// function to attach to the job
func example1(a, b int, c string) {
	log.Println("Example 1: ", c, "is sum", a+b)
}

func main(){
    baseLogger := logrus.New()
    numberWorkers, err := strconv.Atoi(os.Getenv("WORKERSNUMBER"))
    failOnError(err, baseLogger)
    maxJobs, err := strconv.Atoi(os.Getenv("MAXJOBS"))
    failOnError(err, baseLogger)
    scheduler := goschedule.NewScheduler(maxJobs, numberWorkers)
    scheduler.Start()
}
```

### Create a Job
```
    // params are (unique identifier, function signature, ...paramaters)
    j, err := scheduler.FuncJob("exampleFunc1", example1, 1, 2, "Sum")
    failOnError(err, baseLogger)
```

### Schedule a Job
```
    j.At(time.Now().Add(5 * time.Second)).Every(5 * time.Second)
    scheduler.Schedule(j)
```

### Schedule an Event
```
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
    j3, err := scheduler.EventJob("exampleEvent1", ch,  exchange, routingKey, mandatory, immediate,   contentType, body)
    failOnError(err, baseLogger)
    j3.At(time.Now().Add(10 * time.Second)).Every(5 * time  Second)
    scheduler.Schedule(j3)
```

### Get info about the Job

```
    baseLogger.Info(j.JobInfo().RanBefore, j.JobInfo(). FunctionInfo.LatestExecutionTime, j.JobInfo().FunctionInfo.  Results)
```

### Safely stop the Scheduler

```
    scheduler.Stop()
```

### You have control

We don't create the functions or rabbit servers/events/queue for you. You supply us with them. We just safely run them at the scheduled time.

You're responsible for setting up the functions properly or testing RabbitMQ connectivity.

## Contributions Priority

- [ ] Use Redis to avoid pririty queue race conditions if ran on a distributed system. 
- [ ] Use Redis to store jobs unique identifiers and parameters, and ensure identifiers aren't duplicated.
- [ ] Consider time-zones when scheduling jobs.
- [ ] Edit/Remove Jobs.
- [ ] Central Logging.
- [ ] Add more message brokers (Currently RabbitMQ is supperted, Kafka will be integratetd next).
- [ ] Create a circuit-breaker to retry failed jobs a number of times then remove it from the queue until updated.  
- [ ] Implement a dead-job-queue to store failing jobs.


## Inspired from

1. https://github.com/jasonlvhit/gocron

## License

MIT
