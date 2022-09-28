// package device provides support for the driver implementation.
package device

import (
	"math"
	"sync"

	"fmt"

	"gitlab.com/akita/util/ca"
)

// A MemoryAllocator can allocate memory on the CPU and GPUs
type MemoryAllocator interface {
	RegisterDevice(device *Device)
	GetDeviceIDByPAddr(pAddr uint64) int
	AllocateVirtualChunk(pid ca.PID, byteSize uint64, deviceID int) uint64
	allocatePageTablePage(pid ca.PID, deviceID int, vAddr, pAddr uint64) uint64
	AllocateUnified(pid ca.PID, byteSize uint64) uint64
	AllocateVirtualChunkLASP(pid ca.PID, byteSize uint64, deviceID int, partition string) uint64
	Free(pid ca.PID, vAddr uint64)
	Remap(pid ca.PID, pageVAddr, byteSize uint64, deviceID int)
	RemovePage(pid ca.PID, vAddr uint64)
	AllocatePageWithGivenVAddr(
		pid ca.PID,
		deviceID int,
		vAddr uint64,
		unified bool,
	) Page
}

// NewMemoryAllocator creates a new memory allocator.
func NewMemoryAllocator(
	pageTable PageTable,
	log2PageSize uint64,
) MemoryAllocator {
	a := &memoryAllocatorImpl{
		pageTable:            pageTable,
		totalStorageByteSize: 1 << log2PageSize, // Starting with a page to avoid 0 address.
		// totalStorageByteSize: 0, // starting page with 0 because, just like that, our wish (Neha's)
		log2PageSize:        log2PageSize,
		processMemoryStates: make(map[ca.PID]*processMemoryState),
		devices:             make(map[int]*Device),
		// memAllocatorType:    memAllocatorType,
	}
	return a
}

type processMemoryState struct {
	pid       ca.PID
	nextVAddr uint64
}

// A memoryAllocatorImpl provides the default implementation for
// memoryAllocator
type memoryAllocatorImpl struct {
	sync.Mutex
	pageTable            PageTable
	log2PageSize         uint64
	processMemoryStates  map[ca.PID]*processMemoryState
	devices              map[int]*Device
	totalStorageByteSize uint64
	// memAllocatorType     AllocatorType
}

func (a *memoryAllocatorImpl) RegisterDevice(device *Device) {
	a.Lock()
	defer a.Unlock()

	state := device.MemState
	fmt.Println("****", a.totalStorageByteSize, float64(a.totalStorageByteSize)/float64((uint64(1)<<30)), uint64(2)*(uint64(1)<<30))
	state.setInitialAddress(a.totalStorageByteSize)

	a.totalStorageByteSize += state.getStorageSize()

	a.devices[device.ID] = device
}

func (a *memoryAllocatorImpl) GetDeviceIDByPAddr(pAddr uint64) int {
	a.Lock()
	defer a.Unlock()

	return a.deviceIDByPAddr(pAddr)
}

func (a *memoryAllocatorImpl) deviceIDByPAddr(pAddr uint64) int {
	for id, dev := range a.devices {
		state := dev.MemState
		if isPAddrOnDevice(pAddr, state) {
			return id
		}
	}

	panic("device not found")
}

func isPAddrOnDevice(
	pAddr uint64,
	state DeviceMemoryState,
) bool {
	return pAddr >= state.getInitialAddress() &&
		pAddr < state.getInitialAddress()+state.getStorageSize()
}

func (a *memoryAllocatorImpl) AllocateVirtualChunk(
	pid ca.PID,
	byteSize uint64,
	deviceID int,
) uint64 {
	a.Lock()
	defer a.Unlock()
	pageSize := uint64(1 << a.log2PageSize)
	numPages := (byteSize-1)/pageSize + 1
	return a.allocatePages(int(numPages), pid, deviceID, false)
}

