package org

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem/dram/internal/signal"
)

// Banks is indexed by rank, bank-group, bank.
type Banks [][][]Bank

func (b Banks) GetSize() (rank, bankGroup, bank uint64) {
	return uint64(len(b)), uint64(len(b[0])), uint64(len(b[0][0]))
}

func (b Banks) GetBank(rank, bankGroup, bank uint64) Bank {
	return b[rank][bankGroup][bank]
}

func MakeBanks(numRank, numBankGroup, numBank uint64) Banks {
	b := make(Banks, numRank)

	for i := uint64(0); i < numRank; i++ {
		b[i] = make([][]Bank, numBankGroup)

		for j := uint64(0); j < numBankGroup; j++ {
			b[i][j] = make([]Bank, numBank)

			for k := uint64(0); k < numBank; k++ {
				b[i][j][k] = NewBankImpl("")
			}
		}
	}

	return b
}

type Channel interface {
	GetReadyCommand(
		now akita.VTimeInSec,
		cmd *signal.Command,
	) *signal.Command

	StartCommand(
		now akita.VTimeInSec,
		cmd *signal.Command,
	)

	UpdateTiming(
		now akita.VTimeInSec,
		cmd *signal.Command,
	)

	Tick(now akita.VTimeInSec) (madeProgress bool)
}

type ChannelImpl struct {
	Banks  Banks
	Timing Timing
}

func (cs *ChannelImpl) Tick(now akita.VTimeInSec) (madeProgress bool) {
	for i := 0; i < len(cs.Banks); i++ {
		for j := 0; j < len(cs.Banks[0]); j++ {
			for k := 0; k < len(cs.Banks[0][0]); k++ {
				madeProgress = cs.Banks[i][j][k].Tick(now) || madeProgress
			}
		}
	}

	return madeProgress
}

func (cs *ChannelImpl) GetReadyCommand(
	now akita.VTimeInSec,
	cmd *signal.Command,
) *signal.Command {
	readyCmd := cs.Banks.
		GetBank(cmd.Rank, cmd.BankGroup, cmd.Bank).
		GetReadyCommand(now, cmd)

	return readyCmd
}

func (cs *ChannelImpl) StartCommand(now akita.VTimeInSec, cmd *signal.Command) {
	cs.Banks.
		GetBank(cmd.Rank, cmd.BankGroup, cmd.Bank).
		StartCommand(now, cmd)
}

func (cs *ChannelImpl) UpdateTiming(now akita.VTimeInSec, cmd *signal.Command) {
	switch cmd.Kind {
	case signal.CmdKindActivate:
		fallthrough
	case signal.CmdKindRead, signal.CmdKindReadPrecharge,
		signal.CmdKindWrite, signal.CmdKindWritePrecharge,
		signal.CmdKindPrecharge, signal.CmdKindRefreshBank:
		cs.updateAllBankTiming(now, cmd)
	}
}

func (cs *ChannelImpl) updateAllBankTiming(
	now akita.VTimeInSec,
	cmd *signal.Command,
) {
	rank, bankGroup, bank := cs.Banks.GetSize()
	for i := uint64(0); i < rank; i++ {
		for j := uint64(0); j < bankGroup; j++ {
			for k := uint64(0); k < bank; k++ {
				cs.updateBankTiming(now, cmd, i, j, k)
			}
		}
	}
}

func (cs *ChannelImpl) updateBankTiming(
	now akita.VTimeInSec,
	cmd *signal.Command,
	rank, bankGroup, bank uint64,
) {
	timingTable := cs.Timing.OtherRanks
	if cmd.Rank == rank {
		timingTable = cs.Timing.SameRank

		if cmd.BankGroup == bankGroup {
			timingTable = cs.Timing.OtherBanksInBankGroup

			if cmd.Bank == bank {
				timingTable = cs.Timing.SameBank
			}
		}
	}

	for _, entry := range timingTable[cmd.Kind] {
		cs.Banks.GetBank(rank, bankGroup, bank).
			UpdateTiming(entry.NextCmdKind, entry.MinCycleInBetween)
	}
}
