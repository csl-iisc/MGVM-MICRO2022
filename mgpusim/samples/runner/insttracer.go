package runner

import (
	"github.com/tebeka/atexit"
	"gitlab.com/akita/akita"
	"gitlab.com/akita/util/tracing"
)

// instTracer can trace the number of instruction completed.
type instTracer struct {
	count    uint64
	maxCount uint64

	lastWGFinishTime akita.VTimeInSec

	inflightInst map[string]tracing.Task
}

// newInstTracer creates a tracer that can count the number of instructions.
func newInstTracer() *instTracer {
	t := &instTracer{
		inflightInst: map[string]tracing.Task{},
	}
	return t
}

// newInstStopper with stop the execution after a given number of instructions
// is retired.
func newInstStopper(maxInst uint64) *instTracer {
	t := &instTracer{
		maxCount:     maxInst,
		inflightInst: map[string]tracing.Task{},
	}
	return t
}

func (t *instTracer) StartTask(task tracing.Task) {
	if task.Kind != "inst" {
		return
	}

	t.inflightInst[task.ID] = task
}

func (t *instTracer) StepTask(task tracing.Task) {
	// Do nothing
}

func (t *instTracer) EndTask(task tracing.Task) {
	_, found := t.inflightInst[task.ID]
	if !found {
		return
	}

	delete(t.inflightInst, task.ID)

	t.count++

	t.lastWGFinishTime = task.EndTime

	if t.maxCount > 0 && t.count >= t.maxCount {
		atexit.Exit(0)
	}
}
