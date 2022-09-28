// Package runner defines how default benchmark samples are executed.
package runner

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	// Enable profiling
	_ "net/http/pprof"
	"strconv"
	"strings"
	"sync"

	"github.com/tebeka/atexit"
	"gitlab.com/akita/akita"

	// ram "gitlab.com/akita/mem/dram"
	"gitlab.com/akita/mem/idealmemcontroller"
	"gitlab.com/akita/mem/vm/mmu"
	"gitlab.com/akita/mem/vm/tlb"
	"gitlab.com/akita/mgpusim/benchmarks"
	"gitlab.com/akita/mgpusim/driver"
	"gitlab.com/akita/mgpusim/platform"
	"gitlab.com/akita/mgpusim/rdma"
	"gitlab.com/akita/mgpusim/remotetranslation"
	"gitlab.com/akita/util/tracing"
)

var timingFlag = flag.Bool("timing", false, "Run detailed timing simulation.")
var maxInstCount = flag.Uint64("max-inst", 0,
	"Terminate the simulation after the given number of instructions is retired.")
var parallelFlag = flag.Bool("parallel", false,
	"Run the simulation in parallel.")
var isaDebug = flag.Bool("debug-isa", false, "Generate the ISA debugging file.")
var visTracing = flag.Bool("trace-vis", false,
	"Generate trace for visualization purposes.")
var l2tlbsqlTracingFlag = flag.Bool("trace-l2sql", false,
	"Generate trace from L2 TLB into a new MySql database.")
var verifyFlag = flag.Bool("verify", false, "Verify the emulation result.")
var memTracing = flag.Bool("trace-mem", false, "Generate memory trace")
var tlbTracing = flag.Bool("trace-tlb", false, "Generate memory trace")
var disableProgressBar = flag.Bool("no-progress-bar", false,
	"Disables the progress bar")
var instCountReportFlag = flag.Bool("report-inst-count", false,
	"Report the number of instructions executed in each compute unit.")
var cacheLatencyReportFlag = flag.Bool("report-cache-latency", false,
	"Report the average cache latency.")
var cacheHitRateReportFlag = flag.Bool("report-cache-hit-rate", false,
	"Report the cache hit rate of each cache.")
var tlbLatencyReportFlag = flag.Bool("report-tlb-latency", false,
	"Report the average tlb latency.")
var tlbConditionalStatsReportFlag = flag.Bool("report-tlb-conditional", false,
	"Report the tlb stats about how much is remote and how much is local.")
var tlbHitRateReportFlag = flag.Bool("report-tlb-hit-rate", false,
	"Report the tlb hit rate of each tlb.")
var tlbCoalesceReportFlag = flag.Bool("report-tlb-coalesce", false,
	"Report the average tlb coalesce.")
var rtuCoalesceReportFlag = flag.Bool("report-rtu-coalesce", false,
	"Report the average amount of coalescing possible at the RTU.")
var pagewalkLatencyReportFlag = flag.Bool("report-pagewalk-latency", false,
	"Report the average page walk latency.")
var mmuConditionalStatsReportFlag = flag.Bool("report-mmu-conditional", false,
	"Report the mmu stats about how much is remote and how much is local.")
var activeWalkerReportFlag = flag.Bool("report-active-walkers", false,
	"Report the average acitve page walker count.")
var dramLatencyReportFlag = flag.Bool("report-dram-latency", false,
	"Report the average dram latency.")
var pwcHitRateReportFlag = flag.Bool("report-pwc-hit-rate", false,
	"Report the pwc hit rate of each level.")
var rdmaTransactionCountReportFlag = flag.Bool("report-rdma-transaction-count",
	false, "Report the number of transactions going through the RDMA engines.")
var cdmaTransactionCountReportFlag = flag.Bool("report-cdma-transaction-count",
	false, "Report the number of transactions through the chip RDMA engines.")
var rtuTransactionCountReportFlag = flag.Bool("report-rtu-transaction-count",
	false, "Report the number of transactions through the remote trans. engines.")
var dramTransactionCountReportFlag = flag.Bool("report-dram-transaction-count",
	false, "Report the number of transactions accessing the DRAMs.")
var translationReqLatencyFlag = flag.Bool("report-translation-request-latency",
	false, "Report the latency of  requests thrugh the remote trans. engines.")
var addressTranslatorLatencyFlag = flag.Bool("report-address-translator-latency",
	false, "Report the latency of  requests thrugh the addr. trans.")
var l2tlbQueueingImbalanceFlag = flag.Bool("report-l2tlb-queueing-imbalance",
	false, "Report the imbalnce of  requests thrugh the L2 TLB .")
var pageWalkerImbalanceFlag = flag.Bool("report-pagewalker-imbalance",
	false, "Report the imbalnce of  use of the Page walker.")
var entropyFlag = flag.Bool("report-entropy",
	false, "Report the entropy of requests are various L2 TLBs")
var referenceTracingFlag = flag.Bool("report-reference-tracing",
	false, "Report the number of references to differnet addresses")
var gpuFlag = flag.String("gpus", "",
	"The GPUs to use, use a format like 1,2,3,4. By default, GPU 1 is used.")
var unifiedGPUFlag = flag.String("unified-gpus", "",
	`Run multi-GPU benchmark in a unified mode.
Use a format like 1,2,3,4. Cannot coexist with -gpus.`)
var useUnifiedMemoryFlag = flag.Bool("use-unified-memory", false,
	"Run benchmark with Unified Memory or not")
var reportAll = flag.Bool("report-all", false, "Report all metrics to .csv file.")
var reportEssential = flag.Bool("report-essential", false, "Report essential metrics to .csv file.")
var filenameFlag = flag.String("metric-file-name", "metrics",
	"Modify the name of the output csv file.")
var platformType = flag.String("platform-type", "mcmgpu",
	"Name of the platform to build (ideal, mcmgpu, distlb, etc)")
var schedulingAlg = flag.String("scheduling", "chiplet-aware",
	"Name of the scheduling algorithm to use (chiplet-aware, greedy, round-robin, lasp)")
var schedulingPartition = flag.String("sched-partition", "Xdiv",
	"Name of the scheduling partition to use (Xdiv4, Ydiv4, Xblk4, Yblk4)")
var useCoalescingTLBPort = flag.Bool("use-coalescing-tlb-port", false,
	"Specify if the L2TLB port should be coalescing.")
var useCoalescingRTU = flag.Bool("use-coalescing-rtu", false,
	"Specify if the RTU should be coalescing.")
var l2TlbStriping = flag.Uint64("l2-tlb-striping", 512,
	"Specify the granularity at which entries are striped across the L2 TLB.")
var memAllocatorType = flag.String("mem-allocator-type", "interleaved",
	"Specify the type of memory allocator.")
var useLASPMemoryAllocFlag = flag.Bool("use-lasp-mem-alloc", false,
	"Set this if we are going to use the LASP memory allocation API in applications.")
var useLASPHSLMemoryAllocFlag = flag.Bool("use-lasp-hsl-mem-alloc", false,
	"Set this if we are going to use HSL directed rearranged of data.")
var useSwitching = flag.Bool("switch-tlb-striping", false,
	"Specify if the L2 TLB striping should be switched")
var ptCaching = flag.Bool("use-pt-caching", false,
	"Specify if page table pages should be cached locally")
var log2PageSize = flag.Uint64("log2-page-size", 12,
	"Specify the page size")
var useCustomHSL = flag.Bool("use-custom-hsl", false,
	"Specify if to use custom HSL as defined in program")
var customHSL = flag.Uint64("custom-hsl", 1,
	"Specify the value of custom HSL directly to the builder")

type verificationPreEnablingBenchmark interface {
	benchmarks.Benchmark

	EnableVerification()
}

type instCountTracer struct {
	tracer *instTracer
	cu     akita.Component
}

type cacheLatencyTracer struct {
	tracer *tracing.AverageTimeTracer
	cache  akita.Component
}

type cacheHitRateTracer struct {
	tracer *tracing.StepCountTracer
	cache  akita.Component
}

type TLBLatencyTracer struct {
	tracer *tracing.AverageTimeTracer
	tlb    tracing.NamedHookable
}

type TLBPipelineLatencyTracer struct {
	tracer   *tracing.AverageTimeTracer
	pipeline tracing.NamedHookable
}

type L2PipelineLatencyTracer struct {
	tracer   *tracing.AverageTimeTracer
	pipeline tracing.NamedHookable
}

type ConditionalRemoteTLBLatencyTracer struct {
	tracer *tracing.ConditionalAverageTimeTracer
	tlb    akita.Component
}

type AddressTranslatorLatencyTracer struct {
	tracer            *tracing.ConditionalAverageTimeTracer
	addressTranslator akita.Named
}

type MMUMemoryLatencyTracer struct {
	tracer *tracing.ConditionalAverageTimeTracer
	mmu    akita.Component
}

type TLBHitRateTracer struct {
	tracer *tracing.StepCountTracer
	tlb    tracing.NamedHookable
}

type CDMAAccessTracer struct {
	tracer     *tracing.StepCountTracer
	cdmaEngine akita.Component
}

type PWCHitRateTracer struct {
	tracer *tracing.StepCountTracer
	mmu    akita.Component
}

type PageWalkLatencyTracer struct {
	tracer *tracing.AverageTimeTracer
	mmu    akita.Component
}

type ActivePageWalkerTracer struct {
	tracer *tracing.AverageCountTracer
	mmu    akita.Component
}

type RemoteTLBLatencyTracer struct {
	tracer *tracing.AverageTimeTracer
	rtu    akita.Component
}

type RTUAccessTracer struct {
	tracer *tracing.StepCountTracer
	rtu    remotetranslation.RemoteTranslationUnit
}

type DRAMLatencyTracer struct {
	tracer *tracing.AverageTimeTracer
	dram   akita.Component
}

type dramTransactionCountTracer struct {
	tracer *tracing.AverageTimeTracer
	dram   *idealmemcontroller.Comp
	// dram *ram.MemController
}

type rdmaTransactionCountTracer struct {
	outgoingTracer *tracing.AverageTimeTracer
	incomingTracer *tracing.AverageTimeTracer
	rdmaEngine     *rdma.Engine
}

type rtuTransactionCountTracer struct {
	outgoingTracer *tracing.AverageTimeTracer
	incomingTracer *tracing.AverageTimeTracer
	rtuUnit        remotetranslation.RemoteTranslationUnit
}

type TLBQueueImbalanceTracer struct {
	tracer *tlb.GlobalTLBQueueingTracer
	tlb    tlb.L2TLB
}

type PageWalkerImbalanceTracer struct {
	tracer *mmu.GlobalPageWalkerOccupancyTracer
	mmu    *mmu.MMUImpl
}

type EntropyTracer struct {
	tracer *tlb.EntropyTracer
	tlb    tlb.L2TLB
}

type RemoteReferenceCountTracer struct {
	tracer  *tracing.ReferenceTracer
	rtuUnit remotetranslation.RemoteTranslationUnit
}

type TLBSetMissTracer struct {
	tracer *tracing.ReferenceTracer
	tlb    tlb.L2TLB
}

type TLBMSHRStallTracer struct {
	tracer *tracing.TotalTimeTracer
	tlb    tlb.L2TLB
}

type TLBReqStallTracer struct {
	tracer *tracing.StepCountTracer
	tlb    tlb.L2TLB
}

