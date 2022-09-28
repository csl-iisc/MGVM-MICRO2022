package switching

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/noc"
	"gitlab.com/akita/noc/networking/internal/arbitration"
	"gitlab.com/akita/noc/networking/internal/routing"
	"gitlab.com/akita/util"
	"gitlab.com/akita/util/pipelining"
)

type flitPipelineItem struct {
	taskID string
	flit   *noc.Flit
}

func (f flitPipelineItem) TaskID() string {
	return f.taskID
}

// A portComplex is the infrastructure related to a port.
type portComplex struct {
	// localPort is the port that is equipped on the switch.
	localPort akita.Port

	// remotePort is the port that is connected to the localPort.
	remotePort akita.Port

	// Data arrived at the local port needs to be processed in a pipeline. There
	// is a processing pipeline for each local port.
	pipeline pipelining.Pipeline

	// The flits here are buffered after the pipeline and are waiting to be
	// assigned with an output buffer.
	routeBuffer util.Buffer

	// The flits here are buffered to wait to be forwarded to the output buffer.
	forwardBuffer util.Buffer

	// The flits here are waiting to be sent to the next hop.
	sendOutBuffer util.Buffer
}

// Switch is an Akita component that can forward request to destination.
type Switch struct {
	*akita.TickingComponent

	ports                []akita.Port
	portToComplexMapping map[akita.Port]portComplex
	routingTable         routing.Table
	arbiter              arbitration.Arbiter
	numReqPerCycle       int
}

// addPort adds a new port on the switch.
func (s *Switch) addPort(complex portComplex) {
	s.ports = append(s.ports, complex.localPort)
	s.portToComplexMapping[complex.localPort] = complex
	s.arbiter.AddBuffer(complex.forwardBuffer)
}

// GetRoutingTable returns the routine table used by the switch.
func (s *Switch) GetRoutingTable() routing.Table {
	return s.routingTable
}

// Tick update the Switch's state.
func (s *Switch) Tick(now akita.VTimeInSec) bool {
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

func (s *Switch) startProcessing(now akita.VTimeInSec) (madeProgress bool) {
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

func (s *Switch) movePipeline(now akita.VTimeInSec) (madeProgress bool) {
	for _, port := range s.ports {
		complex := s.portToComplexMapping[port]
		madeProgress = complex.pipeline.Tick(now) || madeProgress
	}

	return madeProgress
}

// assigns a flit the outputbuffer to which it is to be routed
func (s *Switch) route(now akita.VTimeInSec) (madeProgress bool) {
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

// The order in which switches tick might also be important
func (s *Switch) forward(now akita.VTimeInSec) (madeProgress bool) {
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
		// fmt.Println(flit.Src.Name())
		flit.OutputBuf.Push(flit)
		buf.Pop()
		madeProgress = true
	}

	return madeProgress
}

func (s *Switch) sendOut(now akita.VTimeInSec) (madeProgress bool) {
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

func (s *Switch) assignFlitOutputBuf(f *noc.Flit) {
	outPort := s.routingTable.FindPort(f.Msg.Meta().Dst)
	complex := s.portToComplexMapping[outPort]
	f.OutputBuf = complex.sendOutBuffer
}

func (s *Switch) setFlitNextHopDst(f *noc.Flit) {
	f.Src = f.Dst
	f.Dst = s.portToComplexMapping[f.Src].remotePort
}

// SwitchBuilder can build switches
type SwitchBuilder struct {
	engine         akita.Engine
	freq           akita.Freq
	routingTable   routing.Table
	arbiter        arbitration.Arbiter
	numReqPerCycle int
}

// WithEngine sets the engine that the switch to build uses.
func (b SwitchBuilder) WithEngine(engine akita.Engine) SwitchBuilder {
	b.engine = engine
	return b
}

// WithFreq sets the frequency that the switch to build works at.
func (b SwitchBuilder) WithFreq(freq akita.Freq) SwitchBuilder {
	b.freq = freq
	return b
}

// WithArbiter sets the arbiter to be used by the swtich to build.
func (b SwitchBuilder) WithArbiter(arbiter arbitration.Arbiter) SwitchBuilder {
	b.arbiter = arbiter
	return b
}

// WithRoutingTable sets the routing table to be used by the switch to build.
func (b SwitchBuilder) WithRoutingTable(rt routing.Table) SwitchBuilder {
	b.routingTable = rt
	return b
}

func (b SwitchBuilder) WithNumReqPerCycle(numReqPerCycle int) SwitchBuilder {
	b.numReqPerCycle = numReqPerCycle
	return b
}

// Build creates a new switch
func (b SwitchBuilder) Build(name string) *Switch {
	b.engineMustBeGiven()
	b.freqMustNotBeZero()
	b.routingTableMustBeGiven()
	b.arbiterMustBeGiven()

	s := &Switch{}
	s.TickingComponent = akita.NewTickingComponent(name, b.engine, b.freq, s)
	s.routingTable = b.routingTable
	s.arbiter = b.arbiter
	s.portToComplexMapping = make(map[akita.Port]portComplex)
	s.numReqPerCycle = b.numReqPerCycle
	return s
}

func (b SwitchBuilder) engineMustBeGiven() {
	if b.engine == nil {
		panic("engine of switch is not given")
	}
}

func (b SwitchBuilder) freqMustNotBeZero() {
	if b.freq == 0 {
		panic("switch frequency cannot be 0")
	}
}

func (b SwitchBuilder) routingTableMustBeGiven() {
	if b.routingTable == nil {
		panic("switch requires a routing table to operate")
	}
}

func (b SwitchBuilder) arbiterMustBeGiven() {
	if b.arbiter == nil {
		panic("switch requires an arbiter to operate")
	}
}
