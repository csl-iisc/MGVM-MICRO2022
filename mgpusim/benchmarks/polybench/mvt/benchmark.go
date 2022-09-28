// Package mvt implements the mvt benchmark from Polybench.
package mvt

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
	X1                  driver.GPUPtr
	Y1                  driver.GPUPtr
	N                   int32
	Padding             int32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}

// Kernel2Args list first set of kernel arguments
type Kernel2Args struct {
	A                   driver.GPUPtr
	X2                  driver.GPUPtr
	Y2                  driver.GPUPtr
	N                   int32
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
	kernel1 *insts.HsaCo
	kernel2 *insts.HsaCo

	a                []float32
	x1, x2           []float32
	y1, y2           []float32
	N                int
	da               driver.GPUPtr
	dx1, dx2         driver.GPUPtr
	dy1, dy2         driver.GPUPtr
	x1_outputFromGPU []float32
	x2_outputFromGPU []float32

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
		hsacoBytes, "mvt_kernel1")
	if b.kernel1 == nil {
		log.Panic("Failed to load kernel binary")
	}
	b.kernel2 = kernels.LoadProgramFromMemory(
		hsacoBytes, "mvt_kernel2")
	if b.kernel2 == nil {
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
	b.a = make([]float32, b.N*b.N)
	b.x1 = make([]float32, b.N)
	b.x2 = make([]float32, b.N)
	b.y1 = make([]float32, b.N)
	b.y2 = make([]float32, b.N)
	b.x1_outputFromGPU = make([]float32, b.N)
	b.x2_outputFromGPU = make([]float32, b.N)

	for i := 0; i < b.N; i++ {
		b.x1[i] = float32(i) / float32(b.N)
		b.x2[i] = float32(i+1) / float32(b.N)
		b.y1[i] = float32(i+3) / float32(b.N)
		b.y2[i] = float32(i+4) / float32(b.N)
		for j := 0; j < b.N; j++ {
			b.a[i*b.N+j] = float32(i*j) / float32(b.N)
		}
	}

	if b.useUnifiedMemory {
		b.da = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.N*b.N*4))
		b.dx1 = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.N*4))
		b.dx2 = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.N*4))
		b.dy1 = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.N*4))
		b.dy2 = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.N*4))
	} else if b.useLASPMemoryAlloc {
		b.da = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.N*b.N*4), "div4")
		b.dx1 = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.N*4), "div4")
		b.dx2 = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.N*4), "div4")
		b.dy1 = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.N*4), "div4")
		b.dy2 = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.N*4), "div4")
	} else {
		b.da = b.driver.AllocateMemory(b.context,
			uint64(b.N*b.N*4))
		b.dx1 = b.driver.AllocateMemory(b.context,
			uint64(b.N*4))
		b.dx2 = b.driver.AllocateMemory(b.context,
			uint64(b.N*4))
		b.dy1 = b.driver.AllocateMemory(b.context,
			uint64(b.N*4))
		b.dy2 = b.driver.AllocateMemory(b.context,
			uint64(b.N*4))
	}
}

func (b *Benchmark) exec() {
	b.driver.MemCopyH2D(b.context, b.da, b.a)
	b.driver.MemCopyH2D(b.context, b.dx1, b.x1)
	b.driver.MemCopyH2D(b.context, b.dx2, b.x2)
	b.driver.MemCopyH2D(b.context, b.dy1, b.y1)
	b.driver.MemCopyH2D(b.context, b.dy2, b.y2)

	localSize := [3]uint16{32, 8, 1}
	globalSizeX := uint32(((b.N-1)/32 + 1) * 32)
	globalSize := [3]uint32{globalSizeX, 1, 1}

	kernel1Arg := Kernel1Args{
		b.da,
		b.dx1,
		b.dy1,
		int32(b.N),
		0,
		0, 0, 0,
	}
	b.driver.LaunchKernel(b.context, b.kernel1,
		globalSize, localSize, &kernel1Arg)
	// b.driver.MemCopyD2H(b.context, b.a_debug, b.da)

	kernel2Arg := Kernel2Args{
		b.da,
		b.dx2,
		b.dy2,
		int32(b.N),
		0,
		0, 0, 0,
	}
	b.driver.LaunchKernel(b.context, b.kernel2,
		globalSize, localSize, &kernel2Arg)
	// b.driver.MemCopyD2H(b.context, b.x_debug, b.dx)

	b.driver.MemCopyD2H(b.context, b.x1_outputFromGPU, b.dx1)
	b.driver.MemCopyD2H(b.context, b.x2_outputFromGPU, b.dx2)
}

// Verify verifies
func (b *Benchmark) Verify() {
	b.cpumvt()

	for i := 0; i < b.N; i++ {
		if b.x1_outputFromGPU[i] != b.x1[i] {
			log.Panicf("Mismatch at %d, expected %f, but get %f",
				i,
				b.x1[i],
				b.x1_outputFromGPU[i])
		}
	}

	log.Printf("Passed!\n")
}

func (b *Benchmark) cpumvt() {
	for i := 0; i < b.N; i++ {
		for j := 0; j < b.N; j++ {
			b.x1[i] = b.x1[i] + b.a[i*b.N+j]*b.y1[j]
		}
	}
	for i := 0; i < b.N; i++ {
		for j := 0; j < b.N; j++ {
			b.x2[i] = b.x2[i] + b.a[i*b.N+j]*b.y2[j]
		}
	}
}
