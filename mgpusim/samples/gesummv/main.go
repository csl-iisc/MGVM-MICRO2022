package main

import (
	"flag"

	_ "net/http/pprof"

	"gitlab.com/akita/mgpusim/benchmarks/polybench/gesummv"
	"gitlab.com/akita/mgpusim/samples/runner"
)

var nFlag = flag.Uint("n", 512, "The height of the first matrix.")

func main() {
	flag.Parse()

	runner := new(runner.Runner).ParseFlag().Init()

	benchmark := gesummv.NewBenchmark(runner.GPUDriver)
	benchmark.N = int(*nFlag)

	runner.AddBenchmark(benchmark)

	runner.Run()
}
