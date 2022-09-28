package arbitration

import (
	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"gitlab.com/akita/noc"
)

var _ = Describe("XBar", func() {
	var (
		mockCtrl         *gomock.Controller
		buf1, buf1Remote *MockBuffer
		buf2             *MockBuffer
		buf3, buf3Remote *MockBuffer
		buf4, buf4Remote *MockBuffer
		xbar             *xbarArbiter
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		buf1 = NewMockBuffer(mockCtrl)
		buf2 = NewMockBuffer(mockCtrl)
		buf3 = NewMockBuffer(mockCtrl)
		buf4 = NewMockBuffer(mockCtrl)
		buf1Remote = NewMockBuffer(mockCtrl)
		buf3Remote = NewMockBuffer(mockCtrl)
		buf4Remote = NewMockBuffer(mockCtrl)

		xbar = NewXBarArbiter().(*xbarArbiter)
		xbar.AddBuffer(buf1)
		xbar.AddBuffer(buf2)
		xbar.AddBuffer(buf3)
		xbar.AddBuffer(buf4)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should arbitrate", func() {
		flit1 := noc.FlitBuilder{}.
			Build()
		flit1.OutputBuf = buf1Remote
		flit2 := noc.FlitBuilder{}.
			Build()
		flit2.OutputBuf = buf1Remote
		flit3 := noc.FlitBuilder{}.
			Build()
		flit3.OutputBuf = buf3Remote
		flit4 := noc.FlitBuilder{}.
			Build()
		flit4.OutputBuf = buf4Remote
		flit5 := noc.FlitBuilder{}.
			Build()
		flit5.OutputBuf = buf1Remote

		buf1.EXPECT().Peek().Return(flit1)
		buf2.EXPECT().Peek().Return(flit2)
		buf3.EXPECT().Peek().Return(flit3)
		buf4.EXPECT().Peek().Return(flit4)

		bufs := xbar.Arbitrate(10)
		Expect(bufs).To(HaveLen(3))
		Expect(bufs[0]).To(BeIdenticalTo(buf1))
		Expect(bufs[1]).To(BeIdenticalTo(buf3))
		Expect(bufs[2]).To(BeIdenticalTo(buf4))

		buf1.EXPECT().Peek().Return(flit5)
		buf2.EXPECT().Peek().Return(flit2)
		buf3.EXPECT().Peek().Return(nil)
		buf4.EXPECT().Peek().Return(nil)

		bufs = xbar.Arbitrate(10)
		Expect(bufs).To(HaveLen(1))
		Expect(bufs[0]).To(BeIdenticalTo(buf2))
	})
})
