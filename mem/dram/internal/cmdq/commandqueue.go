// Package cmdq provides command queue implementations
package cmdq

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem/dram/internal/signal"
)

type CommandQueue interface {
	GetCommandToIssue(
		now akita.VTimeInSec,
	) *signal.Command
	CanAccept(command *signal.Command) bool
	Accept(command *signal.Command)
}
