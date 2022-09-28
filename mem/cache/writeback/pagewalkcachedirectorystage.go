package writeback

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem"
	"gitlab.com/akita/mem/cache"
	"gitlab.com/akita/util/tracing"
)

type pageWalkCacheDirectoryStage struct {
	cache *PageWalkCache
}

func (ds *pageWalkCacheDirectoryStage) Tick(now akita.VTimeInSec) bool {
	// peek at pipeline
	pipelineitem := ds.cache.lookupBuffer.Peek()
	// item := ds.cache.dirStageBuffer.Peek()
	if pipelineitem == nil {
		return false
	}
	pitem := pipelineitem.(pwcPipelineItem)
	item := pitem.trans

	// trans := item.(*transaction)
	trans := item
	if trans.read != nil {
		return ds.doRead(now, trans)
	}
	return ds.doWrite(now, trans)
}

func (ds *pageWalkCacheDirectoryStage) Reset(now akita.VTimeInSec) {
	ds.cache.dirStageBuffer.Clear()
}

func (ds *pageWalkCacheDirectoryStage) doRead(
	now akita.VTimeInSec,
	trans *transaction,
) bool {
	block := ds.cache.directory.Lookup(
		trans.read.PID, trans.read.Address)
	if block != nil {
		return ds.handleReadHit(now, trans, block)
	}
	return ds.handleReadMiss(now, trans)
}

func (ds *pageWalkCacheDirectoryStage) handleReadHit(
	now akita.VTimeInSec,
	trans *transaction,
	block *cache.Block,
) bool {
	if block.IsLocked {
		return false
	}

	tracing.AddTaskStep(
		tracing.MsgIDAtReceiver(trans.read, ds.cache),
		now, ds.cache,
		"read-hit",
	)

	return ds.readFromBank(trans, block)
}

func (ds *pageWalkCacheDirectoryStage) handleReadMiss(
	now akita.VTimeInSec,
	trans *transaction,
) bool {
	tracing.AddTaskStep(
		tracing.MsgIDAtReceiver(trans.read, ds.cache),
		now, ds.cache,
		"read-miss",
	)
	dataReady := mem.DataReadyRspBuilder{}.
		WithSendTime(now).
		WithSrc(ds.cache.TopPort).
		WithDst(trans.read.Src).
		WithRspTo(trans.read.ID).
		WithData(nil).
		Build()
	ds.cache.topSender.Send(dataReady)

	// pipeline
	ds.cache.lookupBuffer.Pop()
	// ds.cache.dirStageBuffer.Pop()

	return true
}

func (ds *pageWalkCacheDirectoryStage) doWrite(
	now akita.VTimeInSec,
	trans *transaction,
) bool {
	write := trans.write
	write.Address = ds.cache.directory.FormulateWriteAddress(write.Address)
	block := ds.cache.directory.LookupExactAddress(trans.write.PID, write.Address)
	if block != nil {
		ok := ds.doWriteHit(now, trans, block)
		tracing.AddTaskStep(
			tracing.MsgIDAtReceiver(trans.write, ds.cache),
			now, ds.cache,
			"write-hit",
		)
		return ok
	}

	ok := ds.doWriteMiss(now, trans)
	if ok {
		tracing.AddTaskStep(
			tracing.MsgIDAtReceiver(trans.write, ds.cache),
			now, ds.cache,
			"write-miss",
		)
	}

	return ok
}

func (ds *pageWalkCacheDirectoryStage) doWriteHit(
	now akita.VTimeInSec,
	trans *transaction,
	block *cache.Block,
) bool {
	if block.IsLocked {
		return false
	}
	ds.cache.directory.Visit(block)
	done := mem.WriteDoneRspBuilder{}.
		WithSendTime(now).
		WithSrc(ds.cache.TopPort).
		WithDst(trans.write.Src).
		WithRspTo(trans.write.ID).
		Build()
	ds.cache.topSender.Send(done)
	//pipeline
	ds.cache.lookupBuffer.Pop()
	// ds.cache.dirStageBuffer.Pop()
	return true
}

func (ds *pageWalkCacheDirectoryStage) doWriteMiss(
	now akita.VTimeInSec,
	trans *transaction,
) bool {
	victim := ds.cache.directory.FindVictim(trans.write.Address)
	if victim.IsLocked || victim.ReadCount > 0 {
		return false
	}
	return ds.writeToBank(trans, victim)
}

func (ds *pageWalkCacheDirectoryStage) readFromBank(
	trans *transaction,
	block *cache.Block,
) bool {
	numBanks := len(ds.cache.dirToBankBuffers)
	bank := bankID(block, ds.cache.directory.WayAssociativity(), numBanks)
	bankBuf := ds.cache.dirToBankBuffers[bank]

	if !bankBuf.CanPush() {
		return false
	}

	ds.cache.directory.Visit(block)
	block.ReadCount++
	trans.block = block
	trans.action = bankReadHit
	//pipeline
	ds.cache.lookupBuffer.Pop()
	// ds.cache.dirStageBuffer.Pop()
	bankBuf.Push(trans)
	return true
}

func (ds *pageWalkCacheDirectoryStage) writeToBank(
	trans *transaction,
	block *cache.Block,
) bool {
	numBanks := len(ds.cache.dirToBankBuffers)
	bank := bankID(block, ds.cache.directory.WayAssociativity(), numBanks)
	bankBuf := ds.cache.dirToBankBuffers[bank]

	if !bankBuf.CanPush() {
		return false
	}
	// h := fmt.Sprintf("%x", trans.write.Address)
	// fmt.Println("dir:", h)
	ds.cache.directory.Visit(block)
	block.IsLocked = true
	block.Tag = trans.write.Address
	block.IsValid = true
	block.PID = trans.write.PID
	trans.block = block
	trans.action = bankWriteHit
	//pipeline
	ds.cache.lookupBuffer.Pop()
	// ds.cache.dirStageBuffer.Pop()
	bankBuf.Push(trans)

	return true
}
