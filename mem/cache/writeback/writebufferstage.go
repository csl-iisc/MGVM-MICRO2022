package writeback

import (
	// "fmt"
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem"
	"gitlab.com/akita/mem/cache"
	"gitlab.com/akita/util/tracing"
	"strings"
)

type writeBufferStage struct {
	cache *Cache

	writeBufferCapacity int
	maxInflightFetch    int
	maxInflightEviction int

	rdmaPort akita.Port

	pendingEvictions []*transaction
	inflightFetch    []*transaction
	inflightEviction []*transaction

	remoteMemAccesses uint64
}

func (wb *writeBufferStage) Tick(now akita.VTimeInSec) bool {
	madeProgress := false

	madeProgress = wb.write(now) || madeProgress
	madeProgress = wb.processReturnRsp(now) || madeProgress
	madeProgress = wb.processNewTransaction(now) || madeProgress

	return madeProgress
}

func (wb *writeBufferStage) processNewTransaction(now akita.VTimeInSec) bool {
	item := wb.cache.writeBufferBuffer.Peek()
	if item == nil {
		return false
	}

	trans := item.(*transaction)
	switch trans.action {
	case writeBufferFetch:
		return wb.processWriteBufferFetch(now, trans)
	case writeBufferEvictAndWrite:
		return wb.processWriteBufferEvictAndWrite(now, trans)
	case writeBufferEvictAndFetch:
		return wb.processWriteBufferFetchAndEvict(now, trans)
	case writeBufferFlush:
		return wb.processWriteBufferFlush(now, trans, true)
	default:
		panic("unknown transaction action")
	}
}

func (wb *writeBufferStage) processWriteBufferFetch(
	now akita.VTimeInSec,
	trans *transaction,
) bool {
	if wb.findDataLocally(trans) {
		return wb.sendFetchedDataToBank(trans)
	}

	return wb.fetchFromBottom(now, trans)
}

func (wb *writeBufferStage) findDataLocally(trans *transaction) bool {
	for _, e := range wb.inflightEviction {
		if e.evictingAddr == trans.fetchAddress {
			trans.fetchedData = e.evictingData
			return true
		}
	}

	for _, e := range wb.pendingEvictions {
		if e.evictingAddr == trans.fetchAddress {
			trans.fetchedData = e.evictingData
			return true
		}
	}
	return false
}

func (wb *writeBufferStage) sendFetchedDataToBank(trans *transaction) bool {
	bankNum := bankID(trans.block,
		wb.cache.directory.WayAssociativity(),
		len(wb.cache.dirToBankBuffers))
	bankBuf := wb.cache.writeBufferToBankBuffers[bankNum]

	if !bankBuf.CanPush() {
		trans.fetchedData = nil
		return false
	}

	trans.mshrEntry.Data = trans.fetchedData
	trans.action = bankWriteFetched
	wb.combineData(trans.mshrEntry)

	wb.cache.mshr.Remove(trans.mshrEntry.PID, trans.mshrEntry.Address)

	bankBuf.Push(trans)

	wb.cache.writeBufferBuffer.Pop()
	return true
}

