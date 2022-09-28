package writeback

import (
	"fmt"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem"
	"gitlab.com/akita/mem/cache"
	"gitlab.com/akita/util/ca"
	"gitlab.com/akita/util/tracing"
)

type directoryStage struct {
	cache *Cache
}

func (ds *directoryStage) Tick(now akita.VTimeInSec) bool {
	// peek at pipeline
	pipelineitem := ds.cache.lookupBuffer.Peek()
	// item := ds.cache.dirStageBuffer.Peek()
	if pipelineitem == nil {
		return false
	}
	pitem := pipelineitem.(cachePipelineItem)
	item := pitem.trans

	// trans := item.(*transaction)
	trans := item
	if trans.read != nil {
		/*return*/ ds.doRead(now, trans)
		return true
	}

	/*return*/
	ds.doWrite(now, trans)
	return true
}

func (ds *directoryStage) Reset(now akita.VTimeInSec) {
	ds.cache.dirStageBuffer.Clear()
}

func (ds *directoryStage) doRead(
	now akita.VTimeInSec,
	trans *transaction,
) bool {
	cachelineID, _ := getCacheLineID(
		trans.read.Address, ds.cache.log2BlockSize)

	mshrEntry := ds.cache.mshr.Query(trans.read.PID, cachelineID)
	if mshrEntry != nil {
		return ds.handleReadMSHRHit(now, trans, mshrEntry)
	}

	block := ds.cache.directory.Lookup(
		trans.read.PID, cachelineID)
	if block != nil {
		return ds.handleReadHit(now, trans, block)

	}

	return ds.handleReadMiss(now, trans)
}

func (ds *directoryStage) handleReadMSHRHit(
	now akita.VTimeInSec,
	trans *transaction,
	mshrEntry *cache.MSHREntry,
) bool {
	trans.mshrEntry = mshrEntry
	mshrEntry.Requests = append(mshrEntry.Requests, trans)
	// pipeline
	ds.cache.lookupBuffer.Pop()
	// ds.cache.dirStageBuffer.Pop()

	tracing.AddTaskStep(
		tracing.MsgIDAtReceiver(trans.read, ds.cache),
		now, ds.cache,
		"read-mshr-hit",
	)
	if trans.read.Info != nil {
		readReqInfo := trans.read.Info.(*mem.ReadReqInfo)
		if readReqInfo.ReturnAccessInfo {
			readReqInfo.AccessResult = mem.ReadMSHRHit
		}
	}

	return true
}

func (ds *directoryStage) handleReadHit(
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

	if trans.read.Info != nil {
		readReqInfo := trans.read.Info.(*mem.ReadReqInfo)
		if readReqInfo.ReturnAccessInfo {
			readReqInfo.AccessResult = mem.ReadHit
		}
	}

	return ds.readFromBank(trans, block)
}

func (ds *directoryStage) handleReadMiss(
	now akita.VTimeInSec,
	trans *transaction,
) bool {
	req := trans.read
	cacheLineID, _ := getCacheLineID(req.Address, ds.cache.log2BlockSize)

	if ds.cache.mshr.IsFull() {
		return false
	}

	victim := ds.cache.directory.FindVictim(cacheLineID)
	if victim.IsLocked || victim.ReadCount > 0 {
		return false
	}

	if ds.needEviction(victim) {
		ok := ds.evict(now, trans, victim)
		if ok {
			tracing.AddTaskStep(
				tracing.MsgIDAtReceiver(trans.read, ds.cache),
				now, ds.cache,
				"read-miss",
			)
		}
		if trans.read.Info != nil {
			readReqInfo := trans.read.Info.(*mem.ReadReqInfo)
			if readReqInfo.ReturnAccessInfo {
				readReqInfo.AccessResult = mem.ReadMiss
			}
		}
		return ok
	}

	ok := ds.fetch(now, trans, victim)
	if ok {
		tracing.AddTaskStep(
			tracing.MsgIDAtReceiver(trans.read, ds.cache),
			now, ds.cache,
			"read-miss",
		)
		if trans.read.Info != nil {
			readReqInfo := trans.read.Info.(*mem.ReadReqInfo)
			if readReqInfo.ReturnAccessInfo {
				readReqInfo.AccessResult = mem.ReadMiss
			}
		}
	}
	return ok
}

