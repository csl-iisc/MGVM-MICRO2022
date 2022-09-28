package main

import (
	"flag"

	"gitlab.com/akita/mgpusim/benchmarks/polybench/gemver"
	"gitlab.com/akita/mgpusim/samples/runner"
)

var nFlag = flag.Int("n", 512, "Dunno")

func main() {
	flag.Parse()

	runner := new(runner.Runner).ParseFlag().Init()

	benchmark := gemver.NewBenchmark(runner.GPUDriver)
	benchmark.N = *nFlag
	benchmark.Alpha = 1.5 //42123.0
	benchmark.Beta = 1.5  //324.0

	runner.AddBenchmark(benchmark)

	runner.Run()
}
