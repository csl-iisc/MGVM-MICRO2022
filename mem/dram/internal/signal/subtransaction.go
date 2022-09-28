package signal

// SubTransaction is the read and write to a single bank.
type SubTransaction struct {
	ID          string
	Transaction *Transaction
	Address     uint64
	Completed   bool
}

func (st SubTransaction) IsRead() bool {
	return st.Transaction.IsRead()
}
