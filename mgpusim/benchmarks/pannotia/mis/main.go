package mis

import (
	"gitlab.com/akita/mgpusim/benchmarks/matrix/csr"
	"gitlab.com/akita/mgpusim/driver"
	"gitlab.com/akita/mgpusim/insts"
	"gitlab.com/akita/mgpusim/kernels"

	"log"
)

// Kernel1Args list first set of kernel arguments
type Kernel1Args struct {
	SArrayD             driver.GPUPtr
	CArrayD             driver.GPUPtr
	CArrayUD            driver.GPUPtr
	NumNodes            int32
	NumEdges            int32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}

// Kernel2Args list first set of kernel arguments
type Kernel2Args struct {
	RowD                driver.GPUPtr
	ColD                driver.GPUPtr
	NodeValueD          driver.GPUPtr
	SArrayD             driver.GPUPtr
	CArrayD             driver.GPUPtr
	MinArrayD           driver.GPUPtr
	StopD               driver.GPUPtr
	NumNodes            int32
	NumEdges            int32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}

// Kernel3Args list first set of kernel arguments
type Kernel3Args struct {
	RowD                driver.GPUPtr
	ColD                driver.GPUPtr
	NodeValueD          driver.GPUPtr
	SArrayD             driver.GPUPtr
	CArrayD             driver.GPUPtr
	CArrayUD            driver.GPUPtr
	MinArrayD           driver.GPUPtr
	NumNodes            int32
	NumEdges            int32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}

// Kernel4Args list first set of kernel arguments
type Kernel4Args struct {
	CArrayUD            driver.GPUPtr
	CArrayD             driver.GPUPtr
	NumNodes            int32
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

	NumNodes, NumEdges            int
	row, col, node_value, s_array []float32
	stop                          []int32
	row_d, col_d                  driver.GPUPtr
	node_value_d, s_array_d       driver.GPUPtr
	c_array_d, c_array_u_d        driver.GPUPtr
	min_array_d, stop_d           driver.GPUPtr
	matrix                        csr.Matrix

	useUnifiedMemory      bool
	useLASPMemoryAlloc    bool
	useLASPSHLMemoryAlloc bool
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
		hsacoBytes, "init")
	if b.kernel1 == nil {
		log.Panic("Failed to load kernel binary")
	}
	b.kernel2 = kernels.LoadProgramFromMemory(
		hsacoBytes, "mis1")
	if b.kernel2 == nil {
		log.Panic("Failed to load kernel binary")
	}
	b.kernel3 = kernels.LoadProgramFromMemory(
		hsacoBytes, "mis2")
	if b.kernel3 == nil {
		log.Panic("Failed to load kernel binary")
	}
	b.kernel4 = kernels.LoadProgramFromMemory(
		hsacoBytes, "mis3")
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
		MakeMatrixGenerator(uint32(b.NumNodes), uint32(b.NumEdges)).
		GenerateMatrix()

	if b.useUnifiedMemory {
		b.row_d = b.driver.AllocateUnifiedMemory(b.context,
			uint64((b.NumNodes)*4))
		b.col_d = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NumEdges*4))
		b.stop_d = b.driver.AllocateUnifiedMemory(b.context,
			uint64(4))
		b.min_array_d = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NumNodes*4))
		b.c_array_d = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NumNodes*4))
		b.c_array_u_d = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NumNodes*4))
		b.s_array_d = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NumNodes*4))
		b.node_value_d = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.NumNodes*4))
	} else if b.useLASPMemoryAlloc {
		b.col_d = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NumEdges*4), "div4")
		b.row_d = b.driver.AllocateMemoryLASP(b.context,
			uint64((b.NumNodes)*4), "div4")
		b.min_array_d = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NumNodes*4), "div4")
		b.c_array_d = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NumNodes*4), "div4")
		b.c_array_u_d = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NumNodes*4), "div4")
		b.s_array_d = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NumNodes*4), "div4")
		b.node_value_d = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NumNodes*4), "div4")
		b.stop_d = b.driver.AllocateMemoryLASP(b.context,
			uint64(4), "div4")
	} else if b.useLASPSHLMemoryAlloc {
		b.col_d = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NumEdges*4), "div4")
		b.row_d = b.driver.AllocateMemoryLASP(b.context,
			uint64((b.NumNodes)*4), "div4")
		b.min_array_d = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NumNodes*4), "div4")
		b.c_array_d = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NumNodes*4), "div4")
		b.c_array_u_d = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NumNodes*4), "div4")
		b.s_array_d = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NumNodes*4), "div4")
		b.node_value_d = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NumNodes*4), "div4")
		b.stop_d = b.driver.AllocateMemoryLASP(b.context,
			uint64(4), "div4")
	} else {
		b.row_d = b.driver.AllocateMemory(b.context,
			uint64((b.NumNodes)*4))
		b.col_d = b.driver.AllocateMemory(b.context,
			uint64(b.NumEdges*4))
		b.stop_d = b.driver.AllocateMemory(b.context,
			uint64(4))
		b.min_array_d = b.driver.AllocateMemory(b.context,
			uint64(b.NumNodes*4))
		b.c_array_d = b.driver.AllocateMemory(b.context,
			uint64(b.NumNodes*4))
		b.c_array_u_d = b.driver.AllocateMemory(b.context,
			uint64(b.NumNodes*4))
		b.s_array_d = b.driver.AllocateMemory(b.context,
			uint64(b.NumNodes*4))
		b.node_value_d = b.driver.AllocateMemory(b.context,
			uint64(b.NumNodes*4))
	}

	b.stop = make([]int32, 1)
	b.node_value = make([]float32, b.NumNodes)
	for i := 0; i < b.NumNodes; i++ {
		b.node_value[i] = float32(i)
	}
	// TODO: Set the node value as rand

}

