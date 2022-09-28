package cp

import (
	"fmt"
	"math"
	"strings"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem"
	"gitlab.com/akita/mem/cache"
	"gitlab.com/akita/mem/idealmemcontroller"
	"gitlab.com/akita/mem/vm/tlb"
	"gitlab.com/akita/mgpusim/pagemigrationcontroller"
	"gitlab.com/akita/mgpusim/protocol"
	"gitlab.com/akita/mgpusim/rdma"
	"gitlab.com/akita/mgpusim/timing/cp/internal/dispatching"
	"gitlab.com/akita/mgpusim/timing/cp/internal/resource"
	"gitlab.com/akita/util/akitaext"
	"gitlab.com/akita/util/tracing"

	// "stats"
	// "github.com/montanaflynn/stats"
	"strconv"
)

type StatsFromComponents struct {
	// can also use ChipletId instead of Port
	TLBResponseCount       int
	AverageQueueLength     []float64
	PrevAverageQueueLength []float64
	NumAccesses            []uint64
	PrevNumAccess          []uint64
	NumMisses              []uint64
	PrevNumMisses          []uint64
	NumHits                []uint64
	PrevNumHits            []uint64
	PrevHitRate            float64
	StalledReqs            []uint64
	PrevStalledReqs        []uint64
	NumAccess              [4]uint64

	MMUResponseCount             int
	RemoteMemAccessesPerWalk     []float64
	PrevRemoteMemAccessesPerWalk []float64
	WalksEnqueued                []float64
	PrevWalksEnqueued            []float64

	RTUResponseCount int
	IncomingReqs     []uint64
	PrevIncomingReqs []uint64
	OutgoingReqs     []uint64
	PrevOutgoingReqs []uint64

	Hits           uint64
	Misses         uint64
	RemoteRequests []uint64

	InterleavingSize uint64
	Log2PageSize     uint64

	numSwitches        uint64
	switchTo4kHysteris bool

	switchL2TLBStriping bool

	statsCollectionInProgress bool
}

// CommandProcessor is an Akita component that is responsible for receiving
// requests from the driver and dispatch the requests to other parts of the
// GPU.
type CommandProcessor struct {
	*akita.TickingComponent
	StatsFromComponents

	Dispatchers        []dispatching.Dispatcher
	DMAEngine          akita.Port
	Driver             akita.Port
	TLBs               []akita.Port
	MMUs               []akita.Port
	CUs                []akita.Port
	AddressTranslators []akita.Port
	RDMA               akita.Port
	PMC                akita.Port
	L1VCaches          []akita.Port
	L1SCaches          []akita.Port
	L1ICaches          []akita.Port
	L2Caches           []akita.Port
	DRAMControllers    []*idealmemcontroller.Comp
	RTUs               []akita.Port

	ToDriver                   akita.Port
	toDriverSender             akitaext.BufferedSender
	ToDMA                      akita.Port
	toDMASender                akitaext.BufferedSender
	ToCUs                      akita.Port
	toCUsSender                akitaext.BufferedSender
	ToTLBs                     akita.Port
	toTLBsSender               akitaext.BufferedSender
	ToMMUs                     akita.Port
	toMMUsSender               akitaext.BufferedSender
	ToAddressTranslators       akita.Port
	toAddressTranslatorsSender akitaext.BufferedSender
	ToCaches                   akita.Port
	toCachesSender             akitaext.BufferedSender
	ToRDMA                     akita.Port
	toRDMASender               akitaext.BufferedSender
	ToRTU                      akita.Port
	toRTUSender                akitaext.BufferedSender
	ToPMC                      akita.Port
	toPMCSender                akitaext.BufferedSender

	currShootdownRequest *protocol.ShootDownCommand
	currFlushRequest     *protocol.FlushCommand

	numTLBs                      uint64
	numCUAck                     uint64
	numAddrTranslationFlushAck   uint64
	numAddrTranslationRestartAck uint64
	numTLBAck                    uint64
	numCacheACK                  uint64

	shootDownInProcess bool

	bottomKernelLaunchReqIDToTopReqMap map[string]*protocol.LaunchKernelReq
	bottomMemCopyH2DReqIDToTopReqMap   map[string]*protocol.MemCopyH2DReq
	bottomMemCopyD2HReqIDToTopReqMap   map[string]*protocol.MemCopyD2HReq

	customHSLpmdUnits uint64
}

// CUInterfaceForCP defines the interface that a CP requires from CU.
type CUInterfaceForCP interface {
	resource.DispatchableCU

	// ControlPort returns a port on the CU that the CP can send controlling
	// messages to.
	ControlPort() akita.Port
}

// RegisterCU allows the Command Processor to control the CU.
func (p *CommandProcessor) RegisterCU(cu CUInterfaceForCP) {
	p.CUs = append(p.CUs, cu.ControlPort())
	for _, d := range p.Dispatchers {
		d.RegisterCU(cu)
	}
}

