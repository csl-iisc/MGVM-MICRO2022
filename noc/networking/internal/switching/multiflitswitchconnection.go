package switching

import (
	"strings"

	"gitlab.com/akita/akita"
)

type MultiFlitSwitchConnectionEnd struct {
	port    akita.Port
	buf     []akita.Msg
	bufSize int
	busy    bool
}

func newMultiFlitSwitchConnectionEnd(port akita.Port, sourceSideBufSize int) *MultiFlitSwitchConnectionEnd {
	end := &MultiFlitSwitchConnectionEnd{}
	end.port = port
	end.bufSize = sourceSideBufSize
	return end
}

// MultiFlitSwitchConnection connects two components without latency
type MultiFlitSwitchConnection struct {
	*akita.TickingComponent

	engine akita.Engine
	// nextPortID int
	// ports      []akita.Port
	// ends       map[akita.Port]*MultiFlitSwitchConnectionEnd
	switchPort akita.Port
	switchEnd  *MultiFlitSwitchConnectionEnd
	devicePort akita.Port
	deviceEnd  *MultiFlitSwitchConnectionEnd
}

// PlugIn marks the port connects to this MultiFlitSwitchConnection.
func (c *MultiFlitSwitchConnection) PlugIn(port akita.Port, sourceSideBufSize int) {
	c.Lock()
	defer c.Unlock()
	if strings.Contains(port.Name(), "Switch0") && c.switchPort == nil {
		c.switchPort = port
		c.switchEnd = newMultiFlitSwitchConnectionEnd(port, sourceSideBufSize)

	} else if c.devicePort == nil {
		c.devicePort = port
		c.deviceEnd = newMultiFlitSwitchConnectionEnd(port, sourceSideBufSize)

	} else {
		panic("something went wrong")
	}
	port.SetConnection(c)
}

// Unplug marks the port no longer connects to this MultiFlitSwitchConnection.
func (c *MultiFlitSwitchConnection) Unplug(port akita.Port) {
	panic("not implemented")
}

// NotifyAvailable is called by a port to notify that the connection can
// deliver to the port again.
func (c *MultiFlitSwitchConnection) NotifyAvailable(now akita.VTimeInSec, port akita.Port) {
	c.TickNow(now)
}

// Send of a MultiFlitSwitchConnection schedules a DeliveryEvent immediately
func (c *MultiFlitSwitchConnection) Send(msg akita.Msg) *akita.SendError {
	c.Lock()
	defer c.Unlock()
	c.msgMustBeValid(msg)

	// srcEnd := c.ends[msg.Meta().Src]
	src := msg.Meta().Src
	dst := msg.Meta().Dst
	// 1 is the device port
	var srcEnd *MultiFlitSwitchConnectionEnd
	if dst == c.devicePort {
		srcEnd = c.switchEnd //ends[c.ports[0]]
	} else if src == c.devicePort { // forward to the switch port
		srcEnd = c.deviceEnd //c.ends[c.ports[1]]
	} else {
		srcEnd = c.deviceEnd
		// panic("what's happening")
	}
	// if srcEnd == nil {
	// 	srcEnd = c.ends[c.ports[0]]
	// }

	if len(srcEnd.buf) >= srcEnd.bufSize {
		srcEnd.busy = true
		return akita.NewSendError()
	}

	srcEnd.buf = append(srcEnd.buf, msg)

	c.TickNow(msg.Meta().SendTime)

	return nil
}

func (c *MultiFlitSwitchConnection) msgMustBeValid(msg akita.Msg) {
	c.portMustNotBeNil(msg.Meta().Src)
	c.portMustNotBeNil(msg.Meta().Dst)
	// c.portMustBeConnected(msg.Meta().Src)
	// c.portMustBeConnected(msg.Meta().Dst)
	c.srcDstMustNotBeTheSame(msg)
}

func (c *MultiFlitSwitchConnection) portMustNotBeNil(port akita.Port) {
	if port == nil {
		panic("src or dst is not given")
	}
}

// func (c *MultiFlitSwitchConnection) portMustBeConnected(port akita.Port) {
// 	if _, connected := c.ends[port]; !connected {
// 		panic("src or dst is not connected")
// 	}
// }

func (c *MultiFlitSwitchConnection) srcDstMustNotBeTheSame(msg akita.Msg) {
	// if msg.Meta().Src == msg.Meta().Dst {
	// 	panic("sending bacxk to src")
	// }
}

func (c *MultiFlitSwitchConnection) Tick(now akita.VTimeInSec) bool {
	madeProgress := false
	// for i := 0; i < len(c.ports); i++ {
	// 	portID := (i + c.nextPortID) % len(c.ports)
	// 	port := c.ports[portID]
	// 	end := c.ends[port]
	// 	madeProgress = c.forwardMany(end, portID, now) || madeProgress
	// }
	madeProgress = c.forwardMany(c.switchEnd, c.switchPort, now) || madeProgress
	madeProgress = c.forwardMany(c.deviceEnd, c.devicePort, now) || madeProgress
	// not sure nextPortID is needed since it might only be needed when multiple ports are trying to
	// send to the same destination since then they will compete for buffer space and the
	// port that is always going first will get an unfair advantage. here we only have two ports
	// sending to each other and so they do not compete for the same buffer space
	// c.nextPortID = (c.nextPortID + 1) % len(c.ports)
	return madeProgress
}

func (c *MultiFlitSwitchConnection) forwardMany(
	end *MultiFlitSwitchConnectionEnd, port akita.Port,
	now akita.VTimeInSec,
) bool {
	madeProgress := false
	for {
		if len(end.buf) == 0 {
			break
		}

		head := end.buf[0]
		head.Meta().RecvTime = now
		var dst akita.Port
		if port == c.devicePort {
			dst = c.switchPort
		} else if port == c.switchPort {
			dst = c.devicePort
		} else {
			panic("something is wrong")
		}
		err := dst.Recv(head)
		if err != nil {
			break
		}

		madeProgress = true
		end.buf = end.buf[1:]

		if end.busy {
			end.port.NotifyAvailable(now)
			end.busy = false
		}
	}

	return madeProgress
}

// NewMultiFlitSwitchConnection creates a new MultiFlitSwitchConnection object
func NewMultiFlitSwitchConnection(
	name string,
	engine akita.Engine,
	freq akita.Freq,
) *MultiFlitSwitchConnection {
	c := new(MultiFlitSwitchConnection)
	c.TickingComponent = akita.NewSecondaryTickingComponent(name, engine, freq, c)
	// c.ends = make(map[akita.Port]*MultiFlitSwitchConnectionEnd)
	return c
}
