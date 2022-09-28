package device

import (
	// "fmt"
	"gitlab.com/akita/mem"
)

// A deviceHugePageInterleavedMemoryState implements DeviceMemoryState as a interleaved allocator
type deviceHugePageInterleavedMemoryState struct {
	// memoryOnChiplet                uint64
	numBanks uint64
	// bankSize                       uint64
	log2PageSize   uint64
	initialAddress uint64
	storageSize    uint64
	// availablePAddrs                [][]uint64
	// pAddrsForPageTable             [][]uint64
	availablePAddrs    []uint64
	pAddrsForPageTable []uint64

	log2MemoryBankInterleavingSize uint64
	// nextBank                       uint64
}

func (dims *deviceHugePageInterleavedMemoryState) setInitialAddress(addr uint64) {
	dims.initialAddress = addr
	pageSize := uint64(1 << dims.log2PageSize)
	endAddr := dims.initialAddress + dims.storageSize - 1*mem.GB
	// numBanks := int(dims.numChiplets * dims.numBankPerChiplet)
	// bank := uint64(0)
	for addr := dims.initialAddress; addr < endAddr; addr += pageSize {
		dims.availablePAddrs = append(dims.availablePAddrs, addr)
		// dims.availablePAddrs[bank] = append(dims.availablePAddrs[bank], addr)
		// bank = (bank + 16) % dims.numBanks
	}
	// bank = 0
	pageTablePageSize := uint64(1 << dims.log2MemoryBankInterleavingSize)
	for addr := endAddr; addr < endAddr+1*mem.GB; addr += pageTablePageSize {
		dims.pAddrsForPageTable = append(dims.pAddrsForPageTable, addr)
		// bank = (bank + 1) % dims.numBanks
	}
	// fmt.Println("len of pAddrs for bank 31:", len(dims.availablePAddrs[31]))
}

func newdeviceHugePageInterleavedMemoryState(log2pagesize uint64) DeviceMemoryState {
	return &deviceHugePageInterleavedMemoryState{
		log2PageSize:                   log2pagesize,
		numBanks:                       32,
		log2MemoryBankInterleavingSize: 12,
		// nextBank:                       0,
		// memoryOnChiplet:   2 * mem.GB,
		// numChiplets:       4,
		// numBankPerChiplet: 8,
		// bankSize:          256 * mem.MB,
	}
}

func (dims *deviceHugePageInterleavedMemoryState) getInitialAddress() uint64 {
	return dims.initialAddress
}

func (dims *deviceHugePageInterleavedMemoryState) setStorageSize(size uint64) {
	dims.storageSize = size
}

func (dims *deviceHugePageInterleavedMemoryState) getStorageSize() uint64 {
	return dims.storageSize
}

func (dims *deviceHugePageInterleavedMemoryState) addSinglePAddr(addr uint64) {
	panic("don't use this function please")
	// bank := (addr - dims.initialAddress) % dims.numBanks
	// dims.availablePAddrs = append(dims.availablePAddrs[bank], addr)
}

func (dims *deviceHugePageInterleavedMemoryState) popNextAvailablePAddrs() uint64 {
	// nextPAddr := dims.availablePAddrs[dims.nextBank][0]
	// dims.availablePAddrs[dims.nextBank] = dims.availablePAddrs[dims.nextBank][1:]
	// dims.nextBank = (dims.nextBank + 16) % dims.numBanks
	nextPAddr := dims.availablePAddrs[0]
	dims.availablePAddrs = dims.availablePAddrs[1:]
	return nextPAddr
}

func (dims *deviceHugePageInterleavedMemoryState) noAvailablePAddrs() bool {
	return len(dims.availablePAddrs) == 0
}

func (dims *deviceHugePageInterleavedMemoryState) allocateMultiplePages(
	numPages int,
) (pAddrs []uint64) {
	for i := 0; i < numPages; i++ {
		pAddr := dims.popNextAvailablePAddrs()
		pAddrs = append(pAddrs, pAddr)
	}
	return pAddrs
}

func (dims *deviceHugePageInterleavedMemoryState) allocatePageTablePage(vAddr, pAddr uint64) uint64 {
	// nextPAddr := dims.pAddrsForPageTable[dims.nextBank][0]
	// dims.pAddrsForPageTable[dims.nextBank] = dims.pAddrsForPageTable[dims.nextBank][1:]
	// dims.nextBank = (dims.nextBank + 1) % dims.numBanks
	nextPAddr := dims.pAddrsForPageTable[0]
	dims.pAddrsForPageTable = dims.pAddrsForPageTable[1:]
	// dims.nextBank = (dims.nextBank + 1) % dims.numBanks
	return nextPAddr
}

func (dims *deviceHugePageInterleavedMemoryState) allocatePageOnChiplet(chiplet int) uint64 {
	panic("not implemented")
	return 0
}
