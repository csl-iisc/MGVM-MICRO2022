package platform

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem"
	"gitlab.com/akita/mgpusim/builders"
	"gitlab.com/akita/mgpusim/driver"
)

// CustomTLBGPUPlatformBuilder can build a platform that equips DisTLBGPU GPU.
type CustomTLBGPUPlatformBuilder struct {
	CommonPlatformBuilder
	l2TlbStriping       uint64
	switchL2TLBStriping bool
	usePtCaching        bool
}

// Makebuilder creates a EmuBuilder with default parameters.
func MakeCustomTLBGPUPlatformBuilder() CustomTLBGPUPlatformBuilder {
	b := CustomTLBGPUPlatformBuilder{
		CommonPlatformBuilder{
			numGPU:                   1,
			log2PageSize:             uint64(12),
			numCUPerShaderArray:      uint64(4),
			numShaderArrayPerChiplet: uint64(8),
			numMemoryBankPerChiplet:  uint64(16),
			numChiplets:              uint64(4),
			totalMem:                 16 * mem.GB,
			bankSize:                 256 * mem.MB,
			lowAddr:                  4 * mem.GB,
		},
		512,
		false,
		false,
	}
	return b
}

// Build builds a platform with DisTLBGPU GPUs.
func (b CustomTLBGPUPlatformBuilder) WithL2TLBStriping(striping uint64) CustomTLBGPUPlatformBuilder {
	b.l2TlbStriping = striping
	return b
}

// Build builds a platform with DisTLBGPU GPUs.
func (b CustomTLBGPUPlatformBuilder) SwitchL2TLBStriping(switchStriping bool) CustomTLBGPUPlatformBuilder {
	b.switchL2TLBStriping = switchStriping
	return b
}

// Build builds a platform with DisTLBGPU GPUs.
func (b CustomTLBGPUPlatformBuilder) UsePtCaching(ptCaching bool) CustomTLBGPUPlatformBuilder {
	b.usePtCaching = ptCaching
	return b
}

// Build builds a platform with DisTLBGPU GPUs.
func (b CustomTLBGPUPlatformBuilder) Build() (akita.Engine, *driver.Driver) {
	engine := b.createEngine()

	gpuDriver := driver.NewDriver(engine, b.log2PageSize, b.memAllocatorType)
	gpuBuilder := b.createGPUBuilder(engine, gpuDriver)
	pcieConnector, rootComplexID :=
		b.createConnection(engine, gpuDriver)

	rdmaAddressTable := b.createRDMAAddrTable()

	pmcAddressTable := b.createPMCPageTable()

	b.createGPUs(
		rootComplexID, pcieConnector,
		gpuBuilder, gpuDriver,
		rdmaAddressTable, pmcAddressTable)

	return engine, gpuDriver
}

func (b *CustomTLBGPUPlatformBuilder) createGPUBuilder(
	engine akita.Engine,
	gpuDriver *driver.Driver,
) builders.Builder {
	gpuBuilder := builders.MakeCustomTLBGPUBuilder()
	gpuBuilder.WithEngine(engine)
	gpuBuilder.WithNumCUPerShaderArray(int(b.numCUPerShaderArray))
	gpuBuilder.WithNumShaderArrayPerChiplet(int(b.numShaderArrayPerChiplet))
	gpuBuilder.WithNumMemoryBankPerChiplet(int(b.numMemoryBankPerChiplet))
	gpuBuilder.WithNumChiplet(int(b.numChiplets))
	gpuBuilder.WithTotalMem(b.totalMem)
	gpuBuilder.CalculateMemoryParameters()
	gpuBuilder.WithLog2PageSize(b.log2PageSize)
	gpuBuilder.WithPageTable(gpuDriver.PageTable)
	gpuBuilder.WithAlg(b.alg)
	gpuBuilder.WithSchedulingPartition(b.partition)
	gpuBuilder.WithWalkersPerChiplet(b.walkersPerChiplet)
	gpuBuilder.UseCoalescingTLBPort(b.useCoalescingTLBPort)
	gpuBuilder.UseCoalescingRTU(b.useCoalescingRTU)
	gpuBuilder.WithCustomHSL(b.customHSLpmdUnits)

	b.setVisTracer(gpuDriver, gpuBuilder)
	b.setTLBTracer(gpuBuilder)
	b.setMemTracer(gpuBuilder)
	b.setISADebugger(gpuBuilder)
	gpuBuilder.WithRemoteTLB(b.l2TlbStriping)
	gpuBuilder.SwitchL2TLBStriping(b.switchL2TLBStriping)

	if b.disableProgressBar {
		gpuBuilder.WithoutProgressBar()
	}

	return gpuBuilder
}
