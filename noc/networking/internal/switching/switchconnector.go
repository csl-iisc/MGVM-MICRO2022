package switching

import (
	"fmt"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/util"
	"gitlab.com/akita/util/pipelining"
)

// SwitchConnector can connect switches and ports together.
type SwitchConnector struct {
	engine              akita.Engine
	bufferSizeInNumFlit int
	switchLatency       int
	numReqPerCycle      int
}

// NewSwitchConnector creates a new switch connector.
func NewSwitchConnector(engine akita.Engine) *SwitchConnector {
	c := &SwitchConnector{
		engine:              engine,
		bufferSizeInNumFlit: 1,
		switchLatency:       1,
	}
	return c
}

// SetBufferSize sets the default buffer size of the port at each switching for
// incoming messages.
func (c *SwitchConnector) SetBufferSize(numFlit int) {
	c.bufferSizeInNumFlit = numFlit
}

// SetSwitchLatency sets the number of cycles required between when the message
// arrives at a switch and when the message is forwarded.
func (c *SwitchConnector) SetSwitchLatency(numCycles int) {
	c.switchLatency = numCycles
}

func (c *SwitchConnector) SetNumReqPerCycle(numReqPerCycle int) {
	c.numReqPerCycle = numReqPerCycle
}

// ConnectSwitches connect two switches together.
func (c *SwitchConnector) ConnectSwitches(
	a, b *Switch,
	freq akita.Freq,
) (portOnA, portOnB akita.Port) {
	portA := akita.NewLimitNumMsgPort(a, c.bufferSizeInNumFlit,
		fmt.Sprintf("%s.Port%d", a.Name(), len(a.ports)))
	portB := akita.NewLimitNumMsgPort(b, c.bufferSizeInNumFlit,
		fmt.Sprintf("%s.Port%d", b.Name(), len(b.ports)))
	conn := akita.NewDirectConnection(
		fmt.Sprintf("%s-%s", portA.Name(), portB.Name()),
		c.engine, freq)
	conn.PlugIn(portA, 2*c.numReqPerCycle)
	conn.PlugIn(portB, 2*c.numReqPerCycle)

	a.addPort(c.createPortComplex(portA, portB))
	b.addPort(c.createPortComplex(portB, portA))

	return portA, portB
}

func (c *SwitchConnector) createPortComplex(
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
func (c *SwitchConnector) ConnectEndPointToSwitch(
	ep *EndPoint,
	sw *Switch,
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

// // ConnectPortToSwitch connects a port to a Switch.
// func (c *SwitchConnector) ConnectPortToSwitch(
// 	p1 akita.Port,
// 	sw *Switch,
// 	freq akita.Freq,
// ) (switchPort akita.Port) {
// 	p2 := akita.NewLimitNumMsgPort(sw, c.bufferSizeInNumFlit,
// 		fmt.Sprintf("%s.Port%d", sw.Name(), len(sw.ports)))
// 	conn := akita.NewDirectConnection(
// 		fmt.Sprintf("%s-%s", p1.Name(), p2.Name()),
// 		c.engine, freq)
// 	conn.PlugIn(p1, 2*c.numReqPerCycle)
// 	conn.PlugIn(p2, 2*c.numReqPerCycle)
// 	//local and then remote
// 	sw.addPort(c.createPortComplex(p2, p1))
// 	return p2
// }
