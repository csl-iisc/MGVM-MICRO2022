package akitaext

import (
	"log"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/util"
)

// BufferedSender can delegate the sending process.
//
// The most common usage of BufferedSender is to be used as the send stage of
// an Akita Component. It is common that multiple sub-component in a component
// need to send messages out from a port. Another common pattern is that a
// large number of messages are generated in one cycle and the messages are
// sent out one per cycle. In both cases, the msguest generate can push the
// packet into a BufferedSender and the call the BufferedSender's Tick function
// to actually send the msguests out.
type BufferedSender interface {
	// CanSend checks if the buffer has enough space to hold "count" msguests.
	CanSend(count int) bool

	// Send enqueues a message into the buffer and the message will be sent out
	// later with the Tick function.
	Send(msg akita.Msg)

	// Clear removes all the msguests to send
	Clear()

	// Tick tries to send one msguest out. If successful, Tick returns true.
	Tick(now akita.VTimeInSec) bool
}

// NewBufferedSender creates a new BufferedSender with certain buffer capacity
// and send to a certain port.
func NewBufferedSender(port akita.Port, buffer util.Buffer) BufferedSender {
	return &bufferedSenderImpl{
		port:   port,
		buffer: buffer,
	}
}

type bufferedSenderImpl struct {
	port   akita.Port
	buffer util.Buffer
}

func (s *bufferedSenderImpl) CanSend(count int) bool {
	if count > s.buffer.Capacity() {
		log.Panic("trying to send number of reuqests exceeding capacity")
	}

	if count+s.buffer.Size() > s.buffer.Capacity() {
		return false
	}

	return true
}

func (s *bufferedSenderImpl) Send(msg akita.Msg) {
	s.buffer.Push(msg)
}

func (s *bufferedSenderImpl) Clear() {
	s.buffer.Clear()
}

func (s *bufferedSenderImpl) Tick(now akita.VTimeInSec) bool {
	item := s.buffer.Peek()
	if item == nil {
		return false
	}

	msg := item.(akita.Msg)
	msg.Meta().SendTime = now
	err := s.port.Send(msg)
	if err != nil {
		return false
	}

	s.buffer.Pop()

	return true
}
