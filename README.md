# GOSCHEDULE

<img src="assets/logo.jpg" align="right"
     alt="GOSCHEDULE by Salem" width="200" height="200">

GOSCHEDULE is an easy-to-understand in-process scheduler to schedule functions and events.
## How To Use It

1. Define number of Jobs and Workers.
2. Create a Scheduler and start it. 
3. Choose what kind of job do you want; Function or Event.
4. Schedule your Job. 
5. Monitor your Jobs while they are running.

## How It Works

1. Jobs are queued into a priorityqueue based on their next running time. 
2. When a job is due it's sent to a worker.
3. The worker works the job. 
4. Everything happens asynchronously.
## Usage

### Install

```sh
$ go get https://github.com/mahsayedsalem/goschedule
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

### You have control

We don't create the functions or rabbit servers/events/queue for you. You supply us with them. We just safely run them at the scheduled time. 

## Contributions Priority

1. Create distributed logging.
2. Use Redis to avoid race conditions if ran on a distributed system. 
3. Use time-zones for scheduling jobs.
4. More test cases!
5. Add more message brokers (Currently only RabbitMQ is supperted)

## Inspired from

1. https://github.com/jasonlvhit/gocron
