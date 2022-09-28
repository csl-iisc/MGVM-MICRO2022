// Package convolution2d implements the convolution2d benchmark from Polybench.
package convolution2d

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
	NI                  int32
	NJ                  int32
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

	a, b            []float32
	NI, NJ          int
	da, db          driver.GPUPtr
	b_outputFromGPU []float32

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
		hsacoBytes, "Convolution2D_kernel")
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
	b.a = make([]float32, b.NI*b.NJ)
	b.b = make([]float32, b.NI*b.NJ)
	b.b_outputFromGPU = make([]float32, b.NI*b.NJ)

	for i := 0; i < b.NI; i++ {
		for j := 0; j < b.NJ; j++ {
			b.a[i*b.NJ+j] = float32(i*j) / (float32(b.NI * b.NJ))
		}
	}

	if b.useUnifiedMemory {
		b.da = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NI*b.NJ*4))
		b.db = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NI*b.NJ*4))
	} else if b.useLASPMemoryAlloc {
		b.da = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NI*b.NJ*4), "div4")
		b.db = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NI*b.NJ*4), "div4")
	} else if b.useLASPHSLMemoryAlloc {
		b.da = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NI*b.NJ*4), "div4")
		b.db = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NI*b.NJ*4), "div4")
	} else {
		b.da = b.driver.AllocateMemory(b.context,
			uint64(b.NI*b.NJ*4))
		b.db = b.driver.AllocateMemory(b.context,
			uint64(b.NI*b.NJ*4))
	}

	if b.useCustomHSL {
		// define cusotm HSL here
		// the number of TLB entries to stripe at
		b.driver.SetHSL(512 * 128)
	}
}

func (b *Benchmark) exec() {
	b.driver.MemCopyH2D(b.context, b.da, b.a)

	localSize := [3]uint16{32, 8, 1}
	globalSizeX := uint32(((b.NI-1)/32 + 1) * 32)
	globalSizeY := uint32(((b.NJ-1)/8 + 1) * 8)
	globalSize := [3]uint32{globalSizeX, globalSizeY, 1}

	kernel1Arg := Kernel1Args{
		b.da,
		b.db,
		int32(b.NI),
		int32(b.NJ),
		0, 0, 0,
	}
	b.driver.LaunchKernel(b.context, b.kernel1,
		globalSize, localSize, &kernel1Arg)
	// b.driver.MemCopyD2H(b.context, b.a_debug, b.da)

	b.driver.MemCopyD2H(b.context, b.b_outputFromGPU, b.db)
}

// Verify verifies
func (b *Benchmark) Verify() {
	b.cpuconvolution2d()

	for i := 0; i < b.NI; i++ {
		for j := 0; j < b.NJ; j++ {
			if math.Abs(float64(b.b_outputFromGPU[i*b.NJ+j]-b.b[i*b.NJ+j])) > 0.001 {
				log.Panicf("Mismatch at %d, %d, expected %f, but get %f",
					i, /*i*b.NJ+j*/
					j, /*b.b[i*b.NJ+j]*/
					b.b_outputFromGPU[i*b.NJ+j])
			}
		}
	}

	log.Printf("Passed!\n")
}

func (b *Benchmark) cpuconvolution2d() {
	c11 := float32(+0.2)
	c21 := float32(+0.5)
	c31 := float32(-0.8)
	c12 := float32(-0.3)
	c22 := float32(+0.6)
	c32 := float32(-0.9)
	c13 := float32(+0.4)
	c23 := float32(+0.7)
	c33 := float32(+0.10)
	for i := 1; i < b.NI-1; i++ {
		for j := 1; j < b.NJ-1; j++ {
			b.b[i*b.NJ+j] =
				c11*b.a[(i-1)*b.NJ+(j-1)] +
					c12*b.a[(i+0)*b.NJ+(j-1)] +
					c13*b.a[(i+1)*b.NJ+(j-1)] +
					c21*b.a[(i-1)*b.NJ+(j+0)] +
					c22*b.a[(i+0)*b.NJ+(j+0)] +
					c23*b.a[(i+1)*b.NJ+(j+0)] +
					c31*b.a[(i-1)*b.NJ+(j+1)] +
					c32*b.a[(i+0)*b.NJ+(j+1)] +
					c33*b.a[(i+1)*b.NJ+(j+1)]
		}
	}
}

func (b *Benchmark) SetCustomHSL() {
	b.useCustomHSL = true
}
