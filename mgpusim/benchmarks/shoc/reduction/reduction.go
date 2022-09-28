// Package reduction implements the stencil2d benchmark from the SHOC suite.
package reduction

import (
	"log"

	"gitlab.com/akita/mgpusim/driver"
	"gitlab.com/akita/mgpusim/insts"
	"gitlab.com/akita/mgpusim/kernels"
)

// ReductionArgs defines kernel arguments
type ReductionArgs struct {
	InputData           driver.GPUPtr
	OutputData          driver.GPUPtr
	SharedMemSize       driver.LocalPtr
	N                   uint32
	Padding             int32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}

// Benchmark defines a benchmark
type Benchmark struct {
	driver  *driver.Driver
	context *driver.Context
	gpus    []int
	queues  []*driver.CommandQueue

	reductionKernel *insts.HsaCo

	Size          uint32
	Iterations    uint32
	SharedMemSize uint32
	hInputData    []uint32
	hOutputData   []uint32
	dInputData    driver.GPUPtr
	dOutputData   driver.GPUPtr
	dSharedMem    driver.GPUPtr

	useUnifiedMemory      bool
	useLASPMemoryAlloc    bool
	useLASPHSLMemoryAlloc bool
}

// NewBenchmark returns a benchmark
func NewBenchmark(driver *driver.Driver) *Benchmark {
	b := new(Benchmark)
	b.driver = driver
	b.context = driver.Init()
	b.loadProgram()
	return b
}

// SelectGPU selects GPU
func (b *Benchmark) SelectGPU(gpus []int) {
	b.gpus = gpus
}

// SetUnifiedMemory uses Unified Memory
func (b *Benchmark) SetUnifiedMemory() {
	b.useUnifiedMemory = true
}

// SetLASPMemoryAlloc use Unified Memory
func (b *Benchmark) SetLASPMemoryAlloc() {
	b.useLASPMemoryAlloc = true
}

// SetLASPMemoryAlloc use Unified Memory
func (b *Benchmark) SetLASPHSLMemoryAlloc() {
	b.useLASPHSLMemoryAlloc = true
}

func (b *Benchmark) loadProgram() {
	hsacoBytes := _escFSMustByte(false, "/kernels.hsaco")

	b.reductionKernel = kernels.LoadProgramFromMemory(
		hsacoBytes, "reduce")
	if b.reductionKernel == nil {
		log.Panic("Failed to load kernel binary")
	}
}

// Run runs
func (b *Benchmark) Run() {
	for _, gpu := range b.gpus {
		b.driver.SelectGPU(b.context, gpu)
		b.queues = append(b.queues, b.driver.CreateCommandQueue(b.context))
	}

	b.initMem()
	b.exec()
}

func (b *Benchmark) initMem() {
	numData := b.Size
	b.hInputData = make([]uint32, numData)
	b.hOutputData = make([]uint32, numData)

	b.SharedMemSize = 256 * 4

	for i := 0; i < int(numData); i++ {
		b.hInputData[i] = uint32(i)
	}

	if b.useUnifiedMemory {
		b.dInputData = b.driver.AllocateUnifiedMemory(
			b.context, uint64(numData*4))
		b.dOutputData = b.driver.AllocateUnifiedMemory(
			b.context, uint64(numData*4))
	} else if b.useLASPMemoryAlloc {
		b.dInputData = b.driver.AllocateMemoryLASP(
			b.context, uint64(numData*4), "div4")
		b.dOutputData = b.driver.AllocateMemoryLASP(
			b.context, uint64(numData*4), "div4")
	} else if b.useLASPHSLMemoryAlloc {
		b.dInputData = b.driver.AllocateMemoryLASP(
			b.context, uint64(numData*4), "div4")
		b.dOutputData = b.driver.AllocateMemoryLASP(
			b.context, uint64(numData*4), "div4")
	} else {
		b.dInputData = b.driver.AllocateMemory(
			b.context, uint64(numData*4))
		b.dOutputData = b.driver.AllocateMemory(
			b.context, uint64(numData*4))
	}
	b.driver.MemCopyH2D(b.context, b.dInputData, b.hInputData)
}

func (b *Benchmark) exec() {

	for i := 0; i < int(b.Iterations); i++ {
		args := ReductionArgs{
			InputData:           b.dInputData,
			OutputData:          b.dOutputData,
			SharedMemSize:       256 * 4,
			N:                   b.Size,
			Padding:             0,
			HiddenGlobalOffsetX: 0,
			HiddenGlobalOffsetY: 0,
			HiddenGlobalOffsetZ: 0,
		}

		globalSize := [3]uint32{256 * 64, 1, 1}
		localSize := [3]uint16{256, 1, 1}
		b.driver.LaunchKernel(b.context,
			b.reductionKernel,
			globalSize, localSize,
			&args,
		)
	}

	b.driver.MemCopyD2H(b.context, b.hOutputData, b.dOutputData)
}

// Verify verfies
func (b *Benchmark) Verify() {
	log.Printf("Skipping Verification. !\n")
}
