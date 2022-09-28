// Package chipnetwork for interconnect
package chipnetwork

import (
	"fmt"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/noc/networking/internal/arbitration"
	"gitlab.com/akita/noc/networking/internal/routing"
	"gitlab.com/akita/noc/networking/internal/switching"
)

type finalDests struct {
	endswitch *switching.Switch
	ports     []akita.Port
}

// Connector can connect devices into a PCIe network.
type Connector struct {
	networkName       string
	engine            akita.Engine
	freq              akita.Freq
	encodingOverhead  float64
	flitByteSize      int
	switchLatency     int
	numReqPerCycle    int
	switchConnector   *switching.SwitchConnector
	switches          []*switching.Switch
	endPoints         []*switching.EndPoint
	finalDestinations map[*switching.Switch]*finalDests
}

// NewConnector creates a new connector that can help configure PCIe networks.
func NewConnector() *Connector {
	c := &Connector{}
	c.encodingOverhead = 0
	c.flitByteSize = 64
	c.freq = 1 * akita.GHz
	c.numReqPerCycle = 1
	c.switchLatency = 1
	c.finalDestinations = make(map[*switching.Switch]*finalDests)

	return c
}

func (c *Connector) WithFreq(freq akita.Freq) *Connector {
	c.freq = freq
	return c
}

func (c *Connector) WithEngine(engine akita.Engine) *Connector {
	c.engine = engine
	return c
}

func (c *Connector) WithSwitchLatency(numCycles int) *Connector {
	c.switchLatency = numCycles
	return c
}

func (c *Connector) WithFlitByteSize(flitByteSize int) *Connector {
	c.flitByteSize = flitByteSize
	return c
}

func (c *Connector) WithNumReqPerCycle(numReqPerCycle int) *Connector {
	c.numReqPerCycle = numReqPerCycle
	return c
}

// WithNetworkName sets the name of the network and the prefix of all the
// component in the network.
func (c *Connector) WithNetworkName(name string) *Connector {
	c.networkName = name
	return c
}

// CreateNetwork creates a network.
func (c *Connector) CreateNetwork() {
	c.switchConnector = switching.NewSwitchConnector(c.engine)
	c.switchConnector.SetSwitchLatency(c.switchLatency)
	c.switchConnector.SetNumReqPerCycle(c.numReqPerCycle)
	c.switchConnector.SetBufferSize(2 * c.numReqPerCycle)
}

// PlugInChip connects a series of ports to a switch.
func (c *Connector) PlugInChip(devicePorts []akita.Port) {
	chipSwitch := switching.SwitchBuilder{}.
		WithEngine(c.engine).
		WithFreq(c.freq).
		WithArbiter(arbitration.NewXBarArbiter()).
		WithRoutingTable(routing.NewTable()).
		WithNumReqPerCycle(c.numReqPerCycle).
		Build(fmt.Sprintf("%s.Switch%d", c.networkName, len(c.switches)))
	c.switches = append(c.switches, chipSwitch)

	for i, dst := range devicePorts {
		endPoint := switching.MakeEndPointBuilder().
			WithEngine(c.engine).
			WithFreq(c.freq).
			WithFlitByteSize(c.flitByteSize).
			WithEncodingOverhead(c.encodingOverhead).
			WithNetworkPortBufferSize(2 * c.numReqPerCycle).
			WithDevicePorts(devicePorts[i : i+1]).
			WithNumReqPerCycle(c.numReqPerCycle).
			Build(fmt.Sprintf("%s.EndPoint%d.%s", c.networkName, len(c.endPoints), dst.Name()))
		c.endPoints = append(c.endPoints, endPoint)
		port := c.switchConnector.
			ConnectEndPointToSwitch(endPoint, chipSwitch, c.freq)
		rt := chipSwitch.GetRoutingTable()
		rt.DefineRoute(dst, port)
	}

	fd := &finalDests{chipSwitch, devicePorts}
	c.finalDestinations[chipSwitch] = fd

	// rt := chipSwitch.GetRoutingTable()
	// for _, dst := range devicePorts {
	// rt.DefineRoute(dst, port)
	// }
}

// MakeNetwork as the name suggests makes the network. haha.
func (c *Connector) MakeNetwork() {
	for i, switch1 := range c.switches {
		for j, switch2 := range c.switches {
			if i >= j {
				continue
			}
			port1, port2 := c.switchConnector.
				ConnectSwitches(switch1, switch2, c.freq)
			rt1 := switch1.GetRoutingTable()
			for _, dst := range c.finalDestinations[switch2].ports {
				rt1.DefineRoute(dst, port1)
			}
			rt2 := switch2.GetRoutingTable()
			for _, dst := range c.finalDestinations[switch1].ports {
				rt2.DefineRoute(dst, port2)
			}
		}
	}
}
