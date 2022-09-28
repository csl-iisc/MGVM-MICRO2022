// Package remotetranslation provides implementation of a remote translation
// unit.
package remotetranslation

import (
	"fmt"
	"log"

	//	"math"
	"reflect"
	"strconv"
	"strings"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem/cache"
	"gitlab.com/akita/mem/device"
	"gitlab.com/akita/util/tracing"
)

type transaction struct {
	fromInside  akita.Msg
	fromOutside akita.Msg
	toInside    akita.Msg
	toOutside   akita.Msg
}

type RemoteTranslationUnit interface {
	tracing.NamedHookable
	SetResponsePorts([]akita.Port)
	SetRemoteAddressTranslationTable(cache.LowModuleFinder)
	GetRequestPort() akita.Port
	GetResponsePort() akita.Port
	GetL1Port() akita.Port
	GetL2Port() akita.Port
	SetLocalModuleFinder(cache.LowModuleFinder)

	GetControlPort() akita.Port
	SetCommandProcessor(akita.Port)
}

// A DefaultRTU helps construct distributed TLBs.
type DefaultRTU struct {
	*akita.TickingComponent

	RequestPort  akita.Port
	ResponsePort akita.Port

	ToL1       akita.Port
	ToL2       akita.Port
	L2CtrlPort akita.Port

	ResponsePorts []akita.Port

	localModules                  cache.LowModuleFinder
	RemoteAddressTranslationTable cache.LowModuleFinder

	transactionsFromOutside []transaction
	transactionsFromInside  []transaction

	numTransPerCycle int

	ControlPort      akita.Port
	CommandProcessor akita.Port

	NumOutgoing          uint64
	NumIncoming          uint64
	NumOutgoingPrevEpoch uint64
	NumIncomingPrevEpoch uint64

	NumReqsServiced uint64
	LastChecked     uint64

	PrevEpochImbalance bool
}

func (rtu *DefaultRTU) GetL1Port() akita.Port {
	return rtu.ToL1
}

func (rtu *DefaultRTU) GetL2Port() akita.Port {
	return rtu.ToL2
}

func (rtu *DefaultRTU) GetRequestPort() akita.Port {
	return rtu.RequestPort
}

func (rtu *DefaultRTU) GetResponsePort() akita.Port {
	return rtu.ResponsePort
}

func (rtu *DefaultRTU) GetControlPort() akita.Port {
	return rtu.ControlPort
}

func (rtu *DefaultRTU) SetResponsePorts(ports []akita.Port) {
	rtu.ResponsePorts = ports
}

func (rtu *DefaultRTU) SetRemoteAddressTranslationTable(f cache.LowModuleFinder) {
	rtu.RemoteAddressTranslationTable = f
}

// SetLocalModuleFinder sets the table to lookup for local translations.
func (rtu *DefaultRTU) SetLocalModuleFinder(lmf cache.LowModuleFinder) {
	rtu.localModules = lmf
}

func (rtu *DefaultRTU) SetCommandProcessor(cp akita.Port) {
	rtu.CommandProcessor = cp
}

func (rtu *DefaultRTU) Tick(now akita.VTimeInSec) bool {

	madeProgress := false

	madeProgress = rtu.processFromCtrlPort(now) || madeProgress
	rtu.consolidateStats(now)
	madeProgress = rtu.processFromL1(now) || madeProgress
	madeProgress = rtu.processFromL2(now) || madeProgress
	madeProgress = rtu.processReqFromOutside(now) || madeProgress
	madeProgress = rtu.processRspFromOutside(now) || madeProgress
	// return true
	return madeProgress
}

// l
func getImbalance(incoming uint64, outgoing uint64) bool {
	imbalanced := false
	if incoming > 2*outgoing {
		imbalanced = true
	}

	return imbalanced
}

