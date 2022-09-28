package main

import (
	"flag"

	"gitlab.com/akita/mgpusim/benchmarks/shoc/fft"
	"gitlab.com/akita/mgpusim/samples/runner"
)

var mb = flag.Int("MB", 8, "data size (in megabytes)")
var passes = flag.Int("passes", 2, "number of passes")

func main() {
	flag.Parse()

	runner := new(runner.Runner).ParseFlag().Init()

	benchmark := fft.NewBenchmark(runner.GPUDriver)
	benchmark.Bytes = int32(*mb)
	benchmark.Passes = int32(*passes)

	runner.AddBenchmark(benchmark)

	runner.Run()
}