// Runner is a class that helps running the benchmarks in the official samples.
type Runner struct {
	Engine                             akita.Engine
	GPUDriver                          *driver.Driver
	maxInstStopper                     *instTracer
	KernelTimeCounter                  *tracing.BusyTimeTracer
	PerGPUKernelTimeCounter            []*tracing.BusyTimeTracer
	InstCountTracers                   []instCountTracer
	CacheLatencyTracers                []cacheLatencyTracer
	TLBLatencyTracers                  []TLBLatencyTracer
	DownTLBLatencyTracers              []TLBLatencyTracer
	L2PipelineLatencyTracers           []L2PipelineLatencyTracer
	TLBPipelineLatencyTracers          []TLBPipelineLatencyTracer
	MMUMemoryLatencyTracers            []MMUMemoryLatencyTracer
	PageWalkLatencyTracers             []PageWalkLatencyTracer
	DRAMLatencyTracers                 []DRAMLatencyTracer
	AddressTranslatorLatencyTracers    []AddressTranslatorLatencyTracer
	CacheHitRateTracers                []cacheHitRateTracer
	RTUCoalescingTracers               [][][][]*tracing.AverageCountTracer
	TLBHitRateTracers                  []TLBHitRateTracer
	PWCHitRateTracers                  []PWCHitRateTracer
	TranslationReqTracer               *tracing.TranslationReqTracer
	ConditionalRemoteTLBLatencyTracers []ConditionalRemoteTLBLatencyTracer
	RDMATransactionCounters            []rdmaTransactionCountTracer
	CDMATransactionCounters            []rdmaTransactionCountTracer
	PageTransactionCounters            []rdmaTransactionCountTracer
	RTUTransactionCounters             []rtuTransactionCountTracer
	CDMAAccessTracers                  []CDMAAccessTracer
	PageAccessTracers                  []CDMAAccessTracer
	RTUAccessTracers                   []RTUAccessTracer
	DRAMTransactionCounters            []dramTransactionCountTracer
	RemoteTLBLatencyTracers            []RemoteTLBLatencyTracer
	L2TLBBufLenTracers                 []*tracing.AverageCountTracer
	L2TLBBufLenG0Tracers               []*tracing.AverageCountTracer
	L2TLBCoalesceAddrTracers           []*tracing.AverageCountTracer
	L2TLBCoalesceTracers               []*tracing.AverageCountTracer
	L2TLBMSHRLenTracers                []*tracing.AverageCountTracer
	L2TLBMSHRUniqLenTracers            []*tracing.AverageCountTracer
	L2TLBMSHRLenG0Tracers              []*tracing.AverageCountTracer
	L2TLBMSHRUniqLenG0Tracers          []*tracing.AverageCountTracer
	ActivePageWalkerTracers            []ActivePageWalkerTracer
	L2TLBMySQLTracer                   tracing.MySQLTracer
	L2TLBQueueingImbalanceTracers      []TLBQueueImbalanceTracer
	PageWalkerImbalanceTracers         []PageWalkerImbalanceTracer
	EntropyTracers                     []EntropyTracer
	RemoteReferenceCountTracers        []RemoteReferenceCountTracer
	TLBSetMissTracers                  []TLBSetMissTracer
	TLBMSHRStallTracers                []TLBMSHRStallTracer
	TLBReqStallTracers                 []TLBReqStallTracer
	Benchmarks                         []benchmarks.Benchmark
	Timing                             bool
	Verify                             bool
	Parallel                           bool
	ReportInstCount                    bool
	ReportCacheLatency                 bool
	ReportTLBLatency                   bool
	ReportTLBConditionalStats          bool
	ReportMMUConditionalStats          bool
	ReportPageWalkLatency              bool
	ReportDRAMLatency                  bool
	ReportTranslationReqLatency        bool
	ReportAddressTranslatorLatency     bool
	ReportCacheHitRate                 bool
	ReportTLBCoalesce                  bool
	ReportRTUCoalesce                  bool
	ReportTLBHitRate                   bool
	ReportPWCHitRate                   bool
	ReportL2TLBMSHRLen                 bool
	ReportDRAMTransactionCount         bool
	ReportRDMATransactionCount         bool
	ReportCDMATransactionCount         bool
	ReportRTUTransactionCount          bool
	ReportActiveWalkerCount            bool
	L2TLBSQLTracing                    bool
	ReportL2TLBQueueImbalance          bool
	ReportPageWalkerImbalance          bool
	ReportEntropy                      bool
	ReportReferenceTracing             bool
	ReportTLBSetMissTracing            bool
	ReportTLBMSHRStallTracing          bool
	ReportTLBReqStalls                 bool
	UseUnifiedMemory                   bool
	UseLASPMemoryAlloc                 bool
	UseLASPHSLMemoryAlloc              bool
	UseCustomHSL                       bool
	metricsCollector                   *collector

	GPUIDs []int
}

// ParseFlag applies the runner flag to runner object
//nolint:gocyclo
func (r *Runner) ParseFlag() *Runner {
	if *parallelFlag {
		r.Parallel = true
	}

	if *verifyFlag {
		r.Verify = true
	}

	if *timingFlag {
		r.Timing = true
	}

	if *useUnifiedMemoryFlag {
		r.UseUnifiedMemory = true
	}

	if *useLASPMemoryAllocFlag {
		r.UseLASPMemoryAlloc = true
	}

	if *useLASPHSLMemoryAllocFlag {
		r.UseLASPHSLMemoryAlloc = true
	}

	if *useCustomHSL {
		r.UseCustomHSL = true
	}

	if *instCountReportFlag {
		r.ReportInstCount = true
	}

	if *cacheLatencyReportFlag {
		r.ReportCacheLatency = true
	}

	if *tlbLatencyReportFlag {
		r.ReportTLBLatency = true
	}

	if *tlbConditionalStatsReportFlag {
		r.ReportTLBConditionalStats = true
	}

	if *mmuConditionalStatsReportFlag {
		r.ReportMMUConditionalStats = true
	}

	if *pagewalkLatencyReportFlag {
		r.ReportPageWalkLatency = true
	}

	if *dramLatencyReportFlag {
		r.ReportDRAMLatency = true
	}

	if *translationReqLatencyFlag {
		r.ReportTranslationReqLatency = true
	}

	if *addressTranslatorLatencyFlag {
		r.ReportAddressTranslatorLatency = true
	}

	if *cacheHitRateReportFlag {
		r.ReportCacheHitRate = true
	}

	if *tlbHitRateReportFlag {
		r.ReportTLBHitRate = true
	}

	if *rtuCoalesceReportFlag {
		r.ReportRTUCoalesce = true
	}

	if *pwcHitRateReportFlag {
		r.ReportPWCHitRate = true
	}

	if *dramTransactionCountReportFlag {
		r.ReportDRAMTransactionCount = true
	}

	if *rdmaTransactionCountReportFlag {
		r.ReportRDMATransactionCount = true
	}

	if *cdmaTransactionCountReportFlag {
		r.ReportCDMATransactionCount = true
	}

	if *rtuTransactionCountReportFlag {
		r.ReportRTUTransactionCount = true
	}

	if *tlbCoalesceReportFlag {
		r.ReportTLBCoalesce = true
	}

	if *activeWalkerReportFlag {
		r.ReportActiveWalkerCount = true
	}

	if *l2tlbsqlTracingFlag {
		r.L2TLBSQLTracing = true
	}

	if *l2tlbQueueingImbalanceFlag {
		r.ReportL2TLBQueueImbalance = true
	}

	if *pageWalkerImbalanceFlag {
		r.ReportPageWalkerImbalance = true
	}

	if *entropyFlag {
		r.ReportEntropy = true
	}

	if *referenceTracingFlag {
		r.ReportReferenceTracing = true
	}

	if *reportAll {
		r.ReportInstCount = true
		r.ReportCacheLatency = true
		r.ReportTLBLatency = true
		r.ReportTLBConditionalStats = true
		r.ReportMMUConditionalStats = true
		r.ReportPageWalkLatency = true
		r.ReportDRAMLatency = true
		r.ReportCacheHitRate = true
		r.ReportTLBHitRate = true
		r.ReportTLBCoalesce = true
		r.ReportRTUCoalesce = true
		r.ReportPWCHitRate = true
		r.ReportActiveWalkerCount = true
		r.ReportDRAMTransactionCount = true
		r.ReportRDMATransactionCount = true
		r.ReportRTUTransactionCount = true
		r.ReportCDMATransactionCount = true
		//temp
		r.ReportTranslationReqLatency = true
		r.ReportAddressTranslatorLatency = true
		r.ReportL2TLBMSHRLen = true
		r.ReportL2TLBQueueImbalance = true
		r.ReportPageWalkerImbalance = true
		r.ReportEntropy = true
		r.ReportReferenceTracing = true
		// not given an option yet TODO
		r.ReportTLBSetMissTracing = true
		r.ReportTLBMSHRStallTracing = true
		r.ReportTLBReqStalls = true
	}

	// who decides what is essential?
	if *reportEssential {

		// r.ReportInstCount = true
		// r.ReportCacheLatency = true
		// r.ReportTLBLatency = true

		// r.ReportPageWalkLatency = true
		// r.ReportPWCHitRate = true

		// r.ReportCacheHitRate = true
		r.ReportTLBHitRate = true

		// r.ReportActiveWalkerCount = true
		// r.ReportDRAMTransactionCount = true
		// r.ReportRDMATransactionCount = true
		// r.ReportRTUTransactionCount = true
		// r.ReportCDMATransactionCount = true
		// r.ReportTLBCoalesce = true
		// r.ReportRTUCoalesce = true

		// r.ReportTranslationReqLatency = true
	}

	return r
}

func (r *Runner) startProfilingServer() {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	fmt.Println("Profiling server running on:",
		listener.Addr().(*net.TCPAddr).Port)

	panic(http.Serve(listener, nil))
}

// Init initializes the platform simulate
func (r *Runner) Init() *Runner {
	go r.startProfilingServer()

	r.ParseFlag()

	log.SetFlags(log.Llongfile | log.Ldate | log.Ltime)

	if r.Timing {
		r.buildTimingPlatform()
	} else {
		// r.buildEmuPlatform()
		r.buildTimingPlatform()
	}

	r.parseGPUFlag()

	r.metricsCollector = &collector{}
	r.addMaxInstStopper()
	r.addKernelTimeTracer()
	r.addInstCountTracer()
	r.addCacheLatencyTracer()
	r.addTLBLatencyTracer()
	r.addTLBConditionalTracer()
	r.addTLBCoalesceTracer()
	r.addL2TLBMSHRLenTracer()
	r.addRTUCoalesceTracer()
	r.addTLBHitRateTracer()
	r.addDRAMLatencyTracer()
	r.addDRAMTracer()
	r.addPageWalkLatencyTracer()
	r.addMMUConditionalTracer()
	r.addAddressTranslatorLatencyTracer()
	r.addCacheHitRateTracer()
	r.addPWCHitRateTracer()
	r.addRDMAEngineTracer()
	r.addCDMAEngineTracer()
	r.addRTUTracer()
	r.addActiveWalkerTracer()
	r.addL2TLBQueueingImbalanceTracer()
	r.addL2TLBEntropyTracer()
	r.addRemoteReferenceTracer()
	r.addTLBSetMissTracer()
	r.addTLBMSHRStallTracer()
	r.addTLBReqStallTracer()
	r.addPageWalkerImbalanceTracker()
	r.addL2TLBSQLTracer()
	if *platformType != "ideal" {
		r.addTranslationReqTracer()
	}

	atexit.Register(func() { r.reportStats() })

	return r
}

