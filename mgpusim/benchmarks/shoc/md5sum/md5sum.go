// Package md5sum implements the md5sum benchmark from the SHOC suite.
package md5sum

import (
	"log"
	"math/rand"

	"gitlab.com/akita/mgpusim/driver"
	"gitlab.com/akita/mgpusim/insts"
	"gitlab.com/akita/mgpusim/kernels"
)

// MD5SumKernelArgs defines kernel arguments
type MD5SumKernelArgs struct {
	SearchDigest0       uint32
	SearchDigest1       uint32
	SearchDigest2       uint32
	SearchDigest3       uint32
	KeySpace            int32
	ByteLength          int32
	ValsPerByte         int32
	FoundIndex          driver.GPUPtr
	FoundKey            driver.GPUPtr
	FoundDigest         driver.GPUPtr
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

	md5sumKernel *insts.HsaCo

	useUnifiedMemory   bool
	useLASPMemoryAlloc bool
}

// NewBenchmark returns a benchmark
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

// SetLASPMemoryAlloc use Unified Memory
func (b *Benchmark) SetLASPMemoryAlloc() {
	b.useLASPMemoryAlloc = true
}

func (b *Benchmark) loadProgram() {
	hsacoBytes := _escFSMustByte(false, "/kernels.hsaco")

	b.md5sumKernel = kernels.LoadProgramFromMemory(
		hsacoBytes, "FindKeyWithDigest_Kernel")
	if b.md5sumKernel == nil {
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
	b.numPaddedCols = ((b.NumCols-1)/b.pad + 1) * b.pad
	b.paddedDataSize = b.NumRows * b.numPaddedCols

	b.hInput = make([]float32, b.paddedDataSize)
	b.hOutput = make([]float32, b.paddedDataSize)
	for i := 0; i < b.paddedDataSize; i++ {
		// b.hInput[i] = float32(i)
		b.hInput[i] = 1
	}

	if b.useUnifiedMemory {
		b.dData1 = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.paddedDataSize*4))
		b.dData2 = b.driver.AllocateUnifiedMemory(b.context,
			uint64(b.paddedDataSize*4))
	} else if b.useLASPMemoryAlloc {
		b.dData1 = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.paddedDataSize*4), "test")
		b.dData2 = b.driver.AllocateMemoryLASP(b.context,
			uint64(b.paddedDataSize*4), "test")
	} else {
		b.dData1 = b.driver.AllocateMemory(b.context,
			uint64(b.paddedDataSize*4))
		b.dData2 = b.driver.AllocateMemory(b.context,
			uint64(b.paddedDataSize*4))
	}

	b.currData = &b.dData1
	b.newData = &b.dData2
}

func (b *Benchmark) exec() {
	b.driver.MemCopyH2D(b.context, *b.currData, b.hInput)
	b.driver.MemCopyH2D(b.context, *b.newData, b.hInput)

	for i := 0; i < b.NumIteration; i++ {
		ldsSize := (b.localRows + 2) * (b.localCols + 2) * 4

		args := StencilKernelArgs{
			Data:                *b.currData,
			NewData:             *b.newData,
			Alignment:           16,
			WCenter:             b.wCenter,
			WCardinal:           b.wCardinal,
			WDiagonal:           b.wDiagonal,
			Sh:                  driver.LocalPtr(ldsSize),
			Padding:             0,
			HiddenGlobalOffsetX: 0,
			HiddenGlobalOffsetY: 0,
			HiddenGlobalOffsetZ: 0,
		}

		globalSize := [3]uint32{
			uint32((b.NumRows - 2) / b.localRows),
			uint32(b.NumCols - 2),
			1,
		}
		localSize := [3]uint16{1, uint16(b.localCols), 1}
		b.driver.LaunchKernel(b.context,
			b.stencilKernel,
			globalSize, localSize,
			&args,
		)

		b.currData, b.newData = b.newData, b.currData
	}

	b.driver.MemCopyD2H(b.context, b.hOutput, *b.currData)
}

// Verify verfies
func (b *Benchmark) Verify() {
	cpuOutput := b.cpuStencil2D()

	mismatch := false
	for x := 0; x < b.NumRows; x++ {
		for y := 0; y < b.NumCols; y++ {
			index := x*b.numPaddedCols + y
			if b.hOutput[index] != cpuOutput[index] {
				mismatch = true
				log.Printf("not match at (%d,%d), expected %f to equal %f\n",
					x, y,
					b.hOutput[index], cpuOutput[index])
			}
		}
	}

	if mismatch {
		panic("Mismatch!\n")
	}
	log.Printf("Passed!\n")
}

func (b *Benchmark) cpuStencil2D() []float32 {
	cpuOutput := make([]float32, b.paddedDataSize)
	for x := 0; x < b.NumRows; x++ {
		for y := 0; y < b.NumCols; y++ {
			cpuOutput[x*b.numPaddedCols+y] =
				b.hInput[x*b.numPaddedCols+y]
		}
	}

	for i := 0; i < b.NumIteration; i++ {
		for x := 0; x < b.NumRows; x++ {
			for y := 0; y < b.NumCols; y++ {
				if x == 0 || y == 0 ||
					x == b.NumRows-1 || y == b.NumCols-1 {
					continue
				}

				center := cpuOutput[x*b.numPaddedCols+y]
				cardinal := cpuOutput[(x-1)*b.numPaddedCols+y] +
					cpuOutput[(x+1)*b.numPaddedCols+y] +
					cpuOutput[x*b.numPaddedCols+(y+1)] +
					cpuOutput[x*b.numPaddedCols+(y-1)]
				diagonal := cpuOutput[(x-1)*b.numPaddedCols+(y+1)] +
					cpuOutput[(x+1)*b.numPaddedCols+(y-1)] +
					cpuOutput[(x+1)*b.numPaddedCols+(y+1)] +
					cpuOutput[(x-1)*b.numPaddedCols+(y-1)]

				out := b.wCenter*center +
					b.wCardinal*cardinal +
					b.wDiagonal*diagonal

				cpuOutput[x*b.numPaddedCols+y] = out
			}
		}
	}

	return cpuOutput
}
