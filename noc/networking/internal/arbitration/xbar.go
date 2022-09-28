package arbitration

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/noc"
	"gitlab.com/akita/util"
)

// NewXBarArbiter creates a new XBar arbiter.
func NewXBarArbiter() Arbiter {
	return &xbarArbiter{}
}

type xbarArbiter struct {
	buffers    []util.Buffer
	nextPortID int
}

func (a *xbarArbiter) AddBuffer(buf util.Buffer) {
	a.buffers = append(a.buffers, buf)
}

func (a *xbarArbiter) Arbitrate(now akita.VTimeInSec) []util.Buffer {
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
