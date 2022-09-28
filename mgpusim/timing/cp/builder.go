package cp

import (
	"fmt"
	"math"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/mgpusim/protocol"
	"gitlab.com/akita/mgpusim/timing/cp/internal/dispatching"
	"gitlab.com/akita/mgpusim/timing/cp/internal/resource"
	"gitlab.com/akita/util"
	"gitlab.com/akita/util/akitaext"
	"gitlab.com/akita/util/tracing"
)

// Builder can build Command Processors
type Builder struct {
	freq            akita.Freq
	engine          akita.Engine
	visTracer       tracing.Tracer
	showProgressBar bool
	numDispatchers  int
	schedulingAlg   string

	partition         string
	customHSLpmdUnits uint64
}

// MakeBuilder creates a new builder with default configuration values.
func MakeBuilder() Builder {
	b := Builder{
		freq:           1 * akita.GHz,
		numDispatchers: 8,
	}
	return b
}

// WithVisTracer enables tracing for visualzation on the command processor and
// the dispatchers.
func (b Builder) WithVisTracer(tracer tracing.Tracer) Builder {
	b.visTracer = tracer
	return b
}

// WithEngine sets the even-driven simulation engine to use.
func (b Builder) WithEngine(engine akita.Engine) Builder {
	b.engine = engine
	return b
}

// WithFreq sets the frequency that the Command Processor works at.
func (b Builder) WithFreq(freq akita.Freq) Builder {
	b.freq = freq
	return b
}

// ShowProgressBar enables progress bar.
func (b Builder) ShowProgressBar() Builder {
	b.showProgressBar = true
	return b
}

// WithAlg sets the scheduling algorithm.
func (b Builder) WithAlg(alg string) Builder {
	b.schedulingAlg = alg
	return b
}

// WithAlg sets the scheduling algorithm.
func (b Builder) WithPartition(partition string) Builder {
	b.partition = partition
	return b
}

// WithAlg sets the scheduling algorithm.
func (b Builder) WithCustomHSLpmdUnits(customHSLpmdUnits uint64) Builder {
	b.customHSLpmdUnits = customHSLpmdUnits
	return b
}

// Build builds a new Command Processor
func (b Builder) Build(name string) *CommandProcessor {
	cp := new(CommandProcessor)
	cp.TickingComponent = akita.NewTickingComponent(name, b.engine, b.freq, cp)

	unlimited := math.MaxInt32
	cp.ToDriver = akita.NewLimitNumMsgPort(cp, 1, name+".ToDriver")
	cp.toDriverSender = akitaext.NewBufferedSender(
		cp.ToDriver, util.NewBuffer(unlimited))
	cp.ToDMA = akita.NewLimitNumMsgPort(cp, 1, name+".ToDMA")
	cp.toDMASender = akitaext.NewBufferedSender(
		cp.ToDMA, util.NewBuffer(unlimited))
	cp.ToCUs = akita.NewLimitNumMsgPort(cp, 1, name+".ToCUs")
	cp.toCUsSender = akitaext.NewBufferedSender(
		cp.ToCUs, util.NewBuffer(unlimited))
	cp.ToTLBs = akita.NewLimitNumMsgPort(cp, 1, name+".ToTLBs")
	cp.toTLBsSender = akitaext.NewBufferedSender(
		cp.ToTLBs, util.NewBuffer(unlimited))
	cp.ToMMUs = akita.NewLimitNumMsgPort(cp, 1, name+".ToMMUs")
	cp.toMMUsSender = akitaext.NewBufferedSender(
		cp.ToMMUs, util.NewBuffer(unlimited))

	cp.ToRDMA = akita.NewLimitNumMsgPort(cp, 1, name+".ToRDMA")
	cp.toRDMASender = akitaext.NewBufferedSender(
		cp.ToRDMA, util.NewBuffer(unlimited))
	cp.ToPMC = akita.NewLimitNumMsgPort(cp, 1, name+".ToPMC")
	cp.toPMCSender = akitaext.NewBufferedSender(
		cp.ToPMC, util.NewBuffer(unlimited))
	cp.ToAddressTranslators = akita.NewLimitNumMsgPort(cp, 1,
		name+".ToAddressTranslators")
	cp.toAddressTranslatorsSender = akitaext.NewBufferedSender(
		cp.ToAddressTranslators, util.NewBuffer(unlimited))
	cp.ToCaches = akita.NewLimitNumMsgPort(cp, 1, name+".ToCaches")
	cp.toCachesSender = akitaext.NewBufferedSender(
		cp.ToCaches, util.NewBuffer(unlimited))

	cp.ToRTU = akita.NewLimitNumMsgPort(cp, 1, name+".ToRTU")
	cp.toRTUSender = akitaext.NewBufferedSender(
		cp.ToRTU, util.NewBuffer(unlimited))

	cp.bottomKernelLaunchReqIDToTopReqMap =
		make(map[string]*protocol.LaunchKernelReq)
	cp.bottomMemCopyH2DReqIDToTopReqMap =
		make(map[string]*protocol.MemCopyH2DReq)
	cp.bottomMemCopyD2HReqIDToTopReqMap =
		make(map[string]*protocol.MemCopyD2HReq)

	b.buildDispatchers(cp)

	cp.customHSLpmdUnits = b.customHSLpmdUnits

	if b.visTracer != nil {
		tracing.CollectTrace(cp, b.visTracer)
	}

	cp.AverageQueueLength = make([]float64, 4)
	cp.PrevAverageQueueLength = make([]float64, 4)
	cp.NumAccesses = make([]uint64, 4)
	cp.PrevNumAccess = make([]uint64, 4)
	cp.NumMisses = make([]uint64, 4)
	cp.PrevNumMisses = make([]uint64, 4)
	cp.NumHits = make([]uint64, 4)
	cp.PrevNumHits = make([]uint64, 4)

	cp.RemoteMemAccessesPerWalk = make([]float64, 4)
	cp.PrevRemoteMemAccessesPerWalk = make([]float64, 4)
	cp.WalksEnqueued = make([]float64, 4)
	cp.PrevWalksEnqueued = make([]float64, 4)

	cp.IncomingReqs = make([]uint64, 4)
	cp.PrevIncomingReqs = make([]uint64, 4)
	cp.OutgoingReqs = make([]uint64, 4)
	cp.PrevOutgoingReqs = make([]uint64, 4)

	cp.StalledReqs = make([]uint64, 4)
	cp.PrevStalledReqs = make([]uint64, 4)

	cp.RemoteRequests = make([]uint64, 4)

	return cp
}

func (b *Builder) buildDispatchers(cp *CommandProcessor) {
	cuResourcePool := resource.NewCUResourcePool()
	builder := dispatching.MakeBuilder().
		WithCP(cp).
		WithCUResourcePool(cuResourcePool).
		WithDispatchingPort(cp.ToCUs).
		WithRespondingPort(cp.ToDriver).
		WithAlg(b.schedulingAlg).
		WithSchedulingPartition(b.partition)

	if b.showProgressBar {
		builder = builder.WithProgressBar()
	}

	for i := 0; i < b.numDispatchers; i++ {
		disp := builder.Build(fmt.Sprintf("%s.Dispatcher%d", cp.Name(), i))

		if b.visTracer != nil {
			tracing.CollectTrace(disp, b.visTracer)
		}

		cp.Dispatchers = append(cp.Dispatchers, disp)
	}
}
