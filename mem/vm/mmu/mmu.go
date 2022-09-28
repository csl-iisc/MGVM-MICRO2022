package mmu

import (
	"encoding/binary"
	"fmt"
	"log"
	"reflect"
	"strconv"

	"strings"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem"
	"gitlab.com/akita/mem/cache"
	"gitlab.com/akita/mem/device"
	"gitlab.com/akita/util/akitaext"
	"gitlab.com/akita/util/tracing"
)

type transactionState int

const (
	newTransaction transactionState = iota
	sentToPageWalkCache
	pageWalkCacheDone
	sentToMem
	memDone
	transactionFinished
)

type transaction struct {
	req  *device.TranslationReq
	page device.Page
	//cycleLeft int
	//migration *device.PageMigrationReqToDriver
	level             int
	msgID             string
	state             transactionState
	PPN               uint64
	vAddr             uint64
	remoteMemAccesses int
}

// MMUImpl is the default mmu implementation. It is also an akita Component.
type MMUImpl struct {
	akita.TickingComponent

	ToTop            akita.Port
	ControlPort      akita.Port
	CommandProcessor akita.Port

	pageWalkCachePort akita.Port
	PageWalkCache     akita.Port
	topSender         akitaext.BufferedSender

	TranslationPort akita.Port
	lowModuleFinder cache.LowModuleFinder
	numChiplets     uint64

	pageTable *device.PageTableImpl
	//	latency             int
	maxRequestsInFlight int

	walkingTranslations []transaction

	remoteMemAccessesInCurEpoch  uint64
	avgWalksEnqueuedInCurEpoch   float64
	remoteMemAccessesInPrevEpoch uint64
	avgWalksEnqueuedInPrevEpoch  float64
	memAccessesInCurEpoch        uint64
	memAccessesInPrevEpoch       uint64

	numWalksDone        uint64
	numWalksInCurEpoch  uint64
	numWalksInPrevEpoch uint64
	lastChecked         uint64
	sendStateInfo       bool
	// numWalksArrived          uint64
	interleaving uint64
}

// Tick defines how the MMU update state each cycle
func (mmu *MMUImpl) Tick(now akita.VTimeInSec) bool {
	madeProgress := false

	madeProgress = mmu.performCtrlReq(now) || madeProgress

	if mmu.GetNumActiveWalkers() > 0 {
		tracing.StartTask("", "", now, mmu, "imbalance", "", nil)
	}

	madeProgress = mmu.topSender.Tick(now) || madeProgress
	madeProgress = mmu.walkPageTable(now) || madeProgress
	madeProgress = mmu.parseFromPageWalkCache(now) || madeProgress
	madeProgress = mmu.parseFromMem(now) || madeProgress
	madeProgress = mmu.parseFromTop(now) || madeProgress
	return true
	// return madeProgress
}

func (mmu *MMUImpl) performCtrlReq(now akita.VTimeInSec) bool {
	item := mmu.ControlPort.Peek()
	if item == nil {
		return false
	}
	madeProgress := false
	switch req := item.(type) {
	case *akita.TLBIndexingSwitchMsg:
		madeProgress = mmu.switchIndexing(now, req)
	case *akita.SendStatsMsg:
		madeProgress = mmu.sendCollectedStatsToCP(now)
	default:
		log.Panicf("cannot process request %s", reflect.TypeOf(req))
	}
	if madeProgress {
		mmu.ControlPort.Retrieve(now)
		return true
	}
	panic("something is wrong!")
}

func div(x, y float64) float64 {
	if y == 0 {
		return 0
	}
	return x / y
}

func (mmu *MMUImpl) sendCollectedStatsToCP(now akita.VTimeInSec) bool {
	req := akita.CollectedStatsMsgBuilder{}.
		WithSendTime(now).
		WithSrc(mmu.ControlPort).
		WithDst(mmu.CommandProcessor).
		Build()
	err := mmu.ControlPort.Send(req)
	if err != nil {
		panic("could not send stats to CP!")
	}
	mmu.ControlPort.Retrieve(now)
	return true
}

