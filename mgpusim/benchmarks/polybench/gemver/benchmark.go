// Package gemver implements the gemver benchmark from Polybench.
package gemver

import (
	// "fmt"
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
	V1                  driver.GPUPtr
	V2                  driver.GPUPtr
	U1                  driver.GPUPtr
	U2                  driver.GPUPtr
	N                   int32
	Padding             int32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}

// Kernel2Args list first set of kernel arguments
type Kernel2Args struct {
	A    driver.GPUPtr
	X    driver.GPUPtr
	Y    driver.GPUPtr
	Z    driver.GPUPtr
	Beta float32
	// Padding1            int32
	N int32
	// Padding2            int32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}

// Kernel3Args list first set of kernel arguments
type Kernel3Args struct {
	A                   driver.GPUPtr
	X                   driver.GPUPtr
	W                   driver.GPUPtr
	Alpha               float32
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

	a               []float32
	u1, u2, v1, v2  []float32
	w, x, y, z      []float32
	w_outputFromGPU []float32
	a_debug         []float32
	x_debug         []float32

	da                 driver.GPUPtr
	du1, du2, dv1, dv2 driver.GPUPtr
	dw, dx, dy, dz     driver.GPUPtr

	N           int
	Alpha, Beta float32

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
		hsacoBytes, "gemver_kernel1")
	if b.kernel1 == nil {
		log.Panic("Failed to load kernel binary")
	}
	b.kernel2 = kernels.LoadProgramFromMemory(
		hsacoBytes, "gemver_kernel2")
	if b.kernel2 == nil {
		log.Panic("Failed to load kernel binary")
	}
	b.kernel3 = kernels.LoadProgramFromMemory(
		hsacoBytes, "gemver_kernel3")
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
	b.a = make([]float32, b.N*b.N)
	b.u1 = make([]float32, b.N)
	b.u2 = make([]float32, b.N)
	b.v1 = make([]float32, b.N)
	b.v2 = make([]float32, b.N)
	b.w = make([]float32, b.N)
	b.x = make([]float32, b.N)
	b.y = make([]float32, b.N)
	b.z = make([]float32, b.N)
	b.w_outputFromGPU = make([]float32, b.N)
	b.a_debug = make([]float32, b.N*b.N)
	b.x_debug = make([]float32, b.N)

	for i := 0; i < b.N; i++ {
		b.u1[i] = float32(i)
		b.u2[i] = float32(i+1) / float32(b.N) / 2.0
		b.v1[i] = float32(i+1) / float32(b.N) / 4.0
		b.v2[i] = float32(i+1) / float32(b.N) / 6.0
		b.y[i] = float32(i+1) / float32(b.N) / 8.0
		b.z[i] = float32(i+1) / float32(b.N) / 9.0
		b.x[i] = 0.0
		b.w[i] = 0.0

		for j := 0; j < b.N; j++ {
			b.a[i*b.N+j] = float32(i*j) / float32(b.N)
		}
	}

	if b.useUnifiedMemory {
		b.da = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.N*b.N*4))
		b.du1 = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.N*4))
		b.du2 = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.N*4))
		b.dv1 = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.N*4))
		b.dv2 = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.N*4))
		b.dw = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.N*4))
		b.dx = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.N*4))
		b.dy = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.N*4))
		b.dz = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.N*4))
	} else {
		b.da = b.driver.AllocateMemory(b.context,
			uint64(b.N*b.N*4))
		b.du1 = b.driver.AllocateMemory(b.context,
			uint64(b.N*4))
		b.du2 = b.driver.AllocateMemory(b.context,
			uint64(b.N*4))
		b.dv1 = b.driver.AllocateMemory(b.context,
			uint64(b.N*4))
		b.dv2 = b.driver.AllocateMemory(b.context,
			uint64(b.N*4))
		b.dw = b.driver.AllocateMemory(b.context,
			uint64(b.N*4))
		b.dx = b.driver.AllocateMemory(b.context,
			uint64(b.N*4))
		b.dy = b.driver.AllocateMemory(b.context,
			uint64(b.N*4))
		b.dz = b.driver.AllocateMemory(b.context,
			uint64(b.N*4))
	}
}

