package main

import (
	"flag"

	"gitlab.com/akita/mgpusim/benchmarks/amdappsdk/reduction"
	"gitlab.com/akita/mgpusim/samples/runner"
)

var Length = flag.Int("length", 4096, "The number of rows in the input matrix.")

func main() {
	flag.Parse()

	runner := new(runner.Runner).ParseFlag().Init()

	benchmark := reduction.NewBenchmark(runner.GPUDriver)
	benchmark.Length = uint32(*Length)

	runner.AddBenchmark(benchmark)

	runner.Run()
}
