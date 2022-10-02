package builders

import (
	"fmt"
	"math"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem"
	"gitlab.com/akita/mem/cache"
	"gitlab.com/akita/mem/cache/writeback"
	"gitlab.com/akita/mem/device"
	"gitlab.com/akita/mem/idealmemcontroller"
	"gitlab.com/akita/mem/vm/addresstranslator"
	"gitlab.com/akita/mem/vm/mmu"
	"gitlab.com/akita/mem/vm/tlb"
	"gitlab.com/akita/mgpusim"
	"gitlab.com/akita/mgpusim/pagemigrationcontroller"
	"gitlab.com/akita/mgpusim/rdma"
	"gitlab.com/akita/mgpusim/remotetranslation"
	"gitlab.com/akita/mgpusim/timing/caches/l1v"
	"gitlab.com/akita/mgpusim/timing/caches/rob"
	"gitlab.com/akita/mgpusim/timing/cp"
	"gitlab.com/akita/mgpusim/timing/cu"
	"gitlab.com/akita/noc/networking/chipnetwork"
	"gitlab.com/akita/util/tracing"
)

type Builder interface {
	Build(string, uint64) *mgpusim.GPU
	WithMemAddrOffset(uint64)
	WithISADebugging()
	WithTLBTracer(tracing.Tracer)
	WithMemTracer(tracing.Tracer)
	WithVisTracer(tracing.Tracer)

	WithEngine(akita.Engine)
	WithNumCUPerShaderArray(int)
	WithNumShaderArrayPerChiplet(int)
	WithNumMemoryBankPerChiplet(int)
	WithNumChiplet(int)
	WithTotalMem(uint64)
	CalculateMemoryParameters()
	WithLog2PageSize(uint64)
	WithPageTable(*device.PageTableImpl)
	WithAlg(string)
	WithSchedulingPartition(string)
}

// Every component, whether used in a specific builder or not, appears here.
// Further, every componenet appears in gpu.go
// Finally, nil the components that aren't going to be used in a specific builder.
type CommonBuilder struct {
	engine                         akita.Engine
	freq                           akita.Freq
	memAddrOffset                  uint64
	totalMem                       uint64
	mmu                            *mmu.MMUImpl
	numChiplet                     int
	numShaderArrayPerChiplet       int
	numCUPerShaderArray            int
	numMemoryBankPerChiplet        int
	memoryPerChiplet               uint64
	memoryPerBank                  uint64
	log2PageSize                   uint64
	log2CacheLineSize              uint64
	log2MemoryBankInterleavingSize uint64
	remoteTLBInterleavingSize      uint64
	walkersPerChiplet              int
	pageTable                      *device.PageTableImpl

	enableISADebugging bool
	enableMemTracing   bool
	enableTLBTracing   bool
	enableVisTracing   bool
	disableProgressBar bool
	visTracer          tracing.Tracer
	tlbTracer          tracing.Tracer
	memTracer          tracing.Tracer

	gpuName                 string
	gpu                     *mgpusim.GPU
	cp                      *cp.CommandProcessor
	chiplets                []*Chiplet
	cus                     []*cu.ComputeUnit
	l1vReorderBuffers       []*rob.ReorderBuffer
	l1iReorderBuffers       []*rob.ReorderBuffer
	l1sReorderBuffers       []*rob.ReorderBuffer
	l1vCaches               []*l1v.Cache
	l1sCaches               []*l1v.Cache
	l1iCaches               []*l1v.Cache
	l2Caches                []*writeback.Cache
	l1vAddrTrans            []addresstranslator.AddressTranslator
	l1sAddrTrans            []addresstranslator.AddressTranslator
	l1iAddrTrans            []addresstranslator.AddressTranslator
	l1vTLBs                 []*tlb.TLB
	l1sTLBs                 []*tlb.TLB
	l1iTLBs                 []*tlb.TLB
	l2TLBs                  []tlb.L2TLB
	drams                   []*idealmemcontroller.Comp
	lowModuleFinderForL1    cache.LowModuleFinder
	lowModuleFinderForL2    *cache.InterleavedLowModuleFinder
	lowModuleFinderForPMC   *cache.InterleavedLowModuleFinder
	dmaEngine               *cp.DMAEngine
	rdmaEngine              *rdma.Engine
	pageMigrationController *pagemigrationcontroller.PageMigrationController

	internalConn *akita.DirectConnection

	interChipletNetwork *chipnetwork.Connector

	alg                  string
	partition            string
	useCoalescingTLBPort bool
	useCoalescingRTU     bool
}

