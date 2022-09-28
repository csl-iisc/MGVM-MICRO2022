package main

import (
	"flag"

	"gitlab.com/akita/mgpusim/benchmarks/polybench/correlation"
	"gitlab.com/akita/mgpusim/samples/runner"
)

var mFlag = flag.Int("num_m", 512, "Dunno")
var nFlag = flag.Int("num_n", 512, "Dunno")

func main() {
	flag.Parse()

	runner := new(runner.Runner).ParseFlag().Init()

	benchmark := correlation.NewBenchmark(runner.GPUDriver)
	benchmark.M = *mFlag
	benchmark.N = *nFlag

	runner.AddBenchmark(benchmark)

	runner.Run()
}
