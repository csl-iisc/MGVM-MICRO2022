package sssp

import (
	"gitlab.com/akita/mgpusim/benchmarks/matrix/csr"
	"gitlab.com/akita/mgpusim/driver"
	"gitlab.com/akita/mgpusim/insts"
	"gitlab.com/akita/mgpusim/kernels"

	"log"
)

// Kernel1Args list first set of kernel arguments
type Kernel1Args struct {
	VectorD1            driver.GPUPtr
	VectorD2            driver.GPUPtr
	SourceVertex        int32
	NumNodes            int32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}

// Kernel2Args list first set of kernel arguments
type Kernel2Args struct {
	Padding             int32
	NumNodes            int32
	RowD                driver.GPUPtr
	ColD                driver.GPUPtr
	DataD               driver.GPUPtr
	VectorD1            driver.GPUPtr
	VectorD2            driver.GPUPtr
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}

// Kernel3Args list first set of kernel arguments
type Kernel3Args struct {
	VectorD1            driver.GPUPtr
	VectorD2            driver.GPUPtr
	NumNodes            int32
	Padding             int32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}

// Kernel4Args list first set of kernel arguments
type Kernel4Args struct {
	VectorD1            driver.GPUPtr
	VectorD2            driver.GPUPtr
	StopD               driver.GPUPtr
	NumNodes            int32
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
	kernel4 *insts.HsaCo

	row, col, data           []float32
	stop                     []int32
	vector1, vector2         []float32
	NumNodes, NumItems       int
	rowd, cold, datad, stopd driver.GPUPtr
	vectord1, vectord2       driver.GPUPtr
	matrix                   csr.Matrix

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
		hsacoBytes, "vector_init")
	if b.kernel1 == nil {
		log.Panic("Failed to load kernel binary")
	}
	b.kernel2 = kernels.LoadProgramFromMemory(
		hsacoBytes, "spmv_min_dot_plus_kernel")
	if b.kernel2 == nil {
		log.Panic("Failed to load kernel binary")
	}
	b.kernel3 = kernels.LoadProgramFromMemory(
		hsacoBytes, "vector_assign")
	if b.kernel3 == nil {
		log.Panic("Failed to load kernel binary")
	}
	b.kernel4 = kernels.LoadProgramFromMemory(
		hsacoBytes, "vector_diff")
	if b.kernel4 == nil {
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
	b.matrix = csr.
		MakeMatrixGenerator(uint32(b.NumNodes), uint32(b.NumItems)).
		GenerateMatrix()

	if b.useUnifiedMemory {
		b.rowd = b.driver.AllocateUnifiedMemory(b.context,
			uint64((b.NumNodes+1)*4))
		b.cold = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NumItems*4))
		b.datad = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NumItems*4))
		b.stopd = b.driver.AllocateUnifiedMemory(b.context,
			uint64(4))
		b.vectord1 = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NumNodes*4))
		b.vectord2 = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NumNodes*4))
	} else if b.useLASPMemoryAlloc {
		b.rowd = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NumNodes*4), "div4")
		b.cold = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NumItems*4), "div4")
		b.datad = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NumItems*4), "div4")
		b.stopd = b.driver.AllocateMemoryLASP(b.context,
			uint64(4), "div4")
		b.vectord1 = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NumNodes*4), "div4")
		b.vectord2 = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NumNodes*4), "div4")
	} else {
		b.rowd = b.driver.AllocateMemory(b.context,
			uint64(b.NumNodes*4))
		b.cold = b.driver.AllocateMemory(b.context,
			uint64(b.NumItems*4))
		b.datad = b.driver.AllocateMemory(b.context,
			uint64(b.NumItems*4))
		b.stopd = b.driver.AllocateMemory(b.context,
			uint64(4))
		b.vectord1 = b.driver.AllocateMemory(b.context,
			uint64(b.NumNodes*4))
		b.vectord2 = b.driver.AllocateMemory(b.context,
			uint64(b.NumNodes*4))
	}

	b.stop = make([]int32, 1)
}

func (b *Benchmark) exec() {
	b.driver.MemCopyH2D(b.context, b.rowd, b.matrix.RowOffsets)
	b.driver.MemCopyH2D(b.context, b.cold, b.matrix.ColumnNumbers)
	b.driver.MemCopyH2D(b.context, b.datad, b.matrix.Values)

	blockSize := int32(256) // BLOCK_SIZE

	args := Kernel1Args{
		VectorD1:            b.vectord1,
		VectorD2:            b.vectord2,
		SourceVertex:        0,
		NumNodes:            int32(b.NumNodes),
		HiddenGlobalOffsetX: 0,
		HiddenGlobalOffsetY: 0,
		HiddenGlobalOffsetZ: 0,
	}

	globalSize := [3]uint32{uint32(b.NumNodes), 1, 1}
	localSize := [3]uint16{uint16(blockSize), 1, 1}

	b.driver.LaunchKernel(b.context,
		b.kernel1,
		globalSize, localSize,
		&args,
	)

	// for i := 0; i < b.NumNodes; i++ {
	for i := 0; i < 1; i++ {
		args2 := Kernel2Args{
			NumNodes:            int32(b.NumNodes),
			Padding:             0,
			RowD:                b.rowd,
			ColD:                b.cold,
			DataD:               b.datad,
			VectorD1:            b.vectord1,
			VectorD2:            b.vectord2,
			HiddenGlobalOffsetX: 0,
			HiddenGlobalOffsetY: 0,
			HiddenGlobalOffsetZ: 0,
		}

		args3 := Kernel3Args{
			VectorD1:            b.vectord1,
			VectorD2:            b.vectord2,
			NumNodes:            int32(b.NumNodes),
			Padding:             0,
			HiddenGlobalOffsetX: 0,
			HiddenGlobalOffsetY: 0,
			HiddenGlobalOffsetZ: 0,
		}

		args4 := Kernel4Args{
			VectorD1:            b.vectord1,
			VectorD2:            b.vectord2,
			StopD:               b.stopd,
			NumNodes:            int32(b.NumNodes),
			Padding:             0,
			HiddenGlobalOffsetX: 0,
			HiddenGlobalOffsetY: 0,
			HiddenGlobalOffsetZ: 0,
		}

		b.driver.MemCopyH2D(b.context, b.stopd, b.stop)

		b.driver.LaunchKernel(b.context,
			b.kernel3,
			globalSize, localSize,
			&args3,
		)

		b.driver.LaunchKernel(b.context,
			b.kernel2,
			globalSize, localSize,
			&args2,
		)

		b.driver.LaunchKernel(b.context,
			b.kernel4,
			globalSize, localSize,
			&args4,
		)

	}

}

func (b *Benchmark) Verify() {
	return
}

// SetLASPMemoryAlloc use Unified Memory
func (b *Benchmark) SetLASPMemoryAlloc() {
	b.useLASPMemoryAlloc = true
}
