package platform

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem"
	"gitlab.com/akita/mgpusim/builders"
	"gitlab.com/akita/mgpusim/driver"
)

// XORTLBGPUPlatformBuilder can build a platform that equips DisTLBGPU GPU.
type XORTLBGPUPlatformBuilder struct {
	CommonPlatformBuilder
	l2TlbStriping       uint64
	switchL2TLBStriping bool
	usePtCaching        bool
}

// Makebuilder creates a EmuBuilder with default parameters.
func MakeXORTLBGPUPlatformBuilder() XORTLBGPUPlatformBuilder {
	b := XORTLBGPUPlatformBuilder{
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
func (b XORTLBGPUPlatformBuilder) WithL2TLBStriping(striping uint64) XORTLBGPUPlatformBuilder {
	b.l2TlbStriping = striping
	return b
}

// Build builds a platform with DisTLBGPU GPUs.
func (b XORTLBGPUPlatformBuilder) SwitchL2TLBStriping(switchStriping bool) XORTLBGPUPlatformBuilder {
	b.switchL2TLBStriping = switchStriping
	return b
}

// Build builds a platform with DisTLBGPU GPUs.
func (b XORTLBGPUPlatformBuilder) UsePtCaching(ptCaching bool) XORTLBGPUPlatformBuilder {
	b.usePtCaching = ptCaching
	return b
}

// Build builds a platform with DisTLBGPU GPUs.
func (b XORTLBGPUPlatformBuilder) Build() (akita.Engine, *driver.Driver) {
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

func (b *XORTLBGPUPlatformBuilder) createGPUBuilder(
	engine akita.Engine,
	gpuDriver *driver.Driver,
) builders.Builder {
	gpuBuilder := builders.MakeXORTLBGPUBuilder()
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

	b.setVisTracer(gpuDriver, gpuBuilder)
	b.setTLBTracer(gpuBuilder)
	b.setMemTracer(gpuBuilder)
	b.setISADebugger(gpuBuilder)
	gpuBuilder.WithRemoteTLB(b.l2TlbStriping)
	gpuBuilder.SwitchL2TLBStriping(b.switchL2TLBStriping)
	gpuBuilder.UsePtCaching(b.usePtCaching)

	if b.disableProgressBar {
		gpuBuilder.WithoutProgressBar()
	}

	return gpuBuilder
}
