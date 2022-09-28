package device

import (
	// "fmt"
	"gitlab.com/akita/mem"
)

// A deviceXorPTAMemState implements DeviceMemoryState as a interleaved allocator
type deviceXorPTAMemState struct {
	*devicePageTableAwareMemoryState
	numTerms       int
	numBitsPerTerm int
	termMask       uint64
}

func newdeviceXorPTAMemState(log2pagesize uint64) DeviceMemoryState {
	state := &deviceXorPTAMemState{
		&devicePageTableAwareMemoryState{
			log2PageSize:                   log2pagesize,
			memoryOnChiplet:                4 * mem.GB,
			numChiplets:                    4,
			numBankPerChiplet:              16,
			bankSize:                       256 * mem.MB,
			log2MemoryBankInterleavingSize: 12,
		},
		4,
		2,
		0,
	}
	state.termMask = (uint64(1) << state.numBitsPerTerm) - 1
	return state
}

func (dims *deviceXorPTAMemState) allocatePageTablePage(vAddr, pAddr uint64) (pAddrToReturn uint64) {
	chipletNum := uint64(0)
	//vAddr is already correctly shifted
	for i := 0; i < dims.numTerms; i++ {
		chipletNum = chipletNum ^ (vAddr & dims.termMask)
		vAddr = vAddr >> dims.numBitsPerTerm
	}
	pAddrToReturn = dims.allocatePageTableOnChiplet(chipletNum)
	// fmt.Println(fmt.Sprintf("%x", vAddr), chipletNum, fmt.Sprintf("%x", pAddr))
	return pAddrToReturn
}