// MakeCommonBuilder provides a GPU builder that can builds the MCM GPU.
func (b *CommonBuilder) SetDefaultCommonBuilderParams() {
	b.freq = 1 * akita.GHz
	b.numShaderArrayPerChiplet = 16
	b.numCUPerShaderArray = 4
	b.numMemoryBankPerChiplet = 8
	b.log2CacheLineSize = 6
	b.log2PageSize = 12
	b.log2MemoryBankInterleavingSize = 12
	b.walkersPerChiplet = 16

	//this ought not to be in common
	b.remoteTLBInterleavingSize = 512

	// set other defaults here
	b.useCoalescingRTU = false
	b.useCoalescingTLBPort = false

}

// WithEngine sets the engine that the GPU use.
func (b *CommonBuilder) WithEngine(engine akita.Engine) {
	b.engine = engine
}

// WithFreq sets the frequency that the GPU works at.
func (b *CommonBuilder) WithFreq(freq akita.Freq) {
	b.freq = freq
}

// WithMemAddrOffset sets the address of the first byte of the GPU to build.
func (b *CommonBuilder) WithMemAddrOffset(offset uint64) {
	b.memAddrOffset = offset
}

// WithTotalMem sets the total amount of memory in the GPU.
func (b *CommonBuilder) WithTotalMem(memory uint64) {
	b.totalMem = memory
}

// WithMMU sets the MMU component that provides the address translation service
// for the GPU.
func (b *CommonBuilder) WithMMU(mmu *mmu.MMUImpl) {
	b.mmu = mmu
}

// WithNumMemoryBankPerChiplet sets the number of L2 cache modules and number
// of memory controllers in each chiplet.
func (b *CommonBuilder) WithNumMemoryBankPerChiplet(n int) {
	b.numMemoryBankPerChiplet = n
}

// WithNumChiplet sets the number of chiplets
func (b *CommonBuilder) WithNumChiplet(n int) {
	b.numChiplet = n
}

// WithNumShaderArrayPerChiplet sets the number of shader arrays in each chiplet
func (b *CommonBuilder) WithNumShaderArrayPerChiplet(n int) {
	b.numShaderArrayPerChiplet = n
}

// WithNumCUPerShaderArray sets the number of CU and number of L1V caches in
// each Shader Array.
func (b *CommonBuilder) WithNumCUPerShaderArray(n int) {
	b.numCUPerShaderArray = n
}

// WithLog2MemoryBankInterleavingSize sets the number of consecutive bytes that
// are guaranteed to be on a memory bank.
func (b *CommonBuilder) WithLog2MemoryBankInterleavingSize(n uint64) {
	b.log2MemoryBankInterleavingSize = n
}

func (b *CommonBuilder) WithRemoteTLB(interleavingSize uint64) {
	b.remoteTLBInterleavingSize = interleavingSize

}

// WithVisTracer applies a tracer to trace all the tasks of all the GPU
// components
func (b *CommonBuilder) WithVisTracer(t tracing.Tracer) {
	b.enableVisTracing = true
	b.visTracer = t
}

// WithTLBTracer applies a tracer to trace the TLBory transactions.
func (b *CommonBuilder) WithTLBTracer(t tracing.Tracer) {
	b.enableTLBTracing = true
	b.tlbTracer = t
}

// WithMemTracer applies a tracer to trace the memory transactions.
func (b *CommonBuilder) WithMemTracer(t tracing.Tracer) {
	b.enableMemTracing = true
	b.memTracer = t
}

// WithISADebugging enables the GPU to dump instruction execution information.
func (b *CommonBuilder) WithISADebugging() {
	b.enableISADebugging = true
}

// WithoutProgressBar disables the progress bar for kernel execution
func (b *CommonBuilder) WithoutProgressBar() {
	b.disableProgressBar = true
}

// WithLog2CacheLineSize sets the cache line size with the power of 2.
func (b *CommonBuilder) WithLog2CacheLineSize(
	log2CacheLine uint64,
) {
	b.log2CacheLineSize = log2CacheLine
}

// WithLog2PageSize sets the page size with the power of 2.
func (b *CommonBuilder) WithLog2PageSize(log2PageSize uint64) {
	b.log2PageSize = log2PageSize
}

// WithPageTable sets the page size with the power of 2.
func (b *CommonBuilder) WithPageTable(pt *device.PageTableImpl) {
	b.pageTable = pt
}

// WithAlg sets the scheduling algorithm
func (b *CommonBuilder) WithAlg(alg string) {
	b.alg = alg
}

func (b *CommonBuilder) WithSchedulingPartition(partition string) {
	b.partition = partition
}

func (b *CommonBuilder) UseCoalescingTLBPort(u bool) {
	b.useCoalescingTLBPort = u
}

func (b *CommonBuilder) UseCoalescingRTU(u bool) {
	b.useCoalescingRTU = u
}

