package tlb

import (
	"fmt"
	"log"

	// "math"
	"reflect"
	"strconv"

	// "strings"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem/cache"
	"gitlab.com/akita/mem/device"
	"gitlab.com/akita/mem/vm/tlb/internal"
	"gitlab.com/akita/util"
	"gitlab.com/akita/util/pipelining"
	"gitlab.com/akita/util/tracing"
)

type tlbPipelineItem struct {
	taskID         string
	translationReq *device.TranslationReq
}

func (t tlbPipelineItem) TaskID() string {
	return t.taskID
}

type L2TLB interface {
	tracing.NamedHookable
	GetTopPort() akita.Port
	GetBottomPort() akita.Port
	GetControlPort() akita.Port
	SetLowModuleFinder(cache.LowModuleFinder)
	GetPipeline() pipelining.Pipeline
	GetFrontQueueLength() int

	SetCommandProcessor(akita.Port)
}

type TLBStats struct {
	numAccess                uint64
	lastChecked              uint64
	numHit, numMiss, numMSHR uint64
	numSamples               uint64
	averageQueueLength       float64
	missesReturned           uint64

	hitsInPrevEpoch           uint64
	missesInPrevEpoch         uint64
	accessesInPrevEpoch       uint64
	missesReturnedInPrevEpoch uint64
	timesMeasuredInPrevEpoch  uint64
	avgQueueLengthInPrevEpoch float64
	numMSHRStallsInPrevEpoch  uint64
	avgMSHRLenInPrevEpoch     float64
	numStalledInPrevEpoch     uint64

	hitsInCurEpoch           uint64
	missesInCurEpoch         uint64
	missesReturnedInCurEpoch uint64
	accessesInCurEpoch       uint64
	timesMeasuredInCurEpoch  uint64
	avgQueueLengthInCurEpoch float64
	numMSHRStallsInCurEpoch  uint64
	avgMSHRLenInCurEpoch     float64
	numStalledInCurEpoch     uint64

	sendStateInfo bool
	interleaving  uint64

	perEpochWhatIf   [4]uint64
	epochCountWhatIf uint64
}

// A LatTLB is a cache that maintains some page information.
type LatTLB struct {
	*akita.TickingComponent

	TopPort     akita.Port
	BottomPort  akita.Port
	ControlPort akita.Port
	ToRTU       akita.Port

	LowModule       akita.Port
	LowModuleFinder cache.LowModuleFinder
	TLBFinder       cache.LowModuleFinder

	HashFuncHSL func(address uint64) uint64

	log2PageSize   uint64
	log2NumSets    uint64
	setMask        uint64
	numSets        int
	numWays        int
	pageSize       uint64
	numReqPerCycle int
	latency        int
	indexingMask   uint64

	Sets []internal.Set

	pipeline     pipelining.Pipeline
	lookupBuffer util.Buffer

	mshr                mshr
	respondingMSHREntry *mshrEntry

	isPaused bool

	CommandProcessor akita.Port
	stats            TLBStats

	setsAccessed [256]uint64
	mode4Kstrip  bool
	hysterisis   bool
}

// GetPipeline gets the pipeline in the LatTLB
func (tlb *LatTLB) GetPipeline() pipelining.Pipeline {
	return tlb.pipeline
}

// GetPipeline gets the pipeline in the LatTLB
func (tlb *LatTLB) GetTopPort() akita.Port {
	return tlb.TopPort
}

// GetPipeline gets the pipeline in the LatTLB
func (tlb *LatTLB) GetBottomPort() akita.Port {
	return tlb.BottomPort
}

// GetPipeline gets the pipeline in the LatTLB
func (tlb *LatTLB) GetControlPort() akita.Port {
	return tlb.ControlPort
}

// GetNumSets gets the number of sets in the LatTLB
func (tlb *LatTLB) GetNumSets() int {
	return tlb.numSets
}

// GetNumWays gets the number of ways in the LatTLB
func (tlb *LatTLB) GetNumWays() int {
	return tlb.numWays
}

func (tlb *LatTLB) SetLowModuleFinder(lmf cache.LowModuleFinder) {
	tlb.LowModuleFinder = lmf
}

func (tlb *LatTLB) SetCommandProcessor(cp akita.Port) {
	tlb.CommandProcessor = cp
}