// Every epoch of 1000 serviced requests.
func (rtu *DefaultRTU) consolidateStats(now akita.VTimeInSec) {
	if (rtu.NumReqsServiced - rtu.LastChecked) >= 5000 {
		imbalanced := getImbalance(rtu.NumIncoming, rtu.NumOutgoing)
		imbalanced = true
		if imbalanced {
			if rtu.PrevEpochImbalance == true {
				req := akita.TriggerMsgBuilder{}.
					WithSendTime(now).
					WithSrc(rtu.ControlPort).
					WithDst(rtu.CommandProcessor).
					Build()
				rtu.ControlPort.Send(req)
			} else {
				rtu.PrevEpochImbalance = true
			}
		} else {
			rtu.PrevEpochImbalance = false
		}
		rtu.NumOutgoingPrevEpoch = rtu.NumOutgoing
		rtu.NumOutgoing = 0
		rtu.NumIncomingPrevEpoch = rtu.NumIncoming
		rtu.NumIncoming = 0
		rtu.LastChecked = rtu.NumReqsServiced
	} // end of epoch
}

func (rtu *DefaultRTU) sendCollectedStatsToCP(
	now akita.VTimeInSec, req *akita.SendStatsMsg) bool {

	rsp := akita.CollectedStatsMsgBuilder{}.
		WithSendTime(now).
		WithSrc(rtu.ControlPort).
		WithDst(rtu.CommandProcessor).
		WithIncomingReqs(rtu.NumIncomingPrevEpoch).
		WithOutgoingReqs(rtu.NumOutgoingPrevEpoch).
		Build()
	rtu.ControlPort.Send(rsp)

	rtu.NumOutgoingPrevEpoch = rtu.NumOutgoing
	rtu.NumOutgoing = 0
	rtu.NumIncomingPrevEpoch = rtu.NumIncoming
	rtu.NumIncoming = 0
	rtu.LastChecked = rtu.NumReqsServiced
	return true
}

func (rtu *DefaultRTU) processFromCtrlPort(now akita.VTimeInSec) bool {
	req := rtu.ControlPort.Peek()
	if req == nil {
		return false
	}

	req = rtu.ControlPort.Retrieve(now)
	switch req := req.(type) {
	case *akita.TLBIndexingSwitchMsg:
		fmt.Println("Helloooooooooooooo", rtu.Name(), req.TLBInterleaving)
		return rtu.switchIndexing(now, req)
	case *akita.SendStatsMsg:
		return rtu.sendCollectedStatsToCP(now, req)
	default:
		log.Panicf("cannot process request of type %s", reflect.TypeOf(req))
		return false
	}
}

func (rtu *DefaultRTU) switchIndexing(now akita.VTimeInSec,
	req *akita.TLBIndexingSwitchMsg) bool {

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

	lmf := rtu.RemoteAddressTranslationTable.(*cache.HashLowModuleFinder)
	localLmf := rtu.localModules.(*cache.CustomTwoLevelLowModuleFinder)

	if req.TLBIndexingSwitch == akita.TLBIndexingSwitch4K {
		lmf.Hashfunc = HashFunc4KXOR
		localLmf.Hashfunc = HashFunc4KXOR
	} else {
		lmf.Hashfunc = HashFuncHSL
		localLmf.Hashfunc = HashFuncHSL
	}

	fmt.Println(rtu.Name(), "RTU  ", rtu.Name(), "switching to ", req.TLBIndexingSwitch, " sir", now, req.TLBInterleaving)

	rtu.NumOutgoingPrevEpoch = 0
	rtu.NumIncomingPrevEpoch = 0
	rtu.NumOutgoing = 0
	rtu.NumIncoming = 0

	rtu.NumReqsServiced = 0
	rtu.LastChecked = 0

	return true
}

