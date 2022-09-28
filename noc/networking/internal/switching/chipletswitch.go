package switching

import (
	"strconv"
	"strings"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/noc"
	"gitlab.com/akita/noc/networking/internal/arbitration"
	"gitlab.com/akita/noc/networking/internal/routing"
)

// ChipletSwitch is an Akita component that can forward request to destination.
type ChipletSwitch struct {
	*akita.TickingComponent

	ports                     []akita.Port
	portToComplexMapping      map[akita.Port]portComplex
	routingTable              routing.Table
	arbiter                   arbitration.Arbiter
	numReqPerCycle            int
	numChiplets               int
	outgoingReqsPerChiplet    []int
	maxOutgoingReqsPerChiplet int
}

// addPort adds a new port on the ChipletSwitch.
func (s *ChipletSwitch) addPort(complex portComplex) {
	s.ports = append(s.ports, complex.localPort)
	s.portToComplexMapping[complex.localPort] = complex
	s.arbiter.AddBuffer(complex.forwardBuffer)
}

// GetRoutingTable returns the routine table used by the ChipletSwitch.
func (s *ChipletSwitch) GetRoutingTable() routing.Table {
	return s.routingTable
}

// Tick update the ChipletSwitch's state.
func (s *ChipletSwitch) Tick(now akita.VTimeInSec) bool {
	for i := 0; i < s.numChiplets; i++ {
		s.outgoingReqsPerChiplet[i] = 0
	}
	madeProgress := false

	for i := 0; i < s.numReqPerCycle; i++ {
		madeProgress = s.sendOut(now) || madeProgress
		madeProgress = s.forward(now) || madeProgress
		// pop from the pipeline
		madeProgress = s.route(now) || madeProgress
		// }
		// make the pipeline make forward progress
		madeProgress = s.movePipeline(now) || madeProgress
		// for i := 0; i < s.numReqPerCycle; i++ {
		// insert into pipeline
		madeProgress = s.startProcessing(now) || madeProgress
	}

	return madeProgress
}

func (s *ChipletSwitch) startProcessing(now akita.VTimeInSec) (madeProgress bool) {
	for _, port := range s.ports {
		item := port.Peek()
		if item == nil {
			continue
		}

		complex := s.portToComplexMapping[port]
		if !complex.pipeline.CanAccept() {
			continue
		}

		pipelineItem := flitPipelineItem{
			taskID: akita.GetIDGenerator().Generate(),
			flit:   item.(*noc.Flit),
		}
		complex.pipeline.Accept(now, pipelineItem)
		port.Retrieve(now)
		madeProgress = true
	}

	return madeProgress
}

func (s *ChipletSwitch) movePipeline(now akita.VTimeInSec) (madeProgress bool) {
	for _, port := range s.ports {
		complex := s.portToComplexMapping[port]
		madeProgress = complex.pipeline.Tick(now) || madeProgress
	}

	return madeProgress
}

// assigns a flit the outputbuffer to which it is to be routed
func (s *ChipletSwitch) route(now akita.VTimeInSec) (madeProgress bool) {
	for _, port := range s.ports {
		portComplex := s.portToComplexMapping[port]
		routeBuf := portComplex.routeBuffer
		forwardBuf := portComplex.forwardBuffer

		item := routeBuf.Peek()
		if item == nil {
			continue
		}

		if !forwardBuf.CanPush() {
			continue
		}

		pipelineItem := item.(flitPipelineItem)
		flit := pipelineItem.flit
		s.assignFlitOutputBuf(flit)
		routeBuf.Pop()
		forwardBuf.Push(flit)
		madeProgress = true
	}

	return madeProgress
}

func getChipletFromFlit(flit *noc.Flit) (chiplet int) {
	src := strings.Split(flit.Src.Name(), ".")[3]
	chiplet, err := strconv.Atoi(strings.Split(src, "_")[1])
	if err != nil {
		panic("something went wrong")
	}
	return
}

// The order in which ChipletSwitches tick might also be important
func (s *ChipletSwitch) forward(now akita.VTimeInSec) (madeProgress bool) {
	inputBuffers := s.arbiter.Arbitrate(now)

	for _, buf := range inputBuffers {
		item := buf.Peek()
		if item == nil {
			continue
		}
		flit := item.(*noc.Flit)

		if !flit.OutputBuf.CanPush() {
			continue
		}
		chiplet := getChipletFromFlit(flit)
		// fmt.Println(chiplet)
		if s.outgoingReqsPerChiplet[chiplet] >= s.maxOutgoingReqsPerChiplet {
			continue
		}
		flit.OutputBuf.Push(flit)
		buf.Pop()
		s.outgoingReqsPerChiplet[chiplet]++
		madeProgress = true
	}

	return madeProgress
}

