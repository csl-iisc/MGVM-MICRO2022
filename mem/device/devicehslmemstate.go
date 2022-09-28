package device

import (
	"fmt"
	"strconv"

	"gitlab.com/akita/mem"
)

// A devicePartitionedHSLAwareMemState implements DeviceMemoryState as a interleaved allocator
type devicePartitionedHSLAwareMemState struct {
	memoryOnChiplet   uint64
	numChiplets       uint64
	numBankPerChiplet uint64
	bankSize          uint64
	log2PageSize      uint64
	initialAddress    uint64
	storageSize       uint64
	availablePAddrs   [][]uint64
	hslSubType        string
}

func (dims *devicePartitionedHSLAwareMemState) setInitialAddress(addr uint64) {
	dims.initialAddress = addr
	fmt.Println("&&&&&&", dims.initialAddress)
	pageSize := uint64(1 << dims.log2PageSize)
	endAddr := dims.initialAddress + dims.storageSize
	for addr < endAddr {
		for i := 0; i < int(dims.numChiplets); i++ {
			for j := 0; j < int(dims.numBankPerChiplet); j++ {
				dims.availablePAddrs[i] = append(dims.availablePAddrs[i], addr)
				addr += pageSize
			}
		}
	}

}

func newdevicePartitionedHSLAwareMemState(log2pagesize uint64, hslSubType string) DeviceMemoryState {
	return &devicePartitionedHSLAwareMemState{
		log2PageSize:      log2pagesize,
		memoryOnChiplet:   4 * mem.GB,
		numChiplets:       4,
		numBankPerChiplet: 16,
		bankSize:          256 * mem.MB,
		availablePAddrs:   make([][]uint64, 4),
		hslSubType:        hslSubType,
	}
}

func (dims *devicePartitionedHSLAwareMemState) getInitialAddress() uint64 {
	return dims.initialAddress
}

func (dims *devicePartitionedHSLAwareMemState) setStorageSize(size uint64) {
	dims.storageSize = size
}

func (dims *devicePartitionedHSLAwareMemState) getStorageSize() uint64 {
	return dims.storageSize
}

func (dims *devicePartitionedHSLAwareMemState) addSinglePAddr(addr uint64) {
	panic("This should not be called!")
	// dims.availablePAddrs = append(dims.availablePAddrs, addr)
}

func (dims *devicePartitionedHSLAwareMemState) popNextAvailablePAddrs() uint64 {
	// panic("This should not be called!")
	nextPAddr := dims.availablePAddrs[0][0]
	dims.availablePAddrs[0] = dims.availablePAddrs[0][1:]
	return nextPAddr
}

func (dims *devicePartitionedHSLAwareMemState) noAvailablePAddrs() bool {
	return len(dims.availablePAddrs) == 0
}

func (dims *devicePartitionedHSLAwareMemState) allocateMultiplePages(
	numPages int,
) (pAddrs []uint64) {
	pagesPerChiplet := numPages / int(dims.numChiplets)
	fmt.Println("***********numPages:", numPages, pagesPerChiplet, numPages%4)
	for i := 0; i < int(dims.numChiplets); i++ {
		for j := 0; j < pagesPerChiplet; j++ {
			pAddr := dims.availablePAddrs[i][0]
			dims.availablePAddrs[i] = dims.availablePAddrs[i][1:]
			pAddrs = append(pAddrs, pAddr)
		}
	}
	remainingPages := numPages % int(dims.numChiplets)
	chiplet := 0
	for j := 0; j < remainingPages; j++ {
		pAddr := dims.availablePAddrs[chiplet][0]
		dims.availablePAddrs[chiplet] = dims.availablePAddrs[chiplet][1:]
		chiplet++
		pAddrs = append(pAddrs, pAddr)
	}
	return pAddrs
}

func (dims *devicePartitionedHSLAwareMemState) allocatePageTablePage(vAddr, pAddr uint64) (pAddrToReturn uint64) {
	//Set the correct chiplet ID based on HSL Sub Type
	interleaving, _ := strconv.Atoi(dims.hslSubType)
	chiplet := ((vAddr) / uint64(interleaving)) % 4

	fmt.Println("page table page:", vAddr, chiplet)
	pAddrToReturn = dims.availablePAddrs[chiplet][0]
	dims.availablePAddrs[chiplet] = dims.availablePAddrs[chiplet][1:]
	return pAddrToReturn
}

func (dims *devicePartitionedHSLAwareMemState) allocatePageOnChiplet(
	chiplet int) uint64 {
	pAddr := dims.availablePAddrs[chiplet][0]
	dims.availablePAddrs[chiplet] = dims.availablePAddrs[chiplet][1:]
	return pAddr
}