func (r *Runner) buildTimingPlatform() {
	// var b platform.Platform
	switch *platformType {
	case "ideal":
		b := platform.MakeIdealVMGPUBuilder()
		if r.Parallel {
			b = b.WithParallelEngine()
		}

		if *isaDebug {
			b = b.WithISADebugging()
		}

		if *visTracing {
			b = b.WithVisTracing()
		}

		if *memTracing {
			b = b.WithMemTracing()
		}

		if *tlbTracing {
			b = b.WithTLBTracing()
		}

		if *disableProgressBar {
			b = b.WithoutProgressBar()
		}
		b = b.WithAlg(*schedulingAlg)
		b = b.WithMemAllocatorType(*memAllocatorType)
		r.Engine, r.GPUDriver = b.Build()
	case "ideal2":
		b := platform.MakeOldIdealVM2GPUBuilder()
		if r.Parallel {
			b = b.WithParallelEngine()
		}

		if *isaDebug {
			b = b.WithISADebugging()
		}

		if *visTracing {
			b = b.WithVisTracing()
		}

		if *memTracing {
			b = b.WithMemTracing()
		}

		if *tlbTracing {
			b = b.WithTLBTracing()
		}

		if *disableProgressBar {
			b = b.WithoutProgressBar()
		}

		b = b.WithAlg(*schedulingAlg)
		b = b.WithMemAllocatorType(*memAllocatorType)
		r.Engine, r.GPUDriver = b.Build()
	case "i2":
		b := platform.MakeIdealVM2GPUPlatformBuilder()
		if r.Parallel {
			b.WithParallelEngine()
		}

		if *isaDebug {
			b.WithISADebugging()
		}

		if *visTracing {
			b.WithVisTracing()
		}

		if *memTracing {
			b.WithMemTracing()
		}

		if *tlbTracing {
			b.WithTLBTracing()
		}

		if *disableProgressBar {
			b.WithoutProgressBar()
		}

		b.WithAlg(*schedulingAlg)
		b.WithSchedulingPartition(*schedulingPartition)
		b.WithMemAllocatorType(*memAllocatorType)
		b.WithLog2PageSize(*log2PageSize)
		r.Engine, r.GPUDriver = b.Build()
	case "mcmgpu":
		b := platform.MakeMCMGPUBuilder()
		if r.Parallel {
			b = b.WithParallelEngine()
		}

		if *isaDebug {
			b = b.WithISADebugging()
		}

		if *visTracing {
			b = b.WithVisTracing()
		}

		if *memTracing {
			b = b.WithMemTracing()
		}

		if *tlbTracing {
			b = b.WithTLBTracing()
		}

		if *disableProgressBar {
			b = b.WithoutProgressBar()
		}
		b = b.WithAlg(*schedulingAlg)
		b = b.WithMemAllocatorType(*memAllocatorType)
		r.Engine, r.GPUDriver = b.Build()
	case "h1":
		b := platform.MakeH1PlatformBuilder()
		if r.Parallel {
			b.WithParallelEngine()
		}

		if *isaDebug {
			b.WithISADebugging()
		}

		if *visTracing {
			b.WithVisTracing()
		}

		if *memTracing {
			b.WithMemTracing()
		}

		if *tlbTracing {
			b.WithTLBTracing()
		}

		if *disableProgressBar {
			b.WithoutProgressBar()
		}
		b.WithAlg(*schedulingAlg)
		b.WithMemAllocatorType(*memAllocatorType)
		b = b.WithL2TLBStriping(*l2TlbStriping)
		r.Engine, r.GPUDriver = b.Build()
	case "h2":
		b := platform.MakeH2PlatformBuilder()
		if r.Parallel {
			b.WithParallelEngine()
		}

		if *isaDebug {
			b.WithISADebugging()
		}

		if *visTracing {
			b.WithVisTracing()
		}

		if *memTracing {
			b.WithMemTracing()
		}

		if *tlbTracing {
			b.WithTLBTracing()
		}

		if *disableProgressBar {
			b.WithoutProgressBar()
		}
		b.WithAlg(*schedulingAlg)
		b.WithMemAllocatorType(*memAllocatorType)
		b = b.WithL2TLBStriping(*l2TlbStriping)
		r.Engine, r.GPUDriver = b.Build()
	case "h3":
		b := platform.MakeH3PlatformBuilder()
		if r.Parallel {
			b.WithParallelEngine()
		}

		if *isaDebug {
			b.WithISADebugging()
		}

		if *visTracing {
			b.WithVisTracing()
		}

		if *memTracing {
			b.WithMemTracing()
		}

		if *tlbTracing {
			b.WithTLBTracing()
		}

		if *disableProgressBar {
			b.WithoutProgressBar()
		}
		b.WithAlg(*schedulingAlg)
		b.WithMemAllocatorType(*memAllocatorType)
		b.UseCoalescingTLBPort(*useCoalescingTLBPort)
		b = b.WithL2TLBStriping(*l2TlbStriping)
		r.Engine, r.GPUDriver = b.Build()
	case "h4":
		b := platform.MakeH4PlatformBuilder()
		if r.Parallel {
			b.WithParallelEngine()
		}

		if *isaDebug {
			b.WithISADebugging()
		}

		if *visTracing {
			b.WithVisTracing()
		}

		if *memTracing {
			b.WithMemTracing()
		}

		if *tlbTracing {
			b.WithTLBTracing()
		}

		if *disableProgressBar {
			b.WithoutProgressBar()
		}
		b.WithAlg(*schedulingAlg)
		b.WithMemAllocatorType(*memAllocatorType)
		b.UseCoalescingTLBPort(*useCoalescingTLBPort)
		b = b.WithL2TLBStriping(*l2TlbStriping)
		r.Engine, r.GPUDriver = b.Build()
	case "h5":
		b := platform.MakeH5PlatformBuilder()
		if r.Parallel {
			b.WithParallelEngine()
		}

		if *isaDebug {
			b.WithISADebugging()
		}

		if *visTracing {
			b.WithVisTracing()
		}

		if *memTracing {
			b.WithMemTracing()
		}

		if *tlbTracing {
			b.WithTLBTracing()
		}

		if *disableProgressBar {
			b.WithoutProgressBar()
		}
		b.WithAlg(*schedulingAlg)
		b.WithMemAllocatorType(*memAllocatorType)
		b.UseCoalescingTLBPort(*useCoalescingTLBPort)
		b = b.WithL2TLBStriping(*l2TlbStriping)
		r.Engine, r.GPUDriver = b.Build()
	case "h6":
		b := platform.MakeH6GPUPlatformBuilder()
		if r.Parallel {
			b.WithParallelEngine()
		}

		if *isaDebug {
			b.WithISADebugging()
		}

		if *visTracing {
			b.WithVisTracing()
		}

		if *memTracing {
			b.WithMemTracing()
		}

		if *tlbTracing {
			b.WithTLBTracing()
		}

		if *disableProgressBar {
			b.WithoutProgressBar()
		}
		b.WithAlg(*schedulingAlg)
		b.WithMemAllocatorType(*memAllocatorType)
		b.UseCoalescingTLBPort(*useCoalescingTLBPort)
		b = b.WithL2TLBStriping(*l2TlbStriping)
		r.Engine, r.GPUDriver = b.Build()
	case "h7":
		b := platform.MakeH7GPUPlatformBuilder()
		if r.Parallel {
			b.WithParallelEngine()
		}

		if *isaDebug {
			b.WithISADebugging()
		}

		if *visTracing {
			b.WithVisTracing()
		}

		if *memTracing {
			b.WithMemTracing()
		}

		if *tlbTracing {
			b.WithTLBTracing()
		}

		if *disableProgressBar {
			b.WithoutProgressBar()
		}
		b.WithAlg(*schedulingAlg)
		b.WithMemAllocatorType(*memAllocatorType)
		b.UseCoalescingTLBPort(*useCoalescingTLBPort)
		b = b.WithL2TLBStriping(*l2TlbStriping)
		r.Engine, r.GPUDriver = b.Build()
	case "privatetlb":
		b := platform.MakePrivateTLBPlatformBuilder()
		if r.Parallel {
			b.WithParallelEngine()
		}

		if *isaDebug {
			b.WithISADebugging()
		}

		if *visTracing {
			b.WithVisTracing()
		}

		if *memTracing {
			b.WithMemTracing()
		}

		if *tlbTracing {
			b.WithTLBTracing()
		}

		if *disableProgressBar {
			b.WithoutProgressBar()
		}
		b.WithAlg(*schedulingAlg)
		b.WithSchedulingPartition(*schedulingPartition)
		b.WithMemAllocatorType(*memAllocatorType)
		b.WithLog2PageSize(*log2PageSize)
		r.Engine, r.GPUDriver = b.Build()
	case "privateh2tlb":
		b := platform.MakePrivateH2TLBPlatformBuilder()
		if r.Parallel {
			b.WithParallelEngine()
		}

		if *isaDebug {
			b.WithISADebugging()
		}

		if *visTracing {
			b.WithVisTracing()
		}

		if *memTracing {
			b.WithMemTracing()
		}

		if *tlbTracing {
			b.WithTLBTracing()
		}

		if *disableProgressBar {
			b.WithoutProgressBar()
		}
		b.WithAlg(*schedulingAlg)
		b.WithSchedulingPartition(*schedulingPartition)
		b.WithMemAllocatorType(*memAllocatorType)
		r.Engine, r.GPUDriver = b.Build()
	case "privatetlb_ideal":
		b := platform.MakePrivateTLBIdealPlatformBuilder()
		if r.Parallel {
			b.WithParallelEngine()
		}

		if *isaDebug {
			b.WithISADebugging()
		}

		if *visTracing {
			b.WithVisTracing()
		}

		if *memTracing {
			b.WithMemTracing()
		}

		if *tlbTracing {
			b.WithTLBTracing()
		}

		if *disableProgressBar {
			b.WithoutProgressBar()
		}
		b.WithAlg(*schedulingAlg)
		b.WithMemAllocatorType(*memAllocatorType)
		r.Engine, r.GPUDriver = b.Build()
	case "xortlb":
		b := platform.MakeXORTLBGPUPlatformBuilder()
		if r.Parallel {
			b.WithParallelEngine()
		}

		if *isaDebug {
			b.WithISADebugging()
		}

		if *visTracing {
			b.WithVisTracing()
		}

		if *memTracing {
			b.WithMemTracing()
		}

		if *tlbTracing {
			b.WithTLBTracing()
		}

		if *disableProgressBar {
			b.WithoutProgressBar()
		}
		b.WithAlg(*schedulingAlg)
		b.WithSchedulingPartition(*schedulingPartition)
		b.WithMemAllocatorType(*memAllocatorType)
		b.UseCoalescingTLBPort(*useCoalescingTLBPort)
		b.UseCoalescingRTU(*useCoalescingRTU)
		b.WithLog2PageSize(*log2PageSize)
		b = b.WithL2TLBStriping(*l2TlbStriping)
		b = b.SwitchL2TLBStriping(*useSwitching)
		b = b.UsePtCaching(*ptCaching)
		r.Engine, r.GPUDriver = b.Build()
	case "xortlb_h1":
		b := platform.MakeXORTLBH1GPUPlatformBuilder()
		if r.Parallel {
			b.WithParallelEngine()
		}

		if *isaDebug {
			b.WithISADebugging()
		}

		if *visTracing {
			b.WithVisTracing()
		}

		if *memTracing {
			b.WithMemTracing()
		}

		if *tlbTracing {
			b.WithTLBTracing()
		}

		if *disableProgressBar {
			b.WithoutProgressBar()
		}
		b.WithAlg(*schedulingAlg)
		b.WithSchedulingPartition(*schedulingPartition)
		b.WithMemAllocatorType(*memAllocatorType)
		b.UseCoalescingTLBPort(*useCoalescingTLBPort)
		b.UseCoalescingRTU(*useCoalescingRTU)
		b = b.WithL2TLBStriping(*l2TlbStriping)
		r.Engine, r.GPUDriver = b.Build()
	case "xortlb_h2":
		b := platform.MakeXORTLBH2GPUPlatformBuilder()
		if r.Parallel {
			b.WithParallelEngine()
		}

		if *isaDebug {
			b.WithISADebugging()
		}

		if *visTracing {
			b.WithVisTracing()
		}

		if *memTracing {
			b.WithMemTracing()
		}

		if *tlbTracing {
			b.WithTLBTracing()
		}

		if *disableProgressBar {
			b.WithoutProgressBar()
		}
		b.WithAlg(*schedulingAlg)
		b.WithSchedulingPartition(*schedulingPartition)
		b.WithMemAllocatorType(*memAllocatorType)
		b.UseCoalescingTLBPort(*useCoalescingTLBPort)
		b.UseCoalescingRTU(*useCoalescingRTU)
		b = b.WithL2TLBStriping(*l2TlbStriping)
		r.Engine, r.GPUDriver = b.Build()
	case "distributedtlb":
		b := platform.MakeDistributedTLBGPUPlatformBuilder()
		if r.Parallel {
			b.WithParallelEngine()
		}

		if *isaDebug {
			b.WithISADebugging()
		}

		if *visTracing {
			b.WithVisTracing()
		}

		if *memTracing {
			b.WithMemTracing()
		}

		if *tlbTracing {
			b.WithTLBTracing()
		}

		if *disableProgressBar {
			b.WithoutProgressBar()
		}
		b.WithAlg(*schedulingAlg)
		b.WithMemAllocatorType(*memAllocatorType)
		b.UseCoalescingTLBPort(*useCoalescingTLBPort)
		b.UseCoalescingRTU(*useCoalescingRTU)
		b = b.WithL2TLBStriping(*l2TlbStriping)
		r.Engine, r.GPUDriver = b.Build()
	case "distlb":
		b := platform.MakeDisTLBGPUBuilder()
		if r.Parallel {
			b = b.WithParallelEngine()
		}

		if *isaDebug {
			b = b.WithISADebugging()
		}

		if *visTracing {
			b = b.WithVisTracing()
		}

		if *memTracing {
			b = b.WithMemTracing()
		}

		if *tlbTracing {
			b = b.WithTLBTracing()
		}

		if *disableProgressBar {
			b = b.WithoutProgressBar()
		}
		b = b.WithAlg(*schedulingAlg)
		b = b.UseCoalescingTLBPort(*useCoalescingTLBPort)
		b = b.WithMemAllocatorType(*memAllocatorType)
		b = b.WithL2TLBStriping(*l2TlbStriping)
		r.Engine, r.GPUDriver = b.Build()
	case "headroom1":
		b := platform.MakeHeadroom1Builder()
		if r.Parallel {
			b = b.WithParallelEngine()
		}

		if *isaDebug {
			b = b.WithISADebugging()
		}

		if *visTracing {
			b = b.WithVisTracing()
		}

		if *memTracing {
			b = b.WithMemTracing()
		}

		if *tlbTracing {
			b = b.WithTLBTracing()
		}

		if *disableProgressBar {
			b = b.WithoutProgressBar()
		}
		b = b.WithMemAllocatorType(*memAllocatorType)
		b = b.WithL2TLBStriping(*l2TlbStriping)
		r.Engine, r.GPUDriver = b.Build()
	case "headroom2":
		b := platform.MakeHeadroom2Builder()
		if r.Parallel {
			b = b.WithParallelEngine()
		}

		if *isaDebug {
			b = b.WithISADebugging()
		}

		if *visTracing {
			b = b.WithVisTracing()
		}

		if *memTracing {
			b = b.WithMemTracing()
		}

		if *tlbTracing {
			b = b.WithTLBTracing()
		}

		if *disableProgressBar {
			b = b.WithoutProgressBar()
		}
		b = b.WithMemAllocatorType(*memAllocatorType)
		b = b.WithL2TLBStriping(*l2TlbStriping)
		r.Engine, r.GPUDriver = b.Build()
	case "headroom3":
		b := platform.MakeHeadroom3Builder()
		if r.Parallel {
			b = b.WithParallelEngine()
		}

		if *isaDebug {
			b = b.WithISADebugging()
		}

		if *visTracing {
			b = b.WithVisTracing()
		}

		if *memTracing {
			b = b.WithMemTracing()
		}

		if *tlbTracing {
			b = b.WithTLBTracing()
		}

		if *disableProgressBar {
			b = b.WithoutProgressBar()
		}
		b = b.WithAlg(*schedulingAlg)
		b = b.WithMemAllocatorType(*memAllocatorType)
		b = b.WithL2TLBStriping(*l2TlbStriping)
		r.Engine, r.GPUDriver = b.Build()
	case "modtlb":
		b := platform.MakeMODTLBGPUPlatformBuilder()
		if r.Parallel {
			b.WithParallelEngine()
		}

		if *isaDebug {
			b.WithISADebugging()
		}

		if *visTracing {
			b.WithVisTracing()
		}

		if *memTracing {
			b.WithMemTracing()
		}

		if *tlbTracing {
			b.WithTLBTracing()
		}

		if *disableProgressBar {
			b.WithoutProgressBar()
		}
		b.WithAlg(*schedulingAlg)
		b.WithSchedulingPartition(*schedulingPartition)
		b.WithMemAllocatorType(*memAllocatorType)
		b.UseCoalescingTLBPort(*useCoalescingTLBPort)
		b.UseCoalescingRTU(*useCoalescingRTU)
		b.WithLog2PageSize(*log2PageSize)
		b = b.WithL2TLBStriping(*l2TlbStriping)
		b = b.SwitchL2TLBStriping(*useSwitching)
		b = b.UsePtCaching(*ptCaching)
		r.Engine, r.GPUDriver = b.Build()
	case "customtlb":
		b := platform.MakeCustomTLBGPUPlatformBuilder()
		if r.Parallel {
			b.WithParallelEngine()
		}

		if *isaDebug {
			b.WithISADebugging()
		}

		if *visTracing {
			b.WithVisTracing()
		}

		if *memTracing {
			b.WithMemTracing()
		}

		if *tlbTracing {
			b.WithTLBTracing()
		}

		if *disableProgressBar {
			b.WithoutProgressBar()
		}
		b.WithAlg(*schedulingAlg)
		b.WithSchedulingPartition(*schedulingPartition)
		b.WithMemAllocatorType(*memAllocatorType)
		b.WithCustomHSL(*customHSL)
		b.UseCoalescingTLBPort(*useCoalescingTLBPort)
		b.UseCoalescingRTU(*useCoalescingRTU)
		b.WithLog2PageSize(*log2PageSize)
		b = b.WithL2TLBStriping(*l2TlbStriping)
		b = b.SwitchL2TLBStriping(*useSwitching)
		b = b.UsePtCaching(*ptCaching)
		r.Engine, r.GPUDriver = b.Build()
	default:
		panic("oh no!")
	}
}

