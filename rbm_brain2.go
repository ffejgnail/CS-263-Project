package main

import (
	"math/rand"

	"github.com/r9y9/nnet/rbm"
)

type RBMBrain2 struct {
	*rbm.RBM
	Mem [][]float64
}

func NewRBMBrain2() *RBMBrain2 {
	return &RBMBrain2{
		RBM: rbm.New(InputSize+OutputSize, 10),
		Mem: make([][]float64, TrainScope*2),
	}
}

func (b *RBMBrain2) Reward(score float64) {
	return
	if score < 1 || len(b.Mem[0]) != InputSize+OutputSize {
		return
	}
	option := rbm.TrainingOption{
		LearningRate:         0.1,
		Epoches:              int(score),
		OrderOfGibbsSampling: 1,
		MiniBatchSize:        20,
		L2Regularization:     true,
		RegularizationRate:   1.0e-10,
		Monitoring:           false,
	}
	b.Train(b.Mem[:1], option)
}

func (b *RBMBrain2) React(input *BrainInput) *BrainOutput {
	raw := make([]float64, InputSize+OutputSize)
	input.Encode(raw)
	for i := 0; i < OutputSize; i++ {
		raw[InputSize+i] = rand.Float64()
	}
	rawOutput, _ := b.Reconstruct(raw, 10)
	b.Mem = append(b.Mem[1:], rawOutput)
	output := new(BrainOutput)
	output.Decode(rawOutput[InputSize:])
	return output
}

func (b *RBMBrain2) Default() {
	option := rbm.TrainingOption{
		LearningRate:         0.1,
		Epoches:              1000,
		OrderOfGibbsSampling: 1,
		MiniBatchSize:        20,
		L2Regularization:     true,
		RegularizationRate:   1.0e-10,
		Monitoring:           false,
	}
	b.Train(createTrainingData(), option)
}