func (b *CommonBuilder) WithWalkersPerChiplet(w int) {
	b.walkersPerChiplet = w
}

// CalculateMemoryParameters calculates
// -> memoryPerChiplet
// -> memoryPerBank
// based on the totalMem, numChiplets and numMemoryBankPerChiplet
func (b *CommonBuilder) CalculateMemoryParameters() {
	b.memoryPerChiplet = b.totalMem / uint64(b.numChiplet)
	b.memoryPerBank =
		b.totalMem / uint64(b.numChiplet*b.numMemoryBankPerChiplet)
}

func (b *CommonBuilder) createGPU(name string, id uint64) {
	b.gpuName = name

	b.gpu = mgpusim.NewGPU(b.gpuName)

	b.gpu.GPUID = id
}

func (b *CommonBuilder) buildCP() {
	builder := cp.MakeBuilder().
		WithEngine(b.engine).
		WithFreq(b.freq).
		WithAlg(b.alg).
		WithPartition(b.partition)

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

func (b *CommonBuilder) buildDMAEngine() {
	b.dmaEngine = cp.NewDMAEngine(
		fmt.Sprintf("%s.DMA", b.gpuName),
		b.engine,
		nil)

	if b.enableVisTracing {
		tracing.CollectTrace(b.dmaEngine, b.visTracer)
	}
}

// BuildSAs builds shader arrays.
func (b *CommonBuilder) BuildSAs(chiplet *Chiplet) {
	saBuilder := makeShaderArrayBuilder()
	saBuilder.withEngine(b.engine)
	saBuilder.withFreq(b.freq)
	saBuilder.withGPUID(b.gpu.GPUID)
	saBuilder.withLog2CachelineSize(b.log2CacheLineSize)
	saBuilder.withLog2PageSize(b.log2PageSize)
	saBuilder.withNumCU(b.numCUPerShaderArray)

	if b.enableVisTracing {
		saBuilder.withVisTracer(b.visTracer)
	}

	for i := 0; i < b.numShaderArrayPerChiplet; i++ {
		saName := fmt.Sprintf("%s.SA_%02d", chiplet.name, i)
		sa := saBuilder.Build(saName)
		b.collectSAComponents(sa, chiplet)
	}
}

func (b *CommonBuilder) collectSAComponents(
	sa shaderArray,
	chiplet *Chiplet,
) {
	for _, cu := range sa.cus {
		b.gpu.CUs = append(b.gpu.CUs, cu)
		b.cus = append(b.cus, cu)
		chiplet.CUs = append(chiplet.CUs, cu)
	}

	for _, rob := range sa.l1vROBs {
		b.l1vReorderBuffers = append(b.l1vReorderBuffers, rob)
		b.gpu.L1VROBs = append(b.gpu.L1VROBs, rob)
		chiplet.L1VROBs = append(chiplet.L1VROBs, rob)
	}

	for _, tlb := range sa.l1vTLBs {
		b.l1vTLBs = append(b.l1vTLBs, tlb)
		b.gpu.L1VTLBs = append(b.gpu.L1VTLBs, tlb)
		chiplet.L1VTLBs = append(chiplet.L1VTLBs, tlb)
	}

	for _, l1v := range sa.l1vCaches {
		b.l1vCaches = append(b.l1vCaches, l1v)
		b.gpu.L1VCaches = append(b.gpu.L1VCaches, l1v)
		chiplet.L1VCaches = append(chiplet.L1VCaches, l1v)
	}

	for _, at := range sa.l1vATs {
		b.l1vAddrTrans = append(b.l1vAddrTrans, at)
		b.gpu.L1VAddrTranslator = append(b.gpu.L1VAddrTranslator, at)
		chiplet.L1VAddrTranslator = append(chiplet.L1VAddrTranslator, at)
	}

	b.l1sAddrTrans = append(b.l1sAddrTrans, sa.l1sAT)
	b.gpu.L1SAddrTranslator = append(b.gpu.L1SAddrTranslator, sa.l1sAT)
	chiplet.L1SAddrTranslator = append(chiplet.L1SAddrTranslator, sa.l1sAT)
	b.l1sReorderBuffers = append(b.l1sReorderBuffers, sa.l1sROB)
	b.gpu.L1SROBs = append(b.gpu.L1SROBs, sa.l1sROB)
	chiplet.L1SROBs = append(chiplet.L1SROBs, sa.l1sROB)
	b.l1sCaches = append(b.l1sCaches, sa.l1sCache)
	b.gpu.L1SCaches = append(b.gpu.L1SCaches, sa.l1sCache)
	chiplet.L1SCaches = append(chiplet.L1SCaches, sa.l1sCache)
	b.l1sTLBs = append(b.l1sTLBs, sa.l1sTLB)
	b.gpu.L1STLBs = append(b.gpu.L1STLBs, sa.l1sTLB)
	chiplet.L1STLBs = append(chiplet.L1STLBs, sa.l1sTLB)

	b.l1iAddrTrans = append(b.l1iAddrTrans, sa.l1iAT)
	b.gpu.L1IAddrTranslator = append(b.gpu.L1IAddrTranslator, sa.l1iAT)
	chiplet.L1IAddrTranslator = append(chiplet.L1IAddrTranslator, sa.l1iAT)
	b.l1iReorderBuffers = append(b.l1iReorderBuffers, sa.l1iROB)
	b.gpu.L1IROBs = append(b.gpu.L1IROBs, sa.l1iROB)
	chiplet.L1IROBs = append(chiplet.L1IROBs, sa.l1iROB)
	b.l1iCaches = append(b.l1iCaches, sa.l1iCache)
	b.gpu.L1ICaches = append(b.gpu.L1ICaches, sa.l1iCache)
	chiplet.L1ICaches = append(chiplet.L1ICaches, sa.l1iCache)
	b.l1iTLBs = append(b.l1iTLBs, sa.l1iTLB)
	b.gpu.L1ITLBs = append(b.gpu.L1ITLBs, sa.l1iTLB)
	chiplet.L1ITLBs = append(chiplet.L1ITLBs, sa.l1iTLB)
}

func (b *CommonBuilder) buildMemBanks(chiplet *Chiplet) {
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
			Offset:              b.memAddrOffset,
			//  + b.memoryPerChiplet*chiplet.ChipletID,
		}
		// fmt.Println("^^^^^", b.numMemoryBankPerChiplet*int(chiplet.ChipletID)+i)
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
		l2.SetLowModuleFinder(&cache.SingleLowModuleFinder{
			LowModule: dram.ToTop,
		})
		if b.enableVisTracing {
			tracing.CollectTrace(l2, b.visTracer)
		}
	}
}

