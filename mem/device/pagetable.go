package device

import (
	"container/list"
	"fmt"
	"sync"

	"gitlab.com/akita/util/ca"
)

// A Page is an entry in the page table, maintaining the information about how
// to translate a virtual address to a physical address
type Page struct {
	PID         ca.PID
	PAddr       uint64
	VAddr       uint64
	PageSize    uint64
	Valid       bool
	DeviceID    uint64
	Unified     bool
	IsMigrating bool
	IsPinned    bool
}

// A PageTable holds the a list of pages.
type PageTable interface {
	Insert(page Page)
	Remove(pid ca.PID, vAddr uint64)
	Find(pid ca.PID, Addr uint64) (Page, bool)
	Update(page Page)
	FindAddr(pid ca.PID, vAddr uint64, level uint64) uint64
	PageTablePagesAsBytes(pid ca.PID) []uint64
	// GetPageTableAsBuffer(pid ca.PID)
}

// NewPageTable creates a new page table
func NewPageTable(log2PageSize uint64) *PageTableImpl {
	return &PageTableImpl{
		log2PageSize: log2PageSize,
		tables:       make(map[ca.PID]*processTableImpl),
	}
}

// PageTableImpl is the default implementation of a Page Table
type PageTableImpl struct {
	sync.Mutex
	log2PageSize uint64
	tables       map[ca.PID]*processTableImpl
	memAllocator MemoryAllocator
	entries      *list.List
	entriesTable map[uint64]*list.Element
}

func (pt *PageTableImpl) getTable(pid ca.PID) *processTableImpl {
	pt.Lock()
	defer pt.Unlock()
	table, found := pt.tables[pid]
	if !found {
		table = newProcessTable(pt.log2PageSize, 4, 9, pt.memAllocator)
		pt.tables[pid] = table
	}

	return table
}

func (pt *PageTableImpl) FindAddr(pid ca.PID, vAddr uint64, level uint64) uint64 {
	table := pt.getTable(pid)
	return table.findAddr(vAddr, level)
}

func (pt *PageTableImpl) AlignToPage(addr uint64) uint64 {
	return (addr >> pt.log2PageSize) << pt.log2PageSize
}

// Insert put a new page into the PageTable
func (pt *PageTableImpl) Insert(page Page) {
	table := pt.getTable(page.PID)
	// fmt.Println(page.VAddr)
	table.insert(page)
}

// Remove removes the entry in the page table that contains the target
// address.
func (pt *PageTableImpl) Remove(pid ca.PID, vAddr uint64) {
	table := pt.getTable(pid)
	table.remove(vAddr)
}

// Find returns the page that contains the given virtual address. The bool
// return value invicates if the page is found or not.
func (pt *PageTableImpl) Find(pid ca.PID, vAddr uint64) (Page, bool) {
	table := pt.getTable(pid)
	vAddr = pt.AlignToPage(vAddr)
	return table.find(vAddr)
}

// Update changes the field of an existing page. The PID and the VAddr field
// will be used to locate the page to update.
func (pt *PageTableImpl) Update(page Page) {
	table := pt.getTable(page.PID)
	table.update(page)
}

func (pt *PageTableImpl) SetMemoryAllocator(a MemoryAllocator) {
	pt.memAllocator = a
}

func (pt *PageTableImpl) PageTablePagesAsBytes(pid ca.PID) []uint64 {
	table := pt.getTable(pid)
	return table.PageTablePagesAsBytes()
}

func (pt *PageTableImpl) Rearrange(vAddr uint64) uint64 {
	return pt.getTable(0).rearrange(vAddr)
}

func (pt *PageTableImpl) GetRoot(pid ca.PID) uint64 {
	return pt.getTable(pid).root.pAddr
}

func (pt *PageTableImpl) AddOffset(root, vAddr uint64) uint64 {
	return pt.getTable(0).addOffset(root, vAddr)
}

