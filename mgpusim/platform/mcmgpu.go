package platform

import (
	"fmt"
	"log"
	"os"

	memtraces "gitlab.com/akita/mem/trace"
	"gitlab.com/akita/mgpusim"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem"
	"gitlab.com/akita/mem/cache"
	"gitlab.com/akita/mgpusim/driver"
	"gitlab.com/akita/mgpusim/gpubuilder"
	"gitlab.com/akita/noc/networking/pcie"
	"gitlab.com/akita/util/tracing"
)

// MCMGPUPlatformBuilder can build a platform that equips MCMGPU GPU.
type MCMGPUPlatformBuilder struct {
	useParallelEngine        bool
	debugISA                 bool
	traceVis                 bool
	traceMem                 bool
	traceTLB                 bool
	numGPU                   int
	log2PageSize             uint64
	disableProgressBar       bool
	numChiplets              uint64
	numCUPerShaderArray      uint64
	numShaderArrayPerChiplet uint64
	numMemoryBankPerChiplet  uint64
	totalMem                 uint64
	bankSize                 uint64
	lowAddr                  uint64
	alg                      string
	memAllocatorType         string
}

// MakeMCMGPUBuilder creates a EmuBuilder with default parameters.
func MakeMCMGPUBuilder() MCMGPUPlatformBuilder {
	b := MCMGPUPlatformBuilder{
		numGPU:                   1,
		log2PageSize:             uint64(12),
		numCUPerShaderArray:      uint64(4),
		numShaderArrayPerChiplet: uint64(8),
		numMemoryBankPerChiplet:  uint64(8),
		numChiplets:              uint64(4),
		totalMem:                 8 * mem.GB,
		bankSize:                 256 * mem.MB,
		lowAddr:                  2 * mem.GB,
	}
	return b
}

// WithParallelEngine lets the EmuBuilder to use parallel engine.
func (b MCMGPUPlatformBuilder) WithParallelEngine() MCMGPUPlatformBuilder {
	b.useParallelEngine = true
	return b
}

// WithISADebugging enables ISA debugging in the simulation.
func (b MCMGPUPlatformBuilder) WithISADebugging() MCMGPUPlatformBuilder {
	b.debugISA = true
	return b
}

// WithVisTracing lets the platform to record traces for visualization purposes.
func (b MCMGPUPlatformBuilder) WithVisTracing() MCMGPUPlatformBuilder {
	b.traceVis = true
	return b
}

// WithMemTracing lets the platform to trace memory operations.
func (b MCMGPUPlatformBuilder) WithMemTracing() MCMGPUPlatformBuilder {
	b.traceMem = true
	return b
}

// WithMemTracing lets the platform to trace memory operations.
func (b MCMGPUPlatformBuilder) WithTLBTracing() MCMGPUPlatformBuilder {
	b.traceTLB = true
	return b
}

// WithNumGPU sets the number of GPUs to build.
func (b MCMGPUPlatformBuilder) WithNumGPU(n int) MCMGPUPlatformBuilder {
	b.numGPU = n
	return b
}

// WithoutProgressBar disables the progress bar for kernel execution
func (b MCMGPUPlatformBuilder) WithoutProgressBar() MCMGPUPlatformBuilder {
	b.disableProgressBar = true
	return b
}

// WithAlg sets the scheduling algorithm
func (b MCMGPUPlatformBuilder) WithAlg(alg string) MCMGPUPlatformBuilder {
	b.alg = alg
	return b
}

// WithLog2PageSize sets the page size as a power of 2.
func (b MCMGPUPlatformBuilder) WithLog2PageSize(
	n uint64,
) MCMGPUPlatformBuilder {
	b.log2PageSize = n
	return b
}

// WithLog2PageSize sets the page size as a power of 2.
func (b MCMGPUPlatformBuilder) WithMemAllocatorType(memAllocatorType string) MCMGPUPlatformBuilder {
	b.memAllocatorType = memAllocatorType
	return b
}

// // WithNumChiplets sets the number of chiplets in the mcm GPU.
// func (b MCMGPUPlatformBuilder) WithNumChiplets(n uint64) MCMGPUPlatformBuilder {
// 	b.numChiplets = n
// 	return b
// }

// Build builds a platform with MCMGPU GPUs.
func (b MCMGPUPlatformBuilder) Build() (akita.Engine, *driver.Driver) {
	engine := b.createEngine()

	gpuDriver := driver.NewDriver(engine, b.log2PageSize, b.memAllocatorType)
	gpuBuilder := b.createGPUBuilder(engine, gpuDriver)
	pcieConnector, rootComplexID :=
		b.createConnection(engine, gpuDriver) //, mmuComponent)

	// mmuComponent.MigrationServiceProvider = gpuDriver.ToMMU

	rdmaAddressTable := b.createRDMAAddrTable()

	pmcAddressTable := b.createPMCPageTable()

	b.createGPUs(
		rootComplexID, pcieConnector,
		gpuBuilder, gpuDriver,
		rdmaAddressTable, pmcAddressTable)

	return engine, gpuDriver
}

