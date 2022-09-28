package mmu

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem/device"
)

var _ = Describe("MMU", func() {

	var (
		mockCtrl      *gomock.Controller
		engine        *MockEngine
		toTop         *MockPort
		migrationPort *MockPort
		topSender     *MockBufferedSender
		pageTable     *device.PageTableImpl
		mmu           *MMUImpl
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		engine = NewMockEngine(mockCtrl)
		toTop = NewMockPort(mockCtrl)
		migrationPort = NewMockPort(mockCtrl)
		topSender = NewMockBufferedSender(mockCtrl)
		pageTable = device.NewPageTable(12)

		builder := MakeBuilder().WithEngine(engine)
		mmu = builder.Build("mmu")
		mmu.ToTop = toTop
		mmu.topSender = topSender
		mmu.MigrationPort = migrationPort
		mmu.pageTable = pageTable
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("parse top", func() {
		It("should process translation request", func() {
			translationReq := device.TranslationReqBuilder{}.
				WithSendTime(10).
				WithDst(mmu.ToTop).
				WithPID(1).
				WithVAddr(0x100000100).
				WithDeviceID(0).
				Build()
			toTop.EXPECT().
				Retrieve(akita.VTimeInSec(10)).
				Return(translationReq)

			mmu.parseFromTop(10)

			Expect(mmu.walkingTranslations).To(HaveLen(1))

		})

		It("should stall parse from top if MMU is servicing max requests", func() {
			mmu.walkingTranslations = make([]transaction, 16)

			madeProgress := mmu.parseFromTop(10)

			Expect(madeProgress).To(BeFalse())
		})
	})

	Context("walk page table", func() {
		It("should reduce translation cycles", func() {
			req := device.TranslationReqBuilder{}.
				WithSendTime(10).
				WithDst(toTop).
				WithPID(1).
				WithVAddr(0x1020).
				WithDeviceID(0).
				Build()
			walking := transaction{req: req, cycleLeft: 10}
			mmu.walkingTranslations = append(mmu.walkingTranslations, walking)

			madeProgress := mmu.walkPageTable(11)

			Expect(mmu.walkingTranslations[0].cycleLeft).To(Equal(9))
			Expect(madeProgress).To(BeTrue())
		})

		It("should send rsp to top if hit", func() {
			page := device.Page{
				PID:      1,
				VAddr:    0x1000,
				PAddr:    0x0,
				PageSize: 4096,
				Valid:    true,
			}
			req := device.TranslationReqBuilder{}.
				WithSendTime(10).
				WithDst(mmu.ToTop).
				WithPID(1).
				WithVAddr(0x1000).
				WithDeviceID(0).
				Build()
			walking := transaction{req: req, cycleLeft: 0}
			mmu.walkingTranslations = append(mmu.walkingTranslations, walking)

			// pageTable.EXPECT().
			// 	Find(ca.PID(1), uint64(0x1000)).
			// 	Return(page, true)
			topSender.EXPECT().CanSend(1).Return(true)
			topSender.EXPECT().
				Send(gomock.Any()).
				Do(func(rsp *device.TranslationRsp) {
					Expect(rsp.Page).To(Equal(page))
				})

			madeProgress := mmu.walkPageTable(11)

			Expect(madeProgress).To(BeTrue())
			Expect(mmu.walkingTranslations).To(HaveLen(0))
		})

		It("should stall if cannot send to top", func() {
			// page := device.Page{
			// 	PID:      1,
			// 	VAddr:    0x1000,
			// 	PAddr:    0x0,
			// 	PageSize: 4096,
			// 	Valid:    true,
			// }
			req := device.TranslationReqBuilder{}.
				WithSendTime(10).
				WithDst(mmu.ToTop).
				WithPID(1).
				WithVAddr(0x1000).
				WithDeviceID(0).
				Build()
			walking := transaction{req: req, cycleLeft: 0}
			mmu.walkingTranslations =
				append(mmu.walkingTranslations, walking)

			// pageTable.EXPECT().
			// 	Find(ca.PID(1), uint64(0x1000)).
			// 	Return(page, true)
			topSender.EXPECT().CanSend(1).Return(false)

			madeProgress := mmu.walkPageTable(11)

			Expect(madeProgress).To(BeFalse())
		})
	})

	Context("migration", func() {
		var (
			page    device.Page
			req     *device.TranslationReq
			walking transaction
		)

		BeforeEach(func() {
			page = device.Page{
				PID:      1,
				VAddr:    0x1000,
				PAddr:    0x0,
				PageSize: 4096,
				Valid:    true,
				DeviceID: 2,
				Unified:  true,
			}
			// pageTable.EXPECT().
			// 	Find(ca.PID(1), uint64(0x1000)).
			// 	Return(page, true).
			// 	AnyTimes()
			req = device.TranslationReqBuilder{}.
				WithSendTime(10).
				WithDst(mmu.ToTop).
				WithPID(1).
				WithVAddr(0x1000).
				WithDeviceID(0).
				Build()
			walking = transaction{
				req:       req,
				page:      page,
				cycleLeft: 0,
			}
		})

		It("should be placed in the migration queue", func() {
			mmu.walkingTranslations = append(mmu.walkingTranslations, walking)

			updatedPage := page
			updatedPage.IsMigrating = true
			// pageTable.EXPECT().Update(updatedPage)

			madeProgress := mmu.walkPageTable(11)

			Expect(madeProgress).To(BeTrue())
			Expect(mmu.walkingTranslations).To(HaveLen(0))
			Expect(mmu.migrationQueue).To(HaveLen(1))
		})

		It("should place the page in the migration queue if the page is being migrated", func() {
			req.PID = 2
			page.PID = 2
			page.IsMigrating = true
			// pageTable.EXPECT().
			// 	Find(ca.PID(2), uint64(0x1000)).
			// 	Return(page, true)
			mmu.walkingTranslations = append(mmu.walkingTranslations, walking)

			// pageTable.EXPECT().Update(page)

			madeProgress := mmu.walkPageTable(11)

			Expect(madeProgress).To(BeTrue())
			Expect(mmu.walkingTranslations).To(HaveLen(0))
			Expect(mmu.migrationQueue).To(HaveLen(1))
		})

		It("should not send to driver if migration queue is empty", func() {
			madeProgress := mmu.sendMigrationToDriver(11)

			Expect(madeProgress).To(BeFalse())
		})

		It("should wait if mmu is waiting for a migration to finish", func() {
			mmu.migrationQueue = append(mmu.migrationQueue, walking)
			mmu.isDoingMigration = true

			madeProgress := mmu.sendMigrationToDriver(11)

			Expect(madeProgress).To(BeFalse())
			Expect(mmu.migrationQueue).To(ContainElement(walking))
		})

		It("should stall if send failed", func() {
			mmu.migrationQueue = append(mmu.migrationQueue, walking)

			migrationPort.EXPECT().
				Send(gomock.Any()).
				Return(akita.NewSendError())

			madeProgress := mmu.sendMigrationToDriver(11)

			Expect(madeProgress).To(BeFalse())
			Expect(mmu.migrationQueue).To(ContainElement(walking))
		})

		It("should send migration request", func() {
			mmu.migrationQueue = append(mmu.migrationQueue, walking)

			migrationPort.EXPECT().
				Send(gomock.Any()).
				Return(nil)
			updatedPage := page
			updatedPage.IsMigrating = true
			// pageTable.EXPECT().Update(updatedPage)

			madeProgress := mmu.sendMigrationToDriver(11)

			Expect(madeProgress).To(BeTrue())
			Expect(mmu.migrationQueue).NotTo(ContainElement(walking))
			Expect(mmu.isDoingMigration).To(BeTrue())
		})

		It("should reply to the GPU if the page is already on the destination GPU", func() {
			walking.req.DeviceID = 2
			mmu.migrationQueue = append(mmu.migrationQueue, walking)

			updatedPage := page
			updatedPage.IsMigrating = false
			// pageTable.EXPECT().Update(updatedPage)
			topSender.EXPECT().Send(gomock.Any())

			madeProgress := mmu.sendMigrationToDriver(11)

			Expect(madeProgress).To(BeTrue())
			Expect(mmu.migrationQueue).NotTo(ContainElement(walking))
			Expect(mmu.isDoingMigration).To(BeFalse())
		})
	})

	// Context("when received migrated page information", func() {
	// 	var (
	// 		// page          device.Page
	// 		req           *device.TranslationReq
	// 		migrating     transaction
	// 		migrationDone *device.PageMigrationRspFromDriver
	// 	)

	// 	BeforeEach(func() {
	// 		// page = device.Page{
	// 		// 	PID:         1,
	// 		// 	VAddr:       0x1000,
	// 		// 	PAddr:       0x0,
	// 		// 	PageSize:    4096,
	// 		// 	Valid:       true,
	// 		// 	DeviceID:    1,
	// 		// 	Unified:     true,
	// 		// 	IsMigrating: true,
	// 		// }
	// 		// pageTable.EXPECT().
	// 		// 	Find(ca.PID(1), uint64(0x1000)).
	// 		// 	Return(page, true).
	// 		// 	AnyTimes()
	// 		req = device.TranslationReqBuilder{}.
	// 			WithSendTime(10).
	// 			WithDst(mmu.ToTop).
	// 			WithPID(1).
	// 			WithVAddr(0x1000).
	// 			WithDeviceID(0).
	// 			Build()
	// 		migrating = transaction{req: req, cycleLeft: 0}
	// 		mmu.currentOnDemandMigration = migrating
	// 		migrationDone = device.NewPageMigrationRspFromDriver(0, nil, nil)
	// 	})

	// 	// It("should do nothing if no respond", func() {
	// 	// 	migrationPort.EXPECT().Peek().Return(nil)

	// 	// 	madeProgress := mmu.processMigrationReturn(10)

	// 	// 	Expect(madeProgress).To(BeFalse())
	// 	// })

	// 	// It("should stall if send to top failed", func() {
	// 	// 	migrationPort.EXPECT().Peek().Return(migrationDone)
	// 	// 	topSender.EXPECT().CanSend(1).Return(false)

	// 	// 	madeProgress := mmu.processMigrationReturn(10)

	// 	// 	Expect(madeProgress).To(BeFalse())
	// 	// 	Expect(mmu.isDoingMigration).To(BeFalse())
	// 	// })

	// 	// It("should send rsp to top", func() {
	// 	// 	migrationPort.EXPECT().Peek().Return(migrationDone)
	// 	// 	topSender.EXPECT().CanSend(1).Return(true)
	// 	// 	topSender.EXPECT().Send(gomock.Any()).
	// 	// 		Do(func(rsp *device.TranslationRsp) {
	// 	// 			Expect(rsp.Page).To(Equal(page))
	// 	// 		})
	// 	// 	migrationPort.EXPECT().Retrieve(gomock.Any())

	// 	// 	updatedPage := page
	// 	// 	updatedPage.IsMigrating = false
	// 	// 	// pageTable.EXPECT().Update(updatedPage)

	// 	// 	updatedPage.IsPinned = true
	// 	// 	// pageTable.EXPECT().Update(updatedPage)

	// 	// 	madeProgress := mmu.processMigrationReturn(10)

	// 	// 	Expect(madeProgress).To(BeTrue())
	// 	// 	Expect(mmu.isDoingMigration).To(BeFalse())
	// 	// })

	// })
})

var _ = Describe("MMU Integration", func() {
	var (
		mockCtrl   *gomock.Controller
		engine     akita.Engine
		mmu        *MMUImpl
		agent      *MockPort
		connection akita.Connection
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		engine = akita.NewSerialEngine()

		builder := MakeBuilder().WithEngine(engine)
		mmu = builder.Build("mmu")
		agent = NewMockPort(mockCtrl)
		connection = akita.NewDirectConnection("conn", engine, 1*akita.GHz)

		agent.EXPECT().SetConnection(connection)
		connection.PlugIn(agent, 10)
		connection.PlugIn(mmu.ToTop, 10)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should lookup", func() {
		page := device.Page{
			PID:      1,
			VAddr:    0x1000,
			PAddr:    0x2000,
			PageSize: 4096,
			Valid:    true,
			DeviceID: 1,
		}
		mmu.pageTable.Insert(page)

		req := device.TranslationReqBuilder{}.
			WithSendTime(10).
			WithSrc(agent).
			WithDst(mmu.ToTop).
			WithPID(1).
			WithVAddr(0x1000).
			WithDeviceID(0).
			Build()
		req.RecvTime = 10
		mmu.ToTop.Recv(req)

		agent.EXPECT().Recv(gomock.Any()).
			Do(func(rsp *device.TranslationRsp) {
				Expect(rsp.Page).To(Equal(page))
				Expect(rsp.RespondTo).To(Equal(req.ID))
			})

		engine.Run()
	})
})
