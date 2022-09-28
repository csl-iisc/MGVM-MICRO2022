package tracing

import (
	"fmt"
	"strconv"
	"sync"
)

// AverageCountTracer can collect the total time of executing a certain type of
// task. If the execution of two tasks overlaps, this tracer will simply add
// the two task processing time together.
type AverageCountTracer struct {
	filter        TaskFilter
	lock          sync.Mutex
	averageCount  float64
	inflightTasks map[string]Task
	taskCount     uint64
}

// NewAverageCountTracer creates a new AverageCountTracer
func NewAverageCountTracer(filter TaskFilter) *AverageCountTracer {
	t := &AverageCountTracer{
		filter:        filter,
		inflightTasks: make(map[string]Task),
	}
	return t
}

// AverageCount returns the total time has been spent on a certain type of tasks.
func (t *AverageCountTracer) AverageCount() float64 {
	t.lock.Lock()
	count := t.averageCount
	t.lock.Unlock()
	return count
}

// TotalCount returns the total number of tasks.
func (t *AverageCountTracer) TotalCount() uint64 {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.taskCount
}

// StartTask records the task start time
func (t *AverageCountTracer) StartTask(task Task) {
	if !t.filter(task) {
		return
	}
	t.lock.Lock()
	count, err := strconv.ParseInt(task.What, 10, 64)
	if err != nil {
		fmt.Println(count, err)
		panic("oh no!")
	}
	t.averageCount = (t.averageCount*float64(t.taskCount) + float64(count)) / (float64(t.taskCount) + 1)
	t.taskCount++
	t.lock.Unlock()
}

// StepTask does nothing
func (t *AverageCountTracer) StepTask(task Task) {
	// Do nothing
}

// EndTask records the end of the task
func (t *AverageCountTracer) EndTask(task Task) {
}
