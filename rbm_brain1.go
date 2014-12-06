package main

import "github.com/taylorchu/rbm"

type RBMBrain struct {
	*rbm.RBM
	Mem [][]float64
}

func NewRBMBrain() *RBMBrain {
	return &RBMBrain{
		RBM: rbm.New(InputSize + OutputSize),
		Mem: make([][]float64, TrainScope*2),
	}
}

func (b *RBMBrain) Reward(score float64) {
	if score < 1 || len(b.Mem[0]) != InputSize+OutputSize {
		return
	}
	b.Train(b.Mem[:1], 0.1, int(score))
}

func (b *RBMBrain) React(input *BrainInput) *BrainOutput {
	raw := make([]float64, InputSize+OutputSize)
	input.Encode(raw)
	rawOutput := b.Reconstruct(raw, 10)
	b.Mem = append(b.Mem[1:], rawOutput)
	output := new(BrainOutput)
	output.Decode(rawOutput[InputSize:])
	return output
}

func (b *RBMBrain) Default() {
	b.Train(createTrainingData(), 0.1, 1000)
}