func (b MCMGPUPlatformBuilder) createGPUs(
	rootComplexID int,
	pcieConnector *pcie.Connector,
	gpuBuilder gpubuilder.MCMGPUBuilder,
	gpuDriver *driver.Driver,
	rdmaAddressTable *cache.BankedLowModuleFinder,
	pmcAddressTable *cache.BankedLowModuleFinder,
) {
	lastSwitchID := rootComplexID
	for i := 1; i < b.numGPU+1; i++ {
		if i%2 == 1 {
			lastSwitchID = pcieConnector.AddSwitch(lastSwitchID)
		}

		b.createGPU(i, gpuBuilder, gpuDriver,
			rdmaAddressTable, pmcAddressTable,
			pcieConnector, lastSwitchID)
	}
}

func (b MCMGPUPlatformBuilder) createPMCPageTable() *cache.BankedLowModuleFinder {
	pmcAddressTable := new(cache.BankedLowModuleFinder)
	pmcAddressTable.BankSize = 4 * mem.GB
	pmcAddressTable.LowModules = append(pmcAddressTable.LowModules, nil)
	return pmcAddressTable
}

func (b MCMGPUPlatformBuilder) createRDMAAddrTable() *cache.BankedLowModuleFinder {
	rdmaAddressTable := new(cache.BankedLowModuleFinder)
	rdmaAddressTable.BankSize = 4 * mem.GB
	rdmaAddressTable.LowModules = append(rdmaAddressTable.LowModules, nil)
	return rdmaAddressTable
}

func (b MCMGPUPlatformBuilder) createConnection(
	engine akita.Engine,
	gpuDriver *driver.Driver,
) (*pcie.Connector, int) {
	//connection := akita.NewDirectConnection(engine)
	// connection := noc.NewFixedBandwidthConnection(32, engine, 1*akita.GHz)
	// connection.SrcBufferCapacity = 40960000
	pcieConnector := pcie.NewConnector().
		WithEngine(engine).
		WithVersion3().
		WithX16().
		WithSwitchLatency(140).
		WithNetworkName("PCIe")
	pcieConnector.CreateNetwork()
	rootComplexID := pcieConnector.CreateRootComplex(
		[]akita.Port{
			gpuDriver.ToGPUs,
			gpuDriver.ToMMU,
			// mmuComponent.MigrationPort,
			// mmuComponent.ToTop,
		})
	return pcieConnector, rootComplexID
}

func (b MCMGPUPlatformBuilder) createEngine() akita.Engine {
	var engine akita.Engine

	if b.useParallelEngine {
		engine = akita.NewParallelEngine()
	} else {
		engine = akita.NewSerialEngine()
	}
	// engine.AcceptHook(akita.NewEventLogger(log.New(os.Stdout, "", 0)))

	return engine
}

// func (b MCMGPUPlatformBuilder) createMMU(
// 	engine akita.Engine,
// ) (*mmu.MMUImpl, device.PageTable) {
// 	pageTable := device.NewPageTable(b.log2PageSize)
// 	mmuBuilder := mmu.MakeBuilder().
// 		WithEngine(engine).
// 		WithFreq(1 * akita.GHz).
// 		WithLog2PageSize(b.log2PageSize).
// 		WithPageTable(pageTable).
// 		WithNumChiplets(b.numChiplets).
// 		WithLowAddr(b.lowAddr).
// 		WithTotMem(b.totalMem).
// 		WithBankSize(b.bankSize).
// 		WithNumMemoryBankPerChiplet(b.numMemoryBankPerChiplet)
// 	mmuComponent := mmuBuilder.Build("MMU")
// 	return mmuComponent, pageTable
// }

