// Package tpacf
package tpacf

import (
	"fmt"
	"log"

	"gitlab.com/akita/mgpusim/driver"
	"gitlab.com/akita/mgpusim/insts"
	"gitlab.com/akita/mgpusim/kernels"
)

// KernelArgs defines kernel arguments
type KernelArgs struct {
	DevHist             driver.GPUPtr
	DevAllXData         driver.GPUPtr
	DevBins             driver.GPUPtr
	NumSets             int32
	NumElements         int32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}

// Benchmark defines a benchmark
type Benchmark struct {
	driver  *driver.Driver
	context *driver.Context
	gpuIDs  []int
	kernel  *insts.HsaCo

	NumSets     int32
	NumBins     int32
	NumElements int32

	HostHist    []float32
	HostAllData []float32
	HostBins    []float32

	useUnifiedMemory   bool
	useLASPMemoryAlloc bool
}

// NewBenchmark creates a new benchmark
func NewBenchmark(driver *driver.Driver) *Benchmark {
	b := new(Benchmark)
	b.driver = driver
	b.context = driver.Init()
	b.loadProgram()

	b.NumSets = 1024
	b.NumElements = 4194304

	return b
}

func (b *Benchmark) loadProgram() {
	hsacoBytes := _escFSMustByte(false, "/kernels.hsaco")

	b.kernel = kernels.LoadProgramFromMemory(
		hsacoBytes, "gen_hists")
	if b.kernel == nil {
		log.Panic("Failed to load kernel binary")
	}

}

// Run runs
func (b *Benchmark) Run() {
	b.initMem()
	b.exec()
}

func (b *Benchmark) SetUnifiedMemory() {
	b.useUnifiedMemory = true
}

// SetLASPMemoryAlloc use Unified Memory
func (b *Benchmark) SetLASPMemoryAlloc() {
	b.useLASPMemoryAlloc = true
}

func (b *Benchmark) initMem() {
	b.initData()
	b.allocateGPUMem()
}

func (b *Benchmark) initData() {
	b.HostHist = make([]float32, b.NumBins*(b.NumBins*2+1))
	b.HostAllData = make([]float32, b.NumBins*(b.NumBins*2+1))
}

func (b *Benchmark) allocateGPUMem() {
	// first initialize memory

	mem_size := (1 + b.NumSets) * b.NumElements * (12)
	f_mem_size := (1 + b.NumSets) * b.NumElements * (4)

	if b.useUnifiedMemory {
		panic("hello??")
	} else if b.useLASPMemoryAlloc {
		b.DevHist = b.driver.AllocateMemoryLASP(b.context, uint64(b.NumBins*(b.NumSets*2+1)), "div4")
		b.DevAllXData = b.driver.AllocateMemoryLASP(b.context, uint64(mem_size), "div4")
		b.DevBins = b.driver.AllocateMemoryLASP(b.context, uint64(b.NumBins), "div4")
	} else {
		panic("hello??")
	}

	b.driver.MemCopyH2D(b.context, b.DevTable, b.HostTable)
	b.driver.MemCopyH2D(b.context, b.DevStarts, b.HostStarts)
}

func (b *Benchmark) allocate(byteSize uint64) driver.GPUPtr {
	return nil
}

func (b *Benchmark) exec() {
}

func (b *Benchmark) runKernel1() {
	workSize := b.col - 1
	offsetR := 0
	offsetC := 0
	blockWidth := workSize / b.blockSize

	for blk := 1; blk <= workSize/b.blockSize; blk++ {
		globalSize := [3]uint32{uint32(b.blockSize * blk), 1, 1}
		localSize := [3]uint16{uint16(b.blockSize), 1, 1}

		args := KernelArgs{}

		b.driver.LaunchKernel(
			b.context,
			b.kernel1,
			globalSize,
			localSize,
			&args,
		)
	}
}

// Verify verifies
func (b *Benchmark) Verify() {

	fmt.Print("Passed!\n")
}