func (wb *writeBufferStage) fetchFromBottom(
	now akita.VTimeInSec,
	trans *transaction,
) bool {
	if wb.tooManyInflightFetches() {
		return false
	}
	lowModulePort := wb.cache.lowModuleFinder.Find(trans.fetchAddress)
	if strings.Contains(lowModulePort.Name(), "PwPort") {
		// fmt.Println("*****************here")
		// panic("something is wrong")
		if !strings.Contains(trans.read.Src.Name(), "MMU") {
			panic("something is wrong")
		}
		// wb.remoteMemAccesses++
		// if wb.remoteMemAccesses%100 == 0 {
		// fmt.Println(wb.cache.Name(), wb.remoteMemAccesses)
		// }
		// panic("everything is wrong")
	}
	// if strings.Contains(lowModulePort.Name(), "RDMA") {
	// 	// }
	// 	// if trans.read != nil && strings.Contains(trans.read.Src.Name(), "MMU") {
	// 	// fmt.Println(trans.read.Src)
	// 	if !wb.cache.topSender.CanSend(1) {
	// 		return false
	// 	}
	// 	read := mem.ReadReqBuilder{}.
	// 		WithSrc(wb.cache.TopPort).
	// 		WithDst(lowModulePort).
	// 		// WithDst(wb.rdmaPort).
	// 		WithPID(trans.fetchPID).
	// 		WithAddress(trans.fetchAddress).
	// 		WithByteSize(1 << wb.cache.log2BlockSize).
	// 		Build()
	// 	// read.ID = trans.read.ID
	// 	wb.cache.topSender.Send(read)
	// 	trans.fetchReadReq = read
	// 	wb.inflightFetch = append(wb.inflightFetch, trans)
	// 	tracing.TraceReqInitiate(read, now, wb.cache,
	// 		tracing.MsgIDAtReceiver(trans.req(), wb.cache))

	// } else {
	if !wb.cache.bottomSender.CanSend(1) {
		return false
	}
	// lowModulePort := wb.cache.lowModuleFinder.Find(trans.fetchAddress)
	read := mem.ReadReqBuilder{}.
		WithSrc(wb.cache.BottomPort).
		WithDst(lowModulePort).
		WithPID(trans.fetchPID).
		WithAddress(trans.fetchAddress).
		WithByteSize(1 << wb.cache.log2BlockSize).
		Build()
	wb.cache.bottomSender.Send(read)
	trans.fetchReadReq = read
	wb.inflightFetch = append(wb.inflightFetch, trans)
	tracing.TraceReqInitiate(read, now, wb.cache,
		tracing.MsgIDAtReceiver(trans.req(), wb.cache))
	// }

	wb.cache.writeBufferBuffer.Pop()
	return true
}

func (wb *writeBufferStage) processWriteBufferEvictAndWrite(
	now akita.VTimeInSec,
	trans *transaction,
) bool {
	if wb.writeBufferFull() {
		return false
	}

	bankNum := bankID(
		trans.block,
		wb.cache.directory.WayAssociativity(),
		len(wb.cache.dirToBankBuffers),
	)
	bankBuf := wb.cache.writeBufferToBankBuffers[bankNum]

	if !bankBuf.CanPush() {
		return false
	}

	trans.action = bankWriteHit
	bankBuf.Push(trans)

	wb.pendingEvictions = append(wb.pendingEvictions, trans)
	wb.cache.writeBufferBuffer.Pop()

	return true
}

func (wb *writeBufferStage) processWriteBufferFetchAndEvict(
	now akita.VTimeInSec,
	trans *transaction,
) bool {
	ok := wb.processWriteBufferFlush(now, trans, false)
	if ok {
		trans.action = writeBufferFetch
		return true
	}

	return false
}

func (wb *writeBufferStage) processWriteBufferFlush(
	now akita.VTimeInSec,
	trans *transaction,
	popAfterDone bool,
) bool {
	if wb.writeBufferFull() {
		return false
	}

	wb.pendingEvictions = append(wb.pendingEvictions, trans)

	if popAfterDone {
		wb.cache.writeBufferBuffer.Pop()
	}

	return true
}

func (wb *writeBufferStage) write(now akita.VTimeInSec) bool {
	if len(wb.pendingEvictions) == 0 {
		return false
	}

	trans := wb.pendingEvictions[0]

	if wb.tooManyInflightEvictions() {
		return false
	}

	if !wb.cache.bottomSender.CanSend(1) {
		return false
	}

	lowModulePort := wb.cache.lowModuleFinder.Find(trans.evictingAddr)
	write := mem.WriteReqBuilder{}.
		WithSrc(wb.cache.BottomPort).
		WithDst(lowModulePort).
		WithPID(trans.evictingPID).
		WithAddress(trans.evictingAddr).
		WithData(trans.evictingData).
		WithDirtyMask(trans.evictingDirtyMask).
		Build()
	wb.cache.bottomSender.Send(write)

	trans.evictionWriteReq = write
	wb.pendingEvictions = wb.pendingEvictions[1:]
	wb.inflightEviction = append(wb.inflightEviction, trans)

	tracing.TraceReqInitiate(write, now, wb.cache,
		tracing.MsgIDAtReceiver(trans.req(), wb.cache))

	return true
}

