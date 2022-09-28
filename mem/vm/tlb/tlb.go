package tlb

import (
	// "fmt"
	"fmt"
	"log"
	"reflect"
	"strings"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem/cache"
	"gitlab.com/akita/mem/device"
	"gitlab.com/akita/mem/vm/tlb/internal"
	"gitlab.com/akita/util/tracing"
)

// A TLB is a cache that maintains some page information.
type TLB struct {
	*akita.TickingComponent

	TopPort     akita.Port
	BottomPort  akita.Port
	ControlPort akita.Port

	LowModule       akita.Port
	LowModuleFinder cache.LowModuleFinder

	numSets        int
	numWays        int
	pageSize       uint64
	numReqPerCycle int

	Sets []internal.Set

	mshr                mshr
	respondingMSHREntry *mshrEntry

	isPaused bool
}

// GetNumSets gets the number of sets in the TLB
func (tlb *TLB) GetNumSets() int {
	return tlb.numSets
}

// GetNumWays gets the number of ways in the TLB
func (tlb *TLB) GetNumWays() int {
	return tlb.numWays
}

func (tlb *TLB) SetLowModuleFinder(lmf cache.LowModuleFinder) {
	tlb.LowModuleFinder = lmf
}

// Reset sets all the entries int he TLB to be invalid
func (tlb *TLB) reset() {
	tlb.Sets = make([]internal.Set, tlb.numSets)
	for i := 0; i < tlb.numSets; i++ {
		set := internal.NewSet(tlb.numWays)
		tlb.Sets[i] = set
	}
}

// Tick defines how TLB update states at each cycle
func (tlb *TLB) Tick(now akita.VTimeInSec) bool {
	madeProgress := false

	madeProgress = tlb.performCtrlReq(now) || madeProgress

	if !tlb.isPaused {
		for i := 0; i < tlb.numReqPerCycle; i++ {
			madeProgress = tlb.respondMSHREntry(now) || madeProgress
		}

		for i := 0; i < tlb.numReqPerCycle; i++ {
			madeProgress = tlb.lookup(now) || madeProgress
		}

		for i := 0; i < tlb.numReqPerCycle; i++ {
			madeProgress = tlb.parseBottom(now) || madeProgress
		}
	}

	return madeProgress
}

func (tlb *TLB) respondMSHREntry(now akita.VTimeInSec) bool {
	if tlb.respondingMSHREntry == nil {
		return false
	}

	mshrEntry := tlb.respondingMSHREntry
	page := mshrEntry.page
	req := mshrEntry.Requests[0]
	var accessResult device.AccessResult
	if mshrEntry.NumResponded() == 0 {
		accessResult = device.TLBMiss
	} else {
		accessResult = device.TLBMshrHit
	}
	rspToTop := device.TranslationRspBuilder{}.
		WithSendTime(now).
		WithSrc(tlb.TopPort).
		WithDst(req.Src).
		WithRspTo(req.ID).
		WithPage(page).
		WithAccessResult(accessResult).
		//device.TLBMiss).
		Build()
	err := tlb.TopPort.Send(rspToTop)
	if err != nil {
		return false
	}
	mshrEntry.IncNumRespondedByOne()
	mshrEntry.Requests = mshrEntry.Requests[1:]
	if len(mshrEntry.Requests) == 0 {
		tlb.respondingMSHREntry = nil
	}
	tracing.StartTracingNetworkReq(rspToTop, now, tlb, req)
	tracing.TraceReqComplete(req, now, tlb)
	return true
}

func (tlb *TLB) lookup(now akita.VTimeInSec) bool {
	msg := tlb.TopPort.Peek()
	if msg == nil {
		return false
	}
	req := msg.(*device.TranslationReq)

	mshrEntry := tlb.mshr.Query(req.PID, req.VAddr)
	if mshrEntry != nil {
		ok := tlb.processTLBMSHRHit(now, mshrEntry, req)
		if ok {
			tracing.TraceReqReceive(req, now, tlb)
			tracing.AddTaskStep(
				tracing.MsgIDAtReceiver(req /*mshrEntry.Requests[0]*/, tlb),
				now, tlb,
				"tlb-mshr-hit",
			)
			tlb.TopPort.Retrieve(now)
			tracing.TraceReqReceive(req, now, tlb)
			tracing.StopTracingNetworkReq(req, now, tlb)
			return true
		}
		return false
	}

	setID := tlb.vAddrToSetID(req.VAddr)
	set := tlb.Sets[setID]
	wayID, page, found := set.Lookup(req.PID, req.VAddr)
	if found && page.Valid {
		return tlb.handleTranslationHit(now, req, setID, wayID, page)
	}

	return tlb.handleTranslationMiss(now, req)
}

