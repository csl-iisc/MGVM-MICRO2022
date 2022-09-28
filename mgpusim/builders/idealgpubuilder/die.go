package idealgpubuilder

import (
	// "fmt"

	"gitlab.com/akita/akita"
	// "gitlab.com/akita/mem"
	"gitlab.com/akita/mem/cache"
	"gitlab.com/akita/mem/cache/writeback"

	// "gitlab.com/akita/mem/dram"
	"gitlab.com/akita/mem/idealmemcontroller"
	"gitlab.com/akita/mem/vm/addresstranslator"
	"gitlab.com/akita/mem/vm/mmu"

	// "gitlab.com/akita/mem/vm/mmu"
	"gitlab.com/akita/mem/vm/idealtlb"
	"gitlab.com/akita/mem/vm/tlb"
	// "gitlab.com/akita/mgpusim"
	// "gitlab.com/akita/mgpusim/pagemigrationcontroller"
	"gitlab.com/akita/mgpusim/rdma"
	"gitlab.com/akita/mgpusim/remotetranslation"
	"gitlab.com/akita/mgpusim/timing/caches/l1v"
	"gitlab.com/akita/mgpusim/timing/caches/rob"

	// "gitlab.com/akita/mgpusim/timing/cp"
	"gitlab.com/akita/mgpusim/timing/cu"
	// "gitlab.com/akita/util/tracing"
)

// Chiplet represents a single chiplet or die of an MCM GPU.
type Chiplet struct {
	Dispatchers       []akita.Component
	CUs               []*cu.ComputeUnit
	L1VCaches         []*l1v.Cache
	L1ICaches         []*l1v.Cache
	L1SCaches         []*l1v.Cache
	L2Caches          []*writeback.Cache
	L2CacheFinder     cache.LowModuleFinder
	L2TLBs            []tlb.L2TLB
	L1VTLBs           []*idealtlb.IdealTLB
	L1STLBs           []*idealtlb.IdealTLB
	L1ITLBs           []*idealtlb.IdealTLB
	L1VROBs           []*rob.ReorderBuffer
	L1IROBs           []*rob.ReorderBuffer
	L1SROBs           []*rob.ReorderBuffer
	ReorderBuffers    []*rob.ReorderBuffer
	L1VAddrTranslator []addresstranslator.AddressTranslator
	L1IAddrTranslator []addresstranslator.AddressTranslator
	L1SAddrTranslator []addresstranslator.AddressTranslator
	DRAMs             []*idealmemcontroller.Comp
	MMU               *mmu.MMUImpl

	chipRdmaEngine        *rdma.Engine
	lowModuleFinderForL1  *cache.InterleavedLowModuleFinder
	remoteTranslationUnit *remotetranslation.RemoteTranslationUnit

	InternalConnection     akita.Connection
	L1TLBToL2TLBConnection *akita.DirectConnection
	L1ToL2Connection       *akita.DirectConnection
	L2ToDramConnection     *akita.DirectConnection

	name      string
	ChipletID uint64
}

// NewChiplet returna a new Chiplet instance with only name and ID set.
func NewChiplet(name string, id uint64) *Chiplet {
	chiplet := new(Chiplet)
	chiplet.name = name
	chiplet.ChipletID = id

	return chiplet
}

// InterChipletPorts returns the list of ports that are connected to other
// chiplets of the GPU
func (c Chiplet) InterChipletPorts() []akita.Port {
	ports := []akita.Port{
		c.chipRdmaEngine.RequestPort,
		c.chipRdmaEngine.ResponsePort,
	}

	return ports
}
