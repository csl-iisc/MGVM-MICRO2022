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

type transaction struct {
	incomingReqs    []mem.AccessReq
	translationReq  *device.TranslationReq
	translationRsp  *device.TranslationRsp
	translationDone bool
}

type reqToBottom struct {
	reqFromTop  mem.AccessReq
	reqToBottom mem.AccessReq
}

// AddressTranslator is a element that is being simulated in Akita.
type AddressTranslator interface {
	tracing.NamedHookable
	// akita.Component
	GetCtrlPort() akita.Port
	GetTopPort() akita.Port
	GetTranslationPort() akita.Port
	GetBottomPort() akita.Port
	SetLowModuleFinder(cache.LowModuleFinder)
	SetTranslationProvider(interface{})
}

// DefaultAddressTranslator is a component that forwards the read/write requests with
// the address translated from virtual to physical.
type DefaultAddressTranslator struct {
	*akita.TickingComponent

	TopPort         akita.Port
	BottomPort      akita.Port
	TranslationPort akita.Port
	CtrlPort        akita.Port

	lowModuleFinder     cache.LowModuleFinder
	translationProvider akita.Port
	log2PageSize        uint64
	deviceID            uint64
	numReqPerCycle      int

	isFlushing bool

	transactions        []*transaction
	inflightReqToBottom []reqToBottom
}

// SetTranslationProvider sets the remote port that can translate addresses.
func (t *DefaultAddressTranslator) SetTranslationProvider(p interface{}) {
	t.translationProvider = p.(akita.Port)
}

// SetLowModuleFinder sets the table recording where to find an address.
func (t *DefaultAddressTranslator) SetLowModuleFinder(lmf cache.LowModuleFinder) {
	t.lowModuleFinder = lmf
}

// GetCtrlPort sets the table recording where to find an address.
func (t *DefaultAddressTranslator) GetCtrlPort() (port akita.Port) {
	port = t.CtrlPort
	return
}

// GetTopPort sets the table recording where to find an address.
func (t *DefaultAddressTranslator) GetTopPort() (port akita.Port) {
	port = t.TopPort
	return
}

// GetTranslationPort sets the table recording where to find an address.
func (t *DefaultAddressTranslator) GetTranslationPort() (port akita.Port) {
	port = t.TranslationPort
	return
}

// GetBottomPort sets the table recording where to find an address.
func (t *DefaultAddressTranslator) GetBottomPort() (port akita.Port) {
	port = t.BottomPort
	return
}

// Tick updates state at each cycle.
func (t *DefaultAddressTranslator) Tick(now akita.VTimeInSec) bool {
	madeProgress := false

	if !t.isFlushing {
		madeProgress = t.runPipeline(now)
	}

	madeProgress = t.handleCtrlRequest(now) || madeProgress

	return madeProgress
}

func (t *DefaultAddressTranslator) runPipeline(now akita.VTimeInSec) bool {
	madeProgress := false

	for i := 0; i < t.numReqPerCycle; i++ {
		madeProgress = t.respond(now) || madeProgress
	}

	for i := 0; i < t.numReqPerCycle; i++ {
		madeProgress = t.sendToBottom(now) || madeProgress
	}

	for i := 0; i < t.numReqPerCycle; i++ {
		madeProgress = t.parseTranslation(now) || madeProgress
	}

	for i := 0; i < t.numReqPerCycle; i++ {
		madeProgress = t.translate(now) || madeProgress
	}

	return madeProgress
}

// func getKind(t *DefaultAddressTranslator) (s string) {
// 	// fmt.Println(t.Name())
// 	s = "sent-" + strings.Split(strings.Split(t.Name(), ".")[3], "_")[0]
// 	// fmt.Println(s)
// 	return
// }

