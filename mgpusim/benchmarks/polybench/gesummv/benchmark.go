// Package gesummv implements the gesummv benchmark from Polybench.
package gesummv

import (
	"log"
	// "math"
	"math/rand"

	"gitlab.com/akita/mgpusim/driver"
	"gitlab.com/akita/mgpusim/insts"
	"gitlab.com/akita/mgpusim/kernels"
)

// Kernel1Args list first set of kernel arguments
type Kernel1Args struct {
	A                   driver.GPUPtr
	B                   driver.GPUPtr
	X                   driver.GPUPtr
	Y                   driver.GPUPtr
	Tmp                 driver.GPUPtr
	Alpha               float32
	Beta                float32
	N                   int32
	Padding             int32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}

// Benchmark defines a benchmark
type Benchmark struct {
	driver           *driver.Driver
	context          *driver.Context
	gpus             []int
	queues           []*driver.CommandQueue
	kernel1, kernel2 *insts.HsaCo

	N           int
	alpha, beta float32
	a, b        []float32
	x, y        []float32
	yOutput     []float32
	tmp         []float32
	dA, dB      driver.GPUPtr
	dX, dY      driver.GPUPtr
	dTmp        driver.GPUPtr

	useUnifiedMemory   bool
	useLASPMemoryAlloc bool
}

// NewBenchmark makes a new benchmark
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

func (b *Benchmark) loadProgram() {
	hsacoBytes := _escFSMustByte(false, "/kernels.hsaco")

	b.kernel1 = kernels.LoadProgramFromMemory(
		hsacoBytes, "gesummv_kernel")
	if b.kernel1 == nil {
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
	rand.Seed(1)
	b.alpha = 43532.0
	b.beta = 12313.0
	b.a = make([]float32, b.N*b.N)
	b.b = make([]float32, b.N*b.N)
	b.x = make([]float32, b.N)
	b.y = make([]float32, b.N)
	b.tmp = make([]float32, b.N)
	b.yOutput = make([]float32, b.N)

	for i := 0; i < b.N; i++ {
		b.x[i] = float32(i) / float32(b.N)
		for j := 0; j < b.N; j++ {
			b.a[i*b.N+j] = float32(i) * float32(j) / float32(b.N)
			b.b[i*b.N+j] = float32(i) * float32(j) / float32(b.N)
		}
	}

	if b.useUnifiedMemory {
		b.dA = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.N*b.N*4))
		b.dB = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.N*b.N*4))
		b.dX = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.N*4))
		b.dY = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.N*4))
		b.dTmp = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.N*4))
	} else if b.useLASPMemoryAlloc {
		b.dA = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.N*b.N*4), "div4")
		b.dB = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.N*b.N*4), "div4")
		b.dX = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.N*4), "div4")
		b.dY = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.N*4), "div4")
		b.dTmp = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.N*4), "div4")
	} else {
		b.dA = b.driver.AllocateMemory(b.context,
			uint64(b.N*b.N*4))
		b.dB = b.driver.AllocateMemory(b.context,
			uint64(b.N*b.N*4))
		b.dX = b.driver.AllocateMemory(b.context,
			uint64(b.N*4))
		b.dY = b.driver.AllocateMemory(b.context,
			uint64(b.N*4))
		b.dTmp = b.driver.AllocateMemory(b.context,
			uint64(b.N*4))
	}
}

func (b *Benchmark) exec() {
	b.driver.MemCopyH2D(b.context, b.dA, b.a)
	b.driver.MemCopyH2D(b.context, b.dB, b.b)
	b.driver.MemCopyH2D(b.context, b.dX, b.x)

	// width := 256
	width := 64
	localSize := [3]uint16{uint16(width), 1, 1}
	globalSizeX := uint32(((b.N-1)/width + 1) * width)
	// globalSizeY := uint32(((b.N-1)/1 + 1) * 1)
	globalSize := [3]uint32{globalSizeX, 1, 1}

	kernel1Arg := Kernel1Args{
		b.dA,
		b.dB,
		b.dX,
		b.dY,
		b.dTmp,
		float32(b.alpha),
		float32(b.beta),
		int32(b.N),
		0,
		0, 0, 0,
	}
	b.driver.LaunchKernel(b.context, b.kernel1,
		globalSize, localSize, &kernel1Arg)

	b.driver.MemCopyD2H(b.context, b.yOutput, b.dY)
}

// Verify verifies
func (b *Benchmark) Verify() {
	b.cpugesummv()

	for i := 0; i < b.N; i++ {
		if b.y[i] != b.yOutput[i] {
			log.Panicf("Mismatch at %d, expected %f, but get %f",
				i, b.y[i], b.yOutput[i])
		}
	}

	log.Printf("Passed!\n")
}

func (b *Benchmark) cpugesummv() {
	for i := 0; i < b.N; i++ {
		b.tmp[i] = 0.0
		b.y[i] = 0.0
		for j := 0; j < b.N; j++ {
			b.tmp[i] = b.a[i*b.N+j]*b.x[j] + b.tmp[i]
			b.y[i] = b.b[i*b.N+j]*b.x[j] + b.y[i]
		}
		b.y[i] = b.alpha*b.tmp[i] + b.beta*b.y[i]
	}
}