func (pt *PageTableImpl) NextLevel(vAddr uint64) uint64 {
	return pt.getTable(0).nextLevel(vAddr)
}

func (pt *PageTableImpl) MoveToLevel(vAddr uint64, level int) uint64 {
	for i := 0; i < level; i++ {
		vAddr = pt.getTable(0).nextLevel(vAddr)
	}
	return vAddr
}

type processTableImpl struct {
	sync.Mutex
	root             *treeNode
	log2PageSize     uint64
	bitsPerLevel     uint64
	bitsPerLevelMask uint64
	numChildren      uint64
	numLevels        uint64
	memAllocator     MemoryAllocator
	entries          *list.List
	entriesTable     map[uint64]*list.Element
}

type treeNode struct {
	sync.Mutex
	pAddr    uint64
	page     Page
	children []*treeNode
}

func newProcessTable(log2PageSize uint64, numLevels uint64, bitsPerLevel uint64, a MemoryAllocator) *processTableImpl {
	t := new(processTableImpl)
	t.log2PageSize = log2PageSize
	t.numLevels = numLevels
	t.bitsPerLevel = bitsPerLevel
	t.numChildren = uint64(1) << bitsPerLevel
	t.bitsPerLevelMask = ^(^uint64(0) << t.bitsPerLevel)
	t.memAllocator = a
	t.entries = list.New()
	t.entriesTable = make(map[uint64]*list.Element)
	return t
}

func (t *processTableImpl) rearrange(vAddr uint64) uint64 {
	rearrangedVpn := uint64(0)
	vpn := vAddr / (uint64(1) << t.log2PageSize)
	var i uint64 = 0
	for ; i < t.numLevels; i++ {
		temp := vpn & t.bitsPerLevelMask
		rearrangedVpn = (rearrangedVpn << t.bitsPerLevel) | temp
		vpn = vpn >> t.bitsPerLevel
	}
	return rearrangedVpn
}

func (t *processTableImpl) newTreeNode(PID ca.PID, GPUID uint64, vAddr, pAddr uint64, level uint64) *treeNode {
	n := new(treeNode)
	isLeaf := level == t.numLevels-1
	if !isLeaf {
		vAddr = vAddr >> (t.bitsPerLevel*(t.numLevels-level-1) + t.log2PageSize)
		n.pAddr = (t.memAllocator).allocatePageTablePage(PID, int(GPUID), vAddr, pAddr)
		fmt.Println("allocated physical page:", n.pAddr, (((n.pAddr-4096)>>12)%32)/8, level) // this arithmetic is off !!
		n.children = make([]*treeNode, t.numChildren)
	}
	return n
}

func (t *processTableImpl) insert(page Page) {
	t.Lock()
	defer t.Unlock()
	if t.root == nil {
		t.root = t.newTreeNode(page.PID, page.DeviceID, 0, 0, 0)
	}
	n := t.root
	vAddr := t.rearrange(page.VAddr)
	var i uint64 = 0
	for ; i < t.numLevels-1; i++ {
		indexOfChild := vAddr & t.bitsPerLevelMask
		if n.children[indexOfChild] == nil {
			newTreeNode := t.newTreeNode(page.PID, page.DeviceID, page.VAddr, page.PAddr, i)
			n.children[indexOfChild] = newTreeNode
		}
		n = n.children[indexOfChild]
		vAddr = vAddr >> t.bitsPerLevel
	}
	indexOfChild := vAddr & t.bitsPerLevelMask
	if n.children[indexOfChild] != nil {
		panic("page already present")
	}
	n.children[indexOfChild] = t.newTreeNode(page.PID, page.DeviceID, page.VAddr, page.PAddr, i)
	n.children[indexOfChild].pAddr = page.PAddr
	n.children[indexOfChild].page = page
	t.pageMustNotExist(page.VAddr)
	elem := t.entries.PushBack(page)
	t.entriesTable[page.VAddr] = elem
}

