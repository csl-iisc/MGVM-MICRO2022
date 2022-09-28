//Package chipnetwork for interconnect
package chipnetwork

import (
	"fmt"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/noc/networking/internal/arbitration"
	"gitlab.com/akita/noc/networking/internal/routing"
	"gitlab.com/akita/noc/networking/internal/switching"
)

// type finalDests struct {
// 	endswitch *switching.Switch
// 	ports     []akita.Port
// }

// InterChipletConnector can connect devices into a PCIe network.
type InterChipletConnector struct {
	networkName      string
	engine           akita.Engine
	freq             akita.Freq
	encodingOverhead float64
	flitByteSize     int
	switchLatency    int
	numReqPerCycle   int
	chipSwitch       *switching.MultiFlitChipletSwitch
	switchConnector  *switching.ChipSwitchConnector
	endPoints        []*switching.EndPoint
	// finalDestinations map[*switching.Switch]*finalDests
}

// NewConnector creates a new InterChipletConnector that can help configure PCIe networks.
func NewInterChipletConnector() *InterChipletConnector {
	c := &InterChipletConnector{}
	c.encodingOverhead = 0
	c.flitByteSize = 64
	c.freq = 1 * akita.GHz
	c.numReqPerCycle = 1
	c.switchLatency = 1
	// c.finalDestinations = make(map[*switching.Switch]*finalDests)

	return c
}

func (c *InterChipletConnector) WithFreq(freq akita.Freq) *InterChipletConnector {
	c.freq = freq
	return c
}

func (c *InterChipletConnector) WithEngine(engine akita.Engine) *InterChipletConnector {
	c.engine = engine
	return c
}

func (c *InterChipletConnector) WithSwitchLatency(numCycles int) *InterChipletConnector {
	c.switchLatency = numCycles
	return c
}

func (c *InterChipletConnector) WithFlitByteSize(flitByteSize int) *InterChipletConnector {
	c.flitByteSize = flitByteSize
	return c
}

func (c *InterChipletConnector) WithNumReqPerCycle(numReqPerCycle int) *InterChipletConnector {
	c.numReqPerCycle = numReqPerCycle
	return c
}

// WithNetworkName sets the name of the network and the prefix of all the
// component in the network.
func (c *InterChipletConnector) WithNetworkName(name string) *InterChipletConnector {
	c.networkName = name
	return c
}

// CreateNetwork creates a network.
func (c *InterChipletConnector) CreateNetwork() {
	c.chipSwitch = switching.MultiFlitChipletSwitchBuilder{}.
		WithEngine(c.engine).
		WithFreq(c.freq).
		WithArbiter(arbitration.NewFlitAwareXBarArbiter()).
		WithRoutingTable(routing.NewTable()).
		WithNumReqPerCycle(c.numReqPerCycle).
		WithNumChiplets(4).
		WithMaxOutgoingReqsPerChiplet(12).
		Build(fmt.Sprintf("%s.Switch%d", c.networkName, 0))
	c.switchConnector = switching.NewChipSwitchConnector(c.engine)
	c.switchConnector.SetSwitchLatency(c.switchLatency)
	c.switchConnector.SetNumReqPerCycle(c.numReqPerCycle)
	c.switchConnector.SetBufferSize(2 * c.numReqPerCycle)
}

// PlugInChip connects a series of ports to a switch.
func (c *InterChipletConnector) PlugInChip(devicePorts []akita.Port) {
	// chipSwitch := switching.MultiFlitChipletSwitchBuilder{}.
	// 	WithEngine(c.engine).
	// 	WithFreq(c.freq).
	// 	WithArbiter(arbitration.NewXBarArbiter()).
	// 	WithRoutingTable(routing.NewTable()).
	// 	WithNumReqPerCycle(c.numReqPerCycle).
	// 	Build(fmt.Sprintf("%s.Switch%d", c.networkName, len(c.switches)))
	// c.switches = append(c.switches, chipSwitch)

	for _, dst := range devicePorts {
		c.switchConnector.ConnectPortToSwitch(dst, c.chipSwitch, c.freq)
		// endPoint := switching.MakeEndPointBuilder().
		// 	WithEngine(c.engine).
		// 	WithFreq(c.freq).
		// 	WithFlitByteSize(c.flitByteSize).
		// 	WithEncodingOverhead(c.encodingOverhead).
		// 	WithNetworkPortBufferSize(2 * c.numReqPerCycle).
		// 	WithDevicePorts(devicePorts[i : i+1]).
		// 	WithNumReqPerCycle(c.numReqPerCycle).
		// 	Build(fmt.Sprintf("%s.EndPoint%d.%s", c.networkName, len(c.endPoints), dst.Name()))
		// c.endPoints = append(c.endPoints, endPoint)
		// port := c.switchConnector.
		// 	ConnectEndPointToSwitch(endPoint, c.chipSwitch, c.freq)
		// rt := c.chipSwitch.GetRoutingTable()
		// rt.DefineRoute(dst, port)
	}

	// fd := &finalDests{chipSwitch, devicePorts}
	// c.finalDestinations[chipSwitch] = fd

	// rt := chipSwitch.GetRoutingTable()
	// for _, dst := range devicePorts {
	// rt.DefineRoute(dst, port)
	// }
}

// MakeNetwork as the name suggests makes the network. haha.
func (c *InterChipletConnector) MakeNetwork() {
	// for i, switch1 := range c.switches {
	// 	for j, switch2 := range c.switches {
	// 		if i >= j {
	// 			continue
	// 		}
	// 		port1, port2 := c.switchConnector.
	// 			ConnectSwitches(switch1, switch2, c.freq)
	// 		rt1 := switch1.GetRoutingTable()
	// 		for _, dst := range c.finalDestinations[switch2].ports {
	// 			rt1.DefineRoute(dst, port1)
	// 		}
	// 		rt2 := switch2.GetRoutingTable()
	// 		for _, dst := range c.finalDestinations[switch1].ports {
	// 			rt2.DefineRoute(dst, port2)
	// 		}
	// 	}
	// }
}
