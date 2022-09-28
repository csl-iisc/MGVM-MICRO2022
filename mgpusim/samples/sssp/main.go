package main

import (
	"flag"

	"gitlab.com/akita/mgpusim/benchmarks/pannotia/sssp"
	"gitlab.com/akita/mgpusim/samples/runner"
)

var NumNodes = flag.Int("numNodes", 2097152, "The number of rows in the input matrix.")
var NumItems = flag.Int("numItems", 4194304, "The number of rows in the input matrix.")

func main() {
	flag.Parse()

	runner := new(runner.Runner).ParseFlag().Init()

	benchmark := sssp.NewBenchmark(runner.GPUDriver)
	benchmark.NumNodes = *NumNodes
	benchmark.NumItems = *NumItems

	runner.AddBenchmark(benchmark)

	runner.Run()
}
