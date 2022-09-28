package idealtlb

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem/device"
	"gitlab.com/akita/util/tracing"
)

type IdealTLB struct {
	*akita.TickingComponent

	TopPort     akita.Port
	BottomPort  akita.Port
	ControlPort akita.Port

	pageTable device.PageTable

	numReqPerCycle int

	isPaused bool
}

// Tick defines how TLB update states at each cycle
func (tlb *IdealTLB) Tick(now akita.VTimeInSec) bool {
	madeProgress := false

	if !tlb.isPaused {
		for i := 0; i < tlb.numReqPerCycle; i++ {
			madeProgress = tlb.lookup(now) || madeProgress
		}
	}

	return madeProgress
}

func (tlb *IdealTLB) lookup(now akita.VTimeInSec) bool {
	msg := tlb.TopPort.Peek()
	if msg == nil {
		return false
	}

	req := msg.(*device.TranslationReq)

	// get the page magically
	page, found := tlb.pageTable.Find(req.PID, req.VAddr)
	if !found {
		panic("Physical to Virtual address translation not found!!")
	}

	return tlb.handleTranslationHit(now, req, page)
}

func (tlb *IdealTLB) handleTranslationHit(
	now akita.VTimeInSec,
	req *device.TranslationReq,
	page device.Page,
) bool {
	ok := tlb.sendRspToTop(now, req, page)
	if !ok {
		return false
	}

	tlb.TopPort.Retrieve(now)

	tracing.StopTracingNetworkReq(req, now, tlb)

	tracing.TraceReqReceive(req, now, tlb)
	tracing.AddTaskStep(
		tracing.MsgIDAtReceiver(req, tlb),
		now, tlb,
		"tlb-hit",
	)
	tracing.TraceReqComplete(req, now, tlb)

	return true
}

func (tlb *IdealTLB) sendRspToTop(
	now akita.VTimeInSec,
	req *device.TranslationReq,
	page device.Page,
) bool {
	rsp := device.TranslationRspBuilder{}.
		WithSendTime(now).
		WithSrc(tlb.TopPort).
		WithDst(req.Src).
		WithRspTo(req.ID).
		WithPage(page).
		Build()
	err := tlb.TopPort.Send(rsp)
	if err == nil {
		tracing.StartTracingNetworkReq(rsp, now, tlb, req)
		return true
	}

	return false
}