// Reset sets all the entries int he LatTLB to be invalid
func (tlb *LatTLB) reset() {
	tlb.Sets = make([]internal.Set, tlb.numSets)
	for i := 0; i < tlb.numSets; i++ {
		set := internal.NewSet(tlb.numWays)
		tlb.Sets[i] = set
	}
}

// Tick defines how LatTLB update states at each cycle
func (tlb *LatTLB) Tick(now akita.VTimeInSec) bool {
	madeProgress := false

	madeProgress = tlb.performCtrlReq(now) || madeProgress

	tlb.doStatsCollection(now)
	// move the call to doStatsCollection to three places where
	// numAccesses is being incremented.

	// improper use of StartTask
	if tlb.GetFrontQueueLength() > 0 {
		tracing.StartTask("", "", now, tlb, "imbalance", "", nil)
	}

	if !tlb.isPaused {
		for i := 0; i < tlb.numReqPerCycle; i++ {
			madeProgress = tlb.respondMSHREntry(now) || madeProgress
		}

		for i := 0; i < tlb.numReqPerCycle; i++ {
			madeProgress = tlb.parseBottom(now) || madeProgress
		}

		for i := 0; i < tlb.numReqPerCycle; i++ {
			madeProgress = tlb.lookup(now) || madeProgress
		}

		madeProgress = tlb.pipeline.Tick(now) || madeProgress

		// pipeline or queue cycling done here

		for i := 0; i < tlb.numReqPerCycle; i++ {
			madeProgress = tlb.parseFromTop(now) || madeProgress
		}
	}
	return madeProgress
}

