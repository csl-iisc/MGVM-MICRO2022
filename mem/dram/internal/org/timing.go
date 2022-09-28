package org

import "gitlab.com/akita/mem/dram/internal/signal"

type TimeTable [][]TimeTableEntry

type TimeTableEntry struct {
	NextCmdKind       signal.CommandKind
	MinCycleInBetween int
}

func (t TimeTable) getTimeAfter(cmdKind signal.CommandKind) []TimeTableEntry {
	return t[cmdKind]
}

func MakeTimeTable() TimeTable {
	return make([][]TimeTableEntry, signal.NumCmdKind)
}

type Timing struct {
	SameBank              TimeTable
	OtherBanksInBankGroup TimeTable
	SameRank              TimeTable
	OtherRanks            TimeTable
}
