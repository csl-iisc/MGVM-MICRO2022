package addresstranslator

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem"
	"gitlab.com/akita/mem/device"
	"gitlab.com/akita/util/ca"
)

var _ = Describe("Address Translator", func() {
	var (
		mockCtrl        *gomock.Controller
		topPort         *MockPort
		bottomPort      *MockPort
		translationPort *MockPort
		ctrlPort        *MockPort
		lowModuleFinder *MockLowModuleFinder

		t *AddressTranslator
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		topPort = NewMockPort(mockCtrl)
		bottomPort = NewMockPort(mockCtrl)
		ctrlPort = NewMockPort(mockCtrl)
		translationPort = NewMockPort(mockCtrl)
		lowModuleFinder = NewMockLowModuleFinder(mockCtrl)

		builder := MakeBuilder().
			WithLog2PageSize(12).
			WithFreq(1).
			WithLowModuleFinder(lowModuleFinder)
		t = builder.Build("address_translator")
		t.log2PageSize = 12
		t.TopPort = topPort
		t.BottomPort = bottomPort
		t.TranslationPort = translationPort
		t.CtrlPort = ctrlPort
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("translate stage", func() {
		var (
			req *mem.ReadReq
		)

		BeforeEach(func() {
			req = mem.ReadReqBuilder{}.
				WithSendTime(8).
				WithAddress(0x100).
				WithByteSize(4).
				WithPID(1).
				Build()
		})

		It("should do nothing if there is no request", func() {
			topPort.EXPECT().Peek().Return(nil)
			madeProgress := t.translate(10)
			Expect(madeProgress).To(BeFalse())
		})

		It("should send translation", func() {
			var transReqReturn *device.TranslationReq
			transReq := device.TranslationReqBuilder{}.
				WithSendTime(6).
				WithPID(1).
				WithVAddr(0x100).
				WithDeviceID(1).
				Build()

			translation := &transaction{
				translationReq: transReq,
			}
			t.transactions = append(t.transactions, translation)
			req.Address = 0x1040

			topPort.EXPECT().Peek().Return(req)
			topPort.EXPECT().Retrieve(gomock.Any())
			translationPort.EXPECT().Send(gomock.Any()).
				DoAndReturn(func(req *device.TranslationReq) *akita.SendError {
					transReqReturn = req
					return nil
				})

			needTick := t.translate(10)

			Expect(needTick).To(BeTrue())
			Expect(translation.incomingReqs).NotTo(ContainElement(req))
			Expect(t.transactions).To(HaveLen(2))
			Expect(t.transactions[1].translationReq).
				To(BeEquivalentTo(transReqReturn))
		})

		It("should stall if cannot send for translation", func() {
			topPort.EXPECT().Peek().Return(req)
			translationPort.EXPECT().
				Send(gomock.Any()).
				Return(&akita.SendError{})

			needTick := t.translate(10)

			Expect(needTick).To(BeFalse())
			Expect(t.transactions).To(HaveLen(0))
		})
	})

	Context("parse translation", func() {
		var (
			transReq1, transReq2 *device.TranslationReq
			trans1, trans2       *transaction
		)

		BeforeEach(func() {
			transReq1 = device.TranslationReqBuilder{}.
				WithSendTime(0).
				WithPID(1).
				WithVAddr(0x100).
				WithDeviceID(1).
				Build()
			trans1 = &transaction{
				translationReq: transReq1,
			}
			transReq2 = device.TranslationReqBuilder{}.
				WithSendTime(0).
				WithPID(1).
				WithVAddr(0x100).
				WithDeviceID(1).
				Build()
			trans2 = &transaction{
				translationReq: transReq2,
			}
			t.transactions = append(t.transactions, trans1, trans2)
		})

		It("should do nothing if there is no translation return", func() {
			translationPort.EXPECT().Peek().Return(nil)
			needTick := t.parseTranslation(10)
			Expect(needTick).To(BeFalse())
		})

		It("should stall if send failed", func() {
			req := mem.ReadReqBuilder{}.
				WithSendTime(6).
				WithAddress(0x10040).
				WithByteSize(4).
				Build()
			translationRsp := device.TranslationRspBuilder{}.
				WithSendTime(8).
				WithRspTo(transReq1.ID).
				WithPage(device.Page{
					PID:   1,
					VAddr: 0x10000,
					PAddr: 0x20000,
				}).
				Build()

			trans1.incomingReqs = []mem.AccessReq{req}
			trans1.translationRsp = translationRsp
			trans1.translationDone = true

			translationPort.EXPECT().Peek().Return(translationRsp)
			lowModuleFinder.EXPECT().Find(uint64(0x20040))
			bottomPort.EXPECT().Send(gomock.Any()).Return(akita.NewSendError())

			madeProgress := t.parseTranslation(10)

			Expect(madeProgress).To(BeFalse())
		})

		It("should forward read request", func() {
			req := mem.ReadReqBuilder{}.
				WithSendTime(6).
				WithAddress(0x10040).
				WithByteSize(4).
				Build()
			translationRsp := device.TranslationRspBuilder{}.
				WithSendTime(8).
				WithRspTo(transReq1.ID).
				WithPage(device.Page{
					PID:   1,
					VAddr: 0x10000,
					PAddr: 0x20000,
				}).
				Build()

			trans1.incomingReqs = []mem.AccessReq{req}
			trans1.translationRsp = translationRsp
			trans1.translationDone = true

			translationPort.EXPECT().Peek().Return(translationRsp)
			translationPort.EXPECT().Retrieve(akita.VTimeInSec(10))
			lowModuleFinder.EXPECT().Find(uint64(0x20040))
			bottomPort.EXPECT().Send(gomock.Any()).
				Do(func(read *mem.ReadReq) {
					Expect(read).NotTo(BeIdenticalTo(req))
					Expect(read.SendTime).To(Equal(akita.VTimeInSec(10)))
					Expect(read.PID).To(Equal(ca.PID(0)))
					Expect(read.Address).To(Equal(uint64(0x20040)))
					Expect(read.AccessByteSize).To(Equal(uint64(4)))
					Expect(read.Src).To(BeIdenticalTo(bottomPort))
				}).
				Return(nil)

			madeProgress := t.parseTranslation(10)

			Expect(madeProgress).To(BeTrue())
			Expect(t.transactions).NotTo(ContainElement(trans1))
			Expect(t.inflightReqToBottom).To(HaveLen(1))
		})

		It("should forward write request", func() {
			data := []byte{1, 2, 3, 4}
			dirty := []bool{false, true, false, true}
			write := mem.WriteReqBuilder{}.
				WithSendTime(6).
				WithAddress(0x10040).
				WithData(data).
				WithDirtyMask(dirty).
				Build()
			translationRsp := device.TranslationRspBuilder{}.
				WithSendTime(8).
				WithRspTo(transReq1.ID).
				WithPage(device.Page{
					PID:   1,
					VAddr: 0x10000,
					PAddr: 0x20000,
				}).
				Build()
			trans1.incomingReqs = []mem.AccessReq{write}
			trans1.translationRsp = translationRsp
			trans1.translationDone = true

			translationPort.EXPECT().Peek().Return(translationRsp)
			translationPort.EXPECT().Retrieve(akita.VTimeInSec(10))
			lowModuleFinder.EXPECT().Find(uint64(0x20040))
			bottomPort.EXPECT().Send(gomock.Any()).
				Do(func(req *mem.WriteReq) {
					Expect(req).NotTo(BeIdenticalTo(write))
					Expect(req.SendTime).To(Equal(akita.VTimeInSec(10)))
					Expect(req.PID).To(Equal(ca.PID(0)))
					Expect(req.Address).To(Equal(uint64(0x20040)))
					Expect(req.Src).To(BeIdenticalTo(bottomPort))
					Expect(req.Data).To(Equal(data))
					Expect(req.DirtyMask).To(Equal(dirty))
				}).
				Return(nil)

			madeProgress := t.parseTranslation(10)

			Expect(madeProgress).To(BeTrue())
			Expect(t.transactions).NotTo(ContainElement(trans1))
			Expect(t.inflightReqToBottom).To(HaveLen(1))
		})
	})

	Context("respond", func() {
		var (
			readFromTop   *mem.ReadReq
			writeFromTop  *mem.WriteReq
			readToBottom  *mem.ReadReq
			writeToBottom *mem.WriteReq
		)

		BeforeEach(func() {
			readFromTop = mem.ReadReqBuilder{}.
				WithSendTime(8).
				WithAddress(0x10040).
				WithByteSize(4).
				Build()
			readToBottom = mem.ReadReqBuilder{}.
				WithSendTime(8).
				WithAddress(0x20040).
				WithByteSize(4).
				Build()
			writeFromTop = mem.WriteReqBuilder{}.
				WithSendTime(8).
				WithAddress(0x10040).
				Build()
			writeToBottom = mem.WriteReqBuilder{}.
				WithSendTime(8).
				WithAddress(0x10040).
				Build()

			t.inflightReqToBottom = []reqToBottom{
				{reqFromTop: readFromTop, reqToBottom: readToBottom},
				{reqFromTop: writeFromTop, reqToBottom: writeToBottom},
			}

		})

		It("should do nothing if there is no response to process", func() {
			bottomPort.EXPECT().Peek().Return(nil)
			madeProgress := t.respond(10)
			Expect(madeProgress).To(BeFalse())
		})

		It("should respond data ready", func() {
			dataReady := mem.DataReadyRspBuilder{}.
				WithSendTime(10).
				WithRspTo(readToBottom.ID).
				Build()
			bottomPort.EXPECT().Peek().Return(dataReady)
			topPort.EXPECT().Send(gomock.Any()).
				Do(func(dr *mem.DataReadyRsp) {
					Expect(dr.RespondTo).To(Equal(readFromTop.ID))
					Expect(dr.Data).To(Equal(dataReady.Data))
				}).
				Return(nil)
			bottomPort.EXPECT().Retrieve(gomock.Any())

			madeProgress := t.respond(10)

			Expect(madeProgress).To(BeTrue())
			Expect(t.inflightReqToBottom).To(HaveLen(1))
		})

		It("should respond write done", func() {
			done := mem.WriteDoneRspBuilder{}.
				WithSendTime(10).
				WithRspTo(writeToBottom.ID).
				Build()
			bottomPort.EXPECT().Peek().Return(done)
			topPort.EXPECT().Send(gomock.Any()).
				Do(func(done *mem.WriteDoneRsp) {
					Expect(done.RespondTo).To(Equal(writeFromTop.ID))
				}).
				Return(nil)
			bottomPort.EXPECT().Retrieve(gomock.Any())

			madeProgress := t.respond(10)

			Expect(madeProgress).To(BeTrue())
			Expect(t.inflightReqToBottom).To(HaveLen(1))
		})

		It("should stall if TopPort is busy", func() {
			dataReady := mem.DataReadyRspBuilder{}.
				WithSendTime(10).
				WithRspTo(readToBottom.ID).
				Build()
			bottomPort.EXPECT().Peek().Return(dataReady)
			topPort.EXPECT().Send(gomock.Any()).
				Do(func(dr *mem.DataReadyRsp) {
					Expect(dr.RespondTo).To(Equal(readFromTop.ID))
					Expect(dr.Data).To(Equal(dataReady.Data))
				}).
				Return(&akita.SendError{})

			madeProgress := t.respond(10)

			Expect(madeProgress).To(BeFalse())
			Expect(t.inflightReqToBottom).To(HaveLen(2))
		})
	})

	Context("when handling control messages", func() {
		var (
			readFromTop   *mem.ReadReq
			writeFromTop  *mem.WriteReq
			readToBottom  *mem.ReadReq
			writeToBottom *mem.WriteReq
			flushReq      *mem.ControlMsg
			restartReq    *mem.ControlMsg
		)

		BeforeEach(func() {
			readFromTop = mem.ReadReqBuilder{}.
				WithSendTime(8).
				WithAddress(0x10040).
				WithByteSize(4).
				Build()
			readToBottom = mem.ReadReqBuilder{}.
				WithSendTime(8).
				WithAddress(0x20040).
				WithByteSize(4).
				Build()
			writeFromTop = mem.WriteReqBuilder{}.
				WithSendTime(8).
				WithAddress(0x10040).
				Build()
			writeToBottom = mem.WriteReqBuilder{}.
				WithSendTime(8).
				WithAddress(0x10040).
				Build()
			flushReq = mem.ControlMsgBuilder{}.
				WithSendTime(8).
				WithDst(t.CtrlPort).
				ToDiscardTransactions().
				Build()
			restartReq = mem.ControlMsgBuilder{}.
				WithSendTime(8).
				WithDst(t.CtrlPort).
				ToRestart().
				Build()

			t.inflightReqToBottom = []reqToBottom{
				{reqFromTop: readFromTop, reqToBottom: readToBottom},
				{reqFromTop: writeFromTop, reqToBottom: writeToBottom},
			}
		})

		It("should handle flush req", func() {
			ctrlPort.EXPECT().Peek().Return(flushReq)
			ctrlPort.EXPECT().Retrieve(akita.VTimeInSec(8)).Return(flushReq)
			ctrlPort.EXPECT().Send(gomock.Any()).Return(nil)

			madeProgress := t.handleCtrlRequest(8)

			Expect(madeProgress).To(BeTrue())
			Expect(t.isFlushing).To(BeTrue())
			Expect(t.inflightReqToBottom).To(BeNil())
		})

		It("should handle restart req", func() {
			ctrlPort.EXPECT().Peek().Return(restartReq)
			ctrlPort.EXPECT().Retrieve(akita.VTimeInSec(8)).Return(restartReq)
			ctrlPort.EXPECT().Send(gomock.Any()).Return(nil)
			topPort.EXPECT().Retrieve(gomock.Any()).Return(nil)
			bottomPort.EXPECT().Retrieve(gomock.Any()).Return(nil)
			translationPort.EXPECT().Retrieve(gomock.Any()).Return(nil)

			madeProgress := t.handleCtrlRequest(8)

			Expect(madeProgress).To(BeTrue())
			Expect(t.isFlushing).To(BeFalse())
		})

	})
})
