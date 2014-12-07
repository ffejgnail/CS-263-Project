package main

import "github.com/taylorchu/rbm"

type RBMBrain struct {
	*rbm.RBM
}

func NewRBMBrain() *RBMBrain {
	return &RBMBrain{
		RBM: rbm.New(InputSize+OutputSize, (InputSize+OutputSize)/2),
	}
}

func (b *RBMBrain) Reward(data [][]float64, score int) {
	b.Train(data, score)
}

func (b *RBMBrain) React(input *BrainInput) (*BrainOutput, []float64) {
	raw := make([]float64, InputSize)
	input.Encode(raw)
	rawOutput := b.Reconstruct(raw, 10)
	output := new(BrainOutput)
	output.Decode(rawOutput[InputSize:])
	return output, rawOutput
}

func (b *RBMBrain) Default() {
	b.Train(createTrainingData(), 1000)
}

func (b *RBMBrain) Dump(filename string) (err error) {
	return b.RBM.WriteFile(filename)
}

func (b *RBMBrain) Load(filename string) (err error) {
	r, err := rbm.FromFile(filename)
	if err != nil {
		return
	}
	b.RBM = r
	return
}
