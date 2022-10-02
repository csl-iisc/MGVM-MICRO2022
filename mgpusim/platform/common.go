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
	"gitlab.com/akita/mgpusim/builders"
	"gitlab.com/akita/mgpusim/driver"
	"gitlab.com/akita/noc/networking/pcie"
	"gitlab.com/akita/util/tracing"
)

// CommonPlatformBuilder can build a platform that equips DisTLBGPU GPU.
type CommonPlatformBuilder struct {
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
	partition                string
	useCoalescingTLBPort     bool
	useCoalescingRTU         bool
	memAllocatorType         string
	customHSLpmdUnits        uint64
	walkersPerChiplet        int
}

// Makebuilder creates a EmuBuilder with default parameters.
func MakeCommonPlatformBuilder() CommonPlatformBuilder {
	b := CommonPlatformBuilder{
		numGPU:                   1,
		log2PageSize:             uint64(12),
		numCUPerShaderArray:      uint64(4),
		numShaderArrayPerChiplet: uint64(8),
		numMemoryBankPerChiplet:  uint64(8),
		numChiplets:              uint64(4),
		totalMem:                 8 * mem.GB,
		bankSize:                 256 * mem.MB,
		lowAddr:                  2 * mem.GB,
		walkersPerChiplet:        16,
	}
	return b
}

// WithParallelEngine lets the EmuBuilder to use parallel engine.
func (b *CommonPlatformBuilder) WithParallelEngine() {
	b.useParallelEngine = true
}

// WithISADebugging enables ISA debugging in the simulation.
func (b *CommonPlatformBuilder) WithISADebugging() {
	b.debugISA = true
}

// WithVisTracing lets the platform to record traces for visualization purposes.
func (b *CommonPlatformBuilder) WithVisTracing() {
	b.traceVis = true
}

// WithTLBTracing lets the platform to trace memory operations.
func (b *CommonPlatformBuilder) WithTLBTracing() {
	b.traceTLB = true
}

// WithMemTracing lets the platform to trace memory operations.
func (b *CommonPlatformBuilder) WithMemTracing() {
	b.traceMem = true
}

// WithNumGPU sets the number of GPUs to build.
func (b *CommonPlatformBuilder) WithNumGPU(n int) {
	b.numGPU = n
}

// WithoutProgressBar disables the progress bar for kernel execution
func (b *CommonPlatformBuilder) WithoutProgressBar() {
	b.disableProgressBar = true
}

// WithLog2PageSize sets the page size as a power of 2.
func (b *CommonPlatformBuilder) WithLog2PageSize(n uint64) {
	b.log2PageSize = n
}

func (b *CommonPlatformBuilder) WithAlg(alg string) {
	b.alg = alg
}

func (b *CommonPlatformBuilder) WithSchedulingPartition(partition string) {
	b.partition = partition
}

func (b *CommonPlatformBuilder) UseCoalescingTLBPort(u bool) {
	b.useCoalescingTLBPort = u
}

func (b *CommonPlatformBuilder) UseCoalescingRTU(u bool) {
	b.useCoalescingRTU = u
}

func (b *CommonPlatformBuilder) WithMemAllocatorType(allocatorType string) {
	b.memAllocatorType = allocatorType
}

func (b *CommonPlatformBuilder) WithCustomHSL(pmdUnits uint64) {
	b.customHSLpmdUnits = pmdUnits
}

func (b *CommonPlatformBuilder) WithWalkersPerChiplet(n int) {
	b.walkersPerChiplet = n
}

// // WithNumChiplets sets the number of chiplets in the mcm GPU.
// func (b CommonPlatformBuilder) WithNumChiplets(n uint64) CommonPlatformBuilder {
// 	b.numChiplets = n
// 	return b
// }

