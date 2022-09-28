package dram

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem"
	"gitlab.com/akita/mem/dram/internal/addressmapping"
	"gitlab.com/akita/mem/dram/internal/cmdq"
	"gitlab.com/akita/mem/dram/internal/org"
	"gitlab.com/akita/mem/dram/internal/signal"
	"gitlab.com/akita/mem/dram/internal/trans"
	"gitlab.com/akita/util/tracing"
)

// Protocol defines the category of the memory controller.
type Protocol int

const (
	DDR3 Protocol = iota
	DDR4
	GDDR5
	GDDR5X
	GDDR6
	LPDDR
	LPDDR3
	LPDDR4
	HBM
	HBM2
	HMC
)

func (p Protocol) isGDDR() bool {
	return p == GDDR5 || p == GDDR5X || p == GDDR6
}

func (p Protocol) isHBM() bool {
	return p == HBM || p == HBM2
}

// A MemController handles read and write requests.
type MemController struct {
	*akita.TickingComponent

	TopPort akita.Port

	storage             *mem.Storage
	addrConverter       addressmapping.AddressConverter
	subTransSplitter    trans.SubTransSplitter
	addrMapper          addressmapping.Mapper
	subTransactionQueue trans.SubTransactionQueue
	cmdQueue            cmdq.CommandQueue
	channel             org.Channel

	inflightTransactions []*signal.Transaction
}

// Tick updates memory controller's internal state.
func (c *MemController) Tick(now akita.VTimeInSec) (madeProgress bool) {
	madeProgress = c.respond(now) || madeProgress
	madeProgress = c.channel.Tick(now) || madeProgress
	madeProgress = c.issue(now) || madeProgress
	madeProgress = c.subTransactionQueue.Tick(now) || madeProgress
	madeProgress = c.parseTop(now) || madeProgress
	return madeProgress
}

func (c *MemController) parseTop(now akita.VTimeInSec) (madeProgress bool) {
	msg := c.TopPort.Peek()
	if msg == nil {
		return false
	}

	trans := &signal.Transaction{}
	switch msg := msg.(type) {
	case *mem.ReadReq:
		trans.Read = msg
	case *mem.WriteReq:
		trans.Write = msg
	}

	c.assignTransInternalAddress(trans)
	c.subTransSplitter.Split(trans)

	if !c.subTransactionQueue.CanPush(len(trans.SubTransactions)) {
		return false
	}

	c.subTransactionQueue.Push(trans)
	c.inflightTransactions = append(c.inflightTransactions, trans)
	c.TopPort.Retrieve(now)

	tracing.TraceReqReceive(msg, now, c)
	for _, st := range trans.SubTransactions {
		tracing.StartTask(
			st.ID,
			tracing.MsgIDAtReceiver(msg, c),
			now,
			c,
			"sub-trans",
			"sub-trans",
			nil,
		)
	}

	// fmt.Printf("%.10f, %s, start transaction, %s, %x\n",
	// 	now, c.Name(), msg.Meta().ID, trans.InternalAddress)

	return true
}

func (c *MemController) assignTransInternalAddress(trans *signal.Transaction) {
	if c.addrConverter != nil {
		trans.InternalAddress = c.addrConverter.ConvertExternalToInternal(
			trans.GlobalAddress())
		return
	}

	trans.InternalAddress = trans.GlobalAddress()
}

func (c *MemController) issue(now akita.VTimeInSec) (madeProgress bool) {
	cmd := c.cmdQueue.GetCommandToIssue(now)
	if cmd == nil {
		return false
	}

	c.channel.StartCommand(now, cmd)
	c.channel.UpdateTiming(now, cmd)

	return true
}

func (c *MemController) respond(now akita.VTimeInSec) (madeProgress bool) {
	for i, t := range c.inflightTransactions {
		if t.IsCompleted() {
			done := c.finalizeTransaction(now, t, i)
			if done {
				return true
			}
		}
	}

	return false
}

func (c *MemController) finalizeTransaction(
	now akita.VTimeInSec,
	t *signal.Transaction,
	i int,
) (done bool) {
	if t.Write != nil {
		done = c.finalizeWriteTrans(now, t, i)
		if done {
			tracing.TraceReqComplete(t.Write, now, c)
		}
	} else {
		done = c.finalizeReadTrans(now, t, i)
		if done {
			tracing.TraceReqComplete(t.Read, now, c)
		}
	}

	return done
}

func (c *MemController) finalizeWriteTrans(
	now akita.VTimeInSec,
	t *signal.Transaction,
	i int,
) (done bool) {
	err := c.storage.Write(t.InternalAddress, t.Write.Data)
	if err != nil {
		panic(err)
	}

	writeDone := mem.WriteDoneRspBuilder{}.
		WithSrc(c.TopPort).
		WithDst(t.Write.Src).
		WithRspTo(t.Write.ID).
		WithSendTime(now).
		Build()
	sendErr := c.TopPort.Send(writeDone)
	if sendErr == nil {
		c.inflightTransactions = append(
			c.inflightTransactions[:i],
			c.inflightTransactions[i+1:]...)

		// fmt.Printf("%.10f, %s, finish transaction %s, %x\n",
		// 	now, c.Name(), t.Write.ID, t.InternalAddress)
		return true
	}

	return false
}

func (c *MemController) finalizeReadTrans(
	now akita.VTimeInSec,
	t *signal.Transaction,
	i int,
) (done bool) {
	data, err := c.storage.Read(t.InternalAddress, t.Read.AccessByteSize)
	if err != nil {
		panic(err)
	}

	dataReady := mem.DataReadyRspBuilder{}.
		WithSrc(c.TopPort).
		WithDst(t.Read.Src).
		WithData(data).
		WithRspTo(t.Read.ID).
		WithSendTime(now).
		Build()
	sendErr := c.TopPort.Send(dataReady)
	if sendErr == nil {
		c.inflightTransactions = append(
			c.inflightTransactions[:i],
			c.inflightTransactions[i+1:]...)

		// fmt.Printf("%.10f, %s, finish transaction %s, %x\n",
		// 	now, c.Name(), t.Read.ID, t.InternalAddress)
		return true
	}

	return false
}
