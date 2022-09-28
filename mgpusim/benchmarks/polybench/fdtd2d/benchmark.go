// Package fdtd2d implements the fdtd2d benchmark from Polybench.
package fdtd2d

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
	Fict                driver.GPUPtr
	Ex                  driver.GPUPtr
	Ey                  driver.GPUPtr
	Hz                  driver.GPUPtr
	T                   int32
	NX                  int32
	NY                  int32
	Padding             int32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}

// Kernel2Args list first set of kernel arguments
type Kernel2Args struct {
	Ex                  driver.GPUPtr
	Ey                  driver.GPUPtr
	Hz                  driver.GPUPtr
	NX                  int32
	NY                  int32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}

// Kernel3Args list first set of kernel arguments
type Kernel3Args struct {
	Ex                  driver.GPUPtr
	Ey                  driver.GPUPtr
	Hz                  driver.GPUPtr
	NX                  int32
	NY                  int32
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

	fict, ex, ey, hz                                     []float32
	NX, NY, TMax                                         int
	ey_outputFromGPU, ex_outputFromGPU, hz_outputFromGPU []float32
	dfict, dex, dey, dhz                                 driver.GPUPtr

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
		hsacoBytes, "fdtd_kernel1")
	if b.kernel1 == nil {
		log.Panic("Failed to load kernel binary")
	}
	b.kernel2 = kernels.LoadProgramFromMemory(
		hsacoBytes, "fdtd_kernel2")
	if b.kernel2 == nil {
		log.Panic("Failed to load kernel binary")
	}
	b.kernel3 = kernels.LoadProgramFromMemory(
		hsacoBytes, "fdtd_kernel3")
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
	b.fict = make([]float32, b.TMax)
	b.ex = make([]float32, b.NX*b.NY)
	b.ey = make([]float32, b.NX*b.NY)
	b.hz = make([]float32, b.NX*b.NY)

	b.ey_outputFromGPU = make([]float32, b.NX*b.NY)
	b.ex_outputFromGPU = make([]float32, b.NX*b.NY)
	b.hz_outputFromGPU = make([]float32, b.NX*b.NY)

	for i := 0; i < b.TMax; i++ {
		b.fict[i] = float32(i)
	}
	for i := 0; i < b.NX; i++ {
		for j := 0; j < b.NY; j++ {
			b.ex[i*b.NY+j] = float32(i*(j+1)) / float32(b.NX)
			b.ey[i*b.NY+j] = float32((i-1)*(j+2)+2) / float32(b.NX)
			b.hz[i*b.NY+j] = float32((i-9)*(j+4)+3) / float32(b.NX)
			// fmt.Println(i, j, b.ey[i*b.NY+j])
		}
	}

	if b.useUnifiedMemory {
		b.dfict = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.TMax*4))
		b.dex = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NX*b.NY*4))
		b.dey = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NX*b.NY*4))
		b.dhz = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NX*b.NY*4))
	} else {
		b.dfict = b.driver.AllocateMemory(b.context,
			uint64(b.TMax*4))
		b.dex = b.driver.AllocateMemory(b.context,
			uint64(b.NX*b.NY*4))
		b.dey = b.driver.AllocateMemory(b.context,
			uint64(b.NX*b.NY*4))
		b.dhz = b.driver.AllocateMemory(b.context,
			uint64(b.NX*b.NY*4))
	}
}

func (b *Benchmark) exec() {
	b.driver.MemCopyH2D(b.context, b.dfict, b.fict)
	b.driver.MemCopyH2D(b.context, b.dex, b.ex)
	b.driver.MemCopyH2D(b.context, b.dey, b.ey)
	b.driver.MemCopyH2D(b.context, b.dhz, b.hz)

	localSize := [3]uint16{32, 8, 1}
	//this is NOT a typo
	globalSizeX := uint32(((b.NY-1)/32 + 1) * 32)
	globalSizeY := uint32(((b.NX-1)/8 + 1) * 8)
	globalSize := [3]uint32{globalSizeX, globalSizeY, 1}

	for t := 0; t < b.TMax; t++ {
		kernel1Arg := Kernel1Args{
			b.dfict,
			b.dex,
			b.dey,
			b.dhz,
			int32(t),
			int32(b.NX),
			int32(b.NY),
			0,
			0, 0, 0,
		}
		b.driver.LaunchKernel(b.context, b.kernel1,
			globalSize, localSize, &kernel1Arg)
		b.driver.MemCopyD2H(b.context, b.ey_outputFromGPU, b.dey)

		kernel2Arg := Kernel2Args{
			b.dex,
			b.dey,
			b.dhz,
			int32(b.NX),
			int32(b.NY),
			0, 0, 0,
		}
		b.driver.LaunchKernel(b.context, b.kernel2,
			globalSize, localSize, &kernel2Arg)

		kernel3Arg := Kernel3Args{
			b.dex,
			b.dey,
			b.dhz,
			int32(b.NX),
			int32(b.NY),
			0, 0, 0,
		}
		b.driver.LaunchKernel(b.context, b.kernel3,
			globalSize, localSize, &kernel3Arg)
	}
	// // b.driver.MemCopyD2H(b.context, b.ex_outputFromGPU, b.dhz)
	// // b.driver.MemCopyD2H(b.context, b.ey_outputFromGPU, b.dhz)
	b.driver.MemCopyD2H(b.context, b.hz_outputFromGPU, b.dhz)
}

// Verify verifies
func (b *Benchmark) Verify() {
	b.cpufdtd2d()

	// allow some amount of slack (not 0.001).
	for i := 0; i < b.NX; i++ {
		for j := 0; j < b.NY; j++ {
			if math.Abs(float64(b.ey_outputFromGPU[i*b.NY+j]-b.ey[i*b.NY+j])/float64(b.ey[i*b.NY+j])) > 0.01 {
				fmt.Println("Mismatch at %d, expected %f, but get %f",
					i*b.NY+j,
					b.ey[i*b.NY+j],
					b.ey_outputFromGPU[i*b.NY+j])
			}
		}
	}

	log.Printf("Passed!\n")
}

func (b *Benchmark) cpufdtd2d() {
	for t := 0; t < b.TMax; t++ {

		for j := 0; j < b.NY; j++ {
			b.ey[j] = b.fict[t]
		}
		for i := 1; i < b.NX; i++ {
			for j := 0; j < b.NY; j++ {
				// fmt.Println(b.ey[i*b.NY+j])
				b.ey[i*b.NY+j] = b.ey[i*b.NY+j] - 0.5*(b.hz[i*b.NY+j]-b.hz[(i-1)*b.NY+j])
				// fmt.Println("***", b.ey[i*b.NY+j])
			}
		}
		for i := 0; i < b.NX; i++ {
			for j := 1; j < b.NY; j++ {
				b.ex[i*b.NY+j] = b.ex[i*b.NY+j] - 0.5*(b.hz[i*b.NY+j]-b.hz[i*b.NY+(j-1)])
			}
		}
		for i := 0; i < b.NX-1; i++ {
			for j := 0; j < b.NY-1; j++ {
				b.hz[i*b.NY+j] = b.hz[i*b.NY+j] -
					0.7*(b.ex[i*b.NY+(j+1)]-b.ex[i*b.NY+(j)]+
						b.ey[(i+1)*b.NY+j]-b.ey[i*b.NY+j])
			}
		}

	}
}

func (b *Benchmark) SetLASPMemoryAlloc() {
}
