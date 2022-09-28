package writeback

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem"
	"gitlab.com/akita/util/tracing"
)

type topParser struct {
	cache *Cache
}

func (p *topParser) Tick(now akita.VTimeInSec) bool {
	if p.cache.state != cacheStateRunning {
		return true // return false
	}

	req := p.cache.TopPort.Peek()
	if req == nil {
		return true // return false
	}

	// instead wrap in pipeline object and push into pipeline
	if !p.cache.pipeline.CanAccept() {
		return true // return false
	}
	// if !p.cache.dirStageBuffer.CanPush() {
	// 	return true // return false
	// }

	trans := &transaction{}
	switch req := req.(type) {
	// case *mem.DataReadyRsp:
	// 	// trans.read = req
	// 	err := p.cache.BottomPort.Recv(req)
	// 	if err != nil {
	// 		return true
	// 	} else {
	// 		p.cache.TopPort.Retrieve(now)
	// 		return true
	// 	}
	case *mem.ReadReq:
		trans.read = req
	case *mem.WriteReq:
		trans.write = req
	}
	// pipeline
	pipelineItem := cachePipelineItem{
		taskID: akita.GetIDGenerator().Generate(),
		trans:  trans,
	}
	p.cache.pipeline.Accept(now, pipelineItem)
	// p.cache.dirStageBuffer.Push(trans)

	p.cache.inFlightTransactions = append(p.cache.inFlightTransactions, trans)

	tracing.TraceReqReceive(req, now, p.cache)

	p.cache.TopPort.Retrieve(now)

	return true
}
