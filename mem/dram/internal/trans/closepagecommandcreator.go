package trans

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem/dram/internal/addressmapping"
	"gitlab.com/akita/mem/dram/internal/signal"
)

// ClosePageCommandCreator always creates precharge commands as precharge
// commands will be the last command in a row.
type ClosePageCommandCreator struct {
	AddrMapper addressmapping.Mapper
}

func (c *ClosePageCommandCreator) Create(
	subTrans *signal.SubTransaction,
) *signal.Command {
	cmd := &signal.Command{
		ID: akita.GetIDGenerator().Generate(),
	}

	if subTrans.IsRead() {
		cmd.Kind = signal.CmdKindReadPrecharge
	} else {
		cmd.Kind = signal.CmdKindWritePrecharge
	}

	cmd.Location = c.AddrMapper.Map(subTrans.Address)
	cmd.SubTrans = subTrans

	return cmd
}