func (mmu *MMUImpl) trace(now akita.VTimeInSec, what string) {
	ctx := akita.HookCtx{
		Domain: mmu,
		Now:    now,
		Item:   what,
	}

	mmu.InvokeHook(ctx)
}

func (mmu *MMUImpl) walkPageTable(now akita.VTimeInSec) bool {
	numActiveTransactions := len(mmu.walkingTranslations)
	tracing.StartTask("", "", now, mmu, "num_active_walkers", fmt.Sprintf("%d", numActiveTransactions), nil)
	madeProgress := numActiveTransactions > 0
	tmp := mmu.walkingTranslations[:0]
	for i := 0; i < numActiveTransactions; i++ {
		trans := &mmu.walkingTranslations[i]
		switch trans.state {
		case newTransaction:
			mmu.sendToPageWalkCache(now, i)
		case pageWalkCacheDone:
			mmu.sendToMem(now, i)
		case memDone:
			mmu.sendToMem(now, i)
		}
		if trans.state != transactionFinished {
			tmp = append(tmp, *trans)
		} else {
			mmu.numWalksDone++
			if mmu.numWalksDone == 1200 && !mmu.sendStateInfo && mmu.interleaving == 12 {
				mmu.numWalksDone = 0
				mmu.lastChecked = 0
				mmu.sendStateInfo = true
			}
			mmu.avgWalksEnqueuedInCurEpoch = (mmu.avgWalksEnqueuedInCurEpoch*float64(mmu.numWalksInCurEpoch) + float64(len(mmu.ToTop.(akita.MsgBufferContainer).GetBuffer()))) / float64(mmu.numWalksInCurEpoch+1)
			mmu.numWalksInCurEpoch++

		}
	}
	mmu.walkingTranslations = tmp
	return madeProgress
}

func (mmu *MMUImpl) sendMsgToCP(now akita.VTimeInSec) bool {
	if !mmu.sendStateInfo {
		panic("how are we sending a message to CP when we shouldn't be?!")
	}
	req := akita.SendStatsMsgBuilder{}.
		WithSendTime(now).
		WithSrc(mmu.ControlPort).
		WithDst(mmu.CommandProcessor).
		Build()

	err := mmu.ControlPort.Send(req)
	if err != nil {
		fmt.Println(err)
		// panic("oh no")
	}
	return true
}

func (mmu *MMUImpl) switchIndexing(now akita.VTimeInSec,
	req *akita.TLBIndexingSwitchMsg) bool {

	mmu.numWalksDone = 0
	mmu.lastChecked = 0

	mmu.avgWalksEnqueuedInCurEpoch = 0
	mmu.avgWalksEnqueuedInPrevEpoch = 0

	mmu.remoteMemAccessesInCurEpoch = 0
	mmu.remoteMemAccessesInPrevEpoch = 0

	mmu.memAccessesInCurEpoch = 0
	mmu.memAccessesInPrevEpoch = 0

	mmu.numWalksInCurEpoch = 0
	mmu.numWalksInPrevEpoch = 0

	if req.TLBInterleaving == 12 {
		mmu.sendStateInfo = false
	} else {
		mmu.sendStateInfo = true
	}
	mmu.interleaving = req.TLBInterleaving
	// fmt.Println(mmu.Name(), "flushing stats on switch", now, req.TLBInterleaving)
	return true
}

func (mmu *MMUImpl) parseFromPageWalkCache(now akita.VTimeInSec) bool {
	madeProgress := false
	item := mmu.pageWalkCachePort.Peek()
	if item != nil {
		switch msg := item.(type) {
		case *mem.DataReadyRsp:
			mmu.handlePageWalkCacheResponse(msg, now)
		case *mem.WriteDoneRsp:
		default:
			panic("unknown message type")
		}
		madeProgress = true
	}
	mmu.pageWalkCachePort.Retrieve(now)
	return madeProgress
}

func (mmu *MMUImpl) parseFromMem(now akita.VTimeInSec) bool {
	madeProgress := false
	item := mmu.TranslationPort.Peek()
	if item != nil {
		switch msg := item.(type) {
		case *mem.DataReadyRsp:
			mmu.handleMemResponse(msg, now)
		default:
			panic("unknown message type")
		}
		madeProgress = true
	}
	mmu.TranslationPort.Retrieve(now)
	return madeProgress
}

