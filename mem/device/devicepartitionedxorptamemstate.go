package device

import (
	"fmt"

	"gitlab.com/akita/mem"
)

// A devicePartitionedXorPtaMemState implements DeviceMemoryState as a interleaved allocator
type devicePartitionedXorPtaMemState struct {
	memoryOnChiplet   uint64
	numChiplets       uint64
	numBankPerChiplet uint64
	bankSize          uint64
	log2PageSize      uint64
	initialAddress    uint64
	storageSize       uint64
	availablePAddrs   [][]uint64
	numTerms          int
	numBitsPerTerm    int
	termMask          uint64
}

func (dims *devicePartitionedXorPtaMemState) setInitialAddress(addr uint64) {
	dims.initialAddress = addr
	fmt.Println("&&&&&&", dims.initialAddress)
	pageSize := uint64(1 << dims.log2PageSize)
	endAddr := dims.initialAddress + dims.storageSize
	for addr < endAddr {
		for i := 0; i < int(dims.numChiplets); i++ {
			for j := 0; j < int(dims.numBankPerChiplet); j++ {
				// (i*dims.numMemoryBankPerChiplet+j)*uint64(1<<dims.log2MemoryBankInterleavingSize)
				dims.availablePAddrs[i] = append(dims.availablePAddrs[i], addr)
				addr += pageSize
				// fmt.Println(addr, uint64(1<<12), dims.log2MemoryBankInterleavingSize)
			}
		}
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

func newdevicePartitionedXorPtaMemState(log2pagesize uint64) DeviceMemoryState {
	return &devicePartitionedXorPtaMemState{
		log2PageSize:      log2pagesize,
		memoryOnChiplet:   4 * mem.GB,
		numChiplets:       4,
		numBankPerChiplet: 16,
		bankSize:          256 * mem.MB,
		availablePAddrs:   make([][]uint64, 4),
		numTerms:          4,
		numBitsPerTerm:    2,
		termMask:          (uint64(1) << 2) - 1,
	}
}

func (dims *devicePartitionedXorPtaMemState) getInitialAddress() uint64 {
	return dims.initialAddress
}

func (dims *devicePartitionedXorPtaMemState) setStorageSize(size uint64) {
	dims.storageSize = size
}

func (dims *devicePartitionedXorPtaMemState) getStorageSize() uint64 {
	return dims.storageSize
}

func (dims *devicePartitionedXorPtaMemState) addSinglePAddr(addr uint64) {
	panic("This should not be called!")
	// dims.availablePAddrs = append(dims.availablePAddrs, addr)
}

func (dims *devicePartitionedXorPtaMemState) popNextAvailablePAddrs() uint64 {
	// panic("This should not be called!")
	nextPAddr := dims.availablePAddrs[0][0]
	dims.availablePAddrs[0] = dims.availablePAddrs[0][1:]
	return nextPAddr
}

func (dims *devicePartitionedXorPtaMemState) noAvailablePAddrs() bool {
	return len(dims.availablePAddrs) == 0
}

func (dims *devicePartitionedXorPtaMemState) allocateMultiplePages(
	numPages int,
) (pAddrs []uint64) {
	// if numPages%4 != 0 {
	// panic("The number of pages needed are not a multiple of 4!")
	// }
	pagesPerChiplet := numPages / int(dims.numChiplets)
	fmt.Println("***********numPages:", numPages, pagesPerChiplet, numPages%4)
	for i := 0; i < int(dims.numChiplets); i++ {
		for j := 0; j < pagesPerChiplet; j++ {
			pAddr := dims.availablePAddrs[i][0]
			// fmt.Println(((pAddr - dims.initialAddress) >> 12 % 32) / 8)
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

func (dims *devicePartitionedXorPtaMemState) allocatePageTablePage(vAddr, pAddr uint64) (pAddrToReturn uint64) {
	// chiplet := (((pAddr - dims.initialAddress) >> 12) % 32) / 8
	chiplet := uint64(0)

	//vAddr is already correctly shifted
	for i := 0; i < dims.numTerms; i++ {
		chiplet = chiplet ^ (vAddr & dims.termMask)
		vAddr = vAddr >> dims.numBitsPerTerm
	}
	fmt.Println("page table page:", chiplet)
	pAddrToReturn = dims.availablePAddrs[chiplet][0]
	dims.availablePAddrs[chiplet] = dims.availablePAddrs[chiplet][1:]
	return pAddrToReturn
}

func (dms *devicePartitionedXorPtaMemState) allocatePageOnChiplet(chiplet int) uint64 {
	panic("not implemented")
	return 0
}
