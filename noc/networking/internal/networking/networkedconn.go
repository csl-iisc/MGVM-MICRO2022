// Package networking provides the implementation of a NetworkedConnection.
package networking

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/noc/networking/internal/switching"
)

// A NetworkedConnection is a complex connection that is composed of a certain
// number of endpoints, switches, and links (implemented with simpler
// connections).
type NetworkedConnection struct {
	akita.HookableBase

	endPoints         []*switching.EndPoint
	switches          []*switching.Switch
	portToEndPointMap map[akita.Port]*switching.EndPoint
}

// NewNetworkedConnection creates a networked connection.
func NewNetworkedConnection() *NetworkedConnection {
	c := &NetworkedConnection{
		portToEndPointMap: make(map[akita.Port]*switching.EndPoint),
	}
	return c
}

// AddEndPoint adds an End Point to the network.
func (c *NetworkedConnection) AddEndPoint(endPoint *switching.EndPoint) {
	c.endPoints = append(c.endPoints, endPoint)
	for _, p := range endPoint.DevicePorts {
		c.portToEndPointMap[p] = endPoint
	}
}

// AddSwitch adds a new switch to the network.
func (c *NetworkedConnection) AddSwitch(sw *switching.Switch) {
	c.switches = append(c.switches, sw)
}