// func (b *CommonBuilder) createDRAMControllerBuilder() dram.Builder {
// 	memBankSize := 2 * mem.GB / uint64(b.numMemoryBankPerChiplet)
// 	if 2*mem.GB%uint64(b.numMemoryBankPerChiplet) != 0 {
// 		panic("GPU memory size is not a multiple of the number of memory banks")
// 	}
// 	dramCol := 64
// 	dramRow := 4096
// 	dramDeviceWidth := 32
// 	dramBankSize := dramCol * dramRow * dramDeviceWidth
// 	dramBank := 4
// 	dramBankGroup := 1
// 	dramBusWidth := 256
// 	dramDevicePerRank := dramBusWidth / dramDeviceWidth
// 	dramRankSize := dramBankSize * dramDevicePerRank * dramBank
// 	dramRank := int(memBankSize) / dramRankSize

// 	memCtrlBuilder := dram.MakeBuilder().
// 		WithEngine(b.engine).
// 		WithFreq(500 * akita.MHz).
// 		WithProtocol(dram.GDDR5).
// 		WithBurstLength(8).
// 		WithDeviceWidth(dramDeviceWidth).
// 		WithBusWidth(dramBusWidth).
// 		WithNumChannel(1).
// 		WithNumRank(dramRank).
// 		WithNumBankGroup(dramBankGroup).
// 		WithNumBank(dramBank).
// 		WithNumCol(dramCol).
// 		WithNumRow(dramRow).
// 		WithCommandQueueSize(8).
// 		WithTransactionQueueSize(32).
// 		WithTCL(24).
// 		WithTCWL(7).
// 		WithTRCDRD(18).
// 		WithTRCDWR(15).
// 		WithTRP(18).
// 		WithTRAS(42).
// 		WithTREFI(11699).
// 		WithTRRDS(9).
// 		WithTRRDL(9).
// 		WithTWTRS(8).
// 		WithTWTRL(8).
// 		WithTWR(18).
// 		WithTCCDS(2).
// 		WithTCCDL(3).
// 		WithTRTRS(0).
// 		WithTRTP(3).
// 		WithTPPD(2)

// 	if b.visTracer != nil {
// 		memCtrlBuilder = memCtrlBuilder.WithAdditionalTracer(b.visTracer)
// 	}

// 	return memCtrlBuilder
// }

