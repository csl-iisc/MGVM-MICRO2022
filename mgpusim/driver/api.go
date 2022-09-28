package driver

import (
	"fmt"
	"log"
	"math"
	"sync/atomic"

	// embed hsaco files
	_ "embed"

	"github.com/rs/xid"
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem/cache"
	"gitlab.com/akita/mem/device"
	"gitlab.com/akita/mgpusim/kernels"
	"gitlab.com/akita/mgpusim/remotetranslation"
	"gitlab.com/akita/util/ca"
)

var nextPID uint64

// Init creates a context to be used in the following API calls.
func (d *Driver) Init() *Context {
	atomic.AddUint64(&nextPID, 1)

	c := &Context{
		pid:          ca.PID(nextPID),
		currentGPUID: 1,
	}

	d.contextMutex.Lock()
	d.contexts = append(d.contexts, c)
	d.contextMutex.Unlock()

	return c
}

// InitWithExistingPID creates a new context that shares the same process ID
// with the given context.
func (d *Driver) InitWithExistingPID(ctx *Context) *Context {
	c := &Context{
		pid:          ctx.pid,
		currentGPUID: 1,
	}

	d.contextMutex.Lock()
	d.contexts = append(d.contexts, c)
	d.contextMutex.Unlock()

	return c
}

// GetNumGPUs return the number of GPUs in the platform
func (d *Driver) GetNumGPUs() int {
	return len(d.GPUs)
}

// SelectGPU requires the driver to perform the following APIs on a selected
// GPU
func (d *Driver) SelectGPU(c *Context, gpuID int) {
	if gpuID >= len(d.devices) {
		log.Panicf("GPU %d is not available", gpuID)
	}
	c.currentGPUID = gpuID
}

// CreateUnifiedGPU can create a virtual GPU that bundles multiple GPUs
// together. It returns the DeviceID of the created unified multi-GPU device.
func (d *Driver) CreateUnifiedGPU(c *Context, gpuIDs []int) int {
	d.mustNotBeAnEmptyList(gpuIDs)
	d.mustBeAllActualGPUs(gpuIDs)

	dev := &device.Device{
		ID:            len(d.devices),
		Type:          device.DeviceTypeUnifiedGPU,
		UnifiedGPUIDs: gpuIDs,
		MemState:      device.NewDeviceMemoryState(d.Log2PageSize, d.memAllocatorType),
	}

	for _, gpuID := range gpuIDs {
		dev.ActualGPUs = append(dev.ActualGPUs, d.devices[gpuID])
	}

	d.devices = append(d.devices, dev)
	d.memAllocator.RegisterDevice(dev)

	return dev.ID
}

func (d *Driver) mustNotBeAnEmptyList(gpuIDs []int) {
	if len(gpuIDs) == 0 {
		panic("must unify at least 1 GPU")
	}
}

func (d *Driver) mustBeAllActualGPUs(gpuIDs []int) {
	for _, gpuID := range gpuIDs {
		dev := d.devices[gpuID]
		if dev.Type != device.DeviceTypeGPU {
			panic("can only unify GPUs")
		}
	}
}

// CreateCommandQueue creates a command queue in the driver
func (d *Driver) CreateCommandQueue(c *Context) *CommandQueue {
	q := new(CommandQueue)
	q.GPUID = c.currentGPUID
	q.Context = c

	c.queueMutex.Lock()
	c.queues = append(c.queues, q)
	c.queueMutex.Unlock()

	return q
}

// DrainCommandQueue will return when there is no command to execute
func (d *Driver) DrainCommandQueue(q *CommandQueue) {
	listener := q.Subscribe()
	defer q.Unsubscribe(listener)

	d.enqueueSignal <- true

	for {
		if q.NumCommand() == 0 {
			return
		}
		listener.Wait()
	}
}

// AllocateMemory allocates a chunk of memory of size byteSize in storage.
// It returns the pointer pointing to the newly allocated memory in the GPU
// memory space.
func (d *Driver) AllocateMemory(
	ctx *Context,
	byteSize uint64,
) GPUPtr {
	ptr := d.memAllocator.AllocateVirtualChunk(ctx.pid, byteSize, ctx.currentGPUID)
	return GPUPtr(ptr)
}

// AllocateMemory allocates a chunk of memory of size byteSize in storage.
// It returns the pointer pointing to the newly allocated memory in the GPU
// memory space.
func (d *Driver) AllocateMemoryLASP(
	ctx *Context,
	byteSize uint64,
	partition string,
) GPUPtr {
	ptr := d.memAllocator.AllocateVirtualChunkLASP(ctx.pid, byteSize, ctx.currentGPUID, partition)
	return GPUPtr(ptr)
}

//AllocateUnifiedMemory allocates a unified memory. Allocation is done on CPU
func (d *Driver) AllocateUnifiedMemory(
	ctx *Context,
	byteSize uint64,
) GPUPtr {
	return GPUPtr(d.memAllocator.AllocateUnified(ctx.pid, byteSize))
}

// Remap keeps the virtual address unchanged and moves the physical address to
// another GPU
func (d *Driver) Remap(ctx *Context, addr, size uint64, deviceID int) {
	d.memAllocator.Remap(ctx.pid, addr, size, deviceID)
}

// Distribute rearranges a consecutive virtual memory space and re-allocate the
// memory on designated GPUs. This function returns the number of bytes
// allocated to each GPU.
func (d *Driver) Distribute(
	ctx *Context,
	addr GPUPtr,
	byteSize uint64,
	gpuIDs []int,
) []uint64 {
	return d.distributor.Distribute(ctx, uint64(addr), byteSize, gpuIDs)
}

