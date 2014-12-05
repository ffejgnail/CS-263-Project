package main

import "github.com/taylorchu/rbm"

//   | 0 |
// 1 |   | 3
//   | 2 |

type Nearby struct {
	HasAnimat       bool
	OtherFace       Face
	OtherTargetFace Face
	MoreHealth      bool
	MoreFood        bool
}

type BrainInput struct {
	MyFace       Face
	MyTargetFace Face
	Nearby       [4]Nearby
}

const (
	InputSize  = 21 // 1 + 5 * 4
	OutputSize = 3
)

func (bi *BrainInput) Encode(b []uint8) {
	if bi.MyFace == bi.MyTargetFace {
		b[0] = 1
	}
	for i := 0; i < 4; i++ {
		off := 1 + 5*i
		if bi.Nearby[i].OtherFace == bi.MyFace {
			b[off] = 1
		}
		if bi.Nearby[i].OtherTargetFace == bi.MyFace {
			b[off+1] = 1
		}
		if bi.Nearby[i].MoreHealth {
			b[off+2] = 1
		}
		if bi.Nearby[i].MoreFood {
			b[off+3] = 1
		}
		if bi.Nearby[i].HasAnimat {
			b[off+4] = 1
		}
	}
}

func encodeFace(f Face) (uint8, uint8, uint8) {
	return uint8(f) & 7 >> 2, uint8(f) & 3 >> 1, uint8(f) & 1
}

type Move uint8

const (
	Stay Move = iota
	TurnLeft
	TurnRight
	Forward
)

type BrainOutput struct {
	Attack bool
	Move   Move
}

type Brain interface {
	React(*BrainInput) *BrainOutput
	Reward(score float64)
}

type RBMBrain struct {
	*rbm.RBM
	Mem [][]uint8
}

func NewRBMBrain() *RBMBrain {
	return &RBMBrain{
		RBM: rbm.New(InputSize + OutputSize),
		Mem: make([][]uint8, TrainScope*2),
	}
}

func (b *RBMBrain) Reward(score float64) {
	if score < 1 || len(b.Mem[0]) != InputSize+OutputSize {
		return
	}
	b.Train(b.Mem[:1], 0.1, int(score))
}

func (b *RBMBrain) React(input *BrainInput) *BrainOutput {
	raw := make([]uint8, InputSize+OutputSize)
	input.Encode(raw)
	rawOutput := b.Reconstruct(raw, 3)
	b.Mem = append(b.Mem[1:], rawOutput)
	output := new(BrainOutput)
	output.Attack = rawOutput[len(rawOutput)-3] == 1

	if rawOutput[len(rawOutput)-2] == 1 {
		output.Move |= 2
	}
	if rawOutput[len(rawOutput)-1] == 1 {
		output.Move |= 1
	}
	return output
}