func (a *memoryAllocatorImpl) AllocateVirtualChunkLASP(
	pid ca.PID, byteSize uint64, deviceID int, partition string,
) uint64 {
	a.Lock()
	defer a.Unlock()
	pageSize := uint64(1 << a.log2PageSize)
	numPages := (byteSize-1)/pageSize + 1
	return a.allocatePagesLASP(int(numPages), pid, deviceID, false, partition)
}

// I only use this function for allocating page table pages
func (a *memoryAllocatorImpl) allocatePageTablePage(
	pid ca.PID,
	deviceID int,
	vAddr, pAddr uint64,
) uint64 {
	// no locks here to prevent deadlocks
	device := a.devices[deviceID]
	return device.allocatePageTablePage(vAddr, pAddr)
}

func (a *memoryAllocatorImpl) AllocateUnified(
	pid ca.PID,
	byteSize uint64,
) uint64 {
	a.Lock()
	defer a.Unlock()

	pageSize := uint64(1 << a.log2PageSize)
	numPages := (byteSize-1)/pageSize + 1
	return a.allocatePages(int(numPages), pid, 1, true)
}

func (a *memoryAllocatorImpl) allocatePages(
	numPages int, pid ca.PID, deviceID int, unified bool,
) (firstPageVAddr uint64) {
	fmt.Println("num pages", numPages)
	pState, found := a.processMemoryStates[pid]
	if !found {
		a.processMemoryStates[pid] = &processMemoryState{
			pid:       pid,
			nextVAddr: uint64(1 << a.log2PageSize),
		}
		pState = a.processMemoryStates[pid]
	}
	device := a.devices[deviceID]

	pageSize := uint64(1 << a.log2PageSize)
	nextVAddr := pState.nextVAddr
	pAddrs := device.allocateMultiplePages(numPages)
	fmt.Println(numPages)
	for i := 0; i < numPages; i++ {
		pAddr := pAddrs[i]
		vAddr := nextVAddr + uint64(i)*pageSize
		page := Page{
			PID:      pid,
			VAddr:    vAddr,
			PAddr:    pAddr,
			PageSize: pageSize,
			Valid:    true,
			Unified:  unified,
			DeviceID: uint64(a.deviceIDByPAddr(pAddr)),
		}
		if page.DeviceID != uint64(deviceID) {
			panic("gpuid != deviceid")
		}
		fmt.Println("data page:", fmt.Sprintf("%x", pAddr), fmt.Sprintf("%x", pAddr-4096), fmt.Sprintf("%x", (((pAddr-4096)>>12)%32)/8))
		a.pageTable.Insert(page)
	}

	pState.nextVAddr += pageSize * uint64(numPages)
	// fmt.Println("#########", nextVAddr, numPages)
	return nextVAddr
}

func (a *memoryAllocatorImpl) findChipletBasedOnPartition(vAddr uint64,
	startAddr uint64, numPages int, partition string) int {
	// pageSize := uint64(1 << a.log2PageSize)
	VPN := (vAddr - startAddr)
	// / pageSize
	if partition == "mod4" {
		return int(VPN % 4)
	}
	if partition == "div4" {
		num_parts := int(math.Ceil(float64(numPages) / 4.0))
		chiplet := int(VPN / uint64(num_parts))
		// fmt.Println(VPN, chiplet)
		return chiplet
	}
	if partition == "test" {
		return 1
	}
	return int(VPN % 4)
}

func (a *memoryAllocatorImpl) Remap(
	pid ca.PID,
	pageVAddr, byteSize uint64,
	deviceID int,
) {
	a.Lock()
	defer a.Unlock()

	pageSize := uint64(1 << a.log2PageSize)
	addr := pageVAddr
	vAddrs := make([]uint64, 0)
	for addr < pageVAddr+byteSize {
		vAddrs = append(vAddrs, addr)
		addr += pageSize
	}

	a.allocateMultiplePagesWithGivenVAddrs(pid, deviceID, vAddrs, false)
}

func (a *memoryAllocatorImpl) RemovePage(pid ca.PID, vAddr uint64) {
	a.Lock()
	defer a.Unlock()

	a.removePage(pid, vAddr)
}

