package platform

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem"
	"gitlab.com/akita/mgpusim/builders"
	"gitlab.com/akita/mgpusim/driver"
)

// H6GPUPlatformBuilder can build a platform that equips DisTLBGPU GPU.
type H6GPUPlatformBuilder struct {
	CommonPlatformBuilder
	l2TlbStriping uint64
}

// Makebuilder creates a EmuBuilder with default parameters.
func MakeH6GPUPlatformBuilder() H6GPUPlatformBuilder {
	b := H6GPUPlatformBuilder{
		CommonPlatformBuilder{
			numGPU:                   1,
			log2PageSize:             uint64(12),
			numCUPerShaderArray:      uint64(4),
			numShaderArrayPerChiplet: uint64(8),
			numMemoryBankPerChiplet:  uint64(8),
			numChiplets:              uint64(4),
			totalMem:                 8 * mem.GB,
			bankSize:                 256 * mem.MB,
			lowAddr:                  2 * mem.GB,
		},
		512,
	}
	return b
}

// Build builds a platform with DisTLBGPU GPUs.
func (b H6GPUPlatformBuilder) WithL2TLBStriping(striping uint64) H6GPUPlatformBuilder {
	b.l2TlbStriping = striping
	return b
}

// Build builds a platform with DisTLBGPU GPUs.
func (b H6GPUPlatformBuilder) Build() (akita.Engine, *driver.Driver) {
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

func (b *H6GPUPlatformBuilder) createGPUBuilder(
	engine akita.Engine,
	gpuDriver *driver.Driver,
) builders.Builder {
	gpuBuilder := builders.MakeH6GPUBuilder()
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
