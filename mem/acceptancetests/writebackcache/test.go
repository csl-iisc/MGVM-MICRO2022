package main

import (
	"fmt"
	"log"
	"math/rand"

	"time"

	"flag"

	"os"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem"
	"gitlab.com/akita/mem/acceptancetests"
	"gitlab.com/akita/mem/cache"
	"gitlab.com/akita/mem/cache/writeback"
	"gitlab.com/akita/mem/idealmemcontroller"
	"gitlab.com/akita/mem/trace"
	"gitlab.com/akita/util/tracing"
)

var seedFlag = flag.Int64("seed", 0, "Random Seed")
var numAccessFlag = flag.Int("num-access", 100000,
	"Number of accesses to generate")
var maxAddressFlag = flag.Uint64("max-address", 1048576, "Address range to use")
var traceFileFlag = flag.String("trace", "", "Trace file")
var traceWithStdoutFlag = flag.Bool("trace-stdout", false, "Trace with stdout")
var parallelFlag = flag.Bool("parallel", false, "Test with parallel engine")

var engine akita.Engine
var agent *acceptancetests.MemAccessAgent

func main() {
	flag.Parse()

	initSeed()
	buildEnvironment()
	runSimulation()
	allMsgsMustBeSent()
}

func initSeed() {
	var seed int64
	if *seedFlag == 0 {
		seed = time.Now().UnixNano()
	} else {
		seed = *seedFlag
	}
	fmt.Fprintf(os.Stderr, "Seed %d\n", seed)
	rand.Seed(seed)
}

func buildEnvironment() {
	if *parallelFlag {
		engine = akita.NewParallelEngine()
	} else {
		engine = akita.NewSerialEngine()
	}
	//engine.AcceptHook(akita.NewEventLogger(log.New(os.Stdout, "", 0)))

	conn := akita.NewDirectConnection("conn", engine, 1*akita.GHz)

	agent = acceptancetests.NewMemAccessAgent(engine)
	agent.MaxAddress = *maxAddressFlag
	agent.WriteLeft = *numAccessFlag
	agent.ReadLeft = *numAccessFlag

	lowModuleFinder := new(cache.SingleLowModuleFinder)
	builder := writeback.MakeBuilder().
		WithEngine(engine).
		WithLowModuleFinder(lowModuleFinder).
		WithByteSize(16 * mem.KB).
		WithLog2BlockSize(6).
		WithWayAssociativity(4).
		WithNumMSHREntry(4).
		WithNumReqPerCycle(16)
	writeBackCache := builder.Build("Cache")

	if *traceWithStdoutFlag {
		logger := log.New(os.Stdout, "", 0)
		tracer := trace.NewTracer(logger)
		tracing.CollectTrace(writeBackCache, tracer)
	} else if *traceFileFlag != "" {
		traceFile, _ := os.Create(*traceFileFlag)
		logger := log.New(traceFile, "", 0)
		tracer := trace.NewTracer(logger)
		tracing.CollectTrace(writeBackCache, tracer)
	}

	dram := idealmemcontroller.New("DRAM", engine, 4*mem.GB)
	lowModuleFinder.LowModule = dram.ToTop

	agent.LowModule = writeBackCache.TopPort

	conn.PlugIn(agent.ToMem, 16)
	conn.PlugIn(writeBackCache.BottomPort, 16)
	conn.PlugIn(writeBackCache.TopPort, 16)
	conn.PlugIn(dram.ToTop, 16)

	agent.TickLater(0)
}

func runSimulation() {
	err := engine.Run()
	if err != nil {
		panic(err)
	}
}

func allMsgsMustBeSent() {
	if len(agent.PendingWriteReq) > 0 || len(agent.PendingReadReq) > 0 {
		panic("Not all req returned")
	}

	if agent.WriteLeft > 0 || agent.ReadLeft > 0 {
		panic("more requests to send")
	}
}
