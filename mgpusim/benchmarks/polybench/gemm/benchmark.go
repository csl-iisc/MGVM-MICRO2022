// Package gemm implements the gemm benchmark from Polybench.
package gemm

import (
	"log"
	// "math"
	// "fmt"
	"math/rand"

	"gitlab.com/akita/mgpusim/driver"
	"gitlab.com/akita/mgpusim/insts"
	"gitlab.com/akita/mgpusim/kernels"
)

// Kernel1Args list first set of kernel arguments
type KernelArgs struct {
	A     driver.GPUPtr
	B     driver.GPUPtr
	C     driver.GPUPtr
	Alpha float32
	Beta  float32
	NI    int32
	NJ    int32
	NK    int32
}

// Benchmark defines a benchmark
type Benchmark struct {
	driver  *driver.Driver
	context *driver.Context
	gpus    []int
	queues  []*driver.CommandQueue
	kernel  *insts.HsaCo

	a, b, c         []float32
	c_outputFromGPU []float32
	da, db, dc      driver.GPUPtr
	NI, NJ, NK      int
	Alpha, Beta     float32

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

	b.kernel = kernels.LoadProgramFromMemory(
		hsacoBytes, "gemm")
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
	rand.Seed(1)
	b.a = make([]float32, b.NI*b.NK)
	b.b = make([]float32, b.NK*b.NJ)
	b.c = make([]float32, b.NI*b.NJ)
	b.c_outputFromGPU = make([]float32, b.NI*b.NJ)

	for i := 0; i < b.NI; i++ {
		for j := 0; j < b.NK; j++ {
			b.a[i*b.NJ+j] = (float32(i*j) / float32(b.NI))
		}
	}
	for i := 0; i < b.NK; i++ {
		for j := 0; j < b.NJ; j++ {
			b.b[i*b.NJ+j] = (float32(i*j) / float32(b.NI))
		}
	}
	for i := 0; i < b.NI; i++ {
		for j := 0; j < b.NJ; j++ {
			b.c[i*b.NJ+j] = (float32(i*j) / float32(b.NI))
		}
	}

	if b.useUnifiedMemory {
		b.da = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NI*b.NK*4))
		b.db = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NK*b.NJ*4))
		b.dc = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NI*b.NJ*4))
	} else {
		b.da = b.driver.AllocateMemory(b.context,
			uint64(b.NI*b.NK*4))
		b.db = b.driver.AllocateMemory(b.context,
			uint64(b.NK*b.NJ*4))
		b.dc = b.driver.AllocateMemory(b.context,
			uint64(b.NI*b.NJ*4))
	}
}

func (b *Benchmark) exec() {
	b.driver.MemCopyH2D(b.context, b.da, b.a)
	b.driver.MemCopyH2D(b.context, b.db, b.b)
	b.driver.MemCopyH2D(b.context, b.dc, b.c)

	localSize := [3]uint16{256, 1, 1}
	globalSizeX := uint32(((b.NJ-1)/256 + 1) * 256)
	globalSizeY := uint32(((b.NI-1)/256 + 1) * 256)
	globalSize := [3]uint32{globalSizeX, globalSizeY, 1}

	kernelArg := KernelArgs{
		A:     b.da,
		B:     b.db,
		C:     b.dc,
		Alpha: b.Alpha,
		Beta:  b.Beta,
		NI:    int32(b.NI),
		NJ:    int32(b.NJ),
		NK:    int32(b.NK),
	}
	b.driver.LaunchKernel(b.context, b.kernel,
		globalSize, localSize, &kernelArg)

	b.driver.MemCopyD2H(b.context, b.c_outputFromGPU, b.dc)
}

// Verify verifies
func (b *Benchmark) Verify() {
	b.cpuGemm()

	for i := 0; i < b.NI; i++ {
		for j := 0; j < b.NJ; j++ {
			// fmt.Println(b.c_outputFromGPU[i*b.NJ+j], b.c[i*b.NJ+j])
			if b.c_outputFromGPU[i*b.NJ+j] != b.c[i*b.NJ+j] {
				log.Panicf("Mismatch at %d, expected %f, but get %f",
					i*b.NJ+j,
					b.c_outputFromGPU[i*b.NJ+j],
					b.c[i*b.NJ+j])
			}
		}
	}

	log.Printf("Passed!\n")
}

func (b *Benchmark) cpuGemm() {
	for i := 0; i < b.NI; i++ {
		for j := 0; j < b.NJ; j++ {
			b.c[i*b.NJ+j] *= b.Beta
			for k := 0; k < b.NK; k++ {
				b.c[i*b.NJ+j] += b.Alpha * b.a[i*b.NK+k] * b.b[k*b.NJ+j]
			}
		}
	}
}

func (b *Benchmark) SetLASPMemoryAlloc() {
}