func (t *processTableImpl) remove(vAddr uint64) {
	t.Lock()
	defer t.Unlock()
	n := t.root
	var i uint64
	rearrangedVpn := t.rearrange(vAddr)
	for i = 0; i < t.numLevels-1; i++ {
		indexOfChild := rearrangedVpn & t.bitsPerLevelMask
		n = n.children[indexOfChild]
		rearrangedVpn = rearrangedVpn >> t.bitsPerLevel
	}
	indexOfChild := rearrangedVpn & t.bitsPerLevelMask
	if n.children[indexOfChild] == nil {
		panic("virtual page does not exist!")
	}
	if n.children[indexOfChild].page.VAddr != vAddr {
		panic("DFS wrongly implemented")
	}
	n.children[indexOfChild] = nil

	t.pageMustExist(vAddr)
	elem := t.entriesTable[vAddr]
	t.entries.Remove(elem)
	delete(t.entriesTable, vAddr)

}

func (t *processTableImpl) update(page Page) {
	t.Lock()
	defer t.Unlock()
	n := t.root
	vAddr := t.rearrange(page.VAddr)
	var i uint64
	for i = 0; i < t.numLevels; i++ {
		indexOfChild := vAddr & t.bitsPerLevelMask
		n = n.children[indexOfChild]
		vAddr = vAddr >> t.bitsPerLevel
	}
	n.pAddr = page.PAddr
	n.page = page

	t.pageMustExist(page.VAddr)
	elem := t.entriesTable[page.VAddr]
	elem.Value = page

}

func (t *processTableImpl) find(vAddr uint64) (Page, bool) {
	t.Lock()
	defer t.Unlock()

	elem, found := t.entriesTable[vAddr]
	if found {
		return elem.Value.(Page), true
	}

	return Page{}, false
}

func (t *processTableImpl) findAddr(vAddr, level uint64) uint64 {
	if level >= t.numLevels {
		panic("level!")
	}
	n := t.root
	vAddr = t.rearrange(vAddr)
	var i uint64
	for i = 0; i < t.numLevels; i++ {
		indexOfChild := vAddr & t.bitsPerLevelMask
		if i == level {
			if n.children[indexOfChild] == nil {
				panic("oh no!")
			}
			return n.pAddr + indexOfChild*8
		}
		n = n.children[indexOfChild]
		vAddr = vAddr >> t.bitsPerLevel
	}
	panic("boo!")
	return 0
}

func (t *processTableImpl) addOffset(root, vAddr uint64) uint64 {
	// fmt.Println(root, vAddr&t.bitsPerLevelMask, (vAddr&t.bitsPerLevelMask)*8)
	return root + (vAddr&t.bitsPerLevelMask)*8
}

func (t *processTableImpl) nextLevel(vAddr uint64) uint64 {
	return vAddr >> t.bitsPerLevel
}

func (t *processTableImpl) pageMustExist(vAddr uint64) {
	_, found := t.entriesTable[vAddr]
	if !found {
		panic("page does not exist")
	}
}

func (t *processTableImpl) pageMustNotExist(vAddr uint64) {
	_, found := t.entriesTable[vAddr]
	if found {
		panic("page exist")
	}
}

func (t *processTableImpl) PageTablePagesAsBytes() []uint64 {
	pagesAsInts := make([]uint64, 0)
	queue := make([]*treeNode, 0)
	queue = append(queue, t.root)
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		numChildren := len(n.children)
		if numChildren > 0 {
			if numChildren != int(t.numChildren) {
				panic("oh no!")
			}
			pagesAsInts = append(pagesAsInts, n.pAddr)
			for i, child := range n.children {
				pAddrOfChild := uint64(0)
				if child != nil {
					pAddrOfChild = child.pAddr
					queue = append(queue, n.children[i])
				}
				pagesAsInts = append(pagesAsInts, pAddrOfChild)
			}
		}
	}
	// fmt.Println(pagesAsInts)
	return pagesAsInts
}
