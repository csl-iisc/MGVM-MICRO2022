package tlb

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem/cache"
	"gitlab.com/akita/mem/device"
	"gitlab.com/akita/mem/vm/tlb/internal"
	"gitlab.com/akita/util/ca"
)

var _ = Describe("TLB", func() {

	var (
		mockCtrl        *gomock.Controller
		engine          *MockEngine
		tlb             *TLB
		set             *MockSet
		topPort         *MockPort
		bottomPort      *MockPort
		controlPort     *MockPort
		lowModuleFinder *cache.SingleLowModuleFinder
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		engine = NewMockEngine(mockCtrl)
		set = NewMockSet(mockCtrl)
		topPort = NewMockPort(mockCtrl)
		bottomPort = NewMockPort(mockCtrl)
		controlPort = NewMockPort(mockCtrl)
		lowModuleFinder = &cache.SingleLowModuleFinder{}

		tlb = MakeBuilder().WithEngine(engine).Build("tlb")
		tlb.SetLowModuleFinder(lowModuleFinder)
		tlb.TopPort = topPort
		tlb.BottomPort = bottomPort
		tlb.ControlPort = controlPort
		tlb.Sets = []internal.Set{set}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should do nothing if there is no req in TopPort", func() {
		topPort.EXPECT().Peek().Return(nil)

		madeProgress := tlb.lookup(10)

		Expect(madeProgress).To(BeFalse())
	})

	Context("hit", func() {
		var (
			wayID int
			page  device.Page
			req   *device.TranslationReq
		)

		BeforeEach(func() {
			wayID = 1
			page = device.Page{
				PID:   1,
				VAddr: 0x100,
				PAddr: 0x200,
				Valid: true,
			}
			set.EXPECT().Lookup(ca.PID(1), uint64(0x100)).
				Return(wayID, page, true)

			req = device.TranslationReqBuilder{}.
				WithSendTime(5).
				WithPID(1).
				WithVAddr(uint64(0x100)).
				WithDeviceID(1).
				Build()
		})

		It("should respond to top", func() {
			topPort.EXPECT().Peek().Return(req)
			topPort.EXPECT().Retrieve(gomock.Any())
			topPort.EXPECT().Send(gomock.Any())

			set.EXPECT().Visit(wayID)

			madeProgress := tlb.lookup(10)

			Expect(madeProgress).To(BeTrue())
		})

		It("should stall if cannot send to top", func() {
			topPort.EXPECT().Peek().Return(req)
			topPort.EXPECT().Send(gomock.Any()).
				Return(&akita.SendError{})

			madeProgress := tlb.lookup(10)

			Expect(madeProgress).To(BeFalse())
		})
	})

	Context("miss", func() {
		var (
			wayID int
			page  device.Page
			req   *device.TranslationReq
		)

		BeforeEach(func() {
			wayID = 1
			page = device.Page{
				PID:   1,
				VAddr: 0x100,
				PAddr: 0x200,
				Valid: false,
			}
			set.EXPECT().
				Lookup(ca.PID(1), uint64(0x100)).
				Return(wayID, page, true).
				AnyTimes()

			req = device.TranslationReqBuilder{}.
				WithSendTime(5).
				WithPID(1).
				WithVAddr(0x100).
				WithDeviceID(1).
				Build()
		})

		It("should fetch from bottom and add entry to MSHR", func() {
			topPort.EXPECT().Peek().Return(req)
			topPort.EXPECT().Retrieve(gomock.Any())
			bottomPort.EXPECT().Send(gomock.Any()).
				Do(func(req *device.TranslationReq) {
					Expect(req.VAddr).To(Equal(uint64(0x100)))
					Expect(req.PID).To(Equal(ca.PID(1)))
					Expect(req.DeviceID).To(Equal(uint64(1)))
				}).
				Return(nil)

			madeProgress := tlb.lookup(10)

			Expect(madeProgress).To(BeTrue())
			Expect(tlb.mshr.IsEntryPresent(ca.PID(1), uint64(0x100))).To(Equal(true))
		})

		It("should find the entry in MSHR and not request from bottom", func() {
			tlb.mshr.Add(1, 0x100)
			topPort.EXPECT().Peek().Return(req)
			topPort.EXPECT().Retrieve(gomock.Any())

			madeProgress := tlb.lookup(10)
			Expect(tlb.mshr.IsEntryPresent(ca.PID(1), uint64(0x100))).
				To(Equal(true))
			Expect(madeProgress).To(BeTrue())
		})

		It("should stall if bottom is busy", func() {
			topPort.EXPECT().Peek().Return(req)
			bottomPort.EXPECT().Send(gomock.Any()).
				Return(&akita.SendError{})

			madeProgress := tlb.lookup(10)

			Expect(madeProgress).To(BeFalse())
		})
	})

	Context("parse bottom", func() {
		var (
			wayID       int
			req         *device.TranslationReq
			fetchBottom *device.TranslationReq
			page        device.Page
			rsp         *device.TranslationRsp
		)

		BeforeEach(func() {
			wayID = 1
			req = device.TranslationReqBuilder{}.
				WithSendTime(5).
				WithPID(1).
				WithVAddr(0x100).
				WithDeviceID(1).
				Build()
			fetchBottom = device.TranslationReqBuilder{}.
				WithSendTime(5).
				WithPID(1).
				WithVAddr(0x100).
				WithDeviceID(1).
				Build()
			page = device.Page{
				PID:   1,
				VAddr: 0x100,
				PAddr: 0x200,
				Valid: true,
			}
			rsp = device.TranslationRspBuilder{}.
				WithSendTime(5).
				WithRspTo(fetchBottom.ID).
				WithPage(page).
				Build()
		})

		It("should do nothing if no return", func() {
			bottomPort.EXPECT().Peek().Return(nil)

			madeProgress := tlb.parseBottom(10)

			Expect(madeProgress).To(BeFalse())
		})

		It("should stall if the TLB is responding to an MSHR entry", func() {
			mshrEntry := tlb.mshr.Add(1, 0x100)
			mshrEntry.Requests = append(mshrEntry.Requests, req)
			tlb.respondingMSHREntry = mshrEntry

			madeProgress := tlb.parseBottom(10)

			Expect(madeProgress).To(BeFalse())
		})

		It("should parse respond from bottom", func() {
			bottomPort.EXPECT().Peek().Return(rsp)
			bottomPort.EXPECT().Retrieve(gomock.Any())
			mshrEntry := tlb.mshr.Add(1, 0x100)
			mshrEntry.Requests = append(mshrEntry.Requests, req)

			set.EXPECT().Evict().Return(wayID, true)
			set.EXPECT().Update(wayID, page)
			set.EXPECT().Visit(wayID)

			// topPort.EXPECT().Send(gomock.Any()).
			// 	Do(func(rsp *vm.TranslationRsp) {
			// 		Expect(rsp.Page).To(Equal(page))
			// 		Expect(rsp.RespondTo).To(Equal(req.ID))
			// 	})

			madeProgress := tlb.parseBottom(10)

			Expect(madeProgress).To(BeTrue())
			Expect(tlb.respondingMSHREntry).NotTo(BeNil())
			Expect(tlb.mshr.IsEntryPresent(ca.PID(1), uint64(0x100))).
				To(Equal(false))
		})

		It("should respond", func() {
			mshrEntry := tlb.mshr.Add(1, 0x100)
			mshrEntry.Requests = append(mshrEntry.Requests, req)
			tlb.respondingMSHREntry = mshrEntry

			topPort.EXPECT().Send(gomock.Any()).Return(nil)

			madeProgress := tlb.respondMSHREntry(10)

			Expect(madeProgress).To(BeTrue())
			Expect(mshrEntry.Requests).To(HaveLen(0))
			Expect(tlb.respondingMSHREntry).To(BeNil())
		})
	})

	Context("flush related handling", func() {
		var (
		// flushReq   *TLBFlushReq
		// restartReq *TLBRestartReq
		)

		BeforeEach(func() {

			// restartReq = TLBRestartReqBuilder{}.
			// 	WithSrc(nil).
			// 	WithDst(nil).
			// 	WithSendTime(10).
			// 	Build()
		})

		It("should do nothing if no req", func() {
			controlPort.EXPECT().Peek().Return(nil)
			madeProgress := tlb.performCtrlReq(10)
			Expect(madeProgress).To(BeFalse())
		})

		It("should handle flush request", func() {
			flushReq := TLBFlushReqBuilder{}.
				WithSrc(nil).
				WithDst(nil).
				WithSendTime(10).
				WithVAddrs([]uint64{0x1000}).
				WithPID(1).
				Build()
			page := device.Page{
				PID:   1,
				VAddr: 0x1000,
				Valid: true,
			}
			wayID := 1

			set.EXPECT().Lookup(ca.PID(1), uint64(0x1000)).
				Return(wayID, page, true)
			set.EXPECT().Update(wayID, device.Page{
				PID:   1,
				VAddr: 0x1000,
				Valid: false,
			})
			controlPort.EXPECT().Peek().Return(flushReq)
			controlPort.EXPECT().Retrieve(akita.VTimeInSec(10)).Return(flushReq)
			controlPort.EXPECT().Send(gomock.Any())

			madeProgress := tlb.performCtrlReq(10)

			Expect(madeProgress).To(BeTrue())
			Expect(tlb.isPaused).To(BeTrue())
		})

		It("should handle restart request", func() {
			restartReq := TLBRestartReqBuilder{}.
				WithSrc(nil).
				WithDst(nil).
				WithSendTime(10).
				Build()
			controlPort.EXPECT().Peek().
				Return(restartReq)
			controlPort.EXPECT().Retrieve(akita.VTimeInSec(10)).
				Return(restartReq)
			controlPort.EXPECT().Send(gomock.Any())
			topPort.EXPECT().Retrieve(gomock.Any()).Return(nil)
			bottomPort.EXPECT().Retrieve(gomock.Any()).Return(nil)

			madeProgress := tlb.performCtrlReq(10)

			Expect(madeProgress).To(BeTrue())
			Expect(tlb.isPaused).To(BeFalse())
		})
	})
})

