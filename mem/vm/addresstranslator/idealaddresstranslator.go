package addresstranslator

import (
	"log"
	"reflect"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem"
	"gitlab.com/akita/mem/cache"
	"gitlab.com/akita/mem/device"
	"gitlab.com/akita/util/tracing"
)

// AddressTranslator is a component that forwards the read/write requests with
// the address translated from virtual to physical.
type IdealAddressTranslator struct {
	*akita.TickingComponent

	TopPort    akita.Port
	BottomPort akita.Port
	CtrlPort   akita.Port

	lowModuleFinder cache.LowModuleFinder
	log2PageSize    uint64
	deviceID        uint64
	numReqPerCycle  int

	isFlushing bool

	inflightReqToBottom []reqToBottom

	pageTable device.PageTable
}

// SetTranslationProvider sets the remote port that can translate addresses.
func (t *IdealAddressTranslator) SetTranslationProvider(p interface{}) {
	t.pageTable = p.(device.PageTable)
}

// SetLowModuleFinder sets the table recording where to find an address.
func (t *IdealAddressTranslator) SetLowModuleFinder(lmf cache.LowModuleFinder) {
	t.lowModuleFinder = lmf
}

// GetCtrlPort sets the table recording where to find an address.
func (t *IdealAddressTranslator) GetCtrlPort() (port akita.Port) {
	port = t.CtrlPort
	return
}

// GetTopPort sets the table recording where to find an address.
func (t *IdealAddressTranslator) GetTopPort() (port akita.Port) {
	port = t.TopPort
	return
}

// GetTopPort sets the table recording where to find an address.
func (t *IdealAddressTranslator) GetTranslationPort() (port akita.Port) {
	panic("Should not ask for translation port in idealaddresstranslator")
	return
}

// GetBottomPort sets the table recording where to find an address.
func (t *IdealAddressTranslator) GetBottomPort() (port akita.Port) {
	port = t.BottomPort
	return
}

// Tick updates state at each cycle.
func (t *IdealAddressTranslator) Tick(now akita.VTimeInSec) bool {
	madeProgress := false

	if !t.isFlushing {
		madeProgress = t.runPipeline(now)
	}

	madeProgress = t.handleCtrlRequest(now) || madeProgress

	return madeProgress
}

func (t *IdealAddressTranslator) runPipeline(now akita.VTimeInSec) bool {
	madeProgress := false

	for i := 0; i < t.numReqPerCycle; i++ {
		madeProgress = t.respond(now) || madeProgress
	}

	// for i := 0; i < t.numReqPerCycle; i++ {
	// 	madeProgress = t.parseTranslation(now) || madeProgress
	// }

	for i := 0; i < t.numReqPerCycle; i++ {
		madeProgress = t.translate(now) || madeProgress
	}

	return madeProgress
}

func (t *IdealAddressTranslator) translate(now akita.VTimeInSec) bool {
	item := t.TopPort.Peek()
	if item == nil {
		return false
	}

	reqFromTop := item.(mem.AccessReq)
	vAddr := reqFromTop.GetAddress()
	PID := reqFromTop.GetPID()
	vPageID := t.addrToPageID(vAddr)
	page, ok := t.pageTable.Find(PID, vPageID)
	if !ok {
		panic("Could not find virtual to physical address translation!")
	}
	translatedReq := t.createTranslatedReq(
		reqFromTop,
		page)
	translatedReq.Meta().SendTime = now
	err := t.BottomPort.Send(translatedReq)
	if err != nil {
		return false
	}

	t.inflightReqToBottom = append(t.inflightReqToBottom,
		reqToBottom{
			reqFromTop:  reqFromTop,
			reqToBottom: translatedReq,
		})

	t.TopPort.Retrieve(now)

	// tracing.StopTracingNetworkReq(reqFromTop, now, t)
	// tracing.StartTracingNetworkReq(translatedReq, now, t, reqFromTop)
	tracing.StartTask(
		reqFromTop.Meta().ID+"-addr-translator-stats",
		tracing.MsgIDAtReceiver(reqFromTop, t),
		now,
		t,
		"addr-translator-stats",
		reflect.TypeOf(reqFromTop).String(),
		reqFromTop,
	)
	tracing.TraceReqFinalize(reqFromTop, now, t)
	tracing.TraceReqInitiate(translatedReq, now, t,
		tracing.MsgIDAtReceiver(reqFromTop, t))

	return true
}

