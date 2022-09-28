package org

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem/dram/internal/signal"
	"gitlab.com/akita/util/tracing"
)

type Bank interface {
	tracing.NamedHookable

	GetReadyCommand(
		now akita.VTimeInSec,
		cmd *signal.Command,
	) *signal.Command
	StartCommand(now akita.VTimeInSec, cmd *signal.Command)
	UpdateTiming(cmdKind signal.CommandKind, cycleNeeded int)
	Tick(now akita.VTimeInSec) bool
}
