package main

import "github.com/r9y9/nnet/rbm"

type RBMBrain2 struct {
	*rbm.RBM
}

func NewRBMBrain2() *RBMBrain2 {
	return &RBMBrain2{
		RBM: rbm.New(InputSize+OutputSize, (InputSize+OutputSize)/2),
	}
}

func (b *RBMBrain2) Reward(data [][]float64, score float64) {
	if score < 1 {
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
	b.Train(data, option)
}

func (b *RBMBrain2) React(input *BrainInput) (*BrainOutput, []float64) {
	raw := make([]float64, InputSize+OutputSize)
	input.Encode(raw)
	for i := 0; i < OutputSize; i++ {
		raw[InputSize+i] = 0.5
	}
	rawOutput, _ := b.Reconstruct(raw, 10)
	output := new(BrainOutput)
	output.Decode(rawOutput[InputSize:])
	return output, rawOutput
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

func (b *RBMBrain2) Load(filename string) (err error) {
	r, err := rbm.Load(filename)
	if err != nil {
		return
	}
	b.RBM = r
	return
}