func (b *MCMGPUPlatformBuilder) createGPUBuilder(
	engine akita.Engine,
	gpuDriver *driver.Driver,
) gpubuilder.MCMGPUBuilder {
	gpuBuilder := gpubuilder.MakeMCMGPUBuilder().
		WithEngine(engine).
		WithNumCUPerShaderArray(int(b.numCUPerShaderArray)).
		WithNumShaderArrayPerChiplet(int(b.numShaderArrayPerChiplet)).
		WithNumMemoryBankPerChiplet(int(b.numMemoryBankPerChiplet)).
		WithNumChiplet(int(b.numChiplets)).
		WithTotalMem(b.totalMem).
		CalculateMemoryParameters().
		WithLog2PageSize(b.log2PageSize).
		WithPageTable(gpuDriver.PageTable).
		WithAlg(b.alg)

	gpuBuilder = b.setVisTracer(gpuDriver, gpuBuilder)
	gpuBuilder = b.setMemTracer(gpuBuilder)
	gpuBuilder = b.setTLBTracer(gpuBuilder)
	gpuBuilder = b.setISADebugger(gpuBuilder)

	if b.disableProgressBar {
		gpuBuilder = gpuBuilder.WithoutProgressBar()
	}

	return gpuBuilder
}

func (b *MCMGPUPlatformBuilder) setISADebugger(
	gpuBuilder gpubuilder.MCMGPUBuilder,
) gpubuilder.MCMGPUBuilder {
	if !b.debugISA {
		return gpuBuilder
	}

	gpuBuilder = gpuBuilder.WithISADebugging()
	return gpuBuilder
}

func (b *MCMGPUPlatformBuilder) setMemTracer(
	gpuBuilder gpubuilder.MCMGPUBuilder,
) gpubuilder.MCMGPUBuilder {
	if !b.traceMem {
		return gpuBuilder
	}

	file, err := os.Create("mem.trace")
	if err != nil {
		panic(err)
	}
	logger := log.New(file, "", 0)
	memTracer := memtraces.NewTracer(logger)
	gpuBuilder = gpuBuilder.WithMemTracer(memTracer)
	return gpuBuilder
}

func (b *MCMGPUPlatformBuilder) setTLBTracer(
	gpuBuilder gpubuilder.MCMGPUBuilder,
) gpubuilder.MCMGPUBuilder {
	if !b.traceTLB {
		return gpuBuilder
	}

	file, err := os.Create("tlb.trace")
	if err != nil {
		panic(err)
	}
	logger := log.New(file, "", 0)

	tlbTracer := memtraces.NewTracer(logger)
	gpuBuilder = gpuBuilder.WithTLBTracer(tlbTracer)

	return gpuBuilder
}

func (b *MCMGPUPlatformBuilder) setVisTracer(
	gpuDriver *driver.Driver,
	gpuBuilder gpubuilder.MCMGPUBuilder,
) gpubuilder.MCMGPUBuilder {
	if !b.traceVis {
		return gpuBuilder
	}

	tracer := tracing.NewMySQLTracer()
	tracer.Init()
	tracing.CollectTrace(gpuDriver, tracer)

	gpuBuilder = gpuBuilder.WithVisTracer(tracer)
	return gpuBuilder
}

func (b *MCMGPUPlatformBuilder) createGPU(
	index int,
	gpuBuilder gpubuilder.MCMGPUBuilder,
	gpuDriver *driver.Driver,
	rdmaAddressTable *cache.BankedLowModuleFinder,
	pmcAddressTable *cache.BankedLowModuleFinder,
	pcieConnector *pcie.Connector,
	pcieSwitchID int,
) *mgpusim.GPU {
	name := fmt.Sprintf("GPU%d", index)
	memAddrOffset := uint64(index) * 2 * mem.GB
	gpu := gpuBuilder.
		WithMemAddrOffset(memAddrOffset).
		Build(name, uint64(index))
	gpuDriver.RegisterGPU(gpu, 8*mem.GB)
	gpu.CommandProcessor.Driver = gpuDriver.ToGPUs

	// b.configRDMAEngine(gpu, rdmaAddressTable)
	b.configPMC(gpu, gpuDriver, pmcAddressTable)

	pcieConnector.PlugInDevice(pcieSwitchID, gpu.ExternalPorts())

	return gpu
}

// func (b *MCMGPUPlatformBuilder) configRDMAEngine(
// 	gpu *mgpusim.GPU,
// 	addrTable *cache.BankedLowModuleFinder,
// ) {
// 	gpu.RDMAEngine.RemoteRDMAAddressTable = addrTable
// 	// addrTable.LowModules = append(
// 	// 	addrTable.LowModules,
// 	// 	gpu.RDMAEngine.RequestPort)
// }

func (b *MCMGPUPlatformBuilder) configPMC(
	gpu *mgpusim.GPU,
	gpuDriver *driver.Driver,
	addrTable *cache.BankedLowModuleFinder,
) {
	gpu.PMC.RemotePMCAddressTable = addrTable
	addrTable.LowModules = append(
		addrTable.LowModules,
		gpu.PMC.RemotePort)
	gpuDriver.RemotePMCPorts = append(
		gpuDriver.RemotePMCPorts, gpu.PMC.RemotePort)
}
