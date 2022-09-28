package device

// DeviceType marks the type of a device.
type DeviceType int

// Defines supported devices.
const (
	DeviceTypeInvalid DeviceType = iota
	DeviceTypeCPU
	DeviceTypeGPU
	DeviceTypeUnifiedGPU
)

// AllocatorType marks the type of memory allocator
type AllocatorType int

// Defines supported allocation algorithms
const (
	AllocatorTypeDefault AllocatorType = iota
	AllocatorTypeBuddy
	AllocatorTypeInterleaved
	AllocatorTypePageTableAware
	AllocatorTypePartitioned
	AllocatorTypeLASP
	AllocatorTypeLASPTPP
)

// MemoryAllocatorType global flag variable for setting the allocator type
var MemoryAllocatorType AllocatorType = AllocatorTypeInterleaved //AllocatorTypePageTableAware

// Device is a CPU or GPU managed by the driver.
type Device struct {
	ID                 int
	Type               DeviceType
	UnifiedGPUIDs      []int
	ActualGPUs         []*Device
	nextActualGPUIndex int
	MemState           DeviceMemoryState
}

// SetTotalMemSize sets total memory size
func (d *Device) SetTotalMemSize(size uint64) {
	d.MemState.setStorageSize(size)
}

func (d *Device) allocatePage() (pAddr uint64) {
	if d.Type == DeviceTypeUnifiedGPU {
		return d.allocateUnifiedGPUPage()
	}

	d.mustHaveSpaceLeft()
	pAddr = d.MemState.popNextAvailablePAddrs()

	return pAddr
}

func (d *Device) allocatePageOnChiplet(chiplet int) (pAddr uint64) {
	d.mustHaveSpaceLeft()
	pAddr = d.MemState.allocatePageOnChiplet(chiplet)

	return pAddr
}

func (d *Device) allocatePageTablePage(vAddr uint64, pAddr uint64) uint64 {
	if d.Type == DeviceTypeUnifiedGPU {
		panic("oh no!")
	}

	d.mustHaveSpaceLeft()
	pAddrToReturn := d.MemState.allocatePageTablePage(vAddr, pAddr)

	return pAddrToReturn
}

func (d *Device) allocateMultiplePages(numPages int) (pAddrs []uint64) {
	if d.Type == DeviceTypeUnifiedGPU {
		return d.allocateMultipleUnifiedGPUPages(numPages)
	}
	d.mustHaveSpaceLeft()
	pAddrs = d.MemState.allocateMultiplePages(numPages)

	return pAddrs
}

func (d *Device) mustHaveSpaceLeft() {
	if d.MemState.noAvailablePAddrs() {
		panic("out of memory")
	}
}

func (d *Device) allocateUnifiedGPUPage() (pAddr uint64) {
	dev := d.ActualGPUs[d.nextActualGPUIndex]
	pAddr = dev.allocatePage()
	d.nextActualGPUIndex = (d.nextActualGPUIndex + 1) % len(d.ActualGPUs)
	return pAddr
}

func (d *Device) allocateMultipleUnifiedGPUPages(numPages int) (pAddrs []uint64) {
	dev := d.ActualGPUs[d.nextActualGPUIndex]
	pAddrs = dev.allocateMultiplePages(numPages)
	d.nextActualGPUIndex = (d.nextActualGPUIndex + 1) % len(d.ActualGPUs)
	return pAddrs
}
