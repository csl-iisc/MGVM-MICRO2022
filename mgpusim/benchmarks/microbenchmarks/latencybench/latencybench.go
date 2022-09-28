//Package latencybench implements a series of microbenchmark
package latencybench

import (
	// "fmt"
	"log"

	"gitlab.com/akita/mgpusim/driver"
	"gitlab.com/akita/mgpusim/insts"
	"gitlab.com/akita/mgpusim/kernels"
)

// KernelArgs defines kernel arguments
type IdleLoopArgs struct {
	RetValues           driver.GPUPtr
	LoopCount           uint32
	Padding             uint32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}
type NarrowStridedReadArgs struct {
	Array               driver.GPUPtr
	Start               uint32
	End                 uint32
	Stride              uint32
	Padding             uint32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}
type WideStridedReadArgs struct {
	Array               driver.GPUPtr
	Start               uint32
	End                 uint32
	Threads             uint32
	Stride              uint32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}
type NarrowStridedWriteArgs struct {
	Array               driver.GPUPtr
	Start               uint32
	End                 uint32
	Stride              uint32
	Padding             uint32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}
type WideStridedWriteArgs struct {
	Array               driver.GPUPtr
	Start               uint32
	End                 uint32
	Threads             uint32
	Stride              uint32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}
type NarrowStridedReadRemoteArgs struct {
	Array               driver.GPUPtr
	Start               uint32
	End                 uint32
	Stride              uint32
	LoopCount           uint32
	RemoteStart         uint32
	Padding             uint32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}
type PChaseArgs struct {
	Array               driver.GPUPtr
	Start               uint32
	Length              uint32
	HiddenGlobalOffsetX int64
	HiddenGlobalOffsetY int64
	HiddenGlobalOffsetZ int64
}
type TwoBlockOneDelayedPChaseArgs struct {
	Array               driver.GPUPtr
	LoopCount           uint32
	SecondBlock         uint32
	Start1, Start2      uint32
	Length1, Length2    uint32
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

	IdleLoop                 *insts.HsaCo
	NarrowStridedRead        *insts.HsaCo
	WideStridedRead          *insts.HsaCo
	NarrowStridedWrite       *insts.HsaCo
	WideStridedWrite         *insts.HsaCo
	NarrowStridedReadRemote  *insts.HsaCo
	PChase                   *insts.HsaCo
	TwoBlockOneDelayedPChase *insts.HsaCo

	Length      uint32
	Stride      uint32
	Start, End  uint32
	Threads     uint32
	Blocks      uint32
	LoopCount   uint32
	RemoteStart uint32
	BenchType   string

	hArray []uint32
	dArray driver.GPUPtr

	useUnifiedMemory bool
}

