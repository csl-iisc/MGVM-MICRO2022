package main

import (
	"flag"

	"gitlab.com/akita/mgpusim/benchmarks/rodinia/nw"
	"gitlab.com/akita/mgpusim/samples/runner"
)

// var nFlag = flag.Int("n", 128, "Dunno")
// var kFlag = flag.Int("k", 1, "Dunno")

func main() {
	flag.Parse()

	runner := new(runner.Runner).ParseFlag().Init()

	benchmark := nw.NewBenchmark(runner.GPUDriver)
	// benchmark.N = *nFlag
	// benchmark.K = *kFlag

	runner.AddBenchmark(benchmark)

	runner.Run()
}