func (b *CommonBuilder) buildL2TLB(chiplet *Chiplet) {
	numSets := 64 // 128 // 256 // changed this here
	numWays := 8  // 8 // changed this here
	log2NumSets := int(math.Log2(float64(numSets)))

	tlbIndexBitsStart := int(math.Log2(float64(b.remoteTLBInterleavingSize))) + int(b.log2PageSize) + 1
	tlbIndexBitsEnd := tlbIndexBitsStart + int(math.Log2(float64(b.numChiplet))) - 1

	mask := uint64(0)
	t := uint64(1) << b.log2PageSize
	numBitsSet := 0
	for i := int(b.log2PageSize) + 1; i <= 64; i++ {
		if i < tlbIndexBitsStart || i > tlbIndexBitsEnd {
			mask = mask | t
			numBitsSet++
			if numBitsSet == log2NumSets {
				break
			}
		}
		t = t << 1
	}
	//fmt.Println(mask)
	builder := tlb.MakeLatTLBBuilder().
		WithEngine(b.engine).
		WithFreq(b.freq).
		WithNumWays(numWays).
		WithNumSets(numSets).
		WithNumMSHREntry(64).
		WithNumReqPerCycle(4).
		WithLog2PageSize(b.log2PageSize).
		WithLowModule(chiplet.MMU.ToTop).
		WithIndexingMask(mask).
		WithLatency(10)
	fmt.Println("num TLB sets:", numSets)
	fmt.Println("num TLB ways:", numWays)
	if b.useCoalescingTLBPort {
		builder = builder.UseCoalescingTLBPort()
	}
	l2TLB := builder.Build(fmt.Sprintf("%s.L2TLB", chiplet.name))
	l2TLB.SetLowModuleFinder(&cache.SingleLowModuleFinder{
		LowModule: chiplet.MMU.ToTop,
	})

	b.l2TLBs = append(b.l2TLBs, l2TLB)
	b.gpu.L2TLBs = append(b.gpu.L2TLBs, l2TLB)
	chiplet.L2TLBs = append(chiplet.L2TLBs, l2TLB)

	if b.enableVisTracing {
		tracing.CollectTrace(l2TLB, b.visTracer)
	}
}

func (b *CommonBuilder) buildMMU(chiplet *Chiplet) {
	mmuBuilder := mmu.MakeBuilder().
		WithEngine(b.engine).
		WithFreq(1 * akita.GHz).
		WithLog2PageSize(b.log2PageSize).
		WithPageTable(b.pageTable).
		WithNumChiplets(uint64(b.numChiplet)).
		//		WithLowAddr(b.memAddrOffset).
		//		WithTotMem(b.totalMem).
		//		WithBankSize(b.memoryPerChiplet).
		//		WithNumMemoryBankPerChiplet(uint64(b.numMemoryBankPerChiplet)).
		WithMaxNumReqInFlight(b.walkersPerChiplet) // changed this here
		//TODO: try increasing the number of walkers to 16 and see what happens
	chiplet.MMU = mmuBuilder.Build(fmt.Sprintf("%s.MMU", chiplet.name))
	chiplet.MMU.CommandProcessor = b.gpu.CommandProcessor.ToMMUs
	b.gpu.MMUs = append(b.gpu.MMUs, chiplet.MMU)
	// mmu.ToTop =
	// b.l2TLBs = append(b.l2TLBs, l2TLB)
	// b.gpu.L2TLBs = append(b.gpu.L2TLBs, l2TLB)
	// chiplet.L2TLBs = append(chiplet.L2TLBs, l2TLB)

	// if b.enableVisTracing {
	// 	tracing.CollectTrace(l2TLB, b.visTracer)
	// }
}

// func getChipletNumFromName(chipletName string) int {
// 	fmt.Println(chipletName)
// 	chipletNum, err := strconv.Atoi(chipletName[14:15])
// 	if err != nil {
// 		panic("error")
// 	}
// 	return chipletNum
// }

// Though the name of function is configChipRDMAEngine, it actually creates
// the Chip RDMA engine too..
func (b *CommonBuilder) configChipRDMAEngine(
	chiplet *Chiplet,
	addrTable *cache.StripedLowModuleFinder, rdmaResponsePorts []akita.Port) {

	chiplet.chipRdmaEngine = rdma.NewEngine(
		fmt.Sprintf("%s.ChipRDMA", chiplet.name),
		b.engine,
		chiplet.lowModuleFinderForL1,
		nil,
	)
	rdmaResponsePorts[chiplet.ChipletID] = chiplet.chipRdmaEngine.ResponsePort
	chiplet.chipRdmaEngine.ResponsePorts = rdmaResponsePorts
	chiplet.chipRdmaEngine.RemoteRDMAAddressTable = addrTable
	addrTable.LowModules = append(addrTable.LowModules,
		chiplet.chipRdmaEngine.RequestPort)
	b.gpu.ChipRDMAEngines = append(b.gpu.ChipRDMAEngines, chiplet.chipRdmaEngine)
}

