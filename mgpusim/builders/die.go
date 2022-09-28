package builders

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem/cache"
	"gitlab.com/akita/mem/cache/writeback"
	"gitlab.com/akita/mem/idealmemcontroller"
	"gitlab.com/akita/mem/vm/addresstranslator"
	"gitlab.com/akita/mem/vm/mmu"
	"gitlab.com/akita/mem/vm/tlb"
	"gitlab.com/akita/mgpusim/rdma"
	"gitlab.com/akita/mgpusim/remotetranslation"
	"gitlab.com/akita/mgpusim/timing/caches/l1v"
	"gitlab.com/akita/mgpusim/timing/caches/rob"
	"gitlab.com/akita/mgpusim/timing/cu"
)

type Chiplet struct {
	Dispatchers       []akita.Component
	CUs               []*cu.ComputeUnit
	L1VCaches         []*l1v.Cache
	L1ICaches         []*l1v.Cache
	L1SCaches         []*l1v.Cache
	L2Caches          []*writeback.Cache
	L2CacheFinder     cache.LowModuleFinder
	L2TLBs            []tlb.L2TLB
	L1VTLBs           []*tlb.TLB
	L1STLBs           []*tlb.TLB
	L1ITLBs           []*tlb.TLB
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
	pageRdmaEngine        *rdma.Engine
	lowModuleFinderForL1  cache.LowModuleFinder
	remoteTranslationUnit remotetranslation.RemoteTranslationUnit

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
