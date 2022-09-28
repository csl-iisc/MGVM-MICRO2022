package writeback

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem"
	"gitlab.com/akita/mem/cache"
	"gitlab.com/akita/util/tracing"
)

type mshrStage struct {
	cache *Cache

	processingMSHREntry *cache.MSHREntry
}

func (s *mshrStage) Tick(now akita.VTimeInSec) bool {
	if s.processingMSHREntry != nil {
		return s.processOneReq(now)
	}

	item := s.cache.mshrStageBuffer.Pop()
	if item == nil {
		return false
	}

	s.processingMSHREntry = item.(*cache.MSHREntry)
	return s.processOneReq(now)
}

func (s *mshrStage) Reset(now akita.VTimeInSec) {
	s.processingMSHREntry = nil
	s.cache.mshrStageBuffer.Clear()
}

func (s *mshrStage) processOneReq(now akita.VTimeInSec) bool {
	if !s.cache.topSender.CanSend(1) {
		return false
	}

	mshrEntry := s.processingMSHREntry
	trans := mshrEntry.Requests[0].(*transaction)

	transactionPresent := s.findTransaction(trans)

	if transactionPresent {
		s.removeTransaction(trans)

		if trans.read != nil {
			s.respondRead(now, trans.read, mshrEntry.Data)
		} else {
			s.respondWrite(now, trans.write)
		}

		mshrEntry.Requests = mshrEntry.Requests[1:]
		if len(mshrEntry.Requests) == 0 {
			s.processingMSHREntry = nil
		}

		return true
	}

	mshrEntry.Requests = mshrEntry.Requests[1:]
	if len(mshrEntry.Requests) == 0 {
		s.processingMSHREntry = nil
	}

	return true
}

func (s *mshrStage) respondRead(
	now akita.VTimeInSec,
	read *mem.ReadReq,
	data []byte,
) {
	_, offset := getCacheLineID(read.Address, s.cache.log2BlockSize)
	dataReadyRspBuilder := mem.DataReadyRspBuilder{}.
		WithSendTime(now).
		WithSrc(s.cache.TopPort).
		WithDst(read.Src).
		WithRspTo(read.ID).
		WithData(data[offset : offset+read.AccessByteSize])

	if read.Info != nil {
		readReqInfo := read.Info.(*mem.ReadReqInfo)
		if readReqInfo.ReturnAccessInfo {
			dataReadyRspBuilder = dataReadyRspBuilder.WithInfo(&mem.DataReadyRspInfo{AccessResult: readReqInfo.AccessResult, Src: s.cache.Name()})
		}
	}
	dataReady := dataReadyRspBuilder.Build()
	s.cache.topSender.Send(dataReady)

	tracing.TraceReqComplete(read, now, s.cache)
}

func (s *mshrStage) respondWrite(
	now akita.VTimeInSec,
	write *mem.WriteReq,
) {
	writeDoneRsp := mem.WriteDoneRspBuilder{}.
		WithSendTime(now).
		WithSrc(s.cache.TopPort).
		WithDst(write.Src).
		WithRspTo(write.ID).
		Build()
	s.cache.topSender.Send(writeDoneRsp)

	tracing.TraceReqComplete(write, now, s.cache)
}

func (s *mshrStage) removeTransaction(trans *transaction) {
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

func (s *mshrStage) findTransaction(trans *transaction) bool {
	for _, t := range s.cache.inFlightTransactions {
		if trans == t {
			return true
		}
	}
	return false
}
