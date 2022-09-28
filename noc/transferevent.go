package noc

import "gitlab.com/akita/akita"

// A TransferEvent is an event that marks that a message completes transfer.
type TransferEvent struct {
	*akita.EventBase
	msg akita.Msg
	vc  int
}

// NewTransferEvent creates a new TransferEvent.
func NewTransferEvent(
	time akita.VTimeInSec,
	handler akita.Handler,
	msg akita.Msg,
	vc int,
) *TransferEvent {
	evt := new(TransferEvent)
	evt.EventBase = akita.NewEventBase(time, handler)
	evt.msg = msg
	evt.vc = vc
	return evt
}
