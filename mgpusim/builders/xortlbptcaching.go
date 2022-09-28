package builders

import (
	"fmt"
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem"
	"gitlab.com/akita/mem/cache"
	"gitlab.com/akita/mem/cache/writeback"
	"gitlab.com/akita/mem/idealmemcontroller"
	"gitlab.com/akita/mem/vm/tlb"
	"gitlab.com/akita/mgpusim"
	"gitlab.com/akita/mgpusim/remotetranslation"
	"gitlab.com/akita/util/tracing"
	"math"
)

type XorWithPTCachingGPUBuilder struct {
	*CommonBuilder

	// specific components
}

// XOR TLB specific function

// MakeXorWithPTCachingGPUBuilder provides a GPU builder that can builds MCM GPU.
func MakeXorWithPTCachingGPUBuilder() XorWithPTCachingGPUBuilder {
	// TODO: should this be using new? is the object being allocated on the stack?
	cbp := CommonBuilder{}
	b := XorWithPTCachingGPUBuilder{CommonBuilder: &cbp}
	b.SetDefaultCommonBuilderParams()
	return b
}

func (b XorWithPTCachingGPUBuilder) Build(name string, id uint64) *mgpusim.GPU {
	b.createGPU(name, id)

	b.buildCP()

	chipRdmaAddressTable := b.createChipRDMAAddrTable()
	rdmaResponsePorts := make([]akita.Port, 4)
	remoteAddressTranslationTable := b.createRemoteAddrTransTable()
	rtuResponsePorts := make([]akita.Port, 4)

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

func (b *XorWithPTCachingGPUBuilder) createRemoteAddrTransTable() *cache.XORLowModuleFinder {

	log2RemoteTLBInterleaving := int(math.Log2(float64(b.remoteTLBInterleavingSize)))
	remoteAddrTransTable :=
		cache.NewXORLowModuleFinder(b.numChiplet, 8, 2, int(b.log2PageSize)+log2RemoteTLBInterleaving)
	b.cp.InterleavingSize = uint64(log2RemoteTLBInterleaving) //uint64(12 + log2RemoteTLBInterleaving)
	return remoteAddrTransTable
}

func (b *XorWithPTCachingGPUBuilder) configRemoteAddressTranslationUnit(chiplet *Chiplet,
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

func (b *XorWithPTCachingGPUBuilder) connectL1TLBToL2TLB(chiplet *Chiplet) {
	tlbConn := akita.NewDirectConnection(chiplet.name+"L1TLB-L2TLB",
		b.engine, b.freq)
	tlbConn.PlugIn(chiplet.L2TLBs[0].GetTopPort(), 64)

	var lowModuleFinder cache.LowModuleFinder

	// find a better way to do this!!!
	log2RemoteTLBInterleaving := int(math.Log2(float64(b.remoteTLBInterleavingSize)))

	interleavedLowModuleFinder := cache.NewLocalXORLowModuleFinder(
		chiplet.ChipletID, uint64(b.numChiplet), 4, 2, 12+log2RemoteTLBInterleaving)

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

func (b *XorWithPTCachingGPUBuilder) connectMMUToL2(chiplet *Chiplet) {
	lowModuleFinder := cache.NewInterleavedLowModuleFinder(
		1 << b.log2MemoryBankInterleavingSize)
	for _, l2 := range chiplet.L2Caches {
		lowModuleFinder.LowModules = append(lowModuleFinder.LowModules,
			l2.TopPort)
	}
	chiplet.MMU.SetLowModuleFinder(lowModuleFinder)
	chiplet.L1ToL2Connection.PlugIn(chiplet.MMU.TranslationPort, 64)
}

func (b *XorWithPTCachingGPUBuilder) buildMemBanks(chiplet *Chiplet) {
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
			TotalNumOfElements:  b.numMemoryBankPerChiplet,
			CurrentElementIndex: i,
			Offset:              b.memAddrOffset + b.memoryPerChiplet*chiplet.ChipletID,
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
		lowModuleFinder := cache.NewInterleavedLowModuleFinder(
			1 << b.log2MemoryBankInterleavingSize)
		lowModuleFinder.ModuleForOtherAddresses = chiplet.chipRdmaEngine.PwPort
		lowModuleFinder.UseAddressSpaceLimitation = true
		lowModuleFinder.LowAddress = b.memAddrOffset +
			b.memoryPerChiplet*chiplet.ChipletID
		lowModuleFinder.HighAddress = lowModuleFinder.LowAddress + b.memoryPerChiplet
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

func (b *XorWithPTCachingGPUBuilder) connectL2ToDRAM(chiplet *Chiplet) {
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