// FreeMemory frees the memory pointed by ptr. The pointer must be allocated
// with the function AllocateMemory earlier. Error will be returned if the ptr
// provided is invalid.
func (d *Driver) FreeMemory(ctx *Context, ptr GPUPtr) error {
	// d.memAllocator.Free(ctx.pid, uint64(ptr))
	return nil
}

// EnqueueMemCopyH2D registers a MemCopyH2DCommand in the queue.
func (d *Driver) EnqueueMemCopyH2D(
	queue *CommandQueue,
	dst GPUPtr,
	src interface{},
) {
	d.enqueueFlushBeforeMemCopy(queue)
	cmd := &MemCopyH2DCommand{
		ID:  xid.New().String(),
		Dst: dst,
		Src: src,
	}
	d.Enqueue(queue, cmd)
}

// EnqueueWritePageTableH2D registers a WritePageTableH2DCommand in the queue.
func (d *Driver) EnqueueWritePageTableH2D(
	queue *CommandQueue) {
	d.enqueueFlushBeforeMemCopy(queue)
	cmd := &WritePageTableH2DCommand{
		ID: xid.New().String(),
	}
	d.Enqueue(queue, cmd)
}

// EnqueueMemCopyD2H registers a MemCopyD2HCommand in the queue.
func (d *Driver) EnqueueMemCopyD2H(
	queue *CommandQueue,
	dst interface{},
	src GPUPtr,
) {
	d.enqueueFlushBeforeMemCopy(queue)
	cmd := &MemCopyD2HCommand{
		ID:  xid.New().String(),
		Dst: dst,
		Src: src,
	}
	d.Enqueue(queue, cmd)
}

//go:embed memcopy.hsaco
var kernelBytes []byte

// EnqueueMemCopyD2D registers a MemCopyD2DCommand (LaunchKernelCommand) in the
// queue.
func (d *Driver) EnqueueMemCopyD2D(
	queue *CommandQueue,
	dst GPUPtr,
	src GPUPtr,
	num int,
) {
	co := kernels.LoadProgramFromMemory(
		kernelBytes, "copyKernel")
	if co == nil {
		panic("fail to load copyKernel kernel")
	}
	gridSize := [3]uint32{uint32(math.Ceil(float64(num) / float64(4))), 1, 1}
	//total_bytes / (wgSize * 4). Each thread copies 4 bytes.

	wgSize := [3]uint16{64, 1, 1}
	kernelArgs := KernelMemCopyArgs{src, dst, int64(num)}

	d.EnqueueLaunchKernel(queue, co, gridSize, wgSize, &kernelArgs)
}

func (d *Driver) enqueueFlushBeforeMemCopy(queue *CommandQueue) {
	dirty := queue.Context.l2Dirty
	queue.commandsMutex.Lock()
	for _, cmd := range queue.commands {
		switch cmd.(type) {
		case *LaunchKernelCommand:
			dirty = true
		case *FlushCommand:
			dirty = false
		}
	}
	queue.commandsMutex.Unlock()

	if dirty {
		d.enqueueFinalFlush(queue)
	}
}

func (d *Driver) enqueueFinalFlush(queue *CommandQueue) {
	cmd := &FlushCommand{
		ID: akita.GetIDGenerator().Generate(),
	}
	d.Enqueue(queue, cmd)
}

// MemCopyH2D copies a memory from the host to a GPU device.
func (d *Driver) MemCopyH2D(ctx *Context, dst GPUPtr, src interface{}) {
	queue := d.CreateCommandQueue(ctx)
	d.EnqueueMemCopyH2D(queue, dst, src)
	d.DrainCommandQueue(queue)
}

// WritePageTableH2D copies the page table from the host to a GPU device.
func (d *Driver) WritePageTableH2D(ctx *Context) {
	queue := d.CreateCommandQueue(ctx)
	d.EnqueueWritePageTableH2D(queue)
	d.DrainCommandQueue(queue)
}

// MemCopyD2H copies a memory from a GPU device to the host
func (d *Driver) MemCopyD2H(ctx *Context, dst interface{}, src GPUPtr) {
	queue := d.CreateCommandQueue(ctx)
	d.EnqueueMemCopyD2H(queue, dst, src)
	d.DrainCommandQueue(queue)
}

// MemCopyD2D copies a memory from a GPU device to another GPU device. num is
// the total number of bytes.
func (d *Driver) MemCopyD2D(ctx *Context, dst GPUPtr, src GPUPtr, num int) {
	queue := d.CreateCommandQueue(ctx)
	d.EnqueueMemCopyD2D(queue, dst, src, num)
	d.DrainCommandQueue(queue)
}

// Sets the HSL for the entire GPUs
// Requires two (tested functions).
// Absolutely DO NOT change this when the GPU is running some program!!!
func (d *Driver) SetHSL(pmdUnits int) {

	fmt.Println("Setting the HSL for the entire GPU. ")

	// some parameters hardcoded
	hsl := func(address uint64) uint64 {
		return (address / uint64(pmdUnits)) % 4
	}

	twolvl := func(address uint64) uint64 {
		return (address / uint64(pmdUnits)) % 4
	}

	d.GPUs[0].RemoteAddressTranslationUnits[0].(*remotetranslation.DefaultRTU).RemoteAddressTranslationTable.(*cache.HashLowModuleFinder).Hashfunc = hsl

	// Yes. I am indeed hardcoding all these here.
	// Setting on one TLB per chiplet suffices because they all share the same
	for i := 0; i < 4; i++ {
		d.GPUs[0].L1VTLBs[32*i].LowModuleFinder.(*cache.CustomTwoLevelLowModuleFinder).Hashfunc = twolvl
	}

	// d.GPUs[0].L1VTLBs[32*i].LowModuleFinder.(*cache.CustomTwoLevelLowModuleFinder).PrivateMode = true

}