func (p *CommandProcessor) SwitchL2TLBStriping(switchStriping bool) {
	p.switchL2TLBStriping = switchStriping
}

func (p *CommandProcessor) resetSwitchingStats() {
	p.TLBResponseCount = 0
	p.MMUResponseCount = 0
	p.InterleavingSize = 9
	p.numSwitches = 0
	for i := 0; i < 4; i++ {
		p.AverageQueueLength[i] = 0.0
		p.PrevAverageQueueLength[i] = 0.0

		p.NumAccesses[i] = 0
		p.PrevNumAccess[i] = 0

		p.NumMisses[i] = 0
		p.PrevNumMisses[i] = 0

		p.NumHits[i] = 0
		p.PrevNumHits[i] = 0

		p.RemoteMemAccessesPerWalk[i] = 0.0
		p.PrevRemoteMemAccessesPerWalk[i] = 0.0

		p.WalksEnqueued[i] = 0.0
		p.PrevWalksEnqueued[i] = 0.0

		p.IncomingReqs[i] = 0
		p.PrevIncomingReqs[i] = 0

		p.OutgoingReqs[i] = 0
		p.PrevOutgoingReqs[i] = 0

		p.StalledReqs[i] = 0
		p.PrevStalledReqs[i] = 0
	}
}

//Tick ticks
func (p *CommandProcessor) Tick(now akita.VTimeInSec) bool {
	madeProgress := false

	madeProgress = p.sendMsgsOut(now) || madeProgress
	madeProgress = p.tickDispatchers(now) || madeProgress
	madeProgress = p.processReqFromDriver(now) || madeProgress
	madeProgress = p.processRspFromInternal(now) || madeProgress

	return madeProgress
}

func (p *CommandProcessor) sendMsgsOut(now akita.VTimeInSec) bool {
	madeProgress := false

	madeProgress = p.sendMsgsOutFromPort(now, p.toDriverSender) || madeProgress
	madeProgress = p.sendMsgsOutFromPort(now, p.toDMASender) || madeProgress
	madeProgress = p.sendMsgsOutFromPort(now, p.toCUsSender) || madeProgress
	madeProgress = p.sendMsgsOutFromPort(now, p.toTLBsSender) || madeProgress
	madeProgress = p.sendMsgsOutFromPort(now, p.toMMUsSender) || madeProgress
	madeProgress = p.sendMsgsOutFromPort(now, p.toAddressTranslatorsSender) || madeProgress
	madeProgress = p.sendMsgsOutFromPort(now, p.toCachesSender) || madeProgress
	madeProgress = p.sendMsgsOutFromPort(now, p.toRDMASender) || madeProgress
	madeProgress = p.sendMsgsOutFromPort(now, p.toRTUSender) || madeProgress
	madeProgress = p.sendMsgsOutFromPort(now, p.toPMCSender) || madeProgress

	return madeProgress
}

func (p *CommandProcessor) sendMsgsOutFromPort(
	now akita.VTimeInSec,
	sender akitaext.BufferedSender,
) (madeProgress bool) {
	for {
		ok := sender.Tick(now)
		if ok {
			madeProgress = true
		} else {
			return madeProgress
		}
	}
}

func (p *CommandProcessor) tickDispatchers(
	now akita.VTimeInSec,
) (madeProgress bool) {
	for _, d := range p.Dispatchers {
		for i := 0; i < 32; i++ {
			madeProgress = d.Tick(now) || madeProgress
		}
	}

	return madeProgress
}

func (p *CommandProcessor) processReqFromDriver(now akita.VTimeInSec) bool {
	msg := p.ToDriver.Peek()
	if msg == nil {
		return false
	}

	switch req := msg.(type) {
	case *protocol.LaunchKernelReq:
		return p.processLaunchKernelReq(now, req)
	case *protocol.FlushCommand:
		return p.processFlushCommand(now, req)
	case *protocol.MemCopyD2HReq, *protocol.MemCopyH2DReq:
		return p.processMemCopyReq(now, req)
	case *protocol.RDMADrainCmdFromDriver:
		return p.processRDMADrainCmd(now, req)
	case *protocol.RDMARestartCmdFromDriver:
		return p.processRDMARestartCommand(now, req)
	case *protocol.ShootDownCommand:
		return p.processShootdownCommand(now, req)
	case *protocol.GPURestartReq:
		return p.processGPURestartReq(now, req)
	case *protocol.PageMigrationReqToCP:
		return p.processPageMigrationReq(now, req)
	}

	panic("never")
}

