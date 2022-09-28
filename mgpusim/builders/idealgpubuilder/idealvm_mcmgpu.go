package idealgpubuilder

import (
	"fmt"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem"
	"gitlab.com/akita/mem/cache"
	"gitlab.com/akita/mem/cache/writeback"
	"gitlab.com/akita/mem/device"

	"gitlab.com/akita/mem/idealmemcontroller"
	"gitlab.com/akita/mem/vm/addresstranslator"
	"gitlab.com/akita/mem/vm/idealtlb"
	"gitlab.com/akita/mem/vm/mmu"
	"gitlab.com/akita/mgpusim"
	"gitlab.com/akita/mgpusim/pagemigrationcontroller"
	"gitlab.com/akita/mgpusim/rdma"
	"gitlab.com/akita/mgpusim/timing/caches/l1v"
	"gitlab.com/akita/mgpusim/timing/caches/rob"
	"gitlab.com/akita/mgpusim/timing/cp"
	"gitlab.com/akita/mgpusim/timing/cu"
	"gitlab.com/akita/noc/networking/chipnetwork"
	"gitlab.com/akita/util/tracing"
)

// IdealVMGPUBuilder can build MCM GPUs.
type IdealVMGPUBuilder struct {
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
	l1vTLBs                 []*idealtlb.IdealTLB
	l1sTLBs                 []*idealtlb.IdealTLB
	l1iTLBs                 []*idealtlb.IdealTLB
	l2TLBs                  []*idealtlb.IdealTLB
	drams                   []*idealmemcontroller.Comp
	lowModuleFinderForL1    *cache.InterleavedLowModuleFinder
	lowModuleFinderForL2    *cache.InterleavedLowModuleFinder
	lowModuleFinderForPMC   *cache.InterleavedLowModuleFinder
	dmaEngine               *cp.DMAEngine
	rdmaEngine              *rdma.Engine
	pageMigrationController *pagemigrationcontroller.PageMigrationController

	alg string

	internalConn *akita.DirectConnection

	interChipletNetwork *chipnetwork.Connector
}

// MakeIdealVMGPUBuilder provides a GPU builder that can builds the MCM GPU.
func MakeIdealVMGPUBuilder() IdealVMGPUBuilder {
	b := IdealVMGPUBuilder{
		freq:                           1 * akita.GHz,
		numShaderArrayPerChiplet:       16,
		numCUPerShaderArray:            4,
		numMemoryBankPerChiplet:        8,
		log2CacheLineSize:              6,
		log2PageSize:                   12,
		log2MemoryBankInterleavingSize: 12,
	}
	return b
}

// WithEngine sets the engine that the GPU use.
func (b IdealVMGPUBuilder) WithEngine(engine akita.Engine) IdealVMGPUBuilder {
	b.engine = engine
	return b
}

// WithFreq sets the frequency that the GPU works at.
func (b IdealVMGPUBuilder) WithFreq(freq akita.Freq) IdealVMGPUBuilder {
	b.freq = freq
	return b
}

// WithMemAddrOffset sets the address of the first byte of the GPU to build.
func (b IdealVMGPUBuilder) WithMemAddrOffset(
	offset uint64,
) IdealVMGPUBuilder {
	b.memAddrOffset = offset
	return b
}

// WithTotalMem sets the total amount of memory in the GPU.
func (b IdealVMGPUBuilder) WithTotalMem(memory uint64) IdealVMGPUBuilder {
	b.totalMem = memory
	return b
}

// WithMMU sets the MMU component that provides the address translation service
// for the GPU.
func (b IdealVMGPUBuilder) WithMMU(mmu *mmu.MMUImpl) IdealVMGPUBuilder {
	b.mmu = mmu
	return b
}

// WithNumMemoryBankPerChiplet sets the number of L2 cache modules and number
// of memory controllers in each chiplet.
func (b IdealVMGPUBuilder) WithNumMemoryBankPerChiplet(n int) IdealVMGPUBuilder {
	b.numMemoryBankPerChiplet = n
	return b
}

// WithNumChiplet sets the number of chiplets
func (b IdealVMGPUBuilder) WithNumChiplet(n int) IdealVMGPUBuilder {
	b.numChiplet = n
	return b
}