func (mmu *MMUImpl) sendToPageWalkCache(now akita.VTimeInSec, i int) {
	trans := &mmu.walkingTranslations[i]
	transState := trans.state
	if transState != newTransaction {
		panic("this state shouldn't be here!")
	}
	readReq := mem.ReadReqBuilder{}.
		WithSendTime(now).
		WithSrc(mmu.pageWalkCachePort).
		WithDst(mmu.PageWalkCache).
		WithPID(trans.req.PID).
		WithAddress(mmu.pageTable.AlignToPage(trans.req.VAddr)).
		WithByteSize(8).
		Build()
	// fmt.Println(readReq.Info)
	trans.msgID = readReq.ID
	mmu.pageWalkCachePort.Send(readReq)
	trans.state = sentToPageWalkCache
}

func (mmu *MMUImpl) sendToMem(now akita.VTimeInSec, i int) {
	trans := &mmu.walkingTranslations[i]
	transState := trans.state
	if transState != pageWalkCacheDone && transState != memDone {
		panic("this state shouldn't be here!")
	}
	PPN := trans.PPN
	PPNWithOffset := mmu.pageTable.AddOffset(PPN, trans.vAddr)
	// HUGEPAGES
	// h := fmt.Sprintf("%x", trans.vAddr)
	// h1 := fmt.Sprintf("%x", PPN)
	// h2 := fmt.Sprintf("%x", PPNWithOffset)
	// fmt.Println(h, h1, h2)
	// fmt.Println("PPNWithOffset:", PPN, trans.vAddr, PPNWithOffset)

	trans.vAddr = mmu.pageTable.NextLevel(trans.vAddr)
	// HUGEPAGES
	srcPort := mmu.TranslationPort
	readReqInfo := &mem.ReadReqInfo{ReturnAccessInfo: true}
	dstPort := mmu.lowModuleFinder.Find(PPNWithOffset)
	// fmt.Println(fmt.Sprintf("%x", PPNWithOffset), dstPort.Name())
	// if strings.Contains(dstPort.Name(), "RDMA") {
	// trans.remoteMemAccesses++
	// }
	readReq := mem.ReadReqBuilder{}.
		WithSendTime(now).
		WithSrc(srcPort).
		WithDst(dstPort).
		WithPID(trans.req.PID).
		WithAddress(PPNWithOffset).
		WithByteSize(8).
		WithInfo(readReqInfo).
		Build()
	trans.msgID = readReq.ID
	// TODO: fix this with err
	srcPort.Send(readReq)
	tracing.StartTask(
		readReq.Meta().ID+"MMU-mem-latency",
		tracing.MsgIDAtReceiver(readReq, mmu),
		now,
		mmu,
		"MMU_mem_latency",
		reflect.TypeOf(readReq).String(),
		readReq,
	)
	trans.state = sentToMem
	// TODO determine if remote or local by string comparison and trace
	// This is NOT production worthy code.
	// Repeat Warning!!! This is ugly as it can get.
	dstsplits := strings.Split(readReq.Dst.Name(), ".")
	dstType := dstsplits[2]
	// mmusplits := strings.Split(mmu.Name(), ".")
	// mmuChiplet := mmusplits[1]
	if dstType == "ChipRDMA" {
		mmu.remoteMemAccessesInCurEpoch++
		tracing.AddTaskStep(tracing.MsgIDAtReceiver(trans.req, mmu),
			now, mmu, "page_walk_req_remote")
		if trans.level == 3 {
			tracing.AddTaskStep(tracing.MsgIDAtReceiver(trans.req, mmu),
				now, mmu, "pw-level-3-remote-reqs")
		}
	} else {
		tracing.AddTaskStep(tracing.MsgIDAtReceiver(trans.req, mmu),
			now, mmu, "page_walk_req_local")
	}
	mmu.memAccessesInCurEpoch++
}

