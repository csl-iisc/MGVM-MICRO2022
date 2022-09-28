package srad

import (
	"log"
	"math/rand"

	"gitlab.com/akita/mgpusim/driver"
	"gitlab.com/akita/mgpusim/insts"
	"gitlab.com/akita/mgpusim/kernels"
)

// KernelArgs defines kernel arguments
type KernelArgs struct {
	TableSize           uint64
	Table               driver.GPUPtr
	Starts              driver.GPUPtr
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

	ThreadBlockSize, NThreadBlocks uint64
	NUpdates                       uint64
	TableSize                      uint64
	DevTable, DevStarts            driver.GPUPtr
	HostTable, HostStarts          []uint64

	useUnifiedMemory   bool
	useLASPMemoryAlloc bool
}

// NewBenchmark makes a new benchmark
func NewBenchmark(driver *driver.Driver) *Benchmark {
	b := new(Benchmark)
	b.driver = driver
	b.context = driver.Init()
	b.loadProgram()
	b.ThreadBlockSize = 32 // take from gups_kernel.cl
	b.NThreadBlocks = 128  // take from gups_kernel.cl
	b.TableSize = 1024 * b.ThreadBlockSize * b.NThreadBlocks
	// b.TableSize = 1024 * b.ThreadBlockSize * b.NThreadBlocks
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

// SetLASPMemoryAlloc use Unified Memory
func (b *Benchmark) SetLASPMemoryAlloc() {
	b.useLASPMemoryAlloc = true
}

func (b *Benchmark) loadProgram() {
	hsacoBytes := _escFSMustByte(false, "/kernels.hsaco")

	b.kernel = kernels.LoadProgramFromMemory(hsacoBytes, "RandomAccessUpdate")
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

func HPCC_Starts(n int64) uint64 {

	return uint64(rand.Int63())

	var PERIOD int64
	PERIOD = 1317624576693539401
	var POLY uint64
	POLY = 7
	var i int

	var m2 [64]uint64

	for n < 0 {
		n += PERIOD
	}
	for n > PERIOD {
		n -= PERIOD
	}
	if n == 0 {
		return 0x1
	}

	var temp uint64
	temp = 1
	for i = 0; i < 64; i++ {
		m2[i] = temp
		if temp < 0 {
			temp = (temp << 1) ^ POLY
		} else {
			temp = (temp << 1)
		}
	}

	for i = 62; i >= 0; i-- {
		if ((n >> i) & 1) != 0 {
			break
		}
	}

	var ran uint64
	ran = 0x2
	for i > 0 {
		temp = 0
		for j := 0; j < 64; j++ {
			if ((ran >> j) & 1) != 0 {
				temp ^= m2[j]
			}
		}
		ran = temp
		i -= 1
		if ((n >> i) & 1) != 0 {
			if ran < 0 {
				ran = (ran << 1) ^ POLY
			} else {
				ran = (ran << 1)
			}
		}
	}

	return ran
}

func (b *Benchmark) initMem() {

	size := b.TableSize * 8
	startsSize := b.NThreadBlocks * b.ThreadBlockSize * 8

	b.HostTable = make([]uint64, b.TableSize)
	b.HostStarts = make([]uint64, b.NThreadBlocks*b.ThreadBlockSize)
	b.NUpdates = 1 * b.NThreadBlocks * b.ThreadBlockSize

	for i := uint64(0); i < b.ThreadBlockSize*b.NThreadBlocks; i++ {
		b.HostStarts[i] = HPCC_Starts(int64((b.NUpdates/b.NThreadBlocks/b.ThreadBlockSize)*i)) % b.TableSize
	}

	if b.useUnifiedMemory {
		panic("hello??")
	} else if b.useLASPMemoryAlloc {
		b.DevTable = b.driver.AllocateMemoryLASP(b.context, uint64(size), "div4")
		b.DevStarts = b.driver.AllocateMemoryLASP(b.context, uint64(startsSize), "div4")
	} else {
		panic("hello??")
	}

	b.driver.MemCopyH2D(b.context, b.DevTable, b.HostTable)
	b.driver.MemCopyH2D(b.context, b.DevStarts, b.HostStarts)
}

func (b *Benchmark) exec() {

	for _, queue := range b.queues {

		kernArg := KernelArgs{
			b.TableSize,
			b.DevTable,
			b.DevStarts,
			0, 0, 0,
		}

		b.driver.EnqueueLaunchKernel(
			queue,
			b.kernel,
			[3]uint32{uint32(b.ThreadBlockSize * b.NThreadBlocks), uint32(1), 1},
			[3]uint16{uint16(b.ThreadBlockSize), uint16(1), 1},
			&kernArg,
		)
	}

	for _, q := range b.queues {
		b.driver.DrainCommandQueue(q)
	}

}

// Verify verifies
func (b *Benchmark) Verify() {
	log.Printf("How will it pass if it is not implemented at all?")
}