func (tlb *LatTLB) respondMSHREntry(now akita.VTimeInSec) bool {
	if tlb.respondingMSHREntry == nil {
		return false
	}

	mshrEntry := tlb.respondingMSHREntry
	page := mshrEntry.page
	req := mshrEntry.Requests[0]
	var accessResult device.AccessResult
	if mshrEntry.NumResponded() == 0 {
		accessResult = device.TLBMiss
		tlb.stats.missesReturned++
		if tlb.stats.interleaving == 12 && !tlb.stats.sendStateInfo && tlb.stats.missesReturned == 1300 {
			tlb.stats.sendStateInfo = true
		}
		// tlb.stats.missesReturnedInCurEpoch++
		// fmt.Println(tlb.stats.missesReturnedInCurEpoch)
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
		WithSrcL2TLB(tlb.Name()).
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

func collectCoalescingStat(tlb L2TLB, now akita.VTimeInSec) {
	bufContainer := tlb.GetTopPort().(akita.MsgBufferContainer) //.(*akita.LimitNumMsgPort)
	buf := bufContainer.GetBuffer()
	coalMine := make(map[uint64]int)
	for _, msg := range buf {
		t := msg.(*device.TranslationReq)
		addr := t.VAddr
		coalMine[addr] = coalMine[addr] + 1
	}
	count := 0
	totalcoal := 0
	for addr := range coalMine {
		if coalMine[addr] > 1 {
			count++
			totalcoal += coalMine[addr]
		}
	}
	tracing.StartTask("", "", now, tlb,
		"buflen", strconv.Itoa(len(buf)), nil)
	if len(buf) > 0 {
		tracing.StartTask("", "", now, tlb,
			"bufleng0", strconv.Itoa(len(buf)), nil)
		tracing.StartTask("", "", now, tlb,
			"coalesceaddr", strconv.Itoa(count), nil)
		tracing.StartTask("", "", now, tlb,
			"coalesce", strconv.Itoa(totalcoal), nil)
	}
}

func (tlb *LatTLB) collectMSHROccupancy(now akita.VTimeInSec) {
	m := tlb.mshr
	uniqEntries := len(m.AllEntries())
	totalEntries := 0
	for _, me := range m.AllEntries() {
		totalEntries += len(me.Requests)
	}
	// fmt.Println("MSHR len:", uniqEntries)
	tracing.StartTask("", "", now, tlb,
		"MSHRlen", strconv.Itoa(totalEntries), nil)
	tracing.StartTask("", "", now, tlb,
		"MSHRuniq", strconv.Itoa(uniqEntries), nil)
	if uniqEntries > 0 {
		tracing.StartTask("", "", now, tlb,
			"MSHRlen_g0", strconv.Itoa(totalEntries), nil)
		tracing.StartTask("", "", now, tlb,
			"MSHRuniq_g0", strconv.Itoa(uniqEntries), nil)
	}
}

func (tlb *LatTLB) parseFromTop(now akita.VTimeInSec) bool {
	msg := tlb.TopPort.Peek()
	collectCoalescingStat(tlb, now)
	tlb.collectMSHROccupancy(now)
	if msg == nil {
		return false
	}

	req := msg.(*device.TranslationReq)

	// push to waiting queue
	if tlb.pipeline.CanAccept() {
		/*tlb.stats.sendStateInfo && */
		if now-req.SendTime > 1e-9 {

			tlb.stats.numStalledInCurEpoch++

			tracing.AddTaskStep(
				tracing.MsgIDAtReceiver(req, tlb),
				now, tlb,
				"stalled-l2-tlb-req-count",
			)
		}

		tracing.AddTaskStep(
			tracing.MsgIDAtReceiver(req, tlb),
			now, tlb,
			"l2-tlb-req-count",
		)

		pipelineItem := tlbPipelineItem{
			taskID:         akita.GetIDGenerator().Generate(),
			translationReq: req,
		}
		tlb.pipeline.Accept(now, pipelineItem)
		tlb.TopPort.Retrieve(now)
		tracing.TraceReqReceive(req, now, tlb)
		tracing.StopTracingNetworkReq(req, now, tlb)
		tracing.StartTask(strconv.FormatUint(req.VAddr, 10), "", now, tlb, "entropy", "", nil)
		return true
	}
	return false
}

func balanceCheck(perEpochWhatIf [4]uint64) bool {
	sum := uint64(0)
	for i := 0; i < 4; i++ {
		sum += perEpochWhatIf[i]
	}
	if sum == 0 {
		return false
	}
	for i := 0; i < 4; i++ {
		if (float32(perEpochWhatIf[i]) / float32(sum)) > 0.5 {
			return false
		}
	}
	return true
}

func (tlb *LatTLB) updateReverseSwitchStats(vAddr uint64) {
	if (tlb.stats.epochCountWhatIf % 5000) == 0 {
		if balanceCheck(tlb.stats.perEpochWhatIf) {
			if tlb.hysterisis == true {
				fmt.Println("balanced twice")
			} else {
				tlb.hysterisis = true
				fmt.Println("balanced once")
			}
		} else {
			tlb.hysterisis = false
			fmt.Println("not balanced")
		}
		for i := 0; i < 4; i++ {
			fmt.Println(tlb.Name(), tlb.stats.perEpochWhatIf[i])
			tlb.stats.perEpochWhatIf[i] = 0
		}
	}
	tlb.stats.epochCountWhatIf += 1
	tlb.stats.perEpochWhatIf[tlb.HashFuncHSL(vAddr)] += 1

}

func (tlb *LatTLB) lookup(now akita.VTimeInSec) bool {
	// pop from waiting queue
	item := tlb.lookupBuffer.Peek()
	if item == nil {
		return false
	}
	pipelineItem := item.(tlbPipelineItem)
	req := pipelineItem.translationReq
	// if req.VAddr == 49152 {
	// fmt.Println(req.VAddr)
	// panic("oh no! dbchbbhw")
	// }
	mshrEntry := tlb.mshr.Query(req.PID, req.VAddr)
	if mshrEntry != nil {
		ok := tlb.processTLBMSHRHit(now, mshrEntry, req)
		if ok {
			tracing.AddTaskStep(
				tracing.MsgIDAtReceiver(req, tlb),
				now, tlb,
				"tlb-mshr-hit",
			)
			tlb.lookupBuffer.Pop()
			// if tlb.stats.sendStateInfo {
			tlb.stats.numAccess += 1
			tlb.stats.accessesInCurEpoch++
			tlb.stats.numMSHR += 1
			return true
		}
		return false
	}

	setID := tlb.vAddrToSetID(req.VAddr)
	// setID := tlb.vAddrToSetIDxor7(req.VAddr)
	set := tlb.Sets[setID]
	wayID, page, found := set.Lookup(req.PID, req.VAddr)
	if found && page.Valid {
		return tlb.handleTranslationHit(now, req, setID, wayID, page)
	}

	// passing setID for tracing purposes
	// TODO: use the return value of the following function to avoid passing setID
	return tlb.handleTranslationMiss(now, req, setID)
}

func (tlb *LatTLB) handleTranslationHit(
	now akita.VTimeInSec,
	req *device.TranslationReq,
	setID, wayID int,
	page device.Page,
) bool {
	ok := tlb.sendRspToTop(now, req, page)
	if !ok {
		return false
	}
	tlb.visit(setID, wayID)
	tlb.lookupBuffer.Pop()
	// if tlb.stats.sendStateInfo {
	tlb.stats.numAccess += 1
	tlb.stats.accessesInCurEpoch++
	tlb.stats.numHit += 1
	tlb.stats.hitsInCurEpoch++
	// }
	// fmt.Println("hits:", tlb.stats.hitsInCurEpoch)

	tracing.AddTaskStep(
		tracing.MsgIDAtReceiver(req, tlb),
		now, tlb,
		"tlb-hit",
	)
	tracing.TraceReqComplete(req, now, tlb)

	return true
}

func (tlb *LatTLB) handleTranslationMiss(
	now akita.VTimeInSec,
	req *device.TranslationReq,
	setID int,
) bool {
	if tlb.doPageWalk(req.VAddr) {

		if tlb.mshr.IsFull() {
			tracing.StartTask(tlb.Name()+"stall", "", now, tlb, "mshr_stall", "", nil)
			tlb.stats.numMSHRStallsInCurEpoch++
			return false
		}

		fetched := tlb.fetchBottom(now, req)
		if fetched {
			//tlb.TopPort.Retrieve(now)
			tlb.lookupBuffer.Pop()
			// if tlb.stats.sendStateInfo {
			tlb.stats.numAccess += 1
			tlb.stats.accessesInCurEpoch++
			tlb.stats.numMiss += 1
			tlb.stats.missesInCurEpoch++
			// }
			// tracing.TraceReqReceive(req, now, tlb)
			tracing.AddTaskStep(
				tracing.MsgIDAtReceiver(req, tlb),
				now, tlb,
				"tlb-miss",
			)

			// this is the missepoint
			// add a set miss tracer here.
			tracing.StartTask("", "", now, tlb, "set_miss_tracing",
				strconv.FormatUint(uint64(setID), 10), nil)

			return true
		}

		return false
	} else {
		req.Dst = tlb.ToRTU
		req.PutInL2TLBBuffer = true
		err := tlb.TopPort.Send(req)
		if err != nil {
			return false
		}
		tracing.AddTaskStep(
			tracing.MsgIDAtReceiver(req, tlb),
			now, tlb,
			"tlb-redirect",
		)
		tlb.lookupBuffer.Pop()
		return true
	}
}

func (tlb *LatTLB) doPageWalk(vAddr uint64) (doWalk bool) {
	return tlb.Name() == tlb.TLBFinder.Find(vAddr).Name()[0:21]
}

// we are assuimg 4K pages here. TODO: FIX THIS!
func (tlb *LatTLB) vAddrToSetIDxor7(vAddr uint64) (setID int) {
	index := uint64(0)
	sp1 := (vAddr >> 12) & 0x7f
	sp2 := (vAddr >> 19) & 0x7f
	sp3 := (vAddr >> 26) & 0x7f
	sp4 := (vAddr >> 33) & 0x7f
	index = uint64(sp1 ^ sp2 ^ sp3 ^ sp4)
	return int(index)
}

func (tlb *LatTLB) vAddrToSetID(vAddr uint64) (setID int) {
	index := uint64(0)
	vpn := vAddr >> tlb.log2PageSize
	// usefulBits := (uint64(1) << 28) - 1
	for i := 0; i < 4; i++ {
		index ^= (vpn & tlb.setMask)
		vpn >>= tlb.log2NumSets
	}
	setID = int(index)
	tlb.setsAccessed[setID]++
	return
}

func (tlb *LatTLB) sendRspToTop(
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

func (tlb *LatTLB) processTLBMSHRHit(
	now akita.VTimeInSec,
	mshrEntry *mshrEntry,
	req *device.TranslationReq,
) bool {
	// if len(mshrEntry.Requests) < 8 {
	mshrEntry.Requests = append(mshrEntry.Requests, req)
	return true
	// }
	// return false
}

func (tlb *LatTLB) fetchBottom(now akita.VTimeInSec, req *device.TranslationReq) bool {
	dstPort := tlb.LowModuleFinder.Find(req.VAddr)

	fetchBottom := device.TranslationReqBuilder{}.
		WithSendTime(now).
		WithSrc(tlb.BottomPort).
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

	tlb.stats.avgMSHRLenInCurEpoch = (tlb.stats.avgMSHRLenInCurEpoch*float64(tlb.stats.timesMeasuredInCurEpoch) + float64(len(tlb.mshr.AllEntries()))) / float64(tlb.stats.timesMeasuredInCurEpoch+1)
	tlb.stats.timesMeasuredInCurEpoch++

	tracing.TraceReqInitiate(fetchBottom, now, tlb,
		tracing.MsgIDAtReceiver(req, tlb))

	return true
}

func (tlb *LatTLB) parseBottom(now akita.VTimeInSec) bool {
	if tlb.respondingMSHREntry != nil {
		return false
	}

	item := tlb.BottomPort.Peek()
	if item == nil {
		return false
	}

	rsp := item.(*device.TranslationRsp)
	page := rsp.Page

	mshrEntryPresent := tlb.mshr.IsEntryPresent(rsp.Page.PID, rsp.Page.VAddr)
	if !mshrEntryPresent {
		panic("oh no!")
		tlb.BottomPort.Retrieve(now)
		tracing.TraceReqFinalize(rsp, now, tlb)
		return true
	}

	setID := tlb.vAddrToSetID(page.VAddr)
	// setID := tlb.vAddrToSetIDxor7(page.VAddr)
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

	tlb.stats.avgMSHRLenInCurEpoch = (tlb.stats.avgMSHRLenInCurEpoch*float64(tlb.stats.timesMeasuredInCurEpoch) + float64(len(tlb.mshr.AllEntries()))) / float64(tlb.stats.timesMeasuredInCurEpoch+1)
	tlb.stats.timesMeasuredInCurEpoch++

	tlb.BottomPort.Retrieve(now)
	tracing.TraceReqFinalize(mshrEntry.reqToBottom, now, tlb)

	tracing.EndTask(tlb.Name()+"stall", now, tlb)
	return true
}

func (tlb *LatTLB) performCtrlReq(now akita.VTimeInSec) bool {
	item := tlb.ControlPort.Peek()
	if item == nil {
		return false
	}

	var madeProgress bool
	switch req := item.(type) {
	case *TLBFlushReq:
		madeProgress = tlb.handleTLBFlush(now, req)
	case *TLBRestartReq:
		madeProgress = tlb.handleTLBRestart(now, req)
	case *akita.TLBIndexingSwitchMsg:
		madeProgress = tlb.switchIndexing(now, req)
	case *akita.SendStatsMsg:
		madeProgress = tlb.sendCollectedStatsToCP(now)
	default:
		log.Panicf("cannot process request %s", reflect.TypeOf(req))
	}
	if madeProgress {
		tlb.ControlPort.Retrieve(now)
		return true
	}
	return false
}

func (tlb *LatTLB) visit(setID, wayID int) int {
	set := tlb.Sets[setID]
	mruPosition := set.Visit(wayID)
	return mruPosition
}

func (tlb *LatTLB) handleTLBFlush(now akita.VTimeInSec, req *TLBFlushReq) bool {
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
		// setID := tlb.vAddrToSetID(vAddr)
		setID := tlb.vAddrToSetIDxor7(vAddr)
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

func (tlb *LatTLB) handleTLBRestart(now akita.VTimeInSec, req *TLBRestartReq) bool {
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

func (tlb *LatTLB) GetFrontQueueLength() int {
	bufContainer := tlb.GetTopPort().(akita.MsgBufferContainer) //.(*akita.LimitNumMsgPort)
	buf := bufContainer.GetBuffer()
	return len(buf)
}

func (tlb *LatTLB) switchIndexing(now akita.VTimeInSec,
	req *akita.TLBIndexingSwitchMsg) bool {

	tlb.stats.numAccess = 0
	tlb.stats.lastChecked = 0
	tlb.stats.hitsInPrevEpoch = 0
	tlb.stats.missesInPrevEpoch = 0
	tlb.stats.accessesInPrevEpoch = 0
	tlb.stats.missesReturnedInPrevEpoch = 0
	tlb.stats.timesMeasuredInPrevEpoch = 0
	tlb.stats.avgQueueLengthInPrevEpoch = 0
	tlb.stats.numMSHRStallsInPrevEpoch = 0
	tlb.stats.avgMSHRLenInPrevEpoch = 0
	tlb.stats.numStalledInPrevEpoch = 0

	tlb.stats.hitsInCurEpoch = 0
	tlb.stats.missesInCurEpoch = 0
	tlb.stats.missesReturnedInCurEpoch = 0
	tlb.stats.accessesInCurEpoch = 0
	tlb.stats.timesMeasuredInCurEpoch = 0
	tlb.stats.avgQueueLengthInCurEpoch = 0
	tlb.stats.numMSHRStallsInCurEpoch = 0
	tlb.stats.avgMSHRLenInCurEpoch = 0
	tlb.stats.numStalledInCurEpoch = 0

	tlb.stats.missesReturned = 0
	if req.TLBInterleaving == 12 {
		tlb.stats.sendStateInfo = false
	} else {
		tlb.stats.sendStateInfo = true
	}
	tlb.stats.interleaving = req.TLBInterleaving

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

	lmf := tlb.TLBFinder.(*cache.CustomTwoLevelLowModuleFinder)
	if req.TLBIndexingSwitch == akita.TLBIndexingSwitch4K {
		tlb.mode4Kstrip = true
		lmf.Hashfunc = HashFunc4KXOR
		fmt.Println(tlb.Name(), "L2 TLB ", tlb.Name(), "switching to 4K sir", now, req.TLBInterleaving)
	} else {
		tlb.mode4Kstrip = false
		lmf.Hashfunc = HashFuncHSL
		fmt.Println(tlb.Name(), "L2 TLB ", tlb.Name(), "switching to HSL sir", now, req.TLBInterleaving)
	}
	return true
}

func div(x, y float64) float64 {
	if y == 0 {
		return 0
	}
	return x / y
}

func (tlb *LatTLB) sendCollectedStatsToCP(now akita.VTimeInSec) bool {

	// TODO: set hits and misses appropriately
	req := akita.CollectedStatsMsgBuilder{}.
		WithSendTime(now).
		WithSrc(tlb.ControlPort).
		WithDst(tlb.CommandProcessor).
		WithNumMisses(tlb.stats.missesInCurEpoch + tlb.stats.missesInPrevEpoch).
		WithNumHits(tlb.stats.hitsInCurEpoch + tlb.stats.hitsInPrevEpoch).
		Build()
	err := tlb.ControlPort.Send(req)
	if err != nil {
		panic("could not send stats to CP!")
	}
	return true
}

func (tlb *LatTLB) doStatsCollection(now akita.VTimeInSec) bool {
	currentQueueLength := float64(tlb.GetFrontQueueLength())
	tlb.stats.avgQueueLengthInCurEpoch = (float64(tlb.stats.timesMeasuredInCurEpoch)*tlb.stats.avgQueueLengthInCurEpoch + currentQueueLength) / float64(tlb.stats.timesMeasuredInCurEpoch+1)
	if tlb.stats.accessesInCurEpoch > 5000 {
		tlb.stats.accessesInPrevEpoch = tlb.stats.accessesInCurEpoch
		tlb.stats.hitsInPrevEpoch = tlb.stats.hitsInCurEpoch
		tlb.stats.missesInPrevEpoch = tlb.stats.missesInCurEpoch
		tlb.stats.missesReturnedInPrevEpoch = tlb.stats.missesReturnedInCurEpoch
		tlb.stats.avgQueueLengthInPrevEpoch = tlb.stats.avgQueueLengthInCurEpoch
		tlb.stats.timesMeasuredInPrevEpoch = tlb.stats.timesMeasuredInCurEpoch
		tlb.stats.numMSHRStallsInPrevEpoch = tlb.stats.numMSHRStallsInCurEpoch
		tlb.stats.avgMSHRLenInPrevEpoch = tlb.stats.avgMSHRLenInCurEpoch
		tlb.stats.numStalledInPrevEpoch = tlb.stats.numStalledInCurEpoch
		tlb.stats.hitsInCurEpoch = 0
		tlb.stats.accessesInCurEpoch = 0
		tlb.stats.missesInCurEpoch = 0
		tlb.stats.missesReturnedInCurEpoch = 0
		tlb.stats.avgQueueLengthInCurEpoch = 0
		tlb.stats.timesMeasuredInCurEpoch = 0
		tlb.stats.numMSHRStallsInCurEpoch = 0
		tlb.stats.avgMSHRLenInCurEpoch = 0
		tlb.stats.numStalledInCurEpoch = 0
	}
	return true
}
