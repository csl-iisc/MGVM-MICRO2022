package writeback

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem"
	"gitlab.com/akita/mem/cache"
	"gitlab.com/akita/util"
	"gitlab.com/akita/util/akitaext"
	"gitlab.com/akita/util/pipelining"
)

type cacheState int

// don't forget to cast accessReq into the correct type at the end of pipeline
type cachePipelineItem struct {
	taskID string
	trans  *transaction
}

func (c cachePipelineItem) TaskID() string {
	return c.taskID
}

const (
	cacheStateInvalid cacheState = iota
	cacheStateRunning
	cacheStatePreFlushing
	cacheStateFlushing
	cacheStatePaused
)

// A Cache in writeback package  is a cache that performs the write-back policy
type Cache struct {
	*akita.TickingComponent

	TopPort     akita.Port
	BottomPort  akita.Port
	ControlPort akita.Port
	// MMUPort     akita.Port

	dirStageBuffer           util.Buffer
	dirToBankBuffers         []util.Buffer
	writeBufferToBankBuffers []util.Buffer
	mshrStageBuffer          util.Buffer
	writeBufferBuffer        util.Buffer

	topSender         akitaext.BufferedSender
	bottomSender      akitaext.BufferedSender
	controlPortSender akitaext.BufferedSender

	// pipeline can come here
	pipeline     pipelining.Pipeline
	lookupBuffer util.Buffer

	topParser   *topParser
	writeBuffer *writeBufferStage
	dirStage    *directoryStage
	bankStages  []*bankStage
	mshrStage   *mshrStage
	flusher     *flusher

	storage         *mem.Storage
	lowModuleFinder cache.LowModuleFinder
	directory       cache.Directory
	mshr            cache.MSHR
	log2BlockSize   uint64
	numReqPerCycle  int

	state                cacheState
	inFlightTransactions []*transaction
}

func (c *Cache) GetPipeline() pipelining.Pipeline {
	return c.pipeline
}

func (c *Cache) SetLowModuleFinder(lmf cache.LowModuleFinder) {
	c.lowModuleFinder = lmf
}

func (c *Cache) Tick(now akita.VTimeInSec) bool {
	madeProgress := false

	madeProgress = c.controlPortSender.Tick(now) || madeProgress

	if c.state != cacheStatePaused {
		madeProgress = c.runPipeline(now) || madeProgress
	}

	madeProgress = c.flusher.Tick(now) || madeProgress
	return true
	// return madeProgress
}

func (c *Cache) runPipeline(now akita.VTimeInSec) bool {
	madeProgress := false

	madeProgress = c.runStage(now, c.topSender) || madeProgress
	madeProgress = c.runStage(now, c.bottomSender) || madeProgress
	madeProgress = c.runStage(now, c.mshrStage) || madeProgress

	for _, bs := range c.bankStages {
		madeProgress = c.runStage(now, bs) || madeProgress
	}

	madeProgress = c.runStage(now, c.writeBuffer) || madeProgress
	madeProgress = c.runStage(now, c.dirStage) || madeProgress
	// putting pipeline here
	// madeProgress = c.runStage(now, c.pipeline) || madeProgress
	madeProgress = c.pipeline.Tick(now) || madeProgress
	madeProgress = c.runStage(now, c.topParser) || madeProgress

	return madeProgress
}

func (c *Cache) runStage(now akita.VTimeInSec, stage akita.Ticker) bool {
	madeProgress := false
	for i := 0; i < c.numReqPerCycle; i++ {
		madeProgress = stage.Tick(now) || madeProgress
	}
	return madeProgress
}

func (c *Cache) discardInflightTransactions(now akita.VTimeInSec) {
	sets := c.directory.GetSets()
	for _, set := range sets {
		for _, block := range set.Blocks {
			block.ReadCount = 0
			block.IsLocked = false
		}
	}

	c.dirStage.Reset(now)
	for _, bs := range c.bankStages {
		bs.Reset(now)
	}
	c.mshrStage.Reset(now)
	c.writeBuffer.Reset(now)

	clearPort(c.TopPort, now)

	c.topSender.Clear()

	c.inFlightTransactions = nil
}
