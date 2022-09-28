package main

import (
	"flag"

	_ "net/http/pprof"

	"gitlab.com/akita/mgpusim/benchmarks/polybench/syrk"
	"gitlab.com/akita/mgpusim/samples/runner"
)

var niFlag = flag.Uint("ni", 64, "The height of the first matrix.")
var njFlag = flag.Uint("nj", 64, "The height of the first matrix.")

func main() {
	flag.Parse()

	runner := new(runner.Runner).ParseFlag().Init()

	benchmark := syrk.NewBenchmark(runner.GPUDriver)
	benchmark.NI = int(*niFlag)
	benchmark.NJ = int(*njFlag)

	runner.AddBenchmark(benchmark)

	runner.Run()
}
