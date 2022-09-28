//Package addtoarray implements a simple microbenchmark
package addtoarray

import (
	// "fmt"
	"log"
	"math/rand"

	"gitlab.com/akita/mgpusim/driver"
	"gitlab.com/akita/mgpusim/insts"
	"gitlab.com/akita/mgpusim/kernels"
)

// KernelArgs defines kernel arguments
type KernelArgs struct {
	Array               driver.GPUPtr
	Index               uint32
	Padding             uint32
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

	kernel *insts.HsaCo

	index  uint32
	length uint32

	hInputData  []uint32
	hOutputData []uint32
	dArray      driver.GPUPtr

	useUnifiedMemory bool
}

// NewBenchmark makes a new benchmark
func NewBenchmark(driver *driver.Driver) *Benchmark {
	b := new(Benchmark)
	b.driver = driver
	b.context = driver.Init()
	b.loadProgram()
	b.length = 32 * 1
	// b.length = 1
	return b
}

// SelectGPU selects GPU
func (b *Benchmark) SelectGPU(gpus []int) {
	b.gpus = gpus
}

// SetUnifiedMemory use Unified Memory
func (b *Benchmark) SetUnifiedMemory() {
	b.useUnifiedMemory = true
}

func (b *Benchmark) loadProgram() {
	hsacoBytes := _escFSMustByte(false, "/kernels.hsaco")

	b.kernel = kernels.LoadProgramFromMemory(hsacoBytes, "MicroBenchmark")
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
	b.hInputData = make([]uint32, b.length)
	b.hOutputData = make([]uint32, b.length)

	start := rand.Intn(100)
	// fmt.Println(start)
	for i := 0; i < int(b.length); i++ {
		b.hInputData[i] = uint32(start) + uint32(i)
		// b.hInputData[i] = uint32(i)
	}

	b.index = 42

	b.dArray = b.driver.AllocateMemory(
		b.context, uint64(b.length*4))
	b.driver.Distribute(b.context, b.dArray, uint64(b.length*4), b.gpus)

	b.driver.MemCopyH2D(b.context, b.dArray, b.hInputData)
}

func (b *Benchmark) exec() {
	for _, queue := range b.queues {
		kernArg := KernelArgs{
			b.dArray,
			b.index,
			0,
			0, 0, 0,
		}

		b.driver.EnqueueLaunchKernel(
			queue,
			b.kernel,
			[3]uint32{(b.length), 1, 1},
			[3]uint16{uint16(1), 1, 1},
			// [3]uint16{uint16(b.length) + 103, 1, 1},
			&kernArg,
		)
	}

	for _, q := range b.queues {
		b.driver.DrainCommandQueue(q)
	}

	b.driver.MemCopyD2H(b.context, b.hOutputData, b.dArray)
}

// Verify verifies
func (b *Benchmark) Verify() {
	failed := false
	for i := 0; i < int(b.length); i++ {
		actual := b.hOutputData[i]
		expected := b.hInputData[i] + b.index
		if expected != actual {
			log.Printf("mismatch at (%d), expected %d, but get %d\n",
				i, expected, actual)
			failed = true
		}
	}

	if failed {
		panic("failed to verify matrix transpose result")
	}
	log.Printf("Passed!\n")
}