func (mmu *MMUImpl) handlePageWalkCacheResponse(rsp *mem.DataReadyRsp, now akita.VTimeInSec) {
	for i := 0; i < len(mmu.walkingTranslations); i++ {
		trans := &mmu.walkingTranslations[i]
		if trans.msgID == rsp.RespondTo {
			if rsp.Data != nil {
				rspData := binary.LittleEndian.Uint64(rsp.Data)
				trans.PPN = rspData & ^uint64(3)
				level := int(rspData & uint64(3))
				trans.vAddr = mmu.pageTable.MoveToLevel(trans.vAddr, level+1)
				trans.level = level + 1
				// h := fmt.Sprintf("%x", trans.req.VAddr)
				// h1 := fmt.Sprintf("%x", trans.PPN)
				// h2 := fmt.Sprintf("%x", trans.vAddr)
				// fmt.Println("response from pwc:", h, h1, level, h2)
			}
			trans.state = pageWalkCacheDone
			//trace point
			tracing.AddTaskStep(tracing.MsgIDAtReceiver(trans.req, mmu),
				now, mmu, "pwc-hit-level"+strconv.Itoa(trans.level))
		}
	}
}

func getAccessResultString(accessResult mem.AccessResult) (str string) {
	switch accessResult {
	case mem.ReadHit:
		str = "read-hit"
		// fmt.Println("read hit")
	case mem.ReadMiss:
		str = "read-miss"
		// fmt.Println("read miss")
	case mem.ReadMSHRHit:
		str = "read-mshr-hit"
		// fmt.Println("read mshr hit")
	default:
		panic("unknown access type")
	}
	return
}

func getChipletNum(component string) (chipletNum string) {
	chipletNum = strings.Split(component, "_")[1][1:2]
	return
}

func (mmu *MMUImpl) handleMemResponse(rsp *mem.DataReadyRsp, now akita.VTimeInSec) {
	for i := 0; i < len(mmu.walkingTranslations); i++ {
		trans := &mmu.walkingTranslations[i]
		if trans.msgID == rsp.RespondTo {
			rspInfo := rsp.Info.(*mem.DataReadyRspInfo)
			accessResult := rspInfo.AccessResult
			src := rspInfo.Src
			taskStep := fmt.Sprintf("chiplet-%s-level-%d-%s", getChipletNum(src), trans.level, getAccessResultString(accessResult))
			// fmt.Println(taskStep)
			tracing.AddTaskStep(
				trans.msgID+"MMU-mem-latency",
				now, mmu,
				taskStep,
			)
			tracing.EndTask(
				trans.msgID+"MMU-mem-latency",
				now,
				mmu,
			)

			trans.PPN = binary.LittleEndian.Uint64(rsp.Data)
			// fmt.Println("response from mem:", trans.PPN, binary.BigEndian.Uint64(rsp.Data))
			trans.state = memDone
			if trans.level+1 == 4 {
				mmu.finalizeTransaction(now, i)
			} else {
				mmu.fillPageWalkCache(now, i)
			}
			trans.level++
		}
	}
}

func (mmu *MMUImpl) fillPageWalkCache(now akita.VTimeInSec, i int) bool {
	trans := &mmu.walkingTranslations[i]
	level := uint64(trans.level)
	// h := fmt.Sprintf("%x", trans.req.VAddr)
	// h1 := fmt.Sprintf("%x", trans.PPN)
	// fmt.Println("filling the pwc", h, h1, level)
	data := uint64ToBytes(trans.PPN | level)
	writeReq := mem.WriteReqBuilder{}.
		WithSendTime(now).
		WithSrc(mmu.pageWalkCachePort).
		WithDst(mmu.PageWalkCache).
		WithPID(trans.req.PID).
		WithAddress(mmu.pageTable.AlignToPage(trans.req.VAddr) | level). //HUGEPAGES
		WithData(data).
		Build()
	trans.msgID = writeReq.ID
	mmu.pageWalkCachePort.Send(writeReq)
	return true
}

func uint64ToBytes(data uint64) []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, data)
	return bytes
}

