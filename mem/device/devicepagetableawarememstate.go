package device

import (
	"gitlab.com/akita/mem"
)

// A devicePageTableAwareMemoryState implements DeviceMemoryState as a interleaved allocator
type devicePageTableAwareMemoryState struct {
	memoryOnChiplet   uint64
	numChiplets       uint64
	numBankPerChiplet uint64
	bankSize          uint64
	log2PageSize      uint64
	initialAddress    uint64
	storageSize       uint64
	availablePAddrs   []uint64
	// curChipletForData              uint64
	// curBankForData                 uint64
	pAddrsForPageTable             [][]uint64
	log2MemoryBankInterleavingSize uint64
}

func newdevicePageTableAwareMemoryState(log2pagesize uint64) DeviceMemoryState {
	return &devicePageTableAwareMemoryState{
		log2PageSize:                   log2pagesize,
		memoryOnChiplet:                4 * mem.GB,
		numChiplets:                    4,
		numBankPerChiplet:              16,
		bankSize:                       256 * mem.MB,
		log2MemoryBankInterleavingSize: 12,
		//		curChiplet:        0,
	}
}

func (dims *devicePageTableAwareMemoryState) assertAddrIsPageAligned(addr uint64) {
	if ((addr >> dims.log2PageSize) << dims.log2PageSize) != addr {
		panic("oh no!")
	}
}

func (dims *devicePageTableAwareMemoryState) setInitialAddress(addr uint64) {
	dims.assertAddrIsPageAligned(addr)
	dims.initialAddress = addr
	dims.availablePAddrs = make([]uint64, 0)
	// , dims.numChiplets)
	dims.pAddrsForPageTable = make([][]uint64, dims.numChiplets)
	// for i := range dims.availablePAddrs {
	// 	dims.availablePAddrs[i] = make([][]uint64, dims.numBankPerChiplet)
	// 	for j := range dims.availablePAddrs[i] {
	// 		dims.availablePAddrs[i][j] = make([]uint64, 0)

	// 	}
	// }

	// pageSize := uint64(1 << dims.log2PageSize)
	// for chiplet := uint64(0); chiplet < dims.numChiplets; chiplet = chiplet + 1 {
	// 	dims.availablePAddrs[chiplet] = make([]uint64, 0)
	// 	initialAddr := dims.initialAddress + chiplet*dims.numBankPerChiplet*dims.bankSize
	// 	endAddr := initialAddr + dims.bankSize*dims.numBankPerChiplet

	// 	for addr := initialAddr; addr < endAddr; addr += pageSize {
	// 		dims.availablePAddrs[chiplet][bank] = append(dims.availablePAddrs[chiplet], addr)
	// 	}
	// }
	// fjdjf
	// fmt.Println(fmt.Sprintf("%x", dims.initialAddress))
	pageSize := uint64(1 << dims.log2PageSize)
	// fmt.Println(dims.storageSize, float64(dims.storageSize)/float64(1*mem.GB))
	endAddr := dims.initialAddress + dims.storageSize - 1*mem.GB
	// 7*mem.GB
	// dims.bankSize*dims.numBankPerChiplet
	for addr := dims.initialAddress; addr < endAddr; addr += pageSize {
		dims.availablePAddrs = append(dims.availablePAddrs, addr)
		// dims.numBankPerChiplet * pageSize {
		// for chiplet := uint64(0); chiplet < dims.numChiplets; chiplet += 1 {
		// temp := chiplet*dims.memoryOnChiplet + addr
		// for k := uint64(0); k < dims.numBankPerChiplet; k += 1 {
		// dims.availablePAddrs[chiplet][k] = append(dims.availablePAddrs[chiplet][k], temp)
		// temp = temp + pageSize
		// }
		// }
	}
	// memNeededForPageTables := 1 * mem.GB / (dims.numBankPerChiplet * dims.numChiplets)
	for addr := endAddr; addr < endAddr+1*mem.GB; {
		for i := 0; i < int(dims.numChiplets); i++ {
			for j := 0; j < int(dims.numBankPerChiplet); j++ {
				// (i*dims.numMemoryBankPerChiplet+j)*uint64(1<<dims.log2MemoryBankInterleavingSize)
				dims.pAddrsForPageTable[i] = append(dims.pAddrsForPageTable[i], addr)
				addr += uint64(1 << dims.log2MemoryBankInterleavingSize)
				// fmt.Println(addr, uint64(1<<12), dims.log2MemoryBankInterleavingSize)
			}
		}
		// dims.availablePAddrs = append(dims.availablePAddrs, addr)
		// dims.numBankPerChiplet * pageSize {
		// for chiplet := uint64(0); chiplet < dims.numChiplets; chiplet += 1 {
		// temp := chiplet*dims.memoryOnChiplet + addr
		// for k := uint64(0); k < dims.numBankPerChiplet; k += 1 {
		// dims.availablePAddrs[chiplet][k] = append(dims.availablePAddrs[chiplet][k], temp)
		// temp = temp + pageSize
		// }
		// }
	}
}

