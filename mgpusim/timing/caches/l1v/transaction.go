package l1v

import (
	"gitlab.com/akita/mem"
	"gitlab.com/akita/mem/cache"
	"gitlab.com/akita/util/ca"
)

type bankActionType int

const (
	bankActionInvalid bankActionType = iota
	bankActionReadHit
	bankActionWrite
	bankActionWriteFetched
)

type transaction struct {
	id string

	read                *mem.ReadReq
	readToBottom        *mem.ReadReq
	dataReadyFromBottom *mem.DataReadyRsp
	dataReadyToTop      *mem.DataReadyRsp

	write          *mem.WriteReq
	writeToBottom  *mem.WriteReq
	doneFromBottom *mem.WriteDoneRsp
	doneToTop      *mem.WriteDoneRsp

	preCoalesceTransactions []*transaction

	bankAction            bankActionType
	block                 *cache.Block
	data                  []byte
	writeFetchedDirtyMask []bool

	fetchAndWrite bool
	bankDone      bool
	bottomDone    bool
	done          bool
}

func (t *transaction) Address() uint64 {
	if t.read != nil {
		return t.read.Address
	}
	return t.write.Address
}

func (t *transaction) PID() ca.PID {
	if t.read != nil {
		return t.read.PID
	}
	return t.write.PID
}
