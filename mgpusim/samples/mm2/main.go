package main

import (
	"flag"

	_ "net/http/pprof"

	"gitlab.com/akita/mgpusim/benchmarks/polybench/mm2"
	"gitlab.com/akita/mgpusim/samples/runner"
)

var niFlag = flag.Uint("ni", 1024, "The height of the first matrix.")
var njFlag = flag.Uint("nj", 1024, "The height of the first matrix.")
var nkFlag = flag.Uint("nk", 1024, "The height of the first matrix.")
var nlFlag = flag.Uint("nl", 1024, "The height of the first matrix.")

func main() {
	flag.Parse()

	runner := new(runner.Runner).ParseFlag().Init()

	benchmark := mm2.NewBenchmark(runner.GPUDriver)
	benchmark.NI = int(*niFlag)
	benchmark.NJ = int(*njFlag)
	benchmark.NK = int(*nkFlag)
	benchmark.NL = int(*nlFlag)

	runner.AddBenchmark(benchmark)

	runner.Run()
}
