package writeback

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem"
	"gitlab.com/akita/util/tracing"
)

type pageWalkCacheBankStage struct {
	cache   *PageWalkCache
	bankID  int
	latency int

	cycleLeft    int
	currentTrans *transaction
}

func (s *pageWalkCacheBankStage) Tick(now akita.VTimeInSec) bool {
	if s.currentTrans != nil {
		s.cycleLeft--
		if s.cycleLeft < 0 {
			return s.finalizeTrans(now)
		}
		return true
	}
	return s.pullFromBuf()
}

func (s *pageWalkCacheBankStage) Reset(now akita.VTimeInSec) {
	s.cache.dirToBankBuffers[s.bankID].Clear()
	s.currentTrans = nil
}

func (s *pageWalkCacheBankStage) pullFromBuf() bool {
	inBuf := s.cache.dirToBankBuffers[s.bankID]
	trans := inBuf.Pop()
	if trans != nil {
		s.cycleLeft = s.latency
		s.currentTrans = trans.(*transaction)
		return true
	}

	return false
}

func (s *pageWalkCacheBankStage) finalizeTrans(now akita.VTimeInSec) bool {
	switch s.currentTrans.action {
	case bankReadHit:
		return s.finalizeReadHit(now)
	case bankWriteHit:
		return s.finalizeWriteHit(now)
	default:
		panic("bank action not supported")
	}
}

func (s *pageWalkCacheBankStage) finalizeReadHit(now akita.VTimeInSec) bool {
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
	// pAddr := binary.LittleEndian.Uint64(data)
	// h := fmt.Sprintf("%x", read.Address)
	// h1 := fmt.Sprintf("%x", pAddr)
	// fmt.Println("pwc:", h, h1)

	s.removeTransaction(s.currentTrans)

	s.currentTrans = nil
	block.ReadCount--

	dataReady := mem.DataReadyRspBuilder{}.
		WithSendTime(now).
		WithSrc(s.cache.TopPort).
		WithDst(read.Src).
		WithRspTo(read.ID).
		WithData(data).
		Build()
	s.cache.topSender.Send(dataReady)

	tracing.TraceReqComplete(read, now, s.cache)

	return true
}

func (s *pageWalkCacheBankStage) finalizeWriteHit(now akita.VTimeInSec) bool {
	if !s.cache.topSender.CanSend(1) {
		s.cycleLeft = 0
		return false
	}

	write := s.currentTrans.write
	block := s.currentTrans.block
	err := s.cache.storage.Write(block.CacheAddress, write.Data)
	// h := fmt.Sprintf("bank write: %x", block.CacheAddress)
	// fmt.Println(h)
	// fmt.Println(write.Data)
	if err != nil {
		panic(err)
	}

	block.IsValid = true
	block.IsLocked = false
	block.IsDirty = true

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

func (s *pageWalkCacheBankStage) removeTransaction(trans *transaction) {
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
