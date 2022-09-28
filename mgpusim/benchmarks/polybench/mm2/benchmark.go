// Package mm2 implements the mm2 benchmark from Polybench.
package mm2

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
	Tmp                 driver.GPUPtr
	A                   driver.GPUPtr
	B                   driver.GPUPtr
	NI                  int32
	NJ                  int32
	NK                  int32
	NL                  int32
	Alpha               float32
	Beta                float32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}

// Kernel2Args list first set of kernel arguments
type Kernel2Args struct {
	Tmp                 driver.GPUPtr
	C                   driver.GPUPtr
	D                   driver.GPUPtr
	NI                  int32
	NJ                  int32
	NK                  int32
	NL                  int32
	Alpha               float32
	Beta                float32
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

	tmp             []float32
	a, b, c, d      []float32
	alpha, beta     float32
	NI, NJ, NK, NL  int
	dtmp            driver.GPUPtr
	da, db, dc, dd  driver.GPUPtr
	d_outputFromGPU []float32

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
		hsacoBytes, "mm2_kernel1")
	if b.kernel1 == nil {
		log.Panic("Failed to load kernel binary")
	}
	b.kernel2 = kernels.LoadProgramFromMemory(
		hsacoBytes, "mm2_kernel2")
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
	b.tmp = make([]float32, b.NI*b.NJ)
	b.a = make([]float32, b.NI*b.NK)
	b.b = make([]float32, b.NK*b.NJ)
	b.c = make([]float32, b.NL*b.NJ)
	b.d = make([]float32, b.NI*b.NL)
	b.d_outputFromGPU = make([]float32, b.NI*b.NL)
	b.alpha = 32412
	b.beta = 2123

	for i := 0; i < b.NI; i++ {
		for j := 0; j < b.NJ; j++ {
			b.a[i*b.NJ+j] = float32(i*j) / float32(b.NI)
		}
	}
	for i := 0; i < b.NK; i++ {
		for j := 0; j < b.NJ; j++ {
			b.b[i*b.NJ+j] = float32(i*(j+1)) / float32(b.NJ)
		}
	}

	for i := 0; i < b.NL; i++ {
		for j := 0; j < b.NJ; j++ {
			b.a[i*b.NJ+j] = float32(i*(j+3)) / float32(b.NL)
		}
	}

	for i := 0; i < b.NI; i++ {
		for j := 0; j < b.NL; j++ {
			b.a[i*b.NL+j] = float32(i*(j+2)) / float32(b.NK)
		}
	}

	if b.useUnifiedMemory {
		b.dtmp = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NI*b.NJ*4))
		b.da = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NI*b.NK*4))
		b.db = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NK*b.NJ*4))
		b.dc = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NL*b.NJ*4))
		b.dd = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NI*b.NL*4))
	} else if b.useLASPMemoryAlloc {
		b.dtmp = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NI*b.NJ*4), "div4")
		b.da = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NI*b.NK*4), "div4")
		b.db = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NK*b.NJ*4), "div4")
		b.dc = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NL*b.NJ*4), "div4")
		b.dd = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NI*b.NL*4), "div4")
	} else {
		b.dtmp = b.driver.AllocateMemory(b.context,
			uint64(b.NI*b.NJ*4))
		b.da = b.driver.AllocateMemory(b.context,
			uint64(b.NI*b.NK*4))
		b.db = b.driver.AllocateMemory(b.context,
			uint64(b.NK*b.NJ*4))
		b.dc = b.driver.AllocateMemory(b.context,
			uint64(b.NL*b.NJ*4))
		b.dd = b.driver.AllocateMemory(b.context,
			uint64(b.NI*b.NL*4))
	}
}

func (b *Benchmark) exec() {
	b.driver.MemCopyH2D(b.context, b.da, b.a)
	b.driver.MemCopyH2D(b.context, b.db, b.b)
	b.driver.MemCopyH2D(b.context, b.dc, b.c)
	b.driver.MemCopyH2D(b.context, b.dd, b.d)

	localSize := [3]uint16{32, 8, 1}
	globalSizeX := uint32(((b.NI-1)/32 + 1) * 32)
	globalSizeY := uint32(((b.NL-1)/8 + 1) * 8)
	globalSize := [3]uint32{globalSizeX, globalSizeY, 1}

	kernel1Arg := Kernel1Args{
		b.dtmp,
		b.da,
		b.db,
		int32(b.NI),
		int32(b.NJ),
		int32(b.NK),
		int32(b.NL),
		float32(b.alpha),
		float32(b.beta),
		0, 0, 0,
	}
	b.driver.LaunchKernel(b.context, b.kernel1,
		globalSize, localSize, &kernel1Arg)
	// b.driver.MemCopyD2H(b.context, b.a_debug, b.da)

	kernel2Arg := Kernel2Args{
		b.dtmp,
		b.dc,
		b.dd,
		int32(b.NI),
		int32(b.NJ),
		int32(b.NK),
		int32(b.NL),
		float32(b.alpha),
		float32(b.beta),
		0, 0, 0,
	}
	b.driver.LaunchKernel(b.context, b.kernel2,
		globalSize, localSize, &kernel2Arg)
	// b.driver.MemCopyD2H(b.context, b.x_debug, b.dx)

	b.driver.MemCopyD2H(b.context, b.d_outputFromGPU, b.dd)
}

// Verify verifies
func (b *Benchmark) Verify() {
	b.cpumm2()

	for i := 0; i < b.NI*b.NL; i++ {
		if b.d_outputFromGPU[i] != b.d[i] {
			log.Panicf("Mismatch at %d, expected %f, but get %f",
				i,
				b.d[i],
				b.d_outputFromGPU[i])
		}
	}

	log.Printf("Passed!\n")
}

func (b *Benchmark) cpumm2() {
	for i := 0; i < b.NI; i++ {
		for j := 0; j < b.NJ; j++ {
			b.tmp[i*b.NJ+j] = 0.0
			for k := 0; k < b.NK; k++ {
				b.tmp[i*b.NJ+j] = b.alpha * b.a[i*b.NK+k] * b.b[k*b.NJ+j]
			}
		}
	}
	for i := 0; i < b.NI; i++ {
		for j := 0; j < b.NL; j++ {
			b.d[i*b.NJ+j] *= b.beta
			for k := 0; k < b.NJ; k++ {
				b.d[i*b.NJ+j] += b.tmp[i*b.NK+k] * b.c[k*b.NJ+j]
			}
		}
	}
}

// SetLASPMemoryAlloc use Unified Memory
func (b *Benchmark) SetLASPMemoryAlloc() {
	b.useLASPMemoryAlloc = true
}