func (p *CommandProcessor) processRspFromInternal(now akita.VTimeInSec) bool {
	madeProgress := false

	madeProgress = p.processRspFromDMAs(now) || madeProgress
	madeProgress = p.processRspFromRDMAs(now) || madeProgress
	madeProgress = p.processRspFromCUs(now) || madeProgress
	madeProgress = p.processRspFromATs(now) || madeProgress
	madeProgress = p.processRspFromCaches(now) || madeProgress
	madeProgress = p.processRspFromTLBs(now) || madeProgress
	madeProgress = p.processRspFromMMUs(now) || madeProgress
	madeProgress = p.processRspFromRTUs(now) || madeProgress
	madeProgress = p.processRspFromPMC(now) || madeProgress

	return madeProgress
}

func (p *CommandProcessor) processRspFromDMAs(now akita.VTimeInSec) bool {
	msg := p.ToDMA.Peek()
	if msg == nil {
		return false
	}

	switch req := msg.(type) {
	case *protocol.MemCopyD2HReq, *protocol.MemCopyH2DReq:
		return p.processMemCopyRsp(now, req)
	}

	panic("never")
}

func (p *CommandProcessor) processRspFromRDMAs(now akita.VTimeInSec) bool {
	msg := p.ToRDMA.Peek()
	if msg == nil {
		return false
	}

	switch req := msg.(type) {
	case *rdma.DrainRsp:
		return p.processRDMADrainRsp(now, req)
	case *rdma.RestartRsp:
		return p.processRDMARestartRsp(now, req)
	}

	panic("never")
}

func (p *CommandProcessor) processRspFromCUs(now akita.VTimeInSec) bool {
	msg := p.ToCUs.Peek()
	if msg == nil {
		return false
	}

	switch req := msg.(type) {
	case *protocol.CUPipelineFlushRsp:
		return p.processCUPipelineFlushRsp(now, req)
	case *protocol.CUPipelineRestartRsp:
		return p.processCUPipelineRestartRsp(now, req)
	}

	return false
}

func (p *CommandProcessor) processRspFromCaches(now akita.VTimeInSec) bool {
	msg := p.ToCaches.Peek()
	if msg == nil {
		return false
	}

	switch req := msg.(type) {
	case *cache.FlushRsp:
		return p.processCacheFlushRsp(now, req)
	case *cache.RestartRsp:
		return p.processCacheRestartRsp(now, req)
	}

	panic("never")
}

func (p *CommandProcessor) processRspFromATs(now akita.VTimeInSec) bool {
	item := p.ToAddressTranslators.Peek()
	if item == nil {
		return false
	}

	msg := item.(*mem.ControlMsg)

	if p.numAddrTranslationFlushAck > 0 {
		return p.processAddressTranslatorFlushRsp(now, msg)
	} else if p.numAddrTranslationRestartAck > 0 {
		return p.processAddressTranslatorRestartRsp(now, msg)
	}

	panic("never")
}

func (p *CommandProcessor) processRspFromTLBs(now akita.VTimeInSec) bool {
	msg := p.ToTLBs.Peek()
	if msg == nil {
		return false
	}

	switch req := msg.(type) {
	case *tlb.TLBFlushRsp:
		return p.processTLBFlushRsp(now, req)
	case *tlb.TLBRestartRsp:
		return p.processTLBRestartRsp(now, req)
	case *akita.CollectedStatsMsg:
		p.processCollectedStats(now, req)
		p.ToTLBs.Retrieve(now)
		return true
	}

	panic("never")
}

func (p *CommandProcessor) processRspFromMMUs(now akita.VTimeInSec) bool {
	msg := p.ToMMUs.Peek()
	if msg == nil {
		return false
	}

	switch req := msg.(type) {
	case *akita.CollectedStatsMsg:
		p.processCollectedStats(now, req)
		p.ToMMUs.Retrieve(now)
		return true
	}

	panic("never")
}

// RTU is the only one capable of sending Trigger requests right now!
func (p *CommandProcessor) processRspFromRTUs(now akita.VTimeInSec) bool {
	msg := p.ToRTU.Peek()
	if msg == nil {
		return false
	}

	switch req := msg.(type) {
	case *akita.CollectedStatsMsg:
		p.processCollectedStats(now, req)
		p.ToRTU.Retrieve(now)
		return true
	case *akita.TriggerMsg:
		fmt.Println("Hello. Trigger message says hello")
		if !p.statsCollectionInProgress {
			p.requestStatsFromComponents(now, req)
		}
		p.ToRTU.Retrieve(now)
		return true
	}

	panic("never")
}

func (p *CommandProcessor) processRspFromPMC(now akita.VTimeInSec) bool {
	msg := p.ToPMC.Peek()
	if msg == nil {
		return false
	}

	switch req := msg.(type) {
	case *pagemigrationcontroller.PageMigrationRspFromPMC:
		return p.processPageMigrationRsp(now, req)
	}

	panic("never")
}

