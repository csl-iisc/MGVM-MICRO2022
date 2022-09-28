// Package adi implements the adi benchmark from Polybench.
package adi

import (
	"fmt"
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
	X                   driver.GPUPtr
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
	X                   driver.GPUPtr
	N                   int32
	Padding             int32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}

// Kernel3Args list first set of kernel arguments
type Kernel3Args struct {
	A                   driver.GPUPtr
	B                   driver.GPUPtr
	X                   driver.GPUPtr
	N                   int32
	Padding             int32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}

// Kernel4Args list first set of kernel arguments
type Kernel4Args struct {
	A                   driver.GPUPtr
	B                   driver.GPUPtr
	X                   driver.GPUPtr
	I1                  int32
	N                   int32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}

// Kernel5Args list first set of kernel arguments
type Kernel5Args struct {
	A                   driver.GPUPtr
	B                   driver.GPUPtr
	X                   driver.GPUPtr
	N                   int32
	Padding             int32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}

// Kernel6Args list first set of kernel arguments
type Kernel6Args struct {
	A                   driver.GPUPtr
	B                   driver.GPUPtr
	X                   driver.GPUPtr
	I1                  int32
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
	kernel3 *insts.HsaCo
	kernel4 *insts.HsaCo
	kernel5 *insts.HsaCo
	kernel6 *insts.HsaCo

	a, b, x         []float32
	N, Steps        int
	da, db, dx      driver.GPUPtr
	x_outputFromGPU []float32
	b_outputFromGPU []float32

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

func (b *Benchmark) loadProgram() {
	hsacoBytes := _escFSMustByte(false, "/kernels.hsaco")

	b.kernel1 = kernels.LoadProgramFromMemory(
		hsacoBytes, "adi_kernel1")
	if b.kernel1 == nil {
		log.Panic("Failed to load kernel binary")
	}
	b.kernel2 = kernels.LoadProgramFromMemory(
		hsacoBytes, "adi_kernel2")
	if b.kernel2 == nil {
		log.Panic("Failed to load kernel binary")
	}
	b.kernel3 = kernels.LoadProgramFromMemory(
		hsacoBytes, "adi_kernel3")
	if b.kernel3 == nil {
		log.Panic("Failed to load kernel binary")
	}
	b.kernel4 = kernels.LoadProgramFromMemory(
		hsacoBytes, "adi_kernel4")
	if b.kernel4 == nil {
		log.Panic("Failed to load kernel binary")
	}
	b.kernel5 = kernels.LoadProgramFromMemory(
		hsacoBytes, "adi_kernel5")
	if b.kernel5 == nil {
		log.Panic("Failed to load kernel binary")
	}
	b.kernel6 = kernels.LoadProgramFromMemory(
		hsacoBytes, "adi_kernel6")
	if b.kernel6 == nil {
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
	b.x = make([]float32, b.N*b.N)
	b.x_outputFromGPU = make([]float32, b.N*b.N)
	b.b_outputFromGPU = make([]float32, b.N*b.N)

	for i := 0; i < b.N; i++ {
		for j := 0; j < b.N; j++ {
			b.x[i*b.N+j] = float32(i*(j+1)+10) / float32(b.N)
			b.a[i*b.N+j] = float32((i-1)*(j+4)+2) / float32(b.N)
			b.b[i*b.N+j] = float32((i+3)*(j+7)+3) / float32(b.N)
		}
	}

	if b.useUnifiedMemory {
		b.da = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.N*b.N*4))
		b.db = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.N*b.N*4))
		b.dx = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.N*b.N*4))
	} else if b.useLASPMemoryAlloc {
		b.da = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.N*b.N*4), "div4")
		b.db = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.N*b.N*4), "div4")
		b.dx = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.N*b.N*4), "div4")
	} else {
		b.da = b.driver.AllocateMemory(b.context,
			uint64(b.N*b.N*4))
		b.db = b.driver.AllocateMemory(b.context,
			uint64(b.N*b.N*4))
		b.dx = b.driver.AllocateMemory(b.context,
			uint64(b.N*b.N*4))
	}
}

