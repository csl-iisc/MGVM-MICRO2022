package akita

import (
	// "fmt"
	"strings"
)

type directConnectionEnd struct {
	port    Port
	buf     []Msg
	bufSize int
	busy    bool
}

// DirectConnection connects two components without latency
type DirectConnection struct {
	*TickingComponent

	engine     Engine
	nextPortID int
	ports      []Port
	ends       map[Port]*directConnectionEnd

	L2Port Port
}

// PlugIn marks the port connects to this DirectConnection.
func (c *DirectConnection) PlugIn(port Port, sourceSideBufSize int) {
	c.Lock()
	defer c.Unlock()

	c.ports = append(c.ports, port)
	end := &directConnectionEnd{}
	end.port = port
	end.bufSize = sourceSideBufSize
	c.ends[port] = end

	if strings.Contains(port.Name(), "L2TLB.TopPort") {
		c.L2Port = port
	}

	port.SetConnection(c)
}

// Unplug marks the port no longer connects to this DirectConnection.
func (c *DirectConnection) Unplug(port Port) {
	panic("not implemented")
}

// NotifyAvailable is called by a port to notify that the connection can
// deliver to the port again.
func (c *DirectConnection) NotifyAvailable(now VTimeInSec, port Port) {
	c.TickNow(now)
}

// Send of a DirectConnection schedules a DeliveryEvent immediately
func (c *DirectConnection) Send(msg Msg) *SendError {
	c.Lock()
	defer c.Unlock()

	c.msgMustBeValid(msg)

	var srcEnd *directConnectionEnd

	if strings.Contains(c.Name(), "L1TLB-L2TLB") {
		msgMeta := msg.Meta()
		// srcName := msgMeta.Src.Name()
		if msgMeta.PutInL2TLBBuffer {
			// if srcName == msgMeta.Dst.Name() && strings.Contains(srcName, "RTU") {
			srcEnd = c.ends[c.L2Port]
		} else {
			srcEnd = c.ends[msgMeta.Src]
		}
	} else {
		srcEnd = c.ends[msg.Meta().Src]
	}

	if len(srcEnd.buf) >= srcEnd.bufSize {
		srcEnd.busy = true
		return NewSendError()
	}

	srcEnd.buf = append(srcEnd.buf, msg)

	c.TickNow(msg.Meta().SendTime)

	return nil
}

func (c *DirectConnection) msgMustBeValid(msg Msg) {
	c.portMustNotBeNil(msg.Meta().Src)
	c.portMustNotBeNil(msg.Meta().Dst)
	c.portMustBeConnected(msg.Meta().Src)
	c.portMustBeConnected(msg.Meta().Dst)
	c.srcDstMustNotBeTheSame(msg)
}

func (c *DirectConnection) portMustNotBeNil(port Port) {
	if port == nil {
		panic("src or dst is not given")
	}
}

func (c *DirectConnection) portMustBeConnected(port Port) {
	if _, connected := c.ends[port]; !connected {
		panic("src or dst is not connected")
	}
}

func (c *DirectConnection) srcDstMustNotBeTheSame(msg Msg) {
	if msg.Meta().Src == msg.Meta().Dst && !strings.Contains(msg.Meta().Src.Name(), "RTU") {
		panic("sending back to src")
	}
}

func (c *DirectConnection) Tick(now VTimeInSec) bool {
	madeProgress := false
	// for {
	// madeProgress := false
	for i := 0; i < len(c.ports); i++ {
		portID := (i + c.nextPortID) % len(c.ports)
		port := c.ports[portID]
		end := c.ends[port]
		madeProgress = c.forwardMany(end, now) || madeProgress
	}
	// if !madeProgress {
	// break
	// }
	// }
	c.nextPortID = (c.nextPortID + 1) % len(c.ports)
	// return true
	return madeProgress
}

func (c *DirectConnection) forwardMany(
	end *directConnectionEnd,
	now VTimeInSec,
) bool {
	madeProgress := false
	for {
		if len(end.buf) == 0 {
			break
		}

		head := end.buf[0]
		head.Meta().RecvTime = now

		err := head.Meta().Dst.Recv(head)
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

func (c *DirectConnection) forwardOne(
	end *directConnectionEnd,
	now VTimeInSec,
) bool {
	madeProgress := false
	if len(end.buf) == 0 {
		return madeProgress
	}

	head := end.buf[0]
	head.Meta().RecvTime = now

	err := head.Meta().Dst.Recv(head)
	if err != nil {
		return madeProgress
	}

	madeProgress = true
	end.buf = end.buf[1:]

	if end.busy {
		end.port.NotifyAvailable(now)
		end.busy = false
	}

	return madeProgress
}

// NewDirectConnection creates a new DirectConnection object
func NewDirectConnection(
	name string,
	engine Engine,
	freq Freq,
) *DirectConnection {
	c := new(DirectConnection)
	c.TickingComponent = NewSecondaryTickingComponent(name, engine, freq, c)
	c.ends = make(map[Port]*directConnectionEnd)
	return c
}
