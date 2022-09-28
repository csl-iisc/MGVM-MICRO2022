// Package rob implemented an reorder buffer for memory requests.
package rob

import (
	"container/list"
	"fmt"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem"
	"gitlab.com/akita/util/tracing"
)

type transaction struct {
	reqFromTop    mem.AccessReq
	reqToBottom   mem.AccessReq
	rspFromBottom mem.AccessRsp
}

// ReorderBuffer can maintain the returning order of memory transactions.
type ReorderBuffer struct {
	*akita.TickingComponent

	TopPort     akita.Port
	BottomPort  akita.Port
	ControlPort akita.Port

	BottomUnit akita.Port

	bufferSize     int
	numReqPerCycle int

	toBottomReqIDToTransactionTable map[string]*list.Element
	transactions                    *list.List
	isFlushing                      bool
}

// Tick updates the status of the ReorderBuffer.
func (b *ReorderBuffer) Tick(now akita.VTimeInSec) (madeProgress bool) {
	madeProgress = b.processControlMsg(now) || madeProgress

	if !b.isFlushing {
		madeProgress = b.runPipeline(now) || madeProgress
	}

	return madeProgress
}

func (b *ReorderBuffer) processControlMsg(
	now akita.VTimeInSec,
) (madeProgress bool) {
	item := b.ControlPort.Peek()
	if item == nil {
		return false
	}

	msg := item.(*mem.ControlMsg)
	if msg.DiscardTransations {
		return b.discardTransactions(now, msg)
	} else if msg.Restart {
		return b.restart(now, msg)
	}

	panic("never")
}

func (b *ReorderBuffer) discardTransactions(
	now akita.VTimeInSec,
	msg *mem.ControlMsg,
) (madeProgress bool) {
	rsp := mem.ControlMsgBuilder{}.
		WithSrc(b.ControlPort).
		WithDst(msg.Src).
		WithSendTime(now).
		ToNotifyDone().
		Build()

	err := b.ControlPort.Send(rsp)
	if err != nil {
		return false
	}

	b.isFlushing = true
	b.toBottomReqIDToTransactionTable = make(map[string]*list.Element)
	b.transactions.Init()
	b.ControlPort.Retrieve(now)

	// fmt.Printf("%.10f, %s, rob flushed\n", now, b.Name())

	return true
}

func (b *ReorderBuffer) restart(
	now akita.VTimeInSec,
	msg *mem.ControlMsg,
) (madeProgress bool) {
	rsp := mem.ControlMsgBuilder{}.
		WithSrc(b.ControlPort).
		WithDst(msg.Src).
		WithSendTime(now).
		ToNotifyDone().
		Build()

	err := b.ControlPort.Send(rsp)
	if err != nil {
		return false
	}

	b.isFlushing = false
	b.toBottomReqIDToTransactionTable = make(map[string]*list.Element)
	b.transactions.Init()

	for b.TopPort.Retrieve(now) != nil {
	}

	for b.BottomPort.Retrieve(now) != nil {
	}

	b.ControlPort.Retrieve(now)

	// fmt.Printf("%.10f, %s, rob restarted\n", now, b.Name())

	return true
}

func (b *ReorderBuffer) runPipeline(now akita.VTimeInSec) (madeProgress bool) {
	for i := 0; i < b.numReqPerCycle; i++ {
		madeProgress = b.bottomUp(now) || madeProgress
	}

	for i := 0; i < b.numReqPerCycle; i++ {
		madeProgress = b.parseBottom(now) || madeProgress
	}

	for i := 0; i < b.numReqPerCycle; i++ {
		madeProgress = b.topDown(now) || madeProgress
	}

	return madeProgress
}

func (b *ReorderBuffer) topDown(now akita.VTimeInSec) bool {
	if b.isFull() {
		return false
	}

	item := b.TopPort.Peek()
	if item == nil {
		return false
	}

	req := item.(mem.AccessReq)
	trans := b.createTransaction(req)

	trans.reqToBottom.Meta().Src = b.BottomPort
	trans.reqToBottom.Meta().SendTime = now
	err := b.BottomPort.Send(trans.reqToBottom)
	if err != nil {
		return false
	}

	b.addTransaction(trans)
	b.TopPort.Retrieve(now)

	tracing.TraceReqReceive(req, now, b)
	tracing.TraceReqInitiate(trans.reqToBottom, now, b,
		tracing.MsgIDAtReceiver(req, b))

	return true
}