func (b *Benchmark) exec() {
	b.driver.MemCopyH2D(b.context, b.da, b.a)
	b.driver.MemCopyH2D(b.context, b.du1, b.u1)
	b.driver.MemCopyH2D(b.context, b.du2, b.u2)
	b.driver.MemCopyH2D(b.context, b.dv1, b.v1)
	b.driver.MemCopyH2D(b.context, b.dv2, b.v2)
	b.driver.MemCopyH2D(b.context, b.dx, b.x)
	b.driver.MemCopyH2D(b.context, b.dy, b.y)
	b.driver.MemCopyH2D(b.context, b.dz, b.z)
	b.driver.MemCopyH2D(b.context, b.dw, b.w)

	localSize := [3]uint16{32, 8, 1}
	globalSizeX := uint32(((b.N-1)/32 + 1) * 32)
	globalSizeY := uint32(((b.N-1)/8 + 1) * 8)
	globalSize := [3]uint32{globalSizeX, globalSizeY, 1}

	kernel1Arg := Kernel1Args{
		b.da,
		b.dv1,
		b.dv2,
		b.du1,
		b.du2,
		int32(b.N),
		0,
		0, 0, 0,
	}
	b.driver.LaunchKernel(b.context, b.kernel1,
		globalSize, localSize, &kernel1Arg)
	b.driver.MemCopyD2H(b.context, b.a_debug, b.da)

	localSize = [3]uint16{256, 1, 1}
	globalSizeX = uint32(((b.N-1)/256 + 1) * 256)
	globalSizeY = uint32(b.N)
	globalSize = [3]uint32{globalSizeX, globalSizeY, 1}

	kernel2Arg := Kernel2Args{
		b.da,
		b.dx,
		b.dy,
		b.dz,
		b.Beta,
		// 0,
		int32(b.N),
		// 0,
		0, 0, 0,
	}
	b.driver.LaunchKernel(b.context, b.kernel2,
		globalSize, localSize, &kernel2Arg)
	b.driver.MemCopyD2H(b.context, b.x_debug, b.dx)

	localSize = [3]uint16{256, 1, 1}
	globalSizeX = uint32(((b.N-1)/256 + 1) * 256)
	globalSizeY = uint32(b.N)
	globalSize = [3]uint32{globalSizeX, globalSizeY, 1}

	kernel3Arg := Kernel3Args{
		b.da,
		b.dx,
		b.dw,
		b.Alpha,
		int32(b.N),
		0, 0, 0,
	}
	b.driver.LaunchKernel(b.context, b.kernel3,
		globalSize, localSize, &kernel3Arg)

	b.driver.MemCopyD2H(b.context, b.w_outputFromGPU, b.dw)
}

// Verify verifies
func (b *Benchmark) Verify() {

	b.cpuGemver()
	// fmt.Println(b.Alpha)
	for i := 0; i < b.N; i++ {
		for j := 0; j < b.N; j++ {
			// fmt.Println(i, j, b.a_debug[i*b.N+j], b.a[i*b.N+j])
			if math.Abs(float64(b.a_debug[i*b.N+j]-b.a[i*b.N+j])) > 0.001 {
				log.Panicf("Mismatch at %d, %d, expected %f, but get %f",
					i, j,
					b.a[i*b.N+j],
					b.a_debug[i*b.N+j])
			}
		}
	}

	for i := 0; i < b.N; i++ {
		// for j := 0; j < b.N; j++ {
		// fmt.Println(i, b.x_debug[i], b.x[i])
		if math.Abs(float64(b.x_debug[i]-b.x[i])) > 0.001 {
			log.Panicf("Mismatch at %d, expected %f, but get %f",
				i,
				b.x[i],
				b.x_debug[i])
		}
		// }
	}

	for i := 0; i < b.N; i++ {
		// for j := 0; j < b.N; j++ {
		// fmt.Println(i, b.w_outputFromGPU[i], b.w[i])
		if math.Abs(float64(b.w_outputFromGPU[i]-b.w[i])) > 1 {
			log.Panicf("Mismatch at %d, expected %f, but get %f",
				i,
				b.w[i],
				b.w_outputFromGPU[i])
			// }
		}
	}

	log.Printf("Passed!\n")
}

func (b *Benchmark) cpuGemver() {
	for i := 0; i < b.N; i++ {
		for j := 0; j < b.N; j++ {
			b.a[i*b.N+j] = b.a[i*b.N+j] + b.u1[i]*b.v1[j] + b.u2[i]*b.v2[j]
		}
	}
	for i := 0; i < b.N; i++ {
		for j := 0; j < b.N; j++ {
			b.x[i] = b.x[i] + b.Beta*b.a[j*b.N+i]*b.y[j]
		}
		b.x[i] += b.z[i]
	}
	// for i := 0; i < b.N; i++ {
	// 	b.x[i] = b.x[i] + b.z[i]
	// }
	for i := 0; i < b.N; i++ {
		for j := 0; j < b.N; j++ {
			b.w[i] += b.Alpha * b.a[i*b.N+j] * b.x[j]
		}
	}
}

func (b *Benchmark) SetLASPMemoryAlloc() {
}
