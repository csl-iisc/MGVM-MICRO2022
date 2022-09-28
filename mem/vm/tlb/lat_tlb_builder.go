package tlb

import (
	"math"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/util"
	"gitlab.com/akita/util/pipelining"
)

// A LatTLBBuilder can build LatTLBs
type LatTLBBuilder struct {
	engine               akita.Engine
	freq                 akita.Freq
	numReqPerCycle       int
	numSets              int
	numWays              int
	pageSize             uint64
	log2PageSize         uint64
	lowModule            akita.Port
	latency              int
	numMSHREntry         int
	mask                 uint64
	useCoalescingTLBPort bool
}

// MakeLatTLBBuilder returns a LatTLBBuilder
func MakeLatTLBBuilder() LatTLBBuilder {
	return LatTLBBuilder{
		freq:           1 * akita.GHz,
		numReqPerCycle: 4,
		numSets:        1,
		numWays:        32,
		pageSize:       4096,
		numMSHREntry:   4,
		latency:        1,
		mask:           uint64(127) << 12,
	}
}

// WithEngine sets the engine that the LatTLBs to use
func (b LatTLBBuilder) WithEngine(engine akita.Engine) LatTLBBuilder {
	b.engine = engine
	return b
}

// WithFreq sets the freq the LatTLBs use
func (b LatTLBBuilder) WithFreq(freq akita.Freq) LatTLBBuilder {
	b.freq = freq
	return b
}

// WithNumSets sets the number of sets in a LatTLB. Use 1 for fully associated
// LatTLBs.
func (b LatTLBBuilder) WithNumSets(n int) LatTLBBuilder {
	b.numSets = n
	return b
}

// WithNumWays sets the number of ways in a LatTLB. Set this field to the number
// of LatTLB entries for all the functions.
func (b LatTLBBuilder) WithNumWays(n int) LatTLBBuilder {
	b.numWays = n
	return b
}

// WithPageSize sets the page size that the LatTLB works with.
func (b LatTLBBuilder) WithPageSize(n uint64) LatTLBBuilder {
	b.pageSize = n
	return b
}

// WithNumReqPerCycle sets the number of requests per cycle can be processed by
// a LatTLB
func (b LatTLBBuilder) WithNumReqPerCycle(n int) LatTLBBuilder {
	b.numReqPerCycle = n
	return b
}

// WithLowModule sets the port that can provide the address translation in case
// of tlb miss.
func (b LatTLBBuilder) WithLowModule(lowModule akita.Port) LatTLBBuilder {
	b.lowModule = lowModule
	return b
}

// WithNumMSHREntry sets the number of mshr entry
func (b LatTLBBuilder) WithNumMSHREntry(num int) LatTLBBuilder {
	b.numMSHREntry = num
	return b
}

// WithLatency sets the number of mshr entry
func (b LatTLBBuilder) WithLatency(latency int) LatTLBBuilder {
	b.latency = latency
	return b
}

// WithIndexingMask sets the bits to use for indexing
func (b LatTLBBuilder) WithIndexingMask(mask uint64) LatTLBBuilder {
	b.mask = mask
	return b
}

// WithIndexingMask sets the bits to use for indexing
func (b LatTLBBuilder) WithLog2PageSize(log2PageSize uint64) LatTLBBuilder {
	b.log2PageSize = log2PageSize
	return b
}

func (b LatTLBBuilder) UseCoalescingTLBPort() LatTLBBuilder {
	b.useCoalescingTLBPort = true
	return b
}

// Build creates a new LatTLB
func (b LatTLBBuilder) Build(name string) L2TLB {
	tlb := &LatTLB{}
	tlb.TickingComponent =
		akita.NewTickingComponent(name, b.engine, b.freq, tlb)

	tlb.numSets = b.numSets

	tlb.log2NumSets = uint64(math.Log2(float64(tlb.numSets)))
	tlb.setMask = uint64(tlb.numSets - 1)

	tlb.numWays = b.numWays
	tlb.numReqPerCycle = b.numReqPerCycle
	tlb.pageSize = b.pageSize
	tlb.latency = b.latency
	tlb.LowModule = b.lowModule
	tlb.indexingMask = b.mask
	if b.log2PageSize == 0 {
		panic("need to set page size in tlb!")
	}
	tlb.log2PageSize = b.log2PageSize
	if b.useCoalescingTLBPort {
		// tlb.TopPort = NewCoalescingPort(tlb, 16*b.numReqPerCycle,
		// 	name+".TopPort")
		tlb.TopPort = NewCoalescingPort(tlb, 512,
			name+".TopPort")
	} else {
		// tlb.TopPort = akita.NewLimitNumMsgPort(tlb, 16*b.numReqPerCycle,
		// name+".TopPort")
		tlb.TopPort = akita.NewLimitNumMsgPort(tlb, 512, name+".TopPort")
	}
	tlb.BottomPort = akita.NewLimitNumMsgPort(tlb, b.numReqPerCycle,
		name+".BottomPort")
	tlb.ControlPort = akita.NewLimitNumMsgPort(tlb, 1,
		name+".ControlPort")
	tlb.mshr = newMSHR(b.numMSHREntry)
	tlb.lookupBuffer = util.NewBuffer(2 * tlb.numReqPerCycle)
	pipelineBuilder := pipelining.MakeBuilder().WithPipelineWidth(tlb.numReqPerCycle).WithNumStage(tlb.latency).WithCyclePerStage(1).WithPostPipelineBuffer(tlb.lookupBuffer)
	tlb.pipeline = pipelineBuilder.Build(tlb.Name() + "_pipeline")
	tlb.reset()

	return tlb
}