// WithNumShaderArrayPerChiplet sets the number of shader arrays in each chiplet
func (b IdealVMGPUBuilder) WithNumShaderArrayPerChiplet(n int) IdealVMGPUBuilder {
	b.numShaderArrayPerChiplet = n
	return b
}

// WithNumCUPerShaderArray sets the number of CU and number of L1V caches in
// each Shader Array.
func (b IdealVMGPUBuilder) WithNumCUPerShaderArray(n int) IdealVMGPUBuilder {
	b.numCUPerShaderArray = n
	return b
}

// WithLog2MemoryBankInterleavingSize sets the number of consecutive bytes that
// are guaranteed to be on a memory bank.
func (b IdealVMGPUBuilder) WithLog2MemoryBankInterleavingSize(
	n uint64,
) IdealVMGPUBuilder {
	b.log2MemoryBankInterleavingSize = n
	return b
}

// WithVisTracer applies a tracer to trace all the tasks of all the GPU
// components
func (b IdealVMGPUBuilder) WithVisTracer(t tracing.Tracer) IdealVMGPUBuilder {
	b.enableVisTracing = true
	b.visTracer = t
	return b
}

// WithMemTracer applies a tracer to trace the memory transactions.
func (b IdealVMGPUBuilder) WithTLBTracer(t tracing.Tracer) IdealVMGPUBuilder {
	b.enableTLBTracing = true
	b.tlbTracer = t
	return b
}

// WithMemTracer applies a tracer to trace the memory transactions.
func (b IdealVMGPUBuilder) WithMemTracer(t tracing.Tracer) IdealVMGPUBuilder {
	b.enableMemTracing = true
	b.memTracer = t
	return b
}

// WithISADebugging enables the GPU to dump instruction execution information.
func (b IdealVMGPUBuilder) WithISADebugging() IdealVMGPUBuilder {
	b.enableISADebugging = true
	return b
}

// WithoutProgressBar disables the progress bar for kernel execution
func (b IdealVMGPUBuilder) WithoutProgressBar() IdealVMGPUBuilder {
	b.disableProgressBar = true
	return b
}

// WithLog2CacheLineSize sets the cache line size with the power of 2.
func (b IdealVMGPUBuilder) WithLog2CacheLineSize(
	log2CacheLine uint64,
) IdealVMGPUBuilder {
	b.log2CacheLineSize = log2CacheLine
	return b
}

// WithLog2PageSize sets the page size with the power of 2.
func (b IdealVMGPUBuilder) WithLog2PageSize(log2PageSize uint64) IdealVMGPUBuilder {
	b.log2PageSize = log2PageSize
	return b
}

// WithPageTable sets the page size with the power of 2.
func (b IdealVMGPUBuilder) WithPageTable(pt *device.PageTableImpl) IdealVMGPUBuilder {
	b.pageTable = pt
	return b
}

// WithoutProgressBar disables the progress bar for kernel execution
func (b IdealVMGPUBuilder) WithAlg(alg string) IdealVMGPUBuilder {
	b.alg = alg
	return b
}

// CalculateMemoryParameters calculates
// -> memoryPerChiplet
// -> memoryPerBank
// based on the totalMem, numChiplets and numMemoryBankPerChiplet
func (b IdealVMGPUBuilder) CalculateMemoryParameters() IdealVMGPUBuilder {
	b.memoryPerChiplet = b.totalMem / uint64(b.numChiplet)
	b.memoryPerBank =
		b.totalMem / uint64(b.numChiplet*b.numMemoryBankPerChiplet)

	return b
}

// Build creates an MCM GPU.
func (b IdealVMGPUBuilder) Build(name string, id uint64) *mgpusim.GPU {
	b.createGPU(name, id)

	b.buildCP()

	chipRdmaAddressTable := b.createChipRDMAAddrTable()
	rdmaResponsePorts := make([]akita.Port, 4)

	for i := 0; i < b.numChiplet; i++ {
		chipletName := fmt.Sprintf("%s.chiplet_%02d", b.gpuName, i)
		chiplet := NewChiplet(chipletName, uint64(i))

		b.BuildSAs(chiplet)
		b.buildMemBanks(chiplet)

		b.configChipRDMAEngine(chiplet, chipRdmaAddressTable, rdmaResponsePorts)

		b.connectL1ToL2(chiplet)
		b.connectL2ToDRAM(chiplet)

		b.chiplets = append(b.chiplets, chiplet)
	}

	b.buildPageMigrationController()
	b.setupDMA()

	b.connectCP()
	b.setupInterchipNetwork()

	return b.gpu
}

