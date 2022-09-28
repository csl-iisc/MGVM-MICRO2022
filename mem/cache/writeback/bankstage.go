package writeback

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem"
	"gitlab.com/akita/util/tracing"
)

type bankStage struct {
	cache   *Cache
	bankID  int
	latency int

	cycleLeft    int
	currentTrans *transaction
}

func (s *bankStage) Tick(now akita.VTimeInSec) bool {
	if s.currentTrans != nil {
		s.cycleLeft--
		if s.cycleLeft < 0 {
			return s.finalizeTrans(now)
		}
		return true
	}
	return s.pullFromBuf()
}

func (s *bankStage) Reset(now akita.VTimeInSec) {
	s.cache.dirToBankBuffers[s.bankID].Clear()
	s.currentTrans = nil
}

func (s *bankStage) pullFromBuf() bool {
	inBuf := s.cache.writeBufferToBankBuffers[s.bankID]
	trans := inBuf.Pop()
	if trans != nil {
		s.cycleLeft = s.latency
		s.currentTrans = trans.(*transaction)
		return true
	}

	if !s.cache.writeBufferBuffer.CanPush() {
		return false
	}

	inBuf = s.cache.dirToBankBuffers[s.bankID]
	trans = inBuf.Pop()
	if trans != nil {
		t := trans.(*transaction)

		if t.action == writeBufferFetch {
			s.cache.writeBufferBuffer.Push(trans)
			return true
		}

		s.cycleLeft = s.latency
		s.currentTrans = t

		return true
	}

	return false
}

func (s *bankStage) finalizeTrans(now akita.VTimeInSec) bool {
	switch s.currentTrans.action {
	case bankReadHit:
		return s.finalizeReadHit(now)
	case bankWriteHit:
		return s.finalizeWriteHit(now)
	case bankWriteFetched:
		return s.finalizeBankWriteFetched(now)
	case bankEvictAndFetch, bankEvictAndWrite, bankEvict:
		return s.finalizeBankEviction(now)
	default:
		panic("bank action not supported")
	}
}

func (s *bankStage) finalizeReadHit(now akita.VTimeInSec) bool {
	if !s.cache.topSender.CanSend(1) {
		s.cycleLeft = 0
		return false
	}

	read := s.currentTrans.read
	addr := read.Address
	_, offset := getCacheLineID(addr, s.cache.log2BlockSize)
	block := s.currentTrans.block

	data, err := s.cache.storage.Read(
		block.CacheAddress+offset, read.AccessByteSize)
	if err != nil {
		panic(err)
	}

	s.removeTransaction(s.currentTrans)

	s.currentTrans = nil
	block.ReadCount--

	dataReadyRspBuilder := mem.DataReadyRspBuilder{}.
		WithSendTime(now).
		WithSrc(s.cache.TopPort).
		WithDst(read.Src).
		WithRspTo(read.ID).
		WithData(data)

	if read.Info != nil {
		readReqInfo := read.Info.(*mem.ReadReqInfo)
		if readReqInfo.ReturnAccessInfo {
			if readReqInfo.AccessResult != mem.ReadHit {
				panic("oh no")
			}
			dataReadyRspBuilder = dataReadyRspBuilder.WithInfo(&mem.DataReadyRspInfo{AccessResult: mem.ReadHit, Src: s.cache.Name()})
		}
	}

	dataReady := dataReadyRspBuilder.Build()
	s.cache.topSender.Send(dataReady)

	tracing.TraceReqComplete(read, now, s.cache)

	return true
}

func (s *bankStage) finalizeWriteHit(now akita.VTimeInSec) bool {
	if !s.cache.topSender.CanSend(1) {
		s.cycleLeft = 0
		return false
	}

	write := s.currentTrans.write
	addr := write.Address
	_, offset := getCacheLineID(addr, s.cache.log2BlockSize)
	block := s.currentTrans.block

	data, err := s.cache.storage.Read(
		block.CacheAddress, 1<<s.cache.log2BlockSize)
	if err != nil {
		panic(err)
	}
	dirtyMask := block.DirtyMask
	if dirtyMask == nil {
		dirtyMask = make([]bool, 1<<s.cache.log2BlockSize)
	}
	for i := 0; i < len(write.Data); i++ {
		if write.DirtyMask == nil || write.DirtyMask[i] {
			index := offset + uint64(i)
			data[index] = write.Data[i]
			dirtyMask[index] = true
		}
	}
	err = s.cache.storage.Write(block.CacheAddress, data)
	if err != nil {
		panic(err)
	}

	block.IsValid = true
	block.IsLocked = false
	block.IsDirty = true
	block.DirtyMask = dirtyMask

	s.removeTransaction(s.currentTrans)

	s.currentTrans = nil

	done := mem.WriteDoneRspBuilder{}.
		WithSendTime(now).
		WithSrc(s.cache.TopPort).
		WithDst(write.Src).
		WithRspTo(write.ID).
		Build()
	s.cache.topSender.Send(done)

	tracing.TraceReqComplete(write, now, s.cache)

	return true
}

func (s *bankStage) finalizeBankWriteFetched(now akita.VTimeInSec) bool {
	if !s.cache.mshrStageBuffer.CanPush() {
		s.cycleLeft = 0
		return false
	}

	mshrEntry := s.currentTrans.mshrEntry
	block := mshrEntry.Block
	s.cache.mshrStageBuffer.Push(mshrEntry)
	s.cache.storage.Write(block.CacheAddress, mshrEntry.Data)
	block.IsLocked = false
	block.IsValid = true

	s.currentTrans = nil

	return true
}

func (s *bankStage) removeTransaction(trans *transaction) {
	for i, t := range s.cache.inFlightTransactions {
		if trans == t {
			s.cache.inFlightTransactions = append(
				(s.cache.inFlightTransactions)[:i],
				(s.cache.inFlightTransactions)[i+1:]...)
			return
		}
	}
	panic("transaction not found")
}

func (s *bankStage) finalizeBankEviction(now akita.VTimeInSec) bool {
	if !s.cache.writeBufferBuffer.CanPush() {
		s.cycleLeft = 0
		return false
	}

	trans := s.currentTrans
	victim := trans.victim
	data, err := s.cache.storage.Read(
		victim.CacheAddress, 1<<s.cache.log2BlockSize)
	if err != nil {
		panic(err)
	}
	trans.evictingData = data

	switch trans.action {
	case bankEvict:
		trans.action = writeBufferFlush
	case bankEvictAndFetch:
		trans.action = writeBufferEvictAndFetch
	case bankEvictAndWrite:
		trans.action = writeBufferEvictAndWrite
	default:
		panic("unsupported action")
	}

	s.cache.writeBufferBuffer.Push(trans)

	s.currentTrans = nil

	return true
}
