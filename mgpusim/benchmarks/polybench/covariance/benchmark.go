// Package covariance implements the covariance benchmark from Polybench.
package covariance

import (
	"log"
	// "math"
	"math/rand"

	"gitlab.com/akita/mgpusim/driver"
	"gitlab.com/akita/mgpusim/insts"
	"gitlab.com/akita/mgpusim/kernels"
)

// MeanKernelArgs list first set of kernel arguments
type MeanKernelArgs struct {
	Mean                driver.GPUPtr
	Data                driver.GPUPtr
	FloatN              float32
	M                   int32
	N                   int32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}

// ReduceKernelArgs list first set of kernel arguments
type ReduceKernelArgs struct {
	Mean                driver.GPUPtr
	Data                driver.GPUPtr
	M                   int32
	N                   int32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}

// CovarKernelArgs list first set of kernel arguments
type CovarKernelArgs struct {
	Symmat              driver.GPUPtr
	Data                driver.GPUPtr
	M                   int32
	N                   int32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}

// Benchmark defines a benchmark
type Benchmark struct {
	driver        *driver.Driver
	context       *driver.Context
	gpus          []int
	queues        []*driver.CommandQueue
	mean_kernel   *insts.HsaCo
	reduce_kernel *insts.HsaCo
	covar_kernel  *insts.HsaCo

	data, symmat         []float32
	mean                 []float32
	symmat_outputFromGPU []float32
	M, N                 int
	ddata, dsymmat       driver.GPUPtr
	dmean                driver.GPUPtr

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

	b.mean_kernel = kernels.LoadProgramFromMemory(
		hsacoBytes, "mean_kernel")
	if b.mean_kernel == nil {
		log.Panic("Failed to load kernel binary")
	}
	b.reduce_kernel = kernels.LoadProgramFromMemory(
		hsacoBytes, "reduce_kernel")
	if b.reduce_kernel == nil {
		log.Panic("Failed to load kernel binary")
	}
	b.covar_kernel = kernels.LoadProgramFromMemory(
		hsacoBytes, "covar_kernel")
	if b.covar_kernel == nil {
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
	b.data = make([]float32, b.M*b.N)
	b.mean = make([]float32, b.M)
	b.symmat = make([]float32, b.M*b.M)
	b.symmat_outputFromGPU = make([]float32, b.M*b.M)

	for i := 0; i < b.M; i++ {
		for j := 0; j < b.N; j++ {
			b.data[i*b.N+j] = float32(i*j) / float32(b.M)
		}
	}

	if b.useUnifiedMemory {
		b.ddata = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.M*b.N*4))
		b.dmean = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.M*4))
		b.dsymmat = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.M*b.M*4))
	} else {
		b.ddata = b.driver.AllocateMemory(b.context,
			uint64(b.M*b.N*4))
		b.dmean = b.driver.AllocateMemory(b.context,
			uint64(b.M*4))
		b.dsymmat = b.driver.AllocateMemory(b.context,
			uint64(b.M*b.M*4))
	}
}

func (b *Benchmark) exec() {
	b.driver.MemCopyH2D(b.context, b.ddata, b.data)
	// b.driver.MemCopyH2D(b.context, b.dmean, b.mean)

	localSize := [3]uint16{256, 1, 1}
	globalSizeX := uint32(((b.M-1)/256 + 1) * 256)
	globalSize := [3]uint32{globalSizeX, 1, 1}

	meanKernelArg := MeanKernelArgs{
		b.dmean,
		b.ddata,
		float32(b.N),
		int32(b.M),
		int32(b.N),
		0, 0, 0,
	}
	b.driver.LaunchKernel(b.context, b.mean_kernel,
		globalSize, localSize, &meanKernelArg)
	// b.driver.MemCopyD2H(b.context, b.a_debug, b.da)

	localSize = [3]uint16{32, 8, 1}
	globalSizeX = uint32(((b.M-1)/32 + 1) * 32)
	globalSizeY := uint32(((b.N-1)/8 + 1) * 8)
	globalSize = [3]uint32{globalSizeX, globalSizeY, 1}

	reduceKernelArg := ReduceKernelArgs{
		b.dmean,
		b.ddata,
		int32(b.M),
		int32(b.N),
		0, 0, 0,
	}
	b.driver.LaunchKernel(b.context, b.reduce_kernel,
		globalSize, localSize, &reduceKernelArg)
	// b.driver.MemCopyD2H(b.context, b.a_debug, b.da)

	localSize = [3]uint16{256, 1, 1}
	globalSizeX = uint32(((b.M-1)/256 + 1) * 256)
	globalSize = [3]uint32{globalSizeX, 1, 1}

	covarKernelArg := CovarKernelArgs{
		b.dsymmat,
		b.ddata,
		int32(b.M),
		int32(b.N),
		0, 0, 0,
	}
	b.driver.LaunchKernel(b.context, b.covar_kernel,
		globalSize, localSize, &covarKernelArg)
	// b.driver.MemCopyD2H(b.context, b.a_debug, b.da)

	b.driver.MemCopyD2H(b.context, b.symmat_outputFromGPU, b.dsymmat)
}

// Verify verifies
func (b *Benchmark) Verify() {
	b.cpuCovar()

	for i := 0; i < b.M; i++ {
		for j := 0; j < b.N; j++ {
			if b.symmat_outputFromGPU[i*b.M+j] != b.symmat[i*b.M+j] {
				log.Panicf("Mismatch at %d, expected %f, but get %f",
					i*b.M+j,
					b.symmat[i*b.M+j],
					b.symmat_outputFromGPU[i*b.M+j])
			}
		}
	}

	log.Printf("Passed!\n")
}

func (b *Benchmark) meanOfColumnVectors() {
	for j := 0; j < b.M; j++ {
		b.mean[j] = 0.0
		for i := 0; i < b.N; i++ {
			b.mean[j] += b.data[i*b.M+j]
		}
		b.mean[j] /= float32(b.N)
	}
}

func (b *Benchmark) centerColumnVectors() {
	for i := 0; i < b.N; i++ {
		for j := 0; j < b.M; j++ {
			b.data[i*b.M+j] -= b.mean[j]
		}
	}
}

func (b *Benchmark) calculateCovariance() {
	for j1 := 0; j1 < b.M; j1++ {
		for j2 := 0; j2 < b.M; j2++ {
			b.symmat[j1*b.M+j2] = 0.0
			for i := 0; i < b.N; i++ {
				b.symmat[j1*b.M+j2] += b.data[i*b.M+j1] * b.data[i*b.M+j2]
			}
			b.symmat[j2*b.M+j1] = b.symmat[j1*b.M+j2]
		}
	}
}

func (b *Benchmark) cpuCovar() {
	b.meanOfColumnVectors()
	b.centerColumnVectors()
	b.calculateCovariance()
}

func (b *Benchmark) SetLASPMemoryAlloc() {
}
