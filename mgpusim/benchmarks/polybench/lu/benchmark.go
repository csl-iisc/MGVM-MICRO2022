// Package lu implements the lu benchmark from Polybench.
package lu

import (
	"log"
	"math"
	"math/rand"

	"gitlab.com/akita/mgpusim/driver"
	"gitlab.com/akita/mgpusim/insts"
	"gitlab.com/akita/mgpusim/kernels"
)

// Kernel1Args list first set of kernel arguments
type Kernel1Args struct {
	A                   driver.GPUPtr
	K                   int32
	N                   int32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}

// Kernel2Args list first set of kernel arguments
type Kernel2Args struct {
	A                   driver.GPUPtr
	K                   int32
	N                   int32
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

	a               []float32
	K, N            int
	da              driver.GPUPtr
	a_outputFromGPU []float32

	useUnifiedMemory bool
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

func (b *Benchmark) loadProgram() {
	hsacoBytes := _escFSMustByte(false, "/kernels.hsaco")

	b.kernel1 = kernels.LoadProgramFromMemory(
		hsacoBytes, "lu_kernel1")
	if b.kernel1 == nil {
		log.Panic("Failed to load kernel binary")
	}
	b.kernel2 = kernels.LoadProgramFromMemory(
		hsacoBytes, "lu_kernel2")
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
	b.a_outputFromGPU = make([]float32, b.N*b.N)

	for i := 0; i < b.N; i++ {
		for j := 0; j < b.N; j++ {
			b.a[i*b.N+j] = float32(i*j+1) / float32(b.N)
		}
	}

	if b.useUnifiedMemory {
		b.da = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.N*b.N*4))
	} else {
		b.da = b.driver.AllocateMemory(b.context,
			uint64(b.N*b.N*4))
	}
}

func (b *Benchmark) exec() {
	b.driver.MemCopyH2D(b.context, b.da, b.a)

	localSize := [3]uint16{256, 1, 1}
	globalSizeX := uint32(((b.N-1)/256 + 1) * 256)
	globalSize := [3]uint32{globalSizeX, 1, 1}

	kernel1Arg := Kernel1Args{
		b.da,
		int32(b.K),
		int32(b.N),
		0, 0, 0,
	}
	b.driver.LaunchKernel(b.context, b.kernel1,
		globalSize, localSize, &kernel1Arg)
	// b.driver.MemCopyD2H(b.context, b.a_debug, b.da)

	localSize = [3]uint16{32, 8, 1}
	globalSizeX = uint32(((b.N-1)/32 + 1) * 32)
	globalSizeY := uint32(((b.N-1)/8 + 1) * 8)
	globalSize = [3]uint32{globalSizeX, globalSizeY, 1}

	kernel2Arg := Kernel2Args{
		b.da,
		int32(b.K),
		int32(b.N),
		0, 0, 0,
	}
	b.driver.LaunchKernel(b.context, b.kernel2,
		globalSize, localSize, &kernel2Arg)
	// b.driver.MemCopyD2H(b.context, b.x_debug, b.dx)

	b.driver.MemCopyD2H(b.context, b.a_outputFromGPU, b.da)
}

// Verify verifies
func (b *Benchmark) Verify() {
	b.cpulu()

	for i := 0; i < b.N; i++ {
		for j := 0; j < b.N; j++ {
			if math.Abs(float64(b.a_outputFromGPU[i*b.N+j]-b.a[i*b.N+j])) > 0.1 {
				log.Panicf("Mismatch at %d, expected %f, but get %f",
					i*b.N+j,
					b.a[i*b.N+j],
					b.a_outputFromGPU[i*b.N+j])
			}
		}
	}

	log.Printf("Passed!\n")
}

func (b *Benchmark) cpulu() {
	// for k := 0; k < b.K; k++ {
	k := b.K
	n := b.N
	for j := k + 1; j < b.N; j++ {
		b.a[k*b.N+j] = b.a[k*b.N+j] / b.a[k*n+k]
	}
	for i := k + 1; i < b.N; i++ {
		for j := k + 1; j < b.N; j++ {
			b.a[i*b.N+j] = b.a[i*b.N+j] - b.a[i*b.N+k]*b.a[k*b.N+j]
		}
	}
	// }
}

func (b *Benchmark) SetLASPMemoryAlloc() {
}
