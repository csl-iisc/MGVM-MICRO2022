// Package gramschmidt implements the gramschmidt benchmark from Polybench.
package gramschmidt

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
	R                   driver.GPUPtr
	Q                   driver.GPUPtr
	K                   int32
	NI                  int32
	NJ                  int32
	Padding             int32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}

// Kernel2Args list first set of kernel arguments
type Kernel2Args struct {
	A                   driver.GPUPtr
	R                   driver.GPUPtr
	Q                   driver.GPUPtr
	K                   int32
	NI                  int32
	NJ                  int32
	Padding             int32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}

// Kernel3Args list first set of kernel arguments
type Kernel3Args struct {
	A                   driver.GPUPtr
	R                   driver.GPUPtr
	Q                   driver.GPUPtr
	K                   int32
	NI                  int32
	NJ                  int32
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
	kernel3 *insts.HsaCo

	a, r, q         []float32
	K, NI, NJ       int
	da, dr, dq      driver.GPUPtr
	a_outputFromGPU []float32

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
		hsacoBytes, "gramschmidt_kernel1")
	if b.kernel1 == nil {
		log.Panic("Failed to load kernel binary")
	}
	b.kernel2 = kernels.LoadProgramFromMemory(
		hsacoBytes, "gramschmidt_kernel2")
	if b.kernel2 == nil {
		log.Panic("Failed to load kernel binary")
	}
	b.kernel3 = kernels.LoadProgramFromMemory(
		hsacoBytes, "gramschmidt_kernel3")
	if b.kernel3 == nil {
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
	b.a = make([]float32, b.NI*b.NJ)
	b.r = make([]float32, b.NJ*b.NJ)
	b.q = make([]float32, b.NI*b.NJ)

	for i := 0; i < b.NI; i++ {
		for j := 0; j < b.NJ; j++ {
			b.a[i*b.NJ+j] = float32(i*j) / float32(b.NI)
			b.q[i*b.NJ+j] = float32(i*(j+1)) / float32(b.NJ)
		}
	}
	for i := 0; i < b.NJ; i++ {
		for j := 0; j < b.NJ; j++ {
			b.r[i*b.NJ+j] = float32(i*(j+2)) / float32(b.NJ)
		}
	}

	if b.useUnifiedMemory {
		b.da = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NI*b.NJ*4))
		b.dr = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NJ*b.NJ*4))
		b.dq = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NI*b.NJ*4))
	} else if b.useLASPMemoryAlloc {
		b.da = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NI*b.NJ*4), "div4")
		b.dr = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NJ*b.NJ*4), "div4")
		b.dq = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NI*b.NJ*4), "div4")
	} else {
		b.da = b.driver.AllocateMemory(b.context,
			uint64(b.NI*b.NJ*4))
		b.dr = b.driver.AllocateMemory(b.context,
			uint64(b.NJ*b.NJ*4))
		b.dq = b.driver.AllocateMemory(b.context,
			uint64(b.NI*b.NJ*4))
	}
}

func (b *Benchmark) exec() {
	b.driver.MemCopyH2D(b.context, b.da, b.a)

	localSize := [3]uint16{256, 1, 1}

	globalSizeX := uint32(((b.NJ-1)/256 + 1) * 256)
	globalSizeY := uint32(((b.NJ-1)/1 + 1) * 1)
	globalSizeK1 := [3]uint32{256, 1, 1}
	globalSizeK2 := [3]uint32{globalSizeX, globalSizeY, 1}
	globalSizeK3 := [3]uint32{globalSizeX, globalSizeY, 1}

	kernel1Arg := Kernel1Args{
		b.da,
		b.dr,
		b.dq,
		int32(b.K),
		int32(b.NI),
		int32(b.NJ),
		0,
		0, 0, 0,
	}
	b.driver.LaunchKernel(b.context, b.kernel1,
		globalSizeK1, localSize, &kernel1Arg)
	// b.driver.MemCopyD2H(b.context, b.a_debug, b.da)

	kernel2Arg := Kernel2Args{
		b.da,
		b.dr,
		b.dq,
		int32(b.K),
		int32(b.NI),
		int32(b.NJ),
		0,
		0, 0, 0,
	}
	b.driver.LaunchKernel(b.context, b.kernel2,
		globalSizeK2, localSize, &kernel2Arg)
	// b.driver.MemCopyD2H(b.context, b.x_debug, b.dx)

	kernel3Arg := Kernel3Args{
		b.da,
		b.dr,
		b.dq,
		int32(b.K),
		int32(b.NI),
		int32(b.NJ),
		0,
		0, 0, 0,
	}
	b.driver.LaunchKernel(b.context, b.kernel3,
		globalSizeK3, localSize, &kernel3Arg)

	// b.driver.MemCopyD2H(b.context, b.a_outputFromGPU, b.da)
}

// Verify verifies
func (b *Benchmark) Verify() {
	b.cpuGramSchmidt()

	for i := 0; i < b.NI; i++ {
		for j := 0; j < b.NJ; j++ {
			if b.a_outputFromGPU[i*b.NJ+j] != b.a[i*b.NJ+j] {
				log.Panicf("Mismatch at %d, expected %f, but get %f",
					i*b.NJ+j,
					b.a[i*b.NJ+j],
					b.a_outputFromGPU[i*b.NJ+j])
			}
		}
	}

	log.Printf("Passed!\n")
}

func (b *Benchmark) cpuGramSchmidt() {
	fmt.Println("Starting cpu execution")
	norm := float32(0.0)
	for k := 0; k < b.NJ; k++ {
		norm = 0.0
		for i := 0; i < b.NI; i++ {
			norm += b.a[i*b.K+k] * b.a[i*b.K+k]
		}
		b.r[k*b.K+k] = float32(math.Sqrt(float64(norm)))
		for i := 0; i < b.NI; i++ {
			b.q[i*b.K+k] = b.a[i*b.K+k] / b.r[k*b.K+k]
		}
		for j := 0; j < b.NJ; j++ {
			b.r[k*b.NJ+j] = 0
			for i := 0; i < b.NI; i++ {
				b.r[k*b.NJ+j] += b.q[i*b.K*k] * b.a[i*b.NJ+j]
			}
			for i := 0; i < b.NI; i++ {
				b.a[i*b.NJ+j] += b.a[i*b.K*j] - b.q[i*b.K+k]*b.r[k*b.K+j]
			}
		}
	}
}

// SetLASPMemoryAlloc use Unified Memory
func (b *Benchmark) SetLASPMemoryAlloc() {
	b.useLASPMemoryAlloc = true
}
