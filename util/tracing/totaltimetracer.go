package tracing

import (
	"sync"

	"gitlab.com/akita/akita"
)

// TotalTimeTracer can collect the total time of executing a certain type of
// task. If the execution of two tasks overlaps, this tracer will simply add
// the two task processing time together.
type TotalTimeTracer struct {
	filter        TaskFilter
	lock          sync.Mutex
	totalTime     akita.VTimeInSec
	inflightTasks map[string]Task
}

// NewTotalTimeTracer creates a new TotalTimeTracer
func NewTotalTimeTracer(filter TaskFilter) *TotalTimeTracer {
	t := &TotalTimeTracer{
		filter:        filter,
		inflightTasks: make(map[string]Task),
	}
	return t
}

// TotalTime returns the total time has been spent on a certain type of tasks.
func (t *TotalTimeTracer) TotalTime() akita.VTimeInSec {
	t.lock.Lock()
	time := t.totalTime
	t.lock.Unlock()
	return time
}

// StartTask records the task start time
func (t *TotalTimeTracer) StartTask(task Task) {
	if !t.filter(task) {
		return
	}

	t.lock.Lock()
	t.inflightTasks[task.ID] = task
	t.lock.Unlock()
}

// StepTask does nothing
func (t *TotalTimeTracer) StepTask(task Task) {
	// Do nothing
}

// EndTask records the end of the task
func (t *TotalTimeTracer) EndTask(task Task) {
	t.lock.Lock()
	originalTask, ok := t.inflightTasks[task.ID]
	if !ok {
		t.lock.Unlock()
		return
	}

	t.totalTime += task.EndTime - originalTask.StartTime
	delete(t.inflightTasks, task.ID)
	t.lock.Unlock()
}