func (b *IdealVMGPUBuilder) createGPU(name string, id uint64) {
	b.gpuName = name

	b.gpu = mgpusim.NewGPU(b.gpuName)

	b.gpu.GPUID = id
}

func (b *IdealVMGPUBuilder) buildCP() {
	builder := cp.MakeBuilder().
		WithEngine(b.engine).
		WithFreq(b.freq).
		WithAlg(b.alg)

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

func (b *IdealVMGPUBuilder) buildDMAEngine() {
	b.dmaEngine = cp.NewDMAEngine(
		fmt.Sprintf("%s.DMA", b.gpuName),
		b.engine,
		nil)

	if b.enableVisTracing {
		tracing.CollectTrace(b.dmaEngine, b.visTracer)
	}
}

// BuildSAs builds shader arrays.
func (b *IdealVMGPUBuilder) BuildSAs(chiplet *Chiplet) {
	saBuilder := makeShaderArrayBuilder().
		withEngine(b.engine).
		withFreq(b.freq).
		withGPUID(b.gpu.GPUID).
		withLog2CachelineSize(b.log2CacheLineSize).
		withLog2PageSize(b.log2PageSize).
		withNumCU(b.numCUPerShaderArray).
		withPageTable(b.pageTable)

	if b.enableVisTracing {
		saBuilder = saBuilder.withVisTracer(b.visTracer)
	}

	for i := 0; i < b.numShaderArrayPerChiplet; i++ {
		saName := fmt.Sprintf("%s.SA_%02d", chiplet.name, i)
		b.buildSA(saBuilder, saName, chiplet)
	}
}

func (b *IdealVMGPUBuilder) buildSA(
	saBuilder shaderArrayBuilder,
	saName string,
	chiplet *Chiplet,
) {
	sa := saBuilder.Build(saName)

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
		// b.gpu.L1VTLBs = append(b.gpu.L1VTLBs, tlb)
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
	// b.gpu.L1STLBs = append(b.gpu.L1STLBs, sa.l1sTLB)
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
	// b.gpu.L1ITLBs = append(b.gpu.L1ITLBs, sa.l1iTLB)
	chiplet.L1ITLBs = append(chiplet.L1ITLBs, sa.l1iTLB)
}

func (b *IdealVMGPUBuilder) buildMemBanks(chiplet *Chiplet) {
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
		// dram := memCtrlBuilder.WithInterleavingAddrConversion(
		// 	1<<b.log2MemoryBankInterleavingSize,
		// 	b.numMemoryBankPerChiplet,
		// 	i, chiplet.ChipletID*b.memoryPerChiplet+b.memAddrOffset,
		// 	chiplet.ChipletID*b.memoryPerChiplet+b.memAddrOffset+4*mem.GB).
		// 	Build(dramName)
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
		l2.SetLowModuleFinder(&cache.SingleLowModuleFinder{
			LowModule: dram.ToTop,
		})
		if b.enableVisTracing {
			tracing.CollectTrace(l2, b.visTracer)
		}
	}
}

