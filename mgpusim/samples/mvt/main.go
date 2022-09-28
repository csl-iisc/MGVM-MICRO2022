package main

import (
	"flag"

	"gitlab.com/akita/mgpusim/benchmarks/polybench/mvt"
	"gitlab.com/akita/mgpusim/samples/runner"
)

var nFlag = flag.Int("n", 1024, "Dunno")

func main() {
	flag.Parse()

	runner := new(runner.Runner).ParseFlag().Init()

	benchmark := mvt.NewBenchmark(runner.GPUDriver)
	benchmark.N = *nFlag

	runner.AddBenchmark(benchmark)

	runner.Run()
}
