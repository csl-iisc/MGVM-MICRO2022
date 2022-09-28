package switching

import (
	"fmt"
	"math"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/noc"
	"gitlab.com/akita/noc/networking/internal/arbitration"
	"gitlab.com/akita/noc/networking/internal/routing"
	"gitlab.com/akita/util"
	"gitlab.com/akita/util/pipelining"
)

type multiFlitPortComplex struct {
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

	assemblingMsg  akita.Msg
	numFlitReqired int
	numFlitArrived int
	// flitsToSend    []*noc.Flit
	msgInBuf util.Buffer
	flitBuf  util.Buffer
	// []akita.Msg
}

// MultiFlitChipletSwitch is an Akita component that can forward request to destination.
type MultiFlitChipletSwitch struct {
	*akita.TickingComponent

	ports                []akita.Port
	portToComplexMapping map[akita.Port]*multiFlitPortComplex
	routingTable         routing.Table
	arbiter              *arbitration.FlitAwareXBarArbiter
	numReqPerCycle       int
	// numChiplets               int
	// outgoingReqsPerChiplet    []int
	// maxOutgoingReqsPerChiplet int

	flitByteSize     int
	encodingOverhead float64

	switchLatency       int
	bufferSizeInNumFlit int
}

func (s *MultiFlitChipletSwitch) createPortComplexForMultiFlits(remotePort akita.Port,
) multiFlitPortComplex {
	switchConn := NewMultiFlitSwitchConnection(remotePort.Name()+".SwitchConnection", s.Engine, s.Freq)
	localPort := akita.NewLimitNumMsgPort(s, s.bufferSizeInNumFlit,
		fmt.Sprintf("%s.Port%d", s.Name(), len(s.ports)))
	switchConn.PlugIn(localPort, 2*s.numReqPerCycle)
	switchConn.PlugIn(remotePort, 2*s.numReqPerCycle)
	msgInBuf := util.NewBuffer(s.numReqPerCycle)
	flitBuf := util.NewBuffer(2 * s.numReqPerCycle)
	sendOutBuf := util.NewBuffer(2 * s.numReqPerCycle)
	forwardBuf := util.NewBuffer(2 * s.numReqPerCycle)
	routeBuf := util.NewBuffer(2 * s.numReqPerCycle)
	pipeline := pipelining.NewPipeline(
		remotePort.Name()+"pipeline", s.switchLatency, 1, routeBuf)

	pc := multiFlitPortComplex{
		localPort:     localPort,
		remotePort:    remotePort,
		msgInBuf:      msgInBuf,
		flitBuf:       flitBuf,
		pipeline:      pipeline,
		routeBuffer:   routeBuf,
		forwardBuffer: forwardBuf,
		sendOutBuffer: sendOutBuf,
	}

	return pc
}

// addPort adds a new port on the MultiFlitChipletSwitch.
func (s *MultiFlitChipletSwitch) addPort(remotePort akita.Port) {
	complex := s.createPortComplexForMultiFlits(remotePort)
	s.ports = append(s.ports, complex.localPort)
	s.portToComplexMapping[complex.localPort] = &complex
	s.arbiter.AddBuffer(complex.forwardBuffer)
	s.routingTable.DefineRoute(remotePort, complex.localPort)
}

// GetRoutingTable returns the routine table used by the MultiFlitChipletSwitch.
func (s *MultiFlitChipletSwitch) GetRoutingTable() routing.Table {
	return s.routingTable
}

// Tick update the MultiFlitChipletSwitch's state.
func (s *MultiFlitChipletSwitch) Tick(now akita.VTimeInSec) bool {
	// for i := 0; i < s.numChiplets; i++ {
	// 	s.outgoingReqsPerChiplet[i] = 0
	// }
	s.arbiter.Reset()
	madeProgress := false
	for i := 0; i < s.numReqPerCycle; i++ {
		// fmt.Println("boo")
		madeProgress = s.tryDeliver(now) || madeProgress
		madeProgress = s.assemble(now) || madeProgress
		madeProgress = s.forward(now) || madeProgress
		madeProgress = s.route(now) || madeProgress
		madeProgress = s.movePipeline(now) || madeProgress
		madeProgress = s.startProcessing(now) || madeProgress
		madeProgress = s.prepareFlits(now) || madeProgress
		madeProgress = s.recv(now) || madeProgress
	}
	// return true
	return madeProgress
}

