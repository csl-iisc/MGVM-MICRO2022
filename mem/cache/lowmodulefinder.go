package cache

import (
	// "fmt"
	"gitlab.com/akita/akita"
)

// LowModuleFinder helps a cache unit or a akita to find the low module that
// should hold the data at a certain address
type LowModuleFinder interface {
	Find(address uint64) akita.Port
}

// SingleLowModuleFinder is used when a unit is connected with only one
// low module
type SingleLowModuleFinder struct {
	LowModule akita.Port
}

// Find simply returns the solo unit that it connects to
func (f *SingleLowModuleFinder) Find(address uint64) akita.Port {
	return f.LowModule
}

// InterleavedLowModuleFinder helps find the low module when the low modules
// maintains interleaved address space
type InterleavedLowModuleFinder struct {
	UseAddressSpaceLimitation bool
	LowAddress                uint64
	HighAddress               uint64
	InterleavingSize          uint64
	LowModules                []akita.Port
	ModuleForOtherAddresses   akita.Port
}

// Find returns the low module that has the data at provided address
func (f *InterleavedLowModuleFinder) Find(address uint64) akita.Port {
	if f.UseAddressSpaceLimitation &&
		(address >= f.HighAddress || address < f.LowAddress) {
		return f.ModuleForOtherAddresses
	}
	number := address / f.InterleavingSize % uint64(len(f.LowModules))
	return f.LowModules[number]
}

// NewInterleavedLowModuleFinder creates a new finder for interleaved lower
// modules
func NewInterleavedLowModuleFinder(interleavingSize uint64) *InterleavedLowModuleFinder {
	finder := new(InterleavedLowModuleFinder)
	finder.LowModules = make([]akita.Port, 0)
	finder.InterleavingSize = interleavingSize
	return finder
}

// BankedLowModuleFinder defines the lower level modules by address banks
type BankedLowModuleFinder struct {
	MemAddrOffset uint64
	BankSize      uint64
	LowModules    []akita.Port
}

// Find returns the port that can provide the data.
func (f *BankedLowModuleFinder) Find(address uint64) akita.Port {
	address = address - f.MemAddrOffset
	i := address / f.BankSize
	return f.LowModules[i]
}

// unc (f *BankedLowModuleFinder) GetModule(i uint64) akita.Port {
// 	return f.LowModules[i]
// }

// NewBankedLowModuleFinder returns a new BankedLowModuleFinder.
func NewBankedLowModuleFinder(memAddrOffset, bankSize uint64) *BankedLowModuleFinder {
	f := new(BankedLowModuleFinder)
	f.MemAddrOffset = memAddrOffset
	f.BankSize = bankSize
	f.LowModules = make([]akita.Port, 0)
	return f
}

type XORLowModuleFinder struct {
	NumLowModules  int
	NumTerms       int
	NumBitsPerTerm int
	OffsetBits     int
	LowModules     []akita.Port
}

func (f *XORLowModuleFinder) Find(address uint64) akita.Port {
	index := uint64(0)
	address = address >> f.OffsetBits
	mask := (uint64(1) << f.NumBitsPerTerm) - 1
	for i := 0; i < f.NumTerms; i++ {
		index = index ^ (address & mask)
		address = address >> f.NumBitsPerTerm
	}
	return f.LowModules[index]
}

func NewXORLowModuleFinder(numModules int, numTerms int, numBitsPerTerm int,
	offsetBits int) *XORLowModuleFinder {
	f := new(XORLowModuleFinder)
	f.NumLowModules = numModules
	f.NumTerms = numTerms
	f.NumBitsPerTerm = numBitsPerTerm
	f.OffsetBits = offsetBits

	return f
}

// BankedLowModuleFinder defines the lower level modules by address banks
type StripedLowModuleFinder struct {
	MemAddrOffset uint64
	NumBanks      uint64
	Striping      uint64
	LowModules    []akita.Port
}

// Find returns the port that can provide the data.
func (f *StripedLowModuleFinder) Find(address uint64) akita.Port {
	address = address - f.MemAddrOffset
	// i := address / f.Striping
	i := (address % (f.NumBanks * f.Striping)) / f.Striping
	// fmt.Println("&&&&&&&&&&", f.NumBanks, f.Striping)
	// fmt.Println(fmt.Sprintf("%x", f.MemAddrOffset), fmt.Sprintf("%x", address), i)
	return f.LowModules[i]
	// i%f.NumBanks]
}

// NewBankedLowModuleFinder returns a new BankedLowModuleFinder.
func NewStripedLowModuleFinder(memAddrOffset, numBanks, striping uint64) *StripedLowModuleFinder {
	f := new(StripedLowModuleFinder)
	f.MemAddrOffset = memAddrOffset
	// fmt.Println(fmt.Sprintf("%x", f.MemAddrOffset))
	f.NumBanks = numBanks
	f.Striping = striping
	f.LowModules = make([]akita.Port, 0)
	return f
}

// BankedLowModuleFinder defines the lower level modules by address banks
type StripedLocalVRemoteLowModuleFinder struct {
	MemAddrOffset           uint64
	NumBanks                uint64
	Striping                uint64
	LocalBankStart          uint64
	LocalBankEnd            uint64
	LowModules              []akita.Port
	ModuleForOtherAddresses akita.Port
}

// Find returns the port that can provide the data.
func (f *StripedLocalVRemoteLowModuleFinder) Find(address uint64) akita.Port {
	address = address - f.MemAddrOffset
	i := (address % (f.NumBanks * f.Striping)) / f.Striping
	// fmt.Println("&&&&&&&&&&", f.NumBanks, f.Striping)
	// % f.NumBanks
	if f.LocalBankStart <= i && i <= f.LocalBankEnd {
		return f.LowModules[i-f.LocalBankStart]
	}
	return f.ModuleForOtherAddresses
}

// NewBankedLowModuleFinder returns a new BankedLowModuleFinder.
func NewStripedLocalVRemoteLowModuleFinder(memAddrOffset, numBanks, striping, start, end uint64) *StripedLocalVRemoteLowModuleFinder {
	f := new(StripedLocalVRemoteLowModuleFinder)
	f.MemAddrOffset = memAddrOffset
	// fmt.Println("localvremote:", fmt.Sprintf("%x", f.MemAddrOffset))
	f.NumBanks = numBanks
	f.Striping = striping
	f.LocalBankStart = start
	f.LocalBankEnd = end
	f.LowModules = make([]akita.Port, 0)
	return f
}
