package goschedule

import (
	"fmt"
	"testing"
	"time"
)

func example1(a, b int, c string) {
	fmt.Println("Example 1: ", c, "is sum", a+b)
}

func example2() {
	fmt.Println("Example 2: Mahmoud")
}

type param struct {
	x int
	y int
}

func example3(p param) {
	fmt.Println(p.x + p.y)
}

func example4(a float32) {
	fmt.Println("Example 4: ", a)
}

func TestPriorityQueue(t *testing.T) {
	tNow := time.Now()
	j1 := newFunctionJob("example1", example1, []interface{}{1, 2, "Mahmoud"}).At(tNow.Add(2 * time.Minute)).Every(10 * time.Second)
	j2 := newFunctionJob("example2", example2, []interface{}{}).At(tNow.Add(1 * time.Minute)).Every(10 * time.Second)
	j3 := newFunctionJob("example3", example3, []interface{}{param{1, 2}}).At(tNow.Add(3 * time.Minute)).Every(10 * time.Second)
	j4 := newFunctionJob("example4", example4, []interface{}{1.9}).At(tNow.Add(40 * time.Second)).Every(10 * time.Second)

	j5 := newFunctionJob("example5", example1, []interface{}{1, 2, "Mahmoud"}).At(tNow.Add(2 * time.Minute)).Every(10 * time.Second)
	j6 := newFunctionJob("example6", example2, []interface{}{}).At(tNow.Add(12 * time.Minute)).Every(10 * time.Second)
	j7 := newFunctionJob("example7", example3, []interface{}{param{1, 2}}).At(tNow.Add(4 * time.Minute)).Every(10 * time.Second)
	j8 := newFunctionJob("example8", example4, []interface{}{1.9}).At(tNow.Add(5 * time.Second)).Every(10 * time.Second)

	resOrder := []string{
		"example8",
		"example4",
		"example2",
		"example1",
		"example5",
		"example3",
		"example7",
		"example6",
	}

	inputArray := []*Job{j1, j2, j3, j4, j6, j7, j8, j5}
	minHeap := newQueue(len(inputArray))
	minHeap.buildMinHeap()
	for i := 0; i < len(inputArray); i++ {
		minHeap.insert(inputArray[i])
	}

	for i := 0; i < len(inputArray); i++ {
		j := minHeap.remove()
		if j.GetIdentifier() != resOrder[i] {
			t.Error("Got: ", j.GetIdentifier(), " Wanted: ", resOrder[i])
		}
	}
}
