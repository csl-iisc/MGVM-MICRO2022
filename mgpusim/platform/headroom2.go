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
	"gitlab.com/akita/mgpusim/builders/headroom2builder"
	"gitlab.com/akita/mgpusim/driver"
	"gitlab.com/akita/noc/networking/pcie"
	"gitlab.com/akita/util/tracing"
)

// Headroom2PlatformBuilder can build a platform that equips Headroom2 GPU.
type Headroom2PlatformBuilder struct {
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
	memAllocatorType         string
	l2TlbStriping            uint64
}

// MakeHeadroom2Builder creates a EmuBuilder with default parameters.
func MakeHeadroom2Builder() Headroom2PlatformBuilder {
	b := Headroom2PlatformBuilder{
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
func (b Headroom2PlatformBuilder) WithParallelEngine() Headroom2PlatformBuilder {
	b.useParallelEngine = true
	return b
}

// WithISADebugging enables ISA debugging in the simulation.
func (b Headroom2PlatformBuilder) WithISADebugging() Headroom2PlatformBuilder {
	b.debugISA = true
	return b
}

// WithVisTracing lets the platform to record traces for visualization purposes.
func (b Headroom2PlatformBuilder) WithVisTracing() Headroom2PlatformBuilder {
	b.traceVis = true
	return b
}

// WithMemTracing lets the platform to trace memory operations.
func (b Headroom2PlatformBuilder) WithMemTracing() Headroom2PlatformBuilder {
	b.traceMem = true
	return b
}

// WithMemTracing lets the platform to trace memory operations.
func (b Headroom2PlatformBuilder) WithTLBTracing() Headroom2PlatformBuilder {
	b.traceTLB = true
	return b
}

// WithNumGPU sets the number of GPUs to build.
func (b Headroom2PlatformBuilder) WithNumGPU(n int) Headroom2PlatformBuilder {
	b.numGPU = n
	return b
}

// WithoutProgressBar disables the progress bar for kernel execution
func (b Headroom2PlatformBuilder) WithoutProgressBar() Headroom2PlatformBuilder {
	b.disableProgressBar = true
	return b
}

// WithLog2PageSize sets the page size as a power of 2.
func (b Headroom2PlatformBuilder) WithLog2PageSize(
	n uint64,
) Headroom2PlatformBuilder {
	b.log2PageSize = n
	return b
}

// WithLog2PageSize sets the page size as a power of 2.
func (b Headroom2PlatformBuilder) WithMemAllocatorType(memAllocatorType string) Headroom2PlatformBuilder {
	b.memAllocatorType = memAllocatorType
	return b
}

// Build builds a platform with DisTLBGPU GPUs.
func (b Headroom2PlatformBuilder) WithL2TLBStriping(striping uint64) Headroom2PlatformBuilder {
	b.l2TlbStriping = striping
	return b
}

// // WithNumChiplets sets the number of chiplets in the mcm GPU.
// func (b Headroom2PlatformBuilder) WithNumChiplets(n uint64) Headroom2PlatformBuilder {
// 	b.numChiplets = n
// 	return b
// }

// Build builds a platform with Headroom2 GPUs.
func (b Headroom2PlatformBuilder) Build() (akita.Engine, *driver.Driver) {
	engine := b.createEngine()

	gpuDriver := driver.NewDriver(engine, b.log2PageSize, b.memAllocatorType)
	headroom2builder := b.createheadroom2builder(engine, gpuDriver)
	pcieConnector, rootComplexID :=
		b.createConnection(engine, gpuDriver) //, mmuComponent)

	// mmuComponent.MigrationServiceProvider = gpuDriver.ToMMU

	rdmaAddressTable := b.createRDMAAddrTable()

	pmcAddressTable := b.createPMCPageTable()

	b.createGPUs(
		rootComplexID, pcieConnector,
		headroom2builder, gpuDriver,
		rdmaAddressTable, pmcAddressTable)

	return engine, gpuDriver
}

func (b Headroom2PlatformBuilder) createGPUs(
	rootComplexID int,
	pcieConnector *pcie.Connector,
	headroom2builder headroom2builder.Headroom2Builder,
	gpuDriver *driver.Driver,
	rdmaAddressTable *cache.BankedLowModuleFinder,
	pmcAddressTable *cache.BankedLowModuleFinder,
) {
	lastSwitchID := rootComplexID
	for i := 1; i < b.numGPU+1; i++ {
		if i%2 == 1 {
			lastSwitchID = pcieConnector.AddSwitch(lastSwitchID)
		}

		b.createGPU(i, headroom2builder, gpuDriver,
			rdmaAddressTable, pmcAddressTable,
			pcieConnector, lastSwitchID)
	}
}

func (b Headroom2PlatformBuilder) createPMCPageTable() *cache.BankedLowModuleFinder {
	pmcAddressTable := new(cache.BankedLowModuleFinder)
	pmcAddressTable.BankSize = 4 * mem.GB
	pmcAddressTable.LowModules = append(pmcAddressTable.LowModules, nil)
	return pmcAddressTable
}

func (b Headroom2PlatformBuilder) createRDMAAddrTable() *cache.BankedLowModuleFinder {
	rdmaAddressTable := new(cache.BankedLowModuleFinder)
	rdmaAddressTable.BankSize = 4 * mem.GB
	rdmaAddressTable.LowModules = append(rdmaAddressTable.LowModules, nil)
	return rdmaAddressTable
}

func (b Headroom2PlatformBuilder) createConnection(
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

func (b Headroom2PlatformBuilder) createEngine() akita.Engine {
	var engine akita.Engine

	if b.useParallelEngine {
		engine = akita.NewParallelEngine()
	} else {
		engine = akita.NewSerialEngine()
	}
	// engine.AcceptHook(akita.NewEventLogger(log.New(os.Stdout, "", 0)))

	return engine
}

// func (b Headroom2PlatformBuilder) createMMU(
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

func (b *Headroom2PlatformBuilder) createheadroom2builder(
	engine akita.Engine,
	gpuDriver *driver.Driver,
) headroom2builder.Headroom2Builder {
	headroom2builder := headroom2builder.MakeHeadroom2Builder().
		WithEngine(engine).
		WithNumCUPerShaderArray(int(b.numCUPerShaderArray)).
		WithNumShaderArrayPerChiplet(int(b.numShaderArrayPerChiplet)).
		WithNumMemoryBankPerChiplet(int(b.numMemoryBankPerChiplet)).
		WithNumChiplet(int(b.numChiplets)).
		WithTotalMem(b.totalMem).
		CalculateMemoryParameters().
		WithLog2PageSize(b.log2PageSize).
		WithPageTable(gpuDriver.PageTable)

	headroom2builder = b.setVisTracer(gpuDriver, headroom2builder)
	headroom2builder = b.setMemTracer(headroom2builder)
	headroom2builder = b.setTLBTracer(headroom2builder)
	headroom2builder = b.setISADebugger(headroom2builder)

	if b.disableProgressBar {
		headroom2builder = headroom2builder.WithoutProgressBar()
	}

	return headroom2builder
}

func (b *Headroom2PlatformBuilder) setISADebugger(
	headroom2builder headroom2builder.Headroom2Builder,
) headroom2builder.Headroom2Builder {
	if !b.debugISA {
		return headroom2builder
	}

	headroom2builder = headroom2builder.WithISADebugging()
	return headroom2builder
}

func (b *Headroom2PlatformBuilder) setMemTracer(
	headroom2builder headroom2builder.Headroom2Builder,
) headroom2builder.Headroom2Builder {
	if !b.traceMem {
		return headroom2builder
	}

	file, err := os.Create("mem.trace")
	if err != nil {
		panic(err)
	}
	logger := log.New(file, "", 0)
	memTracer := memtraces.NewTracer(logger)
	headroom2builder = headroom2builder.WithMemTracer(memTracer)
	return headroom2builder
}

func (b *Headroom2PlatformBuilder) setTLBTracer(
	headroom2builder headroom2builder.Headroom2Builder,
) headroom2builder.Headroom2Builder {
	if !b.traceTLB {
		return headroom2builder
	}

	file, err := os.Create("tlb.trace")
	if err != nil {
		panic(err)
	}
	logger := log.New(file, "", 0)

	tlbTracer := memtraces.NewTracer(logger)
	headroom2builder = headroom2builder.WithTLBTracer(tlbTracer)

	return headroom2builder
}

func (b *Headroom2PlatformBuilder) setVisTracer(
	gpuDriver *driver.Driver,
	headroom2builder headroom2builder.Headroom2Builder,
) headroom2builder.Headroom2Builder {
	if !b.traceVis {
		return headroom2builder
	}

	tracer := tracing.NewMySQLTracer()
	tracer.Init()
	tracing.CollectTrace(gpuDriver, tracer)

	headroom2builder = headroom2builder.WithVisTracer(tracer)
	return headroom2builder
}

func (b *Headroom2PlatformBuilder) createGPU(
	index int,
	headroom2builder headroom2builder.Headroom2Builder,
	gpuDriver *driver.Driver,
	rdmaAddressTable *cache.BankedLowModuleFinder,
	pmcAddressTable *cache.BankedLowModuleFinder,
	pcieConnector *pcie.Connector,
	pcieSwitchID int,
) *mgpusim.GPU {
	name := fmt.Sprintf("GPU%d", index)
	memAddrOffset := uint64(index) * 2 * mem.GB
	gpu := headroom2builder.
		WithMemAddrOffset(memAddrOffset).
		Build(name, uint64(index))
	gpuDriver.RegisterGPU(gpu, 8*mem.GB)
	gpu.CommandProcessor.Driver = gpuDriver.ToGPUs

	// b.configRDMAEngine(gpu, rdmaAddressTable)
	b.configPMC(gpu, gpuDriver, pmcAddressTable)

	pcieConnector.PlugInDevice(pcieSwitchID, gpu.ExternalPorts())

	return gpu
}

// func (b *Headroom2PlatformBuilder) configRDMAEngine(
// 	gpu *mgpusim.GPU,
// 	addrTable *cache.BankedLowModuleFinder,
// ) {
// 	gpu.RDMAEngine.RemoteRDMAAddressTable = addrTable
// 	// addrTable.LowModules = append(
// 	// 	addrTable.LowModules,
// 	// 	gpu.RDMAEngine.RequestPort)
// }

func (b *Headroom2PlatformBuilder) configPMC(
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