func (r *Runner) addMaxInstStopper() {
	if *maxInstCount == 0 {
		return
	}

	r.maxInstStopper = newInstStopper(*maxInstCount)
	for _, gpu := range r.GPUDriver.GPUs {
		for _, cu := range gpu.CUs {
			tracing.CollectTrace(cu.(tracing.NamedHookable), r.maxInstStopper)
		}
	}
}

func (r *Runner) addKernelTimeTracer() {
	r.KernelTimeCounter = tracing.NewBusyTimeTracer(
		func(task tracing.Task) bool {
			return task.What == "*driver.LaunchKernelCommand"
		})
	tracing.CollectTrace(r.GPUDriver, r.KernelTimeCounter)

	for _, gpu := range r.GPUDriver.GPUs {
		gpuKernelTimeCounter := tracing.NewBusyTimeTracer(
			func(task tracing.Task) bool {
				return task.What == "*protocol.LaunchKernelReq"
			})
		r.PerGPUKernelTimeCounter = append(
			r.PerGPUKernelTimeCounter, gpuKernelTimeCounter)
		tracing.CollectTrace(gpu.CommandProcessor, gpuKernelTimeCounter)
	}
}

func (r *Runner) addInstCountTracer() {
	if !r.ReportInstCount {
		return
	}

	for _, gpu := range r.GPUDriver.GPUs {
		for _, cu := range gpu.CUs {
			tracer := newInstTracer()
			r.InstCountTracers = append(r.InstCountTracers,
				instCountTracer{
					tracer: tracer,
					cu:     cu,
				})
			tracing.CollectTrace(cu.(tracing.NamedHookable), tracer)
		}
	}
}

func (r *Runner) addCacheLatencyTracer() {
	if !r.ReportCacheLatency {
		return
	}

	for _, gpu := range r.GPUDriver.GPUs {
		for _, cache := range gpu.L1ICaches {
			tracer := tracing.NewAverageTimeTracer(
				func(task tracing.Task) bool {
					return task.Kind == "req_in"
				})
			r.CacheLatencyTracers = append(r.CacheLatencyTracers,
				cacheLatencyTracer{tracer: tracer, cache: cache})
			tracing.CollectTrace(cache, tracer)
		}

		for _, cache := range gpu.L1SCaches {
			tracer := tracing.NewAverageTimeTracer(
				func(task tracing.Task) bool {
					return task.Kind == "req_in"
				})
			r.CacheLatencyTracers = append(r.CacheLatencyTracers,
				cacheLatencyTracer{tracer: tracer, cache: cache})
			tracing.CollectTrace(cache, tracer)
		}

		for _, cache := range gpu.L1VCaches {
			tracer := tracing.NewAverageTimeTracer(
				func(task tracing.Task) bool {
					return task.Kind == "req_in"
				})
			r.CacheLatencyTracers = append(r.CacheLatencyTracers,
				cacheLatencyTracer{tracer: tracer, cache: cache})
			tracing.CollectTrace(cache, tracer)
		}

		for _, cache := range gpu.L2Caches {
			tracer := tracing.NewAverageTimeTracer(
				func(task tracing.Task) bool {
					return task.Kind == "req_in"
				})
			r.CacheLatencyTracers = append(r.CacheLatencyTracers,
				cacheLatencyTracer{tracer: tracer, cache: cache})
			tracing.CollectTrace(cache, tracer)
		}

		for _, cache := range gpu.L2Caches {
			pipeline := cache.GetPipeline()
			tracer := tracing.NewAverageTimeTracer(
				func(task tracing.Task) bool {
					return task.Kind == "pipeline"
				})
			r.L2PipelineLatencyTracers = append(r.L2PipelineLatencyTracers,
				L2PipelineLatencyTracer{tracer: tracer, pipeline: pipeline})
			tracing.CollectTrace(pipeline, tracer)
		}

	}
}

func (r *Runner) addTLBLatencyTracer() {
	if !r.ReportTLBLatency {
		return
	}

	for _, gpu := range r.GPUDriver.GPUs {
		for _, tlb := range gpu.L1ITLBs {
			tracer := tracing.NewAverageTimeTracer(
				func(task tracing.Task) bool {
					return task.Kind == "req_in"
				})
			r.TLBLatencyTracers = append(r.TLBLatencyTracers,
				TLBLatencyTracer{tracer: tracer, tlb: tlb})
			tracing.CollectTrace(tlb, tracer)
		}

		for _, tlb := range gpu.L1STLBs {
			tracer := tracing.NewAverageTimeTracer(
				func(task tracing.Task) bool {
					return task.Kind == "req_in"
				})
			r.TLBLatencyTracers = append(r.TLBLatencyTracers,
				TLBLatencyTracer{tracer: tracer, tlb: tlb})
			tracing.CollectTrace(tlb, tracer)
		}

		for _, tlb := range gpu.L1VTLBs {
			tracer := tracing.NewAverageTimeTracer(
				func(task tracing.Task) bool {
					return task.Kind == "req_in"
				})
			r.TLBLatencyTracers = append(r.TLBLatencyTracers,
				TLBLatencyTracer{tracer: tracer, tlb: tlb})
			tracing.CollectTrace(tlb, tracer)
		}

		for _, tlb := range gpu.L2TLBs {
			tracer := tracing.NewAverageTimeTracer(
				func(task tracing.Task) bool {
					return task.Kind == "req_in"
				})
			r.TLBLatencyTracers = append(r.TLBLatencyTracers,
				TLBLatencyTracer{tracer: tracer, tlb: tlb})
			tracing.CollectTrace(tlb, tracer)
		}

		for _, tlb := range gpu.L2TLBs {
			pipeline := tlb.GetPipeline()
			tracer := tracing.NewAverageTimeTracer(
				func(task tracing.Task) bool {
					return task.Kind == "pipeline"
				})
			r.TLBPipelineLatencyTracers = append(r.TLBPipelineLatencyTracers,
				TLBPipelineLatencyTracer{tracer: tracer, pipeline: pipeline})
			tracing.CollectTrace(pipeline, tracer)
		}

		numL1VTLBs := len(gpu.L1VTLBs) // + len(gpu.L1STLBs) + len(gpu.L1ITLBs)
		allL1TLBs := make([]*tlb.TLB, numL1VTLBs)
		_ = copy(allL1TLBs, gpu.L1VTLBs)
		allL1TLBs = append(allL1TLBs, gpu.L1STLBs...)
		allL1TLBs = append(allL1TLBs, gpu.L1ITLBs...)
		for _, tlb := range allL1TLBs {
			tracer := tracing.NewAverageTimeTracer(
				func(task tracing.Task) bool {
					return task.Kind == "req_out"
				})
			r.DownTLBLatencyTracers = append(r.DownTLBLatencyTracers,
				TLBLatencyTracer{tracer: tracer, tlb: tlb})
			tracing.CollectTrace(tlb, tracer)
		}

		for _, tlb := range gpu.L2TLBs {
			tracer := tracing.NewAverageTimeTracer(
				func(task tracing.Task) bool {
					return task.Kind == "req_out"
				})
			r.DownTLBLatencyTracers = append(r.DownTLBLatencyTracers,
				TLBLatencyTracer{tracer: tracer, tlb: tlb})
			tracing.CollectTrace(tlb, tracer)
		}

	}
}

func (r *Runner) addTLBConditionalTracer() {
	if !r.ReportTLBConditionalStats {
		return
	}

	for _, gpu := range r.GPUDriver.GPUs {
		numL1VTLBs := len(gpu.L1VTLBs) // + len(gpu.L1STLBs) + len(gpu.L1ITLBs)
		allL1TLBs := make([]*tlb.TLB, numL1VTLBs)
		_ = copy(allL1TLBs, gpu.L1VTLBs)
		allL1TLBs = append(allL1TLBs, gpu.L1STLBs...)
		allL1TLBs = append(allL1TLBs, gpu.L1ITLBs...)
		for _, tlb := range allL1TLBs {
			tracer := tracing.NewConditionalAverageTimeTracer(
				func(task tracing.Task) bool {
					return task.Kind == "L2TLB_stats" // || task.Kind == "local_tlb"
				})
			r.ConditionalRemoteTLBLatencyTracers = append(r.ConditionalRemoteTLBLatencyTracers,
				ConditionalRemoteTLBLatencyTracer{tracer: tracer, tlb: tlb})
			tracing.CollectTrace(tlb, tracer)
		}
	}
}

func (r *Runner) addTLBCoalesceTracer() {
	if !r.ReportTLBCoalesce {
		return
	}

	for _, gpu := range r.GPUDriver.GPUs {
		for _, tlb := range gpu.L2TLBs {
			tracer := tracing.NewAverageCountTracer(
				func(task tracing.Task) bool {
					return task.Kind == "buflen"
				})
			r.L2TLBBufLenTracers = append(r.L2TLBBufLenTracers, tracer)
			tracing.CollectTrace(tlb, tracer)
		}
		for _, tlb := range gpu.L2TLBs {
			tracer := tracing.NewAverageCountTracer(
				func(task tracing.Task) bool {
					return task.Kind == "bufleng0"
				})
			r.L2TLBBufLenG0Tracers = append(r.L2TLBBufLenG0Tracers, tracer)
			tracing.CollectTrace(tlb, tracer)
		}
		for _, tlb := range gpu.L2TLBs {
			tracer := tracing.NewAverageCountTracer(
				func(task tracing.Task) bool {
					return task.Kind == "coalesceaddr"
				})
			r.L2TLBCoalesceAddrTracers = append(r.L2TLBCoalesceAddrTracers, tracer)
			tracing.CollectTrace(tlb, tracer)
		}
		for _, tlb := range gpu.L2TLBs {
			tracer := tracing.NewAverageCountTracer(
				func(task tracing.Task) bool {
					return task.Kind == "coalesce"
				})
			r.L2TLBCoalesceTracers = append(r.L2TLBCoalesceTracers, tracer)
			tracing.CollectTrace(tlb, tracer)
		}
	}
}

func (r *Runner) addL2TLBMSHRLenTracer() {
	if !r.ReportL2TLBMSHRLen {
		return
	}

	for _, gpu := range r.GPUDriver.GPUs {
		for _, tlb := range gpu.L2TLBs {
			tracer := tracing.NewAverageCountTracer(
				func(task tracing.Task) bool {
					return task.Kind == "MSHRlen"
				})
			r.L2TLBMSHRLenTracers = append(r.L2TLBMSHRLenTracers, tracer)
			tracing.CollectTrace(tlb, tracer)
		}
		for _, tlb := range gpu.L2TLBs {
			tracer := tracing.NewAverageCountTracer(
				func(task tracing.Task) bool {
					return task.Kind == "MSHRlen_g0"
				})
			r.L2TLBMSHRLenG0Tracers = append(r.L2TLBMSHRLenG0Tracers, tracer)
			tracing.CollectTrace(tlb, tracer)
		}
		for _, tlb := range gpu.L2TLBs {
			tracer := tracing.NewAverageCountTracer(
				func(task tracing.Task) bool {
					return task.Kind == "MSHRuniq"
				})
			r.L2TLBMSHRUniqLenTracers = append(r.L2TLBMSHRUniqLenTracers, tracer)
			tracing.CollectTrace(tlb, tracer)
		}
		for _, tlb := range gpu.L2TLBs {
			tracer := tracing.NewAverageCountTracer(
				func(task tracing.Task) bool {
					return task.Kind == "MSHRuniq_g0"
				})
			r.L2TLBMSHRUniqLenG0Tracers = append(r.L2TLBMSHRUniqLenG0Tracers, tracer)
			tracing.CollectTrace(tlb, tracer)
		}
	}
}

