package main

import (
	"flag"

	"gitlab.com/akita/mgpusim/benchmarks/polybench/gramschmidt"
	"gitlab.com/akita/mgpusim/samples/runner"
)

var iFlag = flag.Int("ni", 1024, "Dunno")
var jFlag = flag.Int("nj", 1024, "Dunno")
var kFlag = flag.Int("k", 1, "Dunno")

func main() {
	flag.Parse()

	runner := new(runner.Runner).ParseFlag().Init()

	benchmark := gramschmidt.NewBenchmark(runner.GPUDriver)
	benchmark.NI = *iFlag
	benchmark.NJ = *jFlag
	benchmark.K = *kFlag

	runner.AddBenchmark(benchmark)

	runner.Run()
}
