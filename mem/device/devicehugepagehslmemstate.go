package device

import (
	"fmt"
	"strconv"

	"gitlab.com/akita/mem"
)

type deviceHugePageHSLMemState struct {
	memoryOnChiplet   uint64
	numChiplets       uint64
	numBankPerChiplet uint64
	bankSize          uint64
	log2PageSize      uint64
	initialAddress    uint64
	storageSize       uint64
	availablePAddrs   [][]uint64

	pageTablePages [][]uint64
	hslSubType     string
}

func (dims *deviceHugePageHSLMemState) setInitialAddress(addr uint64) {
	dims.initialAddress = addr
	pageSize := uint64(1 << dims.log2PageSize)
	endAddrData := dims.initialAddress + dims.storageSize - 1*mem.GB
	for addr < endAddrData {
		for i := 0; i < int(dims.numChiplets); i++ {
			for j := 0; j < int(dims.numBankPerChiplet); j++ {
				dims.availablePAddrs[i] = append(dims.availablePAddrs[i], addr)
				addr += pageSize
			}
		}
	}
	pageSize = uint64(1 << 12) // page table page size is 4K
	endAddr := dims.initialAddress + dims.storageSize
	for addr < endAddr {
		for i := 0; i < int(dims.numChiplets); i++ {
			for j := 0; j < int(dims.numBankPerChiplet); j++ {
				dims.pageTablePages[i] = append(dims.pageTablePages[i], addr)
				addr += pageSize
			}
		}
	}
}

func newdeviceHugePageHSLMemState(log2pagesize uint64, hslSubType string) DeviceMemoryState {
	return &deviceHugePageHSLMemState{
		log2PageSize:      log2pagesize,
		memoryOnChiplet:   4 * mem.GB,
		numChiplets:       4,
		numBankPerChiplet: 16,
		bankSize:          256 * mem.MB,
		availablePAddrs:   make([][]uint64, 4),
		pageTablePages:    make([][]uint64, 4),
		hslSubType:        hslSubType,
	}
}

func (dims *deviceHugePageHSLMemState) getInitialAddress() uint64 {
	return dims.initialAddress
}

func (dims *deviceHugePageHSLMemState) setStorageSize(size uint64) {
	dims.storageSize = size
}

func (dims *deviceHugePageHSLMemState) getStorageSize() uint64 {
	return dims.storageSize
}

func (dims *deviceHugePageHSLMemState) addSinglePAddr(addr uint64) {
	panic("This should not be called!")
	// dims.availablePAddrs = append(dims.availablePAddrs, addr)
}

func (dims *deviceHugePageHSLMemState) popNextAvailablePAddrs() uint64 {
	// panic("This should not be called!")
	nextPAddr := dims.availablePAddrs[0][0]
	dims.availablePAddrs[0] = dims.availablePAddrs[0][1:]
	return nextPAddr
}

func (dims *deviceHugePageHSLMemState) noAvailablePAddrs() bool {
	return len(dims.availablePAddrs) == 0
}

func (dims *deviceHugePageHSLMemState) allocateMultiplePages(
	numPages int) (pAddrs []uint64) {
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

func (dims *deviceHugePageHSLMemState) allocatePageTablePage(
	vAddr, pAddr uint64) (pAddrToReturn uint64) {
	interleaving, _ := strconv.Atoi(dims.hslSubType)
	chiplet := ((vAddr) / uint64(interleaving)) % 4

	fmt.Println("page table page:", vAddr, chiplet)
	pAddrToReturn = dims.pageTablePages[chiplet][0]
	dims.pageTablePages[chiplet] = dims.pageTablePages[chiplet][1:]
	return pAddrToReturn
}

func (dims *deviceHugePageHSLMemState) allocatePageOnChiplet(
	chiplet int) uint64 {
	pAddr := dims.availablePAddrs[chiplet][0]
	dims.availablePAddrs[chiplet] = dims.availablePAddrs[chiplet][1:]
	return pAddr
}