func getCoalescingTracer(prefix string, statName string) *tracing.AverageCountTracer {
	return tracing.NewAverageCountTracer(func(task tracing.Task) bool {
		return strings.Contains(task.Kind, prefix+"-"+statName)
	})
}

func (r *Runner) addRTUCoalesceTracer() {
	if !r.ReportRTUCoalesce {
		return
	}
	r.RTUCoalescingTracers = make([][][][]*tracing.AverageCountTracer, 2)
	for i, prefix := range [2]string{"incoming", "outgoing"} {
		r.RTUCoalescingTracers[i] = make([][][]*tracing.AverageCountTracer, 4)
		for j, statName := range [4]string{"buflen", "buflen-g0", "coalesceable-addr", "replication-count"} {
			numGPUs := len(r.GPUDriver.GPUs)
			r.RTUCoalescingTracers[i][j] = make([][]*tracing.AverageCountTracer, numGPUs)
			for k, gpu := range r.GPUDriver.GPUs {
				numRTUs := len(gpu.RemoteAddressTranslationUnits)
				r.RTUCoalescingTracers[i][j][k] = make([]*tracing.AverageCountTracer, numRTUs)
				for l, rtu := range gpu.RemoteAddressTranslationUnits {
					tracer := getCoalescingTracer(prefix, statName)
					r.RTUCoalescingTracers[i][j][k][l] = tracer
					tracing.CollectTrace(rtu, tracer)
				}
			}
		}
	}
	return
}

func (r *Runner) addAddressTranslatorLatencyTracer() {
	if !r.ReportAddressTranslatorLatency {
		return
	}
	for _, gpu := range r.GPUDriver.GPUs {
		for _, l1vAddrTrans := range gpu.L1VAddrTranslator {
			tracer := tracing.NewConditionalAverageTimeTracer(
				func(task tracing.Task) bool {
					return task.Kind == "addr-translator-stats" // || task.Kind == "local_tlb"
				})
			tracing.CollectTrace(l1vAddrTrans, tracer)
			r.AddressTranslatorLatencyTracers = append(r.AddressTranslatorLatencyTracers,
				AddressTranslatorLatencyTracer{tracer: tracer, addressTranslator: l1vAddrTrans})
		}
	}
	for _, gpu := range r.GPUDriver.GPUs {
		for _, l1sAddrTrans := range gpu.L1SAddrTranslator {
			tracer := tracing.NewConditionalAverageTimeTracer(
				func(task tracing.Task) bool {
					return task.Kind == "addr-translator-stats" // || task.Kind == "local_tlb"
				})

			tracing.CollectTrace(l1sAddrTrans, tracer)
			r.AddressTranslatorLatencyTracers = append(r.AddressTranslatorLatencyTracers,
				AddressTranslatorLatencyTracer{tracer: tracer, addressTranslator: l1sAddrTrans})
		}
	}
	for _, gpu := range r.GPUDriver.GPUs {
		for _, l1iAddrTrans := range gpu.L1IAddrTranslator {
			tracer := tracing.NewConditionalAverageTimeTracer(
				func(task tracing.Task) bool {
					return task.Kind == "addr-translator-stats" // || task.Kind == "local_tlb"
				})
			tracing.CollectTrace(l1iAddrTrans, tracer)
			r.AddressTranslatorLatencyTracers = append(r.AddressTranslatorLatencyTracers,
				AddressTranslatorLatencyTracer{tracer: tracer, addressTranslator: l1iAddrTrans})
		}
	}
}

func (r *Runner) addPageWalkLatencyTracer() {
	if !r.ReportPageWalkLatency {
		return
	}
	for _, gpu := range r.GPUDriver.GPUs {
		for _, mmu := range gpu.MMUs {
			tracer := tracing.NewAverageTimeTracer(
				func(task tracing.Task) bool {
					return task.Kind == "req_in"
				})
			r.PageWalkLatencyTracers = append(r.PageWalkLatencyTracers,
				PageWalkLatencyTracer{tracer: tracer, mmu: mmu})
			tracing.CollectTrace(mmu, tracer)
		}
	}
}

func (r *Runner) addMMUConditionalTracer() {
	if !r.ReportMMUConditionalStats {
		return
	}
	for _, gpu := range r.GPUDriver.GPUs {
		for _, mmu := range gpu.MMUs {
			tracer := tracing.NewConditionalAverageTimeTracer(
				func(task tracing.Task) bool {
					return task.Kind == "MMU_mem_latency"
				})
			r.MMUMemoryLatencyTracers = append(r.MMUMemoryLatencyTracers,
				MMUMemoryLatencyTracer{tracer: tracer, mmu: mmu})
			tracing.CollectTrace(mmu, tracer)
		}
	}
}

func (r *Runner) addActiveWalkerTracer() {
	if !r.ReportActiveWalkerCount {
		return
	}

	for _, gpu := range r.GPUDriver.GPUs {
		for _, mmu := range gpu.MMUs {
			tracer := tracing.NewAverageCountTracer(
				func(task tracing.Task) bool {
					return task.Kind == "num_active_walkers"
				})
			r.ActivePageWalkerTracers = append(r.ActivePageWalkerTracers,
				ActivePageWalkerTracer{tracer: tracer, mmu: mmu})
			tracing.CollectTrace(mmu, tracer)
		}
	}

}

func (r *Runner) addDRAMLatencyTracer() {
	if !r.ReportDRAMLatency {
		return
	}

	for _, gpu := range r.GPUDriver.GPUs {
		for _, dram := range gpu.MemoryControllers {
			tracer := tracing.NewAverageTimeTracer(
				func(task tracing.Task) bool {
					return task.Kind == "req_in"
				})
			r.DRAMLatencyTracers = append(r.DRAMLatencyTracers,
				DRAMLatencyTracer{tracer: tracer, dram: dram})
			tracing.CollectTrace(dram, tracer)
		}

	}
}

func (r *Runner) addPWCHitRateTracer() {
	if !r.ReportPWCHitRate {
		return
	}

	for _, gpu := range r.GPUDriver.GPUs {
		for _, mmu := range gpu.MMUs {
			tracer := tracing.NewStepCountTracer(
				func(task tracing.Task) bool { return task.ID != "" })
			r.PWCHitRateTracers = append(r.PWCHitRateTracers,
				PWCHitRateTracer{tracer: tracer, mmu: mmu})
			tracing.CollectTrace(mmu, tracer)
		}
	}
}

func (r *Runner) addCacheHitRateTracer() {
	if !r.ReportCacheHitRate {
		return
	}

	for _, gpu := range r.GPUDriver.GPUs {
		for _, cache := range gpu.L1VCaches {
			tracer := tracing.NewStepCountTracer(
				//	func(task tracing.Task) bool { return true })
				func(task tracing.Task) bool {
					return task.Kind == "l1_transaction"
				})
			r.CacheHitRateTracers = append(r.CacheHitRateTracers,
				cacheHitRateTracer{tracer: tracer, cache: cache})
			tracing.CollectTrace(cache, tracer)
		}

		for _, cache := range gpu.L1SCaches {
			tracer := tracing.NewStepCountTracer(
				//func(task tracing.Task) bool { return true })
				func(task tracing.Task) bool {
					return task.Kind == "l1_transaction"
				})
			r.CacheHitRateTracers = append(r.CacheHitRateTracers,
				cacheHitRateTracer{tracer: tracer, cache: cache})
			tracing.CollectTrace(cache, tracer)
		}

		for _, cache := range gpu.L1ICaches {
			tracer := tracing.NewStepCountTracer(
				//func(task tracing.Task) bool { return true })
				func(task tracing.Task) bool {
					return task.Kind == "l1_transaction"
				})
			r.CacheHitRateTracers = append(r.CacheHitRateTracers,
				cacheHitRateTracer{tracer: tracer, cache: cache})
			tracing.CollectTrace(cache, tracer)
		}

		for _, cache := range gpu.L2Caches {
			tracer := tracing.NewStepCountTracer(
				func(task tracing.Task) bool {
					return task.Kind == "req_in"
				})
			r.CacheHitRateTracers = append(r.CacheHitRateTracers,
				cacheHitRateTracer{tracer: tracer, cache: cache})
			tracing.CollectTrace(cache, tracer)
		}
	}
}

func (r *Runner) addTLBHitRateTracer() {
	if !r.ReportTLBHitRate {
		return
	}

	for _, gpu := range r.GPUDriver.GPUs {
		for _, tlb := range gpu.L1VTLBs {
			tracer := tracing.NewStepCountTracer(
				func(task tracing.Task) bool { /*return true*/
					return task.Kind == "req_in"
				})
			r.TLBHitRateTracers = append(r.TLBHitRateTracers,
				TLBHitRateTracer{tracer: tracer, tlb: tlb})
			tracing.CollectTrace(tlb, tracer)
		}

		for _, tlb := range gpu.L1STLBs {
			tracer := tracing.NewStepCountTracer(
				func(task tracing.Task) bool { /*return true })*/
					return task.Kind == "req_in"
				})
			r.TLBHitRateTracers = append(r.TLBHitRateTracers,
				TLBHitRateTracer{tracer: tracer, tlb: tlb})
			tracing.CollectTrace(tlb, tracer)
		}

		for _, tlb := range gpu.L1ITLBs {
			tracer := tracing.NewStepCountTracer(
				func(task tracing.Task) bool { /*return true })*/
					return task.Kind == "req_in"
				})
			r.TLBHitRateTracers = append(r.TLBHitRateTracers,
				TLBHitRateTracer{tracer: tracer, tlb: tlb})
			tracing.CollectTrace(tlb, tracer)
		}

		for _, tlb := range gpu.L2TLBs {
			tracer := tracing.NewStepCountTracer(
				func(task tracing.Task) bool { return true })
			r.TLBHitRateTracers = append(r.TLBHitRateTracers,
				TLBHitRateTracer{tracer: tracer, tlb: tlb})
			tracing.CollectTrace(tlb, tracer)
		}

	}
}

func (r *Runner) addRDMAEngineTracer() {
	if !r.ReportRDMATransactionCount {
		return
	}

	for _, gpu := range r.GPUDriver.GPUs {
		t := rdmaTransactionCountTracer{}
		t.rdmaEngine = gpu.RDMAEngine
		t.incomingTracer = tracing.NewAverageTimeTracer(
			func(task tracing.Task) bool {
				if task.Kind != "req_in" {
					return false
				}

				isFromOutside := strings.Contains(
					task.Detail.(akita.Msg).Meta().Src.Name(), "RDMA")
				if !isFromOutside {
					return false
				}

				return true
			})
		t.outgoingTracer = tracing.NewAverageTimeTracer(
			func(task tracing.Task) bool {
				if task.Kind != "req_in" {
					return false
				}

				isFromOutside := strings.Contains(
					task.Detail.(akita.Msg).Meta().Src.Name(), "RDMA")
				if isFromOutside {
					return false
				}

				return true
			})

		tracing.CollectTrace(t.rdmaEngine, t.incomingTracer)
		tracing.CollectTrace(t.rdmaEngine, t.outgoingTracer)

		r.RDMATransactionCounters = append(r.RDMATransactionCounters, t)
	}
}

func (r *Runner) addRTUTracer() {
	if !r.ReportRTUTransactionCount {
		return
	}

	for _, gpu := range r.GPUDriver.GPUs {
		for _, rtuUnit := range gpu.RemoteAddressTranslationUnits {
			t := rtuTransactionCountTracer{}
			t.rtuUnit = rtuUnit
			t.incomingTracer = tracing.NewAverageTimeTracer(
				func(task tracing.Task) bool {
					if task.Kind != "req_in" {
						return false
					}

					isFromOutside := strings.Contains(
						task.Detail.(akita.Msg).Meta().Src.Name(), "RTU")
					if !isFromOutside {
						return false
					}

					return true
				})
			t.outgoingTracer = tracing.NewAverageTimeTracer(
				func(task tracing.Task) bool {
					if task.Kind != "req_in" {
						return false
					}

					isFromOutside := strings.Contains(
						task.Detail.(akita.Msg).Meta().Src.Name(), "RTU")
					if isFromOutside {
						return false
					}

					return true
				})

			tracing.CollectTrace(t.rtuUnit, t.incomingTracer)
			tracing.CollectTrace(t.rtuUnit, t.outgoingTracer)

			r.RTUTransactionCounters = append(r.RTUTransactionCounters, t)

			rtuAccessTracer := tracing.NewStepCountTracer(
				func(task tracing.Task) bool {
					/*return true*/
					return task.Kind == "req_in"
				})
			r.RTUAccessTracers = append(r.RTUAccessTracers,
				RTUAccessTracer{tracer: rtuAccessTracer, rtu: rtuUnit})
			tracing.CollectTrace(rtuUnit, rtuAccessTracer)

		}
	}
}

