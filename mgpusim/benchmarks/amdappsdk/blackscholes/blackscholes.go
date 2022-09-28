// Package matrixtranspose implements the matrix transpose benchmark from
// AMDAPPSDK.
package blackscholes

import (
	"log"

	"gitlab.com/akita/mgpusim/driver"
	"gitlab.com/akita/mgpusim/insts"
	"gitlab.com/akita/mgpusim/kernels"
)

// KernelArgs defines kernel arguments
type KernelArgs struct {
	RandArray           driver.GPUPtr
	Width               uint32
	Padding             uint32
	CallArray           driver.GPUPtr
	PutArray            driver.GPUPtr
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

	Width, Height int

	RandArray                 []float32
	CallArray, PutArray       []float32
	DevRandArray              driver.GPUPtr
	DevCallArray, DevPutArray driver.GPUPtr

	useUnifiedMemory   bool
	useLASPMemoryAlloc bool
}

// NewBenchmark makes a new benchmark
func NewBenchmark(driver *driver.Driver) *Benchmark {
	b := new(Benchmark)
	b.driver = driver
	b.context = driver.Init()
	b.loadProgram()
	b.Width = 1
	b.Height = 1
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

// SetLASPMemoryAlloc use Unified Memory
func (b *Benchmark) SetLASPMemoryAlloc() {
	b.useLASPMemoryAlloc = true
}

func (b *Benchmark) loadProgram() {
	hsacoBytes := _escFSMustByte(false, "/kernels.hsaco")

	b.kernel = kernels.LoadProgramFromMemory(hsacoBytes, "blackScholes")
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
	numData := b.Height * b.Width * 4

	b.RandArray = make([]float32, numData)
	b.CallArray = make([]float32, numData)
	b.PutArray = make([]float32, numData)

	for i := 0; i < numData; i++ {
		b.RandArray[i] = float32(i)
		b.CallArray[i] = 0.0
		b.PutArray[i] = 0.0
	}

	if b.useUnifiedMemory {
		b.DevRandArray = b.driver.AllocateUnifiedMemory(
			b.context, uint64(numData*4))
		b.DevCallArray = b.driver.AllocateUnifiedMemory(
			b.context, uint64(numData*4))
		b.DevPutArray = b.driver.AllocateUnifiedMemory(
			b.context, uint64(numData*4))
	} else if b.useLASPMemoryAlloc {
		b.DevRandArray = b.driver.AllocateMemoryLASP(
			b.context, uint64(numData*4), "div4")
		b.DevCallArray = b.driver.AllocateMemoryLASP(
			b.context, uint64(numData*4), "div4")
		b.DevPutArray = b.driver.AllocateMemoryLASP(
			b.context, uint64(numData*4), "div4")
	} else {
		b.DevRandArray = b.driver.AllocateMemory(
			b.context, uint64(numData*4))
		b.DevCallArray = b.driver.AllocateMemory(
			b.context, uint64(numData*4))
		b.DevPutArray = b.driver.AllocateMemory(
			b.context, uint64(numData*4))
		b.driver.Distribute(b.context, b.DevRandArray, uint64(numData*4), b.gpus)
		b.driver.Distribute(b.context, b.DevCallArray, uint64(numData*4), b.gpus)
		b.driver.Distribute(b.context, b.DevPutArray, uint64(numData*4), b.gpus)
	}

	b.driver.MemCopyH2D(b.context, b.DevRandArray, b.RandArray)
	b.driver.MemCopyH2D(b.context, b.DevCallArray, b.CallArray)
	b.driver.MemCopyH2D(b.context, b.DevPutArray, b.PutArray)
}

func (b *Benchmark) exec() {

	for _, queue := range b.queues {

		kernArg := KernelArgs{
			b.DevRandArray,
			uint32(b.Width),
			0,
			b.DevCallArray,
			b.DevPutArray,
			0, 0, 0,
		}

		b.driver.EnqueueLaunchKernel(
			queue,
			b.kernel,
			[3]uint32{uint32(b.Width), uint32(b.Height), 1},
			[3]uint16{uint16(256), uint16(1), 1},
			&kernArg,
		)
	}

	for _, q := range b.queues {
		b.driver.DrainCommandQueue(q)
	}

	b.driver.MemCopyD2H(b.context, b.CallArray, b.DevCallArray)
	b.driver.MemCopyD2H(b.context, b.PutArray, b.DevPutArray)
}

// Verify verifies
func (b *Benchmark) Verify() {
	log.Printf("How will it pass if it is not implemented at all?")
}
