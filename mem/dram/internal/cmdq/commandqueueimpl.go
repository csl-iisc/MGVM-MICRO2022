package cmdq

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem/dram/internal/org"
	"gitlab.com/akita/mem/dram/internal/signal"
)

type Queue []*signal.Command

type CommandQueueImpl struct {
	Queues           []Queue
	CapacityPerQueue int
	nextQueueIndex   int
	Channel          org.Channel
}

func (q *CommandQueueImpl) GetCommandToIssue(
	now akita.VTimeInSec,
) *signal.Command {
	for i := 0; i < len(q.Queues); i++ {
		queueIndex, _ := q.getNextQueue()
		readyCmd := q.getFirstReadyInQueue(now, queueIndex)

		if readyCmd != nil {
			return readyCmd
		}
	}

	return nil
}

func (q *CommandQueueImpl) getNextQueue() (queueIndex int, queue Queue) {
	queueIndex = q.nextQueueIndex
	retQueue := q.Queues[q.nextQueueIndex]
	q.nextQueueIndex = (q.nextQueueIndex + 1) % len(q.Queues)
	return queueIndex, retQueue
}

func (q *CommandQueueImpl) getFirstReadyInQueue(
	now akita.VTimeInSec,
	queueIndex int,
) *signal.Command {
	for i, cmd := range q.Queues[queueIndex] {
		readyCmd := q.Channel.GetReadyCommand(now, cmd)

		if readyCmd != nil {
			if cmd.Kind == readyCmd.Kind {
				q.Queues[queueIndex] = append(
					q.Queues[queueIndex][:i], q.Queues[queueIndex][i+1:]...)
			}
			return readyCmd
		}
	}

	return nil
}

func (q *CommandQueueImpl) CanAccept(cmd *signal.Command) bool {
	queueIndex := q.getQueueIndex(cmd)
	queue := q.Queues[queueIndex]

	if len(queue) < q.CapacityPerQueue {
		return true
	}

	return false
}

func (q *CommandQueueImpl) Accept(cmd *signal.Command) {
	queueIndex := q.getQueueIndex(cmd)
	queue := q.Queues[queueIndex]

	if len(queue) >= q.CapacityPerQueue {
		panic("command queue overflow")
	}

	q.Queues[queueIndex] = append(queue, cmd)
}

func (q *CommandQueueImpl) getQueueIndex(cmd *signal.Command) int {
	return int(cmd.Rank)
}
