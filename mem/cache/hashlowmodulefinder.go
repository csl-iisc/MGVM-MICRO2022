package cache

import "gitlab.com/akita/akita"

// HashLowModuleFinder determines the lower module based on a
// provided hashing function
// TODO: add a way to set the hash function dynamically
type HashLowModuleFinder struct {
	Hashfunc   func(address uint64) uint64
	LowModules []akita.Port
	offsetBits uint64
}

// Find returns the module the hash function determines
func (d *HashLowModuleFinder) Find(address uint64) akita.Port {
	address = address >> d.offsetBits
	port := d.Hashfunc(address)
	return d.LowModules[port]
}

// NewHashLowModuleFinder creates a new object of the
// HashLowModuleFinder with defaults
func NewHashLowModuleFinder(offsetBits uint64) *HashLowModuleFinder {
	d := new(HashLowModuleFinder)
	d.Hashfunc = func(address uint64) uint64 {
		return address % uint64(len(d.LowModules))
	}
	d.LowModules = make([]akita.Port, 0)
	d.offsetBits = offsetBits
	return d
}

func (d *HashLowModuleFinder) WithHashFunction(hashfunc func(address uint64) uint64) {
	d.Hashfunc = hashfunc
}
