// Package remotetranslation provides implementation of a remote translation
// unit.
package remotetranslation

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem/cache"
	"gitlab.com/akita/mem/device"
	"gitlab.com/akita/util/tracing"
)

// A CoalescingRemoteTranslationUnit helps construct distributed TLBs.
type CoalescingRemoteTranslationUnit struct {
	*DefaultRTU
	mshr                mshr
	respondingMSHREntry *mshrEntry
}

// Tick ticks. Ha.
func (rtu *CoalescingRemoteTranslationUnit) Tick(now akita.VTimeInSec) bool {
	madeProgress := false
	for i := 0; i < rtu.numTransPerCycle; i++ {
		madeProgress = rtu.respondMSHREntry(now) || madeProgress
	}
	madeProgress = rtu.processFromL1(now) || madeProgress
	madeProgress = rtu.processFromL2(now) || madeProgress
	madeProgress = rtu.processReqFromOutside(now) || madeProgress
	madeProgress = rtu.processRspFromOutside(now) || madeProgress
	return madeProgress
}

func (rtu *CoalescingRemoteTranslationUnit) processFromL1(now akita.VTimeInSec) bool {
	madeProgress := false
	for i := 0; i < rtu.numTransPerCycle; i++ {
		msg := rtu.ToL1.Peek()
		if msg == nil {
			break
		}
		req := msg.(*device.TranslationReq)
		madeProgress = rtu.processReqFromL1(now, req) || madeProgress
	}
	rtu.collectCoalescingStats(now, rtu.transactionsFromInside, "outgoing")
	return madeProgress
}

func (rtu *CoalescingRemoteTranslationUnit) processRspFromOutside(now akita.VTimeInSec) bool {
	madeProgress := false
	for i := 0; i < rtu.numTransPerCycle; i++ {
		msg := rtu.ResponsePort.Peek()
		if msg == nil {
			return madeProgress
		}
		req := msg.(*device.TranslationRsp)
		madeProgress = rtu.processRspFromOutsideAux(now, req) || madeProgress
	}
	return madeProgress
}

func (rtu *CoalescingRemoteTranslationUnit) processReqFromL1(
	now akita.VTimeInSec,
	req *device.TranslationReq,
) bool {
	mshrEntry := rtu.mshr.Query(req.PID, req.VAddr)
	if mshrEntry != nil {
		mshrEntry.Requests = append(mshrEntry.Requests, req)
		rtu.ToL1.Retrieve(now)
		return true
	} else {
		reqSent := rtu.DefaultRTU.processReqFromL1(now, req)
		if reqSent {
			mshrEntry := rtu.mshr.Add(req.PID, req.VAddr)
			mshrEntry.Requests = append(mshrEntry.Requests, req)
			return true
		}
	}
	return false
}

func (rtu *CoalescingRemoteTranslationUnit) respondMSHREntry(now akita.VTimeInSec) bool {
	if rtu.respondingMSHREntry == nil {
		return false
	}

	mshrEntry := rtu.respondingMSHREntry
	req := mshrEntry.Requests[0]
	rsp := mshrEntry.rspToL1

	rspToInside := rtu.cloneRsp(rsp, req.Meta().ID)
	rspToInside.Meta().SendTime = now
	rspToInside.Meta().Src = rtu.ToL1
	rspToInside.Meta().Dst = req.Meta().Src

	err := rtu.ToL1.Send(rspToInside)
	if err == nil {

		mshrEntry.Requests = mshrEntry.Requests[1:]
		if len(mshrEntry.Requests) == 0 {
			rtu.respondingMSHREntry = nil
		}

		tracing.StartTracingNetworkReq(rspToInside, now, rtu, rsp)
		tracing.TraceReqComplete(req, now, rtu)
		return true
	}
	return false
}

func (rtu *CoalescingRemoteTranslationUnit) processRspFromOutsideAux(
	now akita.VTimeInSec,
	rsp *device.TranslationRsp,
) bool {
	if rtu.respondingMSHREntry != nil {
		return false
	}
	transactionIndex := rtu.findTransactionByRspToID(
		rsp.RespondTo, rtu.transactionsFromInside)
	trans := rtu.transactionsFromInside[transactionIndex]

	mshrEntry := rtu.mshr.GetEntry(rsp.Page.PID, rsp.Page.VAddr)
	rtu.respondingMSHREntry = mshrEntry
	mshrEntry.rspToL1 = rsp
	mshrEntry.page = rsp.Page

	rtu.mshr.Remove(rsp.Page.PID, rsp.Page.VAddr)

	rtu.ResponsePort.Retrieve(now)

	tracing.StopTracingNetworkReq(rsp, now, rtu)
	tracing.TraceReqFinalize(trans.toOutside, now, rtu)

	rtu.transactionsFromInside =
		append(rtu.transactionsFromInside[:transactionIndex],
			rtu.transactionsFromInside[transactionIndex+1:]...)

	return true

}

func NewCoalescingRemoteTranslationUnit(
	name string,
	engine akita.Engine,
	localModules cache.LowModuleFinder,
	remoteModules cache.LowModuleFinder,
) *CoalescingRemoteTranslationUnit {
	rtu := new(CoalescingRemoteTranslationUnit)
	rtu.DefaultRTU = NewRemoteTranslationUnit(name, engine, localModules, remoteModules).(*DefaultRTU)
	rtu.TickingComponent = akita.NewTickingComponent(name, engine, 1*akita.GHz, rtu)
	rtu.mshr = newMSHR(360)
	return rtu
}