func (b *Benchmark) exec() {
	b.driver.MemCopyH2D(b.context, b.row_d, b.matrix.RowOffsets)
	b.driver.MemCopyH2D(b.context, b.col_d, b.matrix.ColumnNumbers)
	b.driver.MemCopyH2D(b.context, b.node_value_d, b.node_value)

	blockSize := int32(128) // BLOCK_SIZE

	args := Kernel1Args{
		SArrayD:             b.s_array_d,
		CArrayD:             b.c_array_d,
		CArrayUD:            b.c_array_u_d,
		NumNodes:            int32(b.NumNodes),
		NumEdges:            int32(b.NumEdges),
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

	b.stop[0] = 1
	for b.stop[0] != 0 {
		args2 := Kernel2Args{
			RowD:                b.row_d,
			ColD:                b.col_d,
			NodeValueD:          b.node_value_d,
			SArrayD:             b.s_array_d,
			CArrayD:             b.c_array_d,
			MinArrayD:           b.min_array_d,
			StopD:               b.stop_d,
			NumNodes:            int32(b.NumNodes),
			NumEdges:            int32(b.NumEdges),
			HiddenGlobalOffsetX: 0,
			HiddenGlobalOffsetY: 0,
			HiddenGlobalOffsetZ: 0,
		}

		args3 := Kernel3Args{
			RowD:                b.row_d,
			ColD:                b.col_d,
			NodeValueD:          b.node_value_d,
			SArrayD:             b.s_array_d,
			CArrayD:             b.c_array_d,
			CArrayUD:            b.c_array_u_d,
			MinArrayD:           b.min_array_d,
			NumNodes:            int32(b.NumNodes),
			NumEdges:            int32(b.NumEdges),
			HiddenGlobalOffsetX: 0,
			HiddenGlobalOffsetY: 0,
			HiddenGlobalOffsetZ: 0,
		}

		args4 := Kernel4Args{
			CArrayUD:            b.c_array_u_d,
			CArrayD:             b.c_array_d,
			NumNodes:            int32(b.NumNodes),
			HiddenGlobalOffsetX: 0,
			HiddenGlobalOffsetY: 0,
			HiddenGlobalOffsetZ: 0,
		}

		b.driver.LaunchKernel(b.context,
			b.kernel2,
			globalSize, localSize,
			&args2,
		)

		b.driver.LaunchKernel(b.context,
			b.kernel3,
			globalSize, localSize,
			&args3,
		)

		b.driver.LaunchKernel(b.context,
			b.kernel4,
			globalSize, localSize,
			&args4,
		)

		b.driver.MemCopyD2H(b.context, b.stop, b.stop_d)

		b.stop[0] = 0
	}

}

func (b *Benchmark) Verify() {
	return
}

func (b *Benchmark) SetLASPMemoryAlloc() {
	b.useLASPMemoryAlloc = true
}

func (b *Benchmark) SetLASPHSLMemoryAlloc() {
	b.useLASPSHLMemoryAlloc = true
}
