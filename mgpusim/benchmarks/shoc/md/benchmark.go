// Package md include the benchmark of sparse matrix-vector matiplication.
package md

import (
	"fmt"
	"log"
	"math/rand"

	"gitlab.com/akita/mgpusim/benchmarks/matrix/csr"
	"gitlab.com/akita/mgpusim/driver"
	"gitlab.com/akita/mgpusim/insts"
	"gitlab.com/akita/mgpusim/kernels"
)

// KernelArgs sets up kernel arguments
type KernelArgs struct {
	Force         driver.GPUPtr
	Position      driver.GPUPtr
	Cols          driver.GPUPtr
	RowDelimiters driver.GPUPtr
	Dim           int32
	// VecWidth            int32
	// PartialSums         driver.LocalPtr
	Padding             int32
	Out                 driver.GPUPtr
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}

//Benchmark set up test parameters
type Benchmark struct {
	driver             *driver.Driver
	context            *driver.Context
	gpus               []int
	queues             []*driver.CommandQueue
	useUnifiedMemory   bool
	useLASPMemoryAlloc bool
	mdkernel           *insts.HsaCo

	Dim       int32
	Sparsity  float64
	dValData  driver.GPUPtr
	dVecData  driver.GPUPtr
	dColsData driver.GPUPtr
	dRowDData driver.GPUPtr
	dOutData  driver.GPUPtr
	nItems    int32
	vec       []float32
	out       []float32
	maxval    float32
	matrix    csr.Matrix
}

// NewBenchmark creates a new benchmark
func NewBenchmark(driver *driver.Driver) *Benchmark {
	b := new(Benchmark)
	b.driver = driver
	b.context = driver.Init()
	b.loadProgram()
	b.maxval = 10
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

// SetLASPMemoryAlloc use Unified Memory
func (b *Benchmark) SetLASPMemoryAlloc() {
	b.useLASPMemoryAlloc = true
}

func (b *Benchmark) loadProgram() {
	hsacoBytes := _escFSMustByte(false, "/kernels.hsaco")

	b.mdkernel = kernels.LoadProgramFromMemory(
		hsacoBytes, "compute_lj_force")
	if b.mdkernel == nil {
		log.Panic("Failed to load kernel binary")
	}
}

// Run runs the benchmark
func (b *Benchmark) Run() {
	for _, gpu := range b.gpus {
		b.driver.SelectGPU(b.context, gpu)
		b.queues = append(b.queues, b.driver.CreateCommandQueue(b.context))
	}

	b.initMem()
	b.exec()
}

func (b *Benchmark) initMem() {
	b.nItems = int32(float64(b.Dim) * float64(b.Dim) * b.Sparsity)
	fmt.Printf("Number of non-zero elements %d\n", b.nItems)

	b.matrix = csr.
		MakeMatrixGenerator(uint32(b.Dim), uint32(b.nItems)).
		GenerateMatrix()
	b.vec = make([]float32, b.Dim)
	b.out = make([]float32, b.Dim)

	for j := int32(0); j < b.Dim; j++ {
		b.vec[j] = (rand.Float32() * b.maxval)
	}

	if b.useUnifiedMemory {
		b.dValData = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.nItems*4))
		b.dVecData = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.Dim*4))
		b.dColsData = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.nItems*4))
		b.dRowDData = b.driver.AllocateUnifiedMemory(b.context,
			uint64((b.Dim+1)*4))
		b.dOutData = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.Dim*4))
	} else if b.useLASPMemoryAlloc {
		b.dValData = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.nItems*4), "div4")
		b.dVecData = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.Dim*4), "div4")
		b.dColsData = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.nItems*4), "div4")
		b.dRowDData = b.driver.AllocateMemoryLASP(b.context,
			uint64((b.Dim+1)*4), "div4")
		b.dOutData = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.Dim*4), "div4")
	} else {
		b.dValData = b.driver.AllocateMemory(b.context,
			uint64(b.nItems*4))
		b.dVecData = b.driver.AllocateMemory(b.context,
			uint64(b.Dim*4))
		b.dColsData = b.driver.AllocateMemory(b.context,
			uint64(b.nItems*4))
		b.dRowDData = b.driver.AllocateMemory(b.context,
			uint64((b.Dim+1)*4))
		b.dOutData = b.driver.AllocateMemory(b.context,
			uint64(b.Dim*4))
	}
}

func (b *Benchmark) exec() {
	b.driver.MemCopyH2D(b.context, b.dValData, b.matrix.Values)
	b.driver.MemCopyH2D(b.context, b.dVecData, b.vec)
	b.driver.MemCopyH2D(b.context, b.dColsData, b.matrix.ColumnNumbers)
	b.driver.MemCopyH2D(b.context, b.dRowDData, b.matrix.RowOffsets)
	b.driver.MemCopyH2D(b.context, b.dOutData, b.out)

	//TODO: Review vecWidth, blockSize, and maxwidth
	// vecWidth := int32(64)    // PreferredWorkGroupSizeMultiple
	// maxLocal := int32(64)    // MaxWorkGroupSize
	blockSize := int32(128) // BLOCK_SIZE

	// localWorkSize := vecWidth
	// for ok := true; ok; ok = ((localWorkSize+vecWidth <= maxLocal) && localWorkSize+vecWidth <= blockSize) {
	//	localWorkSize += vecWidth
	// }

	// vectorGlobalWSize := b.Dim * vecWidth // 1 warp per row

	args := KernelArgs{
		Val:           b.dValData,
		Vec:           b.dVecData,
		Cols:          b.dColsData,
		RowDelimiters: b.dRowDData,
		Dim:           b.Dim,
		//VecWidth:            vecWidth,
		//PartialSums:         driver.LocalPtr(blockSize*4), //hardcoded value in spmv.cl
		Padding:             0,
		Out:                 b.dOutData,
		HiddenGlobalOffsetX: 0,
		HiddenGlobalOffsetY: 0,
		HiddenGlobalOffsetZ: 0,
	}

	globalSize := [3]uint32{uint32(b.Dim), 1, 1}
	localSize := [3]uint16{uint16(blockSize), 1, 1}
	//globalSize := [3]uint32{uint32(vectorGlobalWSize), 1, 1}
	//localSize := [3]uint16{uint16(localWorkSize), 1, 1}

	b.driver.LaunchKernel(b.context,
		b.mdkernel,
		globalSize, localSize,
		&args,
	)

	b.driver.MemCopyD2H(b.context, b.out, b.dOutData)
}

// Verify verifies results
func (b *Benchmark) Verify() {
	log.Printf("Passed!\n")
}
