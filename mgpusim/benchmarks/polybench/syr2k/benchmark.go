// Package syr2k implements the syr2k benchmark from Polybench.
package syr2k

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
	C                   driver.GPUPtr
	Alpha               float32
	Beta                float32
	NI                  int32
	NJ                  int32
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

	NI, NJ      int
	alpha, beta float32
	a, b, c     []float32
	cOutput     []float32
	dA, dB, dC  driver.GPUPtr

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

func (b *Benchmark) loadProgram() {
	hsacoBytes := _escFSMustByte(false, "/kernels.hsaco")

	b.kernel1 = kernels.LoadProgramFromMemory(
		hsacoBytes, "syr2k_kernel")
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
	b.alpha = 32412.0
	b.beta = 2123.0
	b.a = make([]float32, b.NI*b.NJ)
	b.b = make([]float32, b.NI*b.NJ)
	b.c = make([]float32, b.NI*b.NI)
	b.cOutput = make([]float32, b.NI*b.NI)

	for i := 0; i < b.NI; i++ {
		for j := 0; j < b.NJ; j++ {
			b.a[i*b.NJ+j] = (float32(i) * float32(j)) / float32(b.NI)
			b.b[i*b.NJ+j] = (float32(i) * float32(j)) / float32(b.NI)
		}
	}
	for i := 0; i < b.NI; i++ {
		for j := 0; j < b.NI; j++ {
			b.a[i*b.NI+j] = (float32(i) * float32(j)) / float32(b.NI)
		}
	}

	if b.useUnifiedMemory {
		b.dA = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NI*b.NJ*4))
		b.dB = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NI*b.NJ*4))
		b.dC = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NI*b.NI*4))
	} else if b.useLASPMemoryAlloc {
		b.dA = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NI*b.NJ*4), "div4")
		b.dB = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NI*b.NJ*4), "div4")
		b.dC = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NI*b.NI*4), "div4")
	} else if b.useLASPHSLMemoryAlloc {
		b.dA = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NI*b.NJ*4), "div4")
		b.dB = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NI*b.NJ*4), "div4")
		b.dC = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NI*b.NI*4), "div4")
	} else {
		b.dA = b.driver.AllocateMemory(b.context,
			uint64(b.NI*b.NJ*4))
		b.dB = b.driver.AllocateMemory(b.context,
			uint64(b.NI*b.NJ*4))
		b.dC = b.driver.AllocateMemory(b.context,
			uint64(b.NI*b.NI*4))
	}
}

func (b *Benchmark) exec() {
	b.driver.MemCopyH2D(b.context, b.dA, b.a)
	b.driver.MemCopyH2D(b.context, b.dB, b.b)
	b.driver.MemCopyH2D(b.context, b.dC, b.c)

	localSize := [3]uint16{32, 8, 1}
	globalSizeX := uint32(((b.NJ-1)/32 + 1) * 32)
	globalSizeY := uint32(((b.NI-1)/8 + 1) * 8)
	globalSize := [3]uint32{globalSizeX, globalSizeY, 1}

	kernel1Arg := Kernel1Args{
		b.dA,
		b.dB,
		b.dC,
		float32(b.alpha),
		float32(b.beta),
		int32(b.NI),
		int32(b.NJ),
		0, 0, 0,
	}
	b.driver.LaunchKernel(b.context, b.kernel1,
		globalSize, localSize, &kernel1Arg)

	b.driver.MemCopyD2H(b.context, b.cOutput, b.dC)
}

// Verify verifies
func (b *Benchmark) Verify() {
	b.cpusyr2k()

	for i := 0; i < b.NI*b.NI; i++ {
		if math.Abs(float64((b.c[i]-b.cOutput[i])/b.c[i])) > 0.001 {
			log.Panicf("Mismatch at %d, expected %f, but get %f",
				i, b.c[i], b.cOutput[i])
		}
	}

	log.Printf("Passed!\n")
}

func (b *Benchmark) cpusyr2k() {
	for i := 0; i < b.NI; i++ {
		for j := 0; j < b.NI; j++ {
			b.c[i*b.NI+j] *= b.beta
		}
	}
	for i := 0; i < b.NI; i++ {
		for j := 0; j < b.NI; j++ {
			for k := 0; k < b.NJ; k++ {
				b.c[i*b.NI+j] += b.alpha * b.a[i*b.NJ+k] * b.b[j*b.NJ+k]
				b.c[i*b.NI+j] += b.alpha * b.b[i*b.NJ+k] * b.a[j*b.NJ+k]
			}
		}
	}
}

// SetLASPMemoryAlloc use Unified Memory
func (b *Benchmark) SetLASPMemoryAlloc() {
	b.useLASPMemoryAlloc = true
}

// SetLASPMemoryAlloc use Unified Memory
func (b *Benchmark) SetLASPHSLMemoryAlloc() {
	b.useLASPHSLMemoryAlloc = true
}
