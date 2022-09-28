package cmdq

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem/dram/internal/addressmapping"
	"gitlab.com/akita/mem/dram/internal/signal"
)

var _ = Describe("CommandQueueImpl", func() {
	var (
		mockCtrl *gomock.Controller
		channel  *MockChannel
		q        CommandQueueImpl
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		channel = NewMockChannel(mockCtrl)
		q = CommandQueueImpl{
			Queues:           make([]Queue, 8),
			CapacityPerQueue: 8,
			nextQueueIndex:   0,
			Channel:          channel,
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should get the next command to issue", func() {
		cmd1 := &signal.Command{
			ID:   "1",
			Kind: signal.CmdKindRead,
			Location: addressmapping.Location{
				Rank: 0,
				Bank: 0,
			},
		}
		q.Queues[0] = append(q.Queues[0], cmd1)

		cmd2 := &signal.Command{
			ID:   "2",
			Kind: signal.CmdKindRead,
			Location: addressmapping.Location{
				Rank: 0,
				Bank: 0,
			},
		}
		q.Queues[0] = append(q.Queues[0], cmd2)

		cmd3 := &signal.Command{
			ID:   "3",
			Kind: signal.CmdKindRead,
			Location: addressmapping.Location{
				Rank: 0,
				Bank: 1,
			},
		}
		q.Queues[1] = append(q.Queues[1], cmd3)

		channel.EXPECT().
			GetReadyCommand(akita.VTimeInSec(10), cmd1).
			Return(nil)
		channel.EXPECT().
			GetReadyCommand(akita.VTimeInSec(10), cmd2).
			Return(cmd2)

		readyCmd := q.GetCommandToIssue(10)

		Expect(readyCmd).To(BeIdenticalTo(cmd2))
		Expect(q.Queues[0]).NotTo(ContainElement(cmd2))
	})

	It("should accept new commands", func() {
		cmd := &signal.Command{}

		Expect(q.CanAccept(cmd)).To(BeTrue())

		q.Accept(cmd)

		Expect(q.Queues[0]).To(ContainElement(cmd))
	})
})
