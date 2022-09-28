package mem

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/akita/akita"
)

var _ = Describe("Write Buffer", func() {
	var (
		mockCtrl *gomock.Controller
		port     *MockPort
		wb       *writeBufferImpl
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		port = NewMockPort(mockCtrl)
		wb = NewWriteBuffer(2, port).(*writeBufferImpl)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should enqueue", func() {
		write := WriteReqBuilder{}.
			WithSendTime(10).
			WithAddress(0x100).
			Build()

		wb.Enqueue(write)

		Expect(wb.buf).To(ContainElement(write))
	})

	It("should write combine", func() {
		write1 := WriteReqBuilder{}.
			WithSendTime(10).
			WithAddress(0x104).
			WithData([]byte{1, 2, 3, 4}).
			Build()
		wb.Enqueue(write1)

		write2 := WriteReqBuilder{}.
			WithSendTime(11).
			WithAddress(0x120).
			WithData([]byte{1, 2, 3, 4}).
			Build()
		wb.Enqueue(write2)

		Expect(wb.buf).To(HaveLen(1))
		combinedWrite := wb.buf[0]
		Expect(combinedWrite.Data).To(Equal([]byte{
			0, 0, 0, 0, 1, 2, 3, 4,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			1, 2, 3, 4, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
		}))
		Expect(combinedWrite.DirtyMask).To(Equal([]bool{
			false, false, false, false, true, true, true, true,
			false, false, false, false, false, false, false, false,
			false, false, false, false, false, false, false, false,
			false, false, false, false, false, false, false, false,
			true, true, true, true, false, false, false, false,
			false, false, false, false, false, false, false, false,
			false, false, false, false, false, false, false, false,
			false, false, false, false, false, false, false, false,
		}))
	})

	It("should write combine", func() {
		write1 := WriteReqBuilder{}.
			WithSendTime(10).
			WithAddress(0x104).
			WithData([]byte{1, 2, 3, 4, 0, 0, 0, 0, 5, 6, 7, 8, 9, 9, 9, 9}).
			WithDirtyMask([]bool{
				true, true, true, true,
				false, false, false, false,
				true, true, true, true,
				true, true, true, true,
			}).
			Build()
		wb.Enqueue(write1)

		write2 := WriteReqBuilder{}.
			WithSendTime(11).
			WithAddress(0x108).
			WithData([]byte{1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4}).
			WithDirtyMask([]bool{
				true, true, true, true,
				true, true, true, true,
				false, false, false, false,
			}).
			Build()
		wb.Enqueue(write2)

		Expect(wb.buf).To(HaveLen(1))
		combinedWrite := wb.buf[0]
		Expect(combinedWrite.Data).To(Equal([]byte{
			0, 0, 0, 0, 1, 2, 3, 4,
			1, 2, 3, 4, 1, 2, 3, 4,
			9, 9, 9, 9, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
		}))
		Expect(combinedWrite.DirtyMask).To(Equal([]bool{
			false, false, false, false, true, true, true, true,
			true, true, true, true, true, true, true, true,
			true, true, true, true, false, false, false, false,
			false, false, false, false, false, false, false, false,
			false, false, false, false, false, false, false, false,
			false, false, false, false, false, false, false, false,
			false, false, false, false, false, false, false, false,
			false, false, false, false, false, false, false, false,
		}))
	})

	It("should not combine write from different PID", func() {
		write1 := WriteReqBuilder{}.
			WithSendTime(10).
			WithAddress(0x104).
			WithPID(1).
			WithData([]byte{1, 2, 3, 4}).
			Build()
		wb.Enqueue(write1)

		write2 := WriteReqBuilder{}.
			WithSendTime(11).
			WithAddress(0x120).
			WithPID(2).
			WithData([]byte{1, 2, 3, 4}).
			Build()
		wb.Enqueue(write2)

		Expect(wb.buf).To(HaveLen(2))
	})

	It("should panic when trying to enqueue over capacity", func() {
		wb.buf = make([]*WriteReq, 2)
		write := WriteReqBuilder{}.Build()

		Expect(func() { wb.Enqueue(write) }).To(Panic())
	})

	It("should check if it can enqueue", func() {
		Expect(wb.CanEnqueue()).To(BeTrue())

		write1 := WriteReqBuilder{}.WithPID(1).Build()
		wb.Enqueue(write1)
		Expect(wb.CanEnqueue()).To(BeTrue())

		write2 := WriteReqBuilder{}.WithPID(2).Build()
		wb.Enqueue(write2)
		Expect(wb.CanEnqueue()).To(BeFalse())
	})

	It("should query", func() {
		write := WriteReqBuilder{}.
			WithSendTime(10).
			WithAddress(0x100).
			WithData(
				[]byte{
					1, 2, 3, 4, 5, 6, 7, 8,
					1, 2, 3, 4, 5, 6, 7, 8,
					1, 2, 3, 4, 5, 6, 7, 8,
					1, 2, 3, 4, 5, 6, 7, 8,
					1, 2, 3, 4, 5, 6, 7, 8,
					1, 2, 3, 4, 5, 6, 7, 8,
					1, 2, 3, 4, 5, 6, 7, 8,
				}).
			WithPID(1).
			Build()
		wb.Enqueue(write)

		read := ReadReqBuilder{}.
			WithSendTime(11).
			WithAddress(0x100).
			WithByteSize(4).
			WithPID(1).
			Build()

		ret := wb.Query(read)

		Expect(ret).To(BeIdenticalTo(write))
	})

	It("should not return write if the read is from another PID", func() {
		write := WriteReqBuilder{}.
			WithSendTime(10).
			WithAddress(0x100).
			WithPID(1).
			WithData(
				[]byte{
					1, 2, 3, 4, 5, 6, 7, 8,
					1, 2, 3, 4, 5, 6, 7, 8,
					1, 2, 3, 4, 5, 6, 7, 8,
					1, 2, 3, 4, 5, 6, 7, 8,
					1, 2, 3, 4, 5, 6, 7, 8,
					1, 2, 3, 4, 5, 6, 7, 8,
					1, 2, 3, 4, 5, 6, 7, 8,
				}).
			Build()
		wb.Enqueue(write)

		read := ReadReqBuilder{}.
			WithSendTime(11).
			WithAddress(0x100).
			WithByteSize(4).
			WithPID(2).
			Build()

		ret := wb.Query(read)

		Expect(ret).To(BeNil())
	})

	It("should do nothing if there is no write to send", func() {
		ret := wb.Tick(10)
		Expect(ret).To(BeFalse())
	})

	It("should stall if send failed", func() {
		write := WriteReqBuilder{}.
			WithSendTime(10).
			WithAddress(0x100).
			Build()
		wb.Enqueue(write)
		port.EXPECT().Send(write).Return(&akita.SendError{})

		ret := wb.Tick(10)

		Expect(ret).To(BeFalse())
		Expect(wb.buf).To(ContainElement(write))
	})

	It("should send requst to bottom", func() {
		write := WriteReqBuilder{}.
			WithSendTime(10).
			WithAddress(0x100).
			Build()
		wb.Enqueue(write)
		port.EXPECT().Send(write).Return(nil)

		ret := wb.Tick(10)

		Expect(ret).To(BeTrue())
		Expect(wb.buf).To(HaveLen(0))
	})
})
