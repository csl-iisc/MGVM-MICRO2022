// Package pagerank implements the PageRank benchmark form Hetero-Mark.
package pagerank

import (
	"fmt"
	"log"
	"math"

	"gitlab.com/akita/mgpusim/benchmarks/matrix/csr"
	"gitlab.com/akita/mgpusim/driver"
	"gitlab.com/akita/mgpusim/insts"
	"gitlab.com/akita/mgpusim/kernels"
)

// KernelArgs defines kernel arguments
type KernelArgs struct {
	NumRows   uint32
	Padding   uint32
	RowOffset driver.GPUPtr
	Col       driver.GPUPtr
	Val       driver.GPUPtr
	Vals      driver.LocalPtr
	Padding2  uint32
	X         driver.GPUPtr
	Y         driver.GPUPtr
}

// Benchmark defines a benchmark
type Benchmark struct {
	driver  *driver.Driver
	context *driver.Context
	gpus    []int
	queues  []*driver.CommandQueue
	kernel  *insts.HsaCo

	NumNodes       uint32
	NumConnections uint32
	MaxIterations  uint32

	hMatrix         csr.Matrix
	hPageRank       []float32
	verPageRank     []float32
	verPageRankTemp []float32

	dPageRank      driver.GPUPtr
	dPageRankTemp  driver.GPUPtr
	dRowOffsets    driver.GPUPtr
	dColumnNumbers driver.GPUPtr
	dValues        driver.GPUPtr
	dLocalValues   driver.LocalPtr

	useUnifiedMemory      bool
	useLASPMemoryAlloc    bool
	useLASPHSLMemoryAlloc bool
}

// NewBenchmark returns a benchmark
func NewBenchmark(driver *driver.Driver) *Benchmark {
	b := new(Benchmark)
	b.driver = driver
	b.context = driver.Init()
	b.loadProgram()
	return b
}

// SelectGPU select GPU
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

// SetLASPMemoryAlloc use Unified Memory
func (b *Benchmark) SetLASPHSLMemoryAlloc() {
	b.useLASPHSLMemoryAlloc = true
}