func (r *Runner) addCDMAEngineTracer() {
	if !r.ReportCDMATransactionCount {
		return
	}

	for _, gpu := range r.GPUDriver.GPUs {
		for _, cdmaEngine := range gpu.ChipRDMAEngines {
			t := rdmaTransactionCountTracer{}
			t.rdmaEngine = cdmaEngine
			t.incomingTracer = tracing.NewAverageTimeTracer(
				func(task tracing.Task) bool {
					if task.Kind != "req_in" {
						return false
					}

					isFromOutside := strings.Contains(
						task.Detail.(akita.Msg).Meta().Src.Name(), "ChipRDMA")
					if !isFromOutside {
						return false
					}

					return true
				})
			t.outgoingTracer = tracing.NewAverageTimeTracer(
				func(task tracing.Task) bool {
					if task.Kind != "req_in" {
						return false
					}

					isFromOutside := strings.Contains(
						task.Detail.(akita.Msg).Meta().Src.Name(), "ChipRDMA")
					if isFromOutside {
						return false
					}

					return true
				})

			tracing.CollectTrace(t.rdmaEngine, t.incomingTracer)
			tracing.CollectTrace(t.rdmaEngine, t.outgoingTracer)

			r.CDMATransactionCounters = append(r.CDMATransactionCounters, t)

			tracer := tracing.NewStepCountTracer(
				func(task tracing.Task) bool { /*return true*/
					return task.Kind == "req_in"
				})
			r.CDMAAccessTracers = append(r.CDMAAccessTracers,
				CDMAAccessTracer{tracer: tracer, cdmaEngine: cdmaEngine})
			tracing.CollectTrace(cdmaEngine, tracer)

		}
	}

	for _, gpu := range r.GPUDriver.GPUs {
		for _, cdmaEngine := range gpu.PageRDMAEngines {
			t := rdmaTransactionCountTracer{}
			t.rdmaEngine = cdmaEngine
			t.incomingTracer = tracing.NewAverageTimeTracer(
				func(task tracing.Task) bool {
					if task.Kind != "req_in" {
						return false
					}

					isFromOutside := strings.Contains(
						task.Detail.(akita.Msg).Meta().Src.Name(), "PageRDMA")
					if !isFromOutside {
						return false
					}

					return true
				})
			t.outgoingTracer = tracing.NewAverageTimeTracer(
				func(task tracing.Task) bool {
					if task.Kind != "req_in" {
						return false
					}

					isFromOutside := strings.Contains(
						task.Detail.(akita.Msg).Meta().Src.Name(), "PageRDMA")
					if isFromOutside {
						return false
					}

					return true
				})

			tracing.CollectTrace(t.rdmaEngine, t.incomingTracer)
			tracing.CollectTrace(t.rdmaEngine, t.outgoingTracer)

			r.PageTransactionCounters = append(r.PageTransactionCounters, t)

			tracer := tracing.NewStepCountTracer(
				func(task tracing.Task) bool { /*return true*/
					return task.Kind == "req_in"
				})
			r.PageAccessTracers = append(r.PageAccessTracers,
				CDMAAccessTracer{tracer: tracer, cdmaEngine: cdmaEngine})
			tracing.CollectTrace(cdmaEngine, tracer)

		}
	}
}

func (r *Runner) addDRAMTracer() {
	if !r.ReportDRAMTransactionCount {
		return
	}

	for _, gpu := range r.GPUDriver.GPUs {
		for _, dram := range gpu.MemoryControllers {
			t := dramTransactionCountTracer{}
			// t.dram = dram.(*idealmemcontroller.Comp)
			// t.dram = dram.(*ram.MemController)
			t.dram = dram
			t.tracer = tracing.NewAverageTimeTracer(
				func(task tracing.Task) bool {
					return true
				})

			tracing.CollectTrace(t.dram, t.tracer)

			r.DRAMTransactionCounters = append(r.DRAMTransactionCounters, t)
		}
	}
}

func (r *Runner) addL2TLBSQLTracer() {
	if !r.L2TLBSQLTracing {
		return
	}
	tracer := tracing.NewMySQLTracer()
	tracer.Init()
	for _, gpu := range r.GPUDriver.GPUs {
		for _, l2tlb := range gpu.L2TLBs {
			tracing.CollectTrace(l2tlb, tracer)
		}
	}
	for _, gpu := range r.GPUDriver.GPUs {
		for _, l1tlb := range gpu.L1VTLBs {
			tracing.CollectTrace(l1tlb, tracer)
		}
	}
	for _, gpu := range r.GPUDriver.GPUs {
		for _, l1tlb := range gpu.L1ITLBs {
			tracing.CollectTrace(l1tlb, tracer)
		}
	}
	for _, gpu := range r.GPUDriver.GPUs {
		for _, l1tlb := range gpu.L1STLBs {
			tracing.CollectTrace(l1tlb, tracer)
		}
	}
}

func (r *Runner) addTranslationReqTracer() {
	if !r.ReportTranslationReqLatency {
		return
	}
	tracer := tracing.NewTranslationReqTracer(
		func(task tracing.Task) bool {
			if task.Kind == "trace-trans-req" {
				return true
			}
			return false
		})
	r.TranslationReqTracer = tracer
	for _, gpu := range r.GPUDriver.GPUs {
		for _, l2tlb := range gpu.L2TLBs {
			tracing.CollectTrace(l2tlb, tracer)
		}
	}
	for _, gpu := range r.GPUDriver.GPUs {
		for _, l1tlb := range gpu.L1VTLBs {
			tracing.CollectTrace(l1tlb, tracer)
		}
	}
	for _, gpu := range r.GPUDriver.GPUs {
		for _, l1tlb := range gpu.L1ITLBs {
			tracing.CollectTrace(l1tlb, tracer)
		}
	}
	for _, gpu := range r.GPUDriver.GPUs {
		for _, l1tlb := range gpu.L1STLBs {
			tracing.CollectTrace(l1tlb, tracer)
		}
	}

	for _, gpu := range r.GPUDriver.GPUs {
		for _, l1vAddrTrans := range gpu.L1VAddrTranslator {
			tracing.CollectTrace(l1vAddrTrans, tracer)
		}
	}
	for _, gpu := range r.GPUDriver.GPUs {
		for _, l1sAddrTrans := range gpu.L1SAddrTranslator {
			tracing.CollectTrace(l1sAddrTrans, tracer)
		}
	}
	for _, gpu := range r.GPUDriver.GPUs {
		for _, l1iAddrTrans := range gpu.L1IAddrTranslator {
			tracing.CollectTrace(l1iAddrTrans, tracer)
		}
	}
	for _, gpu := range r.GPUDriver.GPUs {
		for _, rtu := range gpu.RemoteAddressTranslationUnits {
			tracing.CollectTrace(rtu, tracer)
		}
	}
}

func (r *Runner) addL2TLBQueueingImbalanceTracer() {
	if !r.ReportL2TLBQueueImbalance {
		return
	}
	for _, gpu := range r.GPUDriver.GPUs {
		for _, this := range gpu.L2TLBs {
			tracer := tlb.NewGlobalTLBQueueingTracer(
				func(task tracing.Task) bool {
					if task.Kind == "imbalance" {
						return true
					}
					return false
				})
			tracer.AddThis(this)
			for _, tlb := range gpu.L2TLBs {
				tracer.AddTLB(tlb)
			}
			tracing.CollectTrace(this, tracer)
			r.L2TLBQueueingImbalanceTracers = append(r.L2TLBQueueingImbalanceTracers,
				TLBQueueImbalanceTracer{tracer: tracer, tlb: this})
		}
	}
	return
}

func (r *Runner) addL2TLBEntropyTracer() {
	if !r.ReportEntropy {
		return
	}
	for _, gpu := range r.GPUDriver.GPUs {
		for _, this := range gpu.L2TLBs {
			tracer := tlb.NewEntropyTracer(
				func(task tracing.Task) bool {
					if task.Kind == "entropy" {
						return true
					}
					return false
				})
			tracing.CollectTrace(this, tracer)
			r.EntropyTracers = append(r.EntropyTracers,
				EntropyTracer{tracer: tracer, tlb: this})
		}
	}
	return
}

func (r *Runner) addRemoteReferenceTracer() {
	if !r.ReportReferenceTracing {
		return
	}
	for _, gpu := range r.GPUDriver.GPUs {
		for _, rtuUnit := range gpu.RemoteAddressTranslationUnits {
			tracer := tracing.NewReferenceTracer(
				func(task tracing.Task) bool {
					if task.Kind == "reference_counting" {
						return true
					}
					return false
				})
			tracing.CollectTrace(rtuUnit, tracer)
			r.RemoteReferenceCountTracers = append(r.RemoteReferenceCountTracers,
				RemoteReferenceCountTracer{tracer: tracer, rtuUnit: rtuUnit})
		}
	}
	return
}

func (r *Runner) addTLBSetMissTracer() {
	if !r.ReportTLBSetMissTracing {
		return
	}
	for _, gpu := range r.GPUDriver.GPUs {
		for _, tlb := range gpu.L2TLBs {
			tracer := tracing.NewReferenceTracer(
				func(task tracing.Task) bool {
					if task.Kind == "set_miss_tracing" {
						return true
					}
					return false
				})
			tracing.CollectTrace(tlb, tracer)
			r.TLBSetMissTracers = append(r.TLBSetMissTracers,
				TLBSetMissTracer{tracer: tracer, tlb: tlb})
		}
	}
	return
}

func (r *Runner) addTLBMSHRStallTracer() {
	if !r.ReportTLBMSHRStallTracing {
		return
	}
	for _, gpu := range r.GPUDriver.GPUs {
		for _, tlb := range gpu.L2TLBs {
			tracer := tracing.NewTotalTimeTracer(
				func(task tracing.Task) bool {
					if task.Kind == "mshr_stall" {
						return true
					}
					return false
				})
			tracing.CollectTrace(tlb, tracer)
			r.TLBMSHRStallTracers = append(r.TLBMSHRStallTracers,
				TLBMSHRStallTracer{tracer: tracer, tlb: tlb})
		}
	}
	return
}

func (r *Runner) addTLBReqStallTracer() {
	if !r.ReportTLBReqStalls {
		return
	}
	for _, gpu := range r.GPUDriver.GPUs {
		for _, tlb := range gpu.L2TLBs {
			tracer := tracing.NewStepCountTracer(
				func(task tracing.Task) bool {
					return task.Kind == "stalled-l2-tlb-req-count" || task.Kind == "l2-tlb-req-count"
				})
			tracing.CollectTrace(tlb, tracer)
			r.TLBReqStallTracers = append(r.TLBReqStallTracers,
				TLBReqStallTracer{tracer: tracer, tlb: tlb})
		}
	}
}

func (r *Runner) addPageWalkerImbalanceTracker() {
	if !r.ReportPageWalkerImbalance {
		return
	}
	for _, gpu := range r.GPUDriver.GPUs {
		for _, this := range gpu.MMUs {
			tracer := mmu.NewGlobalPageWalkerOccupancyTracer(
				func(task tracing.Task) bool {
					if task.Kind == "imbalance" {
						return true
					}
					return false
				})
			tracer.AddThis(this)
			for _, mmu := range gpu.MMUs {
				tracer.AddMMU(mmu)
			}
			tracing.CollectTrace(this, tracer)
			r.PageWalkerImbalanceTracers = append(r.PageWalkerImbalanceTracers,
				PageWalkerImbalanceTracer{tracer: tracer, mmu: this})
		}
	}
	return
}

func (r *Runner) parseGPUFlag() {
	if *gpuFlag == "" && *unifiedGPUFlag == "" {
		r.GPUIDs = []int{1}
		return
	}

	if *gpuFlag != "" && *unifiedGPUFlag != "" {
		panic("cannot use -gpus and -unified-gpus together")
	}

	if *unifiedGPUFlag != "" {
		gpuIDs := r.gpuIDStringToList(*unifiedGPUFlag)
		unifiedGPUID := r.GPUDriver.CreateUnifiedGPU(nil, gpuIDs)
		r.GPUIDs = []int{unifiedGPUID}
		return
	}

	gpuIDs := r.gpuIDStringToList(*gpuFlag)
	r.GPUIDs = gpuIDs
}

func (r *Runner) gpuIDStringToList(gpuIDsString string) []int {
	gpuIDs := make([]int, 0)
	gpuIDTokens := strings.Split(gpuIDsString, ",")
	for _, t := range gpuIDTokens {
		gpuID, err := strconv.Atoi(t)
		if err != nil {
			panic(err)
		}
		gpuIDs = append(gpuIDs, gpuID)
	}
	return gpuIDs
}

