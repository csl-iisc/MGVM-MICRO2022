// Package correlation implements the correlation benchmark from Polybench.
package correlation

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

// StdKernelArgs list first set of kernel arguments
type StdKernelArgs struct {
	Mean                driver.GPUPtr
	Std                 driver.GPUPtr
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
	Std                 driver.GPUPtr
	Data                driver.GPUPtr
	FloatN              float32
	M                   int32
	N                   int32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}

// CorrKernelArgs list first set of kernel arguments
type CorrKernelArgs struct {
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
	std_kernel    *insts.HsaCo
	reduce_kernel *insts.HsaCo
	corr_kernel   *insts.HsaCo

	data, symmat         []float32
	mean, stddev         []float32
	symmat_outputFromGPU []float32
	M, N                 int
	ddata, dsymmat       driver.GPUPtr
	dmean, dstddev       driver.GPUPtr

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
	b.std_kernel = kernels.LoadProgramFromMemory(
		hsacoBytes, "std_kernel")
	if b.std_kernel == nil {
		log.Panic("Failed to load kernel binary")
	}
	b.reduce_kernel = kernels.LoadProgramFromMemory(
		hsacoBytes, "reduce_kernel")
	if b.reduce_kernel == nil {
		log.Panic("Failed to load kernel binary")
	}
	b.corr_kernel = kernels.LoadProgramFromMemory(
		hsacoBytes, "corr_kernel")
	if b.corr_kernel == nil {
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

	for i := 0; i < b.M; i++ {
		for j := 0; j < b.N; j++ {
			b.data[i*b.N+j] = float32(i*j+1) / float32(b.M)
		}
	}

	if b.useUnifiedMemory {
		b.ddata = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.M*b.N*4))
		b.dmean = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.M*4))
		b.dstddev = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.M*4))
		b.dsymmat = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.M*b.N*4))
	} else {
		b.ddata = b.driver.AllocateMemory(b.context,
			uint64(b.M*b.N*4))
		b.dmean = b.driver.AllocateMemory(b.context,
			uint64(b.M*4))
		b.dstddev = b.driver.AllocateMemory(b.context,
			uint64(b.M*4))
		b.dsymmat = b.driver.AllocateMemory(b.context,
			uint64(b.M*b.N*4))
	}
}

func (b *Benchmark) exec() {
	b.driver.MemCopyH2D(b.context, b.ddata, b.data)
	// b.driver.MemCopyH2D(b.context, b.dmean, b.mean)
	// b.driver.MemCopyH2D(b.context, b.dstddev, b.stddev)

	localSize := [3]uint16{32, 1, 1}
	globalSizeX := uint32(((b.N-1)/32 + 1) * 32)
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

	localSize = [3]uint16{32, 1, 1}
	globalSizeX = uint32(((b.N-1)/32 + 1) * 32)
	globalSize = [3]uint32{globalSizeX, 1, 1}

	stdKernelArg := StdKernelArgs{
		b.dmean,
		b.dstddev,
		b.ddata,
		float32(b.N),
		int32(b.M),
		int32(b.N),
		0, 0, 0,
	}
	b.driver.LaunchKernel(b.context, b.std_kernel,
		globalSize, localSize, &stdKernelArg)
	// b.driver.MemCopyD2H(b.context, b.a_debug, b.da)

	localSize = [3]uint16{32, 1, 1}
	globalSizeX = uint32(((b.N-1)/32 + 1) * 32)
	globalSize = [3]uint32{globalSizeX, 1, 1}

	reduceKernelArg := ReduceKernelArgs{
		b.dmean,
		b.dstddev,
		b.ddata,
		float32(b.N),
		int32(b.M),
		int32(b.N),
		0, 0, 0,
	}
	b.driver.LaunchKernel(b.context, b.reduce_kernel,
		globalSize, localSize, &reduceKernelArg)
	// b.driver.MemCopyD2H(b.context, b.a_debug, b.da)

	localSize = [3]uint16{32, 1, 1}
	globalSizeX = uint32(((b.N-1)/32 + 1) * 32)
	globalSize = [3]uint32{globalSizeX, 1, 1}

	corrKernelArg := CorrKernelArgs{
		b.dsymmat,
		b.ddata,
		int32(b.M),
		int32(b.N),
		0, 0, 0,
	}
	b.driver.LaunchKernel(b.context, b.corr_kernel,
		globalSize, localSize, &corrKernelArg)
	// b.driver.MemCopyD2H(b.context, b.a_debug, b.da)

	b.driver.MemCopyD2H(b.context, b.symmat_outputFromGPU, b.dsymmat)
}

// Verify verifies
func (b *Benchmark) Verify() {
	// b.cpulu()

	// for i := 0; i < b.N; i++ {
	// 	for j := 0; j < b.N; j++ {
	// 		if b.a_outputFromGPU[i*b.N+j] != b.a[i*b.N+j] {
	// 			log.Panicf("Mismatch at %d, expected %f, but get %f",
	// 				i*b.N+j,
	// 				b.a[i*b.N+j],
	// 				b.a_outputFromGPU[i*b.N+j])
	// 		}
	// 	}
	// }

	log.Printf("Failed!\n")
}

func (b *Benchmark) SetLASPMemoryAlloc() {
}
