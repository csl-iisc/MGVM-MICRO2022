package builders

import (
	"fmt"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem/cache"
	"gitlab.com/akita/mgpusim"
	"gitlab.com/akita/noc/networking/chipnetwork"
)

type H3Builder struct {
	*CommonBuilder

	// specific componenets
}

// Distributed TLB specific function

// MakeDistributedTLBGPUBuilder provides a GPU builder that can builds MCM GPU.
func MakeH3Builder() H3Builder {
	// TODO: should this be using new? is the object being allocated on the stack?
	cbp := CommonBuilder{}
	b := H3Builder{CommonBuilder: &cbp}
	b.SetDefaultCommonBuilderParams()
	return b
}

func (b H3Builder) Build(name string, id uint64) *mgpusim.GPU {
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
		b.buildL2TLB(chiplet)

		b.configChipRDMAEngine(chiplet, chipRdmaAddressTable, rdmaResponsePorts)
		// b.configRemoteAddressTranslationUnit(chiplet, remoteAddressTranslationTable, rtuResponsePorts)

		b.connectL1ToL2(chiplet)
		b.connectL2ToDRAM(chiplet)
		// b.connectL1TLBToL2TLB(chiplet)
		b.connectL2TLBTOMMU(chiplet)
		b.connectMMUToL2(chiplet)

		b.chiplets = append(b.chiplets, chiplet)
	}

	b.connectL1TLBToL2TLBMagic()

	b.buildPageMigrationController()
	b.setupDMA()

	// b.setupMMUs()
	b.connectCP()
	b.setupInterchipNetwork()

	return b.gpu
}

func (b *H3Builder) connectL1TLBToL2TLBMagic() {
	tlbConn := akita.NewDirectConnection("L1TLB-L2TLB",
		b.engine, b.freq)

	var lowModuleFinder cache.LowModuleFinder

	interleavedLowModuleFinder := cache.NewInterleavedLowModuleFinder(
		b.remoteTLBInterleavingSize * 4096)

	for _, l2tlb := range b.gpu.L2TLBs {
		tlbConn.PlugIn(l2tlb.GetTopPort(), 64)
		interleavedLowModuleFinder.LowModules = append(interleavedLowModuleFinder.LowModules, l2tlb.GetTopPort())
	}

	lowModuleFinder = interleavedLowModuleFinder

	for _, l1vTLB := range b.gpu.L1VTLBs {
		l1vTLB.SetLowModuleFinder(lowModuleFinder)
		tlbConn.PlugIn(l1vTLB.BottomPort, 16)
	}

	for _, l1iTLB := range b.gpu.L1ITLBs {
		l1iTLB.SetLowModuleFinder(lowModuleFinder)
		tlbConn.PlugIn(l1iTLB.BottomPort, 16)
	}

	for _, l1sTLB := range b.gpu.L1STLBs {
		l1sTLB.SetLowModuleFinder(lowModuleFinder)
		tlbConn.PlugIn(l1sTLB.BottomPort, 16)
	}
}

func (b *H3Builder) setupInterchipNetwork() {
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
	// b.interChipletNetwork = chipConnector
	// b.gpu.InterChipletNetwork = chipConnector
}

func (b *H3Builder) InterChipletPorts(c *Chiplet) []akita.Port {
	ports := []akita.Port{
		c.chipRdmaEngine.RequestPort,
		c.chipRdmaEngine.ResponsePort,
		// c.remoteTranslationUnit.RequestPort,
		// c.remoteTranslationUnit.ResponsePort,
	}
	return ports
}