func (b *CommonBuilder) configRemoteAddressTranslationUnit(
	chiplet *Chiplet,
	remoteAddressTranslationTable *cache.InterleavedLowModuleFinder, rtuResponsePorts []akita.Port) {
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

	chiplet.remoteTranslationUnit.SetRemoteAddressTranslationTable(
		remoteAddressTranslationTable)
	remoteAddressTranslationTable.LowModules =
		append(remoteAddressTranslationTable.LowModules,
			chiplet.remoteTranslationUnit.GetRequestPort())
	b.gpu.RemoteAddressTranslationUnits =
		append(b.gpu.RemoteAddressTranslationUnits, chiplet.remoteTranslationUnit)
}

func (b *CommonBuilder) connectCP() {
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
	b.connectCPWithTLBs()
	b.connectCPWithCaches()
	b.connectCPWithRTUs()
	b.connectCPWithMMUs()
}

func (b *CommonBuilder) connectMMUToL2(chiplet *Chiplet) {
	chiplet.MMU.SetLowModuleFinder(chiplet.lowModuleFinderForL1)
	chiplet.L1ToL2Connection.PlugIn(chiplet.MMU.TranslationPort, 64)
}

func (b *CommonBuilder) connectCPWithCUs() {
	for _, chiplet := range b.chiplets {
		for _, cu := range chiplet.CUs {
			b.cp.RegisterCU(cu)
			b.internalConn.PlugIn(cu.ToACE, 1)
			b.internalConn.PlugIn(cu.ToCP, 1)
		}
	}
}

func (b *CommonBuilder) connectCPWithAddressTranslators() {
	for _, chiplet := range b.chiplets {
		for _, at := range chiplet.L1VAddrTranslator {
			b.cp.AddressTranslators = append(b.cp.AddressTranslators, at.GetCtrlPort())
			b.internalConn.PlugIn(at.GetCtrlPort(), 1)
		}

		for _, at := range chiplet.L1SAddrTranslator {
			b.cp.AddressTranslators = append(b.cp.AddressTranslators, at.GetCtrlPort())
			b.internalConn.PlugIn(at.GetCtrlPort(), 1)
		}

		for _, at := range chiplet.L1IAddrTranslator {
			b.cp.AddressTranslators = append(b.cp.AddressTranslators, at.GetCtrlPort())
			b.internalConn.PlugIn(at.GetCtrlPort(), 1)
		}

		for _, rob := range chiplet.L1VROBs {
			b.cp.AddressTranslators = append(
				b.cp.AddressTranslators, rob.ControlPort)
			b.internalConn.PlugIn(rob.ControlPort, 1)
		}

		for _, rob := range chiplet.L1IROBs {
			b.cp.AddressTranslators = append(
				b.cp.AddressTranslators, rob.ControlPort)
			b.internalConn.PlugIn(rob.ControlPort, 1)
		}

		for _, rob := range chiplet.L1SROBs {
			b.cp.AddressTranslators = append(
				b.cp.AddressTranslators, rob.ControlPort)
			b.internalConn.PlugIn(rob.ControlPort, 1)
		}
	}
}

func (b *CommonBuilder) connectCPWithTLBs() {
	for _, chiplet := range b.chiplets {
		for _, tlb := range chiplet.L2TLBs {
			b.cp.TLBs = append(b.cp.TLBs, tlb.GetControlPort())
			b.internalConn.PlugIn(tlb.GetControlPort(), 10)
			// added to setup reverse comm from L2 TLB to CP
			tlb.SetCommandProcessor(b.cp.ToTLBs)
		}

		// 	for _, tlb := range chiplet.L1VTLBs {
		// 		b.cp.TLBs = append(b.cp.TLBs, tlb.ControlPort)
		// 		b.internalConn.PlugIn(tlb.ControlPort, 1)
		// 	}

		// 	for _, tlb := range chiplet.L1STLBs {
		// 		b.cp.TLBs = append(b.cp.TLBs, tlb.ControlPort)
		// 		b.internalConn.PlugIn(tlb.ControlPort, 1)
		// 	}

		// 	for _, tlb := range chiplet.L1ITLBs {
		// 		b.cp.TLBs = append(b.cp.TLBs, tlb.ControlPort)
		// 		b.internalConn.PlugIn(tlb.ControlPort, 1)
		// 	}
	}
}

func (b *CommonBuilder) connectCPWithCaches() {
	for _, chiplet := range b.chiplets {
		for _, c := range chiplet.L1ICaches {
			b.cp.L1ICaches = append(b.cp.L1ICaches, c.ControlPort)
			b.internalConn.PlugIn(c.ControlPort, 1)
		}

		for _, c := range chiplet.L1VCaches {
			b.cp.L1VCaches = append(b.cp.L1VCaches, c.ControlPort)
			b.internalConn.PlugIn(c.ControlPort, 1)
		}

		for _, c := range chiplet.L1SCaches {
			b.cp.L1SCaches = append(b.cp.L1SCaches, c.ControlPort)
			b.internalConn.PlugIn(c.ControlPort, 1)
		}

		for _, c := range chiplet.L2Caches {
			b.cp.L2Caches = append(b.cp.L2Caches, c.ControlPort)
			b.internalConn.PlugIn(c.ControlPort, 1)
		}
	}
}

