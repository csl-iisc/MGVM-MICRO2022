package mmu

import (
	"fmt"
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem/cache/writeback"
	"gitlab.com/akita/mem/device"
	"gitlab.com/akita/util"
	"gitlab.com/akita/util/akitaext"
)

// A Builder can build MMU component
type Builder struct {
	engine                   akita.Engine
	freq                     akita.Freq
	log2PageSize             uint64
	pageTable                *device.PageTableImpl
	migrationServiceProvider akita.Port
	maxNumReqInFlight        int
	pageWalkingLatency       int
	numChiplets              uint64
	//	lowAddr                  uint64
	//	totMem                   uint64
	//	bankSize                 uint64
	//	numMemoryBanksPerChiplet uint64
}

// MakeBuilder creates a new builder
func MakeBuilder() Builder {
	return Builder{
		freq:              1 * akita.GHz,
		log2PageSize:      12,
		maxNumReqInFlight: 8, //16,
	}
}

// WithEngine sets the engine to be used with the MMU
func (b Builder) WithEngine(engine akita.Engine) Builder {
	b.engine = engine
	return b
}

// WithFreq sets the frequency that the MMU to work at
func (b Builder) WithFreq(freq akita.Freq) Builder {
	b.freq = freq
	return b
}

// WithLog2PageSize sets the page size that the mmu support.
func (b Builder) WithLog2PageSize(log2PageSize uint64) Builder {
	b.log2PageSize = log2PageSize
	return b
}

// WithPageTable sets the page table that the MMU uses.
func (b Builder) WithPageTable(pageTable *device.PageTableImpl) Builder {
	b.pageTable = pageTable
	return b
}

/*
// WithMigrationServiceProvider sets the destination port that can perform
// page migration.
func (b Builder) WithMigrationServiceProvider(p akita.Port) Builder {
	b.migrationServiceProvider = p
	return b
}
*/
// WithMaxNumReqInFlight sets the number of requests can be concurrently
// processed by the MMU.
func (b Builder) WithMaxNumReqInFlight(n int) Builder {
	b.maxNumReqInFlight = n
	return b
}

/*
// WithPageWalkingLatency sets the number of cycles required for walking a page
// table.
func (b Builder) WithPageWalkingLatency(n int) Builder {
	b.pageWalkingLatency = n
	return b
}
*/
// WithNumChiplets sets the number of cycles required for walking a page
// table.
func (b Builder) WithNumChiplets(n uint64) Builder {
	b.numChiplets = n
	return b
}

/*
// WithLowAddr sets the number of cycles required for walking a page
// table.
func (b Builder) WithLowAddr(la uint64) Builder {
	b.lowAddr = la
	return b
}

// WithTotMem sets the number of cycles required for walking a page
// table.
func (b Builder) WithTotMem(ha uint64) Builder {
	b.totMem = ha
	return b
}

// WithBankSize sets the number of cycles required for walking a page
// table.
func (b Builder) WithBankSize(n uint64) Builder {
	b.bankSize = n
	return b
}

// WithNumMemoryBankPerChiplet sets the number of cycles required for walking a page
// table.
func (b Builder) WithNumMemoryBankPerChiplet(n uint64) Builder {
	b.numMemoryBanksPerChiplet = n
	return b
}
*/
// Build returns a newly created MMU component
func (b Builder) Build(name string) *MMUImpl {
	mmu := new(MMUImpl)
	mmu.TickingComponent = *akita.NewTickingComponent(
		name, b.engine, b.freq, mmu)
	//mmu.migrationQueueSize = 4096

	mmu.ToTop = akita.NewLimitNumMsgPort(mmu, 4096, name+".ToTop")
	mmu.ControlPort = akita.NewLimitNumMsgPort(mmu, 1, name+".ControlPort")

	//mmu.MigrationPort = akita.NewLimitNumMsgPort(mmu, 1, name+".MigrationPort")
	//might want to change capacity later
	mmu.TranslationPort = akita.NewLimitNumMsgPort(mmu, 16, name+".TranslationPort")
	//mmu.MigrationServiceProvider = b.migrationServiceProvider

	mmu.topSender = akitaext.NewBufferedSender(mmu.ToTop, util.NewBuffer(4096))
	if b.pageTable != nil {
		mmu.pageTable = b.pageTable
	} else {
		panic("no page table!")
	}
	mmu.maxRequestsInFlight = b.maxNumReqInFlight
	fmt.Println("num walkers:", mmu.maxRequestsInFlight)
	//mmu.latency = b.pageWalkingLatency
	//mmu.PageAccesedByDeviceID = make(map[uint64][]uint64)
	pageWalkCacheBuilder := writeback.MakePageWalkCacheBuilder().
		WithEngine(b.engine).
		WithLog2PageSize(b.log2PageSize).
		WithBitsPerLevel(9)
	pageWalkCache := pageWalkCacheBuilder.Build("PageWalkCache")
	mmu.PageWalkCache = pageWalkCache.TopPort
	mmu.pageWalkCachePort = akita.NewLimitNumMsgPort(mmu, 4096, name+".ToTop")
	mmuToPageWalkCache := akita.NewDirectConnection("MMUToPageWalkCache", b.engine, b.freq)
	mmuToPageWalkCache.PlugIn(pageWalkCache.TopPort, 4)
	mmuToPageWalkCache.PlugIn(mmu.pageWalkCachePort, 4)
	mmu.sendStateInfo = false
	return mmu
}
