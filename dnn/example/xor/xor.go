package main

import (
	"math/rand"

	"gitlab.com/akita/dnn/tensor"
	"gitlab.com/akita/dnn/training/optimization"

	"gitlab.com/akita/dnn/layers"
	"gitlab.com/akita/dnn/training"
)

func main() {
	rand.Seed(1)
	to := tensor.CPUOperator{}

	network := training.Network{
		Layers: []layers.Layer{
			layers.NewFullyConnectedLayer(to, 2, 4),
			layers.NewReluLayer(to),
			layers.NewFullyConnectedLayer(to, 4, 2),
		},
	}
	trainer := training.Trainer{
		DataSource:      NewDataSource(to),
		Network:         network,
		LossFunc:        training.NewSoftmaxCrossEntropy(to),
		OptimizationAlg: optimization.NewSGD(to, 0.03),
		Epoch:           50,
		BatchSize:       4,
		ShowBatchInfo:   true,
	}

	for _, l := range network.Layers {
		l.Randomize()
	}

	trainer.Train()
}
