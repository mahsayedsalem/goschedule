package goschedule

import (
	"fmt"
	"testing"
	"time"
)

func example1_job(a, b int, c string) {
	fmt.Println("Example 1: ", c, "is sum", a+b)
}

func TestJob(t *testing.T) {
	tNow := time.Now()
	nextRun := tNow.Add(2 * time.Minute)
	j1 := newFunctionJob("example1", example1_job, []interface{}{1, 2, "Mahmoud"}).At(nextRun).Every(10 * time.Second)
	info := j1.JobInfo()
	if info.Identifier != "example1" {
		t.Error("Got: ", info.Identifier, " Wanted: ", "example1")
	}
	if info.NextRunAt != nextRun {
		t.Error("Got: ", info.NextRunAt, " Wanted: ", nextRun)
	}
	if info.RanBefore != false {
		t.Error("Got: ", info.RanBefore, " Wanted: ", false)
	}
	if info.Every != 10*time.Second {
		t.Error("Got: ", info.Every, " Wanted: ", 10*time.Second)
	}
}
