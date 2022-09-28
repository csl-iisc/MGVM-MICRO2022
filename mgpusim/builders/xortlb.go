package builders

import (
	"fmt"
	"math"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem"
	"gitlab.com/akita/mem/cache"
	"gitlab.com/akita/mem/cache/writeback"
	"gitlab.com/akita/mem/idealmemcontroller"
	"gitlab.com/akita/mem/vm/tlb"
	"gitlab.com/akita/mgpusim"
	"gitlab.com/akita/mgpusim/remotetranslation"
	"gitlab.com/akita/util/tracing"
)

type XORTLBGPUBuilder struct {
	*CommonBuilder
	switchL2TLBStriping bool
	usePtCaching        bool
	// specific components
}

// XOR TLB specific function

// MakeXORTLBGPUBuilder provides a GPU builder that can builds MCM GPU.
func MakeXORTLBGPUBuilder() XORTLBGPUBuilder {
	// TODO: should this be using new? is the object being allocated on the stack?
	cbp := CommonBuilder{}
	b := XORTLBGPUBuilder{CommonBuilder: &cbp, switchL2TLBStriping: false, usePtCaching: false}
	b.SetDefaultCommonBuilderParams()
	return b
}

func (b XORTLBGPUBuilder) Build(name string, id uint64) *mgpusim.GPU {
	// fmt.Println("log2 page size: ", b.log2PageSize)
	b.createGPU(name, id)

	b.buildCP()
	b.cp.SwitchL2TLBStriping(b.switchL2TLBStriping)

	chipRdmaAddressTable := b.createChipRDMAAddrTable()
	rdmaResponsePorts := make([]akita.Port, b.numChiplet)
	remoteAddressTranslationTable := b.createRemoteAddrTransTable()
	rtuResponsePorts := make([]akita.Port, b.numChiplet)

	for i := 0; i < b.numChiplet; i++ {
		chipletName := fmt.Sprintf("%s.chiplet_%02d", b.gpuName, i)
		chiplet := NewChiplet(chipletName, uint64(i))
		b.BuildSAs(chiplet)
		b.configChipRDMAEngine(chiplet, chipRdmaAddressTable, rdmaResponsePorts)
		if b.usePtCaching {
			b.buildMemBanks(chiplet)
		} else {
			b.CommonBuilder.buildMemBanks(chiplet)
		}

		b.buildMMU(chiplet)
		b.buildL2TLB(chiplet)

		b.configRemoteAddressTranslationUnit(chiplet, remoteAddressTranslationTable, rtuResponsePorts)

		b.connectL1ToL2(chiplet)
		if b.usePtCaching {
			b.connectL2ToDRAM(chiplet)
		} else {
			b.CommonBuilder.connectL2ToDRAM(chiplet)
		}
		b.connectL1TLBToL2TLB(chiplet)
		b.connectL2TLBTOMMU(chiplet)
		if b.usePtCaching {
			b.connectMMUToL2(chiplet)
		} else {
			b.CommonBuilder.connectMMUToL2(chiplet)
		}

		b.chiplets = append(b.chiplets, chiplet)
	}

	b.buildPageMigrationController()
	b.setupDMA()

	b.connectCP()
	b.setupInterchipNetwork()

	return b.gpu
}

func (b *XORTLBGPUBuilder) createRemoteAddrTransTable() *cache.XORLowModuleFinder {

	log2RemoteTLBInterleaving := int(math.Log2(float64(b.remoteTLBInterleavingSize)))
	numBits := int(math.Log2(float64(b.numChiplet)))
	remoteAddrTransTable :=
		cache.NewXORLowModuleFinder(b.numChiplet, 4, numBits, int(b.log2PageSize)+log2RemoteTLBInterleaving)
	b.cp.InterleavingSize = uint64(log2RemoteTLBInterleaving)
	// uint64(int(b.log2PageSize) + log2RemoteTLBInterleaving)
	b.cp.Log2PageSize = b.log2PageSize
	fmt.Println("setting the page size in the cp:", b.cp.InterleavingSize, b.cp.Log2PageSize)
	return remoteAddrTransTable
}

func (b *XORTLBGPUBuilder) configRemoteAddressTranslationUnit(chiplet *Chiplet,
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
	if b.useCoalescingRTU {
		chiplet.L2TLBs[0].(*tlb.LatTLB).ToRTU = chiplet.remoteTranslationUnit.(*remotetranslation.CoalescingRemoteTranslationUnit).ToL2
	} else {
		chiplet.L2TLBs[0].(*tlb.LatTLB).ToRTU = chiplet.remoteTranslationUnit.(*remotetranslation.DefaultRTU).ToL2
	}
	chiplet.remoteTranslationUnit.SetRemoteAddressTranslationTable(
		remoteAddressTranslationTable)
	remoteAddressTranslationTable.LowModules =
		append(remoteAddressTranslationTable.LowModules,
			chiplet.remoteTranslationUnit.GetRequestPort())
	b.gpu.RemoteAddressTranslationUnits =
		append(b.gpu.RemoteAddressTranslationUnits, chiplet.remoteTranslationUnit)
}

