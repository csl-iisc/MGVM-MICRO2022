package signal

import "gitlab.com/akita/mem"

// Transaction is the state associated with the processing of a read or write
// request.
type Transaction struct {
	Read  *mem.ReadReq
	Write *mem.WriteReq

	InternalAddress uint64
	SubTransactions []*SubTransaction
}

func (t *Transaction) GlobalAddress() uint64 {
	if t.Read != nil {
		return t.Read.Address
	}

	return t.Write.Address
}

func (t *Transaction) AccessByteSize() uint64 {
	if t.Read != nil {
		return t.Read.AccessByteSize
	}

	return uint64(len(t.Write.Data))
}

func (t *Transaction) IsRead() bool {
	return t.Read != nil
}

func (t *Transaction) IsWrite() bool {
	return t.Write != nil
}

func (t *Transaction) IsCompleted() bool {
	for _, st := range t.SubTransactions {
		if !st.Completed {
			return false
		}
	}

	return true
}