func (p *CommandProcessor) processLaunchKernelReq(
	now akita.VTimeInSec,
	req *protocol.LaunchKernelReq,
) bool {

	p.resetSwitchingStats()
	// p.InterleavingSize = 0
	// p.sendTLBIndexingSwitchMsg(now)
	p.sendTLBIndexingResetMsg(now)

	d := p.findAvailableDispatcher()

	if d == nil {
		return false
	}

	d.StartDispatching(req)
	p.ToDriver.Retrieve(now)

	tracing.TraceReqReceive(req, now, p)
	// tracing.TraceReqInitiate(&reqToBottom, now, p,
	// 	tracing.MsgIDAtReceiver(req, p))

	return true
}

func (p *CommandProcessor) findAvailableDispatcher() dispatching.Dispatcher {
	for _, d := range p.Dispatchers {
		if !d.IsDispatching() {
			return d
		}
	}

	return nil
}
func (p *CommandProcessor) processRDMADrainCmd(
	now akita.VTimeInSec,
	cmd *protocol.RDMADrainCmdFromDriver,
) bool {
	req := rdma.DrainReqBuilder{}.
		WithSendTime(now).
		WithSrc(p.ToRDMA).
		WithDst(p.RDMA).
		Build()

	p.toRDMASender.Send(req)
	p.ToDriver.Retrieve(now)

	return true
}

func (p *CommandProcessor) processRDMADrainRsp(
	now akita.VTimeInSec,
	rsp *rdma.DrainRsp,
) bool {
	req := protocol.NewRDMADrainRspToDriver(now, p.ToDriver, p.Driver)

	p.toDriverSender.Send(req)
	p.ToRDMA.Retrieve(now)

	return true
}

func (p *CommandProcessor) processShootdownCommand(
	now akita.VTimeInSec,
	cmd *protocol.ShootDownCommand,
) bool {
	if p.shootDownInProcess == true {
		return false
	}

	p.currShootdownRequest = cmd
	p.shootDownInProcess = true

	for i := 0; i < len(p.CUs); i++ {
		p.numCUAck++
		req := protocol.CUPipelineFlushReqBuilder{}.
			WithSendTime(now).
			WithSrc(p.ToCUs).
			WithDst(p.CUs[i]).
			Build()
		p.toCUsSender.Send(req)
	}

	p.ToDriver.Retrieve(now)

	return true
}

func (p *CommandProcessor) processCUPipelineFlushRsp(
	now akita.VTimeInSec,
	rsp *protocol.CUPipelineFlushRsp,
) bool {
	p.numCUAck--

	if p.numCUAck == 0 {
		for i := 0; i < len(p.AddressTranslators); i++ {
			req := mem.ControlMsgBuilder{}.
				WithSendTime(now).
				WithSrc(p.ToAddressTranslators).
				WithDst(p.AddressTranslators[i]).
				ToDiscardTransactions().
				Build()
			p.toAddressTranslatorsSender.Send(req)
			p.numAddrTranslationFlushAck++
		}
	}

	p.ToCUs.Retrieve(now)

	return true
}

func (p *CommandProcessor) processAddressTranslatorFlushRsp(
	now akita.VTimeInSec,
	msg *mem.ControlMsg,
) bool {
	p.numAddrTranslationFlushAck--

	if p.numAddrTranslationFlushAck == 0 {
		for _, port := range p.L1SCaches {
			p.flushAndResetL1Cache(now, port)
		}

		for _, port := range p.L1VCaches {
			p.flushAndResetL1Cache(now, port)
		}

		for _, port := range p.L1ICaches {
			p.flushAndResetL1Cache(now, port)
		}

		for _, port := range p.L2Caches {
			p.flushAndResetL2Cache(now, port)
		}
	}

	p.ToAddressTranslators.Retrieve(now)

	return true
}

func (p *CommandProcessor) flushAndResetL1Cache(
	now akita.VTimeInSec,
	port akita.Port,
) {
	req := cache.FlushReqBuilder{}.
		WithSendTime(now).
		WithSrc(p.ToCaches).
		WithDst(port).
		PauseAfterFlushing().
		DiscardInflight().
		InvalidateAllCacheLines().
		Build()

	p.toCachesSender.Send(req)
	p.numCacheACK++
}

func (p *CommandProcessor) flushAndResetL2Cache(now akita.VTimeInSec, port akita.Port) {
	req := cache.FlushReqBuilder{}.
		WithSendTime(now).
		WithSrc(p.ToCaches).
		WithDst(port).
		PauseAfterFlushing().
		DiscardInflight().
		InvalidateAllCacheLines().
		Build()

	p.toCachesSender.Send(req)
	p.numCacheACK++
}

func (p *CommandProcessor) processCacheFlushRsp(
	now akita.VTimeInSec,
	rsp *cache.FlushRsp,
) bool {
	p.numCacheACK--
	if p.numCacheACK == 0 {
		if p.shootDownInProcess {
			return p.processCacheFlushCausedByTLBShootdown(now, rsp)
		}
		return p.processRegularCacheFlush(now, rsp)
	}

	p.ToCaches.Retrieve(now)
	return true
}