func (tlb *TLB) handleTranslationHit(
	now akita.VTimeInSec,
	req *device.TranslationReq,
	setID, wayID int,
	page device.Page,
) bool {
	ok := tlb.sendRspToTop(now, req, page)
	if !ok {
		return false
	}

	// mruPos := tlb.visit(setID, wayID)
	tlb.visit(setID, wayID)
	tlb.TopPort.Retrieve(now)

	tracing.TraceReqReceive(req, now, tlb)
	tracing.StopTracingNetworkReq(req, now, tlb)
	tracing.AddTaskStep(
		tracing.MsgIDAtReceiver(req, tlb),
		now, tlb,
		"tlb-hit",
	)
	// tracing.AddTaskStep(
	// 	tracing.MsgIDAtReceiver(req, tlb),
	// 	now, tlb,
	// 	fmt.Sprintf("tlb-hit-at-set-%d-mru-pos-%d", setID, mruPos),
	// )
	// s := fmt.Sprintf("tlb-hit-at-set-%d-mru-pos-%d\n", setID, mruPos)
	// fmt.Println(s)
	tracing.TraceReqComplete(req, now, tlb)

	return true
}

func (tlb *TLB) handleTranslationMiss(
	now akita.VTimeInSec,
	req *device.TranslationReq,
) bool {
	// if tlb.shouldDoPageWalk(req.VAddr) {
	if tlb.mshr.IsFull() {
		return false
	}
	fetched := tlb.fetchBottom(now, req)
	if fetched {
		tlb.TopPort.Retrieve(now)
		tracing.StopTracingNetworkReq(req, now, tlb)
		tracing.TraceReqReceive(req, now, tlb)
		tracing.AddTaskStep(
			tracing.MsgIDAtReceiver(req, tlb),
			now, tlb,
			"tlb-miss",
		)
		return true
	}
	return false
}

func (tlb *TLB) vAddrToSetID(vAddr uint64) (setID int) {
	return int(vAddr / tlb.pageSize % uint64(tlb.numSets))
}

func (tlb *TLB) sendRspToTop(
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
		Build()

	err := tlb.TopPort.Send(rsp)
	if err == nil {
		tracing.StartTracingNetworkReq(rsp, now, tlb, req)
		return true
	}
	return false
}

func (tlb *TLB) processTLBMSHRHit(
	now akita.VTimeInSec,
	mshrEntry *mshrEntry,
	req *device.TranslationReq,
) bool {
	// if len(mshrEntry.Requests) < 8 {
	// 	mshrEntry.Requests = append(mshrEntry.Requests, req)
	// 	return true
	// }
	// return false
	mshrEntry.Requests = append(mshrEntry.Requests, req)
	return true
	/*	tracing.AddTaskStep(
		tracing.MsgIDAtReceiver(req /*mshrEntry.Requests[0], tlb),
		now, tlb,
		"tlb-mshr-hit",
	) */
}

func (tlb *TLB) fetchBottom(now akita.VTimeInSec, req *device.TranslationReq) bool {
	dstPort := tlb.LowModuleFinder.Find(req.VAddr)

	fetchBottom := device.TranslationReqBuilder{}.
		WithSendTime(now).
		WithSrc(tlb.BottomPort).
		// WithDst(tlb.LowModule).
		WithDst(dstPort).
		WithPID(req.PID).
		WithVAddr(req.VAddr).
		WithDeviceID(req.DeviceID).
		Build()
	err := tlb.BottomPort.Send(fetchBottom)
	if err != nil {
		return false
	}

	mshrEntry := tlb.mshr.Add(req.PID, req.VAddr)
	mshrEntry.Requests = append(mshrEntry.Requests, req)
	mshrEntry.reqToBottom = fetchBottom

	tracing.StartTask(
		fetchBottom.Meta().ID+"_L2TLB_stats",
		tracing.MsgIDAtReceiver(req, tlb),
		now,
		tlb,
		"L2TLB_stats",
		reflect.TypeOf(req).String(),
		req,
	)
	tracing.StartTracingNetworkReq(fetchBottom, now, tlb, req)
	// tracing.TraceReqReceive(req, now, tlb)
	tracing.TraceReqInitiate(fetchBottom, now, tlb,
		tracing.MsgIDAtReceiver(req, tlb))

	return true
}

func (tlb *TLB) getRemoteVLocal(rspSrc string) (str string) {
	if getChipletNum(tlb.Name()) != getChipletNum(rspSrc) {
		str = "remote"
	} else {
		str = "local"
	}
	return
}

func getChipletNum(srcL2TLB string) (i string) {
	i = "chiplet-" + strings.Split(srcL2TLB, "_")[1][1:2]
	return
}

func getTaskStep(origin string, accessResult device.AccessResult) (step string) {
	step = origin + "-"
	switch accessResult {
	case device.TLBHit:
		step += "TLBHit"
	case device.TLBMiss:
		step += "TLBMiss"
	case device.TLBMshrHit:
		step += "TLBMshrHit"
	}
	return
}

