package device

import "strings"

// A DeviceMemoryState handles the internal memory allocation algorithms
type DeviceMemoryState interface {
	setInitialAddress(addr uint64)
	getInitialAddress() uint64
	setStorageSize(size uint64)
	getStorageSize() uint64
	addSinglePAddr(addr uint64)
	popNextAvailablePAddrs() uint64
	noAvailablePAddrs() bool
	allocateMultiplePages(numPages int) []uint64
	allocatePageOnChiplet(chiplet int) uint64
	allocatePageTablePage(vAddr, pAddr uint64) uint64
}

// NewDeviceMemoryState creates a new device memory state based on allocator type.
func NewDeviceMemoryState(log2pagesize uint64, memAllocatorType string) DeviceMemoryState {

	// This particular part is in poor taste, but is the faster way to get HSL
	// type information to the memory allocator right now
	hslSubType := ""
	// the order is important!!!
	if strings.Contains(memAllocatorType, "hugehslaware") {
		hslSubType = strings.Split(memAllocatorType, "-")[1]
		memAllocatorType = "hugehslaware"
	} else if strings.Contains(memAllocatorType, "hslaware") {
		hslSubType = strings.Split(memAllocatorType, "-")[1]
		memAllocatorType = "hslaware"
	}

	switch memAllocatorType {
	case "default" /*AllocatorTypeDefault*/ :
		return newDeviceRegularMemoryState(log2pagesize)
	case "buddy" /*AllocatorTypeBuddy*/ :
		return newDeviceBuddyMemoryState(log2pagesize)
	case "interleaved" /*AllocatorTypeInterleaved*/ :
		return newDeviceInterleavedMemoryState(log2pagesize)
	case "pta" /*AllocatorTypePageTableAware*/ :
		return newdevicePageTableAwareMemoryState(log2pagesize)
	case "xorpta":
		return newdeviceXorPTAMemState(log2pagesize)
	case "hugepagesinterleaved":
		return newdeviceHugePageInterleavedMemoryState(log2pagesize)
	case "partitioned":
		return newdevicePartitionedMemState(log2pagesize)
	case "partitionedpta":
		return newdevicePartitionedXorPtaMemState(log2pagesize)
	case "lasp":
		return newdeviceLASPMemState(log2pagesize)
	case "lasptpp":
		return newdeviceLaspTppMemState(log2pagesize)
	case "laspnaivept":
		return newdeviceLASPnaivePTMemState(log2pagesize)
	case "hslaware":
		return newdevicePartitionedHSLAwareMemState(log2pagesize, hslSubType)
	case "hugelasp":
		return newdeviceHugePageLASPMemState(log2pagesize)
	case "hugehslaware":
		return newdeviceHugePageHSLMemState(log2pagesize, hslSubType)
	default:
		panic("Invalid memory allocator type")
	}
}

func newDeviceRegularMemoryState(log2pagesize uint64) DeviceMemoryState {
	return &deviceMemoryStateImpl{
		log2PageSize: log2pagesize,
	}
}

//original implementation of DeviceMemoryState holding free addresses in array
type deviceMemoryStateImpl struct {
	log2PageSize    uint64
	initialAddress  uint64
	storageSize     uint64
	availablePAddrs []uint64
}

func (dms *deviceMemoryStateImpl) setInitialAddress(addr uint64) {
	dms.initialAddress = addr

	pageSize := uint64(1 << dms.log2PageSize)
	endAddr := dms.initialAddress + dms.storageSize
	for addr := dms.initialAddress; addr < endAddr; addr += pageSize {
		dms.addSinglePAddr(addr)
	}
}

func (dms *deviceMemoryStateImpl) getInitialAddress() uint64 {
	return dms.initialAddress
}

func (dms *deviceMemoryStateImpl) setStorageSize(size uint64) {
	dms.storageSize = size
}

func (dms *deviceMemoryStateImpl) getStorageSize() uint64 {
	return dms.storageSize
}

func (dms *deviceMemoryStateImpl) addSinglePAddr(addr uint64) {
	dms.availablePAddrs = append(dms.availablePAddrs, addr)
}

func (dms *deviceMemoryStateImpl) popNextAvailablePAddrs() uint64 {
	nextPAddr := dms.availablePAddrs[0]
	dms.availablePAddrs = dms.availablePAddrs[1:]
	return nextPAddr
}

func (dms *deviceMemoryStateImpl) noAvailablePAddrs() bool {
	return len(dms.availablePAddrs) == 0
}

func (dms *deviceMemoryStateImpl) allocateMultiplePages(
	numPages int,
) (pAddrs []uint64) {
	for i := 0; i < numPages; i++ {
		pAddr := dms.popNextAvailablePAddrs()
		pAddrs = append(pAddrs, pAddr)
	}
	return pAddrs
}

func (dms *deviceMemoryStateImpl) allocatePageTablePage(vAddr, pAddr uint64) uint64 {
	return dms.allocateMultiplePages(1)[0]
}

func (dms *deviceMemoryStateImpl) allocatePageOnChiplet(chiplet int) uint64 {
	panic("not implemented")
	return 0
}
