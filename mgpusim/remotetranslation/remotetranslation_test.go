package remotetranslation

import (
	"log"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem/cache"
	"gitlab.com/akita/mem/device"
)

//go:generate mockgen -destination "mock_akita_test.go" -package remotetranslation -write_package_comment=false gitlab.com/akita/akita Port,Engine

func TestRemoteTranslationUnit(t *testing.T) {
	log.SetOutput(GinkgoWriter)
	RegisterFailHandler(Fail)
	RunSpecs(t, "RemoteAddressTranslation")
}

var _ = Describe("RemoteTranslationUnit", func() {
	var (
		mockCtrl *gomock.Controller

		engine                *MockEngine
		remoteTranslationUnit *RemoteTranslationUnit
		toL1                  *MockPort
		toL2                  *MockPort
		toOutside             *MockPort
		localModules          *cache.SingleLowModuleFinder
		remoteModules         *cache.SingleLowModuleFinder
		localTLB              *MockPort
		remoteTLB             *MockPort
		numTransPerCycle      int
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())

		engine = NewMockEngine(mockCtrl)
		localTLB = NewMockPort(mockCtrl)
		remoteTLB = NewMockPort(mockCtrl)
		localModules = new(cache.SingleLowModuleFinder)
		localModules.LowModule = localTLB
		remoteModules = new(cache.SingleLowModuleFinder)
		remoteModules.LowModule = remoteTLB

		remoteTranslationUnit = NewRemoteTranslationUnit("RemoteTranslationUnit",
			engine, localModules, remoteModules)

		toL1 = NewMockPort(mockCtrl)
		toL2 = NewMockPort(mockCtrl)
		toOutside = NewMockPort(mockCtrl)
		remoteTranslationUnit.ToL1 = toL1
		remoteTranslationUnit.ToL2 = toL2
		remoteTranslationUnit.ToOutside = toOutside

		numTransPerCycle = 6

	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("Read from inside", func() {
		var read *device.TranslationReq

		BeforeEach(func() {
			read = device.TranslationReqBuilder{}.
				WithSendTime(6).
				WithSrc(localTLB).
				WithDst(remoteTranslationUnit.ToOutside).
				WithVAddr(0x100).
				WithPID(42).
				Build()
		})

		It("should send read to outside", func() {
			toL1.EXPECT().Peek().Return(read).Times(numTransPerCycle)
			toOutside.EXPECT().
				Send(gomock.AssignableToTypeOf(&device.TranslationReq{})).
				Return(nil)
			toL1.EXPECT().Retrieve(akita.VTimeInSec(10)).Return(read)

			remoteTranslationUnit.processFromL1(10)

			Expect(remoteTranslationUnit.transactionsFromInside).To(HaveLen(1))
		})

		It("should wait if outside connection is busy", func() {
			toL1.EXPECT().Peek().Return(read).Times(numTransPerCycle)
			toOutside.EXPECT().
				Send(gomock.AssignableToTypeOf(&device.TranslationReq{})).
				Return(akita.NewSendError()).Times(numTransPerCycle)

			remoteTranslationUnit.processFromL1(10)

			Expect(remoteTranslationUnit.transactionsFromInside).To(HaveLen(0))
		})
	})

	Context("Read from outside", func() {
		var read *device.TranslationReq

		BeforeEach(func() {
			read = device.TranslationReqBuilder{}.
				WithSendTime(6).
				WithSrc(localTLB).
				WithDst(remoteTranslationUnit.ToOutside).
				WithVAddr(0x100).
				WithPID(42).
				Build()
		})

		It("should send read to outside", func() {
			toOutside.EXPECT().Peek().Return(read).Times(numTransPerCycle)
			toL2.EXPECT().
				Send(gomock.AssignableToTypeOf(&device.TranslationReq{})).
				Return(nil)
			toOutside.EXPECT().Retrieve(akita.VTimeInSec(10)).Return(read)

			remoteTranslationUnit.processFromOutside(10)

			Expect(remoteTranslationUnit.transactionsFromOutside).To(HaveLen(1))
		})

		It("should wait if outside connection is busy", func() {
			toOutside.EXPECT().Peek().Return(read).Times(numTransPerCycle)
			toL2.EXPECT().
				Send(gomock.AssignableToTypeOf(&device.TranslationReq{})).
				Return(akita.NewSendError()).Times(numTransPerCycle)

			remoteTranslationUnit.processFromOutside(10)

			Expect(remoteTranslationUnit.transactionsFromInside).To(HaveLen(0))
		})
	})

	Context("DataReady from outside", func() {
		var (
			readFromInside *device.TranslationReq
			read           *device.TranslationReq
			rsp            *device.TranslationRsp
		)

		BeforeEach(func() {
			readFromInside = device.TranslationReqBuilder{}.
				WithSendTime(4).
				WithSrc(localTLB).
				WithDst(remoteTranslationUnit.ToL1).
				WithVAddr(0x100).
				WithPID(42).
				Build()
			read = device.TranslationReqBuilder{}.
				WithSendTime(6).
				WithSrc(remoteTranslationUnit.ToOutside).
				WithDst(remoteTLB).
				WithVAddr(0x100).
				WithPID(42).
				Build()
			rsp = device.TranslationRspBuilder{}.
				WithSendTime(9).
				WithSrc(remoteTLB).
				WithDst(remoteTranslationUnit.ToOutside).
				WithRspTo(read.ID).
				Build()

			remoteTranslationUnit.transactionsFromInside = append(
				remoteTranslationUnit.transactionsFromInside,
				transaction{
					fromInside: readFromInside,
					toOutside:  read,
				})
		})

		It("should send rsp to inside", func() {
			toOutside.EXPECT().Peek().Return(rsp).Times(numTransPerCycle)
			toL2.EXPECT().
				Send(gomock.AssignableToTypeOf(&device.TranslationRsp{})).
				Return(nil)
			toOutside.EXPECT().Retrieve(akita.VTimeInSec(10)).Return(read)

			remoteTranslationUnit.processFromOutside(10)

			Expect(remoteTranslationUnit.transactionsFromInside).To(HaveLen(0))
		})

		It("should not send rsp to inside if busy", func() {
			toOutside.EXPECT().Peek().Return(rsp).Times(numTransPerCycle)
			toL2.EXPECT().
				Send(gomock.AssignableToTypeOf(&device.TranslationRsp{})).
				Return(akita.NewSendError()).Times(numTransPerCycle)

			remoteTranslationUnit.processFromOutside(10)

			Expect(remoteTranslationUnit.transactionsFromInside).To(HaveLen(1))
		})
	})

	Context("DataReady from inside", func() {
		var (
			readFromOutside *device.TranslationReq
			read            *device.TranslationReq
			rsp             *device.TranslationRsp
		)

		BeforeEach(func() {
			readFromOutside = device.TranslationReqBuilder{}.
				WithSendTime(4).
				WithSrc(localTLB).
				WithDst(remoteTranslationUnit.ToL2).
				WithVAddr(0x100).
				WithPID(42).
				Build()
			read = device.TranslationReqBuilder{}.
				WithSendTime(6).
				WithSrc(remoteTranslationUnit.ToOutside).
				WithDst(remoteTLB).
				WithVAddr(0x100).
				WithPID(42).
				Build()
			rsp = device.TranslationRspBuilder{}.
				WithSendTime(9).
				WithSrc(remoteTLB).
				WithDst(remoteTranslationUnit.ToOutside).
				WithRspTo(read.ID).
				Build()
			remoteTranslationUnit.transactionsFromOutside = append(
				remoteTranslationUnit.transactionsFromInside,
				transaction{
					fromOutside: readFromOutside,
					toInside:    read,
				})
		})

		It("should send rsp to outside", func() {
			toL2.EXPECT().Peek().Return(rsp).Times(numTransPerCycle)
			toOutside.EXPECT().
				Send(gomock.AssignableToTypeOf(&device.TranslationRsp{})).
				Return(nil)
			toL2.EXPECT().Retrieve(akita.VTimeInSec(10)).Return(read)

			remoteTranslationUnit.processFromL2(10)

			Expect(remoteTranslationUnit.transactionsFromOutside).To(HaveLen(0))
		})

		It("should  not send rsp to outside", func() {
			toL2.EXPECT().Peek().Return(rsp).Times(numTransPerCycle)
			toOutside.EXPECT().
				Send(gomock.AssignableToTypeOf(&device.TranslationRsp{})).
				Return(akita.NewSendError()).Times(numTransPerCycle)

			remoteTranslationUnit.processFromL2(10)

			Expect(remoteTranslationUnit.transactionsFromOutside).To(HaveLen(1))
		})
	})
})
