// Package internal provides the definition required for defining TLB.
package internal

import (
	"fmt"

	"github.com/google/btree"
	"gitlab.com/akita/mem/device"
	"gitlab.com/akita/util/ca"
)

// A Set holds a certain number of pages.
type Set interface {
	Lookup(pid ca.PID, vAddr uint64) (wayID int, page device.Page, found bool)
	Update(wayID int, page device.Page)
	Evict() (wayID int, ok bool)
	Visit(wayID int) (lruRank int)
}

func NewSet(numWays int) Set {
	s := &setImpl{}
	s.blocks = make([]*block, numWays)
	s.visitTree = btree.New(2)
	s.vAddrWayIDMap = make(map[string]int)
	for i := range s.blocks {
		b := &block{}
		s.blocks[i] = b
		b.wayID = i
		s.Visit(i)
	}
	return s
}

type block struct {
	page      device.Page
	wayID     int
	lastVisit uint64
}

func (b *block) Less(anotherBlock btree.Item) bool {
	return b.lastVisit < anotherBlock.(*block).lastVisit
}

type setImpl struct {
	blocks        []*block
	vAddrWayIDMap map[string]int
	visitTree     *btree.BTree
	visitCount    uint64
}

func (s *setImpl) keyString(pid ca.PID, vAddr uint64) string {
	return fmt.Sprintf("%d%016x", pid, vAddr)
}

func (s *setImpl) Lookup(pid ca.PID, vAddr uint64) (
	wayID int,
	page device.Page,
	found bool,
) {
	key := s.keyString(pid, vAddr)
	wayID, ok := s.vAddrWayIDMap[key]
	if !ok {
		return 0, device.Page{}, false
	}

	block := s.blocks[wayID]

	return block.wayID, block.page, true
}

func (s *setImpl) Update(wayID int, page device.Page) {
	block := s.blocks[wayID]
	key := s.keyString(block.page.PID, block.page.VAddr)
	delete(s.vAddrWayIDMap, key)

	block.page = page
	key = s.keyString(page.PID, page.VAddr)
	s.vAddrWayIDMap[key] = wayID
}

func (s *setImpl) Evict() (wayID int, ok bool) {
	if s.hasNothingToEvict() {
		return 0, false
	}

	wayID = s.visitTree.DeleteMin().(*block).wayID
	return wayID, true
}

func (s *setImpl) Visit(wayID int) int {
	visitedBlock := s.blocks[wayID]

	rank := 0
	findRank := func(i btree.Item) bool {
		if i.(*block).lastVisit > visitedBlock.lastVisit {
			rank++
		}
		return true
	}
	s.visitTree.AscendGreaterOrEqual(visitedBlock, findRank)

	s.visitTree.Delete(visitedBlock)

	s.visitCount++
	visitedBlock.lastVisit = s.visitCount
	s.visitTree.ReplaceOrInsert(visitedBlock)

	// fmt.Println(rank)

	return rank
}

func (s *setImpl) hasNothingToEvict() bool {
	return s.visitTree.Len() == 0
}