func (t *IdealAddressTranslator) respond(now akita.VTimeInSec) bool {
	rsp := t.BottomPort.Peek()
	if rsp == nil {
		return false
	}

	reqInBottom := false

	var reqFromTop mem.AccessReq
	var reqToBottomCombo reqToBottom
	var rspToTop mem.AccessRsp
	switch rsp := rsp.(type) {
	case *mem.DataReadyRsp:
		reqInBottom = t.isReqInBottomByID(rsp.RespondTo)
		if reqInBottom {
			reqToBottomCombo = t.findReqToBottomByID(rsp.RespondTo)
			reqFromTop = reqToBottomCombo.reqFromTop
			drToTop := mem.DataReadyRspBuilder{}.
				WithSendTime(now).
				WithSrc(t.TopPort).
				WithDst(reqFromTop.Meta().Src).
				WithRspTo(reqFromTop.Meta().ID).
				WithData(rsp.Data).
				Build()
			rspToTop = drToTop
		}
	case *mem.WriteDoneRsp:
		reqInBottom = t.isReqInBottomByID(rsp.RespondTo)
		if reqInBottom {
			reqToBottomCombo = t.findReqToBottomByID(rsp.RespondTo)
			reqFromTop = reqToBottomCombo.reqFromTop
			rspToTop = mem.WriteDoneRspBuilder{}.
				WithSendTime(now).
				WithSrc(t.TopPort).
				WithDst(reqFromTop.Meta().Src).
				WithRspTo(reqFromTop.Meta().ID).
				Build()
		}
	default:
		log.Panicf("cannot handle respond of type %s", reflect.TypeOf(rsp))
	}

	if reqInBottom {
		err := t.TopPort.Send(rspToTop)
		if err != nil {
			return false
		}

		t.removeReqToBottomByID(rsp.(mem.AccessRsp).GetRespondTo())

		tracing.AddTaskStep(
			reqToBottomCombo.reqFromTop.Meta().ID+"-addr-translator-stats",
			now, t,
			"req-latency",
		)
		tracing.EndTask(
			reqToBottomCombo.reqFromTop.Meta().ID+"-addr-translator-stats",
			now,
			t,
		)
		tracing.TraceReqFinalize(reqToBottomCombo.reqToBottom, now, t)
		tracing.TraceReqComplete(reqToBottomCombo.reqFromTop, now, t)
	}

	t.BottomPort.Retrieve(now)
	return true
}

func (t *IdealAddressTranslator) createTranslatedReq(
	req mem.AccessReq,
	page device.Page,
) mem.AccessReq {
	switch req := req.(type) {
	case *mem.ReadReq:
		return t.createTranslatedReadReq(req, page)
	case *mem.WriteReq:
		return t.createTranslatedWriteReq(req, page)
	default:
		log.Panicf("cannot translate request of type %s", reflect.TypeOf(req))
		return nil
	}
}

func (t *IdealAddressTranslator) createTranslatedReadReq(
	req *mem.ReadReq,
	page device.Page,
) *mem.ReadReq {
	offset := req.Address % (1 << t.log2PageSize)
	addr := page.PAddr + offset
	clone := mem.ReadReqBuilder{}.
		WithSrc(t.BottomPort).
		WithDst(t.lowModuleFinder.Find(addr)).
		WithAddress(addr).
		WithByteSize(req.AccessByteSize).
		WithPID(0).
		WithInfo(req.Info).
		Build()
	clone.CanWaitForCoalesce = req.CanWaitForCoalesce
	return clone
}