func (p *CommandProcessor) processRegularCacheFlush(
	now akita.VTimeInSec,
	flushRsp *cache.FlushRsp,
) bool {
	p.currFlushRequest.Src, p.currFlushRequest.Dst =
		p.currFlushRequest.Dst, p.currFlushRequest.Src
	p.currFlushRequest.SendTime = now

	p.toDriverSender.Send(p.currFlushRequest)

	p.ToCaches.Retrieve(now)

	return true
}

func (p *CommandProcessor) processCacheFlushCausedByTLBShootdown(
	now akita.VTimeInSec,
	flushRsp *cache.FlushRsp,
) bool {
	for i := 0; i < len(p.TLBs); i++ {
		shootDownCmd := p.currShootdownRequest
		req := tlb.TLBFlushReqBuilder{}.
			WithSendTime(now).
			WithSrc(p.ToTLBs).
			WithDst(p.TLBs[i]).
			WithPID(shootDownCmd.PID).
			WithVAddrs(shootDownCmd.VAddr).
			Build()

		p.toTLBsSender.Send(req)
		p.numTLBAck++
	}

	p.ToCaches.Retrieve(now)
	return true
}

func (p *CommandProcessor) processTLBFlushRsp(
	now akita.VTimeInSec,
	rsp *tlb.TLBFlushRsp,
) bool {
	p.numTLBAck--

	if p.numTLBAck == 0 {
		req := protocol.NewShootdownCompleteRsp(now, p.ToDriver, p.Driver)
		p.toDriverSender.Send(req)

		p.shootDownInProcess = false
	}

	p.ToTLBs.Retrieve(now)

	return true
}

func (p *CommandProcessor) processRDMARestartCommand(
	now akita.VTimeInSec,
	cmd *protocol.RDMARestartCmdFromDriver,
) bool {
	req := rdma.RestartReqBuilder{}.
		WithSrc(p.ToRDMA).
		WithDst(p.RDMA).
		WithSendTime(now).
		Build()

	p.toRDMASender.Send(req)

	p.ToDriver.Retrieve(now)

	return true
}

func (p *CommandProcessor) processRDMARestartRsp(now akita.VTimeInSec, rsp *rdma.RestartRsp) bool {
	req := protocol.NewRDMARestartRspToDriver(now, p.ToDriver, p.Driver)
	p.toDriverSender.Send(req)
	p.ToRDMA.Retrieve(now)

	return true
}

func (p *CommandProcessor) processGPURestartReq(
	now akita.VTimeInSec,
	cmd *protocol.GPURestartReq,
) bool {
	for _, port := range p.L2Caches {
		p.restartCache(now, port)
	}
	for _, port := range p.L1ICaches {
		p.restartCache(now, port)
	}
	for _, port := range p.L1SCaches {
		p.restartCache(now, port)
	}

	for _, port := range p.L1VCaches {
		p.restartCache(now, port)
	}

	p.ToDriver.Retrieve(now)

	return true
}

func (p *CommandProcessor) restartCache(now akita.VTimeInSec, port akita.Port) {
	req := cache.RestartReqBuilder{}.
		WithSendTime(now).
		WithSrc(p.ToCaches).
		WithDst(port).
		Build()

	p.toCachesSender.Send(req)

	p.numCacheACK++
}

func (p *CommandProcessor) processCacheRestartRsp(
	now akita.VTimeInSec,
	rsp *cache.RestartRsp,
) bool {
	p.numCacheACK--
	if p.numCacheACK == 0 {
		for i := 0; i < len(p.TLBs); i++ {
			p.numTLBAck++

			req := tlb.TLBRestartReqBuilder{}.
				WithSendTime(now).
				WithSrc(p.ToTLBs).
				WithDst(p.TLBs[i]).
				Build()
			p.toTLBsSender.Send(req)
		}
	}

	p.ToCaches.Retrieve(now)

	return true
}

func (p *CommandProcessor) processTLBRestartRsp(
	now akita.VTimeInSec,
	rsp *tlb.TLBRestartRsp,
) bool {
	p.numTLBAck--

	if p.numTLBAck == 0 {
		for i := 0; i < len(p.AddressTranslators); i++ {
			req := mem.ControlMsgBuilder{}.
				WithSendTime(now).
				WithSrc(p.ToAddressTranslators).
				WithDst(p.AddressTranslators[i]).
				ToRestart().
				Build()
			p.toAddressTranslatorsSender.Send(req)

			// fmt.Printf("Restarting %s\n", p.AddressTranslators[i].Name())

			p.numAddrTranslationRestartAck++
		}
	}

	p.ToTLBs.Retrieve(now)

	return true
}

