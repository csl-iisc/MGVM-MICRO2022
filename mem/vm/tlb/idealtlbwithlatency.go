package tlb

import (
	"log"
	"reflect"

	//"strings"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem/cache"
	"gitlab.com/akita/mem/device"
	"gitlab.com/akita/util"
	"gitlab.com/akita/util/pipelining"
	"gitlab.com/akita/util/tracing"
)

// A LatTLB is a cache that maintains some page information.
type IdealLatTLB struct {
	*akita.TickingComponent

	TopPort     akita.Port
	BottomPort  akita.Port
	ControlPort akita.Port

	numReqPerCycle int
	latency        int

	pipeline     pipelining.Pipeline
	lookupBuffer util.Buffer

	pageTable device.PageTable

	isPaused bool

	CommandProcessor akita.Port
}

// GetPipeline gets the pipeline in the LatTLB
func (tlb *IdealLatTLB) GetPipeline() pipelining.Pipeline {
	return tlb.pipeline
}

// GetPipeline gets the pipeline in the LatTLB
func (tlb *IdealLatTLB) GetTopPort() akita.Port {
	return tlb.TopPort
}

// GetPipeline gets the pipeline in the LatTLB
func (tlb *IdealLatTLB) GetBottomPort() akita.Port {
	return tlb.BottomPort
}

// GetPipeline gets the pipeline in the LatTLB
func (tlb *IdealLatTLB) GetControlPort() akita.Port {
	return tlb.ControlPort
}

func (tlb *IdealLatTLB) SetLowModuleFinder(lmf cache.LowModuleFinder) {
	panic("idealLatTlb has no low module!")
}

// Tick defines how LatTLB update states at each cycle
func (tlb *IdealLatTLB) Tick(now akita.VTimeInSec) bool {
	madeProgress := false

	madeProgress = tlb.performCtrlReq(now) || madeProgress

	if !tlb.isPaused {
		for i := 0; i < tlb.numReqPerCycle; i++ {
			madeProgress = tlb.lookup(now) || madeProgress
		}
		madeProgress = tlb.pipeline.Tick(now) || madeProgress
		for i := 0; i < tlb.numReqPerCycle; i++ {
			madeProgress = tlb.parseFromTop(now) || madeProgress
		}
	}

	return madeProgress
}

// func (tlb *IdealLatTLB) collectCoalescingStat(now akita.VTimeInSec) {
// 	bufContainer := tlb.TopPort.(akita.MsgBufferContainer) //.(*akita.LimitNumMsgPort)
// 	buf := bufContainer.GetBuffer()
// 	coalMine := make(map[uint64]int)
// 	for _, msg := range buf {
// 		t := msg.(*device.TranslationReq)
// 		addr := t.VAddr
// 		coalMine[addr] = coalMine[addr] + 1
// 	}
// 	count := 0
// 	totalcoal := 0
// 	for addr := range coalMine {
// 		if coalMine[addr] > 1 {
// 			count++
// 			totalcoal += coalMine[addr]
// 		}
// 	}
// 	tracing.StartTask("", "", now, tlb,
// 		"buflen", strconv.Itoa(len(buf)), nil)
// 	if len(buf) > 0 {
// 		tracing.StartTask("", "", now, tlb,
// 			"bufleng0", strconv.Itoa(len(buf)), nil)
// 		tracing.StartTask("", "", now, tlb,
// 			"coalesceaddr", strconv.Itoa(count), nil)
// 		tracing.StartTask("", "", now, tlb,
// 			"coalesce", strconv.Itoa(totalcoal), nil)
// 	}
// }

func (tlb *IdealLatTLB) parseFromTop(now akita.VTimeInSec) bool {
	msg := tlb.TopPort.Peek()
	collectCoalescingStat(tlb, now)
	if msg == nil {
		return false
	}

	req := msg.(*device.TranslationReq)
	// push to waiting queue
	if tlb.pipeline.CanAccept() {
		pipelineItem := tlbPipelineItem{
			taskID:         akita.GetIDGenerator().Generate(),
			translationReq: req,
		}
		tlb.pipeline.Accept(now, pipelineItem)
		tlb.TopPort.Retrieve(now)
		tracing.TraceReqReceive(req, now, tlb)
		tracing.StopTracingNetworkReq(req, now, tlb)
		return true
	}
	return false
}

