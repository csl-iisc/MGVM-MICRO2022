package writeback

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem"
	"gitlab.com/akita/mem/cache"
	"gitlab.com/akita/util"
	"gitlab.com/akita/util/akitaext"
	"gitlab.com/akita/util/pipelining"
)

// A PageWalkCacheBuilder can build writeback caches
type PageWalkCacheBuilder struct {
	engine           akita.Engine
	freq             akita.Freq
	lowModuleFinder  cache.LowModuleFinder
	wayAssociativity int
	log2BlockSize    uint64
	byteSize         uint64
	numReqPerCycle   int
	pipelineLatency  int
	log2PageSize     uint64
	bitsPerLevel     uint64
}

func MakePageWalkCacheBuilder() PageWalkCacheBuilder {
	return PageWalkCacheBuilder{
		freq:             1 * akita.GHz,
		wayAssociativity: 32, //16,
		log2BlockSize:    4,
		byteSize:         512, //256, //bytes
		numReqPerCycle:   4,
		pipelineLatency:  10,
	}
}

func (b PageWalkCacheBuilder) WithEngine(engine akita.Engine) PageWalkCacheBuilder {
	b.engine = engine
	return b
}

func (b PageWalkCacheBuilder) WithFreq(freq akita.Freq) PageWalkCacheBuilder {
	b.freq = freq
	return b
}

func (b PageWalkCacheBuilder) WithWayAssociativity(n int) PageWalkCacheBuilder {
	b.wayAssociativity = n
	return b
}

func (b PageWalkCacheBuilder) WithLog2BlockSize(n uint64) PageWalkCacheBuilder {
	b.log2BlockSize = n
	return b
}

func (b PageWalkCacheBuilder) WithLog2PageSize(n uint64) PageWalkCacheBuilder {
	b.log2PageSize = n
	return b
}

func (b PageWalkCacheBuilder) WithBitsPerLevel(n uint64) PageWalkCacheBuilder {
	b.bitsPerLevel = n
	return b
}

func (b PageWalkCacheBuilder) WithNumReqPerCycle(n int) PageWalkCacheBuilder {
	b.numReqPerCycle = n
	return b
}

func (b PageWalkCacheBuilder) WithByteSize(byteSize uint64) PageWalkCacheBuilder {
	b.byteSize = byteSize
	return b
}

// Build creates a usable writeback cache.
func (b *PageWalkCacheBuilder) Build(name string) *PageWalkCache {
	cache := new(PageWalkCache)
	cache.TickingComponent = akita.NewTickingComponent(
		name, b.engine, b.freq, cache)

	b.configureCache(cache)
	b.createPorts(cache)
	b.createPortSenders(cache)
	b.createInternalStages(cache)
	b.createInternalBuffers(cache)

	return cache
}

func (b *PageWalkCacheBuilder) configureCache(cacheModule *PageWalkCache) {
	blockSize := 1 << b.log2BlockSize
	vimctimFinder := cache.NewLRUVictimFinder()
	numSet := int(b.byteSize / uint64(b.wayAssociativity*blockSize))
	directory := cache.NewPageWalkCacheDirectory(
		numSet, b.wayAssociativity, blockSize, vimctimFinder, b.log2PageSize, b.bitsPerLevel)
	storage := mem.NewStorage(b.byteSize)

	cacheModule.log2BlockSize = b.log2BlockSize
	cacheModule.numReqPerCycle = b.numReqPerCycle
	cacheModule.directory = directory
	cacheModule.storage = storage
	cacheModule.state = cacheStateRunning
}

func (b *PageWalkCacheBuilder) createPorts(cache *PageWalkCache) {
	cache.TopPort = akita.NewLimitNumMsgPort(cache,
		cache.numReqPerCycle*2, cache.Name()+".ToTop")
}

func (b *PageWalkCacheBuilder) createPortSenders(cache *PageWalkCache) {
	cache.topSender = akitaext.NewBufferedSender(
		cache.TopPort, util.NewBuffer(cache.numReqPerCycle*4))
}

func (b *PageWalkCacheBuilder) createInternalStages(cache *PageWalkCache) {
	cache.topParser = &pageWalkCacheTopParser{cache: cache}
	cache.dirStage = &pageWalkCacheDirectoryStage{cache: cache}
	cache.bankStages = make([]*pageWalkCacheBankStage, 1)
	cache.bankStages[0] = &pageWalkCacheBankStage{
		cache:   cache,
		bankID:  0,
		latency: 1,
	}
	cache.lookupBuffer = util.NewBuffer(2 * b.numReqPerCycle)
	pipelineBuilder := pipelining.MakeBuilder().
		WithPipelineWidth(b.numReqPerCycle).WithNumStage(b.pipelineLatency).
		WithCyclePerStage(1).WithPostPipelineBuffer(cache.lookupBuffer)
	cache.pipeline = pipelineBuilder.Build(cache.Name() + "_pipeline")
}

func (b *PageWalkCacheBuilder) createInternalBuffers(cache *PageWalkCache) {
	cache.dirStageBuffer = util.NewBuffer(cache.numReqPerCycle)
	cache.dirToBankBuffers = make([]util.Buffer, 1)
	cache.dirToBankBuffers[0] = util.NewBuffer(cache.numReqPerCycle)
}