func (rtu *DefaultRTU) collectCoalescingStats(now akita.VTimeInSec, transactions []transaction, prefix string) {
	vAddrCountMap := make(map[uint64]int)
	for _, t := range transactions {
		var msg akita.Msg
		if t.fromInside != nil {
			msg = t.fromInside
		} else if t.fromOutside != nil {
			msg = t.fromOutside
		} else {
			panic("something is wrong")
		}
		vAddr := msg.(*device.TranslationReq).VAddr
		vAddrCountMap[vAddr]++
	}
	coalescableAddrs := 0
	replicationCount := 0
	for vAddr := range vAddrCountMap {
		if vAddrCountMap[vAddr] > 1 {
			coalescableAddrs++
			replicationCount += vAddrCountMap[vAddr]
		}
	}
	numTransactions := len(transactions)
	tracing.StartTask("", "", now, rtu,
		prefix+"-buflen", strconv.Itoa(numTransactions), nil)
	if numTransactions > 0 {
		tracing.StartTask("", "", now, rtu,
			prefix+"-buflen-g0", strconv.Itoa(numTransactions), nil)
		tracing.StartTask("", "", now, rtu,
			prefix+"-coalesceable-addr", strconv.Itoa(coalescableAddrs), nil)
		tracing.StartTask("", "", now, rtu,
			prefix+"-replication-count", strconv.Itoa(replicationCount), nil)
	}
}