func (b *Benchmark) loadProgram() {
	hsacoBytes := _escFSMustByte(false, "/kernels.hsaco")

	b.kernel = kernels.LoadProgramFromMemory(hsacoBytes, "PageRankUpdateGpu")
	if b.kernel == nil {
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
	initData := float32(1.0) / float32(b.NumNodes)
	b.hPageRank = make([]float32, b.NumNodes)
	b.verPageRank = make([]float32, b.NumNodes)
	b.verPageRankTemp = make([]float32, b.NumNodes)
	b.hMatrix = csr.MakeMatrixGenerator(b.NumNodes, b.NumConnections).
		GenerateMatrix()

	for i := uint32(0); i < b.NumNodes; i++ {
		b.hPageRank[i] = initData
		b.verPageRank[i] = initData
	}

	if b.useUnifiedMemory {
		b.dPageRank = b.driver.AllocateUnifiedMemory(
			b.context, uint64(b.NumNodes*4))
		b.dPageRankTemp = b.driver.AllocateUnifiedMemory(
			b.context, uint64(b.NumNodes*4))
		b.dRowOffsets = b.driver.AllocateUnifiedMemory(
			b.context, uint64((b.NumNodes+1)*4))
		b.dColumnNumbers = b.driver.AllocateUnifiedMemory(
			b.context, uint64(b.NumConnections*4))
		b.dValues = b.driver.AllocateUnifiedMemory(
			b.context, uint64(b.NumConnections*4))
	} else if b.useLASPMemoryAlloc {
		b.dPageRank = b.driver.AllocateMemoryLASP(
			b.context, uint64(b.NumNodes*4), "div4")
		b.dPageRankTemp = b.driver.AllocateMemoryLASP(
			b.context, uint64(b.NumNodes*4), "div4")
		b.dRowOffsets = b.driver.AllocateMemoryLASP(
			b.context, uint64((b.NumNodes+1)*4), "div4")
		b.dColumnNumbers = b.driver.AllocateMemoryLASP(
			b.context, uint64(b.NumConnections*4), "div4")
		b.dValues = b.driver.AllocateMemoryLASP(
			b.context, uint64(b.NumConnections*4), "div4")
	} else if b.useLASPHSLMemoryAlloc {
		b.dColumnNumbers = b.driver.AllocateMemoryLASP(
			b.context, uint64(b.NumConnections*4), "div4")
		b.dValues = b.driver.AllocateMemoryLASP(
			b.context, uint64(b.NumConnections*4), "div4")
		b.dPageRank = b.driver.AllocateMemoryLASP(
			b.context, uint64(b.NumNodes*4), "div4")
		b.dPageRankTemp = b.driver.AllocateMemoryLASP(
			b.context, uint64(b.NumNodes*4), "div4")
		b.dRowOffsets = b.driver.AllocateMemoryLASP(
			b.context, uint64((b.NumNodes+1)*4), "div4")
	} else {
		b.dPageRank = b.driver.AllocateMemory(
			b.context, uint64(b.NumNodes*4))
		b.dPageRankTemp = b.driver.AllocateMemory(
			b.context, uint64(b.NumNodes*4))
		b.dRowOffsets = b.driver.AllocateMemory(
			b.context, uint64((b.NumNodes+1)*4))
		b.dColumnNumbers = b.driver.AllocateMemory(
			b.context, uint64(b.NumConnections*4))
		b.dValues = b.driver.AllocateMemory(
			b.context, uint64(b.NumConnections*4))
	}
}

func printMatrix(matrix [][]float32, n uint32) {
	for i := uint32(0); i < n; i++ {
		for j := uint32(0); j < n; j++ {
			fmt.Printf("%f ", matrix[i][j])
		}
		fmt.Printf("\n")
	}
}

func (b *Benchmark) exec() {
	b.driver.MemCopyH2D(b.context, b.dPageRank, b.hPageRank)
	b.driver.MemCopyH2D(b.context, b.dRowOffsets,
		b.hMatrix.RowOffsets)
	b.driver.MemCopyH2D(b.context, b.dColumnNumbers,
		b.hMatrix.ColumnNumbers)
	b.driver.MemCopyH2D(b.context, b.dValues,
		b.hMatrix.Values)

	b.dLocalValues = driver.LocalPtr(256)

	localWorkSize := 64
	i := uint32(0)

	for i = 0; i < b.MaxIterations; i++ {
		var kernArg KernelArgs
		if i%2 == 0 {
			kernArg = KernelArgs{
				NumRows:   b.NumNodes,
				RowOffset: b.dRowOffsets,
				Col:       b.dColumnNumbers,
				Val:       b.dValues,
				Vals:      b.dLocalValues,
				X:         b.dPageRank,
				Y:         b.dPageRankTemp,
			}
		} else {
			kernArg = KernelArgs{
				NumRows:   b.NumNodes,
				RowOffset: b.dRowOffsets,
				Col:       b.dColumnNumbers,
				Val:       b.dValues,
				Vals:      b.dLocalValues,
				X:         b.dPageRankTemp,
				Y:         b.dPageRank,
			}
		}

		b.driver.LaunchKernel(
			b.context,
			b.kernel,
			[3]uint32{b.NumNodes * 64, 1, 1},
			[3]uint16{uint16(localWorkSize), 1, 1},
			&kernArg,
		)
	}

	if i%2 != 0 {
		b.driver.MemCopyD2H(b.context, b.hPageRank, b.dPageRankTemp)
	} else {
		b.driver.MemCopyD2H(b.context, b.hPageRank, b.dPageRank)
	}
}

// Verify verifies
func (b *Benchmark) Verify() {
	var i uint32
	m := b.hMatrix
	for i = 0; i < b.MaxIterations; i++ {
		for i := uint32(0); i < b.NumNodes; i++ {
			newValue := float32(0)
			for j := m.RowOffsets[i]; j < m.RowOffsets[i+1]; j++ {
				newValue += m.Values[j] * b.verPageRank[m.ColumnNumbers[j]]
			}
			b.verPageRankTemp[i] = newValue
		}
		copy(b.verPageRank, b.verPageRankTemp)
	}

	for i := uint32(0); i < b.NumNodes; i++ {
		if math.Abs(float64(b.verPageRank[i]-b.hPageRank[i])) > 1e-5 {
			log.Panicf("Mismatch at %d, expected %f, but get %f\n",
				i, b.verPageRank[i], b.hPageRank[i])
		}
	}

	log.Printf("Passed!\n")
}