func (dims *devicePageTableAwareMemoryState) getInitialAddress() uint64 {
	return dims.initialAddress
}

func (dims *devicePageTableAwareMemoryState) setStorageSize(size uint64) {
	dims.storageSize = size
}

func (dims *devicePageTableAwareMemoryState) getStorageSize() uint64 {
	return dims.storageSize
}

func (dims *devicePageTableAwareMemoryState) addSinglePAddr(addr uint64) {
	panic("best not to call this function!")
	dims.assertAddrIsPageAligned(addr)
	// addr = addr - dims.initialAddress
	// chiplet := addr / dims.memoryOnChiplet
	// bank := (addr % dims.memoryOnChiplet) / dims.bankSize
	dims.availablePAddrs = append(dims.availablePAddrs, addr)
}

func (dims *devicePageTableAwareMemoryState) popNextAvailablePAddrs() uint64 {
	// nextPAddr := dims.availablePAddrs[dims.curChipletForData][dims.curBankForData][0]
	// dims.availablePAddrs[dims.curChipletForData][dims.curBankForData] = dims.availablePAddrs[dims.curChipletForData][dims.curBankForData][1:]
	// dims.curBankForData = (dims.curBankForData + 1) % dims.numBankPerChiplet
	// if dims.curBankForData == 0 {
	// 	dims.curChipletForData = (dims.curChipletForData + 1) % dims.numChiplets
	// }
	// Need to take care of the following: what happens if there is no space on  a chiplet to allocate a page table page?
	//	for ; len(dims.availablePAddrs[dims.curChiplet]) == 0; dims.curChiplet = (dims.curChiplet + 1) % dims.numChiplets {
	//	}
	nextPAddr := dims.availablePAddrs[0]
	dims.availablePAddrs = dims.availablePAddrs[1:]
	return nextPAddr
}

func (dims *devicePageTableAwareMemoryState) noAvailablePAddrs() bool {
	return len(dims.availablePAddrs) == 0
}

func (dims *devicePageTableAwareMemoryState) allocateMultiplePages(
	numPages int,
) (pAddrs []uint64) {
	for i := 0; i < numPages; i++ {
		pAddr := dims.popNextAvailablePAddrs()
		pAddrs = append(pAddrs, pAddr)
	}
	return pAddrs
}

func (dims *devicePageTableAwareMemoryState) allocatePageTablePage(vAddr, pAddr uint64) uint64 {
	chipletNum := vAddr & 3
	return dims.allocatePageTableOnChiplet(chipletNum)
}

func (dims *devicePageTableAwareMemoryState) allocatePageTableOnChiplet(
	chiplet uint64,
) (pAddr uint64) {
	// todo: want to distribute page table pages across banks
	// bank := dims.curBankForPageTable[chiplet]
	// pAddr = dims.availablePAddrs[chiplet][bank][0]
	// dims.availablePAddrs[chiplet][bank] = dims.availablePAddrs[chiplet][bank][1:]
	// dims.curBankForPageTable[chiplet] = (bank + 1) % dims.numBankPerChiplet
	pAddr = dims.pAddrsForPageTable[chiplet][0]
	dims.pAddrsForPageTable[chiplet] = dims.pAddrsForPageTable[chiplet][1:]
	// fmt.Println(chiplet, pAddr)
	return pAddr
}

func (dims *devicePageTableAwareMemoryState) allocatePageOnChiplet(chiplet int) uint64 {
	panic("not implemented")
	return 0
}
