// Package jacobi2d implements the jacobi2d benchmark from Polybench.
package jacobi2d

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
	B                   driver.GPUPtr
	N                   int32
	Padding             int32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}

// Kernel2Args list first set of kernel arguments
type Kernel2Args struct {
	A                   driver.GPUPtr
	B                   driver.GPUPtr
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

	a, b            []float32
	N, Steps        int
	da, db          driver.GPUPtr
	a_outputFromGPU []float32
	b_outputFromGPU []float32

	useUnifiedMemory      bool
	useLASPMemoryAlloc    bool
	useLASPHSLMemoryAlloc bool
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

// SetLASPMemoryAlloc use Unified Memory
func (b *Benchmark) SetLASPHSLMemoryAlloc() {
	b.useLASPHSLMemoryAlloc = true
}

func (b *Benchmark) loadProgram() {
	hsacoBytes := _escFSMustByte(false, "/kernels.hsaco")

	b.kernel1 = kernels.LoadProgramFromMemory(
		hsacoBytes, "runJacobi2D_kernel1")
	if b.kernel1 == nil {
		log.Panic("Failed to load kernel binary")
	}
	b.kernel2 = kernels.LoadProgramFromMemory(
		hsacoBytes, "runJacobi2D_kernel2")
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
	b.b = make([]float32, b.N*b.N)
	b.a_outputFromGPU = make([]float32, b.N*b.N)
	b.b_outputFromGPU = make([]float32, b.N*b.N)

	for i := 0; i < b.N; i++ {
		for j := 0; j < b.N; j++ {
			b.a[i] = float32(i*(j+2)+10) / float32(b.N)
			b.b[i] = float32((i-4)*(j-1)+11) / float32(b.N)
		}
	}

	if b.useUnifiedMemory {
		b.da = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.N*b.N*4))
		b.db = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.N*b.N*4))
	} else if b.useLASPMemoryAlloc {
		b.da = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.N*b.N*4), "div4")
		b.db = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.N*b.N*4), "div4")
	} else if b.useLASPHSLMemoryAlloc {
		b.da = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.N*b.N*4), "div4")
		b.db = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.N*b.N*4), "div4")
	} else {
		b.da = b.driver.AllocateMemory(b.context,
			uint64(b.N*b.N*4))
		b.db = b.driver.AllocateMemory(b.context,
			uint64(b.N*b.N*4))
	}
}

func (b *Benchmark) exec() {
	b.driver.MemCopyH2D(b.context, b.da, b.a)
	b.driver.MemCopyH2D(b.context, b.db, b.b)

	localSize := [3]uint16{32, 8, 1}
	globalSizeX := uint32(((b.N-1)/32 + 1) * 32)
	globalSizeY := uint32(((b.N-1)/8 + 1) * 8)
	globalSize := [3]uint32{globalSizeX, globalSizeY, 1}

	for t := 0; t < b.Steps; t++ {
		kernel1Arg := Kernel1Args{
			b.da,
			b.db,
			int32(b.N),
			0,
			0, 0, 0,
		}
		b.driver.LaunchKernel(b.context, b.kernel1,
			globalSize, localSize, &kernel1Arg)

		kernel2Arg := Kernel2Args{
			b.da,
			b.db,
			int32(b.N),
			0,
			0, 0, 0,
		}
		b.driver.LaunchKernel(b.context, b.kernel2,
			globalSize, localSize, &kernel2Arg)
	}

	b.driver.MemCopyD2H(b.context, b.a_outputFromGPU, b.da)
	b.driver.MemCopyD2H(b.context, b.b_outputFromGPU, b.db)
}

// Verify verifies
func (b *Benchmark) Verify() {
	b.cpujacobi2d()

	// allow some amount of slack (not 0.001).
	for i := 1; i < b.N-1; i++ {
		if math.Abs(float64(b.a_outputFromGPU[i]-b.a[i])/float64(b.a[i])) > 0.01 {
			log.Panicf("Mismatch at %d, expected %f, but get %f",
				i,
				b.a[i],
				b.a_outputFromGPU[i])
		}
	}
	for i := 1; i < b.N-1; i++ {
		if math.Abs(float64(b.b_outputFromGPU[i]-b.b[i])/float64(b.b[i])) > 0.01 {
			log.Panicf("Mismatch at %d, expected %f, but get %f",
				i,
				b.b[i],
				b.b_outputFromGPU[i])
		}
	}

	log.Printf("Passed!\n")
}

func (b *Benchmark) cpujacobi2d() {
	for t := 0; t < b.Steps; t++ {
		for i := 1; i < b.N-1; i++ {
			for j := 1; j < b.N-1; j++ {
				b.b[i*b.N+j] = 0.2 * (b.a[i*b.N+j] + b.a[i*b.N+(j-1)] +
					b.a[i*b.N+(j+1)] + b.a[(i-1)*b.N+j] + b.a[(i+1)*b.N+j])
			}
		}
		for i := 1; i < b.N-1; i++ {
			for j := 1; j < b.N-1; j++ {
				b.a[i*b.N+j] = b.b[i*b.N+j]
			}
		}
	}
}
