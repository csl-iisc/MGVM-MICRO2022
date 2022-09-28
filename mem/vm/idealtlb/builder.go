package idealtlb

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem/device"
)

// A Builder can build ideal TLBs
type Builder struct {
	engine         akita.Engine
	freq           akita.Freq
	numReqPerCycle int
	pageTable      device.PageTable
}

// MakeBuilder returns a Builder
func MakeBuilder() Builder {
	return Builder{
		freq:           1 * akita.GHz,
		numReqPerCycle: 1024,
	}
}

// WithEngine sets the engine that the TLBs to use
func (b Builder) WithEngine(engine akita.Engine) Builder {
	b.engine = engine
	return b
}

// WithFreq sets the freq the TLBs use
func (b Builder) WithFreq(freq akita.Freq) Builder {
	b.freq = freq
	return b
}

// WithNumReqPerCycle sets the number of requests per cycle can be processed by
// a TLB
func (b Builder) WithNumReqPerCycle(n int) Builder {
	b.numReqPerCycle = n
	return b
}

// WithPageTable sets the page table (for magix lookup).
func (b Builder) WithPageTable(pageTable device.PageTable) Builder {
	b.pageTable = pageTable
	return b
}

// Build creates a new Ideal TLB
func (b Builder) Build(name string) *IdealTLB {
	idealtlb := &IdealTLB{}
	idealtlb.TickingComponent =
		akita.NewTickingComponent(name, b.engine, b.freq, idealtlb)

	idealtlb.pageTable = b.pageTable
	idealtlb.numReqPerCycle = b.numReqPerCycle
	idealtlb.isPaused = false

	idealtlb.TopPort = akita.NewLimitNumMsgPort(idealtlb, b.numReqPerCycle,
		name+".TopPort")
	idealtlb.BottomPort = akita.NewLimitNumMsgPort(idealtlb, b.numReqPerCycle,
		name+".BottomPort")
	idealtlb.ControlPort = akita.NewLimitNumMsgPort(idealtlb, 1,
		name+".ControlPort")

	return idealtlb
}
