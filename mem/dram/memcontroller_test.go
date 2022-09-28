package dram

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem"
	"gitlab.com/akita/mem/dram/internal/signal"
)

var _ = Describe("MemController", func() {
	var (
		mockCtrl *gomock.Controller

		topPort             *MockPort
		addrConverter       *MockAddressConverter
		subTransSplitter    *MockSubTransSplitter
		subTransactionQueue *MockSubTransactionQueue
		cmdQueue            *MockCommandQueue
		channel             *MockChannel
		storage             *mem.Storage

		memCtrl *MemController
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())

		topPort = NewMockPort(mockCtrl)
		subTransactionQueue = NewMockSubTransactionQueue(mockCtrl)
		subTransSplitter = NewMockSubTransSplitter(mockCtrl)
		addrConverter = NewMockAddressConverter(mockCtrl)
		cmdQueue = NewMockCommandQueue(mockCtrl)
		channel = NewMockChannel(mockCtrl)
		storage = mem.NewStorage(4 * mem.GB)

		memCtrl = MakeBuilder().Build("memCtrl")
		memCtrl.TopPort = topPort
		memCtrl.subTransactionQueue = subTransactionQueue
		memCtrl.subTransSplitter = subTransSplitter
		memCtrl.addrConverter = addrConverter
		memCtrl.cmdQueue = cmdQueue
		memCtrl.channel = channel
		memCtrl.storage = storage
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("parse top", func() {
		It("should do nothing if no message", func() {
			topPort.EXPECT().Peek().Return(nil)

			madeProgress := memCtrl.parseTop(10)

			Expect(madeProgress).To(BeFalse())
		})

		It("should stall if substransaction queue is full", func() {
			read := mem.ReadReqBuilder{}.
				WithAddress(0x1000).
				Build()

			topPort.EXPECT().Peek().Return(read)
			addrConverter.EXPECT().ConvertExternalToInternal(uint64(0x1000))
			subTransSplitter.EXPECT().
				Split(gomock.Any()).
				Do(func(t *signal.Transaction) {
					Expect(t.Read).To(BeIdenticalTo(read))
					t.SubTransactions = make([]*signal.SubTransaction, 3)
				})
			subTransactionQueue.EXPECT().CanPush(3).Return(false)

			madeProgress := memCtrl.parseTop(10)

			Expect(madeProgress).To(BeFalse())
		})

		It("should push sub-transactions to subtrans queue", func() {
			read := mem.ReadReqBuilder{}.
				WithAddress(0x1000).
				Build()

			topPort.EXPECT().Peek().Return(read)
			topPort.EXPECT().Retrieve(gomock.Any()).Return(read)
			addrConverter.EXPECT().ConvertExternalToInternal(uint64(0x1000))
			subTransSplitter.EXPECT().
				Split(gomock.Any()).
				Do(func(t *signal.Transaction) {
					Expect(t.Read).To(BeIdenticalTo(read))
					for i := 0; i < 3; i++ {
						st := &signal.SubTransaction{}
						t.SubTransactions = append(t.SubTransactions, st)
					}
				})
			subTransactionQueue.EXPECT().CanPush(3).Return(true)
			subTransactionQueue.EXPECT().Push(gomock.Any())

			madeProgress := memCtrl.parseTop(10)

			Expect(madeProgress).To(BeTrue())
			Expect(memCtrl.inflightTransactions).To(HaveLen(1))
		})

	})

	Context("issue", func() {
		It("should not issue if nothing is ready", func() {
			cmdQueue.EXPECT().
				GetCommandToIssue(akita.VTimeInSec(10)).
				Return(nil)

			madeProgress := memCtrl.issue(10)

			Expect(madeProgress).To(BeFalse())
		})

		It("should issue", func() {
			cmd := &signal.Command{}
			cmdQueue.EXPECT().
				GetCommandToIssue(akita.VTimeInSec(10)).
				Return(cmd)
			channel.EXPECT().StartCommand(akita.VTimeInSec(10), cmd)
			channel.EXPECT().UpdateTiming(akita.VTimeInSec(10), cmd)

			madeProgress := memCtrl.issue(10)

			Expect(madeProgress).To(BeTrue())
		})
	})

	Context("respond", func() {
		It("should do nothing if there is no transaction", func() {
			madeProgress := memCtrl.respond(10)

			Expect(madeProgress).To(BeFalse())
		})

		It("should do nothing if there is no completed transaction",
			func() {
				trans := &signal.Transaction{}
				subTransaction := &signal.SubTransaction{
					Transaction: trans,
					Completed:   false,
				}
				trans.SubTransactions = append(trans.SubTransactions,
					subTransaction)
				memCtrl.inflightTransactions = append(
					memCtrl.inflightTransactions, trans)

				madeProgress := memCtrl.respond(10)

				Expect(madeProgress).To(BeFalse())
			})

		It("should send write done response", func() {
			write := mem.WriteReqBuilder{}.
				WithAddress(0x40).
				WithData([]byte{1, 2, 3, 4}).
				Build()
			trans := &signal.Transaction{
				InternalAddress: 0x40,
				Write:           write,
			}
			subTransaction := &signal.SubTransaction{
				Transaction: trans,
				Completed:   true,
			}
			trans.SubTransactions = append(trans.SubTransactions,
				subTransaction)
			memCtrl.inflightTransactions = append(memCtrl.inflightTransactions,
				trans)

			topPort.EXPECT().Send(gomock.Any()).Return(nil)

			madeProgress := memCtrl.respond(10)

			Expect(madeProgress).To(BeTrue())
			data, _ := storage.Read(0x40, 4)
			Expect(data).To(Equal([]byte{1, 2, 3, 4}))
			Expect(memCtrl.inflightTransactions).NotTo(ContainElement(trans))
		})

		It("should send data ready response", func() {
			storage.Write(0x40, []byte{1, 2, 3, 4})
			read := mem.ReadReqBuilder{}.
				WithAddress(0x40).
				WithByteSize(4).
				Build()
			trans := &signal.Transaction{
				InternalAddress: 0x40,
				Read:            read,
			}
			subTransaction := &signal.SubTransaction{
				Transaction: trans,
				Completed:   true,
			}
			trans.SubTransactions = append(trans.SubTransactions,
				subTransaction)
			memCtrl.inflightTransactions = append(memCtrl.inflightTransactions,
				trans)

			topPort.EXPECT().Send(gomock.Any()).Do(func(dr *mem.DataReadyRsp) {
				Expect(dr.Data).To(Equal([]byte{1, 2, 3, 4}))
			}).Return(nil)

			madeProgress := memCtrl.respond(10)

			Expect(madeProgress).To(BeTrue())
			Expect(memCtrl.inflightTransactions).NotTo(ContainElement(trans))
		})
	})
})
