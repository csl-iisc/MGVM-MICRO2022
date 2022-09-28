package switching

import (
	"fmt"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/util"
	"gitlab.com/akita/util/pipelining"
)

// ChipSwitchConnector can connect switches and ports together.
type ChipSwitchConnector struct {
	engine              akita.Engine
	bufferSizeInNumFlit int
	switchLatency       int
	numReqPerCycle      int
}

// NewChipSwitchConnector creates a new switch connector.
func NewChipSwitchConnector(engine akita.Engine) *ChipSwitchConnector {
	c := &ChipSwitchConnector{
		engine:              engine,
		bufferSizeInNumFlit: 1,
		switchLatency:       1,
	}
	return c
}

// SetBufferSize sets the default buffer size of the port at each switching for
// incoming messages.
func (c *ChipSwitchConnector) SetBufferSize(numFlit int) {
	c.bufferSizeInNumFlit = numFlit
}

// SetSwitchLatency sets the number of cycles required between when the message
// arrives at a switch and when the message is forwarded.
func (c *ChipSwitchConnector) SetSwitchLatency(numCycles int) {
	fmt.Println("setting switch connector latency:", numCycles)
	c.switchLatency = numCycles
}

func (c *ChipSwitchConnector) SetNumReqPerCycle(numReqPerCycle int) {
	c.numReqPerCycle = numReqPerCycle
}

// ConnectSwitches connect two switches together.
// func (c *ChipSwitchConnector) ConnectSwitches(
// 	a, b *Switch,
// 	freq akita.Freq,
// ) (portOnA, portOnB akita.Port) {
// 	portA := akita.NewLimitNumMsgPort(a, c.bufferSizeInNumFlit,
// 		fmt.Sprintf("%s.Port%d", a.Name(), len(a.ports)))
// 	portB := akita.NewLimitNumMsgPort(b, c.bufferSizeInNumFlit,
// 		fmt.Sprintf("%s.Port%d", b.Name(), len(b.ports)))
// 	conn := akita.NewDirectConnection(
// 		fmt.Sprintf("%s-%s", portA.Name(), portB.Name()),
// 		c.engine, freq)
// 	conn.PlugIn(portA, 2*c.numReqPerCycle)
// 	conn.PlugIn(portB, 2*c.numReqPerCycle)

// 	a.addPort(c.createPortComplex(portA, portB))
// 	b.addPort(c.createPortComplex(portB, portA))

// 	return portA, portB
// }

func (c *ChipSwitchConnector) createPortComplex(
	local, remote akita.Port,
) portComplex {

	sendOutBuf := util.NewBuffer(2 * c.numReqPerCycle)
	forwardBuf := util.NewBuffer(2 * c.numReqPerCycle)
	routeBuf := util.NewBuffer(2 * c.numReqPerCycle)
	pipeline := pipelining.NewPipeline(
		local.Name()+"pipeline", c.switchLatency, 1, routeBuf)

	pc := portComplex{
		localPort:     local,
		remotePort:    remote,
		pipeline:      pipeline,
		routeBuffer:   routeBuf,
		forwardBuffer: forwardBuf,
		sendOutBuffer: sendOutBuf,
	}

	return pc
}

// ConnectEndPointToSwitch connects an EndPoint to a Switch.
func (c *ChipSwitchConnector) ConnectEndPointToSwitch(
	ep *EndPoint,
	sw *ChipletSwitch,
	freq akita.Freq,
) (switchPort akita.Port) {
	port := akita.NewLimitNumMsgPort(sw, c.bufferSizeInNumFlit,
		fmt.Sprintf("%s.Port%d", sw.Name(), len(sw.ports)))
	conn := akita.NewDirectConnection(
		fmt.Sprintf("%s-%s", ep.NetworkPort.Name(), port.Name()),
		c.engine, freq)
	conn.PlugIn(port, 2*c.numReqPerCycle)
	conn.PlugIn(ep.NetworkPort, 2*c.numReqPerCycle)

	sw.addPort(c.createPortComplex(port, ep.NetworkPort))

	ep.DefaultSwitchDst = port

	return port
}

// ConnectPortToSwitch connects a port to a Switch.
func (c *ChipSwitchConnector) ConnectPortToSwitch(
	remotePort akita.Port,
	sw *MultiFlitChipletSwitch,
	freq akita.Freq,
) {
	// endPoint := switching.MakeChipletSwitchEndPointBuilder().
	// 	WithEngine(c.engine).
	// 	WithFreq(c.freq).
	// 	WithFlitByteSize(c.flitByteSize).
	// 	WithEncodingOverhead(c.encodingOverhead).
	// 	WithNetworkPortBufferSize(2 * c.numReqPerCycle).
	// 	WithDevicePorts(devicePorts[i : i+1]).
	// 	WithNumReqPerCycle(c.numReqPerCycle).
	// 	Build(fmt.Sprintf("%s.EndPoint%d.%s", c.networkName, len(c.endPoints), dst.Name()))

	// localPort := akita.NewLimitNumMsgPort(sw, c.bufferSizeInNumFlit,
	// fmt.Sprintf("%s.Port%d", sw.Name(), len(sw.ports)))
	// conn := akita.NewDirectConnection(
	// 	fmt.Sprintf("%s-%s", localPort.Name(), remotePort.Name()),
	// 	c.engine, freq)
	// sw.PlugIn(localPort, 2*c.numReqPerCycle)
	sw.addPort(remotePort) //, 2*c.numReqPerCycle)
	//local and then remote
	// sw.addPort(c.createPortComplexForMultiFlits(localPort, remotePort))
	// return p2
}