func (b *XORTLBGPUBuilder) connectL1TLBToL2TLB(chiplet *Chiplet) {
	tlbConn := akita.NewDirectConnection(chiplet.name+"L1TLB-L2TLB",
		b.engine, b.freq)
	tlbConn.PlugIn(chiplet.L2TLBs[0].GetTopPort(), 64)

	var lowModuleFinder cache.LowModuleFinder

	// find a better way to do this!!!
	log2RemoteTLBInterleaving := int(math.Log2(float64(b.remoteTLBInterleavingSize)))

	numBits := int(math.Log2(float64(b.numChiplet)))
	interleavedLowModuleFinder := cache.NewLocalXORLowModuleFinder(
		chiplet.ChipletID, uint64(b.numChiplet), 4, numBits, int(b.log2PageSize)+log2RemoteTLBInterleaving)

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

func (b *XORTLBGPUBuilder) SwitchL2TLBStriping(switchStriping bool) {
	b.switchL2TLBStriping = switchStriping
}

func (b *XORTLBGPUBuilder) UsePtCaching(ptCaching bool) {
	b.usePtCaching = ptCaching
}

func (b *XORTLBGPUBuilder) connectMMUToL2(chiplet *Chiplet) {
	lowModuleFinder := cache.NewInterleavedLowModuleFinder(
		1 << b.log2MemoryBankInterleavingSize)
	for _, l2 := range chiplet.L2Caches {
		lowModuleFinder.LowModules = append(lowModuleFinder.LowModules,
			l2.TopPort)
	}
	chiplet.MMU.SetLowModuleFinder(lowModuleFinder)
	chiplet.L1ToL2Connection.PlugIn(chiplet.MMU.TranslationPort, 64)
}

func (b *XORTLBGPUBuilder) buildMemBanks(chiplet *Chiplet) {
	// memCtrlBuilder := b.createDRAMControllerBuilder()
	l2Builder := writeback.MakeBuilder().
		WithEngine(b.engine).
		WithFreq(b.freq).
		WithLog2BlockSize(b.log2CacheLineSize).
		WithWayAssociativity(16).
		WithByteSize(256 * mem.KB).
		WithNumMSHREntry(64).
		WithNumReqPerCycle(2)

	for i := 0; i < b.numMemoryBankPerChiplet; i++ {
		dramName := fmt.Sprintf("%s.DRAM_%d", chiplet.name, i)
		dram := idealmemcontroller.New(
			dramName, b.engine, 512*mem.MB)
		addrConverter := idealmemcontroller.InterleavingConverter{
			InterleavingSize:    1 << b.log2MemoryBankInterleavingSize,
			TotalNumOfElements:  b.numChiplet * b.numMemoryBankPerChiplet,
			CurrentElementIndex: b.numMemoryBankPerChiplet*int(chiplet.ChipletID) + i,
			Offset:              b.memAddrOffset, // + b.memoryPerChiplet*chiplet.ChipletID,
		}
		dram.AddressConverter = addrConverter

		b.drams = append(b.drams, dram)
		b.gpu.MemoryControllers = append(b.gpu.MemoryControllers, dram)
		chiplet.DRAMs = append(chiplet.DRAMs, dram)

		if b.enableVisTracing {
			tracing.CollectTrace(dram, b.visTracer)
		}
		cacheName := fmt.Sprintf("%s.L2_%d", chiplet.name, i)
		l2 := l2Builder.Build(cacheName)
		b.l2Caches = append(b.l2Caches, l2)
		b.gpu.L2Caches = append(b.gpu.L2Caches, l2)
		chiplet.L2Caches = append(chiplet.L2Caches, l2)
		lowModuleFinder := cache.NewStripedLocalVRemoteLowModuleFinder(b.memAddrOffset, uint64(b.numChiplet*b.numMemoryBankPerChiplet),
			1<<b.log2MemoryBankInterleavingSize, 8*chiplet.ChipletID+uint64(i), 8*chiplet.ChipletID+uint64(i))
		// lowModuleFinder.ModuleForOtherAddresses = chiplet.chipRdmaEngine.ToL1

		// lowModuleFinder := cache.NewInterleavedLowModuleFinder(
		// 1 << b.log2MemoryBankInterleavingSize)
		lowModuleFinder.ModuleForOtherAddresses = chiplet.chipRdmaEngine.PwPort
		lowModuleFinder.LowModules = append(lowModuleFinder.LowModules, dram.ToTop)
		l2.SetLowModuleFinder(lowModuleFinder)

		// l2.SetLowModuleFinder(&cache.SingleLowModuleFinder{
		// 	LowModule: dram.ToTop,
		// })
		if b.enableVisTracing {
			tracing.CollectTrace(l2, b.visTracer)
		}
	}
}

func (b *XORTLBGPUBuilder) connectL2ToDRAM(chiplet *Chiplet) {
	chiplet.L2ToDramConnection = akita.NewDirectConnection(
		chiplet.name+".L2-DRAM", b.engine, b.freq)

	lowModuleFinder := cache.NewInterleavedLowModuleFinder(
		1 << b.log2MemoryBankInterleavingSize)

	for _, l2 := range chiplet.L2Caches {
		chiplet.L2ToDramConnection.PlugIn(l2.BottomPort, 64)
	}

	for _, dram := range chiplet.DRAMs {
		chiplet.L2ToDramConnection.PlugIn(dram.ToTop, 64)
		lowModuleFinder.LowModules = append(lowModuleFinder.LowModules,
			dram.ToTop)
	}
	chiplet.L2ToDramConnection.PlugIn(chiplet.chipRdmaEngine.PwPort, 64)
	b.pageMigrationController.MemCtrlFinder = lowModuleFinder
	chiplet.L2ToDramConnection.PlugIn(b.pageMigrationController.LocalMemPort, 16)
}