// AddBenchmark adds an benchmark that the driver runs
func (r *Runner) AddBenchmark(b benchmarks.Benchmark) {
	b.SelectGPU(r.GPUIDs)
	if r.UseUnifiedMemory {
		b.SetUnifiedMemory()
	}
	if r.UseLASPMemoryAlloc {
		b.SetLASPMemoryAlloc()
	}
	if r.UseLASPHSLMemoryAlloc {
		b.SetLASPHSLMemoryAlloc()
	}
	// if r.UseCustomHSL {
	// 	b.SetCustomHSL()
	// }
	r.Benchmarks = append(r.Benchmarks, b)
}

// AddBenchmarkWithoutSettingGPUsToUse allows for user specified GPUs for
// the benchmark to run.
func (r *Runner) AddBenchmarkWithoutSettingGPUsToUse(b benchmarks.Benchmark) {
	if r.UseUnifiedMemory {
		b.SetUnifiedMemory()
	}
	if r.UseLASPMemoryAlloc {
		b.SetLASPMemoryAlloc()
	}
	r.Benchmarks = append(r.Benchmarks, b)
}

// Run runs the benchmark on the simulator
func (r *Runner) Run() {
	r.GPUDriver.Run()

	var wg sync.WaitGroup
	for _, b := range r.Benchmarks {
		wg.Add(1)
		go func(b benchmarks.Benchmark, wg *sync.WaitGroup) {
			if r.Verify {
				if b, ok := b.(verificationPreEnablingBenchmark); ok {
					b.EnableVerification()
				}
			}

			b.Run()

			if r.Verify {
				b.Verify()
			}
			wg.Done()
		}(b, &wg)
	}
	wg.Wait()

	r.GPUDriver.Terminate()
	r.Engine.Finished()

	//r.reportStats()

	atexit.Exit(0)
}

func (r *Runner) reportStats() {
	r.reportExecutionTime()
	r.reportInstCount()
	r.reportCacheLatency()
	r.reportCacheHitRate()
	r.reportTLBLatency()
	r.reportTLBConditionalStats()
	r.reportPageWalkLatency()
	r.reportMMUConditionalStats()
	r.reportActiveWalkerCount()
	r.reportDRAMLatency()
	r.reportTLBHitRate()
	r.reportPWCHitRate()
	r.reportRDMATransactionCount()
	r.reportCDMATransactionCount()
	r.reportRTUTransactionCount()
	// r.ReportRTUCoalescing()
	r.reportDRAMTransactionCount()
	r.reportL2TLBQueueImbalance()
	r.reportPageWalkerImbalance()
	r.reportL2TLBEntropy()
	r.reportReferenceTracing()
	r.reportTLBSetMissTracing()
	r.reportTLBMSHRStallTracing()
	r.reportTLBReqStalls()

	if *platformType != "ideal" {
		r.reportTranslationReqLatency()
	}
	r.reportAddressTranslatorLatency()
	r.dumpMetrics()
}

func (r *Runner) reportInstCount() {
	for _, t := range r.InstCountTracers {
		r.metricsCollector.Collect(
			t.cu.Name(), "inst_count", float64(t.tracer.count))
	}
	for _, t := range r.InstCountTracers {
		r.metricsCollector.Collect(
			t.cu.Name(), "last_WG_finish_time", float64(t.tracer.lastWGFinishTime))
	}
}

func (r *Runner) reportExecutionTime() {
	if r.Timing {
		r.metricsCollector.Collect(
			r.GPUDriver.Name(),
			"kernel_time", float64(r.KernelTimeCounter.BusyTime()))
		r.metricsCollector.Collect(
			r.GPUDriver.Name(),
			"total_time", float64(r.Engine.CurrentTime()))
		kernelTimes := r.KernelTimeCounter.KernelTimes()
		for i, k := range kernelTimes {
			r.metricsCollector.Collect(
				r.GPUDriver.Name(),
				"kernel_time "+strconv.Itoa(i), float64(k))
		}

		for i, c := range r.PerGPUKernelTimeCounter {
			r.metricsCollector.Collect(
				r.GPUDriver.GPUs[i].CommandProcessor.Name(),
				"kernel_time", float64(c.BusyTime()))
			kernelTimes := r.KernelTimeCounter.KernelTimes()
			for j, k := range kernelTimes {
				r.metricsCollector.Collect(
					r.GPUDriver.GPUs[i].CommandProcessor.Name(),
					"kernel_time "+strconv.Itoa(j), float64(k))
			}

			// Only should be called at Exit()
			r.KernelTimeCounter.TerminateAllTasks(r.Engine.CurrentTime())

			kernelTimesForceStop := r.KernelTimeCounter.KernelTimesForceStop(r.Engine.CurrentTime())
			for j, k := range kernelTimesForceStop {
				r.metricsCollector.Collect(
					r.GPUDriver.GPUs[i].CommandProcessor.Name(),
					"kernel_time (force stop) "+strconv.Itoa(j), float64(k))
			}
		}
	}
}

func (r *Runner) reportCacheLatency() {
	for _, tracer := range r.CacheLatencyTracers {
		if tracer.tracer.AverageTime() == 0 {
			continue
		}

		r.metricsCollector.Collect(
			tracer.cache.Name(),
			"req_average_latency",
			float64(tracer.tracer.AverageTime()),
		)
	}
	for _, tracer := range r.L2PipelineLatencyTracers {
		if tracer.tracer.AverageTime() == 0 {
			continue
		}

		r.metricsCollector.Collect(
			tracer.pipeline.Name(),
			"req_average_pipeline_latency",
			float64(tracer.tracer.AverageTime()),
		)
	}
}

func (r *Runner) reportTLBLatency() {
	for _, tracer := range r.TLBLatencyTracers {
		if tracer.tracer.AverageTime() == 0 {
			continue
		}

		r.metricsCollector.Collect(
			tracer.tlb.Name(),
			"req_average_latency",
			float64(tracer.tracer.AverageTime()),
		)
	}
	for _, tracer := range r.DownTLBLatencyTracers {
		if tracer.tracer.AverageTime() == 0 {
			continue
		}

		r.metricsCollector.Collect(
			tracer.tlb.Name(),
			"down_req_average_latency",
			float64(tracer.tracer.AverageTime()),
		)
	}

	for i, tracer := range r.L2TLBBufLenTracers {
		if tracer.AverageCount() == 0 {
			continue
		}
		r.metricsCollector.Collect(
			"L2TLB"+strconv.Itoa(i),
			"average_buf_len",
			float64(tracer.AverageCount()),
		)
	}
	for i, tracer := range r.L2TLBBufLenG0Tracers {
		if tracer.AverageCount() == 0 {
			continue
		}
		r.metricsCollector.Collect(
			"L2TLB"+strconv.Itoa(i),
			"average_buf_len_g0",
			float64(tracer.AverageCount()),
		)
	}
	for i, tracer := range r.L2TLBCoalesceAddrTracers {
		if tracer.AverageCount() == 0 {
			continue
		}
		r.metricsCollector.Collect(
			"L2TLB"+strconv.Itoa(i),
			"coalesce_addr",
			float64(tracer.AverageCount()),
		)
	}
	for i, tracer := range r.L2TLBCoalesceTracers {
		if tracer.AverageCount() == 0 {
			continue
		}
		r.metricsCollector.Collect(
			"L2TLB"+strconv.Itoa(i),
			"coalesce",
			float64(tracer.AverageCount()),
		)
	}

	for _, tracer := range r.TLBPipelineLatencyTracers {
		if tracer.tracer.AverageTime() == 0 {
			continue
		}

		r.metricsCollector.Collect(
			tracer.pipeline.Name(),
			"req_average_pipeline_latency",
			float64(tracer.tracer.AverageTime()),
		)
	}

	for i, tracer := range r.L2TLBMSHRLenTracers {
		if tracer.AverageCount() == 0 {
			continue
		}
		r.metricsCollector.Collect(
			"L2TLB"+strconv.Itoa(i),
			"average_mshr_len",
			float64(tracer.AverageCount()),
		)
	}
	for i, tracer := range r.L2TLBMSHRLenG0Tracers {
		if tracer.AverageCount() == 0 {
			continue
		}
		r.metricsCollector.Collect(
			"L2TLB"+strconv.Itoa(i),
			"average_mshr_len_g0",
			float64(tracer.AverageCount()),
		)
	}

	for i, tracer := range r.L2TLBMSHRUniqLenTracers {
		if tracer.AverageCount() == 0 {
			continue
		}
		r.metricsCollector.Collect(
			"L2TLB"+strconv.Itoa(i),
			"average_mshr_uniq_len",
			float64(tracer.AverageCount()),
		)
	}

	for i, tracer := range r.L2TLBMSHRUniqLenG0Tracers {
		if tracer.AverageCount() == 0 {
			continue
		}
		r.metricsCollector.Collect(
			"L2TLB"+strconv.Itoa(i),
			"average_mshr_uniq_len_g0",
			float64(tracer.AverageCount()),
		)
	}

}

func (r *Runner) reportTLBConditionalStats() {
	for _, tracer := range r.ConditionalRemoteTLBLatencyTracers {
		tlbTracer := tracer.tracer
		for _, stepName := range tlbTracer.GetStepNames() {
			if tlbTracer.AverageTime(stepName) == 0 {
				continue
			}

			r.metricsCollector.Collect(
				tracer.tlb.Name(),
				stepName+"-latency",
				float64(tlbTracer.AverageTime(stepName)),
			)

			r.metricsCollector.Collect(
				tracer.tlb.Name(),
				stepName+"-num",
				float64(tlbTracer.TotalCount(stepName)),
			)
		}
	}
}

func (r *Runner) reportPageWalkLatency() {
	for _, tracer := range r.PageWalkLatencyTracers {
		if tracer.tracer.AverageTime() == 0 {
			continue
		}

		r.metricsCollector.Collect(
			tracer.mmu.Name(),
			"req_average_latency",
			float64(tracer.tracer.AverageTime()),
		)
	}

}

func (r *Runner) reportMMUConditionalStats() {
	for _, tracer := range r.MMUMemoryLatencyTracers {
		mmuTracer := tracer.tracer
		for _, stepName := range mmuTracer.GetStepNames() {
			if mmuTracer.AverageTime(stepName) == 0 {
				continue
			}

			r.metricsCollector.Collect(
				tracer.mmu.Name(),
				stepName+"-latency",
				float64(mmuTracer.AverageTime(stepName)),
			)

			r.metricsCollector.Collect(
				tracer.mmu.Name(),
				stepName+"-num",
				float64(mmuTracer.TotalCount(stepName)),
			)
		}
	}
}

func (r *Runner) reportActiveWalkerCount() {
	for _, tracer := range r.ActivePageWalkerTracers {
		if tracer.tracer.AverageCount() == 0 {
			continue
		}
		r.metricsCollector.Collect(
			tracer.mmu.Name(),
			"average active walkers",
			float64(tracer.tracer.AverageCount()),
		)
	}
}

func (r *Runner) ReportRTUCoalescing() {
	for i, prefix := range [2]string{"incoming", "outgoing"} {
		for j, statName := range [4]string{"buflen", "buflen-g0", "coalesceable-addr", "replication-count"} {
			for k, gpu := range r.GPUDriver.GPUs {
				for l, rtu := range gpu.RemoteAddressTranslationUnits {
					tracer := r.RTUCoalescingTracers[i][j][k][l]
					//	if tracer.AverageCount() == 0 {
					//		continue
					//	}
					fullStatName := prefix + "-" + statName
					r.metricsCollector.Collect(
						rtu.Name(),
						fullStatName,
						float64(tracer.AverageCount()),
					)
					r.metricsCollector.Collect(
						rtu.Name(),
						fullStatName+"-count",
						float64(tracer.TotalCount()),
					)
				}

			}
		}
	}
}

func (r *Runner) reportDRAMLatency() {
	for _, tracer := range r.DRAMLatencyTracers {
		if tracer.tracer.AverageTime() == 0 {
			continue
		}

		r.metricsCollector.Collect(
			tracer.dram.Name(),
			"req_average_latency",
			float64(tracer.tracer.AverageTime()),
		)
	}
}

func (r *Runner) reportCacheHitRate() {
	for _, tracer := range r.CacheHitRateTracers {
		//fmt.Println(tracer.cache.Name())
		readHit := tracer.tracer.GetStepCount("read-hit")
		readMiss := tracer.tracer.GetStepCount("read-miss")
		readMSHRHit := tracer.tracer.GetStepCount("read-mshr-hit")
		writeHit := tracer.tracer.GetStepCount("write-hit")
		writeMiss := tracer.tracer.GetStepCount("write-miss")
		writeMSHRHit := tracer.tracer.GetStepCount("write-mshr-hit")

		totalTransaction := readHit + readMiss + readMSHRHit +
			writeHit + writeMiss + writeMSHRHit

		if totalTransaction == 0 {
			continue
		}

		r.metricsCollector.Collect(
			tracer.cache.Name(), "read-hit", float64(readHit))
		r.metricsCollector.Collect(
			tracer.cache.Name(), "read-miss", float64(readMiss))
		r.metricsCollector.Collect(
			tracer.cache.Name(), "read-mshr-hit", float64(readMSHRHit))
		r.metricsCollector.Collect(
			tracer.cache.Name(), "write-hit", float64(writeHit))
		r.metricsCollector.Collect(
			tracer.cache.Name(), "write-miss", float64(writeMiss))
		r.metricsCollector.Collect(
			tracer.cache.Name(), "write-mshr-hit", float64(writeMSHRHit))
	}
}

