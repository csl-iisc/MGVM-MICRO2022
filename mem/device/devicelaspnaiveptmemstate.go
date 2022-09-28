package device

import (
	"fmt"

	"gitlab.com/akita/mem"
)

type deviceLASPnaivePTMemState struct {
	memoryOnChiplet   uint64
	numChiplets       uint64
	numBankPerChiplet uint64
	bankSize          uint64
	log2PageSize      uint64
	initialAddress    uint64
	storageSize       uint64
	availablePAddrs   [][]uint64

	lastChipletForPT uint64
}

func (dims *deviceLASPnaivePTMemState) setInitialAddress(addr uint64) {
	dims.initialAddress = addr
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

func newdeviceLASPnaivePTMemState(log2pagesize uint64) DeviceMemoryState {
	return &deviceLASPnaivePTMemState{
		log2PageSize:      log2pagesize,
		memoryOnChiplet:   4 * mem.GB,
		numChiplets:       4,
		numBankPerChiplet: 16,
		bankSize:          256 * mem.MB,
		availablePAddrs:   make([][]uint64, 4),
		lastChipletForPT:  0,
	}
}

func (dims *deviceLASPnaivePTMemState) getInitialAddress() uint64 {
	return dims.initialAddress
}

func (dims *deviceLASPnaivePTMemState) setStorageSize(size uint64) {
	dims.storageSize = size
}

func (dims *deviceLASPnaivePTMemState) getStorageSize() uint64 {
	return dims.storageSize
}

func (dims *deviceLASPnaivePTMemState) addSinglePAddr(addr uint64) {
	panic("This should not be called!")
	// dims.availablePAddrs = append(dims.availablePAddrs, addr)
}

func (dims *deviceLASPnaivePTMemState) popNextAvailablePAddrs() uint64 {
	// panic("This should not be called!")
	nextPAddr := dims.availablePAddrs[0][0]
	dims.availablePAddrs[0] = dims.availablePAddrs[0][1:]
	return nextPAddr
}

func (dims *deviceLASPnaivePTMemState) noAvailablePAddrs() bool {
	return len(dims.availablePAddrs) == 0
}

func (dims *deviceLASPnaivePTMemState) allocateMultiplePages(
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

func (dims *deviceLASPnaivePTMemState) allocatePageTablePage(
	vAddr, pAddr uint64) (pAddrToReturn uint64) {
	// chiplet := ((pAddr >> 12) % 64) / 16
	chiplet := dims.lastChipletForPT
	dims.lastChipletForPT = (dims.lastChipletForPT + 1) % dims.numChiplets
	fmt.Println("page table page:", chiplet)
	pAddrToReturn = dims.availablePAddrs[chiplet][0]
	dims.availablePAddrs[chiplet] = dims.availablePAddrs[chiplet][1:]
	return pAddrToReturn
}

func (dims *deviceLASPnaivePTMemState) allocatePageOnChiplet(
	chiplet int) uint64 {
	pAddr := dims.availablePAddrs[chiplet][0]
	dims.availablePAddrs[chiplet] = dims.availablePAddrs[chiplet][1:]
	return pAddr
}
