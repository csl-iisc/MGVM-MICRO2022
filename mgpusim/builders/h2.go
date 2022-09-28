package builders

import (
	"fmt"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem/cache"
	"gitlab.com/akita/mgpusim"
	"gitlab.com/akita/mgpusim/rdma"
	"gitlab.com/akita/noc/networking/chipnetwork"
)

type H2Builder struct {
	*CommonBuilder

	// specific componenets
	interChipletMagicNetwork *akita.DirectConnection
}

// Distributed TLB specific function

// MakeDistributedTLBGPUBuilder provides a GPU builder that can builds MCM GPU.
func MakeH2Builder() H2Builder {
	// TODO: should this be using new? is the object being allocated on the stack?
	cbp := CommonBuilder{}
	b := H2Builder{CommonBuilder: &cbp}
	b.SetDefaultCommonBuilderParams()
	return b
}

func (b H2Builder) Build(name string, id uint64) *mgpusim.GPU {
	b.createGPU(name, id)

	b.buildCP()

	chipRdmaAddressTable := b.createChipRDMAAddrTable()
	rdmaResponsePorts := make([]akita.Port, 4)
	pageRdmaAddressTable := b.createChipRDMAAddrTable()
	pageRdmaResponsePorts := make([]akita.Port, 4)
	remoteAddressTranslationTable := b.createRemoteAddrTransTable()
	rtuResponsePorts := make([]akita.Port, 4)

	for i := 0; i < b.numChiplet; i++ {
		chipletName := fmt.Sprintf("%s.chiplet_%02d", b.gpuName, i)
		chiplet := NewChiplet(chipletName, uint64(i))

		b.BuildSAs(chiplet)
		b.buildMemBanks(chiplet)
		b.buildMMU(chiplet)
		b.buildL2TLB(chiplet)

		b.configChipRDMAEngine(chiplet, chipRdmaAddressTable, rdmaResponsePorts)
		b.configPageRDMAEngine(chiplet, pageRdmaAddressTable, pageRdmaResponsePorts)
		b.configRemoteAddressTranslationUnit(chiplet, remoteAddressTranslationTable, rtuResponsePorts)

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
	b.setupInterchipMagicNetwork()

	return b.gpu
}

func (b *H2Builder) configPageRDMAEngine(
	chiplet *Chiplet,
	addrTable *cache.StripedLowModuleFinder, pageRdmaResponsePorts []akita.Port) {
	chiplet.pageRdmaEngine = rdma.NewEngine(
		fmt.Sprintf("%s.PageRDMA", chiplet.name),
		b.engine,
		chiplet.lowModuleFinderForL1,
		nil,
	)
	pageRdmaResponsePorts[chiplet.ChipletID] = chiplet.pageRdmaEngine.ResponsePort
	chiplet.pageRdmaEngine.ResponsePorts = pageRdmaResponsePorts
	chiplet.pageRdmaEngine.RemoteRDMAAddressTable = addrTable
	addrTable.LowModules = append(addrTable.LowModules,
		chiplet.pageRdmaEngine.RequestPort)
	b.gpu.PageRDMAEngines = append(b.gpu.PageRDMAEngines, chiplet.pageRdmaEngine)
}

// call only aftet connectL1ToL2
func (b *H2Builder) connectMMUToL2(chiplet *Chiplet) {
	pageTableLowModuleFinder := cache.NewInterleavedLowModuleFinder(
		1 << b.log2MemoryBankInterleavingSize)
	pageTableLowModuleFinder.ModuleForOtherAddresses = chiplet.pageRdmaEngine.ToL1
	pageTableLowModuleFinder.UseAddressSpaceLimitation = true
	pageTableLowModuleFinder.LowAddress = b.memAddrOffset +
		b.memoryPerChiplet*chiplet.ChipletID
	pageTableLowModuleFinder.HighAddress = pageTableLowModuleFinder.LowAddress + b.memoryPerChiplet

	chiplet.pageRdmaEngine.SetLocalModuleFinder(pageTableLowModuleFinder)

	l1ToL2Conn := chiplet.L1ToL2Connection
	l1ToL2Conn.PlugIn(chiplet.pageRdmaEngine.ToL1, 64)
	l1ToL2Conn.PlugIn(chiplet.pageRdmaEngine.ToL2, 64)

	for _, l2 := range chiplet.L2Caches {
		pageTableLowModuleFinder.LowModules = append(pageTableLowModuleFinder.LowModules,
			l2.TopPort)
	}

	chiplet.MMU.SetLowModuleFinder(pageTableLowModuleFinder)
	l1ToL2Conn.PlugIn(chiplet.MMU.TranslationPort, 64)
}
func (b *H2Builder) setupInterchipNetwork() {
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

func (b *H2Builder) InterChipletPorts(c *Chiplet) []akita.Port {
	ports := []akita.Port{
		c.remoteTranslationUnit.GetRequestPort(),
		c.remoteTranslationUnit.GetResponsePort(),
	}
	return ports
}

func (b *H2Builder) InterChipletMagicPorts(c *Chiplet) []akita.Port {
	ports := []akita.Port{
		c.chipRdmaEngine.RequestPort,
		c.chipRdmaEngine.ResponsePort,
	}
	return ports
}

func (b *H2Builder) setupInterchipMagicNetwork() {
	interchipDirectConnection := akita.NewDirectConnection("magic",
		b.engine, 1*akita.GHz)
	for _, chiplet := range b.chiplets {
		for _, port := range b.InterChipletMagicPorts(chiplet) {
			interchipDirectConnection.PlugIn(port, 64)
		}
	}
	b.interChipletMagicNetwork = interchipDirectConnection
	b.gpu.InterChipletMagicNetwork = interchipDirectConnection
}