func (b *CommonBuilder) connectCPWithRTUs() {
	for _, chiplet := range b.chiplets {
		rtu := chiplet.remoteTranslationUnit
		b.cp.RTUs = append(b.cp.RTUs, rtu.GetControlPort())
		b.internalConn.PlugIn(rtu.GetControlPort(), 10)
		rtu.SetCommandProcessor(b.cp.ToRTU)
	}
}

func (b *CommonBuilder) connectCPWithMMUs() {
	for _, chiplet := range b.chiplets {
		mmu := chiplet.MMU
		b.cp.MMUs = append(b.cp.MMUs, mmu.ControlPort)
		b.internalConn.PlugIn(mmu.ControlPort, 10)
	}
}

func (b *CommonBuilder) connectL1ToL2(chiplet *Chiplet) {
	fmt.Println("memory address offset:", b.memAddrOffset)
	lowModuleFinder := cache.NewStripedLocalVRemoteLowModuleFinder(b.memAddrOffset, uint64(b.numChiplet*b.numMemoryBankPerChiplet),
		1<<b.log2MemoryBankInterleavingSize, uint64(b.numMemoryBankPerChiplet)*chiplet.ChipletID, uint64(b.numMemoryBankPerChiplet)*chiplet.ChipletID+uint64(b.numMemoryBankPerChiplet-1))
	lowModuleFinder.ModuleForOtherAddresses = chiplet.chipRdmaEngine.ToL1
	// lowModuleFinder.UseAddressSpaceLimitation = true
	// lowModuleFinder.LowAddress = b.memAddrOffset +
	// b.memoryPerChiplet*chiplet.ChipletID
	// lowModuleFinder.HighAddress = lowModuleFinder.LowAddress + b.memoryPerChiplet

	l1ToL2Conn := akita.NewDirectConnection(chiplet.name+".L1-L2",
		b.engine, b.freq)

	chiplet.chipRdmaEngine.SetLocalModuleFinder(lowModuleFinder)
	l1ToL2Conn.PlugIn(chiplet.chipRdmaEngine.ToL1, 64)
	l1ToL2Conn.PlugIn(chiplet.chipRdmaEngine.ToL2, 64)

	for _, l2 := range chiplet.L2Caches {
		lowModuleFinder.LowModules = append(lowModuleFinder.LowModules,
			l2.TopPort)
		l1ToL2Conn.PlugIn(l2.TopPort, 64)
	}

	for _, l1v := range chiplet.L1VCaches {
		l1v.SetLowModuleFinder(lowModuleFinder)
		l1ToL2Conn.PlugIn(l1v.BottomPort, 16)
	}

	for _, l1s := range chiplet.L1SCaches {
		l1s.SetLowModuleFinder(lowModuleFinder)
		l1ToL2Conn.PlugIn(l1s.BottomPort, 16)
	}

	for _, l1iAT := range chiplet.L1IAddrTranslator {
		l1iAT.SetLowModuleFinder(lowModuleFinder)
		l1ToL2Conn.PlugIn(l1iAT.GetBottomPort(), 16)
	}
	chiplet.lowModuleFinderForL1 = lowModuleFinder
	// fmt.Println(lowModuleFinder)
	chiplet.L1ToL2Connection = l1ToL2Conn
}

func (b *CommonBuilder) connectL2ToDRAM(chiplet *Chiplet) {
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

	// b.dmaEngine.SetLocalDataSource(lowModuleFinder)
	// chiplet.L2ToDramConnection.PlugIn(b.dmaEngine.ToMem, 64)

	b.pageMigrationController.MemCtrlFinder = lowModuleFinder
	chiplet.L2ToDramConnection.PlugIn(b.pageMigrationController.LocalMemPort, 16)
}

