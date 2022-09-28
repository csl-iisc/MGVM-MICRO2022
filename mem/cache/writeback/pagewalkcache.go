package writeback

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem"
	"gitlab.com/akita/mem/cache"
	"gitlab.com/akita/util"
	"gitlab.com/akita/util/akitaext"
	"gitlab.com/akita/util/pipelining"
)

type pwcPipelineItem struct {
	taskID string
	trans  *transaction
}

func (p pwcPipelineItem) TaskID() string {
	return p.taskID
}

// A PageWalkCache in writeback package  is a PageWalkCache that performs the write-back policy
type PageWalkCache struct {
	*akita.TickingComponent

	TopPort akita.Port

	dirStageBuffer   util.Buffer
	dirToBankBuffers []util.Buffer

	topSender akitaext.BufferedSender

	lookupBuffer util.Buffer
	pipeline     pipelining.Pipeline

	topParser  *pageWalkCacheTopParser
	dirStage   *pageWalkCacheDirectoryStage
	bankStages []*pageWalkCacheBankStage

	storage        *mem.Storage
	directory      *cache.PageWalkCacheDirectoryImpl
	log2BlockSize  uint64
	numReqPerCycle int

	state                cacheState
	inFlightTransactions []*transaction
}

func (c *PageWalkCache) Tick(now akita.VTimeInSec) bool {
	madeProgress := false
	if c.state != cacheStatePaused {
		madeProgress = c.runPipeline(now) || madeProgress
	}
	return madeProgress
}

func (c *PageWalkCache) runPipeline(now akita.VTimeInSec) bool {
	madeProgress := false

	madeProgress = c.runStage(now, c.topSender) || madeProgress
	for _, bs := range c.bankStages {
		madeProgress = c.runStage(now, bs) || madeProgress
	}
	madeProgress = c.runStage(now, c.dirStage) || madeProgress
	// putting pipeline here
	// madeProgress = c.runStage(now, c.pipeline) || madeProgress
	madeProgress = c.pipeline.Tick(now) || madeProgress
	madeProgress = c.runStage(now, c.topParser) || madeProgress
	return madeProgress
}

func (c *PageWalkCache) runStage(now akita.VTimeInSec, stage akita.Ticker) bool {
	madeProgress := false
	for i := 0; i < c.numReqPerCycle; i++ {
		madeProgress = stage.Tick(now) || madeProgress
	}
	return madeProgress
}
