package cache

import (
	// "fmt"
	"gitlab.com/akita/util/ca"
)

// NewDirectory returns a new directory object
func NewPageWalkCacheDirectory(
	set, way, blockSize int,
	victimFinder VictimFinder,
	log2PageSize uint64,
	bitsPerLevel uint64,
) *PageWalkCacheDirectoryImpl {
	d := new(PageWalkCacheDirectoryImpl)
	d.victimFinder = victimFinder
	d.Sets = make([]Set, set)

	d.NumSets = set
	d.NumWays = way
	d.BlockSize = blockSize
	d.log2PageSize = log2PageSize
	d.bitsPerLevel = bitsPerLevel
	// fmt.Println(d.log2PageSize)
	curMask := ^((uint64(1) << d.log2PageSize) - 1)
	curMask = curMask << d.bitsPerLevel
	for i := 2; i >= 0; i-- {
		d.masks[i] = curMask
		// fmt.Println(fmt.Sprintf("%x", curMask))
		// h := fmt.Sprintf("%x", curMask)
		// h2 := fmt.Sprintf("%x", d.masks[i])
		// fmt.Println(i, h, h2)
		curMask = curMask << d.bitsPerLevel
	}
	d.Reset()

	return d
}

type PageWalkCacheDirectoryImpl struct {
	NumSets      int
	NumWays      int
	BlockSize    int
	masks        [3]uint64
	log2PageSize uint64
	bitsPerLevel uint64

	Sets []Set

	victimFinder VictimFinder
}

// TotalSize returns the maximum number of bytes can be stored in the cache
func (d *PageWalkCacheDirectoryImpl) TotalSize() uint64 {
	return uint64(d.NumSets) * uint64(d.NumWays) * uint64(d.BlockSize)
}

// Get the set that a certain address should store at
func (d *PageWalkCacheDirectoryImpl) getSet(reqAddr uint64) (set *Set, setID int) {
	setID = int(reqAddr / uint64(d.BlockSize) % uint64(d.NumSets))
	set = &d.Sets[setID]
	return
}

// Lookup finds the block that reqAddr. If the reqAddr is valid
// in the cache, return the block information. Otherwise, return nil
func (d *PageWalkCacheDirectoryImpl) Lookup(PID ca.PID, reqAddr uint64) *Block {
	set, _ := d.getSet(reqAddr) // only one set at the moment, so doesn't really matter overmuch
	var longestMatch *Block = nil
	curLevel := -1
	for _, block := range set.Blocks {
		if block.IsValid && block.PID == PID {
			tag := block.Tag & ^uint64(3)
			level := int(block.Tag & uint64(3))
			mask := d.masks[level]
			maskedReqAddr := mask & reqAddr
			if tag == maskedReqAddr && level > curLevel {
				longestMatch = block
				curLevel = level
				// h := fmt.Sprintf("%x", maskedReqAddr)
				// h1 := fmt.Sprintf("%x", tag)
				// h2 := fmt.Sprintf("%x", reqAddr)
				// fmt.Println("EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE")
				// fmt.Println(h, h1, h2, level, curLevel)
			}
		}
	}
	return longestMatch
	// reqAddrAsString := strconv.FormatUint(reqAddr, 64)
	// lengthOfLongestMatch := 0
	// var longestMatch *Block = nil
	// for _, block := range set.Blocks {
	// 	if block.IsValid && block.PID == PID {
	// 		tagAsString := strconv.FormatUint(block.Tag, 64)
	// 		if strings.HasPrefix(reqAddrAsString, tagAsString) {
	// 			if len(tagAsString) > lengthOfLongestMatch {
	// 				lengthOfLongestMatch = len(tagAsString)
	// 				longestMatch = block
	// 			}
	// 		}
	// 	}
	// }
	// return longestMatch
}

func (d *PageWalkCacheDirectoryImpl) LookupExactAddress(PID ca.PID, reqAddr uint64) *Block {
	set, _ := d.getSet(reqAddr) // only one set at the moment, so doesn't really matter overmuch
	var toReturn *Block = nil
	for _, block := range set.Blocks {
		if block.IsValid && block.PID == PID && block.Tag == reqAddr {
			toReturn = block
		}
	}
	return toReturn
}

// FindVictim returns a block that can be used to stored data at address addr.
//
// If it is valid, the cache controller need to decide what to do to evict the
// the data in the block
func (d *PageWalkCacheDirectoryImpl) FindVictim(addr uint64) *Block {
	set, _ := d.getSet(addr)
	block := d.victimFinder.FindVictim(set)
	return block
}

// Visit moves the block to the end of the LRUQueue
func (d *PageWalkCacheDirectoryImpl) Visit(block *Block) {
	set := d.Sets[block.SetID]
	for i, b := range set.LRUQueue {
		if b == block {
			set.LRUQueue = append(set.LRUQueue[:i], set.LRUQueue[i+1:]...)
			break
		}
	}
	set.LRUQueue = append(set.LRUQueue, block)
}

// GetSets returns all the sets in a directory
func (d *PageWalkCacheDirectoryImpl) GetSets() []Set {
	return d.Sets
}

// Reset will mark all the blocks in the directory invalid
func (d *PageWalkCacheDirectoryImpl) Reset() {
	d.Sets = make([]Set, d.NumSets)
	for i := 0; i < d.NumSets; i++ {
		for j := 0; j < d.NumWays; j++ {
			block := new(Block)
			block.IsValid = false
			block.SetID = i
			block.WayID = j
			block.CacheAddress = uint64(i*d.NumWays+j) * uint64(d.BlockSize)
			d.Sets[i].Blocks = append(d.Sets[i].Blocks, block)
			d.Sets[i].LRUQueue = append(d.Sets[i].LRUQueue, block)
		}
	}
}

// WayAssociativity returns the number of ways per set in the cache.
func (d *PageWalkCacheDirectoryImpl) WayAssociativity() int {
	return d.NumWays
}

// WayAssociativity returns the number of ways per set in the cache.
func (d *PageWalkCacheDirectoryImpl) FormulateWriteAddress(req uint64) uint64 {
	vAddr := req & ^uint64(3)
	level := req & uint64(3)
	mask := d.masks[level]
	maskedVAddr := vAddr & mask
	toReturn := (maskedVAddr | level)
	return toReturn
}
