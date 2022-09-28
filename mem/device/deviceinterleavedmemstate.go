package device

import (
	"gitlab.com/akita/mem"
)

// A deviceInterleavedMemoryState implements DeviceMemoryState as a interleaved allocator
type deviceInterleavedMemoryState struct {
	memoryOnChiplet   uint64
	numChiplets       uint64
	numBankPerChiplet uint64
	bankSize          uint64
	log2PageSize      uint64
	initialAddress    uint64
	storageSize       uint64
	availablePAddrs   []uint64
}

func (dims *deviceInterleavedMemoryState) setInitialAddress(addr uint64) {
	dims.initialAddress = addr
	pageSize := uint64(1 << dims.log2PageSize)
	endAddr := dims.initialAddress + dims.storageSize
	for addr := dims.initialAddress; addr < endAddr; addr += pageSize {
		dims.availablePAddrs = append(dims.availablePAddrs, addr)
	}

	// for addr := dims.initialAddress; addr < endAddr; addr += dims.numBankPerChiplet * pageSize {
	// 	for chiplet := uint64(0); chiplet < dims.numChiplets; chiplet += 1 {
	// 		temp := chiplet*dims.memoryOnChiplet + addr
	// 		for k := uint64(0); k < dims.numBankPerChiplet; k += 1 {
	// 			dims.addSinglePAddr(temp)
	// 			temp = temp + pageSize
	// 		}
	// 	}
	// }
}

func newDeviceInterleavedMemoryState(log2pagesize uint64) DeviceMemoryState {
	return &deviceInterleavedMemoryState{
		log2PageSize:      log2pagesize,
		memoryOnChiplet:   4 * mem.GB,
		numChiplets:       4,
		numBankPerChiplet: 16,
		bankSize:          256 * mem.MB,
	}
}

func (dims *deviceInterleavedMemoryState) getInitialAddress() uint64 {
	return dims.initialAddress
}

func (dims *deviceInterleavedMemoryState) setStorageSize(size uint64) {
	dims.storageSize = size
}

func (dims *deviceInterleavedMemoryState) getStorageSize() uint64 {
	return dims.storageSize
}

func (dims *deviceInterleavedMemoryState) addSinglePAddr(addr uint64) {
	dims.availablePAddrs = append(dims.availablePAddrs, addr)
}

func (dims *deviceInterleavedMemoryState) popNextAvailablePAddrs() uint64 {
	nextPAddr := dims.availablePAddrs[0]
	dims.availablePAddrs = dims.availablePAddrs[1:]
	return nextPAddr
}

func (dims *deviceInterleavedMemoryState) noAvailablePAddrs() bool {
	return len(dims.availablePAddrs) == 0
}

func (dims *deviceInterleavedMemoryState) allocateMultiplePages(
	numPages int,
) (pAddrs []uint64) {
	for i := 0; i < numPages; i++ {
		pAddr := dims.popNextAvailablePAddrs()
		pAddrs = append(pAddrs, pAddr)
	}
	return pAddrs
}

func (dims *deviceInterleavedMemoryState) allocatePageTablePage(vAddr, pAddr uint64) uint64 {
	return dims.allocateMultiplePages(1)[0]
}

func (dims *deviceInterleavedMemoryState) allocatePageOnChiplet(chiplet int) uint64 {
	panic("not implemented")
	return 0
}
