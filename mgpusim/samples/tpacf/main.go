package main

import (
	"flag"

	"gitlab.com/akita/mgpusim/benchmarks/parboil/tpacf"
	"gitlab.com/akita/mgpusim/samples/runner"
)

func main() {
	flag.Parse()

	runner := new(runner.Runner).ParseFlag().Init()

	benchmark := tpacf.NewBenchmark(runner.GPUDriver)

	runner.AddBenchmark(benchmark)

	runner.Run()
}
