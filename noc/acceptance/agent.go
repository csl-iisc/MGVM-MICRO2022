package acceptance

import (
	"fmt"

	"gitlab.com/akita/akita"
)

// Agent can send and receive request.
type Agent struct {
	*akita.TickingComponent
	test         *Test
	Ports        []akita.Port
	MsgsToSend   []akita.Msg
	sendBytes    uint64
	recvBytes    uint64
	totalLatency akita.VTimeInSec
	numRecvMsgs  uint64
}

// NewAgent creates a new agent.
func NewAgent(
	engine akita.Engine,
	freq akita.Freq,
	name string,
	numPorts int,
	test *Test,
) *Agent {
	a := &Agent{}
	a.test = test
	a.TickingComponent = akita.NewTickingComponent(name, engine, freq, a)
	for i := 0; i < numPorts; i++ {
		p := akita.NewLimitNumMsgPort(a, 24, fmt.Sprintf("%s.Port%d", name, i))
		a.Ports = append(a.Ports, p)
	}
	return a
}

// Tick tries to receive requests and send requests out.
func (a *Agent) Tick(now akita.VTimeInSec) bool {
	madeProgress := false
	for i := 0; i < 12; i++ {
		madeProgress = a.send(now) || madeProgress
	}
	for i := 0; i < 12; i++ {
		madeProgress = a.recv(now) || madeProgress
	}
	return madeProgress
}

func (a *Agent) send(now akita.VTimeInSec) bool {
	if len(a.MsgsToSend) == 0 {
		return false
	}

	msg := a.MsgsToSend[0]
	msg.Meta().SendTime = now
	err := msg.Meta().Src.Send(msg)
	if err == nil {
		a.MsgsToSend = a.MsgsToSend[1:]
		a.sendBytes += uint64(msg.Meta().TrafficBytes)
		return true
	}

	return false
}

func (a *Agent) recv(now akita.VTimeInSec) bool {
	madeProgress := false
	for _, port := range a.Ports {
		msg := port.Retrieve(now)
		if msg != nil {
			a.test.receiveMsg(msg, port)
			a.recvBytes += uint64(msg.Meta().TrafficBytes)
			a.totalLatency += now - msg.Meta().SendTime
			a.numRecvMsgs += 1
			madeProgress = true
		}
	}
	return madeProgress
}
