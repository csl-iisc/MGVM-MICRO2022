package builders

import (
	"fmt"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem/cache"
	"gitlab.com/akita/mem/vm/tlb"
	"gitlab.com/akita/mgpusim"
	"gitlab.com/akita/mgpusim/rdma"
	"gitlab.com/akita/noc/networking/chipnetwork"
)

type PrivateH2TLBGPUBuilder struct {
	*CommonBuilder

	// specific componenets
	interChipletMagicNetwork *akita.DirectConnection
}

// Distributed TLB specific function

// MakeDistributedTLBGPUBuilder provides a GPU builder that can builds MCM GPU.
func MakePrivateH2TLBGPUBuilder() PrivateH2TLBGPUBuilder {
	// TODO: should this be using new? is the object being allocated on the stack?
	cbp := CommonBuilder{}
	b := PrivateH2TLBGPUBuilder{CommonBuilder: &cbp}
	b.SetDefaultCommonBuilderParams()
	return b
}

func (b PrivateH2TLBGPUBuilder) Build(name string, id uint64) *mgpusim.GPU {
	b.createGPU(name, id)

	b.buildCP()

	chipRdmaAddressTable := b.createChipRDMAAddrTable()
	rdmaResponsePorts := make([]akita.Port, 4)
	pageRdmaAddressTable := b.createChipRDMAAddrTable()
	pageRdmaResponsePorts := make([]akita.Port, 4)

	for i := 0; i < b.numChiplet; i++ {
		chipletName := fmt.Sprintf("%s.chiplet_%02d", b.gpuName, i)
		chiplet := NewChiplet(chipletName, uint64(i))

		b.BuildSAs(chiplet)
		b.buildMemBanks(chiplet)
		b.buildMMU(chiplet)
		// TODO: this may have to be overridden
		b.buildL2TLB(chiplet)

		b.configChipRDMAEngine(chiplet, chipRdmaAddressTable, rdmaResponsePorts)
		b.configPageRDMAEngine(chiplet, pageRdmaAddressTable, pageRdmaResponsePorts)
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
	b.setupInterchipMagicNetwork()

	return b.gpu
}

func (b *PrivateH2TLBGPUBuilder) connectCP() {
	b.internalConn = akita.NewDirectConnection(
		b.gpuName+"InternalConn", b.engine, b.freq)
	b.gpu.InternalConnection = b.internalConn

	b.internalConn.PlugIn(b.cp.ToDriver, 1)
	b.internalConn.PlugIn(b.cp.ToDMA, 128)
	b.internalConn.PlugIn(b.cp.ToCaches, 128)
	b.internalConn.PlugIn(b.cp.ToCUs, 128)
	b.internalConn.PlugIn(b.cp.ToTLBs, 128)
	b.internalConn.PlugIn(b.cp.ToAddressTranslators, 128)
	b.internalConn.PlugIn(b.cp.ToRDMA, 4)
	b.internalConn.PlugIn(b.cp.ToPMC, 4)

	b.internalConn.PlugIn(b.cp.ToRTU, 4)
	b.internalConn.PlugIn(b.cp.ToMMUs, 4)

	b.cp.RDMA = b.rdmaEngine.CtrlPort
	b.internalConn.PlugIn(b.cp.RDMA, 1)

	b.cp.DMAEngine = b.dmaEngine.ToCP
	b.internalConn.PlugIn(b.dmaEngine.ToCP, 1)

	b.cp.PMC = b.pageMigrationController.CtrlPort
	b.internalConn.PlugIn(b.pageMigrationController.CtrlPort, 1)

	b.connectCPWithCUs()
	b.connectCPWithAddressTranslators()
	b.connectCPWithCaches()
	b.connectCPWithMMUs()
	b.connectCPWithTLBs()
}

func (b *PrivateH2TLBGPUBuilder) connectL1TLBToL2TLB(chiplet *Chiplet) {
	tlbConn := akita.NewDirectConnection(chiplet.name+"L1TLB-L2TLB",
		b.engine, b.freq)
	tlbConn.PlugIn(chiplet.L2TLBs[0].GetTopPort(), 64)

	var lowModuleFinder cache.LowModuleFinder

	singeLowModuleFinder := new(cache.SingleLowModuleFinder)
	singeLowModuleFinder.LowModule = chiplet.L2TLBs[0].GetTopPort()

	chiplet.L2TLBs[0].(*tlb.LatTLB).TLBFinder = singeLowModuleFinder

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

func (b *PrivateH2TLBGPUBuilder) setupInterchipNetwork() {
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

func (b *PrivateH2TLBGPUBuilder) InterChipletPorts(c *Chiplet) []akita.Port {
	ports := []akita.Port{
		c.chipRdmaEngine.RequestPort,
		c.chipRdmaEngine.ResponsePort,
	}
	return ports
}
func (b *PrivateH2TLBGPUBuilder) configPageRDMAEngine(
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
func (b *PrivateH2TLBGPUBuilder) connectMMUToL2(chiplet *Chiplet) {
	pageTableLowModuleFinder := cache.NewStripedLocalVRemoteLowModuleFinder(b.memAddrOffset, uint64(b.numChiplet*b.numMemoryBankPerChiplet),
		1<<b.log2MemoryBankInterleavingSize, uint64(b.numMemoryBankPerChiplet)*chiplet.ChipletID, uint64(b.numMemoryBankPerChiplet)*chiplet.ChipletID+uint64(b.numMemoryBankPerChiplet-1))
	pageTableLowModuleFinder.ModuleForOtherAddresses = chiplet.pageRdmaEngine.ToL1

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

func (b *PrivateH2TLBGPUBuilder) setupInterchipMagicNetwork() {
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

func (b *PrivateH2TLBGPUBuilder) InterChipletMagicPorts(c *Chiplet) []akita.Port {
	ports := []akita.Port{
		c.pageRdmaEngine.RequestPort,
		c.pageRdmaEngine.ResponsePort,
	}
	return ports
}
