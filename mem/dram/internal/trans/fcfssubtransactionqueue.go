package trans

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem/dram/internal/cmdq"
	"gitlab.com/akita/mem/dram/internal/signal"
)

// A FCFSSubTransactionQueue returns sub-transactions in a
// first-come-first-serve way.
type FCFSSubTransactionQueue struct {
	Capacity   int
	Queue      []*signal.SubTransaction
	CmdCreator CommandCreator
	CmdQueue   cmdq.CommandQueue
}

func (q *FCFSSubTransactionQueue) CanPush(n int) bool {
	if n >= q.Capacity {
		panic("queue size not large enough to handle a single transaction")
	}

	if len(q.Queue)+n > q.Capacity {
		return false
	}
	return true
}

func (q *FCFSSubTransactionQueue) Push(t *signal.Transaction) {
	if len(q.Queue)+len(t.SubTransactions) > q.Capacity {
		panic("pushing too many subtransactions into queue.")
	}

	q.Queue = append(q.Queue, t.SubTransactions...)
}

func (q *FCFSSubTransactionQueue) Tick(now akita.VTimeInSec) bool {
	for i, subTrans := range q.Queue {
		cmd := q.CmdCreator.Create(subTrans)

		if q.CmdQueue.CanAccept(cmd) {
			q.CmdQueue.Accept(cmd)
			q.Queue = append(q.Queue[:i], q.Queue[i+1:]...)

			// fmt.Printf("Command Pushed: %#v\n", cmd)

			return true
		}
	}

	return false
}
