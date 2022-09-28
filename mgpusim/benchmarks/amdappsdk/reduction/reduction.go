// Package matrixtranspose implements the matrix transpose benchmark from
// AMDAPPSDK.
package reduction

import (
	"log"

	"gitlab.com/akita/mgpusim/driver"
	"gitlab.com/akita/mgpusim/insts"
	"gitlab.com/akita/mgpusim/kernels"
)

// KernelArgs defines kernel arguments
type KernelArgs struct {
	Input               driver.GPUPtr
	Output              driver.GPUPtr
	SData               driver.GPUPtr
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

	kernel *insts.HsaCo

	GroupSize, VectorSize, Multiply uint32
	Length                          uint32

	Input, Output, SData          []int32
	DevInput, DevOutput, DevSData driver.GPUPtr

	useUnifiedMemory      bool
	useLASPMemoryAlloc    bool
	useLASPHSLMemoryAlloc bool
	useCustomHSL          bool
}

// NewBenchmark makes a new benchmark
func NewBenchmark(driver *driver.Driver) *Benchmark {
	b := new(Benchmark)
	b.driver = driver
	b.context = driver.Init()
	b.loadProgram()
	b.Length = 4096
	b.GroupSize = 256
	b.VectorSize = 1
	b.Multiply = 1
	return b
}

// SelectGPU selects GPU
func (b *Benchmark) SelectGPU(gpus []int) {
	b.gpus = gpus
}

// SetUnifiedMemory use Unified Memory
func (b *Benchmark) SetUnifiedMemory() {
	b.useUnifiedMemory = true
}

func (b *Benchmark) loadProgram() {
	hsacoBytes := _escFSMustByte(false, "/kernels.hsaco")

	b.kernel = kernels.LoadProgramFromMemory(hsacoBytes, "reduce")
	if b.kernel == nil {
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

	b.Length = b.Length / b.VectorSize
	b.Input = make([]int32, b.Length)

	numData := b.Length

	for i := 0; i < int(numData); i++ {
		b.Input[i] = int32(i)
	}

	if b.useUnifiedMemory {
		b.DevInput = b.driver.AllocateUnifiedMemory(
			b.context, uint64(numData*4))
		b.DevOutput = b.driver.AllocateUnifiedMemory(
			b.context, uint64(numData*4))
		b.DevSData = b.driver.AllocateUnifiedMemory(
			b.context, uint64(numData*4))
	} else if b.useLASPMemoryAlloc {
		b.DevInput = b.driver.AllocateMemoryLASP(
			b.context, uint64(numData*4), "div4")
		b.DevOutput = b.driver.AllocateMemoryLASP(
			b.context, uint64(numData*4), "div4")
		b.DevSData = b.driver.AllocateMemoryLASP(
			b.context, uint64(numData*4), "div4")
	} else if b.useLASPHSLMemoryAlloc {
		b.DevInput = b.driver.AllocateMemoryLASP(
			b.context, uint64(numData*4), "div4")
		b.DevOutput = b.driver.AllocateMemoryLASP(
			b.context, uint64(numData*4), "div4")
		b.DevSData = b.driver.AllocateMemoryLASP(
			b.context, uint64(numData*4), "div4")
	} else {
		b.DevInput = b.driver.AllocateMemory(
			b.context, uint64(numData*4))
		b.DevOutput = b.driver.AllocateMemory(
			b.context, uint64(numData*4))
		b.DevSData = b.driver.AllocateMemory(
			b.context, uint64(numData*4))
		// b.driver.Distribute(b.context, b.DevInput, uint64(numData*4), b.gpus)
		// b.driver.Distribute(b.context, b.DevOutput, uint64(numData*4), b.gpus)
		// b.driver.Distribute(b.context, b.DevSData, uint64(numData*4), b.gpus)
	}

	b.driver.MemCopyH2D(b.context, b.DevInput, b.Input)
}

func (b *Benchmark) exec() {

	for _, queue := range b.queues {

		kernArg := KernelArgs{
			b.DevInput,
			b.DevOutput,
			b.DevSData,
			0, 0, 0,
		}

		b.driver.EnqueueLaunchKernel(
			queue,
			b.kernel,
			[3]uint32{uint32(b.Length / b.Multiply), 1, 1},
			[3]uint16{uint16(b.GroupSize), 1, 1},
			&kernArg,
		)
	}

	for _, q := range b.queues {
		b.driver.DrainCommandQueue(q)
	}

	b.driver.MemCopyD2H(b.context, b.Output, b.DevOutput)
}

// Verify verifies
func (b *Benchmark) Verify() {
	log.Printf("How will it pass if it is not implemented at all?")
}

// SetLASPMemoryAlloc use Unified Memory
func (b *Benchmark) SetLASPMemoryAlloc() {
	b.useLASPMemoryAlloc = true
}

// SetLASPMemoryAlloc use Unified Memory
func (b *Benchmark) SetLASPHSLMemoryAlloc() {
	b.useLASPHSLMemoryAlloc = true
}