func (t *DefaultAddressTranslator) translate(now akita.VTimeInSec) bool {
	// if len(t.transactions) >= 64 {
	// return false
	// }
	item := t.TopPort.Peek()
	if item == nil {
		return false
	}

	req := item.(mem.AccessReq)
	vAddr := req.GetAddress()
	vPageID := t.addrToPageID(vAddr)

	transReq := device.TranslationReqBuilder{}.
		WithSendTime(now).
		WithSrc(t.TranslationPort).
		WithDst(t.translationProvider).
		WithPID(req.GetPID()).
		WithVAddr(vPageID).
		WithDeviceID(t.deviceID).
		Build()
	err := t.TranslationPort.Send(transReq)
	if err != nil {
		return false
	}

	translation := &transaction{
		incomingReqs:   []mem.AccessReq{req},
		translationReq: transReq,
	}
	t.transactions = append(t.transactions, translation)

	tracing.StartTask(
		transReq.Meta().ID+"-addr-translator-stats",
		tracing.MsgIDAtReceiver(req, t),
		now,
		t,
		"addr-translator-stats",
		reflect.TypeOf(transReq).String(),
		transReq,
	)
	tracing.StartTask(
		req.Meta().ID+"-addr-translator-stats",
		tracing.MsgIDAtReceiver(req, t),
		now,
		t,
		"addr-translator-stats",
		reflect.TypeOf(req).String(),
		req,
	)
	tracing.StartTracingNetworkReq(transReq, now, t, req)
	tracing.TraceReqReceive(req, now, t)
	tracing.TraceReqInitiate(transReq, now, t, tracing.MsgIDAtReceiver(req, t))

	t.TopPort.Retrieve(now)

	return true
}

func (t *DefaultAddressTranslator) parseTranslation(now akita.VTimeInSec) bool {
	rsp := t.TranslationPort.Peek()
	if rsp == nil {
		return false
	}

	transRsp := rsp.(*device.TranslationRsp)
	transaction := t.findTranslationByReqID(transRsp.RespondTo)
	if transaction == nil {
		t.TranslationPort.Retrieve(now)
		return true
	}

	transaction.translationRsp = transRsp
	transaction.translationDone = true

	t.TranslationPort.Retrieve(now)

	tracing.AddTaskStep(
		transaction.translationReq.Meta().ID+"-addr-translator-stats",
		now, t,
		"translation-latency",
	)
	tracing.EndTask(
		transaction.translationReq.Meta().ID+"-addr-translator-stats",
		now,
		t,
	)
	tracing.StopTracingNetworkReq(transRsp, now, t)
	tracing.TraceReqFinalize(transaction.translationReq, now, t)

	return true
}
func (t *DefaultAddressTranslator) sendToBottom(now akita.VTimeInSec) bool {
	for _, transaction := range t.transactions {
		if transaction.translationDone {
			reqFromTop := transaction.incomingReqs[0]
			translatedReq := t.createTranslatedReq(
				reqFromTop,
				transaction.translationRsp.Page)
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
			transaction.incomingReqs = transaction.incomingReqs[1:]
			if len(transaction.incomingReqs) == 0 {
				t.removeExistingTranslation(transaction)
			}
			tracing.TraceReqInitiate(translatedReq, now, t,
				tracing.MsgIDAtReceiver(reqFromTop, t))

			return true
		}
	}
	return false
}

func (t *DefaultAddressTranslator) respond(now akita.VTimeInSec) bool {
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

func (t *DefaultAddressTranslator) createTranslatedReq(
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

func (t *DefaultAddressTranslator) createTranslatedReadReq(
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

func (t *DefaultAddressTranslator) createTranslatedWriteReq(
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

func (t *DefaultAddressTranslator) addrToPageID(addr uint64) uint64 {
	return (addr >> t.log2PageSize) << t.log2PageSize
}

func (t *DefaultAddressTranslator) findTranslationByReqID(id string) *transaction {
	for _, t := range t.transactions {
		if t.translationReq.ID == id {
			return t
		}
	}
	return nil
}

func (t *DefaultAddressTranslator) removeExistingTranslation(trans *transaction) {
	for i, tr := range t.transactions {
		if tr == trans {
			t.transactions = append(t.transactions[:i], t.transactions[i+1:]...)
			return
		}
	}
	panic("translation not found")
}

func (t *DefaultAddressTranslator) isReqInBottomByID(id string) bool {
	for _, r := range t.inflightReqToBottom {
		if r.reqToBottom.Meta().ID == id {
			return true
		}
	}
	return false
}

func (t *DefaultAddressTranslator) findReqToBottomByID(id string) reqToBottom {
	for _, r := range t.inflightReqToBottom {
		if r.reqToBottom.Meta().ID == id {
			return r
		}
	}
	panic("req to bottom not found")
}

func (t *DefaultAddressTranslator) removeReqToBottomByID(id string) {
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

func (t *DefaultAddressTranslator) handleCtrlRequest(now akita.VTimeInSec) bool {
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

func (t *DefaultAddressTranslator) handleFlushReq(
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

	t.transactions = nil
	t.inflightReqToBottom = nil
	t.isFlushing = true

	return true
}

func (t *DefaultAddressTranslator) handleRestartReq(
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

	for t.TranslationPort.Retrieve(now) != nil {
	}

	t.isFlushing = false

	t.CtrlPort.Retrieve(now)

	return true
}