func (p *CommandProcessor) processAddressTranslatorRestartRsp(
	now akita.VTimeInSec,
	rsp *mem.ControlMsg,
) bool {
	p.numAddrTranslationRestartAck--

	if p.numAddrTranslationRestartAck == 0 {
		for i := 0; i < len(p.CUs); i++ {
			req := protocol.CUPipelineRestartReqBuilder{}.
				WithSendTime(now).
				WithSrc(p.ToCUs).
				WithDst(p.CUs[i]).
				Build()
			p.toCUsSender.Send(req)

			p.numCUAck++
		}
	}

	p.ToAddressTranslators.Retrieve(now)

	return true
}

func (p *CommandProcessor) processCUPipelineRestartRsp(
	now akita.VTimeInSec,
	rsp *protocol.CUPipelineRestartRsp,
) bool {
	p.numCUAck--

	if p.numCUAck == 0 {
		rsp := protocol.NewGPURestartRsp(now, p.ToDriver, p.Driver)
		p.toDriverSender.Send(rsp)
	}

	p.ToCUs.Retrieve(now)

	return true
}

func (p *CommandProcessor) processPageMigrationReq(
	now akita.VTimeInSec,
	cmd *protocol.PageMigrationReqToCP,
) bool {
	req := pagemigrationcontroller.PageMigrationReqToPMCBuilder{}.
		WithSendTime(now).
		WithSrc(p.ToPMC).
		WithDst(p.PMC).
		WithPageSize(cmd.PageSize).
		WithPMCPortOfRemoteGPU(cmd.DestinationPMCPort).
		WithReadFrom(cmd.ToReadFromPhysicalAddress).
		WithWriteTo(cmd.ToWriteToPhysicalAddress).
		Build()
	p.toPMCSender.Send(req)

	p.ToDriver.Retrieve(now)

	return true
}

func (p *CommandProcessor) processPageMigrationRsp(
	now akita.VTimeInSec,
	rsp *pagemigrationcontroller.PageMigrationRspFromPMC,
) bool {
	req := protocol.NewPageMigrationRspToDriver(now, p.ToDriver, p.Driver)

	p.toDriverSender.Send(req)

	p.ToPMC.Retrieve(now)

	return true
}

func (p *CommandProcessor) processFlushCommand(
	now akita.VTimeInSec,
	cmd *protocol.FlushCommand,
) bool {
	if p.numCacheACK > 0 {
		return false
	}

	for _, port := range p.L1ICaches {
		p.flushCache(now, port)
	}

	for _, port := range p.L1SCaches {
		p.flushCache(now, port)
	}

	for _, port := range p.L1VCaches {
		p.flushCache(now, port)
	}

	for _, port := range p.L2Caches {
		p.flushCache(now, port)
	}

	p.currFlushRequest = cmd
	if p.numCacheACK == 0 {
		p.currFlushRequest.Src, p.currFlushRequest.Dst =
			p.currFlushRequest.Dst, p.currFlushRequest.Src
		p.currFlushRequest.SendTime = now
		p.toDriverSender.Send(p.currFlushRequest)
	}

	p.ToDriver.Retrieve(now)

	return true
}

func (p *CommandProcessor) flushCache(now akita.VTimeInSec, port akita.Port) {
	flushReq := cache.FlushReqBuilder{}.
		WithSendTime(now).
		WithSrc(p.ToCaches).
		WithDst(port).
		Build()
	p.toCachesSender.Send(flushReq)
	p.numCacheACK++
}

func (p *CommandProcessor) cloneMemCopyH2DReq(
	req *protocol.MemCopyH2DReq,
) *protocol.MemCopyH2DReq {
	cloned := *req
	cloned.ID = akita.GetIDGenerator().Generate()
	p.bottomMemCopyH2DReqIDToTopReqMap[cloned.ID] = req
	return &cloned
}

func (p *CommandProcessor) cloneMemCopyD2HReq(
	req *protocol.MemCopyD2HReq,
) *protocol.MemCopyD2HReq {
	cloned := *req
	cloned.ID = akita.GetIDGenerator().Generate()
	p.bottomMemCopyD2HReqIDToTopReqMap[cloned.ID] = req
	return &cloned
}

func (p *CommandProcessor) processMemCopyReq(
	now akita.VTimeInSec,
	req akita.Msg,
) bool {
	var cloned akita.Msg
	switch req := req.(type) {
	case *protocol.MemCopyH2DReq:
		cloned = p.cloneMemCopyH2DReq(req)
	case *protocol.MemCopyD2HReq:
		cloned = p.cloneMemCopyD2HReq(req)
	default:
		panic("unknown type")
	}

	cloned.Meta().Dst = p.DMAEngine
	cloned.Meta().Src = p.ToDMA
	cloned.Meta().SendTime = now

	p.toDMASender.Send(cloned)
	p.ToDriver.Retrieve(now)

	tracing.TraceReqReceive(req, now, p)
	tracing.TraceReqInitiate(cloned, now, p, tracing.MsgIDAtReceiver(req, p))

	return true
}