func (b *ReorderBuffer) parseBottom(now akita.VTimeInSec) bool {
	item := b.BottomPort.Peek()
	if item == nil {
		return false
	}

	rsp := item.(mem.AccessRsp)
	rspTo := rsp.GetRespondTo()
	if rspTo == "5818306" {
		fmt.Printf("Here\n")
	}
	transElement, found := b.toBottomReqIDToTransactionTable[rspTo]

	if found {
		trans := transElement.Value.(*transaction)
		trans.rspFromBottom = rsp

		tracing.TraceReqFinalize(trans.reqToBottom, now, b)
	}

	b.BottomPort.Retrieve(now)

	return true
}

func (b *ReorderBuffer) bottomUp(now akita.VTimeInSec) bool {
	elem := b.transactions.Front()
	if elem == nil {
		return false
	}

	trans := elem.Value.(*transaction)
	if trans.rspFromBottom == nil {
		return false
	}

	rsp := b.duplicateRsp(trans.rspFromBottom, trans.reqFromTop.Meta().ID)
	rsp.Meta().Dst = trans.reqFromTop.Meta().Src
	rsp.Meta().Src = b.TopPort
	rsp.Meta().SendTime = now

	err := b.TopPort.Send(rsp)
	if err != nil {
		return false
	}

	b.deleteTransaction(elem)

	tracing.TraceReqComplete(trans.reqFromTop, now, b)

	return true
}

func (b *ReorderBuffer) isFull() bool {
	return b.transactions.Len() >= b.bufferSize
}

func (b *ReorderBuffer) createTransaction(req mem.AccessReq) *transaction {
	return &transaction{
		reqFromTop:  req,
		reqToBottom: b.duplicateReq(req),
	}
}

func (b *ReorderBuffer) addTransaction(trans *transaction) {
	elem := b.transactions.PushBack(trans)
	b.toBottomReqIDToTransactionTable[trans.reqToBottom.Meta().ID] = elem
}

func (b *ReorderBuffer) deleteTransaction(elem *list.Element) {
	trans := elem.Value.(*transaction)
	b.transactions.Remove(elem)
	delete(b.toBottomReqIDToTransactionTable, trans.reqToBottom.Meta().ID)
}

func (b *ReorderBuffer) duplicateReq(req mem.AccessReq) mem.AccessReq {
	switch req := req.(type) {
	case *mem.ReadReq:
		return b.duplicateReadReq(req)
	case *mem.WriteReq:
		return b.duplicateWriteReq(req)
	default:
		panic("unsupported type")
	}
}

func (b *ReorderBuffer) duplicateReadReq(req *mem.ReadReq) *mem.ReadReq {
	return mem.ReadReqBuilder{}.
		WithAddress(req.Address).
		WithByteSize(req.AccessByteSize).
		WithPID(req.PID).
		WithDst(b.BottomUnit).
		Build()
}

func (b *ReorderBuffer) duplicateWriteReq(req *mem.WriteReq) *mem.WriteReq {
	return mem.WriteReqBuilder{}.
		WithAddress(req.Address).
		WithPID(req.PID).
		WithData(req.Data).
		WithDirtyMask(req.DirtyMask).
		WithDst(b.BottomUnit).
		Build()
}

func (b *ReorderBuffer) duplicateRsp(
	rsp mem.AccessRsp,
	rspTo string,
) mem.AccessRsp {
	switch rsp := rsp.(type) {
	case *mem.DataReadyRsp:
		return b.duplicateDataReadyRsp(rsp, rspTo)
	case *mem.WriteDoneRsp:
		return b.duplicateWriteDoneRsp(rsp, rspTo)
	default:
		panic("type not supported")
	}
}

func (b *ReorderBuffer) duplicateDataReadyRsp(
	rsp *mem.DataReadyRsp,
	rspTo string,
) *mem.DataReadyRsp {
	return mem.DataReadyRspBuilder{}.
		WithData(rsp.Data).
		WithRspTo(rspTo).
		Build()
}

func (b *ReorderBuffer) duplicateWriteDoneRsp(
	rsp *mem.WriteDoneRsp,
	rspTo string,
) *mem.WriteDoneRsp {
	return mem.WriteDoneRspBuilder{}.
		WithRspTo(rspTo).
		Build()
}
