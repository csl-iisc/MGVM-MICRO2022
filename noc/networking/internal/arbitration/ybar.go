package arbitration

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/noc"
	"gitlab.com/akita/util"
)

// NewYBarArbiter creates a new YBar arbiter.
func NewYBarArbiter() Arbiter {
	return &ybarArbiter{}
}

type ybarArbiter struct {
	msgsSentOut []int
	maxMsgs     int
	resetCycles int
	buffers     []util.Buffer
	nextPortID  int
}

func (a *ybarArbiter) AddBuffer(buf util.Buffer) {
	a.buffers = append(a.buffers, buf)
}

func (a *ybarArbiter) Arbitrate(now akita.VTimeInSec) []util.Buffer {
	startingPortID := a.nextPortID
	selectedPort := make([]util.Buffer, 0)
	occupiedOutputPort := make(map[util.Buffer]bool)

	for i := 0; i < len(a.buffers); i++ {
		currPortID := (startingPortID + i) % len(a.buffers)
		buf := a.buffers[currPortID]
		item := buf.Peek()
		if item == nil {
			continue
		}

		flit := item.(*noc.Flit)
		if _, ok := occupiedOutputPort[flit.OutputBuf]; ok {
			continue
		}

		selectedPort = append(selectedPort, buf)
		occupiedOutputPort[flit.OutputBuf] = true
	}

	a.nextPortID = (a.nextPortID + 1) % len(a.buffers)

	return selectedPort
}
