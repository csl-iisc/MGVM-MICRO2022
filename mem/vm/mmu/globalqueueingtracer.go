package mmu

import (
	"sync"

	"gitlab.com/akita/util/tracing"
)

type GlobalPageWalkerOccupancyTracer struct {
	filter           tracing.TaskFilter
	lock             sync.Mutex
	this             *MMUImpl
	all_mmus         []*MMUImpl
	averageImbalance float64
	count            int
	threshold        int
}

func NewGlobalPageWalkerOccupancyTracer(filter tracing.TaskFilter) *GlobalPageWalkerOccupancyTracer {
	t := &GlobalPageWalkerOccupancyTracer{
		filter:           filter,
		averageImbalance: 0.0,
		count:            0,
		// TODO: fix
		threshold: 4,
	}
	t.all_mmus = make([]*MMUImpl, 0)
	return t
}

// the Mmu  corrosponding to this tracer
func (t *GlobalPageWalkerOccupancyTracer) AddThis(mmu *MMUImpl) {
	t.this = mmu
}

// list of all Mmu
func (t *GlobalPageWalkerOccupancyTracer) AddMMU(mmu *MMUImpl) {
	t.all_mmus = append(t.all_mmus, mmu)
}

func (t *GlobalPageWalkerOccupancyTracer) ReportImbalance() float64 {
	return t.averageImbalance
}

func (t *GlobalPageWalkerOccupancyTracer) ReportCount() int {
	return t.count
}

func (t *GlobalPageWalkerOccupancyTracer) StartTask(task tracing.Task) {
	if !t.filter(task) {
		return
	}
	if t.this.GetNumActiveWalkers() < t.threshold {
		return
	}
	t.lock.Lock()
	queuelengthsum := 0.0
	for _, mmu := range t.all_mmus {
		// if t.this.Name() == tlb.Name() {
		// 	continue
		// }
		queuelengthsum += float64(mmu.GetNumActiveWalkers())
	}
	imbalance := float64(t.this.GetNumActiveWalkers()) / queuelengthsum
	t.averageImbalance = (t.averageImbalance*float64(t.count) + imbalance) / float64(t.count+1)
	t.count += 1
	t.lock.Unlock()
}

func (t *GlobalPageWalkerOccupancyTracer) StepTask(task tracing.Task) {
	return
}

func (t *GlobalPageWalkerOccupancyTracer) EndTask(task tracing.Task) {
	return
}
