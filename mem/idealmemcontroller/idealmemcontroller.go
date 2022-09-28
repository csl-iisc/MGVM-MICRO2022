package idealmemcontroller

import (
	"log"
	"reflect"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem"
	"gitlab.com/akita/util/tracing"
)

type readRespondEvent struct {
	*akita.EventBase

	req *mem.ReadReq
}

func newReadRespondEvent(time akita.VTimeInSec, handler akita.Handler,
	req *mem.ReadReq,
) *readRespondEvent {
	return &readRespondEvent{akita.NewEventBase(time, handler), req}
}

type writeRespondEvent struct {
	*akita.EventBase

	req *mem.WriteReq
}

func newWriteRespondEvent(time akita.VTimeInSec, handler akita.Handler,
	req *mem.WriteReq,
) *writeRespondEvent {
	return &writeRespondEvent{akita.NewEventBase(time, handler), req}
}

// An Comp is an ideal memory controller that can perform read and write
//
// Ideal memory controller always respond to the request in a fixed number of
// cycles. There is no limitation on the concurrency of this unit.
type Comp struct {
	*akita.TickingComponent

	ToTop              akita.Port
	Storage            *mem.Storage
	Latency            int
	AddressConverter   AddressConverter
	MaxNumTransaction  int
	currNumTransaction int
}

// Handle defines how the Comp handles event
func (c *Comp) Handle(e akita.Event) error {
	switch e := e.(type) {
	case *readRespondEvent:
		return c.handleReadRespondEvent(e)
	case *writeRespondEvent:
		return c.handleWriteRespondEvent(e)
	case akita.TickEvent:
		return c.TickingComponent.Handle(e)
	default:
		log.Panicf("cannot handle event of %s", reflect.TypeOf(e))
	}

	return nil
}

// Tick updates ideal memory controller state.
func (c *Comp) Tick(now akita.VTimeInSec) bool {
	if c.currNumTransaction >= c.MaxNumTransaction {
		return false
	}

	msg := c.ToTop.Retrieve(now)
	if msg == nil {
		return false
	}

	tracing.TraceReqReceive(msg, now, c)
	c.currNumTransaction++

	switch msg := msg.(type) {
	case *mem.ReadReq:
		c.handleReadReq(now, msg)
		return true
	case *mem.WriteReq:
		c.handleWriteReq(now, msg)
		return true
	default:
		log.Panicf("cannot handle request of type %s", reflect.TypeOf(msg))
	}

	return false
}

func (c *Comp) handleReadReq(now akita.VTimeInSec, req *mem.ReadReq) {
	timeToSchedule := c.Freq.NCyclesLater(c.Latency, now)
	respondEvent := newReadRespondEvent(timeToSchedule, c, req)
	c.Engine.Schedule(respondEvent)
}

func (c *Comp) handleWriteReq(now akita.VTimeInSec, req *mem.WriteReq) {
	timeToSchedule := c.Freq.NCyclesLater(c.Latency, now)
	respondEvent := newWriteRespondEvent(timeToSchedule, c, req)
	c.Engine.Schedule(respondEvent)
}

func (c *Comp) handleReadRespondEvent(e *readRespondEvent) error {
	now := e.Time()
	req := e.req

	addr := req.Address
	if c.AddressConverter != nil {
		addr = c.AddressConverter.ConvertExternalToInternal(addr)
	}

	data, err := c.Storage.Read(addr, req.AccessByteSize)
	if err != nil {
		log.Panic(err)
	}

	rsp := mem.DataReadyRspBuilder{}.
		WithSendTime(now).
		WithSrc(c.ToTop).
		WithDst(req.Src).
		WithRspTo(req.ID).
		WithData(data).
		Build()

	networkErr := c.ToTop.Send(rsp)
	if networkErr != nil {
		retry := newReadRespondEvent(c.Freq.NextTick(now), c, req)
		c.Engine.Schedule(retry)
		return nil
	}

	tracing.TraceReqComplete(req, now, c)
	c.currNumTransaction--
	c.TickLater(now)

	return nil
}

func (c *Comp) handleWriteRespondEvent(e *writeRespondEvent) error {
	now := e.Time()
	req := e.req

	rsp := mem.WriteDoneRspBuilder{}.
		WithSendTime(now).
		WithSrc(c.ToTop).
		WithDst(req.Src).
		WithRspTo(req.ID).
		Build()

	networkErr := c.ToTop.Send(rsp)
	if networkErr != nil {
		retry := newWriteRespondEvent(c.Freq.NextTick(now), c, req)
		c.Engine.Schedule(retry)
		return nil
	}

	addr := req.Address

	if c.AddressConverter != nil {
		addr = c.AddressConverter.ConvertExternalToInternal(addr)
	}

	if req.DirtyMask == nil {
		err := c.Storage.Write(addr, req.Data)
		if err != nil {
			log.Panic(err)
		}
	} else {
		data, err := c.Storage.Read(addr, uint64(len(req.Data)))
		if err != nil {
			panic(err)
		}
		for i := 0; i < len(req.Data); i++ {
			if req.DirtyMask[i] == true {
				data[i] = req.Data[i]
			}
		}
		err = c.Storage.Write(addr, data)
		if err != nil {
			panic(err)
		}
	}

	tracing.TraceReqComplete(req, now, c)
	c.currNumTransaction--
	c.TickLater(now)

	return nil
}

// New creates a new ideal memory controller
func New(
	name string,
	engine akita.Engine,
	capacity uint64,
) *Comp {
	c := new(Comp)
	c.TickingComponent = akita.NewTickingComponent(name, engine, 1*akita.GHz, c)
	c.Latency = 100
	c.MaxNumTransaction = 100 //100

	c.Storage = mem.NewStorage(capacity)
	//TODO: is this sufficient number of ports?
	c.ToTop = akita.NewLimitNumMsgPort(c, 16, name+".ToTop")
	return c
}
