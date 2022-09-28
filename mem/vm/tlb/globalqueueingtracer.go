package tlb

import (
	"sync"

	"gitlab.com/akita/util/tracing"
)

type GlobalTLBQueueingTracer struct {
	filter           tracing.TaskFilter
	lock             sync.Mutex
	this             L2TLB
	all_tlbs         []L2TLB
	averageImbalance float64
	imbalanceBuckets []uint64
	count            int
	threshold        int
}

func NewGlobalTLBQueueingTracer(filter tracing.TaskFilter) *GlobalTLBQueueingTracer {
	t := &GlobalTLBQueueingTracer{
		filter:           filter,
		averageImbalance: 0.0,
		count:            0,
		// TODO: fix
		threshold:        32,
		imbalanceBuckets: make([]uint64, 11),
	}
	t.all_tlbs = make([]L2TLB, 0)
	return t
}

// the TLB  corrosponding to this tracer
func (t *GlobalTLBQueueingTracer) AddThis(tlb L2TLB) {
	t.this = tlb
}

// list of all TLBs
func (t *GlobalTLBQueueingTracer) AddTLB(tlb L2TLB) {
	t.all_tlbs = append(t.all_tlbs, tlb)
}

func (t *GlobalTLBQueueingTracer) ReportImbalance() float64 {
	return t.averageImbalance
}

func (t *GlobalTLBQueueingTracer) ReportImbalanceBuckets() []uint64 {
	return t.imbalanceBuckets
}

func (t *GlobalTLBQueueingTracer) ReportCount() int {
	return t.count
}

func (t *GlobalTLBQueueingTracer) StartTask(task tracing.Task) {
	if !t.filter(task) {
		return
	}
	if t.this.GetFrontQueueLength() < t.threshold {
		return
	}
	t.lock.Lock()
	queuelengthsum := 0.0
	for _, tlb := range t.all_tlbs {
		// if t.this.Name() == tlb.Name() {
		// 	continue
		// }
		queuelengthsum += float64(tlb.GetFrontQueueLength())
	}
	imbalance := float64(t.this.GetFrontQueueLength()) / queuelengthsum
	t.averageImbalance = (t.averageImbalance*float64(t.count) + imbalance) / float64(t.count+1)
	t.imbalanceBuckets[int(imbalance*100)/10]++
	t.count += 1
	t.lock.Unlock()
}

func (t *GlobalTLBQueueingTracer) StepTask(task tracing.Task) {
	return
}

func (t *GlobalTLBQueueingTracer) EndTask(task tracing.Task) {
	return
}
