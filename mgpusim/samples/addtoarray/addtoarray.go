package main

import (
	"flag"

	"gitlab.com/akita/mgpusim/benchmarks/microbenchmarks/addtoarray"
	"gitlab.com/akita/mgpusim/samples/runner"
)

func main() {
	flag.Parse()

	runner := new(runner.Runner).ParseFlag().Init()

	for i := 0; i < 1; i++ {
		benchmark := addtoarray.NewBenchmark(runner.GPUDriver)

		runner.AddBenchmark(benchmark)
	}

	runner.Run()
}
