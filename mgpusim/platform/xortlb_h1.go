package platform

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem"
	"gitlab.com/akita/mgpusim/builders"
	"gitlab.com/akita/mgpusim/driver"
)

// XORTLBH1GPUPlatformBuilder can build a platform that equips DisTLBGPU GPU.
type XORTLBH1GPUPlatformBuilder struct {
	CommonPlatformBuilder
	l2TlbStriping uint64
}

// Makebuilder creates a EmuBuilder with default parameters.
func MakeXORTLBH1GPUPlatformBuilder() XORTLBH1GPUPlatformBuilder {
	b := XORTLBH1GPUPlatformBuilder{
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
	}
	return b
}

// Build builds a platform with DisTLBGPU GPUs.
func (b XORTLBH1GPUPlatformBuilder) WithL2TLBStriping(striping uint64) XORTLBH1GPUPlatformBuilder {
	b.l2TlbStriping = striping
	return b
}

// Build builds a platform with DisTLBGPU GPUs.
func (b XORTLBH1GPUPlatformBuilder) Build() (akita.Engine, *driver.Driver) {
	engine := b.createEngine()

	gpuDriver := driver.NewDriver(engine, b.log2PageSize, b.memAllocatorType)
	// fmt.Println(b.l2TlbStriping, "boo")
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

func (b *XORTLBH1GPUPlatformBuilder) createGPUBuilder(
	engine akita.Engine,
	gpuDriver *driver.Driver,
) builders.Builder {
	gpuBuilder := builders.MakeXORTLBH1GPUBuilder()
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
	gpuBuilder.UseCoalescingTLBPort(b.useCoalescingTLBPort)
	gpuBuilder.UseCoalescingRTU(b.useCoalescingRTU)

	b.setVisTracer(gpuDriver, gpuBuilder)
	b.setTLBTracer(gpuBuilder)
	b.setMemTracer(gpuBuilder)
	b.setISADebugger(gpuBuilder)
	// b.setUseCoalescingTLBPort(gpuBuilder)
	// fmt.Println(b.l2TlbStriping)
	gpuBuilder.WithRemoteTLB(b.l2TlbStriping)

	if b.disableProgressBar {
		gpuBuilder.WithoutProgressBar()
	}

	return gpuBuilder
}