var _ = Describe("TLB Integration", func() {
	var (
		mockCtrl        *gomock.Controller
		engine          akita.Engine
		tlb             *TLB
		lowModule       *MockPort
		lowModuleFinder *cache.SingleLowModuleFinder
		agent           *MockPort
		connection      akita.Connection
		page            device.Page
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		engine = akita.NewSerialEngine()
		lowModule = NewMockPort(mockCtrl)
		lowModuleFinder = &cache.SingleLowModuleFinder{}
		lowModuleFinder.LowModule = lowModule
		agent = NewMockPort(mockCtrl)
		connection = akita.NewDirectConnection("Conn", engine, 1*akita.GHz)
		tlb = MakeBuilder().WithEngine(engine).Build("tlb")
		tlb.LowModule = lowModule
		tlb.SetLowModuleFinder(lowModuleFinder)

		agent.EXPECT().SetConnection(connection)
		lowModule.EXPECT().SetConnection(connection)
		connection.PlugIn(agent, 10)
		connection.PlugIn(lowModule, 10)
		connection.PlugIn(tlb.TopPort, 10)
		connection.PlugIn(tlb.BottomPort, 10)
		connection.PlugIn(tlb.ControlPort, 10)

		page = device.Page{
			PID:   1,
			VAddr: 0x1000,
			PAddr: 0x2000,
			Valid: true,
		}
		lowModule.EXPECT().Recv(gomock.Any()).
			Do(func(req *device.TranslationReq) {
				rsp := device.TranslationRspBuilder{}.
					WithSendTime(req.RecvTime + 1).
					WithSrc(lowModule).
					WithDst(req.Src).
					WithPage(page).
					WithRspTo(req.ID).
					Build()
				connection.Send(rsp)
			}).
			AnyTimes()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should do tlb miss", func() {
		req := device.TranslationReqBuilder{}.
			WithSendTime(10).
			WithSrc(agent).
			WithDst(tlb.TopPort).
			WithPID(1).
			WithVAddr(0x1000).
			WithDeviceID(1).
			Build()
		req.RecvTime = 10
		tlb.TopPort.Recv(req)

		agent.EXPECT().Recv(gomock.Any()).
			Do(func(rsp *device.TranslationRsp) {
				Expect(rsp.Page).To(Equal(page))
			})

		engine.Run()
	})

	It("should have faster hit than miss", func() {
		time1 := akita.VTimeInSec(10)
		req := device.TranslationReqBuilder{}.
			WithSendTime(time1).
			WithSrc(agent).
			WithDst(tlb.TopPort).
			WithPID(1).
			WithVAddr(0x1000).
			WithDeviceID(1).
			Build()
		req.RecvTime = time1
		tlb.TopPort.Recv(req)

		agent.EXPECT().Recv(gomock.Any()).
			Do(func(rsp *device.TranslationRsp) {
				Expect(rsp.Page).To(Equal(page))
			})

		engine.Run()

		time2 := engine.CurrentTime()

		req.RecvTime = time2
		tlb.TopPort.Recv(req)

		agent.EXPECT().Recv(gomock.Any()).
			Do(func(rsp *device.TranslationRsp) {
				Expect(rsp.Page).To(Equal(page))
			})

		engine.Run()

		time3 := engine.CurrentTime()

		Expect(time3 - time2).To(BeNumerically("<", time2-time1))
	})

	/*It("should have miss after shootdown ", func() {
		time1 := akita.VTimeInSec(10)
		req := vm.NewTranslationReq(time1, agent, tlb.TopPort, 1, 0x1000, 1)
		req.SetRecvTime(time1)
		tlb.TopPort.Recv(*req)
		agent.EXPECT().Recv(gomock.Any()).
			Do(func(rsp vm.TranslationReadyRsp) {
				Expect(rsp.Page).To(Equal(&page))
			})
		engine.Run()

		time2 := engine.CurrentTime()
		shootdownReq := vm.NewPTEInvalidationReq(
			time2, agent, tlb.ControlPort, 1, []uint64{0x1000})
		shootdownReq.SetRecvTime(time2)
		tlb.ControlPort.Recv(*shootdownReq)
		agent.EXPECT().Recv(gomock.Any()).
			Do(func(rsp vm.InvalidationCompleteRsp) {
				Expect(rsp.RespondTo).To(Equal(shootdownReq.ID))
			})
		engine.Run()

		time3 := engine.CurrentTime()
		req.SetRecvTime(time3)
		tlb.TopPort.Recv(*req)
		agent.EXPECT().Recv(gomock.Any()).
			Do(func(rsp vm.TranslationReadyRsp) {
				Expect(rsp.Page).To(Equal(&page))
			})
		engine.Run()
		time4 := engine.CurrentTime()

		Expect(time4 - time3).To(BeNumerically("~", time2-time1))
	})*/

})
