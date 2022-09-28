package builders

import (
	"fmt"
	"math"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem/cache"
	"gitlab.com/akita/mem/vm/tlb"
	"gitlab.com/akita/mgpusim"
	"gitlab.com/akita/mgpusim/rdma"
	"gitlab.com/akita/mgpusim/remotetranslation"
	"gitlab.com/akita/noc/networking/chipnetwork"
)

type XORTLBH2GPUBuilder struct {
	*CommonBuilder

	// specific componenets
	interChipletMagicNetwork *akita.DirectConnection
}

// Distributed TLB specific function

// MakeDistributedTLBGPUBuilder provides a GPU builder that can builds MCM GPU.
func MakeXORTLBH2GPUBuilder() XORTLBH2GPUBuilder {
	// TODO: should this be using new? is the object being allocated on the stack?
	cbp := CommonBuilder{}
	b := XORTLBH2GPUBuilder{CommonBuilder: &cbp}
	b.SetDefaultCommonBuilderParams()
	return b
}

func (b XORTLBH2GPUBuilder) Build(name string, id uint64) *mgpusim.GPU {
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

func (b *XORTLBH2GPUBuilder) configPageRDMAEngine(
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
func (b *XORTLBH2GPUBuilder) connectMMUToL2(chiplet *Chiplet) {
	pageTableLowModuleFinder := cache.NewStripedLocalVRemoteLowModuleFinder(b.memAddrOffset, uint64(b.numChiplet*b.numMemoryBankPerChiplet),
		1<<b.log2MemoryBankInterleavingSize, uint64(b.numMemoryBankPerChiplet)*chiplet.ChipletID, uint64(b.numMemoryBankPerChiplet)*chiplet.ChipletID+uint64(b.numMemoryBankPerChiplet-1))
	pageTableLowModuleFinder.ModuleForOtherAddresses = chiplet.pageRdmaEngine.ToL1
	// pageTableLowModuleFinder.UseAddressSpaceLimitation = true
	// pageTableLowModuleFinder.LowAddress = b.memAddrOffset +
	// 	b.memoryPerChiplet*chiplet.ChipletID
	// pageTableLowModuleFinder.HighAddress = pageTableLowModuleFinder.LowAddress + b.memoryPerChiplet

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

func (b *XORTLBH2GPUBuilder) setupInterchipNetwork() {
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

func (b *XORTLBH2GPUBuilder) InterChipletPorts(c *Chiplet) []akita.Port {
	ports := []akita.Port{
		c.chipRdmaEngine.RequestPort,
		c.chipRdmaEngine.ResponsePort,
		c.remoteTranslationUnit.GetRequestPort(),
		c.remoteTranslationUnit.GetResponsePort(),
	}
	return ports
}

func (b *XORTLBH2GPUBuilder) InterChipletMagicPorts(c *Chiplet) []akita.Port {
	ports := []akita.Port{
		c.pageRdmaEngine.RequestPort,
		c.pageRdmaEngine.ResponsePort,
	}
	return ports
}

func (b *XORTLBH2GPUBuilder) setupInterchipMagicNetwork() {
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

func (b *XORTLBH2GPUBuilder) createRemoteAddrTransTable() *cache.XORLowModuleFinder {

	log2RemoteTLBInterleaving := int(math.Log2(float64(b.remoteTLBInterleavingSize)))
	remoteAddrTransTable :=
		cache.NewXORLowModuleFinder(b.numChiplet, 4, 2, int(b.log2PageSize)+log2RemoteTLBInterleaving)
	b.cp.InterleavingSize = uint64(log2RemoteTLBInterleaving) //uint64(int(b.log2PageSize) + log2RemoteTLBInterleaving)
	return remoteAddrTransTable
}

func (b *XORTLBH2GPUBuilder) configRemoteAddressTranslationUnit(chiplet *Chiplet,
	remoteAddressTranslationTable *cache.XORLowModuleFinder,
	rtuResponsePorts []akita.Port) {
	if b.useCoalescingRTU {
		chiplet.remoteTranslationUnit = remotetranslation.NewCoalescingRemoteTranslationUnit(
			fmt.Sprintf("%s.RTU", chiplet.name),
			b.engine,
			nil,
			nil,
		)
	} else {
		chiplet.remoteTranslationUnit = remotetranslation.NewRemoteTranslationUnit(
			fmt.Sprintf("%s.RTU", chiplet.name),
			b.engine,
			nil,
			nil,
		)
	}
	rtuResponsePorts[chiplet.ChipletID] = chiplet.remoteTranslationUnit.GetResponsePort()
	chiplet.remoteTranslationUnit.SetResponsePorts(rtuResponsePorts)
	chiplet.remoteTranslationUnit.(*remotetranslation.DefaultRTU).L2CtrlPort = chiplet.L2TLBs[0].GetControlPort()

	chiplet.remoteTranslationUnit.SetRemoteAddressTranslationTable(
		remoteAddressTranslationTable)
	remoteAddressTranslationTable.LowModules =
		append(remoteAddressTranslationTable.LowModules,
			chiplet.remoteTranslationUnit.GetRequestPort())
	b.gpu.RemoteAddressTranslationUnits =
		append(b.gpu.RemoteAddressTranslationUnits, chiplet.remoteTranslationUnit)
}

func (b *XORTLBH2GPUBuilder) connectL1TLBToL2TLB(chiplet *Chiplet) {
	tlbConn := akita.NewDirectConnection(chiplet.name+"L1TLB-L2TLB",
		b.engine, b.freq)
	tlbConn.PlugIn(chiplet.L2TLBs[0].GetTopPort(), 64)

	var lowModuleFinder cache.LowModuleFinder

	// find a better way to do this!!!
	log2RemoteTLBInterleaving := int(math.Log2(float64(b.remoteTLBInterleavingSize)))

	interleavedLowModuleFinder := cache.NewLocalXORLowModuleFinder(
		chiplet.ChipletID, uint64(b.numChiplet), 4, 2, int(b.log2PageSize)+log2RemoteTLBInterleaving)

	interleavedLowModuleFinder.LocalLowModule = chiplet.L2TLBs[0].GetTopPort()
	chiplet.L2TLBs[0].(*tlb.LatTLB).TLBFinder = interleavedLowModuleFinder

	interleavedLowModuleFinder.RemoteLowModule = chiplet.remoteTranslationUnit.GetL1Port()

	lowModuleFinder = interleavedLowModuleFinder

	chiplet.remoteTranslationUnit.SetLocalModuleFinder(lowModuleFinder)
	tlbConn.PlugIn(chiplet.remoteTranslationUnit.GetL1Port(), 64)
	tlbConn.PlugIn(chiplet.remoteTranslationUnit.GetL2Port(), 64)

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
