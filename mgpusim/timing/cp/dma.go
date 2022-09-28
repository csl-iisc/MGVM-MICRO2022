package cp

import (
	"log"
	"reflect"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem"
	"gitlab.com/akita/mem/cache"
	"gitlab.com/akita/mgpusim/protocol"
	"gitlab.com/akita/util/tracing"
)

// A DMAEngine is responsible for accessing data that does not belongs to
// the GPU that the DMAEngine works in.
type DMAEngine struct {
	*akita.TickingComponent

	Log2AccessSize uint64

	localDataSource cache.LowModuleFinder

	processingReq akita.Msg

	toSendToMem []akita.Msg
	toSendToCP  []akita.Msg
	pendingReqs []akita.Msg

	ToCP             akita.Port
	ToMem            []akita.Port
	DramToChipletMap map[akita.Port]akita.Port
}

// SetLocalDataSource sets the table that maps from addresses to port that can
// provide the data.
func (dma *DMAEngine) SetLocalDataSource(s cache.LowModuleFinder) {
	dma.localDataSource = s
}

// Tick ticks
func (dma *DMAEngine) Tick(now akita.VTimeInSec) bool {
	madeProgress := false

	madeProgress = dma.sendToCP(now, &dma.toSendToCP) || madeProgress
	madeProgress = dma.sendToMem(now, &dma.toSendToMem) || madeProgress
	madeProgress = dma.parseFromMem(now) || madeProgress
	madeProgress = dma.parseFromCP(now) || madeProgress

	return madeProgress
}

func (dma *DMAEngine) sendToCP(
	now akita.VTimeInSec,
	reqs *[]akita.Msg,
) bool {
	if len(*reqs) == 0 {
		return false
	}

	req := (*reqs)[0]
	req.Meta().SendTime = now
	port := dma.ToCP
	err := port.Send(req)
	if err == nil {
		*reqs = (*reqs)[1:]
		return true
	}

	return false
}
func (dma *DMAEngine) sendToMem(
	now akita.VTimeInSec,
	reqs *[]akita.Msg,
) bool {
	if len(*reqs) == 0 {
		return false
	}

	req := (*reqs)[0]
	req.Meta().SendTime = now
	port := dma.DramToChipletMap[req.Meta().Dst]
	err := port.Send(req)
	if err == nil {
		*reqs = (*reqs)[1:]
		return true
	}

	return false
}

func (dma *DMAEngine) parseFromMem(now akita.VTimeInSec) bool {
	for _, port := range dma.ToMem {
		req := port.Retrieve(now)
		if req == nil {
			continue
		}

		switch req := req.(type) {
		case *mem.DataReadyRsp:
			dma.processDataReadyRsp(now, req)
		case *mem.WriteDoneRsp:
			dma.processDoneRsp(now, req)
		default:
			log.Panicf("cannot handle request of type %s", reflect.TypeOf(req))
		}

		return true
	}
	return false
}

func (dma *DMAEngine) processDataReadyRsp(
	now akita.VTimeInSec,
	rsp *mem.DataReadyRsp,
) {
	req := dma.removeReqFromPendingReqList(rsp.RespondTo).(*mem.ReadReq)
	tracing.TraceReqFinalize(req, now, dma)

	processing := dma.processingReq.(*protocol.MemCopyD2HReq)

	offset := req.Address - processing.SrcAddress
	copy(processing.DstBuffer[offset:], rsp.Data)
	// fmt.Printf("Dma DataReady %x, %v\n", req.Address, rsp.Data)

	if len(dma.pendingReqs) == 0 {
		tracing.TraceReqComplete(dma.processingReq, now, dma)
		dma.processingReq = nil
		processing.Src, processing.Dst = processing.Dst, processing.Src
		dma.toSendToCP = append(dma.toSendToCP, processing)
	}
}

func (dma *DMAEngine) processDoneRsp(
	now akita.VTimeInSec,
	rsp *mem.WriteDoneRsp,
) {
	r := dma.removeReqFromPendingReqList(rsp.RespondTo)
	tracing.TraceReqFinalize(r, now, dma)

	processing := dma.processingReq.(*protocol.MemCopyH2DReq)
	if len(dma.pendingReqs) == 0 {
		tracing.TraceReqComplete(dma.processingReq, now, dma)
		dma.processingReq = nil
		processing.Src, processing.Dst = processing.Dst, processing.Src
		dma.toSendToCP = append(dma.toSendToCP, processing)
	}
}

