package writeback

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/akita/mem"
	"gitlab.com/akita/mem/cache"
	"gitlab.com/akita/util"
	ca "gitlab.com/akita/util/ca"
)

var _ = Describe("DirectoryStage", func() {

	var (
		mockCtrl          *gomock.Controller
		ds                *directoryStage
		cacheModule       *Cache
		mshr              *MockMSHR
		dirBuf            *MockBuffer
		directory         *MockDirectory
		bankBuf           *MockBuffer
		writeBufferBuffer *MockBuffer
		lowModuleFinder   *MockLowModuleFinder
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		dirBuf = NewMockBuffer(mockCtrl)
		mshr = NewMockMSHR(mockCtrl)
		directory = NewMockDirectory(mockCtrl)
		directory.EXPECT().WayAssociativity().Return(4).AnyTimes()
		writeBufferBuffer = NewMockBuffer(mockCtrl)
		bankBuf = NewMockBuffer(mockCtrl)
		lowModuleFinder = NewMockLowModuleFinder(mockCtrl)

		builder := MakeBuilder()
		cacheModule = builder.Build("cache")
		cacheModule.dirStageBuffer = dirBuf
		cacheModule.mshr = mshr
		cacheModule.directory = directory
		cacheModule.writeBufferBuffer = writeBufferBuffer
		cacheModule.dirToBankBuffers = []util.Buffer{bankBuf}
		cacheModule.lowModuleFinder = lowModuleFinder

		ds = &directoryStage{
			cache: cacheModule,
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should return if no transaction", func() {
		dirBuf.EXPECT().Peek().Return(nil)
		ret := ds.Tick(10)
		Expect(ret).To(BeFalse())
	})

	Context("read", func() {
		var (
			read  *mem.ReadReq
			trans *transaction
		)

		BeforeEach(func() {
			read = mem.ReadReqBuilder{}.
				WithSendTime(10).
				WithAddress(0x100).
				WithPID(1).
				WithByteSize(64).
				Build()
			trans = &transaction{
				read: read,
			}
			dirBuf.EXPECT().Peek().Return(trans)
		})

		Context("mshr hit", func() {
			var (
				mshrEntry *cache.MSHREntry
			)

			BeforeEach(func() {
				mshrEntry = &cache.MSHREntry{}
				mshr.EXPECT().
					Query(ca.PID(1), uint64(0x100)).
					Return(mshrEntry)
			})

			It("should add to MSHR", func() {
				dirBuf.EXPECT().Pop()

				ret := ds.Tick(10)

				Expect(ret).To(BeTrue())
				Expect(mshrEntry.Requests).To(HaveLen(1))
			})
		})

		Context("hit", func() {
			var (
				block *cache.Block
			)

			BeforeEach(func() {
				mshr.EXPECT().
					Query(ca.PID(1), uint64(0x100)).
					Return(nil)

				block = &cache.Block{
					Tag: 0x100,
				}
				directory.EXPECT().
					Lookup(ca.PID(1), uint64(0x100)).
					Return(block)
			})

			It("should stall is bank is busy", func() {
				bankBuf.EXPECT().CanPush().Return(false)

				ret := ds.Tick(10)

				Expect(ret).To(BeFalse())
			})

			It("should stall if block is locked", func() {
				block.IsLocked = true

				ret := ds.Tick(10)

				Expect(ret).To(BeFalse())
			})

			It("should pass transaction to bank", func() {
				bankBuf.EXPECT().CanPush().Return(true)
				bankBuf.EXPECT().Push(gomock.Any()).
					Do(func(trans *transaction) {
						Expect(trans.read).To(BeIdenticalTo(read))
						Expect(trans.block).To(BeIdenticalTo(block))
					})
				dirBuf.EXPECT().Pop()
				directory.EXPECT().Visit(block)

				ret := ds.Tick(10)

				Expect(ret).To(BeTrue())
				Expect(block.ReadCount).To(Equal(1))
				Expect(trans.action).To(Equal(bankReadHit))
			})
		})

		Context("miss, mshr miss, mshr full", func() {
			It("should stall", func() {
				directory.EXPECT().
					Lookup(ca.PID(1), uint64(0x100)).
					Return(nil)
				mshr.EXPECT().Query(ca.PID(1), uint64(0x100)).Return(nil)
				mshr.EXPECT().IsFull().Return(true)

				ret := ds.Tick(10)

				Expect(ret).To(BeFalse())
			})
		})

		Context("miss, mshr miss, no need to evict", func() {
			var (
				block *cache.Block
			)

			BeforeEach(func() {
				block = &cache.Block{
					PID:     2,
					Tag:     0x200,
					IsValid: true,
					IsDirty: false,
				}

				directory.EXPECT().
					Lookup(ca.PID(1), uint64(0x100)).
					Return(nil)
				directory.EXPECT().FindVictim(uint64(0x100)).Return(block)
				mshr.EXPECT().Query(ca.PID(1), uint64(0x100)).Return(nil)
				mshr.EXPECT().IsFull().Return(false)
			})

			It("should stall if WriteBuffer buffer if full", func() {
				bankBuf.EXPECT().CanPush().Return(false)

				ret := ds.Tick(10)

				Expect(ret).To(BeFalse())
			})

			It("should create mshr entry and read from bottom", func() {
				mshrEntry := &cache.MSHREntry{}
				bankBuf.EXPECT().CanPush().Return(true)
				bankBuf.EXPECT().Push(gomock.Any()).
					Do(func(transaction *transaction) {
						Expect(transaction.action).To(Equal(writeBufferFetch))
						Expect(trans.fetchPID).To(Equal(ca.PID(1)))
						Expect(transaction.fetchAddress).
							To(Equal(uint64(0x100)))
					})
				mshr.EXPECT().Add(ca.PID(1), uint64(0x100)).Return(mshrEntry)
				dirBuf.EXPECT().Pop()
				directory.EXPECT().Visit(block)

				ret := ds.Tick(10)

				Expect(ret).To(BeTrue())
				Expect(block.Tag).To(Equal(uint64(0x100)))
				Expect(block.IsValid).To(BeTrue())
				Expect(block.IsLocked).To(BeTrue())
				Expect(block.PID).To(Equal(ca.PID(1)))
				Expect(trans.block).To(BeIdenticalTo(block))
				Expect(mshrEntry.Requests).To(ContainElement(trans))
				Expect(mshrEntry.Block).To(BeIdenticalTo(block))
			})
		})

		Context("miss, mshr miss, need eviction", func() {
			var (
				block *cache.Block
			)

			BeforeEach(func() {
				block = &cache.Block{
					PID:          2,
					Tag:          0x200,
					CacheAddress: 0x300,
					IsValid:      true,
					IsDirty:      true,
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

				directory.EXPECT().
					Lookup(ca.PID(1), uint64(0x100)).
					Return(nil)
				directory.EXPECT().FindVictim(uint64(0x100)).Return(block)
				mshr.EXPECT().Query(ca.PID(1), uint64(0x100)).Return(nil)
				mshr.EXPECT().IsFull().Return(false)
			})

			It("should stall if bank buffer is full", func() {
				bankBuf.EXPECT().CanPush().Return(false)

				ret := ds.Tick(10)

				Expect(ret).To(BeFalse())
			})

			It("should stall if victim is locked", func() {
				block.IsLocked = true

				ret := ds.Tick(10)

				Expect(ret).To(BeFalse())
			})

			It("should do evict", func() {
				bankBuf.EXPECT().CanPush().Return(true)
				bankBuf.EXPECT().
					Push(gomock.Any()).
					Do(func(trans *transaction) {
						Expect(trans.victim.Tag).To(Equal(uint64(0x200)))
						Expect(trans.victim.CacheAddress).
							To(Equal(uint64(0x300)))
					})
				mshrEntry := &cache.MSHREntry{}
				mshr.EXPECT().Add(ca.PID(1), uint64(0x100)).Return(mshrEntry)
				dirBuf.EXPECT().Pop()

				ret := ds.Tick(10)

				Expect(ret).To(BeTrue())
				Expect(block.Tag).To(Equal(uint64(0x100)))
				Expect(block.IsLocked).To(BeTrue())
				Expect(block.IsValid).To(BeTrue())
				Expect(block.IsDirty).To(BeFalse())
				Expect(trans.action).To(Equal(bankEvictAndFetch))
				Expect(trans.block).To(BeIdenticalTo(block))
				Expect(trans.victim.Tag).To(Equal(uint64(0x200)))
				Expect(trans.victim.CacheAddress).To(Equal(uint64(0x300)))
				Expect(trans.victim.DirtyMask).To(Equal([]bool{
					true, true, true, true, false, false, false, false,
					true, true, true, true, false, false, false, false,
					true, true, true, true, false, false, false, false,
					true, true, true, true, false, false, false, false,
					true, true, true, true, false, false, false, false,
					true, true, true, true, false, false, false, false,
					true, true, true, true, false, false, false, false,
					true, true, true, true, false, false, false, false,
				}))
				Expect(trans.evictingPID).To(Equal(ca.PID(2)))
				Expect(trans.evictingAddr).To(Equal(uint64(0x200)))
				Expect(trans.evictingDirtyMask).To(Equal([]bool{
					true, true, true, true, false, false, false, false,
					true, true, true, true, false, false, false, false,
					true, true, true, true, false, false, false, false,
					true, true, true, true, false, false, false, false,
					true, true, true, true, false, false, false, false,
					true, true, true, true, false, false, false, false,
					true, true, true, true, false, false, false, false,
					true, true, true, true, false, false, false, false,
				}))
				Expect(trans.fetchPID).To(Equal(ca.PID(1)))
				Expect(trans.fetchAddress).To(Equal(uint64(0x100)))
				Expect(mshrEntry.Block).To(BeIdenticalTo(block))
				Expect(mshrEntry.Requests).To(ContainElement(trans))
			})
		})
	})

	Context("write", func() {
		var (
			write *mem.WriteReq
			trans *transaction
		)

		BeforeEach(func() {
			write = mem.WriteReqBuilder{}.
				WithSendTime(10).
				WithAddress(0x100).
				WithPID(1).
				Build()
			write.PID = 1
			trans = &transaction{
				write: write,
			}
			dirBuf.EXPECT().Peek().Return(trans)
		})

		Context("mshr hit", func() {
			var (
				mshrEntry *cache.MSHREntry
			)

			BeforeEach(func() {
				mshrEntry = &cache.MSHREntry{}
				mshr.EXPECT().
					Query(ca.PID(1), uint64(0x100)).
					Return(mshrEntry)
			})

			It("should add to MSHR", func() {
				dirBuf.EXPECT().Pop()

				ret := ds.Tick(10)

				Expect(ret).To(BeTrue())
				Expect(mshrEntry.Requests).To(HaveLen(1))
			})
		})

		Context("hit", func() {
			var (
				block *cache.Block
			)

			BeforeEach(func() {
				block = &cache.Block{
					Tag:     0x100,
					IsValid: true,
				}

				mshr.EXPECT().
					Query(ca.PID(1), uint64(0x100)).
					Return(nil)

				directory.EXPECT().
					Lookup(ca.PID(1), uint64(0x100)).
					Return(block)
			})

			It("should stall is bank is busy", func() {
				bankBuf.EXPECT().CanPush().Return(false)

				ret := ds.Tick(10)

				Expect(ret).To(BeFalse())
			})

			It("should stall is block is loked", func() {
				block.IsLocked = true

				ret := ds.Tick(10)

				Expect(ret).To(BeFalse())
			})

			It("should stall if block is being read", func() {
				block.ReadCount = 1

				ret := ds.Tick(10)

				Expect(ret).To(BeFalse())
			})

			It("should send to bank", func() {
				bankBuf.EXPECT().CanPush().Return(true)
				bankBuf.EXPECT().Push(gomock.Any()).
					Do(func(trans *transaction) {
						Expect(trans.block).To(BeIdenticalTo(block))
					})
				dirBuf.EXPECT().Pop()
				directory.EXPECT().Visit(block)

				ret := ds.Tick(10)

				Expect(ret).To(BeTrue())
				Expect(block.IsLocked).To(BeTrue())
				Expect(trans.action).To(Equal(bankWriteHit))
			})
		})

		Context("miss, write full line, no eviction", func() {
			var (
				block *cache.Block
			)

			BeforeEach(func() {
				block = &cache.Block{
					Tag:     0x200,
					IsValid: false,
					IsDirty: false,
				}

				write.Data = []byte{
					1, 2, 3, 4, 5, 6, 7, 8,
					1, 2, 3, 4, 5, 6, 7, 8,
					1, 2, 3, 4, 5, 6, 7, 8,
					1, 2, 3, 4, 5, 6, 7, 8,
					1, 2, 3, 4, 5, 6, 7, 8,
					1, 2, 3, 4, 5, 6, 7, 8,
					1, 2, 3, 4, 5, 6, 7, 8,
					1, 2, 3, 4, 5, 6, 7, 8,
				}
				directory.EXPECT().
					Lookup(ca.PID(1), uint64(0x100)).
					Return(nil)
				directory.EXPECT().FindVictim(uint64(0x100)).Return(block)
				mshr.EXPECT().Query(ca.PID(1), uint64(0x100)).Return(nil)
			})

			It("should stall if victim is locked", func() {
				block.IsLocked = true
				ret := ds.Tick(10)
				Expect(ret).To(BeFalse())
			})

			It("should stall if victim is being read", func() {
				block.ReadCount = 1
				ret := ds.Tick(10)
				Expect(ret).To(BeFalse())
			})

			It("should stall is bank is busy", func() {
				bankBuf.EXPECT().CanPush().Return(false)

				ret := ds.Tick(10)

				Expect(ret).To(BeFalse())
			})

			It("should send to bank", func() {
				bankBuf.EXPECT().CanPush().Return(true)
				bankBuf.EXPECT().Push(gomock.Any()).
					Do(func(trans *transaction) {
						Expect(trans.block).To(BeIdenticalTo(block))
					})
				dirBuf.EXPECT().Pop()
				directory.EXPECT().Visit(block)

				ret := ds.Tick(10)

				Expect(ret).To(BeTrue())
				Expect(block.IsLocked).To(BeTrue())
				Expect(block.Tag).To(Equal(uint64(0x100)))
				Expect(block.IsValid).To(BeTrue())
				Expect(block.PID).To(Equal(ca.PID(1)))
				Expect(trans.action).To(Equal(bankWriteHit))
			})
		})

		Context("miss, write full line, need eviction", func() {
			var (
				block *cache.Block
			)

			BeforeEach(func() {
				block = &cache.Block{
					Tag:          0x200,
					CacheAddress: 0x300,
					IsValid:      true,
					IsDirty:      true,
				}

				directory.EXPECT().
					Lookup(ca.PID(1), uint64(0x100)).
					Return(nil)
				mshr.EXPECT().Query(ca.PID(1), uint64(0x100)).Return(nil)
				directory.EXPECT().FindVictim(uint64(0x100)).Return(block)
				write.Data = make([]byte, 64)
			})

			It("should stall if evictor buffer is full", func() {
				bankBuf.EXPECT().CanPush().Return(false)
				ret := ds.Tick(10)
				Expect(ret).To(BeFalse())
			})

			It("should send to evictor", func() {
				bankBuf.EXPECT().CanPush().Return(true)
				bankBuf.EXPECT().
					Push(gomock.Any()).
					Do(func(trans *transaction) {
						Expect(trans.victim.Tag).To(Equal(uint64(0x200)))
						Expect(trans.victim.CacheAddress).
							To(Equal(uint64(0x300)))
					})
				dirBuf.EXPECT().Pop()

				ret := ds.Tick(10)

				Expect(ret).To(BeTrue())
				Expect(block.Tag).To(Equal(uint64(0x100)))
				Expect(block.IsLocked).To(BeTrue())
				Expect(block.IsValid).To(BeTrue())
				Expect(trans.action).To(Equal(bankEvictAndWrite))
			})
		})

		Context("miss, write partial line, need eviction", func() {
			var (
				block *cache.Block
			)

			BeforeEach(func() {
				block = &cache.Block{
					Tag:          0x200,
					CacheAddress: 0x300,
					IsValid:      true,
					IsDirty:      true,
				}

				write.Data = make([]byte, 4)
				directory.EXPECT().
					Lookup(ca.PID(1), uint64(0x100)).
					Return(nil)
				mshr.EXPECT().Query(ca.PID(1), uint64(0x100)).Return(nil)
			})

			It("should stall if mshr is full", func() {
				mshr.EXPECT().IsFull().Return(true)
				ret := ds.Tick(10)
				Expect(ret).To(BeFalse())
			})

			It("should stall if victim block is locked", func() {
				mshr.EXPECT().IsFull().Return(false)
				directory.EXPECT().FindVictim(uint64(0x100)).Return(block)
				block.IsLocked = true
				ret := ds.Tick(10)
				Expect(ret).To(BeFalse())
			})

			It("should stall if evictor buffer is full", func() {
				mshr.EXPECT().IsFull().Return(false)
				directory.EXPECT().FindVictim(uint64(0x100)).Return(block)
				bankBuf.EXPECT().CanPush().Return(false)
				ret := ds.Tick(10)
				Expect(ret).To(BeFalse())
			})

			It("should send to write buffer and create mshr entry", func() {
				mshrEntry := &cache.MSHREntry{}
				mshr.EXPECT().IsFull().Return(false)
				directory.EXPECT().FindVictim(uint64(0x100)).Return(block)
				bankBuf.EXPECT().CanPush().Return(true)
				bankBuf.EXPECT().
					Push(gomock.Any()).
					Do(func(trans *transaction) {
						Expect(trans.victim.Tag).To(Equal(uint64(0x200)))
						Expect(trans.victim.CacheAddress).
							To(Equal(uint64(0x300)))
					})
				mshr.EXPECT().Add(ca.PID(1), uint64(0x100)).Return(mshrEntry)
				dirBuf.EXPECT().Pop()

				ret := ds.Tick(10)

				Expect(ret).To(BeTrue())
				Expect(block.PID).To(Equal(ca.PID(1)))
				Expect(block.Tag).To(Equal(uint64(0x100)))
				Expect(block.IsLocked).To(BeTrue())
				Expect(block.IsValid).To(BeTrue())
				Expect(block.IsDirty).To(BeFalse())
				Expect(trans.action).To(Equal(bankEvictAndFetch))
			})
		})
	})
})
