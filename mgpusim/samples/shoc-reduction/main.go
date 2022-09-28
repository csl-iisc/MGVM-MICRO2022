package main

import (
	"flag"

	"gitlab.com/akita/mgpusim/benchmarks/shoc/reduction"
	"gitlab.com/akita/mgpusim/samples/runner"
)

var Size = flag.Int("Size", 25165824, "The number of rows in the input matrix.")
var Iterations = flag.Int("Iterations", 1, "The number of rows in the input matrix.")

func main() {
	flag.Parse()

	runner := new(runner.Runner).ParseFlag().Init()

	benchmark := reduction.NewBenchmark(runner.GPUDriver)
	benchmark.Size = uint32(*Size)
	benchmark.Iterations = uint32(*Iterations)

	runner.AddBenchmark(benchmark)

	runner.Run()
}
