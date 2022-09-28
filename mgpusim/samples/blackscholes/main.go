package main

import (
	"flag"

	"gitlab.com/akita/mgpusim/benchmarks/amdappsdk/blackscholes"
	"gitlab.com/akita/mgpusim/samples/runner"
)

var Width = flag.Int("width", 1024, "The number of rows in the input matrix.")
var Height = flag.Int("height", 1024, "The number of rows in the input matrix.")

func main() {
	flag.Parse()

	runner := new(runner.Runner).ParseFlag().Init()

	benchmark := blackscholes.NewBenchmark(runner.GPUDriver)
	benchmark.Width = *Width
	benchmark.Height = *Height

	runner.AddBenchmark(benchmark)

	runner.Run()
}
