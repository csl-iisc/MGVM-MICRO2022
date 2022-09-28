package main

import (
	"flag"

	_ "net/http/pprof"

	"gitlab.com/akita/mgpusim/benchmarks/polybench/mm3"
	"gitlab.com/akita/mgpusim/samples/runner"
)

var niFlag = flag.Uint("ni", 256, "The height of the first matrix.")
var njFlag = flag.Uint("nj", 256, "The height of the first matrix.")
var nkFlag = flag.Uint("nk", 256, "The height of the first matrix.")
var nlFlag = flag.Uint("nl", 256, "The height of the first matrix.")
var nmFlag = flag.Uint("nm", 256, "The height of the first matrix.")

func main() {
	flag.Parse()

	runner := new(runner.Runner).ParseFlag().Init()

	benchmark := mm3.NewBenchmark(runner.GPUDriver)
	benchmark.NI = int(*niFlag)
	benchmark.NJ = int(*njFlag)
	benchmark.NK = int(*nkFlag)
	benchmark.NL = int(*nlFlag)
	benchmark.NM = int(*nmFlag)

	runner.AddBenchmark(benchmark)

	runner.Run()
}
