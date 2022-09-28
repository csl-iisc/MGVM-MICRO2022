package tracing

import (
	"math"
	"sync"

	"gitlab.com/akita/akita"
)

// Count references to
type ReferenceTracer struct {
	filter            TaskFilter
	lock              sync.Mutex
	reference         map[string]uint64
	lastAccess        map[string]akita.VTimeInSec
	refInterval       map[string]akita.VTimeInSec
	refIntervalHist   map[string]map[uint]uint
	intervalhistogram [1024]uint
	counthistogram    [32]uint
}

func NewReferenceTracer(filter TaskFilter) *ReferenceTracer {
	t := &ReferenceTracer{
		filter:          filter,
		reference:       make(map[string]uint64),
		lastAccess:      make(map[string]akita.VTimeInSec),
		refInterval:     make(map[string]akita.VTimeInSec),
		refIntervalHist: make(map[string]map[uint]uint),
	}
	return t
}

func (t *ReferenceTracer) StartTask(task Task) {
	if !t.filter(task) {
		return
	}

	// shamelessly repurposing What to send the VPN
	address := task.What

	t.lock.Lock()

	if val, ok := t.lastAccess[address]; ok {
		delta := task.StartTime - val
		t.refInterval[address] = ((t.refInterval[address] * akita.VTimeInSec(t.reference[address])) + delta) / akita.VTimeInSec(t.reference[address]+1)

		deltaInCycles := float64(delta) * 1 * float64(akita.GHz)
		logDeltaInCycles := uint(math.Log2(math.Ceil(float64(deltaInCycles))))
		// 2^300 is largeeee. delta must be zero, but delta is floating pt.
		if logDeltaInCycles > 300 {
			logDeltaInCycles = 0
		}

		t.refIntervalHist[address][logDeltaInCycles] += 1
	} else {
		t.refIntervalHist[address] = make(map[uint]uint)
	}

	t.reference[address] += 1
	t.lastAccess[address] = task.StartTime // API abuse

	t.lock.Unlock()
}

func (t *ReferenceTracer) StepTask(task Task) {
}

func (t *ReferenceTracer) EndTask(task Task) {
}

func (t *ReferenceTracer) ReturnAverageReferenceCount() float64 {
	average := uint64(0)
	count := 0
	for ref := range t.reference {
		average = average + t.reference[ref]
		count += 1
	}
	if count == 0 {
		return 0.0
	}
	return float64(average) / float64(count)
}

func (t *ReferenceTracer) ReturnAverageReferenceInterval() float64 {
	average := akita.VTimeInSec(0)
	count := 0
	for ref := range t.refInterval {
		average = average + t.refInterval[ref]
		count += 1
	}
	if count == 0 {
		return 0.0
	}
	return float64(average) / float64(count)
}

func (t *ReferenceTracer) CalculateHistogramSummary() {
	for ref := range t.refIntervalHist {
		for interval := range t.refIntervalHist[ref] {
			if interval < 1024 {
				t.intervalhistogram[interval] += 1
			}
		}
	}
	for ref := range t.reference {
		logcount := uint(math.Log2(math.Ceil(float64(t.reference[ref]))))
		// overflow condition
		if logcount > 32 {
			logcount = 0
		}
		t.counthistogram[logcount] += 1
	}
}

func (t *ReferenceTracer) ReturnHistogramSummary(interval uint) uint {
	return t.intervalhistogram[interval]
}

func (t *ReferenceTracer) ReturnCountHistogramSummary(count uint) uint {
	return t.counthistogram[count]
}
