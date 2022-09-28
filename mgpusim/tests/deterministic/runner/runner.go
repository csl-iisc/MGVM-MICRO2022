package runner

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	// Enable debugging
	_ "net/http/pprof"
	"strconv"
	"strings"
	"sync"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/mgpusim/benchmarks"
	"gitlab.com/akita/mgpusim/driver"
	"gitlab.com/akita/mgpusim/platform"
	"gitlab.com/akita/util/tracing"
)

var verifyFlag = flag.Bool("verify", false, "Verify the emulation result.")
var cacheLatencyReportFlag = flag.Bool("report-cache-latency", false, "Report the average cache latency.")
var gpuFlag = flag.String("gpus", "1",
	"The GPUs to use, use a format like 1,2,3,4")

type cacheLatencyTracer struct {
	tracer *tracing.AverageTimeTracer
	cache  akita.Component
}

// Runner is a class that helps running the benchmarks in the official samples.
type Runner struct {
	Engine                  akita.Engine
	GPUDriver               *driver.Driver
	KernelTimeCounter       *tracing.BusyTimeTracer
	PerGPUKernelTimeCounter []*tracing.BusyTimeTracer
	CacheLatencyTracers     []cacheLatencyTracer
	Benchmarks              []benchmarks.Benchmark
	Timing                  bool
	Verify                  bool
	Parallel                bool
	ReportCacheLatency      bool

	GPUIDs []int
}

// ParseFlag applies the runner flag to runner object
func (r *Runner) ParseFlag() *Runner {
	if *verifyFlag {
		r.Verify = true
	}

	r.Timing = true

	if *cacheLatencyReportFlag {
		r.ReportCacheLatency = true
	}

	r.parseGPUFlag()
	return r
}

func (r *Runner) startProfilingServer() {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	fmt.Println("Profiling server running on:", listener.Addr().(*net.TCPAddr).Port)

	panic(http.Serve(listener, nil))
}

// Init initializes the platform simulate
func (r *Runner) Init() *Runner {
	go r.startProfilingServer()

	builder := platform.MakeMCMGPUBuilder()
	r.Engine, r.GPUDriver = builder.Build()

	r.addKernelTimeTracer()
	r.addCacheLatencyTracer()

	return r
}

func (r *Runner) addKernelTimeTracer() {
	r.KernelTimeCounter = tracing.NewBusyTimeTracer(
		func(task tracing.Task) bool {
			return task.What == "*driver.LaunchKernelCommand"
		})
	tracing.CollectTrace(r.GPUDriver, r.KernelTimeCounter)

	for _, gpu := range r.GPUDriver.GPUs {
		gpuKernelTimeCountner := tracing.NewBusyTimeTracer(
			func(task tracing.Task) bool {
				return task.What == "Launch Kernel"
			})
		r.PerGPUKernelTimeCounter = append(
			r.PerGPUKernelTimeCounter, gpuKernelTimeCountner)
		tracing.CollectTrace(gpu.CommandProcessor, gpuKernelTimeCountner)
	}
}

func (r *Runner) addCacheLatencyTracer() {
	if !r.ReportCacheLatency {
		return
	}

	for _, gpu := range r.GPUDriver.GPUs {
		for _, cache := range gpu.L2Caches {
			tracer := tracing.NewAverageTimeTracer(
				func(task tracing.Task) bool {
					return task.Kind == "req_in"
				})
			r.CacheLatencyTracers = append(r.CacheLatencyTracers,
				cacheLatencyTracer{tracer: tracer, cache: cache})
			tracing.CollectTrace(cache, tracer)
		}
	}
}

func (r *Runner) parseGPUFlag() {
	r.GPUIDs = make([]int, 0)
	gpuIDTokens := strings.Split(*gpuFlag, ",")
	for _, t := range gpuIDTokens {
		gpuID, err := strconv.Atoi(t)
		if err != nil {
			log.Fatal(err)
		}
		r.GPUIDs = append(r.GPUIDs, gpuID)
	}
}

// AddBenchmark adds an benchmark that the driver runs
func (r *Runner) AddBenchmark(b benchmarks.Benchmark) {
	b.SelectGPU(r.GPUIDs)
	r.Benchmarks = append(r.Benchmarks, b)
}

// AddBenchmarkWithoutSettingGPUsToUse allows for user specified GPUs for
// the benchmark to run.
func (r *Runner) AddBenchmarkWithoutSettingGPUsToUse(b benchmarks.Benchmark) {
	r.Benchmarks = append(r.Benchmarks, b)
}

// Run runs the benchmark on the simulator
func (r *Runner) Run() {
	r.GPUDriver.Run()

	var wg sync.WaitGroup
	for _, b := range r.Benchmarks {
		wg.Add(1)
		go func(b benchmarks.Benchmark, wg *sync.WaitGroup) {
			b.Run()
			if r.Verify {
				b.Verify()
			}
			wg.Done()
		}(b, &wg)
	}
	wg.Wait()

	r.GPUDriver.Terminate()
	r.Engine.Finished()

	fmt.Printf("Kernel time: %.12f\n", r.KernelTimeCounter.BusyTime())
	fmt.Printf("Total time: %.12f\n", r.Engine.CurrentTime())
	for i, c := range r.PerGPUKernelTimeCounter {
		fmt.Printf("GPU %d kernel time: %.12f\n", i+1, c.BusyTime())
	}

	for _, tracer := range r.CacheLatencyTracers {
		fmt.Printf("Cache %s average latency %.12f\n",
			tracer.cache.Name(),
			tracer.tracer.AverageTime(),
		)
	}
}
