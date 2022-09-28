package trans

import (
	"gitlab.com/akita/mem/dram/internal/signal"
)

// A CommandCreator can convert a subtransaction to a command.
type CommandCreator interface {
	Create(subTrans *signal.SubTransaction) *signal.Command
}
