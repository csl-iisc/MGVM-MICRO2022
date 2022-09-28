// Package mm3 implements the mm3 benchmark from Polybench.
package mm3

import (
	"log"
	// "math"
	// "math/rand"

	"gitlab.com/akita/mgpusim/driver"
	"gitlab.com/akita/mgpusim/insts"
	"gitlab.com/akita/mgpusim/kernels"
)

// Kernel1Args list first set of kernel arguments
type Kernel1Args struct {
	A                   driver.GPUPtr
	B                   driver.GPUPtr
	E                   driver.GPUPtr
	NI                  int32
	NJ                  int32
	NK                  int32
	Padding             int32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}

// Kernel2Args list first set of kernel arguments
type Kernel2Args struct {
	C                   driver.GPUPtr
	D                   driver.GPUPtr
	F                   driver.GPUPtr
	NJ                  int32
	NL                  int32
	NM                  int32
	Padding             int32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}

// Kernel3Args list first set of kernel arguments
type Kernel3Args struct {
	E                   driver.GPUPtr
	F                   driver.GPUPtr
	G                   driver.GPUPtr
	NI                  int32
	NL                  int32
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

	a, b, c, d         []float32
	e, f, g            []float32
	NI, NJ, NK, NL, NM int
	da, db, dc, dd     driver.GPUPtr
	de, df, dg         driver.GPUPtr

	g_outputFromGPU []float32

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
		hsacoBytes, "mm3_kernel1")
	if b.kernel1 == nil {
		log.Panic("Failed to load kernel binary")
	}
	b.kernel2 = kernels.LoadProgramFromMemory(
		hsacoBytes, "mm3_kernel2")
	if b.kernel2 == nil {
		log.Panic("Failed to load kernel binary")
	}
	b.kernel3 = kernels.LoadProgramFromMemory(
		hsacoBytes, "mm3_kernel3")
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
	b.e = make([]float32, b.NI*b.NJ)
	b.a = make([]float32, b.NI*b.NK)
	b.b = make([]float32, b.NK*b.NJ)
	b.f = make([]float32, b.NJ*b.NL)
	b.c = make([]float32, b.NJ*b.NM)
	b.d = make([]float32, b.NM*b.NL)
	b.g = make([]float32, b.NI*b.NL)
	b.g_outputFromGPU = make([]float32, b.NI*b.NL)

	for i := 0; i < b.NI; i++ {
		for j := 0; j < b.NK; j++ {
			b.a[i*b.NK+j] = float32(i*j) / float32(b.NI)
		}
	}
	for i := 0; i < b.NK; i++ {
		for j := 0; j < b.NJ; j++ {
			b.b[i*b.NJ+j] = float32(i*(j+1)) / float32(b.NJ)
		}
	}

	for i := 0; i < b.NJ; i++ {
		for j := 0; j < b.NM; j++ {
			b.c[i*b.NM+j] = float32(i*(j+3)) / float32(b.NL)
		}
	}

	for i := 0; i < b.NM; i++ {
		for j := 0; j < b.NL; j++ {
			b.d[i*b.NL+j] = float32(i*(j+2)) / float32(b.NK)
		}
	}

	if b.useUnifiedMemory {
		b.da = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NI*b.NK*4))
		b.db = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NK*b.NJ*4))
		b.dc = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NJ*b.NM*4))
		b.dd = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NM*b.NL*4))
		b.de = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NI*b.NJ*4))
		b.df = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NJ*b.NL*4))
		b.dg = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NI*b.NL*4))
	} else {
		b.da = b.driver.AllocateMemory(b.context,
			uint64(b.NI*b.NK*4))
		b.db = b.driver.AllocateMemory(b.context,
			uint64(b.NK*b.NJ*4))
		b.dc = b.driver.AllocateMemory(b.context,
			uint64(b.NJ*b.NM*4))
		b.dd = b.driver.AllocateMemory(b.context,
			uint64(b.NM*b.NL*4))
		b.de = b.driver.AllocateMemory(b.context,
			uint64(b.NI*b.NJ*4))
		b.df = b.driver.AllocateMemory(b.context,
			uint64(b.NJ*b.NL*4))
		b.dg = b.driver.AllocateMemory(b.context,
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
		b.da,
		b.db,
		b.de,
		int32(b.NI),
		int32(b.NJ),
		int32(b.NK),
		0,
		0, 0, 0,
	}
	b.driver.LaunchKernel(b.context, b.kernel1,
		globalSize, localSize, &kernel1Arg)

	globalSizeX = uint32(((b.NL-1)/32 + 1) * 32)
	globalSizeY = uint32(((b.NJ-1)/8 + 1) * 8)
	globalSize = [3]uint32{globalSizeX, globalSizeY, 1}

	kernel2Arg := Kernel2Args{
		b.dc,
		b.dd,
		b.df,
		int32(b.NJ),
		int32(b.NL),
		int32(b.NM),
		0,
		0, 0, 0,
	}
	b.driver.LaunchKernel(b.context, b.kernel2,
		globalSize, localSize, &kernel2Arg)

	globalSizeX = uint32(((b.NL-1)/32 + 1) * 32)
	globalSizeY = uint32(((b.NI-1)/8 + 1) * 8)
	globalSize = [3]uint32{globalSizeX, globalSizeY, 1}

	kernel3Arg := Kernel3Args{
		b.de,
		b.df,
		b.dg,
		int32(b.NI),
		int32(b.NL),
		int32(b.NJ),
		0,
		0, 0, 0,
	}
	b.driver.LaunchKernel(b.context, b.kernel3,
		globalSize, localSize, &kernel3Arg)
	// b.driver.MemCopyD2H(b.context, b.x_debug, b.dx)

	b.driver.MemCopyD2H(b.context, b.g_outputFromGPU, b.dg)
}

// Verify verifies
func (b *Benchmark) Verify() {
	b.cpumm3()

	for i := 0; i < b.NI*b.NL; i++ {
		if b.g_outputFromGPU[i] != b.g[i] {
			log.Panicf("Mismatch at %d, expected %f, but get %f",
				i,
				b.g[i],
				b.g_outputFromGPU[i])
		}
	}

	log.Printf("Passed!\n")
}

func (b *Benchmark) cpumm3() {
	for i := 0; i < b.NI; i++ {
		for j := 0; j < b.NJ; j++ {
			b.e[i*b.NJ+j] = 0.0
			for k := 0; k < b.NK; k++ {
				b.e[i*b.NJ+j] += b.a[i*b.NK+k] * b.b[k*b.NJ+j]
			}
		}
	}
	for i := 0; i < b.NJ; i++ {
		for j := 0; j < b.NL; j++ {
			b.f[i*b.NL+j] = 0.0
			for k := 0; k < b.NK; k++ {
				b.f[i*b.NL+j] += b.c[i*b.NK+k] * b.d[k*b.NL+j]
			}
		}
	}
	for i := 0; i < b.NI; i++ {
		for j := 0; j < b.NL; j++ {
			b.g[i*b.NL+j] = 0.0
			for k := 0; k < b.NJ; k++ {
				b.g[i*b.NL+j] += b.e[i*b.NJ+k] * b.f[k*b.NL+j]
			}
		}
	}
}

func (b *Benchmark) SetLASPMemoryAlloc() {
}
