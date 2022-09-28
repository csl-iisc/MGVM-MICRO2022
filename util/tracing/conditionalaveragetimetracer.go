package tracing

import (
	"sort"
	"strings"
	"sync"

	//"fmt"

	"gitlab.com/akita/akita"
)

// ConditionalAverageTimeTracer can collect the total time of executing a certain type of
// task. If the execution of two tasks overlaps, this tracer will simply add
// the two task processing time together.
type ConditionalAverageTimeTracer struct {
	filter        TaskFilter
	lock          sync.Mutex
	averageTime   map[string]akita.VTimeInSec
	inflightTasks map[string]Task
	taskCount     map[string]uint64
}

// NewConditionalAverageTimeTracer creates a new ConditionalAverageTimeTracer
func NewConditionalAverageTimeTracer(filter TaskFilter) *ConditionalAverageTimeTracer {
	t := &ConditionalAverageTimeTracer{
		filter:        filter,
		inflightTasks: make(map[string]Task),
		averageTime:   make(map[string]akita.VTimeInSec),
		taskCount:     make(map[string]uint64),
	}
	return t
}

// AverageTime returns the total time has been spent on a certain type of tasks.
func (t *ConditionalAverageTimeTracer) AverageTime(taskType string) akita.VTimeInSec {
	t.lock.Lock()
	time := t.averageTime[taskType]
	t.lock.Unlock()
	return time
}

// TotalCount returns the total number of tasks.
func (t *ConditionalAverageTimeTracer) TotalCount(taskType string) uint64 {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.taskCount[taskType]
}

// StartTask records the task start time
func (t *ConditionalAverageTimeTracer) StartTask(task Task) {
	if !t.filter(task) {
		return
	}

	t.lock.Lock()
	t.inflightTasks[task.ID] = task
	t.lock.Unlock()
}

// StepTask does nothing
func (t *ConditionalAverageTimeTracer) StepTask(task Task) {
	// Do nothing
	t.lock.Lock()
	originalTask, ok := t.inflightTasks[task.ID]
	if !ok {
		if strings.Contains(task.ID, "addr-translator-stats") {
			panic("Task not found!")
		}
		t.lock.Unlock()
		return
	}
	taskType := task.Steps[0].What
	//fmt.Println(taskType)
	taskTime := task.Steps[0].Time - originalTask.StartTime
	// if taskType == "translation-latency" && strings.Contains(originalTask.Where, "L1VAddrTrans") {
	// fmt.Println(taskType, originalTask.Where, taskTime*1000000000)
	// }
	t.averageTime[taskType] = akita.VTimeInSec(
		(float64(t.averageTime[taskType])*float64(t.taskCount[taskType]) + float64(taskTime)) /
			float64(t.taskCount[taskType]+1))
	t.taskCount[taskType]++
	t.lock.Unlock()
}

// EndTask records the end of the task
func (t *ConditionalAverageTimeTracer) EndTask(task Task) {

	t.lock.Lock()
	delete(t.inflightTasks, task.ID)
	t.lock.Unlock()
}

// GetStepNames returns all the step names collected.
func (t *ConditionalAverageTimeTracer) GetStepNames() (keys []string) {
	keys = make([]string, len(t.averageTime))
	i := 0
	for k := range t.averageTime {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return
}