func (a *memoryAllocatorImpl) removePage(pid ca.PID, vAddr uint64) {
	page, ok := a.pageTable.Find(pid, vAddr)

	if !ok {
		panic("page not found")
	}

	deviceID := a.deviceIDByPAddr(page.PAddr)
	dState := a.devices[deviceID].MemState
	dState.addSinglePAddr(page.PAddr)

	a.pageTable.Remove(page.PID, page.VAddr)
}

func (a *memoryAllocatorImpl) AllocatePageWithGivenVAddr(
	pid ca.PID,
	deviceID int,
	vAddr uint64,
	isUnified bool,
) Page {
	a.Lock()
	defer a.Unlock()

	return a.allocatePageWithGivenVAddr(pid, deviceID, vAddr, isUnified)
}

func (a *memoryAllocatorImpl) allocatePageWithGivenVAddr(
	pid ca.PID,
	deviceID int,
	vAddr uint64,
	isUnified bool,
) Page {
	pageSize := uint64(1 << a.log2PageSize)

	device := a.devices[deviceID]
	pAddr := device.allocatePage()

	page := Page{
		PID:      pid,
		VAddr:    vAddr,
		PAddr:    pAddr,
		PageSize: pageSize,
		Valid:    true,
		DeviceID: uint64(deviceID),
		Unified:  isUnified,
	}
	a.pageTable.Update(page)

	return page
}

func (a *memoryAllocatorImpl) allocateMultiplePagesWithGivenVAddrs(
	pid ca.PID,
	deviceID int,
	vAddrs []uint64,
	isUnified bool,
) (pages []Page) {
	pageSize := uint64(1 << a.log2PageSize)

	device := a.devices[deviceID]
	pAddrs := device.allocateMultiplePages(len(vAddrs))

	for i, vAddr := range vAddrs {
		page := Page{
			PID:      pid,
			VAddr:    vAddr,
			PAddr:    pAddrs[i],
			PageSize: pageSize,
			Valid:    true,
			DeviceID: uint64(deviceID),
			Unified:  isUnified,
		}
		a.pageTable.Update(page)
		pages = append(pages, page)
	}

	return pages
}

func (a *memoryAllocatorImpl) Free(pid ca.PID, ptr uint64) {
	a.Lock()
	defer a.Unlock()

	a.removePage(pid, ptr)
}

func (a *memoryAllocatorImpl) allocatePagesLASP(
	numPages int, pid ca.PID, deviceID int, unified bool, partition string,
) (firstPageVAddr uint64) {
	fmt.Println("num pages", numPages)
	pState, found := a.processMemoryStates[pid]
	if !found {
		a.processMemoryStates[pid] = &processMemoryState{
			pid: pid,
			// nextVAddr: uint64(1 << a.log2PageSize),
			nextVAddr: 0,
		}
		pState = a.processMemoryStates[pid]
	}
	device := a.devices[deviceID]

	pageSize := uint64(1 << a.log2PageSize)
	nextVAddr := pState.nextVAddr
	for i := 0; i < numPages; i++ {
		vAddr := nextVAddr + uint64(i)*pageSize
		// maybe replace this to be based on i instead of vAddr if we want data structs to start at first chiplet always
		chiplet := a.findChipletBasedOnPartition(uint64(i), 0, numPages, partition)
		pAddr := device.allocatePageOnChiplet(chiplet)
		page := Page{
			PID:      pid,
			VAddr:    vAddr,
			PAddr:    pAddr,
			PageSize: pageSize,
			Valid:    true,
			Unified:  unified,
			DeviceID: uint64(a.deviceIDByPAddr(pAddr)),
		}
		if page.DeviceID != uint64(deviceID) {
			panic("gpuid != deviceid")
		}
		// fmt.Println("data page:", ((pAddr>>12)%32)/8)
		fmt.Println("data page:", chiplet)
		a.pageTable.Insert(page)
	}

	pState.nextVAddr += pageSize * uint64(numPages)
	return nextVAddr
}