func (tlb *IdealLatTLB) lookup(now akita.VTimeInSec) bool {

	// pop from waiting queue
	item := tlb.lookupBuffer.Peek()
	if item == nil {
		return false
	}
	pipelineItem := item.(tlbPipelineItem)
	req := pipelineItem.translationReq
	page, ok := tlb.pageTable.Find(req.PID, req.VAddr)
	if !ok {
		panic("something went wrong!")
	}
	return tlb.handleTranslationHit(now, req, page)
}

func (tlb *IdealLatTLB) handleTranslationHit(
	now akita.VTimeInSec,
	req *device.TranslationReq,
	page device.Page,
) bool {
	ok := tlb.sendRspToTop(now, req, page)
	if !ok {
		return false
	}
	tlb.lookupBuffer.Pop()

	tracing.AddTaskStep(
		tracing.MsgIDAtReceiver(req, tlb),
		now, tlb,
		"tlb-hit",
	)
	tracing.TraceReqComplete(req, now, tlb)

	return true
}

func (tlb *IdealLatTLB) sendRspToTop(
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
		WithAccessResult(device.TLBHit).
		WithSrcL2TLB(tlb.Name()).
		Build()

	err := tlb.TopPort.Send(rsp)
	if err == nil {
		tracing.StartTracingNetworkReq(rsp, now, tlb, req)
		return true
	}
	return false
}

func (tlb *IdealLatTLB) performCtrlReq(now akita.VTimeInSec) bool {
	item := tlb.ControlPort.Peek()
	if item == nil {
		return false
	}

	item = tlb.ControlPort.Retrieve(now)

	switch req := item.(type) {
	case *TLBFlushReq:
		return tlb.handleTLBFlush(now, req)
	case *TLBRestartReq:
		return tlb.handleTLBRestart(now, req)
	default:
		log.Panicf("cannot process request %s", reflect.TypeOf(req))
	}

	return true
}

// func (tlb *IdealLatTLB) visit(setID, wayID int) int {
// 	set := tlb.Sets[setID]
// 	mruPosition := set.Visit(wayID)
// 	return mruPosition
// }

func (tlb *IdealLatTLB) handleTLBFlush(now akita.VTimeInSec, req *TLBFlushReq) bool {
	rsp := TLBFlushRspBuilder{}.
		WithSrc(tlb.ControlPort).
		WithDst(req.Src).
		WithSendTime(now).
		Build()

	err := tlb.ControlPort.Send(rsp)
	if err != nil {
		return false
	}

	// for _, vAddr := range req.VAddr {
	// 	setID := tlb.vAddrToSetID(vAddr)
	// 	set := tlb.Sets[setID]
	// 	wayID, page, found := set.Lookup(req.PID, vAddr)
	// 	if !found {
	// 		continue
	// 	}

	// 	page.Valid = false
	// 	set.Update(wayID, page)
	// }

	// tlb.mshr.Reset()
	// tlb.isPaused = true
	return true
}

func (tlb *IdealLatTLB) handleTLBRestart(now akita.VTimeInSec, req *TLBRestartReq) bool {
	rsp := TLBRestartRspBuilder{}.
		WithSendTime(now).
		WithSrc(tlb.ControlPort).
		WithDst(req.Src).
		Build()

	err := tlb.ControlPort.Send(rsp)
	if err != nil {
		return false
	}

	tlb.isPaused = false

	for tlb.TopPort.Retrieve(now) != nil {
		tlb.TopPort.Retrieve(now)
	}

	// for tlb.BottomPort.Retrieve(now) != nil {
	// 	tlb.BottomPort.Retrieve(now)
	// }

	return true
}

func (tlb *IdealLatTLB) GetFrontQueueLength() int {
	bufContainer := tlb.GetTopPort().(akita.MsgBufferContainer) //.(*akita.LimitNumMsgPort)
	buf := bufContainer.GetBuffer()
	return len(buf)
}

func (tlb *IdealLatTLB) SetCommandProcessor(cp akita.Port) {
	tlb.CommandProcessor = cp
}