// Though the name of function is configChipRDMAEngine, it actually creates
// the Chip RDMA engine too..
func (b *IdealVMGPUBuilder) configChipRDMAEngine(
	chiplet *Chiplet,
	addrTable *cache.BankedLowModuleFinder, rdmaResponsePorts []akita.Port) {

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

func (b *IdealVMGPUBuilder) connectCP() {
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
}

func (b *IdealVMGPUBuilder) connectCPWithCUs() {
	for _, chiplet := range b.chiplets {
		for _, cu := range chiplet.CUs {
			b.cp.RegisterCU(cu)
			b.internalConn.PlugIn(cu.ToACE, 1)
			b.internalConn.PlugIn(cu.ToCP, 1)
		}
	}
}

func (b *IdealVMGPUBuilder) connectCPWithAddressTranslators() {
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

func (b *IdealVMGPUBuilder) connectCPWithTLBs() {
	for _, chiplet := range b.chiplets {
		// for _, tlb := range chiplet.L2TLBs {
		// 	b.cp.TLBs = append(b.cp.TLBs, tlb.ControlPort)
		// 	b.internalConn.PlugIn(tlb.ControlPort, 1)
		// }

		for _, tlb := range chiplet.L1VTLBs {
			b.cp.TLBs = append(b.cp.TLBs, tlb.ControlPort)
			b.internalConn.PlugIn(tlb.ControlPort, 1)
		}

		for _, tlb := range chiplet.L1STLBs {
			b.cp.TLBs = append(b.cp.TLBs, tlb.ControlPort)
			b.internalConn.PlugIn(tlb.ControlPort, 1)
		}

		for _, tlb := range chiplet.L1ITLBs {
			b.cp.TLBs = append(b.cp.TLBs, tlb.ControlPort)
			b.internalConn.PlugIn(tlb.ControlPort, 1)
		}
	}
}

func (b *IdealVMGPUBuilder) connectCPWithCaches() {
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

func (b *IdealVMGPUBuilder) connectL1ToL2(chiplet *Chiplet) {
	lowModuleFinder := cache.NewInterleavedLowModuleFinder(
		1 << b.log2MemoryBankInterleavingSize)
	lowModuleFinder.ModuleForOtherAddresses = chiplet.chipRdmaEngine.ToL1
	lowModuleFinder.UseAddressSpaceLimitation = true
	lowModuleFinder.LowAddress = b.memAddrOffset +
		b.memoryPerChiplet*chiplet.ChipletID
	lowModuleFinder.HighAddress = lowModuleFinder.LowAddress + b.memoryPerChiplet

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

	return
	// chiplet.MMU.SetLowModuleFinder(lowModuleFinder)
	// l1ToL2Conn.PlugIn(chiplet.MMU.TranslationPort, 64)
}

func (b *IdealVMGPUBuilder) connectL2ToDRAM(chiplet *Chiplet) {
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

func (b *IdealVMGPUBuilder) buildPageMigrationController() {
	b.pageMigrationController =
		pagemigrationcontroller.NewPageMigrationController(
			fmt.Sprintf("%s.PMC", b.gpuName),
			b.engine,
			b.lowModuleFinderForPMC,
			nil)

	b.gpu.PMC = b.pageMigrationController
}

func (b *IdealVMGPUBuilder) setupDMA() {
	lowModuleFinder := cache.NewTwoLevelLowModuleFinder(
		b.memAddrOffset, b.memAddrOffset+b.totalMem,
		b.memoryPerChiplet, uint64(b.numChiplet),
		4096, uint64(b.numMemoryBankPerChiplet))

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

func (b *IdealVMGPUBuilder) setupInterchipNetwork() {
	chipConnector := chipnetwork.NewInterChipletConnector().
		WithEngine(b.engine).
		WithSwitchLatency(360).
		WithFreq(1 * akita.GHz).
		WithFlitByteSize(64).
		WithNumReqPerCycle(12).
		WithNetworkName("ICN")
	chipConnector.CreateNetwork()
	for _, chiplet := range b.chiplets {
		chipConnector.PlugInChip(chiplet.InterChipletPorts())
	}
	chipConnector.MakeNetwork()
}

func (b *IdealVMGPUBuilder) buildRDMAEngine() {
	b.rdmaEngine = rdma.NewEngine(
		fmt.Sprintf("%s.RDMA", b.gpuName),
		b.engine,
		b.lowModuleFinderForL1,
		nil,
	)
	b.gpu.RDMAEngine = b.rdmaEngine
}

func (b *IdealVMGPUBuilder) createChipRDMAAddrTable() *cache.BankedLowModuleFinder {
	chipRdmaAddressTable := new(cache.BankedLowModuleFinder)
	chipRdmaAddressTable.BankSize = b.memoryPerChiplet
	chipRdmaAddressTable.LowModules = append(chipRdmaAddressTable.LowModules, nil)
	return chipRdmaAddressTable
}

func (b *IdealVMGPUBuilder) createRemoteAddrTransTable() *cache.InterleavedLowModuleFinder {
	return nil
}