func (t *IdealAddressTranslator) createTranslatedWriteReq(
	req *mem.WriteReq,
	page device.Page,
) *mem.WriteReq {
	offset := req.Address % (1 << t.log2PageSize)
	addr := page.PAddr + offset
	clone := mem.WriteReqBuilder{}.
		WithSrc(t.BottomPort).
		WithDst(t.lowModuleFinder.Find(addr)).
		WithData(req.Data).
		WithDirtyMask(req.DirtyMask).
		WithAddress(addr).
		WithPID(0).
		WithInfo(req.Info).
		Build()
	clone.CanWaitForCoalesce = req.CanWaitForCoalesce
	return clone
}

func (t *IdealAddressTranslator) addrToPageID(addr uint64) uint64 {
	return (addr >> t.log2PageSize) << t.log2PageSize
}

// func (t *IdealAddressTranslator) findTranslationByReqID(id string) *transaction {
// 	for _, t := range t.transactions {
// 		if t.translationReq.ID == id {
// 			return t
// 		}
// 	}
// 	return nil
// }

// func (t *IdealAddressTranslator) removeExistingTranslation(trans *transaction) {
// 	for i, tr := range t.transactions {
// 		if tr == trans {
// 			t.transactions = append(t.transactions[:i], t.transactions[i+1:]...)
// 			return
// 		}
// 	}
// 	panic("translation not found")
// }

func (t *IdealAddressTranslator) isReqInBottomByID(id string) bool {
	for _, r := range t.inflightReqToBottom {
		if r.reqToBottom.Meta().ID == id {
			return true
		}
	}
	return false
}

func (t *IdealAddressTranslator) findReqToBottomByID(id string) reqToBottom {
	for _, r := range t.inflightReqToBottom {
		if r.reqToBottom.Meta().ID == id {
			return r
		}
	}
	panic("req to bottom not found")
}

func (t *IdealAddressTranslator) removeReqToBottomByID(id string) {
	for i, r := range t.inflightReqToBottom {
		if r.reqToBottom.Meta().ID == id {
			t.inflightReqToBottom = append(
				t.inflightReqToBottom[:i],
				t.inflightReqToBottom[i+1:]...)
			return
		}
	}
	panic("req to bottom not found")
}

func (t *IdealAddressTranslator) handleCtrlRequest(now akita.VTimeInSec) bool {
	req := t.CtrlPort.Peek()
	if req == nil {
		return false
	}

	msg := req.(*mem.ControlMsg)

	if msg.DiscardTransations {
		return t.handleFlushReq(now, msg)
	} else if msg.Restart {
		return t.handleRestartReq(now, msg)
	}

	panic("never")
}

func (t *IdealAddressTranslator) handleFlushReq(
	now akita.VTimeInSec,
	req *mem.ControlMsg,
) bool {
	rsp := mem.ControlMsgBuilder{}.
		WithSrc(t.CtrlPort).
		WithDst(req.Src).
		WithSendTime(now).
		ToNotifyDone().
		Build()

	err := t.CtrlPort.Send(rsp)
	if err != nil {
		return false
	}

	t.CtrlPort.Retrieve(now)

	// t.transactions = nil
	t.inflightReqToBottom = nil
	t.isFlushing = true

	return true
}

func (t *IdealAddressTranslator) handleRestartReq(
	now akita.VTimeInSec,
	req *mem.ControlMsg,
) bool {
	rsp := mem.ControlMsgBuilder{}.
		WithSrc(t.CtrlPort).
		WithDst(req.Src).
		WithSendTime(now).
		ToNotifyDone().
		Build()

	err := t.CtrlPort.Send(rsp)

	if err != nil {
		return false
	}

	for t.TopPort.Retrieve(now) != nil {
	}

	for t.BottomPort.Retrieve(now) != nil {
	}

	// for t.TranslationPort.Retrieve(now) != nil {
	// }

	t.isFlushing = false

	t.CtrlPort.Retrieve(now)

	return true
}
