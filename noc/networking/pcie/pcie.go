// Package pcie provides a Connector and establishes a PCIe connection.
package pcie

import (
	"fmt"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/noc/networking/internal/arbitration"
	"gitlab.com/akita/noc/networking/internal/networking"
	"gitlab.com/akita/noc/networking/internal/routing"
	"gitlab.com/akita/noc/networking/internal/switching"
)

type switchNode struct {
	parentSwitchNode   *switchNode
	currSwitch         *switching.Switch
	localPortToParent  akita.Port
	remotePortToParent akita.Port
}

// Connector can connect devices into a PCIe network.
type Connector struct {
	networkName      string
	engine           akita.Engine
	freq             akita.Freq
	encodingOverhead float64
	flitByteSize     int
	switchLatency    int
	switchConnector  *switching.SwitchConnector
	network          *networking.NetworkedConnection
	switches         []*switching.Switch
	rootComplexes    []*switching.Switch
	switchNodes      []*switchNode // The same switches, but as trees.
	endPoints        []*switching.EndPoint
	numReqPerCycle   int
}

// NewConnector creates a new connector that can help configure PCIe networks.
func NewConnector() *Connector {
	c := &Connector{}
	c = c.WithVersion3().
		WithX16().
		WithSwitchLatency(140)
	c.numReqPerCycle = 1
	return c
}

// WithEngine sets the event-driven simulation engine that the PCIe connection
// uses.
func (c *Connector) WithEngine(engine akita.Engine) *Connector {
	c.engine = engine
	return c
}

// WithVersion3 defines the PCIe connection to use the PCIe-v3 standard.
func (c *Connector) WithVersion3() *Connector {
	c.freq = 1 * akita.GHz
	c.encodingOverhead = (130 - 128) / 130
	return c
}

// WithVersion4 defines the PCIe connection to use the PCIe-v4 standard.
func (c *Connector) WithVersion4() *Connector {
	c.freq = 2 * akita.GHz
	c.encodingOverhead = (130 - 128) / 130
	return c
}

func (c *Connector) WithSwitchLatency(numCycles int) *Connector {
	c.switchLatency = numCycles
	return c
}

// WithX1 sets PCIe connection to a X1 PCIe connection.
func (c *Connector) WithX1() *Connector {
	c.flitByteSize = 1
	return c
}

// WithX2 sets PCIe connection to a X1 PCIe connection.
func (c *Connector) WithX2() *Connector {
	c.flitByteSize = 2
	return c
}

// WithX4 sets PCIe connection to a X1 PCIe connection.
func (c *Connector) WithX4() *Connector {
	c.flitByteSize = 4
	return c
}

// WithX8 sets PCIe connection to a X1 PCIe connection.
func (c *Connector) WithX8() *Connector {
	c.flitByteSize = 8
	return c
}

// WithX16 sets PCIe connection to a X1 PCIe connection.
func (c *Connector) WithX16() *Connector {
	c.flitByteSize = 16
	return c
}

// WithNetworkName sets the name of the network and the prefix of all the
// component in the network.
func (c *Connector) WithNetworkName(name string) *Connector {
	c.networkName = name
	return c
}

// CreateNetwork creates a network. This function should be called before
// creating root complexes.
func (c *Connector) CreateNetwork() {
	c.network = networking.NewNetworkedConnection()
	c.switchConnector = switching.NewSwitchConnector(c.engine)
	c.switchConnector.SetSwitchLatency(c.switchLatency)
	c.switchConnector.SetNumReqPerCycle(c.numReqPerCycle)
	c.switchConnector.SetBufferSize(2 * c.numReqPerCycle)
}

