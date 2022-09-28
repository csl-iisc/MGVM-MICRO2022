package tlb

import (
	// "fmt"
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem/device"
	"gitlab.com/akita/util"
	"gitlab.com/akita/util/pipelining"
)

// A IdealLatTLBBuilder can build LatTLBs
type IdealLatTLBBuilder struct {
	engine               akita.Engine
	freq                 akita.Freq
	numReqPerCycle       int
	latency              int
	pageTable            device.PageTable
	useCoalescingTLBPort bool
}

// MakeIdealLatTLBBuilder returns a IdealLatTLBBuilder
func MakeIdealLatTLBBuilder() IdealLatTLBBuilder {
	return IdealLatTLBBuilder{
		freq:           1 * akita.GHz,
		numReqPerCycle: 4,
		latency:        1,
	}
}

// WithEngine sets the engine that the LatTLBs to use
func (b IdealLatTLBBuilder) WithEngine(engine akita.Engine) IdealLatTLBBuilder {
	b.engine = engine
	return b
}

// WithFreq sets the freq the LatTLBs use
func (b IdealLatTLBBuilder) WithFreq(freq akita.Freq) IdealLatTLBBuilder {
	b.freq = freq
	return b
}

// // WithNumSets sets the number of sets in a LatTLB. Use 1 for fully associated
// // LatTLBs.
// func (b IdealLatTLBBuilder) WithNumSets(n int) IdealLatTLBBuilder {
// 	b.numSets = n
// 	return b
// }

// // WithNumWays sets the number of ways in a LatTLB. Set this field to the number
// // of LatTLB entries for all the functions.
// func (b IdealLatTLBBuilder) WithNumWays(n int) IdealLatTLBBuilder {
// 	b.numWays = n
// 	return b
// }

// // WithPageSize sets the page size that the LatTLB works with.
// func (b IdealLatTLBBuilder) WithPageSize(n uint64) IdealLatTLBBuilder {
// 	b.pageSize = n
// 	return b
// }

// WithNumReqPerCycle sets the number of requests per cycle can be processed by
// a LatTLB
func (b IdealLatTLBBuilder) WithNumReqPerCycle(n int) IdealLatTLBBuilder {
	b.numReqPerCycle = n
	return b
}

// // WithLowModule sets the port that can provide the address translation in case
// // of tlb miss.
// func (b IdealLatTLBBuilder) WithLowModule(lowModule akita.Port) IdealLatTLBBuilder {
// 	b.lowModule = lowModule
// 	return b
// }

// // WithNumMSHREntry sets the number of mshr entry
// func (b IdealLatTLBBuilder) WithNumMSHREntry(num int) IdealLatTLBBuilder {
// 	b.numMSHREntry = num
// 	return b
// }

// WithLatency sets the number of mshr entry
func (b IdealLatTLBBuilder) WithLatency(latency int) IdealLatTLBBuilder {
	b.latency = latency
	return b
}

func (b IdealLatTLBBuilder) WithPageTable(pt device.PageTable) IdealLatTLBBuilder {
	b.pageTable = pt
	return b
}

func (b IdealLatTLBBuilder) UseCoalescingTLBPort() IdealLatTLBBuilder {
	b.useCoalescingTLBPort = true
	return b
}

// Build creates a new LatTLB
func (b IdealLatTLBBuilder) Build(name string) L2TLB {
	tlb := &IdealLatTLB{}
	tlb.TickingComponent =
		akita.NewTickingComponent(name, b.engine, b.freq, tlb)

	tlb.numReqPerCycle = b.numReqPerCycle
	tlb.latency = b.latency
	if b.useCoalescingTLBPort {
		tlb.TopPort = NewCoalescingPort(tlb, 16*b.numReqPerCycle,
			name+".TopPort")
	} else {
		tlb.TopPort = akita.NewLimitNumMsgPort(tlb, 16*b.numReqPerCycle,
			name+".TopPort")
	}
	tlb.BottomPort = akita.NewLimitNumMsgPort(tlb, b.numReqPerCycle,
		name+".BottomPort")
	tlb.ControlPort = akita.NewLimitNumMsgPort(tlb, 1,
		name+".ControlPort")
	tlb.lookupBuffer = util.NewBuffer(2 * tlb.numReqPerCycle)
	pipelineBuilder := pipelining.MakeBuilder().WithPipelineWidth(tlb.numReqPerCycle).WithNumStage(tlb.latency).WithCyclePerStage(1).WithPostPipelineBuffer(tlb.lookupBuffer)
	tlb.pipeline = pipelineBuilder.Build(tlb.Name() + "_pipeline")
	// added this for consistency
	tlb.pageTable = b.pageTable
	// tlb.SetPageTable(b.pageTable)
	// tlb.reset()

	return tlb
}
