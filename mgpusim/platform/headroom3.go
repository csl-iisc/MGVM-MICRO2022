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
	"gitlab.com/akita/mgpusim/builders/headroom3builder"
	"gitlab.com/akita/mgpusim/driver"
	"gitlab.com/akita/noc/networking/pcie"
	"gitlab.com/akita/util/tracing"
)

// Headroom3PlatformBuilder can build a platform that equips Headroom3 GPU.
type Headroom3PlatformBuilder struct {
	useParallelEngine        bool
	debugISA                 bool
	traceVis                 bool
	traceTLB                 bool
	traceMem                 bool
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
	l2TlbStriping            uint64
}

// MakeHeadroom3Builder creates a EmuBuilder with default parameters.
func MakeHeadroom3Builder() Headroom3PlatformBuilder {
	b := Headroom3PlatformBuilder{
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
func (b Headroom3PlatformBuilder) WithParallelEngine() Headroom3PlatformBuilder {
	b.useParallelEngine = true
	return b
}

// WithISADebugging enables ISA debugging in the simulation.
func (b Headroom3PlatformBuilder) WithISADebugging() Headroom3PlatformBuilder {
	b.debugISA = true
	return b
}

// WithVisTracing lets the platform to record traces for visualization purposes.
func (b Headroom3PlatformBuilder) WithVisTracing() Headroom3PlatformBuilder {
	b.traceVis = true
	return b
}

// WithTLBTracing lets the platform to trace memory operations.
func (b Headroom3PlatformBuilder) WithTLBTracing() Headroom3PlatformBuilder {
	b.traceTLB = true
	return b
}

// WithMemTracing lets the platform to trace memory operations.
func (b Headroom3PlatformBuilder) WithMemTracing() Headroom3PlatformBuilder {
	b.traceMem = true
	return b
}

// WithNumGPU sets the number of GPUs to build.
func (b Headroom3PlatformBuilder) WithNumGPU(n int) Headroom3PlatformBuilder {
	b.numGPU = n
	return b
}

// WithoutProgressBar disables the progress bar for kernel execution
func (b Headroom3PlatformBuilder) WithoutProgressBar() Headroom3PlatformBuilder {
	b.disableProgressBar = true
	return b
}

// WithLog2PageSize sets the page size as a power of 2.
func (b Headroom3PlatformBuilder) WithLog2PageSize(
	n uint64,
) Headroom3PlatformBuilder {
	b.log2PageSize = n
	return b
}

func (b Headroom3PlatformBuilder) WithAlg(alg string) Headroom3PlatformBuilder {
	b.alg = alg
	return b
}

// WithLog2PageSize sets the page size as a power of 2.
func (b Headroom3PlatformBuilder) WithMemAllocatorType(memAllocatorType string) Headroom3PlatformBuilder {
	b.memAllocatorType = memAllocatorType
	return b
}

// Build builds a platform with DisTLBGPU GPUs.
func (b Headroom3PlatformBuilder) WithL2TLBStriping(striping uint64) Headroom3PlatformBuilder {
	b.l2TlbStriping = striping
	return b
}

// // WithNumChiplets sets the number of chiplets in the mcm GPU.
// func (b Headroom3PlatformBuilder) WithNumChiplets(n uint64) Headroom3PlatformBuilder {
// 	b.numChiplets = n
// 	return b
// }

// Build builds a platform with Headroom3 GPUs.
func (b Headroom3PlatformBuilder) Build() (akita.Engine, *driver.Driver) {
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

func (b Headroom3PlatformBuilder) createGPUs(
	rootComplexID int,
	pcieConnector *pcie.Connector,
	gpuBuilder headroom3builder.Headroom3Builder,
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

func (b Headroom3PlatformBuilder) createPMCPageTable() *cache.BankedLowModuleFinder {
	pmcAddressTable := new(cache.BankedLowModuleFinder)
	pmcAddressTable.BankSize = 4 * mem.GB
	pmcAddressTable.LowModules = append(pmcAddressTable.LowModules, nil)
	return pmcAddressTable
}

func (b Headroom3PlatformBuilder) createRDMAAddrTable() *cache.BankedLowModuleFinder {
	rdmaAddressTable := new(cache.BankedLowModuleFinder)
	rdmaAddressTable.BankSize = 4 * mem.GB
	rdmaAddressTable.LowModules = append(rdmaAddressTable.LowModules, nil)
	return rdmaAddressTable
}

func (b Headroom3PlatformBuilder) createConnection(
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

func (b Headroom3PlatformBuilder) createEngine() akita.Engine {
	var engine akita.Engine

	if b.useParallelEngine {
		engine = akita.NewParallelEngine()
	} else {
		engine = akita.NewSerialEngine()
	}
	// engine.AcceptHook(akita.NewEventLogger(log.New(os.Stdout, "", 0)))

	return engine
}

// func (b Headroom3PlatformBuilder) createMMU(
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

func (b *Headroom3PlatformBuilder) createGPUBuilder(
	engine akita.Engine,
	gpuDriver *driver.Driver,
) headroom3builder.Headroom3Builder {
	gpuBuilder := headroom3builder.MakeHeadroom3Builder().
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
	gpuBuilder = b.setTLBTracer(gpuBuilder)
	gpuBuilder = b.setMemTracer(gpuBuilder)
	gpuBuilder = b.setISADebugger(gpuBuilder)

	gpuBuilder = gpuBuilder.WithRemoteTLB(b.l2TlbStriping)

	if b.disableProgressBar {
		gpuBuilder = gpuBuilder.WithoutProgressBar()
	}

	return gpuBuilder
}

func (b *Headroom3PlatformBuilder) setISADebugger(
	gpuBuilder headroom3builder.Headroom3Builder,
) headroom3builder.Headroom3Builder {
	if !b.debugISA {
		return gpuBuilder
	}

	gpuBuilder = gpuBuilder.WithISADebugging()
	return gpuBuilder
}

func (b *Headroom3PlatformBuilder) setTLBTracer(
	gpuBuilder headroom3builder.Headroom3Builder,
) headroom3builder.Headroom3Builder {
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

func (b *Headroom3PlatformBuilder) setMemTracer(
	gpuBuilder headroom3builder.Headroom3Builder,
) headroom3builder.Headroom3Builder {
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

func (b *Headroom3PlatformBuilder) setVisTracer(
	gpuDriver *driver.Driver,
	gpuBuilder headroom3builder.Headroom3Builder,
) headroom3builder.Headroom3Builder {
	if !b.traceVis {
		return gpuBuilder
	}

	tracer := tracing.NewMySQLTracer()
	tracer.Init()
	tracing.CollectTrace(gpuDriver, tracer)

	gpuBuilder = gpuBuilder.WithVisTracer(tracer)
	return gpuBuilder
}

func (b *Headroom3PlatformBuilder) createGPU(
	index int,
	gpuBuilder headroom3builder.Headroom3Builder,
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

// func (b *Headroom3PlatformBuilder) configRDMAEngine(
// 	gpu *mgpusim.GPU,
// 	addrTable *cache.BankedLowModuleFinder,
// ) {
// 	gpu.RDMAEngine.RemoteRDMAAddressTable = addrTable
// 	// addrTable.LowModules = append(
// 	// 	addrTable.LowModules,
// 	// 	gpu.RDMAEngine.RequestPort)
// }

func (b *Headroom3PlatformBuilder) configPMC(
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