func (tlb *TLB) parseBottom(now akita.VTimeInSec) bool {
	if tlb.respondingMSHREntry != nil {
		return false
	}

	item := tlb.BottomPort.Peek()
	if item == nil {
		return false
	}

	rsp := item.(*device.TranslationRsp)
	page := rsp.Page
	//	fmt.Println(rsp.Meta().ID)//, req.Meta().ID)

	mshrEntryPresent := tlb.mshr.IsEntryPresent(rsp.Page.PID, rsp.Page.VAddr)
	if !mshrEntryPresent {
		panic("oh no!")
		tlb.BottomPort.Retrieve(now)
		tracing.TraceReqFinalize(rsp, now, tlb)
		return true
	}

	setID := tlb.vAddrToSetID(page.VAddr)
	set := tlb.Sets[setID]
	wayID, ok := tlb.Sets[setID].Evict()
	if !ok {
		panic("failed to evict")
	}
	set.Update(wayID, page)
	set.Visit(wayID)

	mshrEntry := tlb.mshr.GetEntry(rsp.Page.PID, rsp.Page.VAddr)
	tlb.respondingMSHREntry = mshrEntry
	mshrEntry.page = page

	tlb.mshr.Remove(rsp.Page.PID, rsp.Page.VAddr)
	tlb.BottomPort.Retrieve(now)

	tracing.StopTracingNetworkReq(rsp, now, tlb)
	tracing.TraceReqFinalize(mshrEntry.reqToBottom, now, tlb)
	//fmt.Println("here")
	taskStepRemoteVLocal := getTaskStep(tlb.getRemoteVLocal(rsp.SrcL2TLB), rsp.HitOrMiss)
	taskStepSrcL2TLB := getTaskStep(getChipletNum(rsp.SrcL2TLB), rsp.HitOrMiss)
	//fmt.Println(tlb.Name())
	tracing.AddTaskStep(
		mshrEntry.reqToBottom.Meta().ID+"_L2TLB_stats",
		now, tlb,
		taskStepRemoteVLocal,
	)
	tracing.AddTaskStep(
		mshrEntry.reqToBottom.Meta().ID+"_L2TLB_stats",
		now, tlb,
		taskStepSrcL2TLB,
	)
	tracing.EndTask(
		mshrEntry.reqToBottom.Meta().ID+"_L2TLB_stats",
		now,
		tlb,
	)

	return true
}

func (tlb *TLB) performCtrlReq(now akita.VTimeInSec) bool {
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
	case *akita.TLBIndexingSwitchMsg:
		return tlb.switchIndexing(now, req)
	default:
		log.Panicf("cannot process request %s", reflect.TypeOf(req))
	}

	return true
}

func (tlb *TLB) visit(setID, wayID int) int {
	set := tlb.Sets[setID]
	mruPosition := set.Visit(wayID)
	return mruPosition
}

func (tlb *TLB) handleTLBFlush(now akita.VTimeInSec, req *TLBFlushReq) bool {
	rsp := TLBFlushRspBuilder{}.
		WithSrc(tlb.ControlPort).
		WithDst(req.Src).
		WithSendTime(now).
		Build()

	err := tlb.ControlPort.Send(rsp)
	if err != nil {
		return false
	}

	for _, vAddr := range req.VAddr {
		setID := tlb.vAddrToSetID(vAddr)
		set := tlb.Sets[setID]
		wayID, page, found := set.Lookup(req.PID, vAddr)
		if !found {
			continue
		}

		page.Valid = false
		set.Update(wayID, page)
	}

	tlb.mshr.Reset()
	tlb.isPaused = true
	return true
}

func (tlb *TLB) handleTLBRestart(now akita.VTimeInSec, req *TLBRestartReq) bool {
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

	for tlb.BottomPort.Retrieve(now) != nil {
		tlb.BottomPort.Retrieve(now)
	}

	return true
}

func (tlb *TLB) switchIndexing(now akita.VTimeInSec,
	req *akita.TLBIndexingSwitchMsg) bool {
	lmf := tlb.LowModuleFinder.(*cache.CustomTwoLevelLowModuleFinder)

	HashFunc4KXOR := func(address uint64) uint64 {
		NumBitsPerTerm := 2
		NumTerms := 4
		index := uint64(0)
		mask := (uint64(1) << NumBitsPerTerm) - 1
		for i := 0; i < NumTerms; i++ {
			index = index ^ (address & mask)
			address = address >> NumBitsPerTerm
		}
		return index
	}

	// some parameters hardcoded
	HashFuncHSL := func(address uint64) uint64 {
		return (address / uint64(req.TLBInterleaving)) % 4
	}

	if req.TLBIndexingSwitch == akita.TLBIndexingSwitch4K {
		lmf.Hashfunc = HashFunc4KXOR
	} else {
		lmf.Hashfunc = HashFuncHSL
	}

	fmt.Println(tlb.Name(), "L1 TLB ", tlb.Name(), "switching to ", req.TLBIndexingSwitch, " sir", now, req.TLBInterleaving)
	return true
}
