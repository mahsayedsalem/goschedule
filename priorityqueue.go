package goschedule

import "fmt"

type PrioirityQueue struct {
	jobs []*Job
	size int
}

func newQueue(maxsize int) *PrioirityQueue {
	pq := &PrioirityQueue{
		jobs: []*Job{},
		size: 0,
	}
	return pq
}

func (m *PrioirityQueue) leaf(index int) bool {
	if index >= (m.size/2) && index <= m.size {
		return true
	}
	return false
}

func (m *PrioirityQueue) parent(index int) int {
	return (index - 1) / 2
}

func (m *PrioirityQueue) leftchild(index int) int {
	return 2*index + 1
}

func (m *PrioirityQueue) rightchild(index int) int {
	return 2*index + 2
}

func (m *PrioirityQueue) insert(newJob *Job) {
	m.jobs = append(m.jobs, newJob)
	m.size++
	m.upHeapify(m.size - 1)
}

func (m *PrioirityQueue) swap(first, second int) {
	temp := m.jobs[first]
	m.jobs[first] = m.jobs[second]
	m.jobs[second] = temp
}

func (m *PrioirityQueue) upHeapify(index int) {
	for m.jobs[index].GetNextRunUnixNanoTime() < m.jobs[m.parent(index)].GetNextRunUnixNanoTime() {
		m.swap(index, m.parent(index))
		index = m.parent(index)
	}
}

func (m *PrioirityQueue) downHeapify(current int) {
	if m.leaf(current) {
		return
	}
	smallest := current
	leftChildIndex := m.leftchild(current)
	rightRightIndex := m.rightchild(current)
	//If current is smallest then return
	if leftChildIndex < m.size && m.jobs[leftChildIndex].GetNextRunUnixNanoTime() < m.jobs[smallest].GetNextRunUnixNanoTime() {
		smallest = leftChildIndex
	}
	if rightRightIndex < m.size && m.jobs[rightRightIndex].GetNextRunUnixNanoTime() < m.jobs[smallest].GetNextRunUnixNanoTime() {
		smallest = rightRightIndex
	}
	if smallest != current {
		m.swap(current, smallest)
		m.downHeapify(smallest)
	}
	return
}
func (m *PrioirityQueue) buildMinHeap() {
	for index := ((m.size / 2) - 1); index >= 0; index-- {
		m.downHeapify(index)
	}
}

func (m *PrioirityQueue) printJobIdentifiersSorted() {
	for _, j := range m.jobs {
		fmt.Println(j.GetIdentifier())
	}
}

func (m *PrioirityQueue) remove() *Job {
	top := m.jobs[0]
	m.jobs[0] = m.jobs[m.size-1]
	m.jobs = m.jobs[:(m.size)-1]
	m.size--
	m.downHeapify(0)
	return top
}
