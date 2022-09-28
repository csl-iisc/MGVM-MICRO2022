package color

import (
	"fmt"

	"gitlab.com/akita/mgpusim/benchmarks/matrix/csr"
	"gitlab.com/akita/mgpusim/driver"
	"gitlab.com/akita/mgpusim/insts"
	"gitlab.com/akita/mgpusim/kernels"

	"log"
)

// Kernel1Args list first set of kernel arguments
type Kernel1Args struct {
	Row                 driver.GPUPtr
	Col                 driver.GPUPtr
	Node_value          driver.GPUPtr
	Color_array         driver.GPUPtr
	Stop                driver.GPUPtr
	Max_d               driver.GPUPtr
	Color               int32
	Num_nodes           int32
	Num_edges           int32
	Padding             int32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}

// Kernel2Args list first set of kernel arguments
type Kernel2Args struct {
	Node_value          driver.GPUPtr
	Color_array         driver.GPUPtr
	Max_d               driver.GPUPtr
	Color               int32
	Num_nodes           int32
	Num_edges           int32
	Padding             int32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}

// Benchmark defines a benchmark
type Benchmark struct {
	driver                                      *driver.Driver
	context                                     *driver.Context
	gpus                                        []int
	queues                                      []*driver.CommandQueue
	kernel1                                     *insts.HsaCo
	kernel2                                     *insts.HsaCo
	matrix                                      csr.Matrix
	row, col, nodeValue, color, max             []float32
	NumNodes, NumEdges, stop                    int32
	rowD, colD, nodeValueD, stopD, colorD, maxD driver.GPUPtr

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
		hsacoBytes, "color")
	if b.kernel1 == nil {
		log.Panic("Failed to load kernel binary")
	}
	b.kernel2 = kernels.LoadProgramFromMemory(
		hsacoBytes, "color2")
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
	if b.useUnifiedMemory {
		panic("oh no!")
	} else if b.useLASPMemoryAlloc {
		b.matrix = csr.
			MakeMatrixGenerator(uint32(b.NumNodes), uint32(b.NumEdges)).
			GenerateMatrix()
			// TODO: Values need to be on nodes, not edges
		b.colorD = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NumNodes*4), "div4")
		b.maxD = b.driver.AllocateMemoryLASP(b.context,
			uint64((b.NumNodes)*4), "div4") //numNodes+1 was here
		b.rowD = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NumNodes*4), "div4")
		b.colD = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NumEdges*4), "div4")
		b.nodeValueD = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.NumEdges*4), "div4")
		b.stopD = b.driver.AllocateMemoryLASP(b.context,
			uint64(4), "div4")
	} else {
		b.matrix = csr.
			MakeMatrixGenerator(uint32(b.NumNodes), uint32(b.NumEdges)).
			GenerateMatrix()
		b.rowD = b.driver.AllocateMemory(b.context,
			uint64(b.NumNodes*4))

		b.colD = b.driver.AllocateMemory(b.context,
			uint64(b.NumEdges*4))

		b.stopD = b.driver.AllocateMemory(b.context,
			uint64(4))

		b.colorD = b.driver.AllocateMemory(b.context,
			uint64(b.NumNodes*4))
		b.nodeValueD = b.driver.AllocateMemory(b.context,
			uint64(b.NumNodes*4))
		b.maxD = b.driver.AllocateMemory(b.context,
			uint64(b.NumNodes*4)) //numNodes+1 was here

	}
}

func (b *Benchmark) exec() {

	b.color = make([]float32, b.NumNodes)
	var n int32
	for n = 0; n < b.NumNodes; n++ {
		b.color[n] = -1
	}

	b.driver.MemCopyH2D(b.context, b.colorD, b.color)
	b.driver.MemCopyH2D(b.context, b.maxD, b.color)
	b.driver.MemCopyH2D(b.context, b.rowD, b.matrix.RowOffsets)
	b.driver.MemCopyH2D(b.context, b.colD, b.matrix.ColumnNumbers)
	b.driver.MemCopyH2D(b.context, b.nodeValueD, b.matrix.Values)

	blockSize := int32(64) // BLOCKSIZE
	globalSize := b.NumNodes
	if b.NumNodes%blockSize > 0 {
		globalSize = (b.NumNodes/blockSize + 1) * blockSize
	}

	globalWork := [3]uint32{uint32(globalSize), 1, 1}
	localWork := [3]uint16{uint16(blockSize), 1, 1}

	stop := int32(1)
	graphColor := int32(1)

	args1 := Kernel1Args{
		Row:                 b.rowD,
		Col:                 b.colD,
		Node_value:          b.nodeValueD,
		Color_array:         b.colorD,
		Stop:                b.stopD,
		Max_d:               b.maxD,
		Num_nodes:           b.NumNodes,
		Num_edges:           b.NumEdges,
		HiddenGlobalOffsetX: 0,
		HiddenGlobalOffsetY: 0,
		HiddenGlobalOffsetZ: 0,
	}

	args2 := Kernel2Args{
		Node_value:          b.nodeValueD,
		Color_array:         b.colorD,
		Max_d:               b.maxD,
		Num_nodes:           b.NumNodes,
		Num_edges:           b.NumEdges,
		HiddenGlobalOffsetX: 0,
		HiddenGlobalOffsetY: 0,
		HiddenGlobalOffsetZ: 0,
	}

	for stop > 0 {
		stop = 0
		b.driver.MemCopyH2D(b.context, b.stopD, &stop)
		args1.Color = graphColor
		args2.Color = graphColor
		b.driver.LaunchKernel(b.context,
			b.kernel1,
			globalWork, localWork,
			&args1,
		)
		fmt.Println("Launched first kernel")
		b.driver.LaunchKernel(b.context,
			b.kernel1,
			globalWork, localWork,
			&args2,
		)
		b.driver.MemCopyD2H(b.context, &stop, b.stopD)
		graphColor++
	}

}

func (b *Benchmark) Verify() {
	return
}

// SetLASPMemoryAlloc use Unified Memory
func (b *Benchmark) SetLASPMemoryAlloc() {
	b.useLASPMemoryAlloc = true
}
