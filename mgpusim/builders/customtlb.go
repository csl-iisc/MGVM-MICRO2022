package builders

import (
	"fmt"
	"math"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem/cache"
	"gitlab.com/akita/mem/vm/tlb"
	"gitlab.com/akita/mgpusim"
	"gitlab.com/akita/mgpusim/remotetranslation"
	"gitlab.com/akita/mgpusim/timing/cp"
)

type CustomTLBGPUBuilder struct {
	*CommonBuilder
	switchL2TLBStriping bool
	usePtCaching        bool
	// specific components

	customHSLpmdUnits uint64
}

// Custom TLB specific function

// MakeCustomTLBGPUBuilder provides a GPU builder that can builds MCM GPU.
func MakeCustomTLBGPUBuilder() CustomTLBGPUBuilder {
	cbp := CommonBuilder{}
	b := CustomTLBGPUBuilder{CommonBuilder: &cbp, switchL2TLBStriping: false, usePtCaching: false}
	b.SetDefaultCommonBuilderParams()
	return b
}

func (b CustomTLBGPUBuilder) Build(name string, id uint64) *mgpusim.GPU {
	b.createGPU(name, id)

	b.buildCP()
	b.cp.SwitchL2TLBStriping(b.switchL2TLBStriping)

	chipRdmaAddressTable := b.createChipRDMAAddrTable()
	rdmaResponsePorts := make([]akita.Port, b.numChiplet)
	remoteAddressTranslationTable := b.createRemoteAddrTransTable()
	rtuResponsePorts := make([]akita.Port, b.numChiplet)

	// some parameters hardcoded
	hsl := func(address uint64) uint64 {
		return (address / uint64(b.customHSLpmdUnits)) % uint64(b.numChiplet)
	}
	remoteAddressTranslationTable.Hashfunc = hsl

	for i := 0; i < b.numChiplet; i++ {
		chipletName := fmt.Sprintf("%s.chiplet_%02d", b.gpuName, i)
		chiplet := NewChiplet(chipletName, uint64(i))
		b.BuildSAs(chiplet)
		b.configChipRDMAEngine(chiplet, chipRdmaAddressTable, rdmaResponsePorts)
		b.buildMemBanks(chiplet)

		b.buildMMU(chiplet)
		b.buildL2TLB(chiplet)

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

	b.connectCP()
	b.setupInterchipNetwork()

	return b.gpu
}

func (b *CustomTLBGPUBuilder) createRemoteAddrTransTable() *cache.HashLowModuleFinder {

	log2RemoteTLBInterleaving := int(math.Log2(float64(b.remoteTLBInterleavingSize)))
	remoteAddrTransTable := cache.NewHashLowModuleFinder(b.log2PageSize)
	b.cp.InterleavingSize = uint64(log2RemoteTLBInterleaving)
	b.cp.Log2PageSize = b.log2PageSize
	fmt.Println("setting the page size in the cp:", b.cp.InterleavingSize, b.cp.Log2PageSize)
	return remoteAddrTransTable
}

func (b *CustomTLBGPUBuilder) configRemoteAddressTranslationUnit(chiplet *Chiplet,
	remoteAddressTranslationTable *cache.HashLowModuleFinder,
	rtuResponsePorts []akita.Port) {
	chiplet.remoteTranslationUnit = remotetranslation.NewRemoteTranslationUnit(
		fmt.Sprintf("%s.RTU", chiplet.name),
		b.engine,
		nil,
		nil,
	)
	rtuResponsePorts[chiplet.ChipletID] = chiplet.remoteTranslationUnit.GetResponsePort()
	chiplet.remoteTranslationUnit.SetResponsePorts(rtuResponsePorts)
	chiplet.remoteTranslationUnit.(*remotetranslation.DefaultRTU).L2CtrlPort = chiplet.L2TLBs[0].GetControlPort()
	chiplet.L2TLBs[0].(*tlb.LatTLB).ToRTU = chiplet.remoteTranslationUnit.(*remotetranslation.DefaultRTU).ToL2
	chiplet.remoteTranslationUnit.SetRemoteAddressTranslationTable(
		remoteAddressTranslationTable)
	remoteAddressTranslationTable.LowModules =
		append(remoteAddressTranslationTable.LowModules,
			chiplet.remoteTranslationUnit.GetRequestPort())
	b.gpu.RemoteAddressTranslationUnits =
		append(b.gpu.RemoteAddressTranslationUnits, chiplet.remoteTranslationUnit)
}

func (b *CustomTLBGPUBuilder) connectL1TLBToL2TLB(chiplet *Chiplet) {
	tlbConn := akita.NewDirectConnection(chiplet.name+"L1TLB-L2TLB",
		b.engine, b.freq)
	tlbConn.PlugIn(chiplet.L2TLBs[0].GetTopPort(), 64)

	var lowModuleFinder cache.LowModuleFinder

	interleavedLowModuleFinder := cache.NewCustomTwoLevelLowModuleFinder(
		b.log2PageSize, chiplet.ChipletID)

	// some parameters hardcoded
	hsl := func(address uint64) uint64 {
		return (address / uint64(b.customHSLpmdUnits)) % uint64(b.numChiplet)
	}
	interleavedLowModuleFinder.Hashfunc = hsl

	interleavedLowModuleFinder.LocalLowModule = chiplet.L2TLBs[0].GetTopPort()
	chiplet.L2TLBs[0].(*tlb.LatTLB).TLBFinder = interleavedLowModuleFinder
	chiplet.L2TLBs[0].(*tlb.LatTLB).HashFuncHSL = hsl

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

func (b *CustomTLBGPUBuilder) SwitchL2TLBStriping(switchStriping bool) {
	b.switchL2TLBStriping = switchStriping
}

func (b *CustomTLBGPUBuilder) WithCustomHSL(customHSLpmdUnits uint64) {
	b.customHSLpmdUnits = customHSLpmdUnits
}

// over riding buildCP as CustomTLB needs customHSLpmdUnits as part of CP
func (b *CustomTLBGPUBuilder) buildCP() {
	builder := cp.MakeBuilder().
		WithEngine(b.engine).
		WithFreq(b.freq).
		WithAlg(b.alg).
		WithPartition(b.partition).
		WithCustomHSLpmdUnits(b.customHSLpmdUnits)

	if !b.disableProgressBar {
		builder = builder.ShowProgressBar()
	}

	if b.enableVisTracing {
		builder = builder.WithVisTracer(b.visTracer)
	}

	b.cp = builder.Build(b.gpuName + ".CommandProcessor")
	b.gpu.CommandProcessor = b.cp

	//TODO: are these per die or per GPU??
	b.buildDMAEngine()
	b.buildRDMAEngine()
	b.buildPageMigrationController()
}