func (wb *writeBufferStage) processReturnRsp(now akita.VTimeInSec) bool {
	msg := wb.cache.BottomPort.Peek()
	if msg == nil {
		return false
	}

	switch msg := msg.(type) {
	case *mem.DataReadyRsp:
		return wb.processDataReadyRsp(now, msg)
	case *mem.WriteDoneRsp:
		return wb.processWriteDoneRsp(now, msg)
	default:
		panic("unknown msg type")
	}
}

func (wb *writeBufferStage) processDataReadyRsp(
	now akita.VTimeInSec,
	dataReady *mem.DataReadyRsp,
) bool {
	trans := wb.findInflightFetchByFetchReadReqID(dataReady.RespondTo)
	bankIndex := bankID(
		trans.block,
		wb.cache.directory.WayAssociativity(),
		len(wb.cache.dirToBankBuffers),
	)
	bankBuf := wb.cache.writeBufferToBankBuffers[bankIndex]

	if !bankBuf.CanPush() {
		return false
	}

	trans.fetchedData = dataReady.Data
	trans.action = bankWriteFetched
	trans.mshrEntry.Data = dataReady.Data
	wb.combineData(trans.mshrEntry)

	wb.cache.mshr.Remove(trans.mshrEntry.PID, trans.mshrEntry.Address)

	bankBuf.Push(trans)

	wb.removeInflightFetch(trans)
	wb.cache.BottomPort.Retrieve(now)

	tracing.TraceReqFinalize(trans.fetchReadReq, now, wb.cache)

	return true
}

func (wb *writeBufferStage) combineData(mshrEntry *cache.MSHREntry) {
	mshrEntry.Block.DirtyMask = make([]bool, 1<<wb.cache.log2BlockSize)
	for _, t := range mshrEntry.Requests {
		trans := t.(*transaction)
		if trans.read != nil {
			continue
		}

		mshrEntry.Block.IsDirty = true
		write := trans.write
		_, offset := getCacheLineID(write.Address, wb.cache.log2BlockSize)
		for i := 0; i < len(write.Data); i++ {
			if write.DirtyMask == nil || write.DirtyMask[i] {
				index := offset + uint64(i)
				mshrEntry.Data[index] = write.Data[i]
				mshrEntry.Block.DirtyMask[index] = true
			}
		}
	}
}

func (wb *writeBufferStage) findInflightFetchByFetchReadReqID(
	id string,
) *transaction {
	for _, t := range wb.inflightFetch {
		if t.fetchReadReq.ID == id {
			return t
		}
	}

	panic("inflight read not found")
}

func (wb *writeBufferStage) removeInflightFetch(f *transaction) {
	for i, trans := range wb.inflightFetch {
		if trans == f {
			wb.inflightFetch = append(
				wb.inflightFetch[:i],
				wb.inflightFetch[i+1:]...,
			)
			return
		}
	}

	panic("not found")
}

func (wb *writeBufferStage) processWriteDoneRsp(
	now akita.VTimeInSec,
	writeDone *mem.WriteDoneRsp,
) bool {
	for i, e := range wb.inflightEviction {
		if e.evictionWriteReq.ID == writeDone.RespondTo {
			wb.inflightEviction = append(
				wb.inflightEviction[:i],
				wb.inflightEviction[i+1:]...,
			)
			wb.cache.BottomPort.Retrieve(now)
			tracing.TraceReqFinalize(e.evictionWriteReq, now, wb.cache)
			return true
		}
	}

	panic("write request not found")
}

func (wb *writeBufferStage) writeBufferFull() bool {
	numEntry := len(wb.pendingEvictions) + len(wb.inflightEviction)
	return numEntry >= wb.writeBufferCapacity
}

func (wb *writeBufferStage) tooManyInflightFetches() bool {
	return len(wb.inflightFetch) >= wb.maxInflightFetch
}

func (wb *writeBufferStage) tooManyInflightEvictions() bool {
	return len(wb.inflightEviction) >= wb.maxInflightEviction
}

func (wb *writeBufferStage) Reset(now akita.VTimeInSec) {
	wb.cache.writeBufferBuffer.Clear()
}