func (p *CommandProcessor) findAndRemoveOriginalMemCopyRequest(
	rsp akita.Msg,
) akita.Msg {
	switch rsp := rsp.(type) {
	case *protocol.MemCopyH2DReq:
		origionalReq := p.bottomMemCopyH2DReqIDToTopReqMap[rsp.ID]
		delete(p.bottomMemCopyH2DReqIDToTopReqMap, rsp.ID)
		return origionalReq
	case *protocol.MemCopyD2HReq:
		originalReq := p.bottomMemCopyD2HReqIDToTopReqMap[rsp.ID]
		delete(p.bottomMemCopyD2HReqIDToTopReqMap, rsp.ID)
		return originalReq
	default:
		panic("unknown type")
	}
}

func (p *CommandProcessor) processMemCopyRsp(
	now akita.VTimeInSec,
	req akita.Msg,
) bool {
	originalReq := p.findAndRemoveOriginalMemCopyRequest(req)
	originalReq.Meta().Dst = p.Driver
	originalReq.Meta().Src = p.ToDriver
	originalReq.Meta().SendTime = now
	p.toDriverSender.Send(originalReq)
	p.ToDMA.Retrieve(now)

	tracing.TraceReqComplete(originalReq, now, p)
	tracing.TraceReqFinalize(req, now, p)

	return true
}

func (p *CommandProcessor) sendTLBIndexingResetMsg(now akita.VTimeInSec) bool {
	// go over the RTU  ports and send TLBIndexingSwitchMsgs
	if !p.switchL2TLBStriping {
		return true
	}
	if p.numSwitches >= 10 {
		return true
	}
	switchTo := akita.TLBIndexingSwitchHSL
	interleaving := p.customHSLpmdUnits

	for i := 0; i < len(p.RTUs); i++ {
		req := akita.TLBIndexingSwitchMsgBuilder{}.
			WithSendTime(now).
			WithSrc(p.ToRTU).
			WithDst(p.RTUs[i]).
			WithInterleaving(interleaving).
			WithSwitchTo(switchTo).
			Build()
		p.toRTUSender.Send(req)
	}
	for i := 0; i < len(p.TLBs); i++ {
		req := akita.TLBIndexingSwitchMsgBuilder{}.
			WithSendTime(now).
			WithSrc(p.ToTLBs).
			WithDst(p.TLBs[i]).
			WithInterleaving(interleaving).
			WithSwitchTo(switchTo).
			Build()
		p.toTLBsSender.Send(req)
	}
	for i := 0; i < len(p.MMUs); i++ {
		req := akita.TLBIndexingSwitchMsgBuilder{}.
			WithSendTime(now).
			WithSrc(p.ToMMUs).
			WithDst(p.MMUs[i]).
			WithInterleaving(interleaving).
			WithSwitchTo(switchTo).
			Build()
		p.toMMUsSender.Send(req)
	}
	p.numSwitches++
	return true
}

func (p *CommandProcessor) sendTLBIndexingSwitchMsg(now akita.VTimeInSec) bool {
	// go over the RTU  ports and send TLBIndexingSwitchMsgs
	if !p.switchL2TLBStriping {
		return true
	}
	if p.numSwitches >= 10 {
		return true
	}
	switchTo := akita.TLBIndexingSwitch4K
	interleaving := p.customHSLpmdUnits

	for i := 0; i < len(p.RTUs); i++ {
		req := akita.TLBIndexingSwitchMsgBuilder{}.
			WithSendTime(now).
			WithSrc(p.ToRTU).
			WithDst(p.RTUs[i]).
			WithInterleaving(interleaving).
			WithSwitchTo(switchTo).
			Build()
		p.toRTUSender.Send(req)
	}
	for i := 0; i < len(p.TLBs); i++ {
		req := akita.TLBIndexingSwitchMsgBuilder{}.
			WithSendTime(now).
			WithSrc(p.ToTLBs).
			WithDst(p.TLBs[i]).
			WithInterleaving(interleaving).
			WithSwitchTo(switchTo).
			Build()
		p.toTLBsSender.Send(req)
	}
	for i := 0; i < len(p.MMUs); i++ {
		req := akita.TLBIndexingSwitchMsgBuilder{}.
			WithSendTime(now).
			WithSrc(p.ToMMUs).
			WithDst(p.MMUs[i]).
			WithInterleaving(interleaving).
			WithSwitchTo(switchTo).
			Build()
		p.toMMUsSender.Send(req)
	}
	p.numSwitches++
	return true
}

func getChipletNum(component string) int {
	// fmt.Println(component, component[14:15])
	chiplet, _ := strconv.Atoi(component[14:15])
	return chiplet
}

func sum(A []uint64) (s uint64) {
	s = 0
	for _, a := range A {
		s += uint64(a)
	}
	return s
}

func maxUint64(A []uint64) (m uint64) {
	m = A[0]
	for _, a := range A {
		if a > m {
			m = a
		}
	}
	return m
}

