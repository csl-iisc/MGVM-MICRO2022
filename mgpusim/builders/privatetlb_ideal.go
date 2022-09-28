package builders

import (
	"fmt"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem/cache"
	"gitlab.com/akita/mem/vm/tlb"
	"gitlab.com/akita/mgpusim"
	"gitlab.com/akita/noc/networking/chipnetwork"
	"gitlab.com/akita/util/tracing"
)

type PrivateTLBIdealGPUBuilder struct {
	*CommonBuilder

	// specific componenets
}

// Distributed TLB specific function

// MakeDistributedTLBGPUBuilder provides a GPU builder that can builds MCM GPU.
func MakePrivateTLBIdealGPUBuilder() PrivateTLBIdealGPUBuilder {
	// TODO: should this be using new? is the object being allocated on the stack?
	cbp := CommonBuilder{}
	b := PrivateTLBIdealGPUBuilder{CommonBuilder: &cbp}
	b.SetDefaultCommonBuilderParams()
	return b
}

func (b PrivateTLBIdealGPUBuilder) Build(name string, id uint64) *mgpusim.GPU {
	b.createGPU(name, id)

	b.buildCP()

	chipRdmaAddressTable := b.createChipRDMAAddrTable()
	rdmaResponsePorts := make([]akita.Port, 4)
	// remoteAddressTranslationTable := b.createRemoteAddrTransTable()
	// rtuResponsePorts := make([]akita.Port, 4)

	for i := 0; i < b.numChiplet; i++ {
		chipletName := fmt.Sprintf("%s.chiplet_%02d", b.gpuName, i)
		chiplet := NewChiplet(chipletName, uint64(i))

		b.BuildSAs(chiplet)
		b.buildMemBanks(chiplet)
		b.buildMMU(chiplet)
		// TODO: this may have to be overridden
		b.buildL2TLB(chiplet)

		b.configChipRDMAEngine(chiplet, chipRdmaAddressTable, rdmaResponsePorts)
		// b.configRemoteAddressTranslationUnit(chiplet, remoteAddressTranslationTable, rtuResponsePorts)

		b.connectL1ToL2(chiplet)
		b.connectL2ToDRAM(chiplet)
		b.connectL1TLBToL2TLB(chiplet)
		b.connectL2TLBTOMMU(chiplet)
		b.connectMMUToL2(chiplet)

		b.chiplets = append(b.chiplets, chiplet)
	}

	b.buildPageMigrationController()
	b.setupDMA()

	// b.setupMMUs()
	b.connectCP()
	b.setupInterchipNetwork()

	return b.gpu
}

func (b *PrivateTLBIdealGPUBuilder) connectL1TLBToL2TLB(chiplet *Chiplet) {
	tlbConn := akita.NewDirectConnection(chiplet.name+"L1TLB-L2TLB",
		b.engine, b.freq)
	tlbConn.PlugIn(chiplet.L2TLBs[0].GetTopPort(), 64)

	var lowModuleFinder cache.LowModuleFinder

	singeLowModuleFinder := new(cache.SingleLowModuleFinder)
	singeLowModuleFinder.LowModule = chiplet.L2TLBs[0].GetTopPort()

	lowModuleFinder = singeLowModuleFinder

	for _, l1vTLB := range chiplet.L1VTLBs {
		l1vTLB.SetLowModuleFinder(lowModuleFinder)
		tlbConn.PlugIn(l1vTLB.BottomPort, 16)
	}

	for _, l1iTLB := range chiplet.L1ITLBs {
		l1iTLB.SetLowModuleFinder(lowModuleFinder)
		tlbConn.PlugIn(l1iTLB.BottomPort, 16)
	}

	for _, l1sTLB := range chiplet.L1STLBs {
		l1sTLB.SetLowModuleFinder(lowModuleFinder)
		tlbConn.PlugIn(l1sTLB.BottomPort, 16)
	}
}

func (b *PrivateTLBIdealGPUBuilder) setupInterchipNetwork() {
	chipConnector := chipnetwork.NewInterChipletConnector().
		WithEngine(b.engine).
		WithSwitchLatency(360).
		WithFreq(1 * akita.GHz).
		WithFlitByteSize(64).
		WithNumReqPerCycle(12).
		WithNetworkName("ICN")
	chipConnector.CreateNetwork()
	for _, chiplet := range b.chiplets {
		chipConnector.PlugInChip(b.InterChipletPorts(chiplet))
	}
	chipConnector.MakeNetwork()
}

func (b *PrivateTLBIdealGPUBuilder) InterChipletPorts(c *Chiplet) []akita.Port {
	ports := []akita.Port{
		c.chipRdmaEngine.RequestPort,
		c.chipRdmaEngine.ResponsePort,
	}
	return ports
}

func (b *PrivateTLBIdealGPUBuilder) buildL2TLB(chiplet *Chiplet) {
	builder := tlb.MakeIdealLatTLBBuilder().
		WithEngine(b.engine).
		WithFreq(b.freq).
		WithNumReqPerCycle(4).
		WithLatency(10).
		WithPageTable(b.pageTable)

	if b.useCoalescingTLBPort {
		builder = builder.UseCoalescingTLBPort()
	}
	l2TLB := builder.Build(fmt.Sprintf("%s.L2TLB", chiplet.name))
	b.l2TLBs = append(b.l2TLBs, l2TLB)
	b.gpu.L2TLBs = append(b.gpu.L2TLBs, l2TLB)
	chiplet.L2TLBs = append(chiplet.L2TLBs, l2TLB)

	if b.enableVisTracing {
		tracing.CollectTrace(l2TLB, b.visTracer)
	}
}