func (r *Runner) reportTLBHitRate() {
	for _, tracer := range r.TLBHitRateTracers {
		//fmt.Println(tracer.tlb.Name())
		tlbHit := tracer.tracer.GetStepCount("tlb-hit")
		tlbMiss := tracer.tracer.GetStepCount("tlb-miss")
		tlbMSHRHit := tracer.tracer.GetStepCount("tlb-mshr-hit")
		totalTransaction := tlbHit + tlbMiss + tlbMSHRHit
		if totalTransaction == 0 {
			continue
		}
		//tlb := tracer.tlb.(*tlb.TLB)
		tlbName := tracer.tlb.Name()
		r.metricsCollector.Collect(
			tlbName, "tlb-hit", float64(tlbHit))
		r.metricsCollector.Collect(
			tlbName, "tlb-miss", float64(tlbMiss))
		r.metricsCollector.Collect(
			tlbName, "tlb-mshr-hit", float64(tlbMSHRHit))

	}
}

func (r *Runner) reportPWCHitRate() {
	for _, tracer := range r.PWCHitRateTracers {
		pwcHit0 := tracer.tracer.GetStepCount("pwc-hit-level0")
		pwcHit1 := tracer.tracer.GetStepCount("pwc-hit-level1")
		pwcHit2 := tracer.tracer.GetStepCount("pwc-hit-level2")
		pwcHit3 := tracer.tracer.GetStepCount("pwc-hit-level3")
		remoteMemReq := tracer.tracer.GetStepCount("page_walk_req_remote")
		localMemReq := tracer.tracer.GetStepCount("page_walk_req_local")
		lastLevelRemoteMemReq := tracer.tracer.GetStepCount("pw-level-3-remote-reqs")
		totalTransaction := pwcHit3 + pwcHit2 + pwcHit1 + pwcHit0

		if totalTransaction == 0 {
			continue
		}

		r.metricsCollector.Collect(
			tracer.mmu.Name(), "pwc-miss", float64(pwcHit0))
		r.metricsCollector.Collect(
			tracer.mmu.Name(), "pwc-hit-level1", float64(pwcHit1))
		r.metricsCollector.Collect(
			tracer.mmu.Name(), "pwc-hit-level2", float64(pwcHit2))
		r.metricsCollector.Collect(
			tracer.mmu.Name(), "pwc-hit-level3", float64(pwcHit3))
		r.metricsCollector.Collect(
			tracer.mmu.Name(), "page_walk_req_local", float64(localMemReq))
		r.metricsCollector.Collect(
			tracer.mmu.Name(), "page_walk_req_remote", float64(remoteMemReq))
		r.metricsCollector.Collect(
			tracer.mmu.Name(), "pw-level-3-remote-reqs", float64(lastLevelRemoteMemReq))
	}
}

func (r *Runner) reportRDMATransactionCount() {
	for _, t := range r.RDMATransactionCounters {
		r.metricsCollector.Collect(
			t.rdmaEngine.Name(),
			"outgoing_trans_count",
			float64(t.outgoingTracer.TotalCount()),
		)
		r.metricsCollector.Collect(
			t.rdmaEngine.Name(),
			"incoming_trans_count",
			float64(t.incomingTracer.TotalCount()),
		)
	}
}

func (r *Runner) reportCDMATransactionCount() {
	for _, t := range r.CDMATransactionCounters {
		r.metricsCollector.Collect(
			t.rdmaEngine.Name(),
			"outgoing_trans_count",
			float64(t.outgoingTracer.TotalCount()),
		)
		r.metricsCollector.Collect(
			t.rdmaEngine.Name(),
			"outgoing_trans_latency",
			float64(t.outgoingTracer.AverageTime()),
		)
		r.metricsCollector.Collect(
			t.rdmaEngine.Name(),
			"incoming_trans_count",
			float64(t.incomingTracer.TotalCount()),
		)
		r.metricsCollector.Collect(
			t.rdmaEngine.Name(),
			"incoming_trans_latency",
			float64(t.incomingTracer.AverageTime()),
		)
	}
	for _, t := range r.CDMAAccessTracers {
		cdmaEngine := t.cdmaEngine
		cdmaEngineTracer := t.tracer
		for _, step := range cdmaEngineTracer.GetStepNames() {
			r.metricsCollector.Collect(
				cdmaEngine.Name(),
				"inter_chiplet_traffic: "+step,
				float64(cdmaEngineTracer.GetStepCount(step)),
			)
		}
	}

	for _, t := range r.PageTransactionCounters {
		r.metricsCollector.Collect(
			t.rdmaEngine.Name(),
			"outgoing_trans_count",
			float64(t.outgoingTracer.TotalCount()),
		)
		r.metricsCollector.Collect(
			t.rdmaEngine.Name(),
			"outgoing_trans_latency",
			float64(t.outgoingTracer.AverageTime()),
		)
		r.metricsCollector.Collect(
			t.rdmaEngine.Name(),
			"incoming_trans_count",
			float64(t.incomingTracer.TotalCount()),
		)
		r.metricsCollector.Collect(
			t.rdmaEngine.Name(),
			"incoming_trans_latency",
			float64(t.incomingTracer.AverageTime()),
		)
	}
	for _, t := range r.PageAccessTracers {
		cdmaEngine := t.cdmaEngine
		cdmaEngineTracer := t.tracer
		for _, step := range cdmaEngineTracer.GetStepNames() {
			r.metricsCollector.Collect(
				cdmaEngine.Name(),
				"inter_chiplet_traffic: "+step,
				float64(cdmaEngineTracer.GetStepCount(step)),
			)
		}
	}
}

func (r *Runner) reportRTUTransactionCount() {
	for _, t := range r.RTUTransactionCounters {
		r.metricsCollector.Collect(
			t.rtuUnit.Name(),
			"outgoing_trans_count",
			float64(t.outgoingTracer.TotalCount()),
		)
		r.metricsCollector.Collect(
			t.rtuUnit.Name(),
			"outgoing_trans_latency",
			float64(t.outgoingTracer.AverageTime()),
		)
		r.metricsCollector.Collect(
			t.rtuUnit.Name(),
			"incoming_trans_count",
			float64(t.incomingTracer.TotalCount()),
		)
		r.metricsCollector.Collect(
			t.rtuUnit.Name(),
			"incoming_trans_latency",
			float64(t.incomingTracer.AverageTime()),
		)
	}

	for _, t := range r.RTUAccessTracers {
		rtu := t.rtu
		rtuTracer := t.tracer
		for _, step := range rtuTracer.GetStepNames() {
			r.metricsCollector.Collect(
				rtu.Name(),
				"inter_chiplet_tlb_traffic: "+step,
				float64(rtuTracer.GetStepCount(step)),
			)
		}
	}
}

func (r *Runner) reportDRAMTransactionCount() {
	for _, t := range r.DRAMTransactionCounters {
		r.metricsCollector.Collect(
			t.dram.Name(),
			"trans_count",
			float64(t.tracer.TotalCount()),
		)
	}
}

func (r *Runner) reportTranslationReqLatency() {
	tracer := r.TranslationReqTracer
	if tracer == nil {
		return
	}
	for _, stepName := range tracer.GetStepNames() {
		// if tracer.AverageTime(stepName) == 0 {
		// 	continue
		// }
		r.metricsCollector.Collect(
			"trace-trans-req",
			stepName+"-latency",
			float64(tracer.AverageTime(stepName)),
		)
		r.metricsCollector.Collect(
			"trace-trans-req",
			stepName+"-num",
			float64(tracer.TotalCount(stepName)),
		)
	}
}

func (r *Runner) reportAddressTranslatorLatency() {
	for _, t := range r.AddressTranslatorLatencyTracers {
		addressTranslator := t.addressTranslator
		tracer := t.tracer
		r.metricsCollector.Collect(
			addressTranslator.Name(),
			"req-latency",
			float64(tracer.AverageTime("req-latency")),
		)
		r.metricsCollector.Collect(
			addressTranslator.Name(),
			"req-num",
			float64(tracer.TotalCount("req-latency")),
		)
		r.metricsCollector.Collect(
			addressTranslator.Name(),
			"translation-latency",
			float64(tracer.AverageTime("translation-latency")),
		)
		r.metricsCollector.Collect(
			addressTranslator.Name(),
			"translation-num",
			float64(tracer.TotalCount("translation-latency")),
		)
	}
}

func (r *Runner) reportL2TLBQueueImbalance() {
	for _, t := range r.L2TLBQueueingImbalanceTracers {
		r.metricsCollector.Collect(
			t.tlb.Name(),
			"imbalance",
			float64(t.tracer.ReportImbalance()),
		)
		r.metricsCollector.Collect(
			t.tlb.Name(),
			"imbalance-count",
			float64(t.tracer.ReportCount()),
		)
		for i, bucket := range t.tracer.ReportImbalanceBuckets() {
			r.metricsCollector.Collect(
				t.tlb.Name(),
				"imbalance-"+strconv.Itoa(i*10),
				float64(bucket),
			)
		}
	}
}

func (r *Runner) reportPageWalkerImbalance() {
	for _, t := range r.PageWalkerImbalanceTracers {
		r.metricsCollector.Collect(
			t.mmu.Name(),
			"imbalance",
			float64(t.tracer.ReportImbalance()),
		)
		r.metricsCollector.Collect(
			t.mmu.Name(),
			"imbalance-count",
			float64(t.tracer.ReportCount()),
		)
	}
}

func (r *Runner) reportL2TLBEntropy() {
	for _, t := range r.EntropyTracers {
		for i := 12; i < 28; i++ {
			r.metricsCollector.Collect(
				t.tlb.Name(),
				"entropy_"+strconv.Itoa(i),
				float64(t.tracer.ReportEntropy()[i]),
			)
		}
	}
}

func (r *Runner) reportReferenceTracing() {
	for _, t := range r.RemoteReferenceCountTracers {
		r.metricsCollector.Collect(
			t.rtuUnit.Name(),
			"reference_count_avg_",
			float64(t.tracer.ReturnAverageReferenceCount()),
		)
	}
	for _, t := range r.RemoteReferenceCountTracers {
		r.metricsCollector.Collect(
			t.rtuUnit.Name(),
			"reference_interval_avg_",
			float64(t.tracer.ReturnAverageReferenceInterval()),
		)
	}
	for _, t := range r.RemoteReferenceCountTracers {
		t.tracer.CalculateHistogramSummary()
		for i := 0; i < 64; i++ {
			r.metricsCollector.Collect(
				t.rtuUnit.Name(),
				"reference_histogram_summary_"+strconv.Itoa(i),
				float64(t.tracer.ReturnHistogramSummary(uint(i))),
			)
		}
	}
	for _, t := range r.RemoteReferenceCountTracers {
		for i := 0; i < 32; i++ {
			r.metricsCollector.Collect(
				t.rtuUnit.Name(),
				"reference_count_histogram_summary_"+strconv.Itoa(i),
				float64(t.tracer.ReturnCountHistogramSummary(uint(i))),
			)
		}
	}

}

func (r *Runner) reportTLBSetMissTracing() {
	for _, t := range r.TLBSetMissTracers {
		r.metricsCollector.Collect(
			t.tlb.Name(),
			"miss_count_",
			float64(t.tracer.ReturnAverageReferenceCount()),
		)
	}
	for _, t := range r.TLBSetMissTracers {
		r.metricsCollector.Collect(
			t.tlb.Name(),
			"miss_interval_avg_",
			float64(t.tracer.ReturnAverageReferenceInterval()),
		)
	}
	for _, t := range r.TLBSetMissTracers {
		t.tracer.CalculateHistogramSummary()
		for i := 0; i < 64; i++ {
			r.metricsCollector.Collect(
				t.tlb.Name(),
				"miss_histogram_summary_"+strconv.Itoa(i),
				float64(t.tracer.ReturnHistogramSummary(uint(i))),
			)
		}
	}
	for _, t := range r.TLBSetMissTracers {
		for i := 0; i < 32; i++ {
			r.metricsCollector.Collect(
				t.tlb.Name(),
				"miss_count_histogram_summary_"+strconv.Itoa(i),
				float64(t.tracer.ReturnCountHistogramSummary(uint(i))),
			)
		}
	}
}

func (r *Runner) reportTLBMSHRStallTracing() {
	for _, t := range r.TLBMSHRStallTracers {
		r.metricsCollector.Collect(
			t.tlb.Name(),
			"mshr_stall",
			float64(t.tracer.TotalTime()),
		)
	}
}

func (r *Runner) reportTLBReqStalls() {
	for _, t := range r.TLBReqStallTracers {
		tracer := t.tracer
		tlb := t.tlb
		for _, stepName := range tracer.GetStepNames() {
			r.metricsCollector.Collect(
				tlb.Name(),
				stepName,
				float64(tracer.GetStepCount(stepName)),
			)
		}
	}
}

func (r *Runner) dumpMetrics() {
	r.metricsCollector.Dump(*filenameFlag)
}