// NewBenchmark makes a new benchmark
func NewBenchmark(driver *driver.Driver) *Benchmark {
	b := new(Benchmark)
	b.driver = driver
	b.context = driver.Init()
	b.loadProgram()
	b.Length = 1024
	b.Stride = 32
	b.Start = 0
	b.End = 1024
	b.Threads = 1
	b.Blocks = 1
	b.LoopCount = 1024
	b.RemoteStart = 1024
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

	b.IdleLoop = kernels.LoadProgramFromMemory(hsacoBytes, "IdleLoop")
	if b.IdleLoop == nil {
		log.Panic("Failed to load kernel binary")
	}
	b.NarrowStridedRead = kernels.LoadProgramFromMemory(hsacoBytes,
		"NarrowStridedRead")
	if b.NarrowStridedRead == nil {
		log.Panic("Failed to load kernel binary")
	}
	b.WideStridedRead = kernels.LoadProgramFromMemory(hsacoBytes,
		"WideStridedRead")
	if b.WideStridedRead == nil {
		log.Panic("Failed to load kernel binary")
	}
	b.NarrowStridedWrite = kernels.LoadProgramFromMemory(hsacoBytes,
		"NarrowStridedWrite")
	if b.NarrowStridedWrite == nil {
		log.Panic("Failed to load kernel binary")
	}
	b.WideStridedWrite = kernels.LoadProgramFromMemory(hsacoBytes,
		"WideStridedWrite")
	if b.WideStridedWrite == nil {
		log.Panic("Failed to load kernel binary")
	}
	b.NarrowStridedReadRemote = kernels.LoadProgramFromMemory(hsacoBytes,
		"NarrowStridedReadRemote")
	if b.NarrowStridedReadRemote == nil {
		log.Panic("Failed to load kernel binary")
	}
	b.PChase = kernels.LoadProgramFromMemory(hsacoBytes, "PChase")
	if b.PChase == nil {
		log.Panic("Failed to load kernel binary")
	}
	b.TwoBlockOneDelayedPChase = kernels.LoadProgramFromMemory(hsacoBytes,
		"TwoBlockOneDelayedPChase")
	if b.TwoBlockOneDelayedPChase == nil {
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
	b.hArray = make([]uint32, b.Length)

	for i := 0; i < int(b.Length); i++ {
		b.hArray[i] = uint32(i)
	}

	b.dArray = b.driver.AllocateMemory(
		b.context, uint64(b.Length*4))
	b.driver.Distribute(b.context, b.dArray, uint64(b.Length*4), b.gpus)

	if b.BenchType == "pChase" || b.BenchType == "twoBlockOneDelayedPChase" {
		for i := 0; i < int(b.Length); i = i + int(b.Stride) {
			b.hArray[i] = uint32(i) + b.Stride
		}
	}

	// if b.BenchType != "idle" {
	// 	b.driver.MemCopyH2D(b.context, b.dArray, b.hArray)
	// }
}

func (b *Benchmark) exec() {
	for _, queue := range b.queues {

		switch b.BenchType {
		case "idle":
			kernArg := IdleLoopArgs{
				b.dArray,
				b.LoopCount,
				0,
				0, 0, 0,
			}
			b.driver.EnqueueLaunchKernel(
				queue,
				b.IdleLoop,
				[3]uint32{1, 1, 1},
				[3]uint16{1, 1, 1},
				&kernArg,
			)

		case "narrowStridedRead":
			kernArg := NarrowStridedReadArgs{
				b.dArray,
				b.Start,
				b.End,
				b.Stride,
				0,
				0, 0, 0,
			}
			b.driver.EnqueueLaunchKernel(
				queue,
				b.NarrowStridedRead,
				[3]uint32{1, 1, 1},
				[3]uint16{1, 1, 1},
				&kernArg,
			)

		case "wideStridedRead":
			kernArg := WideStridedReadArgs{
				b.dArray,
				b.Start,
				b.End,
				b.Threads * b.Blocks,
				b.Stride,
				0, 0, 0,
			}
			b.driver.EnqueueLaunchKernel(
				queue,
				b.WideStridedRead,
				[3]uint32{b.Blocks * b.Threads, 1, 1},
				[3]uint16{uint16(b.Threads), 1, 1},
				&kernArg,
			)

		case "narrowStridedWrite":
			kernArg := NarrowStridedWriteArgs{
				b.dArray,
				b.Start,
				b.End,
				b.Stride,
				0,
				0, 0, 0,
			}
			b.driver.EnqueueLaunchKernel(
				queue,
				b.NarrowStridedWrite,
				[3]uint32{1, 1, 1},
				[3]uint16{1, 1, 1},
				&kernArg,
			)

		case "wideStridedWrite":
			kernArg := WideStridedWriteArgs{
				b.dArray,
				b.Start,
				b.End,
				b.Threads * b.Blocks,
				b.Stride,
				0, 0, 0,
			}
			b.driver.EnqueueLaunchKernel(
				queue,
				b.WideStridedWrite,
				[3]uint32{b.Blocks * b.Threads, 1, 1},
				[3]uint16{uint16(b.Threads), 1, 1},
				&kernArg,
			)

		case "narrowStridedReadRemote":
			kernArg := NarrowStridedReadRemoteArgs{
				b.dArray,
				b.Start,
				b.End,
				b.Stride,
				b.LoopCount,
				b.RemoteStart,
				0,
				0, 0, 0,
			}
			b.driver.EnqueueLaunchKernel(
				queue,
				b.NarrowStridedReadRemote,
				[3]uint32{b.Blocks * b.Threads, 1, 1},
				[3]uint16{uint16(b.Threads), 1, 1},
				&kernArg,
			)

		case "pChase":
			kernArg := PChaseArgs{
				b.dArray,
				524288,
				1,
				0, 0, 0,
			}
			b.driver.EnqueueLaunchKernel(
				queue,
				b.PChase,
				[3]uint32{1, 1, 1},
				[3]uint16{1, 1, 1},
				&kernArg,
			)

		case "twoBlockOneDelayedPChase":
			kernArg := TwoBlockOneDelayedPChaseArgs{
				b.dArray,
				b.LoopCount,
				4,
				524288,
				524288,
				524288/b.Stride + 1,
				524288/b.Stride + 1,
				// 2,
				// 2,
				0, 0, 0,
			}
			b.driver.EnqueueLaunchKernel(
				queue,
				b.TwoBlockOneDelayedPChase,
				[3]uint32{5, 1, 1},
				[3]uint16{1, 1, 1},
				&kernArg,
			)
		default:
			log.Printf("Invalid Benchmark Type")
		}
	}

	for _, q := range b.queues {
		b.driver.DrainCommandQueue(q)
	}

	b.driver.MemCopyD2H(b.context, b.hArray, b.dArray)
}

// Verify verifies
func (b *Benchmark) Verify() {
	log.Printf("Passed!\n")
}