func sumFloat64(A []float64) (s float64) {
	s = 0
	for _, a := range A {
		s += a
	}
	return s
}

func min(A []float64) (m float64) {
	m = A[0]
	for _, a := range A {
		if a < m {
			m = a
		}
	}
	return m
}

func max(A []float64) (m float64) {
	m = A[0]
	for _, a := range A {
		if a > m {
			m = a
		}
	}
	return m
}

func div(a, b float64) (m float64) {
	if b == 0 {
		return 0
	}
	return a / b
}

func (p *CommandProcessor) resetCollectedStats() {
	p.Hits = 0
	p.Misses = 0
	for i := 0; i < 4; i++ {
		p.RemoteRequests[i] = 0
	}
}

// Returns true if any element in array is further from the mean by threshold.
// Array element with index == exception is not included in any calculation.
// Set exception = -1 if you don't want any exceptions.
func getImbalance(array []uint64, count float64) bool {
	sum := float64(0.0)
	imbalanced := false
	for i := 0; i < len(array); i++ {
		sum += float64(array[i])
	}
	for i := 0; i < len(array); i++ {
		if math.Abs((float64(array[i]))/sum) > 0.8 {
			imbalanced = true
		}
	}

	return imbalanced
}

func (p *CommandProcessor) processCollectedStats(now akita.VTimeInSec,
	req *akita.CollectedStatsMsg) {

	srcName := req.Src.Name()
	chipletNum := getChipletNum(req.Src.Name())

	if strings.Contains(srcName, "L2TLB") {
		p.TLBResponseCount++
		p.Hits += req.NumHits
		p.Misses += req.NumMisses
		p.NumAccesses[chipletNum] = req.NumHits + req.NumMisses
	} else if strings.Contains(srcName, "MMU") {
		p.MMUResponseCount++
	} else if strings.Contains(srcName, "RTU") {
		p.RTUResponseCount++
		p.RemoteRequests[chipletNum] += req.IncomingReqs
		p.IncomingReqs[chipletNum] = req.IncomingReqs
	}

	if p.TLBResponseCount == 4 && p.MMUResponseCount == 4 && p.RTUResponseCount == 4 {
		p.statsCollectionInProgress = false

		// fmt.Println("Incoming ", p.IncomingReqs[0], p.IncomingReqs[1], p.IncomingReqs[2], p.IncomingReqs[3])
		fmt.Println("L2 TLB Requests ", p.NumAccesses[0], p.NumAccesses[1], p.NumAccesses[2], p.NumAccesses[3])
		fmt.Println("L2 TLB Hit Rate ", float64(p.Hits)/float64(p.Hits+p.Misses))
		fmt.Println("Switching: Num Incoming ", p.IncomingReqs[0], p.IncomingReqs[1], p.IncomingReqs[2], p.IncomingReqs[3])

		switchTo4k := false
		// Process all the data here. Arbitrary thresholds of 0.5
		if (float64(p.Hits) / float64(p.Hits+p.Misses)) > 0.9 {
			if getImbalance(p.IncomingReqs, 4.0) == true {
				switchTo4k = true
			}
		} else {
			switchTo4k = false
		}

		// if switchTo4k {
		// 	p.sendTLBIndexingSwitchMsg(now)
		// }

		// // Take action if necessary
		if switchTo4k && p.switchTo4kHysteris {
			p.sendTLBIndexingSwitchMsg(now)
		} else if switchTo4k {
			p.switchTo4kHysteris = true
		} else {
			p.switchTo4kHysteris = false
		}

		p.resetCollectedStats()
	}

}

func (p *CommandProcessor) requestStatsFromComponents(now akita.VTimeInSec,
	req *akita.TriggerMsg) bool {

	if p.statsCollectionInProgress {
		panic("stats collection in progress!")
	}
	p.statsCollectionInProgress = true
	for i := 0; i < len(p.TLBs); i++ {
		req := akita.SendStatsMsgBuilder{}.
			WithSendTime(now).
			WithSrc(p.ToTLBs).
			WithDst(p.TLBs[i]).
			Build()
		p.toTLBsSender.Send(req)
	}
	for i := 0; i < len(p.MMUs); i++ {
		req := akita.SendStatsMsgBuilder{}.
			WithSendTime(now).
			WithSrc(p.ToMMUs).
			WithDst(p.MMUs[i]).
			Build()
		p.toMMUsSender.Send(req)
	}
	for i := 0; i < len(p.RTUs); i++ {
		req := akita.SendStatsMsgBuilder{}.
			WithSendTime(now).
			WithSrc(p.ToRTU).
			WithDst(p.RTUs[i]).
			Build()
		p.toRTUSender.Send(req)
	}
	p.TLBResponseCount = 0
	p.MMUResponseCount = 0
	p.RTUResponseCount = 0
	return true
}
