package dram

import (
	"testing"

	"gitlab.com/akita/mem"

	"github.com/golang/mock/gomock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/akita/akita"
)

//go:generate mockgen -destination "mock_akita_test.go" -package $GOPACKAGE -write_package_comment=false gitlab.com/akita/akita Port
//go:generate mockgen -destination "mock_trans_test.go" -package $GOPACKAGE -write_package_comment=false gitlab.com/akita/mem/dram/internal/trans SubTransactionQueue,SubTransSplitter
//go:generate mockgen -destination "mock_addressmapping_test.go" -package $GOPACKAGE -write_package_comment=false gitlab.com/akita/mem/dram/internal/addressmapping AddressConverter,Mapper
//go:generate mockgen -destination "mock_cmdq_test.go" -package $GOPACKAGE -write_package_comment=false gitlab.com/akita/mem/dram/internal/cmdq CommandQueue
//go:generate mockgen -destination "mock_org_test.go" -package $GOPACKAGE -write_package_comment=false gitlab.com/akita/mem/dram/internal/org Channel

func TestDram(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dram Suite")
}

var _ = Describe("DRAM Integration", func() {
	var (
		mockCtrl *gomock.Controller
		engine   akita.Engine
		srcPort  *MockPort
		memCtrl  *MemController
		conn     *akita.DirectConnection
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		engine = akita.NewSerialEngine()
		memCtrl = MakeBuilder().
			WithEngine(engine).
			Build("memCtrl")
		srcPort = NewMockPort(mockCtrl)

		conn = akita.NewDirectConnection("conn", engine, 1*akita.GHz)
		srcPort.EXPECT().SetConnection(conn)
		conn.PlugIn(memCtrl.TopPort, 1)
		conn.PlugIn(srcPort, 1)
	})

	It("should read and write", func() {
		write := mem.WriteReqBuilder{}.
			WithAddress(0x40).
			WithData([]byte{1, 2, 3, 4}).
			WithSrc(srcPort).
			WithDst(memCtrl.TopPort).
			WithSendTime(0).
			Build()

		read := mem.ReadReqBuilder{}.
			WithAddress(0x40).
			WithByteSize(4).
			WithSrc(srcPort).
			WithDst(memCtrl.TopPort).
			WithSendTime(0).
			Build()

		memCtrl.TopPort.Recv(write)
		memCtrl.TopPort.Recv(read)

		ret1 := srcPort.EXPECT().
			Recv(gomock.Any()).
			Do(func(wd *mem.WriteDoneRsp) {
				Expect(wd.RespondTo).To(Equal(write.ID))
			})
		srcPort.EXPECT().
			Recv(gomock.Any()).
			Do(func(dr *mem.DataReadyRsp) {
				Expect(dr.RespondTo).To(Equal(read.ID))
				Expect(dr.Data).To(Equal([]byte{1, 2, 3, 4}))
			}).After(ret1)

		engine.Run()
	})
})
