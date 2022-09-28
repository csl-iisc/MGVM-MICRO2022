package tlb

import (
	"strconv"
	"sync"

	"gitlab.com/akita/util/tracing"
)

type EntropyTracer struct {
	filter   tracing.TaskFilter
	lock     sync.Mutex
	pos_zero [64]float64
	pos_one  [64]float64
}

func NewEntropyTracer(filter tracing.TaskFilter) *EntropyTracer {
	t := &EntropyTracer{
		filter: filter,
	}
	return t
}

func (t *EntropyTracer) ReportEntropy() []float64 {
	entropy := make([]float64, 64)
	for i := 0; i < 64; i++ {
		if t.pos_one[i] == 0 && t.pos_zero[i] == 0 {
			entropy[i] = -1
		} else {
			entropy[i] = t.pos_one[i] / (t.pos_zero[i] + t.pos_one[i])
		}
	}
	return entropy
}

func (t *EntropyTracer) StartTask(task tracing.Task) {
	if !t.filter(task) {
		return
	}
	t.lock.Lock()
	addr, _ := strconv.Atoi(task.ID)
	for i := 0; i < 12; i++ {
		addr = addr >> 1
	}
	for j := 12; j < 64; j++ {
		bit := addr & 0x1
		if bit == 0 {
			t.pos_zero[j] += 1
		} else {
			t.pos_one[j] += 1
		}
		addr = addr >> 1
	}
	t.lock.Unlock()
}

func (t *EntropyTracer) StepTask(task tracing.Task) {
	return
}

func (t *EntropyTracer) EndTask(task tracing.Task) {
	return
}