func (s *ChipletSwitch) sendOut(now akita.VTimeInSec) (madeProgress bool) {
	for _, port := range s.ports {
		complex := s.portToComplexMapping[port]
		sendOutBuf := complex.sendOutBuffer

		item := sendOutBuf.Peek()
		if item == nil {
			continue
		}

		flit := item.(*noc.Flit)
		flit.Meta().Src = complex.localPort
		flit.Meta().Dst = complex.remotePort
		flit.Meta().SendTime = now
		err := complex.localPort.Send(flit)
		if err == nil {
			sendOutBuf.Pop()
			madeProgress = true
		}
	}

	return madeProgress
}

func (s *ChipletSwitch) assignFlitOutputBuf(f *noc.Flit) {
	outPort := s.routingTable.FindPort(f.Msg.Meta().Dst)
	complex := s.portToComplexMapping[outPort]
	f.OutputBuf = complex.sendOutBuffer
}

func (s *ChipletSwitch) setFlitNextHopDst(f *noc.Flit) {
	f.Src = f.Dst
	f.Dst = s.portToComplexMapping[f.Src].remotePort
}

// ChipletSwitchBuilder can build ChipletSwitches
type ChipletSwitchBuilder struct {
	engine                    akita.Engine
	freq                      akita.Freq
	routingTable              routing.Table
	arbiter                   arbitration.Arbiter
	numReqPerCycle            int
	numChiplets               int
	maxOutgoingReqsPerChiplet int
}

// WithEngine sets the engine that the ChipletSwitch to build uses.
func (b ChipletSwitchBuilder) WithEngine(engine akita.Engine) ChipletSwitchBuilder {
	b.engine = engine
	return b
}

// WithFreq sets the frequency that the ChipletSwitch to build works at.
func (b ChipletSwitchBuilder) WithFreq(freq akita.Freq) ChipletSwitchBuilder {
	b.freq = freq
	return b
}

// WithArbiter sets the arbiter to be used by the swtich to build.
func (b ChipletSwitchBuilder) WithArbiter(arbiter arbitration.Arbiter) ChipletSwitchBuilder {
	b.arbiter = arbiter
	return b
}

// WithRoutingTable sets the routing table to be used by the ChipletSwitch to build.
func (b ChipletSwitchBuilder) WithRoutingTable(rt routing.Table) ChipletSwitchBuilder {
	b.routingTable = rt
	return b
}

func (b ChipletSwitchBuilder) WithNumReqPerCycle(numReqPerCycle int) ChipletSwitchBuilder {
	b.numReqPerCycle = numReqPerCycle
	return b
}

func (b ChipletSwitchBuilder) WithNumChiplets(numChiplets int) ChipletSwitchBuilder {
	b.numChiplets = numChiplets
	return b
}

func (b ChipletSwitchBuilder) WithMaxOutgoingReqsPerChiplet(maxOutgoingReqsPerChiplet int) ChipletSwitchBuilder {
	b.maxOutgoingReqsPerChiplet = maxOutgoingReqsPerChiplet
	return b
}

// Build creates a new ChipletSwitch
func (b ChipletSwitchBuilder) Build(name string) *ChipletSwitch {
	b.engineMustBeGiven()
	b.freqMustNotBeZero()
	b.routingTableMustBeGiven()
	b.arbiterMustBeGiven()

	s := &ChipletSwitch{}
	s.TickingComponent = akita.NewTickingComponent(name, b.engine, b.freq, s)
	s.routingTable = b.routingTable
	s.arbiter = b.arbiter
	s.portToComplexMapping = make(map[akita.Port]portComplex)
	s.numReqPerCycle = b.numReqPerCycle
	s.numChiplets = b.numChiplets
	s.maxOutgoingReqsPerChiplet = b.maxOutgoingReqsPerChiplet
	if s.numChiplets == 0 || s.maxOutgoingReqsPerChiplet == 0 {
		panic("must be positive")
	}
	s.outgoingReqsPerChiplet = make([]int, s.numChiplets)
	return s
}

func (b ChipletSwitchBuilder) engineMustBeGiven() {
	if b.engine == nil {
		panic("engine of ChipletSwitch is not given")
	}
}

func (b ChipletSwitchBuilder) freqMustNotBeZero() {
	if b.freq == 0 {
		panic("ChipletSwitch frequency cannot be 0")
	}
}

func (b ChipletSwitchBuilder) routingTableMustBeGiven() {
	if b.routingTable == nil {
		panic("ChipletSwitch requires a routing table to operate")
	}
}

func (b ChipletSwitchBuilder) arbiterMustBeGiven() {
	if b.arbiter == nil {
		panic("ChipletSwitch requires an arbiter to operate")
	}
}