// CreateRootComplex creates a root complex of the PCIe connection. It requires
// a set of port that connects to the CPU in the system.
func (c *Connector) CreateRootComplex(cpuPorts []akita.Port) (switchID int) {
	cpuEndPoint := switching.MakeEndPointBuilder().
		WithEngine(c.engine).
		WithFreq(c.freq).
		WithDevicePorts(cpuPorts).
		WithFlitByteSize(c.flitByteSize).
		WithEncodingOverhead(c.encodingOverhead).
		Build(c.networkName + ".CPUEndPoint")
	c.network.AddEndPoint(cpuEndPoint)
	c.endPoints = append(c.endPoints, cpuEndPoint)

	rootComplexSwitch := switching.SwitchBuilder{}.
		WithEngine(c.engine).
		WithFreq(c.freq).
		WithArbiter(arbitration.NewXBarArbiter()).
		WithRoutingTable(routing.NewTable()).
		WithNumReqPerCycle(c.numReqPerCycle).
		Build(c.networkName + ".RootComplex")
	c.network.AddSwitch(rootComplexSwitch)
	c.switches = append(c.switches, rootComplexSwitch)
	c.rootComplexes = append(c.switches, rootComplexSwitch)
	c.switchNodes = append(c.switchNodes, &switchNode{
		parentSwitchNode:   nil,
		currSwitch:         rootComplexSwitch,
		localPortToParent:  nil,
		remotePortToParent: nil,
	})

	port := c.switchConnector.
		ConnectEndPointToSwitch(cpuEndPoint, rootComplexSwitch, c.freq)

	c.routeRootComplexToCPU(rootComplexSwitch, port, cpuPorts)

	return 0
}

func (c *Connector) routeRootComplexToCPU(
	rc *switching.Switch,
	rcPort akita.Port,
	cpuPorts []akita.Port,
) {
	rt := rc.GetRoutingTable()
	for _, p := range cpuPorts {
		rt.DefineRoute(p, rcPort)
	}
}

// AddSwitch adds a new switch connecting from an existing switch.
func (c *Connector) AddSwitch(baseSwitchID int) (switchID int) {
	switchID = len(c.switches)
	sw := switching.SwitchBuilder{}.
		WithEngine(c.engine).
		WithFreq(c.freq).
		WithRoutingTable(routing.NewTable()).
		WithArbiter(arbitration.NewXBarArbiter()).
		WithNumReqPerCycle(c.numReqPerCycle).
		Build(fmt.Sprintf("%s.Switch%d", c.networkName, switchID))
	c.network.AddSwitch(sw)
	c.switches = append(c.switches, sw)

	baseSwitch := c.switches[baseSwitchID]

	basePort, newPort := c.switchConnector.
		ConnectSwitches(baseSwitch, sw, c.freq)
	sw.GetRoutingTable().DefineDefaultRoute(newPort)
	c.switchNodes = append(c.switchNodes, &switchNode{
		parentSwitchNode:   c.switchNodes[baseSwitchID],
		currSwitch:         sw,
		localPortToParent:  newPort,
		remotePortToParent: basePort,
	})

	return switchID
}

// PlugInDevice connects a series of ports to a switch.
func (c *Connector) PlugInDevice(baseSwitchID int, devicePorts []akita.Port) {
	endPoint := switching.MakeEndPointBuilder().
		WithEngine(c.engine).
		WithFreq(c.freq).
		WithFlitByteSize(c.flitByteSize).
		WithEncodingOverhead(c.encodingOverhead).
		WithNetworkPortBufferSize(2 * c.numReqPerCycle).
		WithDevicePorts(devicePorts).
		WithNumReqPerCycle(c.numReqPerCycle).
		Build(fmt.Sprintf("%s.EndPoint%d", c.networkName, len(c.endPoints)))
	c.endPoints = append(c.endPoints, endPoint)
	c.network.AddEndPoint(endPoint)

	baseSwitch := c.switches[baseSwitchID]
	port := c.switchConnector.
		ConnectEndPointToSwitch(endPoint, baseSwitch, c.freq)

	sn := c.switchNodes[baseSwitchID]
	c.establishRouteToDevice(sn, port, devicePorts)
}

func (c *Connector) establishRouteToDevice(
	sn *switchNode,
	port akita.Port,
	devicePorts []akita.Port,
) {
	if sn == nil {
		return
	}

	rt := sn.currSwitch.GetRoutingTable()
	for _, dst := range devicePorts {
		rt.DefineRoute(dst, port)
	}

	c.establishRouteToDevice(sn.parentSwitchNode,
		sn.remotePortToParent, devicePorts)
}