func (mmu *MMUImpl) finalizeTransaction(
	now akita.VTimeInSec,
	walkingIndex int,
) bool {
	req := mmu.walkingTranslations[walkingIndex].req
	page, found := mmu.pageTable.Find(req.PID, req.VAddr)
	// fmt.Println(req.VAddr, mmu.walkingTranslations[walkingIndex].PPN)
	if !found {
		panic("page not found")
	}
	pAddr := mmu.walkingTranslations[walkingIndex].PPN
	// fmt.Println("translation:", page.VAddr, pAddr)
	if pAddr != page.PAddr {
		panic("addresses don't match!")
	}
	newPage := device.Page{PID: req.PID, VAddr: req.VAddr, PAddr: pAddr, Valid: true}
	mmu.walkingTranslations[walkingIndex].page = newPage
	mmu.walkingTranslations[walkingIndex].state = transactionFinished
	return mmu.doPageWalkHit(now, walkingIndex)
}

func (mmu *MMUImpl) doPageWalkHit(
	now akita.VTimeInSec,
	walkingIndex int,

) bool {
	if !mmu.topSender.CanSend(1) {
		return false
	}
	walking := mmu.walkingTranslations[walkingIndex]

	rsp := device.TranslationRspBuilder{}.
		WithSendTime(now).
		WithSrc(mmu.ToTop).
		WithDst(walking.req.Src).
		WithRspTo(walking.req.ID).
		WithPage(walking.page).
		Build()

	mmu.topSender.Send(rsp)
	// mmu.toRemoveFromPTW = append(mmu.toRemoveFromPTW, walkingIndex)
	// tracing.StartTracingNetworkReq(rsp, now, mmu, walking.req)
	tracing.TraceReqComplete(walking.req, now, mmu)

	return true
}

func (mmu *MMUImpl) sendTranlationRsp(
	now akita.VTimeInSec,
	trans transaction,
) (madeProgress bool) {
	req := trans.req
	page := trans.page

	rsp := device.TranslationRspBuilder{}.
		WithSendTime(now).
		WithSrc(mmu.ToTop).
		WithDst(req.Src).
		WithRspTo(req.ID).
		WithPage(page).
		Build()
	mmu.topSender.Send(rsp)

	return true
}

func (mmu *MMUImpl) parseFromTop(now akita.VTimeInSec) bool {
	if len(mmu.walkingTranslations) >= mmu.maxRequestsInFlight {
		return false
	}

	req := mmu.ToTop.Retrieve(now)
	if req == nil {
		return false
	}
	// tracing.StopTracingNetworkReq(req, now, mmu)
	switch req := req.(type) {
	case *device.TranslationReq:
		mmu.startWalking(req, now)
	default:
		log.Panicf("MMU canot handle request of type %s", reflect.TypeOf(req))
	}
	return true
}

func (mmu *MMUImpl) startWalking(req *device.TranslationReq, now akita.VTimeInSec) {
	// fmt.Println(fmt.Sprintf("%x", req.VAddr))
	rearrangedVAddr := mmu.pageTable.Rearrange(req.VAddr)
	root := mmu.pageTable.GetRoot(req.PID)
	translationInPipeline := transaction{
		req: req,
		//cycleLeft: mmu.latency,
		level: 0,
		msgID: "invalid",
		state: newTransaction,
		vAddr: rearrangedVAddr,
		PPN:   root,
	}

	mmu.walkingTranslations = append(mmu.walkingTranslations, translationInPipeline)

	tracing.TraceReqReceive(req, now, mmu)
}

//SetLowModuleFinder sets the table recording where to find an address.
func (mmu *MMUImpl) SetLowModuleFinder(lmf cache.LowModuleFinder) {
	mmu.lowModuleFinder = lmf
}

func unique(intSlice []uint64) []uint64 {
	keys := make(map[int]bool)
	list := []uint64{}
	for _, entry := range intSlice {
		if _, value := keys[int(entry)]; !value {
			keys[int(entry)] = true
			list = append(list, entry)
		}
	}
	return list
}

func (mmu *MMUImpl) GetNumActiveWalkers() int {
	return len(mmu.walkingTranslations)
}