func (ds *directoryStage) doWrite(
	now akita.VTimeInSec,
	trans *transaction,
) bool {
	write := trans.write
	cachelineID, _ := getCacheLineID(write.Address, ds.cache.log2BlockSize)

	mshrEntry := ds.cache.mshr.Query(write.PID, cachelineID)
	if mshrEntry != nil {
		ok := ds.doWriteMSHRHit(now, trans, mshrEntry)
		tracing.AddTaskStep(
			tracing.MsgIDAtReceiver(trans.write, ds.cache),
			now, ds.cache,
			"write-mshr-hit",
		)

		return ok
	}

	block := ds.cache.directory.Lookup(trans.write.PID, cachelineID)
	if block != nil {
		ok := ds.doWriteHit(trans, block)
		if ok {
			tracing.AddTaskStep(
				tracing.MsgIDAtReceiver(trans.write, ds.cache),
				now, ds.cache,
				"write-hit",
			)
		}

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

func (ds *directoryStage) doWriteMSHRHit(
	now akita.VTimeInSec,
	trans *transaction,
	mshrEntry *cache.MSHREntry,
) bool {
	trans.mshrEntry = mshrEntry
	mshrEntry.Requests = append(mshrEntry.Requests, trans)
	//pipeline
	ds.cache.lookupBuffer.Pop()
	// ds.cache.dirStageBuffer.Pop()

	return true
}

func (ds *directoryStage) doWriteHit(
	trans *transaction,
	block *cache.Block,
) bool {
	if block.IsLocked || block.ReadCount > 0 {
		return false
	}

	return ds.writeToBank(trans, block)
}

func (ds *directoryStage) doWriteMiss(
	now akita.VTimeInSec,
	trans *transaction,
) bool {
	write := trans.write

	if ds.isWritingFullLine(write) {
		return ds.writeFullLineMiss(now, trans)
	}
	return ds.writePartialLineMiss(now, trans)
}

func (ds *directoryStage) writeFullLineMiss(now akita.VTimeInSec, trans *transaction) bool {
	write := trans.write
	cachelineID, _ := getCacheLineID(write.Address, ds.cache.log2BlockSize)

	victim := ds.cache.directory.FindVictim(cachelineID)
	if victim.IsLocked || victim.ReadCount > 0 {
		return false
	}

	if ds.needEviction(victim) {
		return ds.evict(now, trans, victim)
	}

	return ds.writeToBank(trans, victim)
}

func (ds *directoryStage) writePartialLineMiss(
	now akita.VTimeInSec,
	trans *transaction,
) bool {
	write := trans.write
	cachelineID, _ := getCacheLineID(write.Address, ds.cache.log2BlockSize)

	if ds.cache.mshr.IsFull() {
		return false
	}

	victim := ds.cache.directory.FindVictim(cachelineID)
	if victim.IsLocked || victim.ReadCount > 0 {
		return false
	}

	if ds.needEviction(victim) {
		return ds.evict(now, trans, victim)
	}

	return ds.fetch(now, trans, victim)
}

func (ds *directoryStage) readFromBank(
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
	// pipeline
	ds.cache.lookupBuffer.Pop()
	// ds.cache.dirStageBuffer.Pop()
	bankBuf.Push(trans)
	return true
}

func (ds *directoryStage) writeToBank(
	trans *transaction,
	block *cache.Block,
) bool {
	numBanks := len(ds.cache.dirToBankBuffers)
	bank := bankID(block, ds.cache.directory.WayAssociativity(), numBanks)
	bankBuf := ds.cache.dirToBankBuffers[bank]

	if !bankBuf.CanPush() {
		return false
	}

	addr := trans.write.Address
	cachelineID, _ := getCacheLineID(addr, ds.cache.log2BlockSize)

	ds.cache.directory.Visit(block)
	block.IsLocked = true
	block.Tag = cachelineID
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

func (ds *directoryStage) evict(
	now akita.VTimeInSec,
	trans *transaction,
	victim *cache.Block,
) bool {
	bankNum := bankID(victim,
		ds.cache.directory.WayAssociativity(), len(ds.cache.dirToBankBuffers))
	bankBuf := ds.cache.dirToBankBuffers[bankNum]

	if !bankBuf.CanPush() {
		return false
	}

	var addr uint64
	var pid ca.PID
	if trans.read != nil {
		addr = trans.read.Address
		pid = trans.read.PID
	} else {
		addr = trans.write.Address
		pid = trans.write.PID
	}

	cacheLineID, _ := getCacheLineID(addr, ds.cache.log2BlockSize)

	trans.action = bankEvictAndFetch
	trans.victim = &cache.Block{
		PID:          victim.PID,
		Tag:          victim.Tag,
		CacheAddress: victim.CacheAddress,
		DirtyMask:    victim.DirtyMask,
	}
	trans.block = victim
	trans.evictingPID = trans.victim.PID
	trans.evictingAddr = trans.victim.Tag
	trans.evictingDirtyMask = victim.DirtyMask

	if ds.evictionNeedFetch(trans) {
		mshrEntry := ds.cache.mshr.Add(pid, cacheLineID)
		mshrEntry.Block = victim
		mshrEntry.Requests = append(mshrEntry.Requests, trans)
		trans.mshrEntry = mshrEntry
		trans.fetchPID = pid
		trans.fetchAddress = cacheLineID
		trans.action = bankEvictAndFetch
	} else {
		trans.action = bankEvictAndWrite
	}

	victim.Tag = cacheLineID
	victim.PID = pid
	victim.IsLocked = true
	victim.IsDirty = false

	//pipeline
	ds.cache.lookupBuffer.Pop()
	// ds.cache.dirStageBuffer.Pop()
	bankBuf.Push(trans)

	return true
}

func (ds *directoryStage) evictionNeedFetch(t *transaction) bool {
	if t.write == nil {
		return true
	}

	if ds.isWritingFullLine(t.write) {
		return false
	}

	return true
}

func (ds *directoryStage) fetch(
	now akita.VTimeInSec,
	trans *transaction,
	block *cache.Block,
) bool {
	var addr uint64
	var pid ca.PID
	var req mem.AccessReq
	if trans.read != nil {
		req = trans.read
		addr = trans.read.Address
		pid = trans.read.PID
	} else {
		req = trans.write
		addr = trans.write.Address
		pid = trans.write.PID
	}
	cacheLineID, _ := getCacheLineID(addr, ds.cache.log2BlockSize)

	bankNum := bankID(block,
		ds.cache.directory.WayAssociativity(), len(ds.cache.dirToBankBuffers))
	bankBuf := ds.cache.dirToBankBuffers[bankNum]

	if !bankBuf.CanPush() {
		return false
	}

	mshrEntry := ds.cache.mshr.Add(pid, cacheLineID)
	trans.mshrEntry = mshrEntry
	trans.block = block
	block.IsLocked = true
	block.Tag = cacheLineID
	block.PID = pid
	block.IsValid = true
	ds.cache.directory.Visit(block)

	tracing.AddTaskStep(
		tracing.MsgIDAtReceiver(req, ds.cache),
		now, ds.cache,
		fmt.Sprintf("add-mshr-entry-0x%x-0x%x", mshrEntry.Address, block.Tag),
	)

	//pipeline
	ds.cache.lookupBuffer.Pop()
	// ds.cache.dirStageBuffer.Pop()

	trans.action = writeBufferFetch
	trans.fetchPID = pid
	trans.fetchAddress = cacheLineID
	bankBuf.Push(trans)

	mshrEntry.Block = block
	mshrEntry.Requests = append(mshrEntry.Requests, trans)

	return true
}

func (ds *directoryStage) isWritingFullLine(write *mem.WriteReq) bool {
	if len(write.Data) != (1 << ds.cache.log2BlockSize) {
		return false
	}

	if write.DirtyMask != nil {
		for _, dirty := range write.DirtyMask {
			if !dirty {
				return false
			}
		}
	}

	return true
}

func (ds *directoryStage) needEviction(victim *cache.Block) bool {
	return victim.IsValid && victim.IsDirty
}
