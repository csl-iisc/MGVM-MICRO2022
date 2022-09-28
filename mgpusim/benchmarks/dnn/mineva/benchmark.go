// Package mineva implements mineva network
package mineva

import (
	"math"

	"gitlab.com/akita/dnn/dataset/mnist"
	"gitlab.com/akita/dnn/layers"
	"gitlab.com/akita/dnn/training"
	"gitlab.com/akita/dnn/training/optimization"
	simLayers "gitlab.com/akita/mgpusim/benchmarks/dnn/layers"
	"gitlab.com/akita/mgpusim/driver"
)

// Benchmark defines the Mineva network training benchmark.
type Benchmark struct {
	driver  *driver.Driver
	context *driver.Context
	gpus    []int

	network training.Network
	trainer training.Trainer
}

// NewBenchmark creates a new benchmark.
func NewBenchmark(driver *driver.Driver) *Benchmark {
	b := new(Benchmark)

	b.driver = driver
	b.context = b.driver.Init()

	b.network = training.Network{
		Layers: []layers.Layer{
			simLayers.CPUToGPULayer{
				GPUDriver: b.driver,
				GPUCtx:    b.context,
			},
			simLayers.NewFullyConnectedLayer(
				784, 256,
				b.driver, b.context,
				simLayers.NewMatrixOperator(b.driver, b.context),
			),
			simLayers.NewReluLayer(b.driver, b.context),
			simLayers.NewFullyConnectedLayer(
				256, 100,
				b.driver, b.context,
				simLayers.NewMatrixOperator(b.driver, b.context),
			),
			simLayers.NewReluLayer(b.driver, b.context),
			simLayers.NewFullyConnectedLayer(
				100, 100,
				b.driver, b.context,
				simLayers.NewMatrixOperator(b.driver, b.context),
			),
			simLayers.NewReluLayer(b.driver, b.context),
			simLayers.NewFullyConnectedLayer(
				100, 10,
				b.driver, b.context,
				simLayers.NewMatrixOperator(b.driver, b.context),
			),
			simLayers.GPUToCPULayer{
				GPUDriver: b.driver,
				GPUCtx:    b.context,
			},
		},
	}

	b.trainer = training.Trainer{
		DataSource:      mnist.NewTrainingDataSource(),
		Network:         b.network,
		LossFunc:        training.SoftmaxCrossEntropy{},
		OptimizationAlg: optimization.NewAdam(0.001),
		Tester: &training.Tester{
			DataSource: mnist.NewTestDataSource(),
			Network:    b.network,
			BatchSize:  math.MaxInt32,
		},
		Epoch:         1000,
		BatchSize:     128,
		ShowBatchInfo: true,
	}

	b.enableLayerVerification(&b.network)

	return b
}

func (b *Benchmark) enableLayerVerification(network *training.Network) {
	network.Layers[1].(*simLayers.FullyConnectedLayer).EnableVerification()
	network.Layers[2].(*simLayers.ReluLayer).EnableVerification()
	network.Layers[3].(*simLayers.FullyConnectedLayer).EnableVerification()
	network.Layers[4].(*simLayers.ReluLayer).EnableVerification()
	network.Layers[5].(*simLayers.FullyConnectedLayer).EnableVerification()
	network.Layers[6].(*simLayers.ReluLayer).EnableVerification()
}

// SelectGPU selects the GPU to use.
func (b *Benchmark) SelectGPU(gpuIDs []int) {
	if len(gpuIDs) > 1 {
		panic("multi-GPU is not supported by DNN workloads")
	}
}

// Run executes the benchmark.
func (b *Benchmark) Run() {
	for _, l := range b.network.Layers {
		l.Randomize()
	}
	b.trainer.Train()
}

// Verify runs the benchmark on the CPU and checks the result.
func (b *Benchmark) Verify() {
	panic("not implemented")
}

// SetUnifiedMemory asks the benchmark to use unified memory.
func (b *Benchmark) SetUnifiedMemory() {
	panic("unified memory is not supported by dnn workloads")
}
