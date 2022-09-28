package main

import (
	"flag"

	"gitlab.com/akita/mgpusim/benchmarks/mafiaports/gups"
	"gitlab.com/akita/mgpusim/samples/runner"
)

func main() {
	flag.Parse()

	runner := new(runner.Runner).ParseFlag().Init()

	benchmark := gups.NewBenchmark(runner.GPUDriver)

	runner.AddBenchmark(benchmark)

	runner.Run()
}
