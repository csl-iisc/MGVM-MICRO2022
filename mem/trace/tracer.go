// Package trace provides a tracer that can trace memory system tasks.
package trace

import (
	"log"

	"gitlab.com/akita/mem"
	"gitlab.com/akita/util/tracing"
)

// A tracer is a hook that can record the actions of a memory model into
// traces.
type tracer struct {
	logger *log.Logger
}

// StartTask marks the start of a memory transaction
func (t *tracer) StartTask(task tracing.Task) {
	req, ok := task.Detail.(mem.AccessReq)
	if !ok {
		return
	}
	t.logger.Printf("start, %.12f, %s, %s, %s, 0x%x, %d\n",
		task.StartTime, task.Where, task.ID, task.What,
		req.GetAddress(), req.GetByteSize())
}

// StepTask marks the memory transaction has completed a milestone
func (t *tracer) StepTask(task tracing.Task) {
	t.logger.Printf("step, %.12f, %s, %s\n",
		task.Steps[0].Time,
		task.ID,
		task.Steps[0].What)
}

// EndTask marks the end of a memory transaction
func (t *tracer) EndTask(task tracing.Task) {
	t.logger.Printf("end, %.12f, %s\n", task.EndTime, task.ID)
}

// NewTracer creates a new Tracer.
func NewTracer(logger *log.Logger) tracing.Tracer {
	t := new(tracer)
	t.logger = logger
	return t
}
