package main

import (
	"flag"

	"gitlab.com/akita/mgpusim/benchmarks/polybench/fdtd2d"
	"gitlab.com/akita/mgpusim/samples/runner"
)

var nxFlag = flag.Int("nx", 4096, "Dunno")
var nyFlag = flag.Int("ny", 2048, "Dunno")
var tMax = flag.Int("max_steps", 1, "Dunno")

func main() {
	flag.Parse()

	runner := new(runner.Runner).ParseFlag().Init()

	benchmark := fdtd2d.NewBenchmark(runner.GPUDriver)
	benchmark.NX = *nxFlag
	benchmark.NY = *nyFlag
	benchmark.TMax = *tMax

	runner.AddBenchmark(benchmark)

	runner.Run()
}