func (b CommonPlatformBuilder) createGPUs(
	rootComplexID int,
	pcieConnector *pcie.Connector,
	gpuBuilder builders.Builder,
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

func (b CommonPlatformBuilder) createPMCPageTable() *cache.BankedLowModuleFinder {
	pmcAddressTable := new(cache.BankedLowModuleFinder)
	pmcAddressTable.BankSize = 4 * mem.GB
	pmcAddressTable.LowModules = append(pmcAddressTable.LowModules, nil)
	return pmcAddressTable
}

func (b CommonPlatformBuilder) createRDMAAddrTable() *cache.BankedLowModuleFinder {
	rdmaAddressTable := new(cache.BankedLowModuleFinder)
	rdmaAddressTable.BankSize = 4 * mem.GB
	rdmaAddressTable.LowModules = append(rdmaAddressTable.LowModules, nil)
	return rdmaAddressTable
}

func (b CommonPlatformBuilder) createConnection(
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

func (b CommonPlatformBuilder) createEngine() akita.Engine {
	var engine akita.Engine

	if b.useParallelEngine {
		engine = akita.NewParallelEngine()
	} else {
		engine = akita.NewSerialEngine()
	}
	// engine.AcceptHook(akita.NewEventLogger(log.New(os.Stdout, "", 0)))

	return engine
}

func (b *CommonPlatformBuilder) setISADebugger(
	gpuBuilder builders.Builder,
) builders.Builder {
	if !b.debugISA {
		return gpuBuilder
	}

	gpuBuilder.WithISADebugging()
	return gpuBuilder
}

func (b *CommonPlatformBuilder) setTLBTracer(
	gpuBuilder builders.Builder,
) builders.Builder {
	if !b.traceTLB {
		return gpuBuilder
	}

	file, err := os.Create("tlb.trace")
	if err != nil {
		panic(err)
	}
	logger := log.New(file, "", 0)
	tlbTracer := memtraces.NewTracer(logger)
	gpuBuilder.WithTLBTracer(tlbTracer)
	return gpuBuilder
}

func (b *CommonPlatformBuilder) setMemTracer(
	gpuBuilder builders.Builder,
) builders.Builder {
	if !b.traceMem {
		return gpuBuilder
	}

	file, err := os.Create("mem.trace")
	if err != nil {
		panic(err)
	}
	logger := log.New(file, "", 0)
	memTracer := memtraces.NewTracer(logger)
	gpuBuilder.WithMemTracer(memTracer)
	return gpuBuilder
}

func (b *CommonPlatformBuilder) setVisTracer(
	gpuDriver *driver.Driver,
	gpuBuilder builders.Builder,
) builders.Builder {
	if !b.traceVis {
		return gpuBuilder
	}

	tracer := tracing.NewMySQLTracer()
	tracer.Init()
	tracing.CollectTrace(gpuDriver, tracer)

	gpuBuilder.WithVisTracer(tracer)
	return gpuBuilder
}

func (b *CommonPlatformBuilder) createGPU(
	index int,
	gpuBuilder builders.Builder,
	gpuDriver *driver.Driver,
	rdmaAddressTable *cache.BankedLowModuleFinder,
	pmcAddressTable *cache.BankedLowModuleFinder,
	pcieConnector *pcie.Connector,
	pcieSwitchID int,
) *mgpusim.GPU {
	name := fmt.Sprintf("GPU%d", index)
	memAddrOffset := uint64(index)*2*mem.GB + uint64(1<<b.log2PageSize)
	gpuBuilder.WithMemAddrOffset(memAddrOffset)
	gpu := gpuBuilder.Build(name, uint64(index))
	gpuDriver.RegisterGPU(gpu, 8*mem.GB)
	gpu.CommandProcessor.Driver = gpuDriver.ToGPUs

	// b.configRDMAEngine(gpu, rdmaAddressTable)
	b.configPMC(gpu, gpuDriver, pmcAddressTable)

	pcieConnector.PlugInDevice(pcieSwitchID, gpu.ExternalPorts())

	return gpu
}

// func (b *CommonPlatformBuilder) configRDMAEngine(
// 	gpu *mgpusim.GPU,
// 	addrTable *cache.BankedLowModuleFinder,
// ) {
// 	gpu.RDMAEngine.RemoteRDMAAddressTable = addrTable
// 	// addrTable.LowModules = append(
// 	// 	addrTable.LowModules,
// 	// 	gpu.RDMAEngine.RequestPort)
// }

func (b *CommonPlatformBuilder) configPMC(
	gpu *mgpusim.GPU,
	gpuDriver *driver.Driver,
	addrTable *cache.BankedLowModuleFinder,
) {
	gpu.PMC.RemotePMCAddressTable = addrTable
	addrTable.LowModules = append(addrTable.LowModules, gpu.PMC.RemotePort)
	gpuDriver.RemotePMCPorts = append(gpuDriver.RemotePMCPorts, gpu.PMC.RemotePort)
}