func (rtu *DefaultRTU) processFromL1(now akita.VTimeInSec) bool {
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

func (rtu *DefaultRTU) processFromL2(now akita.VTimeInSec) bool {
	madeProgress := false
	for i := 0; i < rtu.numTransPerCycle; i++ {
		msg := rtu.ToL2.Peek()
		if msg == nil {
			return madeProgress
		}
		// req := msg.(*device.TranslationRsp)

		switch req := msg.(type) {
		case *device.TranslationRsp:
			madeProgress = rtu.processRspFromL2(now, req) || madeProgress
		case *device.TranslationReq:
			madeProgress = rtu.processReqFromL2(now, req) || madeProgress
		default:
			panic("unknown req type")
		}
	}
	return madeProgress
}

func (rtu *DefaultRTU) processReqFromOutside(now akita.VTimeInSec) bool {
	madeProgress := false
	for i := 0; i < rtu.numTransPerCycle; i++ {
		msg := rtu.RequestPort.Peek()
		if msg == nil {
			break
		}

		switch req := msg.(type) {
		case *device.TranslationReq:
			madeProgress = rtu.processReqFromOutsideAux(now, req) || madeProgress
		default:
			log.Panicf("cannot process request of type %s", reflect.TypeOf(req))
			return false
		}
	}
	rtu.collectCoalescingStats(now, rtu.transactionsFromOutside, "incoming")
	return madeProgress
}

func (rtu *DefaultRTU) processRspFromOutside(now akita.VTimeInSec) bool {
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

func (rtu *DefaultRTU) processReqFromL1(
	now akita.VTimeInSec,
	req *device.TranslationReq,
) bool {
	dst := rtu.RemoteAddressTranslationTable.Find(req.VAddr)
	var cloned *device.TranslationReq
	var err *akita.SendError
	if getChipletNum(rtu.Name()) == getChipletNum(dst.Name()) {
		dst = rtu.localModules.Find(req.VAddr)
		if !strings.Contains(dst.Name(), "L2") {
			panic("oh no!")
		}
		req.Dst = dst
		err = rtu.ToL2.Send(req)
	} else {
		cloned = rtu.cloneReq(req)
		cloned.Meta().Src = rtu.RequestPort
		cloned.Meta().Dst = dst
		cloned.Meta().SendTime = now
		err = rtu.RequestPort.Send(cloned)
	}

	if err == nil {
		rtu.ToL1.Retrieve(now)
		if getChipletNum(rtu.Name()) != getChipletNum(dst.Name()) {

			rtu.NumReqsServiced++
			rtu.NumOutgoing++

			tracing.StopTracingNetworkReq(req, now, rtu)
			tracing.StartTracingNetworkReq(cloned, now, rtu, req)

			tracing.TraceReqReceive(req, now, rtu)
			tracing.TraceReqInitiate(cloned, now, rtu,
				tracing.MsgIDAtReceiver(req, rtu))

			srcName := req.Meta().Src.Name()
			dstName := dst.Name()
			srcChip := strings.Split(srcName, ".")[1]
			dstChip := strings.Split(dstName, ".")[1]
			if strings.Contains(srcName, "L1") {
				tracing.AddTaskStep(tracing.MsgIDAtReceiver(req, rtu), now, rtu,
					srcChip+":RTUAccessToChiplet:"+dstChip)
			} else {
				panic("Panic!!!")
			}

			trans := transaction{
				fromInside: req,
				toOutside:  cloned,
			}
			rtu.transactionsFromInside = append(rtu.transactionsFromInside, trans)

			// no stranger to abusing tracing API
			tracing.StartTask("", "", now, rtu, "reference_counting",
				strconv.FormatUint(req.VAddr, 10), nil)

		}
		return true
	}
	return false
}

func (rtu *DefaultRTU) processReqFromOutsideAux(
	now akita.VTimeInSec,
	req *device.TranslationReq,
) bool {
	dst := rtu.RemoteAddressTranslationTable.Find(req.VAddr)
	var err *akita.SendError
	var cloned *device.TranslationReq
	if getChipletNum(dst.Name()) != getChipletNum(rtu.Name()) {
		// cloned = rtu.cloneReq(req)
		// cloned.Meta().Src = req.Src
		// cloned.Meta().Dst = dst
		// cloned.Meta().SendTime = now
		req.Dst = dst
		err = rtu.RequestPort.Send(req)
	} else {
		dst = rtu.localModules.Find(req.VAddr)
		if !strings.Contains(dst.Name(), "L2") {
			panic("oh no!")
		}
		if getChipletNum(req.Src.Name()) == getChipletNum(rtu.Name()) {
			transactionIndex := rtu.findTransactionByRspToID(
				req.ID, rtu.transactionsFromInside)
			trans := rtu.transactionsFromInside[transactionIndex]
			originalReq := trans.fromInside.(*device.TranslationReq)
			originalReq.Dst = dst
			err = rtu.ToL2.Send(originalReq)
		} else {
			cloned = rtu.cloneReq(req)
			cloned.Meta().Src = rtu.ToL2
			cloned.Meta().Dst = dst
			cloned.Meta().SendTime = now
			err = rtu.ToL2.Send(cloned)
		}
	}

	if err == nil {
		rtu.RequestPort.Retrieve(now)
		if cloned != nil && strings.Contains(cloned.Src.Name(), "L2") {

			rtu.NumReqsServiced++
			rtu.NumIncoming++

			tracing.StopTracingNetworkReq(req, now, rtu)
			tracing.StartTracingNetworkReq(cloned, now, rtu, req)

			tracing.TraceReqReceive(req, now, rtu)
			tracing.TraceReqInitiate(cloned, now, rtu,
				tracing.MsgIDAtReceiver(req, rtu))

			trans := transaction{
				fromOutside: req,
				toInside:    cloned,
			}
			rtu.transactionsFromOutside =
				append(rtu.transactionsFromOutside, trans)
		}
		return true
	}
	return false
}

func (rtu *DefaultRTU) getResponsePort(dst akita.Port) akita.Port {
	chiplet, _ := strconv.Atoi(dst.Name()[14:15])
	return rtu.ResponsePorts[chiplet]
}

func (rtu *DefaultRTU) processRspFromL2(
	now akita.VTimeInSec,
	rsp *device.TranslationRsp,
) bool {
	transactionIndex := rtu.findTransactionByRspToID(
		rsp.RespondTo, rtu.transactionsFromOutside)
	trans := rtu.transactionsFromOutside[transactionIndex]

	rspToOutside := rtu.cloneRsp(rsp, trans.fromOutside.Meta().ID)
	rspToOutside.Meta().SendTime = now
	rspToOutside.Meta().Src = rtu.ResponsePort
	rspToOutside.Meta().Dst = rtu.getResponsePort(trans.fromOutside.Meta().Src)

	err := rtu.ResponsePort.Send(rspToOutside)
	if err == nil {
		rtu.ToL2.Retrieve(now)

		tracing.StopTracingNetworkReq(rsp, now, rtu)
		tracing.StartTracingNetworkReq(rspToOutside, now, rtu, rsp)

		tracing.TraceReqFinalize(trans.toInside, now, rtu)
		tracing.TraceReqComplete(trans.fromOutside, now, rtu)

		rtu.transactionsFromOutside =
			append(rtu.transactionsFromOutside[:transactionIndex],
				rtu.transactionsFromOutside[transactionIndex+1:]...)
		return true
	}
	return false
}

func (rtu *DefaultRTU) processReqFromL2(
	now akita.VTimeInSec,
	req *device.TranslationReq,
) bool {
	dst := rtu.RemoteAddressTranslationTable.Find(req.VAddr)
	var cloned *device.TranslationReq
	var err *akita.SendError
	if getChipletNum(rtu.Name()) == getChipletNum(dst.Name()) {
		req.Dst = rtu.localModules.Find(req.VAddr)
		if !strings.Contains(req.Dst.Name(), "L2") {
			panic("oh no!")
		}
		err = rtu.ToL2.Send(req)
		// NEHABUG: there is a bug here. we are not starting to trace a network request here
		// though we really should.
		// dst = rtu.ToL2
		// cloned := rtu.cloneReq(req)
		// cloned.Meta().Dst = dst
		// cloned.Meta().SendTime = now
	} else if strings.Contains(req.Src.Name(), "L1") {
		cloned = rtu.cloneReq(req)
		cloned.Meta().Src = rtu.RequestPort
		cloned.Meta().Dst = dst
		cloned.Meta().SendTime = now
		err = rtu.RequestPort.Send(cloned)
	} else if strings.Contains(req.Src.Name(), "RTU") {
		transactionIndex := rtu.findTransactionByRspToID(req.ID, rtu.transactionsFromOutside)
		trans := rtu.transactionsFromOutside[transactionIndex]
		originalReq := trans.fromOutside.(*device.TranslationReq)
		originalReq.Dst = dst
		err = rtu.RequestPort.Send(originalReq)
	} else {
		panic("oh no!")
	}

	if err == nil {
		rtu.ToL2.Retrieve(now)
		if cloned != nil {

			// strings.Contains(req.Src.Name(), "L1") {
			tracing.StopTracingNetworkReq(req, now, rtu)
			tracing.StartTracingNetworkReq(cloned, now, rtu, req)

			tracing.TraceReqReceive(req, now, rtu)
			tracing.TraceReqInitiate(cloned, now, rtu,
				tracing.MsgIDAtReceiver(req, rtu))

			srcName := req.Meta().Src.Name()
			dstName := dst.Name()
			srcChip := strings.Split(srcName, ".")[1]
			dstChip := strings.Split(dstName, ".")[1]
			if strings.Contains(srcName, "L1") {
				tracing.AddTaskStep(tracing.MsgIDAtReceiver(req, rtu), now, rtu,
					srcChip+":RTUAccessToChiplet:"+dstChip)
			} else {
				panic("Panic!!!")
			}

			trans := transaction{
				fromInside: req,
				toOutside:  cloned,
			}
			rtu.transactionsFromInside = append(rtu.transactionsFromInside, trans)
		}
		return true
	}
	return false
}

func (rtu *DefaultRTU) processRspFromOutsideAux(
	now akita.VTimeInSec,
	rsp *device.TranslationRsp,
) bool {
	transactionIndex := rtu.findTransactionByRspToID(
		rsp.RespondTo, rtu.transactionsFromInside)
	trans := rtu.transactionsFromInside[transactionIndex]

	rspToInside := rtu.cloneRsp(rsp, trans.fromInside.Meta().ID)
	rspToInside.Meta().SendTime = now
	rspToInside.Meta().Src = rtu.ToL1
	rspToInside.Meta().Dst = trans.fromInside.Meta().Src

	err := rtu.ToL1.Send(rspToInside)
	if err == nil {
		rtu.ResponsePort.Retrieve(now)

		tracing.StopTracingNetworkReq(rsp, now, rtu)
		tracing.StartTracingNetworkReq(rspToInside, now, rtu, rsp)

		tracing.TraceReqFinalize(trans.toOutside, now, rtu)
		tracing.TraceReqComplete(trans.fromInside, now, rtu)

		rtu.transactionsFromInside =
			append(rtu.transactionsFromInside[:transactionIndex],
				rtu.transactionsFromInside[transactionIndex+1:]...)

		return true
	}

	return false
}

func (rtu *DefaultRTU) findTransactionByRspToID(
	rspTo string,
	transactions []transaction,
) int {
	for i, trans := range transactions {
		if trans.toOutside != nil && trans.toOutside.Meta().ID == rspTo {
			return i
		}

		if trans.toInside != nil && trans.toInside.Meta().ID == rspTo {
			return i
		}
	}

	log.Panicf("transaction %s not found", rspTo)
	return 0
}

func (rtu *DefaultRTU) cloneReq(origin *device.TranslationReq) *device.TranslationReq {
	req := device.TranslationReqBuilder{}.
		WithSendTime(origin.SendTime).
		WithSrc(origin.Src).
		WithDst(origin.Dst).
		WithPID(origin.PID).
		WithDeviceID(origin.DeviceID).
		WithVAddr(origin.VAddr).
		Build()
	// req.OriginalSender = origin.Src.Name()
	return req
}

func getChipletNum(component string) (i int) {
	i, _ = strconv.Atoi(strings.Split(component, ".")[1][9:10])
	return
}

func (rtu *DefaultRTU) cloneRsp(origin *device.TranslationRsp, rspTo string) *device.TranslationRsp {
	if getChipletNum(origin.Src.Name()) != getChipletNum(origin.SrcL2TLB) {
		panic("oh no!")
	}
	rsp := device.TranslationRspBuilder{}.
		WithSendTime(origin.SendTime).
		WithSrc(origin.Src).
		WithDst(origin.Dst).
		WithRspTo(rspTo).
		WithPage(origin.Page).
		WithAccessResult(origin.HitOrMiss).
		WithSrcL2TLB(origin.SrcL2TLB).
		Build()
	return rsp
}

func NewRemoteTranslationUnit(
	name string,
	engine akita.Engine,
	localModules cache.LowModuleFinder,
	remoteModules cache.LowModuleFinder,
) RemoteTranslationUnit {
	rtu := new(DefaultRTU)
	rtu.TickingComponent = akita.NewTickingComponent(name, engine, 1*akita.GHz, rtu)
	rtu.localModules = localModules
	rtu.RemoteAddressTranslationTable = remoteModules

	rtu.ToL1 = akita.NewLimitNumMsgPort(rtu, 2*12, name+".ToL1")
	rtu.ToL2 = akita.NewLimitNumMsgPort(rtu, 2*12, name+".ToL2")
	rtu.RequestPort = akita.NewLimitNumMsgPort(rtu, 2*12, name+".RequestPort")
	rtu.ResponsePort = akita.NewLimitNumMsgPort(rtu, 2*12, name+".ResponsePort")

	rtu.ControlPort = akita.NewLimitNumMsgPort(rtu, 1, name+".ControlPort")

	rtu.numTransPerCycle = 12

	rtu.NumOutgoing = 0
	rtu.NumIncoming = 0
	rtu.NumOutgoingPrevEpoch = 0
	rtu.NumIncomingPrevEpoch = 0

	return rtu
}
