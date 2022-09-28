package main

import (
	"flag"

	"gitlab.com/akita/mgpusim/benchmarks/polybench/jacobi1d"
	"gitlab.com/akita/mgpusim/samples/runner"
)

var nFlag = flag.Int("n", 4096, "Dunno")
var stepsFlag = flag.Int("steps", 32, "Dunno")

func main() {
	flag.Parse()

	runner := new(runner.Runner).ParseFlag().Init()

	benchmark := jacobi1d.NewBenchmark(runner.GPUDriver)
	benchmark.N = *nFlag
	benchmark.Steps = *stepsFlag

	runner.AddBenchmark(benchmark)

	runner.Run()
}
