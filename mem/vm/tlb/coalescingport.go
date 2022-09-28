package tlb

import (
	"sync"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem/device"
)

// CoalescingPort is a type of port that can hold at most a certain number
// of messages.
type CoalescingPort struct {
	akita.HookableBase

	name string
	comp akita.Component
	conn akita.Connection

	buf          []akita.Msg
	bufLock      sync.RWMutex
	bufCapacity  int
	portBusy     bool
	portBusyLock sync.RWMutex
	mshr         mshr
	numReqs      int // when does numReqs get decremented?
}

// SetConnection sets which connection plugged in to this port.
func (p *CoalescingPort) SetConnection(conn akita.Connection) {
	p.conn = conn
}

// Component returns the owner component of the port.
func (p *CoalescingPort) Component() akita.Component {
	return p.comp
}

// Name returns the name of the port.
func (p *CoalescingPort) Name() string {
	return p.name
}

// Send is used to send a message out from a component

// THIS CODE IS NOT YET COMPLETE
// WHAT HAPPENS WHEN WE CAN'T PUSH ALL THE RESPONSES IN ONE GO?
func (p *CoalescingPort) Send(msg akita.Msg) *akita.SendError {
	// err := p.conn.Send(msg)
	rsp := msg.(*device.TranslationRsp)
	page := rsp.Page
	PID := page.PID
	vAddr := page.VAddr
	me := p.mshr.Query(PID, vAddr)
	if me == nil {
		panic("something went wrong")
	}
	var err *akita.SendError
	for len(me.Requests) > 0 {
		req := me.Requests[0]
		tempRsp := device.TranslationRspBuilder{}.
			WithSendTime(rsp.SendTime). // this is fine because a new response will be generated if any of the sends fails and this response will contain the new (latest) sendtime
			WithSrc(rsp.Src).
			WithDst(req.Src).
			WithRspTo(req.ID).
			WithPage(page).
			WithAccessResult(rsp.HitOrMiss).
			WithSrcL2TLB(rsp.SrcL2TLB).
			Build()
		err = p.conn.Send(tempRsp)
		if err != nil {
			break
		}
		p.numReqs--
		me.Requests = me.Requests[1:]
	}
	if len(me.Requests) == 0 {
		p.mshr.Remove(PID, vAddr)
	}
	return err
}

// Recv is used to deliver a message to a component
func (p *CoalescingPort) Recv(msg akita.Msg) *akita.SendError {
	p.bufLock.Lock()
	if p.numReqs >= p.bufCapacity {
		p.portBusyLock.Lock()
		p.portBusy = true
		p.portBusyLock.Unlock()
		p.bufLock.Unlock()
		return akita.NewSendError()
	}
	// if len(p.buf) >= p.bufCapacity {
	// 	p.portBusyLock.Lock()
	// 	p.portBusy = true
	// 	p.portBusyLock.Unlock()
	// 	p.bufLock.Unlock()
	// 	return NewSendError()
	// }

	hookCtx := akita.HookCtx{
		Domain: p,
		Now:    msg.Meta().RecvTime,
		Pos:    akita.HookPosPortMsgRecvd,
		Item:   msg,
	}
	// what is the hook doing?????????????????? There are no hooks (as there shouldn't be)
	p.InvokeHook(hookCtx)
	req := msg.(*device.TranslationReq)
	me := p.mshr.Query(req.PID, req.VAddr)
	if me != nil {
		me.Requests = append(me.Requests, req)
	} else {
		me = p.mshr.Add(req.PID, req.VAddr)
		me.Requests = append(me.Requests, req)
		me.reqToBottom = req
		p.buf = append(p.buf, msg)
		// only append to the buffer if this is a new message
	}
	p.numReqs++
	// increment the number of requests in all cases
	p.bufLock.Unlock()

	if p.comp != nil {
		p.comp.NotifyRecv(msg.Meta().RecvTime, p)
	}
	return nil
}

// Retrieve is used by the component to take a message from the incoming buffer
// Retrieve will pick up a request from buf, which is exactly what we want
// since all unique requests are stored in buf
func (p *CoalescingPort) Retrieve(now akita.VTimeInSec) akita.Msg {
	p.bufLock.Lock()
	defer p.bufLock.Unlock()

	if len(p.buf) == 0 {
		return nil
	}

	msg := p.buf[0]
	p.buf = p.buf[1:]
	hookCtx := akita.HookCtx{
		Domain: p,
		Now:    now,
		Pos:    akita.HookPosPortMsgRetrieve,
		Item:   msg,
	}
	p.InvokeHook(hookCtx)

	p.portBusyLock.Lock()
	if p.portBusy {
		p.portBusy = false
		p.conn.NotifyAvailable(now, p)
	}
	p.portBusyLock.Unlock()

	return msg
}

// Peek returns the first message in the port without removing it.
func (p *CoalescingPort) Peek() akita.Msg {
	p.bufLock.RLock()
	defer p.bufLock.RUnlock()

	if len(p.buf) == 0 {
		return nil
	}

	msg := p.buf[0]
	return msg
}

// NotifyAvailable is called by the connection to notify the port that the
// connection is available again
func (p *CoalescingPort) NotifyAvailable(now akita.VTimeInSec) {
	if p.comp != nil {
		p.comp.NotifyPortFree(now, p)
	}
}

// GetBuffer returns the request queue of the port
func (p *CoalescingPort) GetBuffer() []akita.Msg {
	return p.buf
}

// NewCoalescingPort creates a new port that works for the provided component
func NewCoalescingPort(
	comp akita.Component,
	capacity int,
	name string,
) *CoalescingPort {
	p := new(CoalescingPort)
	p.comp = comp
	p.bufCapacity = capacity
	p.mshr = newMSHR(512)
	p.name = name
	return p
}
