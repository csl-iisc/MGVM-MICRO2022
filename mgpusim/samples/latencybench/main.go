package main

import (
	"flag"

	"gitlab.com/akita/mgpusim/benchmarks/microbenchmarks/latencybench"
	"gitlab.com/akita/mgpusim/samples/runner"
)

// use blocks = 1281, remotestart = 1280 and loopcount sufficient to spill to next chiplet.
// use 52884

var length = flag.Int("length", 1024, "The length of array to sort.")
var stride = flag.Int("stride", 1024, "The stride.")
var start = flag.Int("start", 0, "The start location to stride.")
var end = flag.Int("end", 1024, "The end location  to stride.")
var threads = flag.Int("threads", 1, "The number of threads.")
var blocks = flag.Int("blocks", 1281, "The number of thread blocks.")
var loopcount = flag.Int("loopcount", 256, "Number of iteration of loop.")
var remotestart = flag.Int("remotestart", 1280,
	"Number of blocks to ensure next chiplet .")
var benchtype = flag.String("benchtype", "idle", "Type of micro benchmark.")

func main() {
	flag.Parse()

	runner := new(runner.Runner).ParseFlag().Init()

	benchmark := latencybench.NewBenchmark(runner.GPUDriver)
	benchmark.Length = (uint32)(*length)
	benchmark.Stride = (uint32)(*stride)
	benchmark.Start = (uint32)(*start)
	benchmark.End = (uint32)(*end)
	benchmark.Threads = (uint32)(*threads)
	benchmark.Blocks = (uint32)(*blocks)
	benchmark.LoopCount = (uint32)(*loopcount)
	benchmark.RemoteStart = (uint32)(*remotestart)
	benchmark.BenchType = (string)(*benchtype)

	runner.AddBenchmark(benchmark)

	runner.Run()
}
