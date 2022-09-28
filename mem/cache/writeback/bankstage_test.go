package writeback

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/akita/mem"
	"gitlab.com/akita/mem/cache"
	"gitlab.com/akita/util"
)

var _ = Describe("Bank Stage", func() {
	var (
		mockCtrl          *gomock.Controller
		cacheModule       *Cache
		dirInBuf          *MockBuffer
		writeBufferInBuf  *MockBuffer
		bs                *bankStage
		storage           *mem.Storage
		topSender         *MockBufferedSender
		writeBufferBuffer *MockBuffer
		mshrStageBuffer   *MockBuffer
		lowModuleFinder   *MockLowModuleFinder
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		dirInBuf = NewMockBuffer(mockCtrl)
		writeBufferInBuf = NewMockBuffer(mockCtrl)
		mshrStageBuffer = NewMockBuffer(mockCtrl)
		topSender = NewMockBufferedSender(mockCtrl)
		writeBufferBuffer = NewMockBuffer(mockCtrl)
		lowModuleFinder = NewMockLowModuleFinder(mockCtrl)
		storage = mem.NewStorage(4 * mem.KB)

		builder := MakeBuilder()
		cacheModule = builder.Build("cache")
		cacheModule.dirToBankBuffers = []util.Buffer{dirInBuf}
		cacheModule.writeBufferToBankBuffers = []util.Buffer{writeBufferInBuf}
		cacheModule.mshrStageBuffer = mshrStageBuffer
		cacheModule.topSender = topSender
		cacheModule.writeBufferBuffer = writeBufferBuffer
		cacheModule.lowModuleFinder = lowModuleFinder
		cacheModule.storage = storage
		cacheModule.inFlightTransactions = nil

		bs = &bankStage{
			cache:   cacheModule,
			bankID:  0,
			latency: 10,
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("No transaction running", func() {
		It("should do nothing if there is no transaction", func() {
			writeBufferInBuf.EXPECT().Pop().Return(nil)
			writeBufferBuffer.EXPECT().CanPush().Return(true)
			dirInBuf.EXPECT().Pop().Return(nil)

			ret := bs.Tick(10)

			Expect(ret).To(BeFalse())
		})

		It("should extract transactions from write buffer first", func() {
			trans := &transaction{}

			writeBufferInBuf.EXPECT().Pop().Return(trans)

			ret := bs.Tick(10)

			Expect(ret).To(BeTrue())
			Expect(bs.cycleLeft).To(Equal(10))
			Expect(bs.currentTrans).To(BeIdenticalTo(trans))
		})

		It("should stall if write buffer buffer is full", func() {
			writeBufferInBuf.EXPECT().Pop().Return(nil)
			writeBufferBuffer.EXPECT().CanPush().Return(false)

			ret := bs.Tick(10)

			Expect(ret).To(BeFalse())
		})

		It("should extract transactions from directory", func() {
			trans := &transaction{}

			writeBufferInBuf.EXPECT().Pop().Return(nil)
			writeBufferBuffer.EXPECT().CanPush().Return(true)
			dirInBuf.EXPECT().Pop().Return(trans)

			ret := bs.Tick(10)

			Expect(ret).To(BeTrue())
			Expect(bs.cycleLeft).To(Equal(10))
			Expect(bs.currentTrans).To(BeIdenticalTo(trans))
		})

		It("should directly forward fetch transaction to writebuffer", func() {
			trans := &transaction{
				action: writeBufferFetch,
			}

			writeBufferInBuf.EXPECT().Pop().Return(nil)
			writeBufferBuffer.EXPECT().CanPush().Return(true)
			writeBufferBuffer.EXPECT().Push(trans)
			dirInBuf.EXPECT().Pop().Return(trans)

			ret := bs.Tick(10)

			Expect(ret).To(BeTrue())
			Expect(bs.cycleLeft).To(Equal(0))
			Expect(bs.currentTrans).To(BeNil())
		})
	})

	Context("processing a transaction", func() {
		var (
			trans *transaction
		)

		BeforeEach(func() {
			trans = &transaction{}
			bs.currentTrans = trans
			bs.cycleLeft = 2
		})

		It("should reduce cycle left", func() {
			ret := bs.Tick(10)
			Expect(ret).To(BeTrue())
			Expect(bs.cycleLeft).To(Equal(1))
		})
	})

	Context("completing a read hit transaction", func() {
		var (
			read  *mem.ReadReq
			block *cache.Block
			trans *transaction
		)

		BeforeEach(func() {
			storage.Write(0x40, []byte{1, 2, 3, 4, 5, 6, 7, 8})
			read = mem.ReadReqBuilder{}.
				WithSendTime(6).
				WithAddress(0x104).
				WithByteSize(4).
				Build()
			block = &cache.Block{
				CacheAddress: 0x40,
				ReadCount:    1,
			}
			trans = &transaction{
				read:   read,
				block:  block,
				action: bankReadHit,
			}
			cacheModule.inFlightTransactions = append(
				cacheModule.inFlightTransactions, trans)
			bs.currentTrans = trans
			bs.cycleLeft = 0
		})

		It("should stall if send buffer is full", func() {
			topSender.EXPECT().CanSend(1).Return(false)
			ret := bs.Tick(10)
			Expect(ret).To(BeFalse())
			Expect(bs.cycleLeft).To(Equal(0))
		})

		It("should read and send response", func() {
			topSender.EXPECT().CanSend(1).Return(true)
			topSender.EXPECT().Send(gomock.Any()).
				Do(func(dr *mem.DataReadyRsp) {
					Expect(dr.RespondTo).To(Equal(read.ID))
					Expect(dr.Data).To(Equal([]byte{5, 6, 7, 8}))
				})

			ret := bs.Tick(10)

			Expect(ret).To(BeTrue())
			Expect(bs.currentTrans).To(BeNil())
			Expect(block.ReadCount).To(Equal(0))
			Expect(cacheModule.inFlightTransactions).
				NotTo(ContainElement(trans))
		})
	})

	Context("completing a write-hit transaction", func() {
		var (
			write *mem.WriteReq
			block *cache.Block
			trans *transaction
		)

		BeforeEach(func() {
			write = mem.WriteReqBuilder{}.
				WithSendTime(6).
				WithAddress(0x104).
				WithData([]byte{5, 6, 7, 8}).
				Build()
			block = &cache.Block{
				CacheAddress: 0x40,
				ReadCount:    1,
				IsLocked:     true,
			}
			trans = &transaction{
				write:  write,
				block:  block,
				action: bankWriteHit,
			}
			cacheModule.inFlightTransactions = append(
				cacheModule.inFlightTransactions, trans)
			bs.currentTrans = trans
			bs.cycleLeft = 0
		})

		It("should stall if send buffer is full", func() {
			topSender.EXPECT().CanSend(1).Return(false)
			ret := bs.Tick(10)
			Expect(ret).To(BeFalse())
			Expect(bs.cycleLeft).To(Equal(0))
		})

		It("should write and send response", func() {
			topSender.EXPECT().CanSend(1).Return(true)
			topSender.EXPECT().Send(gomock.Any()).
				Do(func(done *mem.WriteDoneRsp) {
					Expect(done.RespondTo).To(Equal(write.ID))
				})

			ret := bs.Tick(10)

			Expect(ret).To(BeTrue())
			Expect(bs.currentTrans).To(BeNil())
			data, _ := storage.Read(0x44, 4)
			Expect(data).To(Equal([]byte{5, 6, 7, 8}))
			Expect(block.IsValid).To(BeTrue())
			Expect(block.IsLocked).To(BeFalse())
			Expect(block.IsDirty).To(BeTrue())
			Expect(block.DirtyMask).To(Equal([]bool{
				false, false, false, false, true, true, true, true,
				false, false, false, false, false, false, false, false,
				false, false, false, false, false, false, false, false,
				false, false, false, false, false, false, false, false,
				false, false, false, false, false, false, false, false,
				false, false, false, false, false, false, false, false,
				false, false, false, false, false, false, false, false,
				false, false, false, false, false, false, false, false,
			}))
			Expect(cacheModule.inFlightTransactions).
				NotTo(ContainElement(trans))
		})
	})

	Context("completing a write fetched transaction", func() {
		var (
			block     *cache.Block
			mshrEntry *cache.MSHREntry
			trans     *transaction
		)

		BeforeEach(func() {
			block = &cache.Block{
				CacheAddress: 0x40,
				IsLocked:     true,
			}
			mshrEntry = &cache.MSHREntry{
				Data: []byte{
					1, 2, 3, 4, 5, 6, 7, 8,
					1, 2, 3, 4, 5, 6, 7, 8,
					1, 2, 3, 4, 5, 6, 7, 8,
					1, 2, 3, 4, 5, 6, 7, 8,
					1, 2, 3, 4, 5, 6, 7, 8,
					1, 2, 3, 4, 5, 6, 7, 8,
					1, 2, 3, 4, 5, 6, 7, 8,
					1, 2, 3, 4, 5, 6, 7, 8,
				},
				Block: block,
			}
			trans = &transaction{
				mshrEntry: mshrEntry,
				action:    bankWriteFetched,
			}
			bs.currentTrans = trans
			bs.cycleLeft = 0
		})

		It("should stall if the mshr stage buffer is full", func() {
			mshrStageBuffer.EXPECT().CanPush().Return(false)
			ret := bs.Tick(10)
			Expect(ret).To(BeFalse())
			Expect(bs.cycleLeft).To(Equal(0))
		})

		It("should write to storage and send to mshr stage", func() {
			mshrStageBuffer.EXPECT().CanPush().Return(true)
			mshrStageBuffer.EXPECT().Push(mshrEntry)

			ret := bs.Tick(10)

			Expect(ret).To(BeTrue())
			Expect(bs.currentTrans).To(BeNil())
			writtenData, _ := storage.Read(0x40, 64)
			Expect(writtenData).To(Equal(mshrEntry.Data))
			Expect(block.IsLocked).To(BeFalse())
			Expect(block.IsValid).To(BeTrue())
			Expect(bs.currentTrans).To(BeNil())
		})
	})

	Context("finalizing a read for eviction action", func() {
		var (
			victim *cache.Block
			trans  *transaction
		)

		BeforeEach(func() {
			victim = &cache.Block{
				Tag:          0x200,
				CacheAddress: 0x300,
				DirtyMask: []bool{
					true, true, true, true, false, false, false, false,
					true, true, true, true, false, false, false, false,
					true, true, true, true, false, false, false, false,
					true, true, true, true, false, false, false, false,
					true, true, true, true, false, false, false, false,
					true, true, true, true, false, false, false, false,
					true, true, true, true, false, false, false, false,
					true, true, true, true, false, false, false, false,
				},
			}
			trans = &transaction{
				victim: victim,
				action: bankEvictAndFetch,
			}
			bs.currentTrans = trans
			bs.cycleLeft = 0
		})

		It("should stall if the bottom sender is busy", func() {
			writeBufferBuffer.EXPECT().CanPush().Return(false)

			ret := bs.Tick(10)

			Expect(ret).To(BeFalse())
			Expect(bs.cycleLeft).To(Equal(0))
		})

		It("should send write to bottom", func() {
			data := []byte{
				1, 2, 3, 4, 5, 6, 7, 8,
				1, 2, 3, 4, 5, 6, 7, 8,
				1, 2, 3, 4, 5, 6, 7, 8,
				1, 2, 3, 4, 5, 6, 7, 8,
				1, 2, 3, 4, 5, 6, 7, 8,
				1, 2, 3, 4, 5, 6, 7, 8,
				1, 2, 3, 4, 5, 6, 7, 8,
				1, 2, 3, 4, 5, 6, 7, 8,
			}
			storage.Write(0x300, data)
			writeBufferBuffer.EXPECT().CanPush().Return(true)
			writeBufferBuffer.EXPECT().Push(gomock.Any()).
				Do(func(eviction *transaction) {
					Expect(eviction.action).To(Equal(writeBufferEvictAndFetch))
					Expect(eviction.evictingData).To(Equal(data))
				})

			ret := bs.Tick(10)
			Expect(ret).To(BeTrue())
			Expect(bs.currentTrans).To(BeNil())
		})

	})
})
