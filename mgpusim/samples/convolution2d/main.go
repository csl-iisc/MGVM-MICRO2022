package main

import (
	"flag"

	"gitlab.com/akita/mgpusim/benchmarks/polybench/convolution2d"
	"gitlab.com/akita/mgpusim/samples/runner"
)

var niFlag = flag.Int("ni", 32, "Dunno")
var njFlag = flag.Int("nj", 32, "Dunno")

func main() {
	flag.Parse()

	runner := new(runner.Runner).ParseFlag().Init()

	benchmark := convolution2d.NewBenchmark(runner.GPUDriver)
	benchmark.NI = *niFlag
	benchmark.NJ = *njFlag

	runner.AddBenchmark(benchmark)

	runner.Run()
}