func (dma *DMAEngine) removeReqFromPendingReqList(id string) akita.Msg {
	var reqToRet akita.Msg
	newList := make([]akita.Msg, 0, len(dma.pendingReqs)-1)
	for _, r := range dma.pendingReqs {
		if r.Meta().ID == id {
			reqToRet = r
		} else {
			newList = append(newList, r)
		}
	}
	dma.pendingReqs = newList

	if reqToRet == nil {
		panic("not found")
	}

	return reqToRet
}

func (dma *DMAEngine) parseFromCP(now akita.VTimeInSec) bool {
	if dma.processingReq != nil {
		return false
	}

	req := dma.ToCP.Retrieve(now)
	if req == nil {
		return false
	}
	tracing.TraceReqReceive(req, now, dma)

	dma.processingReq = req
	switch req := req.(type) {
	case *protocol.MemCopyH2DReq:
		dma.parseMemCopyH2D(now, req)
	case *protocol.MemCopyD2HReq:
		dma.parseMemCopyD2H(now, req)
	default:
		log.Panicf("cannot process request of type %s", reflect.TypeOf(req))
	}

	return true
}

func (dma *DMAEngine) parseMemCopyH2D(
	now akita.VTimeInSec,
	req *protocol.MemCopyH2DReq,
) {
	offset := uint64(0)
	lengthLeft := uint64(len(req.SrcBuffer))
	addr := req.DstAddress
	// fmt.Println("address:", addr)
	for lengthLeft > 0 {
		addrUnitFirstByte := addr & (^uint64(0) << dma.Log2AccessSize)
		unitOffset := addr - addrUnitFirstByte
		lengthInUnit := (1 << dma.Log2AccessSize) - unitOffset
		// fmt.Println(dma.Log2AccessSize, lengthInUnit, lengthLeft, unitOffset)
		length := lengthLeft
		if lengthInUnit < length {
			length = lengthInUnit
		}

		module := dma.localDataSource.Find(addr)
		// fmt.Println("module:", module.Name(), addr)
		reqToBottom := mem.WriteReqBuilder{}.
			WithSendTime(now).
			WithSrc(dma.DramToChipletMap[module]).
			WithDst(module).
			WithAddress(addr).
			WithData(req.SrcBuffer[offset : offset+length]).
			Build()
		dma.toSendToMem = append(dma.toSendToMem, reqToBottom)
		dma.pendingReqs = append(dma.pendingReqs, reqToBottom)

		tracing.TraceReqInitiate(reqToBottom, now, dma,
			tracing.MsgIDAtReceiver(dma.processingReq, dma))

		addr += length
		lengthLeft -= length
		offset += length
	}
}

func (dma *DMAEngine) parseMemCopyD2H(
	now akita.VTimeInSec,
	req *protocol.MemCopyD2HReq,
) {
	offset := uint64(0)
	lengthLeft := uint64(len(req.DstBuffer))
	addr := req.SrcAddress

	for lengthLeft > 0 {
		addrUnitFirstByte := addr & (^uint64(0) << dma.Log2AccessSize)
		unitOffset := addr - addrUnitFirstByte
		lengthInUnit := (1 << dma.Log2AccessSize) - unitOffset

		length := lengthLeft
		if lengthInUnit < length {
			length = lengthInUnit
		}

		module := dma.localDataSource.Find(addr)
		reqToBottom := mem.ReadReqBuilder{}.
			WithSendTime(now).
			WithSrc(dma.DramToChipletMap[module]).
			WithDst(module).
			WithAddress(addr).
			WithByteSize(length).
			Build()
		dma.toSendToMem = append(dma.toSendToMem, reqToBottom)
		dma.pendingReqs = append(dma.pendingReqs, reqToBottom)

		tracing.TraceReqInitiate(reqToBottom, now, dma,
			tracing.MsgIDAtReceiver(dma.processingReq, dma))

		addr += length
		lengthLeft -= length
		offset += length
	}
}

// NewDMAEngine creates a DMAEngine, injecting a engine and a "LowModuleFinder"
// that helps with locating the module that holds the data.
func NewDMAEngine(
	name string,
	engine akita.Engine,
	localDataSource cache.LowModuleFinder,
) *DMAEngine {
	dma := new(DMAEngine)
	dma.TickingComponent = akita.NewTickingComponent(
		name, engine, 1*akita.GHz, dma)

	dma.Log2AccessSize = 6
	dma.localDataSource = localDataSource

	dma.ToCP = akita.NewLimitNumMsgPort(dma, 40960000, name+".ToCP")
	// dma.ToMem = akita.NewLimitNumMsgPort(dma, 64, name+".ToMem")
	dma.DramToChipletMap = make(map[akita.Port]akita.Port)

	return dma
}
