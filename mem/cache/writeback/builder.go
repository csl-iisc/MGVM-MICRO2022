package writeback

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem"
	"gitlab.com/akita/mem/cache"
	"gitlab.com/akita/util"
	"gitlab.com/akita/util/akitaext"
	"gitlab.com/akita/util/pipelining"
)

// A Builder can build writeback caches
type Builder struct {
	engine              akita.Engine
	freq                akita.Freq
	lowModuleFinder     cache.LowModuleFinder
	wayAssociativity    int
	log2BlockSize       uint64
	byteSize            uint64
	numMSHREntry        int
	numReqPerCycle      int
	writeBufferCapacity int
	maxInflightFetch    int
	maxInflightEviction int
	bankLatency         int
	pipelineLatency     int
	rdmaPort            akita.Port
}

func MakeBuilder() Builder {
	return Builder{
		freq:                1 * akita.GHz,
		wayAssociativity:    4,
		log2BlockSize:       6,
		byteSize:            512 * mem.KB,
		numMSHREntry:        16,
		numReqPerCycle:      1,
		writeBufferCapacity: 1024,
		maxInflightFetch:    128,
		maxInflightEviction: 128,
		bankLatency:         1,
		pipelineLatency:     10,
	}
}

func (b Builder) WithEngine(engine akita.Engine) Builder {
	b.engine = engine
	return b
}

func (b Builder) WithFreq(freq akita.Freq) Builder {
	b.freq = freq
	return b
}

func (b Builder) WithWayAssociativity(n int) Builder {
	b.wayAssociativity = n
	return b
}

func (b Builder) WithLog2BlockSize(n uint64) Builder {
	b.log2BlockSize = n
	return b
}

func (b Builder) WithNumMSHREntry(n int) Builder {
	b.numMSHREntry = n
	return b
}

func (b Builder) WithLowModuleFinder(f cache.LowModuleFinder) Builder {
	b.lowModuleFinder = f
	return b
}

func (b Builder) WithNumReqPerCycle(n int) Builder {
	b.numReqPerCycle = n
	return b
}

func (b Builder) WithByteSize(byteSize uint64) Builder {
	b.byteSize = byteSize
	return b
}

// WithWriteBufferSize sets the number of cachlines that can reside in the
// writebuffer.
func (b Builder) WithWriteBufferSize(n int) Builder {
	b.writeBufferCapacity = n
	return b
}

// WithMaxInflightFetch sets the number of concurrent fetch that the write-back
// cache can issue at the same time.
func (b Builder) WithMaxInflightFetch(n int) Builder {
	b.maxInflightFetch = n
	return b
}

// WithMaxInflightEviction sets the number of concurrent eviction that the
// write buffer can write to a low-level module.
func (b Builder) WithMaxInflightEviction(n int) Builder {
	b.maxInflightEviction = n
	return b
}

// WithMaxInflightEviction sets the number of concurrent eviction that the
// write buffer can write to a low-level module.
func (b Builder) WithRDMAPort(rdmaPort akita.Port) Builder {
	b.rdmaPort = rdmaPort
	return b
}

// Build creates a usable writeback cache.
func (b *Builder) Build(name string) *Cache {
	cache := new(Cache)
	cache.TickingComponent = akita.NewTickingComponent(
		name, b.engine, b.freq, cache)

	b.configureCache(cache)
	b.createPorts(cache)
	b.createPortSenders(cache)
	b.createInternalStages(cache)
	b.createInternalBuffers(cache)

	return cache
}

func (b *Builder) configureCache(cacheModule *Cache) {
	blockSize := 1 << b.log2BlockSize
	vimctimFinder := cache.NewLRUVictimFinder()
	numSet := int(b.byteSize / uint64(b.wayAssociativity*blockSize))
	directory := cache.NewDirectory(
		numSet, b.wayAssociativity, blockSize, vimctimFinder)
	mshr := cache.NewMSHR(b.numMSHREntry)
	storage := mem.NewStorage(b.byteSize)

	cacheModule.log2BlockSize = b.log2BlockSize
	cacheModule.numReqPerCycle = b.numReqPerCycle
	cacheModule.directory = directory
	cacheModule.mshr = mshr
	cacheModule.storage = storage
	cacheModule.lowModuleFinder = b.lowModuleFinder
	cacheModule.state = cacheStateRunning
}

func (b *Builder) createPorts(cache *Cache) {
	cache.TopPort = akita.NewLimitNumMsgPort(cache,
		cache.numReqPerCycle*32, cache.Name()+".ToTop")
	cache.BottomPort = akita.NewLimitNumMsgPort(cache,
		cache.numReqPerCycle*2, cache.Name()+".BottomPort")
	cache.ControlPort = akita.NewLimitNumMsgPort(cache,
		cache.numReqPerCycle*2, cache.Name()+".ControlPort")
	cache.ControlPort = akita.NewLimitNumMsgPort(cache,
		cache.numReqPerCycle*2, cache.Name()+".MMUPort")
}

func (b *Builder) createPortSenders(cache *Cache) {
	cache.topSender = akitaext.NewBufferedSender(
		cache.TopPort, util.NewBuffer(cache.numReqPerCycle*4))
	cache.bottomSender = akitaext.NewBufferedSender(
		cache.BottomPort, util.NewBuffer(cache.numReqPerCycle*4))
	cache.controlPortSender = akitaext.NewBufferedSender(
		cache.ControlPort, util.NewBuffer(cache.numReqPerCycle*4))
}

func (b *Builder) createInternalStages(cache *Cache) {
	cache.topParser = &topParser{cache: cache}
	cache.dirStage = &directoryStage{cache: cache}
	cache.bankStages = make([]*bankStage, 1)
	cache.bankStages[0] = &bankStage{
		cache:   cache,
		bankID:  0,
		latency: b.bankLatency,
	}
	cache.mshrStage = &mshrStage{cache: cache}
	cache.flusher = &flusher{cache: cache}
	cache.writeBuffer = &writeBufferStage{
		cache:               cache,
		writeBufferCapacity: b.writeBufferCapacity,
		maxInflightFetch:    b.maxInflightFetch,
		maxInflightEviction: b.maxInflightEviction,
		rdmaPort:            b.rdmaPort,
	}
	cache.lookupBuffer = util.NewBuffer(2 * b.numReqPerCycle)
	pipelineBuilder := pipelining.MakeBuilder().
		WithPipelineWidth(b.numReqPerCycle).WithNumStage(b.pipelineLatency).
		WithCyclePerStage(1).WithPostPipelineBuffer(cache.lookupBuffer)
	cache.pipeline = pipelineBuilder.Build(cache.Name() + "_pipeline")

}

func (b *Builder) createInternalBuffers(cache *Cache) {
	cache.dirStageBuffer = util.NewBuffer(2 * cache.numReqPerCycle)
	cache.dirToBankBuffers = make([]util.Buffer, 1)
	cache.dirToBankBuffers[0] = util.NewBuffer(2 * cache.numReqPerCycle)
	cache.writeBufferToBankBuffers = make([]util.Buffer, 1)
	cache.writeBufferToBankBuffers[0] = util.NewBuffer(2 * cache.numReqPerCycle)
	cache.mshrStageBuffer = util.NewBuffer(2 * cache.numReqPerCycle)
	cache.writeBufferBuffer = util.NewBuffer(2 * cache.numReqPerCycle)
}
