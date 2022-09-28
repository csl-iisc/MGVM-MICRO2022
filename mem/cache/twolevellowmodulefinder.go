package cache

import (
	"gitlab.com/akita/akita"
)

// TwoLevelLowModuleFinder determines the lower module based on
// two levels. The outer level is banked (like chiplets in an MCM).
// The inner level is  interleaved at some granularity (like memory banks in a
// chiplet of an MCM)
type TwoLevelLowModuleFinder struct {
	UseAddressSpaceLimitation bool
	LowAddress                uint64
	HighAddress               uint64
	OuterBankSize             uint64
	NumOuterBanks             uint64
	InnerInterleavingSize     uint64
	NumInnerBanks             uint64
	LowModules                []akita.Port
	ModuleForOtherAddresses   akita.Port
}

// Find returns the module determined by using the outer banking policy
// and the inner interleaving policy.
func (d *TwoLevelLowModuleFinder) Find(address uint64) akita.Port {
	if d.UseAddressSpaceLimitation &&
		(address >= d.HighAddress || address < d.LowAddress) {
		panic("panicking because addres out of range haha")
	}
	address = address - d.LowAddress
	chipletNumber := address / d.OuterBankSize
	offsetWithinChiplet := address - (chipletNumber * d.OuterBankSize)
	bankNumber := (offsetWithinChiplet / d.InnerInterleavingSize) % d.NumInnerBanks
	return d.LowModules[chipletNumber*d.NumInnerBanks+bankNumber]
}

// NewTwoLevelLowModuleFinder creates a new object of the
// TwoLevelLowModuleFinder with parameters.
func NewTwoLevelLowModuleFinder(
	LowAddress uint64, HighAddress uint64,
	OuterBankSize uint64, NumOuterBanks uint64,
	InnerInterleavingSize uint64, NumInnerBanks uint64,
) *TwoLevelLowModuleFinder {
	d := new(TwoLevelLowModuleFinder)
	d.UseAddressSpaceLimitation = true
	d.LowAddress = LowAddress
	d.HighAddress = HighAddress
	d.OuterBankSize = OuterBankSize
	d.NumOuterBanks = NumOuterBanks
	d.InnerInterleavingSize = InnerInterleavingSize
	d.NumInnerBanks = NumInnerBanks
	return d
}

// LocalInterleavedLowModuleFinder helps a local device to determine whether
// a request must be routed to local low module or remote low modules, assuming
// that the request space is distributed in an interleaved fashion.
type LocalInterleavedLowModuleFinder struct {
	LocalLowModuleID uint64
	NumLowModules    uint64
	InterleavingSize uint64
	LocalLowModule   akita.Port
	RemoteLowModule  akita.Port
}

// Find returns either the local lowmoudle or port to remote modules
func (f *LocalInterleavedLowModuleFinder) Find(address uint64) akita.Port {
	number := address / f.InterleavingSize % f.NumLowModules
	if number == f.LocalLowModuleID {
		return f.LocalLowModule
	} else {
		return f.RemoteLowModule
	}
}

// NewLocalInterleavedLowModuleFinder return a new LocalInterleavedLowModuleFinder
func NewLocalInterleavedLowModuleFinder(id uint64, numModules uint64,
	interleavingSize uint64) *LocalInterleavedLowModuleFinder {
	f := new(LocalInterleavedLowModuleFinder)
	f.LocalLowModuleID = id
	f.NumLowModules = numModules
	f.InterleavingSize = interleavingSize

	return f
}

type LocalXORLowModuleFinder struct {
	LocalLowModuleID uint64
	NumLowModules    uint64
	NumTerms         int
	NumBitsPerTerm   int
	OffsetBits       int
	LocalLowModule   akita.Port
	RemoteLowModule  akita.Port
}

func (f *LocalXORLowModuleFinder) Find(address uint64) akita.Port {
	index := uint64(0)
	address = address >> f.OffsetBits
	mask := (uint64(1) << f.NumBitsPerTerm) - 1
	for i := 0; i < f.NumTerms; i++ {
		index = index ^ (address & mask)
		address = address >> f.NumBitsPerTerm
	}
	if index == f.LocalLowModuleID {
		return f.LocalLowModule
	} else {
		return f.RemoteLowModule
	}
}

func NewLocalXORLowModuleFinder(id uint64, numModules uint64, numTerms int,
	numBitsPerTerm int, offsetBits int) *LocalXORLowModuleFinder {
	f := new(LocalXORLowModuleFinder)
	f.LocalLowModuleID = id
	f.NumLowModules = numModules
	f.NumTerms = numTerms
	f.NumBitsPerTerm = numBitsPerTerm
	f.OffsetBits = offsetBits

	return f
}

type CustomTwoLevelLowModuleFinder struct {
	LocalLowModuleID uint64
	LocalLowModule   akita.Port
	RemoteLowModule  akita.Port
	Hashfunc         func(address uint64) uint64
	PrivateMode      bool
	offsetBits       uint64
}

// Find returns the module the hash function determines
func (f *CustomTwoLevelLowModuleFinder) Find(address uint64) akita.Port {
	if f.PrivateMode {
		return f.LocalLowModule
	}
	address = address >> f.offsetBits
	index := f.Hashfunc(address)
	if index == f.LocalLowModuleID {
		return f.LocalLowModule
	} else {
		return f.RemoteLowModule
	}
}

func NewCustomTwoLevelLowModuleFinder(offsetBits uint64, id uint64) *CustomTwoLevelLowModuleFinder {
	f := new(CustomTwoLevelLowModuleFinder)
	f.LocalLowModuleID = id
	f.Hashfunc = func(address uint64) uint64 {
		return address % 4
	}
	f.PrivateMode = false
	f.offsetBits = offsetBits
	return f
}
