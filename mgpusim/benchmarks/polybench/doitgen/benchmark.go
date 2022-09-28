// Package doitgen implements the doitgen benchmark from Polybench.
package doitgen

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
	NR                  int32
	NQ                  int32
	NP                  int32
	Padding             int32
	A                   driver.GPUPtr
	C4                  driver.GPUPtr
	Sum                 driver.GPUPtr
	R                   int32
	Padding2            int32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}

// Kernel2Args list second set of kernel arguments
type Kernel2Args struct {
	NR                  int32
	NQ                  int32
	NP                  int32
	Padding             int32
	A                   driver.GPUPtr
	C4                  driver.GPUPtr
	Sum                 driver.GPUPtr
	R                   int32
	Padding2            int32
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

	NR, NQ, NP    int
	a, c4, sum    []float32
	dA, dC4, dSum driver.GPUPtr
	cpuSum        []float32

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
		hsacoBytes, "doitgen_kernel1")
	if b.kernel1 == nil {
		log.Panic("Failed to load kernel binary")
	}

	b.kernel2 = kernels.LoadProgramFromMemory(
		hsacoBytes, "doitgen_kernel2")
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
	b.a = make([]float32, b.NR*b.NQ*b.NP)
	b.c4 = make([]float32, b.NP*b.NP)
	b.sum = make([]float32, b.NR*b.NQ*b.NP)

	for r := 0; r < b.NR; r++ {
		for q := 0; q < b.NQ; q++ {
			for p := 0; p < b.NP; p++ {
				b.a[r*b.NQ*b.NP+q*b.NP+p] = (float32(q*p) / float32(b.NP))
			}
		}
	}
	for i := 0; i < b.NP; i++ {
		for j := 0; j < b.NP; j++ {
			b.c4[i*b.NP+j] = float32(i*j) / float32(b.NP)
		}
	}

	if b.useUnifiedMemory {
		b.dA = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NR*b.NQ*b.NP*4))
		b.dC4 = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NP*b.NP*4))
		b.dSum = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NR*b.NQ*b.NP*4))
	} else if b.useLASPMemoryAlloc {
		b.dA = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NR*b.NQ*b.NP*4), "div4")
		b.dC4 = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NP*b.NP*4), "div4")
		b.dSum = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NR*b.NQ*b.NP*4), "div4")

	} else {
		b.dA = b.driver.AllocateMemory(b.context,
			uint64(b.NR*b.NQ*b.NP*4))
		b.dC4 = b.driver.AllocateMemory(b.context,
			uint64(b.NP*b.NP*4))
		b.dSum = b.driver.AllocateMemory(b.context,
			uint64(b.NR*b.NQ*b.NP*4))
	}
}

func (b *Benchmark) exec() {
	b.driver.MemCopyH2D(b.context, b.dA, b.a)
	b.driver.MemCopyH2D(b.context, b.dC4, b.c4)

	localSizeX := 32
	localSizeY := 8
	localSize := [3]uint16{uint16(localSizeX), uint16(localSizeY), 1}
	globalSizeX := uint32(((b.NP-1)/localSizeX + 1) * localSizeX)
	globalSizeY := uint32(((b.NQ-1)/localSizeY + 1) * localSizeY)
	globalSize := [3]uint32{globalSizeX, globalSizeY, 1}

	kernel1Arg := Kernel1Args{
		int32(b.NR),
		int32(b.NQ),
		int32(b.NP),
		0,
		b.dA,
		b.dC4,
		b.dSum,
		1,
		0,
		0, 0, 0,
	}
	b.driver.LaunchKernel(b.context, b.kernel1,
		globalSize, localSize, &kernel1Arg)

	globalSizeX = uint32(((b.NP-1)/localSizeX + 1) * localSizeX)
	globalSizeY = uint32(((b.NQ-1)/localSizeY + 1) * localSizeY)
	globalSize = [3]uint32{globalSizeX, globalSizeY, 1}

	kernel2Arg := Kernel2Args{
		int32(b.NR),
		int32(b.NQ),
		int32(b.NP),
		0,
		b.dA,
		b.dC4,
		b.dSum,
		1,
		0,
		0, 0, 0,
	}
	b.driver.LaunchKernel(b.context, b.kernel2,
		globalSize, localSize, &kernel2Arg)

	b.driver.MemCopyD2H(b.context, b.sum, b.dSum)
}

// Verify verifies
func (b *Benchmark) Verify() {
	b.cpuDoitgen()

	for r := 0; r < b.NR; r++ {
		for q := 0; q < b.NQ; q++ {
			for p := 0; p < b.NP; p++ {
				if b.cpuSum[r*b.NQ*b.NP+q*b.NP+p] != b.sum[r*b.NQ*b.NP+q*b.NP+p] {
					log.Panicf("Mismatch at %d, expected %f, but get %f",
						r*b.NQ*b.NP+q*b.NP+p,
						b.cpuSum[r*b.NQ*b.NP+q*b.NP+p],
						b.sum[r*b.NQ*b.NP+q*b.NP+p])
				}
			}
		}
	}

	log.Printf("Passed!\n")
}

func (b *Benchmark) cpuDoitgen() {
	b.cpuSum = make([]float32, b.NR*b.NP*b.NQ)

	for r := 0; r < b.NR; r++ {
		for q := 0; q < b.NQ; q++ {
			for p := 0; p < b.NP; p++ {
				b.cpuSum[r*b.NQ*b.NP+q*b.NP+p] = 0
				for s := 0; s < b.NP; s++ {
					b.cpuSum[r*b.NQ*b.NP+q*b.NP+p] =
						b.cpuSum[r*b.NQ*b.NP+q*b.NP+p] +
							b.a[r*b.NQ*b.NP+q*b.NP+s] + b.c4[s*b.NP+p]
					// fix please
				}
			}
			for p := 0; p < b.NR; p++ {
				b.a[r*b.NQ*b.NP+q*b.NP+p] = b.cpuSum[r*b.NQ*b.NP+q*b.NP+p]
			}
		}
	}
}

// SetLASPMemoryAlloc use Unified Memory
func (b *Benchmark) SetLASPMemoryAlloc() {
	b.useLASPMemoryAlloc = true
}
