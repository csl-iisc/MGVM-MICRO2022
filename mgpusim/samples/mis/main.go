package main

import (
	"flag"

	"gitlab.com/akita/mgpusim/benchmarks/pannotia/mis"
	"gitlab.com/akita/mgpusim/samples/runner"
)

var NumNodes = flag.Int("numNodes", 128, "The number of rows in the input matrix.")
var NumItems = flag.Int("numItems", 64, "The number of rows in the input matrix.")

func main() {
	flag.Parse()

	runner := new(runner.Runner).ParseFlag().Init()

	benchmark := mis.NewBenchmark(runner.GPUDriver)
	benchmark.NumNodes = *NumNodes
	benchmark.NumEdges = *NumItems

	runner.AddBenchmark(benchmark)

	runner.Run()
}
