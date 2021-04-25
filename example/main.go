package main

import (
	"goschedule"
	"log"
	"time"

	"github.com/sirupsen/logrus"
)

func example1(a, b int, c string) {
	log.Println("Example 1: ", c, "is sum", a+b)
}

func example2(a, b int, c string) int {
	log.Println("Example 2: ", c, "is sum", a+b)
	return a + b
}

func main() {
	baseLogger := logrus.New()
	scheduler, err := goschedule.NewScheduler()
	if err != nil {
		baseLogger.Println(err)
	}
	scheduler.Start()
	j, err := scheduler.FuncJob("example1", example1, 1, 2, "Sum")
	if err != nil {
		baseLogger.Println(err)
	}
	j.At(time.Now().Add(5 * time.Second)).Every(5 * time.Second)
	scheduler.Schedule(j)

	j1, err := scheduler.FuncJob("example2", example2, 5, 2, "Sum")
	if err != nil {
		baseLogger.Println(err)
	}
	j1.At(time.Now().Add(4 * time.Second)).Every(6 * time.Second)
	scheduler.Schedule(j1)

	j2, err := scheduler.FuncJob("example3", example2, 5, 2, "Sum")
	if err != nil {
		baseLogger.Println(err)
	}
	j2.At(time.Now().Add(4 * time.Second))
	scheduler.Schedule(j2)

	time.Sleep(time.Second * 25)
	scheduler.Stop()

	baseLogger.Info(j.JobInfo().RanBefore, j.JobInfo().FunctionInfo.LatestExecutionTime, j.JobInfo().FunctionInfo.Results)
	baseLogger.Info(j1.JobInfo().RanBefore, j1.JobInfo().FunctionInfo.LatestExecutionTime, j1.JobInfo().FunctionInfo.Results)
}
