package main

import (
	"flag"

	"gitlab.com/akita/mgpusim/benchmarks/amdappsdk/bitonicsort"
	"gitlab.com/akita/mgpusim/samples/runner"
)

var length = flag.Int("length", 1024, "The length of array to sort.")
var orderAscending = flag.Bool("order-asc", true, "Sorting in ascending order.")

func main() {
	flag.Parse()

	runner := new(runner.Runner).ParseFlag().Init()

	benchmark := bitonicsort.NewBenchmark(runner.GPUDriver)
	benchmark.Length = *length
	benchmark.OrderAscending = *orderAscending

	runner.AddBenchmark(benchmark)

	runner.Run()
}