func (s *MultiFlitChipletSwitch) assemble(now akita.VTimeInSec) bool {
	madeProgress := false
	for _, port := range s.ports {
		complex := s.portToComplexMapping[port]
		newFlits := make([]*noc.Flit, 0)
		// numFlits := complex.sendOutBuffer.Size()
		for {
			// for _, f := range complex.flitsToAssemble {
			item := complex.sendOutBuffer.Pop()
			if item == nil {
				break
			}
			f := item.(*noc.Flit)
			if complex.assemblingMsg == nil {
				complex.assemblingMsg = f.Msg
				complex.numFlitArrived = 1
				complex.numFlitReqired = f.NumFlitInMsg
				// fmt.Println("here")
				madeProgress = true
			} else if complex.assemblingMsg != nil && f.Msg == complex.assemblingMsg {
				complex.numFlitArrived++
				madeProgress = true
			} else {
				newFlits = append(newFlits, f)
			}
		}
		for i := 0; i < len(newFlits); i++ {
			complex.sendOutBuffer.Push(newFlits[i])
		}
		// complex.flitsToAssemble = newFlits
	}
	return madeProgress
}

func (s *MultiFlitChipletSwitch) tryDeliver(now akita.VTimeInSec) bool {
	madeProgress := false
	for _, port := range s.ports {
		complex := s.portToComplexMapping[port]
		if complex.assemblingMsg == nil {
			continue
		}
		if complex.numFlitArrived < complex.numFlitReqired {
			continue
		}
		complex.assemblingMsg.Meta().SendTime = now
		err := complex.localPort.Send(complex.assemblingMsg)
		// fmt.Println("yooooooooooooooooooooo")
		// complex.assemblingMsg.Meta().Dst.Recv(complex.assemblingMsg)
		if err == nil {
			// log.Printf("%.12f, EP %s, msg %s assembled and deliverd\n",
			// 	now, ep.Name(), ep.assemblingMsg.Meta().ID)
			complex.assemblingMsg = nil
			complex.numFlitReqired = 0
			complex.numFlitArrived = 0
			madeProgress = true
		} //ICN.Switch0.Port6
	}
	return madeProgress
}
func (s *MultiFlitChipletSwitch) prepareFlits(now akita.VTimeInSec) bool {
	madeProgress := false
	for _, port := range s.ports {
		complex := s.portToComplexMapping[port]

		if complex.flitBuf.Size() > 0 {
			continue
		}

		if complex.msgInBuf.Size() == 0 {
			continue
		}

		msg := complex.msgInBuf.Pop().(akita.Msg)
		for _, f := range s.msgToFlits(msg) {
			complex.flitBuf.Push(f)
		}
		madeProgress = true
	}
	return madeProgress
}

func (s *MultiFlitChipletSwitch) msgToFlits(msg akita.Msg) []*noc.Flit {
	numFlit := 1
	if msg.Meta().TrafficBytes > 0 {
		trafficByte := msg.Meta().TrafficBytes
		trafficByte += int(math.Ceil(
			float64(trafficByte) * s.encodingOverhead))
		numFlit = (trafficByte-1)/s.flitByteSize + 1
	}

	flits := make([]*noc.Flit, numFlit)
	for i := 0; i < numFlit; i++ {
		flits[i] = noc.FlitBuilder{}.
			WithSrc(msg.Meta().Src).
			// WithDst(s.DefaultSwitchDst).
			WithSeqID(i).
			WithNumFlitInMsg(numFlit).
			WithMsg(msg).
			Build()
	}
	return flits
}

func (s *MultiFlitChipletSwitch) recv(now akita.VTimeInSec) bool {
	madeProgress := false
	for _, port := range s.ports {
		item := port.Peek()
		if item == nil {
			continue
		}
		complex := s.portToComplexMapping[port]
		if !complex.msgInBuf.CanPush() {
			continue
		}
		complex.msgInBuf.Push(item)
		port.Retrieve(now)
		madeProgress = true
	}
	return madeProgress
}

