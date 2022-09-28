package arbitration

import (
	"strconv"
	"strings"
	"sync"

	// "fmt"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/noc"
	"gitlab.com/akita/util"
	// "gitlab.com/akita/mem"
)

type msgReservation struct {
	msg          akita.Msg
	numFlitInMsg int
	flitsSent    int
}

func newMessageReservation(flit *noc.Flit) *msgReservation {
	m := &msgReservation{
		msg:          flit.Msg,
		numFlitInMsg: flit.NumFlitInMsg,
		flitsSent:    0,
	}
	return m
}

type bufReservation struct {
	msgs          map[akita.Msg]*msgReservation
	reservedSlots int
}

func newBufReservation() *bufReservation {
	r := &bufReservation{
		msgs:          make(map[akita.Msg]*msgReservation),
		reservedSlots: 0,
	}
	return r
}

// NewFlitAwareXBarArbiter creates a new FlitAwareXBar arbiter.
func NewFlitAwareXBarArbiter() *FlitAwareXBarArbiter {
	return &FlitAwareXBarArbiter{
		outBufReservations:        make(map[util.Buffer]*bufReservation),
		numChiplets:               4,
		outgoingReqsPerChiplet:    make([]int, 4),
		maxOutgoingReqsPerChiplet: 12, // this needs to be changed
	}
}

type FlitAwareXBarArbiter struct {
	buffers                   []util.Buffer
	nextPortID                int
	outBufReservations        map[util.Buffer]*bufReservation
	mapLock                   sync.RWMutex
	numChiplets               int
	outgoingReqsPerChiplet    []int
	maxOutgoingReqsPerChiplet int
}

func getChipletFromFlitSrc(flit *noc.Flit) (chiplet int) {
	src := strings.Split(flit.Src.Name(), ".")[1]
	chiplet, err := strconv.Atoi(strings.Split(src, "_")[1])
	if err != nil {
		panic("something went wrong")
	}
	return
}

func (a *FlitAwareXBarArbiter) hasReservation(flit *noc.Flit) bool {
	a.mapLock.Lock()
	defer a.mapLock.Unlock()
	outBuf := flit.OutputBuf
	bufReservation, ok := a.outBufReservations[outBuf]
	if !ok {
		return ok
	}
	_, ok = bufReservation.msgs[flit.Msg]
	return ok
}

func (a *FlitAwareXBarArbiter) makeReservation(flit *noc.Flit) bool {
	a.mapLock.Lock()
	defer a.mapLock.Unlock()
	outBuf := flit.OutputBuf
	outBufReservation, ok := a.outBufReservations[outBuf]
	if !ok {
		r := newBufReservation()
		a.outBufReservations[outBuf] = r
		outBufReservation = r
	}
	msgReservations := outBufReservation.msgs
	if _, ok := msgReservations[flit.Msg]; ok {
		panic("reservation for this msg already exists")

	}
	// it should be alright to exceed the capacity of the buffer for one reservation
	freeSlots := flit.OutputBuf.Capacity() - flit.OutputBuf.Size() - outBufReservation.reservedSlots
	if freeSlots < 0 {
		panic("something went wrong")
	}
	if freeSlots >= flit.NumFlitInMsg {
		m := newMessageReservation(flit)
		outBufReservation.msgs[flit.Msg] = m
		outBufReservation.reservedSlots += flit.NumFlitInMsg
		if _, ok := a.outBufReservations[flit.OutputBuf].msgs[flit.Msg]; !ok {
			panic("hello")
		}
		return true
	}
	return false
}

func (a *FlitAwareXBarArbiter) markOneFlitSent(flit *noc.Flit) {
	a.mapLock.Lock()
	defer a.mapLock.Unlock()
	outBuf := flit.OutputBuf
	outBufReservation := a.outBufReservations[outBuf]
	msgReservation, ok := outBufReservation.msgs[flit.Msg]
	if !ok {
		panic("Msg should be present")
	}
	msgReservation.flitsSent++
	outBufReservation.reservedSlots--
	// fmt.Println("mark one flit sent:", outBufReservation.reservedSlots)
	if msgReservation.flitsSent == msgReservation.numFlitInMsg {
		// fmt.Println("deleting msg", flit, flit.Msg)
		delete(outBufReservation.msgs, flit.Msg)
	}
}

func (a *FlitAwareXBarArbiter) AddBuffer(buf util.Buffer) {
	a.buffers = append(a.buffers, buf)
}

func (a *FlitAwareXBarArbiter) Reset() {
	for i := 0; i < a.numChiplets; i++ {
		a.outgoingReqsPerChiplet[i] = 0
	}
}

func (a *FlitAwareXBarArbiter) Arbitrate(now akita.VTimeInSec) []util.Buffer {
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
		// fmt.Println(flit, flit.OutputBuf, flit.OutputBuf.Size(), flit.NumFlitInMsg)
		// b, ok := a.outBufReservations[flit.OutputBuf]
		// if ok {
		// 	fmt.Println(b.reservedSlots)
		// }
		// switch t := flit.Msg.(type) {
		// case *mem.WriteReq:
		// 	if t.Address == 8590332864 {
		// 		fmt.Println("here")
		// 	}
		// }
		chiplet := getChipletFromFlitSrc(flit)
		if a.outgoingReqsPerChiplet[chiplet] >= a.maxOutgoingReqsPerChiplet {
			continue
		}
		if !flit.OutputBuf.CanPush() {
			continue
		}
		if _, ok := occupiedOutputPort[flit.OutputBuf]; ok {
			continue
		}
		// fmt.Println("here")
		// if flit.NumFlitInMsg > 1 {
		// fmt.Println(i)
		if a.hasReservation(flit) || a.makeReservation(flit) {
			// if _, ok := occupiedOutputPort[flit.OutputBuf]; ok {
			// 	continue
			// }
			// don't need to check if one can push into the output buffer
			// since a preexisting reservation will ensure that there is space
			// inside the buffer and a new reservation will also check if there is
			// a space inside the buffer
			// if flit.OutputBuf.CanPush() {
			// if !flit.OutputBuf.CanPush() {
			// 	// continue
			// 	panic("something is wrong")
			// }
			// fmt.Println("sending out flit:", i, flit, flit.OutputBuf)
			a.markOneFlitSent(flit)
			// a.outgoingReqsPerChiplet[chiplet]++

			// }
		} else {
			continue
		}
		// }
		// else {

		// }
		selectedPort = append(selectedPort, buf)
		a.outgoingReqsPerChiplet[chiplet]++
		occupiedOutputPort[flit.OutputBuf] = true
	}

	a.nextPortID = (a.nextPortID + 1) % len(a.buffers)
	// fmt.Println("********************************************************")
	return selectedPort
}