func (b *Benchmark) exec() {
	b.driver.MemCopyH2D(b.context, b.da, b.a)
	b.driver.MemCopyH2D(b.context, b.db, b.b)
	b.driver.MemCopyH2D(b.context, b.dx, b.x)

	localSize := [3]uint16{256, 1, 1}
	globalSizeX := uint32(((b.N-1)/256 + 1) * 256)
	globalSizeY := uint32(((b.N-1)/1 + 1) * 1)
	globalSize := [3]uint32{globalSizeX, globalSizeY, 1}

	for t := 0; t < b.Steps; t++ {
		kernel1Arg := Kernel1Args{
			b.da,
			b.db,
			b.dx,
			int32(b.N),
			0,
			0, 0, 0,
		}
		b.driver.LaunchKernel(b.context, b.kernel1,
			globalSize, localSize, &kernel1Arg)

		kernel2Arg := Kernel2Args{
			b.da,
			b.db,
			b.dx,
			int32(b.N),
			0,
			0, 0, 0,
		}
		b.driver.LaunchKernel(b.context, b.kernel2,
			globalSize, localSize, &kernel2Arg)

		kernel3Arg := Kernel3Args{
			b.da,
			b.db,
			b.dx,
			int32(b.N),
			0,
			0, 0, 0,
		}
		b.driver.LaunchKernel(b.context, b.kernel3,
			globalSize, localSize, &kernel3Arg)

		for i := 0; i < b.N; i++ {
			kernel4Arg := Kernel4Args{
				b.da,
				b.db,
				b.dx,
				int32(i),
				int32(b.N),
				0, 0, 0,
			}
			b.driver.LaunchKernel(b.context, b.kernel4,
				globalSize, localSize, &kernel4Arg)
		}

		kernel5Arg := Kernel5Args{
			b.da,
			b.db,
			b.dx,
			int32(b.N),
			0,
			0, 0, 0,
		}
		b.driver.LaunchKernel(b.context, b.kernel5,
			globalSize, localSize, &kernel5Arg)

		for i := 0; i < b.N; i++ {
			kernel6Arg := Kernel6Args{
				b.da,
				b.db,
				b.dx,
				int32(i),
				int32(b.N),
				0, 0, 0,
			}
			b.driver.LaunchKernel(b.context, b.kernel6,
				globalSize, localSize, &kernel6Arg)

		}

	}

	b.driver.MemCopyD2H(b.context, b.x_outputFromGPU, b.dx)
	b.driver.MemCopyD2H(b.context, b.b_outputFromGPU, b.db)
}

// Verify verifies
func (b *Benchmark) Verify() {
	b.cpuadi()

	// allow some amount of slack (not 0.001).
	for i := 0; i < b.N; i++ {
		for j := 0; j < b.N; j++ {
			if math.Abs(float64(b.x_outputFromGPU[i*b.N+j]-b.x[i*b.N+j])/float64(b.x[i*b.N+j])) > 0.1 {
				log.Panicf("Mismatch at %d, expected %f, but get %f",
					i*b.N+j,
					b.x[i*b.N+j],
					b.x_outputFromGPU[i*b.N+j])
			}
		}
	}

	// allow some amount of slack (not 0.001).
	for i := 0; i < b.N; i++ {
		for j := 0; j < b.N; j++ {
			if math.Abs(float64(b.b_outputFromGPU[i*b.N+j]-b.b[i*b.N+j])/float64(b.b[i*b.N+j])) > 0.1 {
				log.Panicf("Mismatch at %d, expected %f, but get %f",
					i*b.N+j,
					b.b[i*b.N+j],
					b.b_outputFromGPU[i*b.N+j])
			}
		}
	}

	log.Printf("Passed!\n")
}

func (b *Benchmark) cpuadi() {
	for t := 0; t < b.Steps; t++ {
		for i := 0; i < b.N; i++ {
			for j := 1; j < b.N; j++ {
				// X[i1*n + i2] = X[i1*n + i2] - X[i1*n + (i2-1)] * A[i1*n + i2] / B[i1*n + (i2-1)];
				// B[i1*n + i2] = B[i1*n + i2] - A[i1*n + i2] * A[i1*n + i2] / B[i1*n + (i2-1)];
				fmt.Println(i, j, b.x[i*b.N+j], b.x[i*b.N+(j-1)], b.a[i*b.N+j], b.b[i*b.N+(j-1)], b.x[i*b.N+j]-b.x[i*b.N+(j-1)]*b.a[i*b.N+j]/b.b[i*b.N+(j-1)])
				b.x[i*b.N+j] = b.x[i*b.N+j] - b.x[i*b.N+(j-1)]*b.a[i*b.N+j]/b.b[i*b.N+(j-1)]
				b.b[i*b.N+j] = b.b[i*b.N+j] - b.a[i*b.N+j]*b.a[i*b.N+j]/b.b[i*b.N+(j-1)]
				// 0.2 * (b.a[i*b.N+j] + b.a[i*b.N+(j-1)] +
				// b.a[i*b.N+(j+1)] + b.a[(i-1)*b.N+j] + b.a[(i+1)*b.N+j])
			}
		}
		// 	for i := 1; i < b.N-1; i++ {
		// 		for j := 1; j < b.N-1; j++ {
		// 			b.a[i*b.N+j] = b.b[i*b.N+j]
		// 		}
		// 	}
	}
}

func (b *Benchmark) SetLASPMemoryAlloc() {
	b.useLASPMemoryAlloc = true
}
