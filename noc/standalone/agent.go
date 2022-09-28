package standalone

import (
	"log"
	"reflect"

	"gitlab.com/akita/akita"
)

// TrafficMsg is a type of msguests that only used in standalone network test.
// It has a byte size, but we do not care about the information it carries.
type TrafficMsg struct {
	akita.MsgMeta
}

// Meta returns the meta data of the message.
func (m *TrafficMsg) Meta() *akita.MsgMeta {
	return &m.MsgMeta
}

// NewTrafficMsg creates a new traffic message
func NewTrafficMsg(src, dst akita.Port, byteSize int) *TrafficMsg {
	msg := new(TrafficMsg)
	msg.Src = src
	msg.Dst = dst
	msg.TrafficBytes = byteSize
	return msg
}

// StartSendEvent is an event that triggers an agent to send a message.
type StartSendEvent struct {
	*akita.EventBase
	Msg *TrafficMsg
}

// NewStartSendEvent creates a new StartSendEvent.
func NewStartSendEvent(
	time akita.VTimeInSec,
	src, dst *Agent,
	byteSize int,
	trafficClass int,
) *StartSendEvent {
	e := new(StartSendEvent)
	e.EventBase = akita.NewEventBase(time, src)
	e.Msg = NewTrafficMsg(src.ToOut, dst.ToOut, byteSize)
	e.Msg.Meta().TrafficClass = trafficClass
	return e
}

// Agent is a component that connects the network. It can send and receive
// msguests to/ from the network.
type Agent struct {
	*akita.TickingComponent

	ToOut akita.Port

	Buffer []*TrafficMsg
}

// NotifyRecv notifies that a port has received a message.
func (a *Agent) NotifyRecv(now akita.VTimeInSec, port akita.Port) {
	a.ToOut.Retrieve(now)
	a.TickLater(now)
}

// Handle defines how an agent handles events.
func (a *Agent) Handle(e akita.Event) error {
	switch e := e.(type) {
	case *StartSendEvent:
		a.handleStartSendEvent(e)
	case akita.TickEvent:
		a.TickingComponent.Handle(e)
	default:
		log.Panicf("cannot handle event of type %s", reflect.TypeOf(e))
	}
	return nil
}

func (a *Agent) handleStartSendEvent(e *StartSendEvent) {
	a.Buffer = append(a.Buffer, e.Msg)
	a.TickLater(e.Time())
}

func (a *Agent) Tick(now akita.VTimeInSec) bool {
	return a.sendDataOut(now)
}

func (a *Agent) sendDataOut(now akita.VTimeInSec) bool {
	if len(a.Buffer) == 0 {
		return false
	}

	msg := a.Buffer[0]
	msg.Meta().SendTime = now
	err := a.ToOut.Send(msg)
	if err == nil {
		a.Buffer = a.Buffer[1:]
		return true
	}
	return false
}

// NewAgent creates a new agent.
func NewAgent(name string, engine akita.Engine) *Agent {
	a := new(Agent)
	a.TickingComponent = akita.NewTickingComponent(name, engine, 1*akita.GHz, a)

	a.ToOut = akita.NewLimitNumMsgPort(a, 4, name+".ToOut")

	return a
}
