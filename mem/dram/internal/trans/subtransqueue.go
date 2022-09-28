package trans

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem/dram/internal/signal"
)

// A SubTransactionQueue is a queue for subtransactions.
type SubTransactionQueue interface {
	CanPush(n int) bool
	Push(t *signal.Transaction)
	Tick(now akita.VTimeInSec) bool
}