func (s *MultiFlitChipletSwitch) startProcessing(now akita.VTimeInSec) (madeProgress bool) {
	madeProgress = false
	for _, port := range s.ports {
		complex := s.portToComplexMapping[port]
		if !complex.pipeline.CanAccept() {
			continue
		}
		if complex.flitBuf.Size() == 0 {
			continue
		}
		item := complex.flitBuf.Pop()
		// if item == nil {
		// 	continue
		// }
		pipelineItem := flitPipelineItem{
			taskID: akita.GetIDGenerator().Generate(),
			flit:   item.(*noc.Flit),
		}
		complex.pipeline.Accept(now, pipelineItem)
		madeProgress = true
	}

	return madeProgress
}

func (s *MultiFlitChipletSwitch) movePipeline(now akita.VTimeInSec) (madeProgress bool) {
	for _, port := range s.ports {
		complex := s.portToComplexMapping[port]
		madeProgress = complex.pipeline.Tick(now) || madeProgress
	}

	return madeProgress
}

// assigns a flit the outputbuffer to which it is to be routed
func (s *MultiFlitChipletSwitch) route(now akita.VTimeInSec) (madeProgress bool) {
	for _, port := range s.ports {
		multiFlitPortComplex := s.portToComplexMapping[port]
		routeBuf := multiFlitPortComplex.routeBuffer
		forwardBuf := multiFlitPortComplex.forwardBuffer

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

// func getChipletFromFlitSrc(flit *noc.Flit) (chiplet int) {
// 	src := strings.Split(flit.Src.Name(), ".")[1]
// 	chiplet, err := strconv.Atoi(strings.Split(src, "_")[1])
// 	if err != nil {
// 		panic("something went wrong")
// 	}
// 	return
// }

// The order in which MultiFlitChipletSwitches tick might also be important
func (s *MultiFlitChipletSwitch) forward(now akita.VTimeInSec) (madeProgress bool) {
	inputBuffers := s.arbiter.Arbitrate(now)

	for _, buf := range inputBuffers {
		item := buf.Peek()
		if item == nil {
			continue
		}
		flit := item.(*noc.Flit)

		if !flit.OutputBuf.CanPush() {
			panic("something went wrong")
			// continue
		}

		// chiplet := getChipletFromFlitSrc(flit)

		// if s.outgoingReqsPerChiplet[chiplet] >= s.maxOutgoingReqsPerChiplet {
		// 	continue
		// }
		flit.OutputBuf.Push(flit)
		buf.Pop()
		// s.outgoingReqsPerChiplet[chiplet]++
		madeProgress = true
	}

	return madeProgress
}

// func (s *MultiFlitChipletSwitch) sendOut(now akita.VTimeInSec) (madeProgress bool) {
// 	for _, port := range s.ports {
// 		complex := s.portToComplexMapping[port]
// 		sendOutBuf := complex.sendOutBuffer
// 		item := sendOutBuf.Peek()
// 		if item == nil {
// 			continue
// 		}
// 		flit := item.(*noc.Flit)
// 		flit.Meta().Src = complex.localPort
// 		flit.Meta().Dst = complex.remotePort
// 		flit.Meta().SendTime = now
// 		err := complex.localPort.Send(flit)
// 		if err == nil {
// 			sendOutBuf.Pop()
// 			madeProgress = true
// 		}
// 	}
// 	return madeProgress
// }

func (s *MultiFlitChipletSwitch) assignFlitOutputBuf(f *noc.Flit) {
	outPort := s.routingTable.FindPort(f.Msg.Meta().Dst)
	// fmt.Println(f.Msg.Meta().Dst.Name(), outPort.Name())
	complex := s.portToComplexMapping[outPort]
	f.OutputBuf = complex.sendOutBuffer
}

func (s *MultiFlitChipletSwitch) setFlitNextHopDst(f *noc.Flit) {
	f.Src = f.Dst
	f.Dst = s.portToComplexMapping[f.Src].remotePort
}

// MultiFlitChipletSwitchBuilder can build MultiFlitChipletSwitches
type MultiFlitChipletSwitchBuilder struct {
	engine         akita.Engine
	freq           akita.Freq
	routingTable   routing.Table
	arbiter        arbitration.Arbiter
	numReqPerCycle int
	// numChiplets               int
	// maxOutgoingReqsPerChiplet int
}

// WithEngine sets the engine that the MultiFlitChipletSwitch to build uses.
func (b MultiFlitChipletSwitchBuilder) WithEngine(engine akita.Engine) MultiFlitChipletSwitchBuilder {
	b.engine = engine
	return b
}

// WithFreq sets the frequency that the MultiFlitChipletSwitch to build works at.
func (b MultiFlitChipletSwitchBuilder) WithFreq(freq akita.Freq) MultiFlitChipletSwitchBuilder {
	b.freq = freq
	return b
}

// WithArbiter sets the arbiter to be used by the swtich to build.
func (b MultiFlitChipletSwitchBuilder) WithArbiter(arbiter arbitration.Arbiter) MultiFlitChipletSwitchBuilder {
	b.arbiter = arbiter
	return b
}

// WithRoutingTable sets the routing table to be used by the MultiFlitChipletSwitch to build.
func (b MultiFlitChipletSwitchBuilder) WithRoutingTable(rt routing.Table) MultiFlitChipletSwitchBuilder {
	b.routingTable = rt
	return b
}

func (b MultiFlitChipletSwitchBuilder) WithNumReqPerCycle(numReqPerCycle int) MultiFlitChipletSwitchBuilder {
	b.numReqPerCycle = numReqPerCycle
	return b
}

func (b MultiFlitChipletSwitchBuilder) WithNumChiplets(numChiplets int) MultiFlitChipletSwitchBuilder {
	// b.numChiplets = numChiplets
	return b
}

func (b MultiFlitChipletSwitchBuilder) WithMaxOutgoingReqsPerChiplet(maxOutgoingReqsPerChiplet int) MultiFlitChipletSwitchBuilder {
	// b.maxOutgoingReqsPerChiplet = maxOutgoingReqsPerChiplet
	return b
}

// Build creates a new MultiFlitChipletSwitch
func (b MultiFlitChipletSwitchBuilder) Build(name string) *MultiFlitChipletSwitch {
	b.engineMustBeGiven()
	b.freqMustNotBeZero()
	b.routingTableMustBeGiven()
	b.arbiterMustBeGiven()

	s := &MultiFlitChipletSwitch{}
	s.TickingComponent = akita.NewTickingComponent(name, b.engine, b.freq, s)
	s.routingTable = b.routingTable
	s.arbiter = b.arbiter.(*arbitration.FlitAwareXBarArbiter)
	s.portToComplexMapping = make(map[akita.Port]*multiFlitPortComplex)
	s.numReqPerCycle = b.numReqPerCycle
	// s.numChiplets = b.numChiplets
	// s.maxOutgoingReqsPerChiplet = b.maxOutgoingReqsPerChiplet
	// if s.numChiplets == 0 || s.maxOutgoingReqsPerChiplet == 0 {
	// panic("must be positive")
	// }
	// s.outgoingReqsPerChiplet = make([]int, s.numChiplets)
	s.flitByteSize = 64
	s.encodingOverhead = 0

	s.switchLatency = 360 //720 //180 //360 // changed this here
	fmt.Println("switch latency:", s.switchLatency)
	s.bufferSizeInNumFlit = 64
	return s
}

func (b MultiFlitChipletSwitchBuilder) engineMustBeGiven() {
	if b.engine == nil {
		panic("engine of MultiFlitChipletSwitch is not given")
	}
}

func (b MultiFlitChipletSwitchBuilder) freqMustNotBeZero() {
	if b.freq == 0 {
		panic("MultiFlitChipletSwitch frequency cannot be 0")
	}
}

func (b MultiFlitChipletSwitchBuilder) routingTableMustBeGiven() {
	if b.routingTable == nil {
		panic("MultiFlitChipletSwitch requires a routing table to operate")
	}
}

func (b MultiFlitChipletSwitchBuilder) arbiterMustBeGiven() {
	if b.arbiter == nil {
		panic("MultiFlitChipletSwitch requires an arbiter to operate")
	}
}