func (b *CommonBuilder) connectL1TLBToL2TLB(chiplet *Chiplet) {
	tlbConn := akita.NewDirectConnection(chiplet.name+"L1TLB-L2TLB",
		b.engine, b.freq)
	tlbConn.PlugIn(chiplet.L2TLBs[0].GetTopPort(), 64)

	var lowModuleFinder cache.LowModuleFinder

	interleavedLowModuleFinder := cache.NewLocalInterleavedLowModuleFinder(
		chiplet.ChipletID, uint64(b.numChiplet), b.remoteTLBInterleavingSize*4096)
	interleavedLowModuleFinder.LocalLowModule = chiplet.L2TLBs[0].GetTopPort()

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

func (b *CommonBuilder) connectL2TLBTOMMU(chiplet *Chiplet) {
	tlbToMMUConn := akita.NewDirectConnection(chiplet.name+".L2TLB-MMU",
		b.engine, b.freq)
	tlbToMMUConn.PlugIn(chiplet.MMU.ToTop, 64)
	for _, l2tlb := range chiplet.L2TLBs {
		tlbToMMUConn.PlugIn(l2tlb.GetBottomPort(), 16)
	}

}

func (b *CommonBuilder) buildPageMigrationController() {
	b.pageMigrationController =
		pagemigrationcontroller.NewPageMigrationController(
			fmt.Sprintf("%s.PMC", b.gpuName),
			b.engine,
			b.lowModuleFinderForPMC,
			nil)

	b.gpu.PMC = b.pageMigrationController
}

func (b *CommonBuilder) setupDMA() {
	lowModuleFinder := cache.NewStripedLowModuleFinder(b.memAddrOffset, uint64(b.numChiplet*b.numMemoryBankPerChiplet), 1<<b.log2MemoryBankInterleavingSize)
	// cache.NewStripedLowModuleFinder(b.memAddrOffset, 32, )
	// cache.NewTwoLevelLowModuleFinder(
	// b.memAddrOffset, b.memAddrOffset+b.totalMem,
	// b.memoryPerChiplet, uint64(b.numChiplet),
	// 4096, uint64(b.numMemoryBankPerChiplet))

	for _, chiplet := range b.chiplets {
		toMem := akita.NewLimitNumMsgPort(b.dmaEngine, 64, chiplet.name+".ToMem")
		b.dmaEngine.ToMem = append(b.dmaEngine.ToMem, toMem)

		for _, dram := range chiplet.DRAMs {
			lowModuleFinder.LowModules = append(lowModuleFinder.LowModules,
				dram.ToTop)
			b.dmaEngine.DramToChipletMap[dram.ToTop] = toMem
		}

		chiplet.L2ToDramConnection.PlugIn(toMem, 64)
	}

	b.dmaEngine.SetLocalDataSource(lowModuleFinder)
	// b.mmu.SetLowModuleFinder(lowModuleFinder)

	b.pageMigrationController.MemCtrlFinder = lowModuleFinder
	for _, chiplet := range b.chiplets {
		chiplet.L2ToDramConnection.PlugIn(b.pageMigrationController.LocalMemPort, 16)
	}
}

func (b *CommonBuilder) setupInterchipNetwork() {
	chipConnector := chipnetwork.NewInterChipletConnector().
		WithEngine(b.engine).
		WithSwitchLatency(360). // changed this here //720 // 180 // 360
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

func (b *CommonBuilder) buildRDMAEngine() {
	b.rdmaEngine = rdma.NewEngine(
		fmt.Sprintf("%s.RDMA", b.gpuName),
		b.engine,
		b.lowModuleFinderForL1,
		nil,
	)
	b.gpu.RDMAEngine = b.rdmaEngine
}

func (b *CommonBuilder) createChipRDMAAddrTable() *cache.StripedLowModuleFinder {
	// chipRdmaAddressTable := new(cache.BankedLowModuleFinder)
	// chipRdmaAddressTable.BankSize = 8 * b.log2MemoryBankInterleavingSize
	// chipRdmaAddressTable.MemAddrOffset = 2 * mem.GB
	//b.memoryPerChiplet
	// chipRdmaAddressTable.LowModules = append(chipRdmaAddressTable.LowModules, nil)
	chipRdmaAddressTable := cache.NewStripedLowModuleFinder(b.memAddrOffset, uint64(b.numChiplet), uint64(b.numMemoryBankPerChiplet)*(1<<b.log2MemoryBankInterleavingSize))
	return chipRdmaAddressTable
}

func (b *CommonBuilder) createRemoteAddrTransTable() *cache.InterleavedLowModuleFinder {
	remoteAddrTransTable :=
		cache.NewInterleavedLowModuleFinder(b.remoteTLBInterleavingSize * 4096)
	return remoteAddrTransTable
}

// InterChipletPorts returns the list of ports that are connected to other
// chiplets of the GPU
func (b *CommonBuilder) InterChipletPorts(c *Chiplet) []akita.Port {
	ports := []akita.Port{
		c.chipRdmaEngine.RequestPort,
		c.chipRdmaEngine.ResponsePort,
		c.remoteTranslationUnit.GetRequestPort(),
		c.remoteTranslationUnit.GetResponsePort(),
	}
	return ports
}
