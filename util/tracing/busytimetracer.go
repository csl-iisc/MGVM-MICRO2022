package tracing

import (
	"sync"

	"gitlab.com/akita/akita"
)

type taskTimeStartEnd struct {
	start, end akita.VTimeInSec
}

// BusyTimeTracer traces the that a domain is processing a kind of task. If the
// task processing time overlaps, this tracer only consider one instance of the
// overlapped time.
type BusyTimeTracer struct {
	lock          sync.Mutex
	filter        TaskFilter
	inflightTasks map[string]Task
	taskTimes     []taskTimeStartEnd
}

// NewBusyTimeTracer creates a new BusyTimeTracer
func NewBusyTimeTracer(filter TaskFilter) *BusyTimeTracer {
	t := &BusyTimeTracer{
		filter:        filter,
		inflightTasks: make(map[string]Task),
	}
	return t
}

// BusyTime returns the total time has been spent on a certain type of tasks.
func (t *BusyTimeTracer) KernelTimes() []akita.VTimeInSec {
	kernelTimes := make([]akita.VTimeInSec, len(t.taskTimes))
	for i, t1 := range t.taskTimes {
		// if t1.What == "Launch Kernel" || t.What == "*driver.LaunchKernelCommand" {
		kernelTimes[i] = t1.end - t1.start
		// }
	}
	return kernelTimes
}

func (t *BusyTimeTracer) KernelTimesForceStop(now akita.VTimeInSec) []akita.VTimeInSec {
	kernelTimes := make([]akita.VTimeInSec, len(t.taskTimes))
	for i, t1 := range t.taskTimes {
		kernelTimes[i] = now - t1.start
		// }
	}
	return kernelTimes
}

// BusyTime returns the total time has been spent on a certain type of tasks.
func (t *BusyTimeTracer) BusyTime() akita.VTimeInSec {
	busyTime := akita.VTimeInSec(0.0)
	coveredMask := make(map[int]bool)

	for i, t1 := range t.taskTimes {
		if _, covered := coveredMask[i]; covered {
			continue
		}

		coveredMask[i] = true

		extTime := taskTimeStartEnd{
			start: t1.start,
			end:   t1.end,
		}

		for j, t2 := range t.taskTimes {
			if _, covered := coveredMask[j]; covered {
				continue
			}

			if t.taskTimeOverlap(t1, t2) {
				coveredMask[j] = true
				t.extendTaskTime(&extTime, t2)
			}
		}

		busyTime += extTime.end - extTime.start
	}

	return busyTime
}

// TerminateAllTasks will mark all the tasks as completed.
func (t *BusyTimeTracer) TerminateAllTasks(now akita.VTimeInSec) {
	for k, v := range t.inflightTasks {
		taskTime := taskTimeStartEnd{
			start: v.StartTime,
			end:   now,
		}
		t.taskTimes = append(t.taskTimes, taskTime)
		delete(t.inflightTasks, k)
	}
}

func (t *BusyTimeTracer) taskTimeOverlap(t1, t2 taskTimeStartEnd) bool {
	if t1.start <= t2.start && t1.end >= t2.start {
		return true
	}

	if t1.start <= t2.end && t1.end >= t2.end {
		return true
	}

	if t1.start >= t2.start && t1.end <= t2.end {
		return true
	}

	return false
}

func (t *BusyTimeTracer) extendTaskTime(
	base *taskTimeStartEnd,
	t2 taskTimeStartEnd,
) {
	if t2.start < base.start {
		base.start = t2.start
	}

	if t2.end > base.end {
		base.end = t2.end
	}
}

// StartTask records the task start time
func (t *BusyTimeTracer) StartTask(task Task) {
	if !t.filter(task) {
		return
	}

	t.inflightTasks[task.ID] = task
}

// StepTask does nothing
func (t *BusyTimeTracer) StepTask(task Task) {
	// Do nothing
}

// EndTask records the end of the task
func (t *BusyTimeTracer) EndTask(task Task) {
	originalTask, ok := t.inflightTasks[task.ID]
	if !ok {
		return
	}
	taskTime := taskTimeStartEnd{
		start: originalTask.StartTime,
		end:   task.EndTime,
	}
	t.taskTimes = append(t.taskTimes, taskTime)
	delete(t.inflightTasks, task.ID)
}
