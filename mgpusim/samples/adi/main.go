package main

import (
	"flag"

	"gitlab.com/akita/mgpusim/benchmarks/polybench/adi"
	"gitlab.com/akita/mgpusim/samples/runner"
)

var nFlag = flag.Int("n", 256, "The width of the matrix.")
var steps = flag.Int("steps", 1, "The height of the matrix.")

func main() {
	flag.Parse()

	runner := new(runner.Runner).ParseFlag().Init()

	benchmark := adi.NewBenchmark(runner.GPUDriver)
	benchmark.N = *nFlag
	benchmark.Steps = *steps

	runner.AddBenchmark(benchmark)

	runner.Run()
}
