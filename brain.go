package main

import (
	"fmt"
	"math/rand"
)

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
	InputSize  = 12 // 4 * 3
	OutputSize = 4
)

func (bi *BrainInput) Encode(b []float64) {
	//if bi.MyFace == bi.MyTargetFace {
	//	b[0] = 1
	//}
	for i := 0; i < 4; i++ {
		off := 3 * i
		//if bi.Nearby[i].OtherFace == bi.MyFace {
		//	b[off] = 1
		//}
		//if bi.Nearby[i].OtherTargetFace == bi.MyFace {
		//	b[off+1] = 1
		//}
		if bi.Nearby[i].MoreHealth {
			b[off] = 1
		}
		if bi.Nearby[i].MoreFood {
			b[off+1] = 1
		}
		if bi.Nearby[i].HasAnimat {
			b[off+2] = 1
		}
	}
}

type Move uint8

const (
	MoveUp Move = iota
	MoveLeft
	MoveDown
	MoveRight
	Stay
)

type BrainOutput struct {
	Attack bool
	Move   Move
}

func (bo *BrainOutput) Decode(b []float64) {
	bo.Attack = b[0] == 1
	if b[1] == 1 {
		bo.Move |= 4
	}
	if b[2] == 1 {
		bo.Move |= 2
	}
	if b[3] == 1 {
		bo.Move |= 1
	}
}

func (bo *BrainOutput) Encode(b []float64) {
	if bo.Attack {
		b[0] = 1
	}
	if bo.Move&4 != 0 {
		b[1] = 1
	}
	if bo.Move&2 != 0 {
		b[2] = 1
	}
	if bo.Move&1 != 0 {
		b[3] = 1
	}
}

type Brain interface {
	React(*BrainInput) (*BrainOutput, []float64)
	Reward(data [][]float64, score float64)
	Dump(filename string) error
	Load(filename string) error
}

func expected(input *BrainInput) (output *BrainOutput) {
	output = new(BrainOutput)
	output.Move = Stay
	for i := 0; i < 4; i++ {
		if input.Nearby[i].MoreFood {
			output.Move = Move(i)
			return
		}
	}
	return
}

func createTrainingData() [][]float64 {
	data := make([][]float64, PreTrain)
	for i := 0; i < PreTrain; i++ {
		input := new(BrainInput)
		input.MyFace = Face(rand.Intn(8))
		input.MyTargetFace = Face(rand.Intn(8))
		for i := 0; i < 4; i++ {
			nearby := &input.Nearby[i]
			nearby.MoreFood = rand.Intn(2) == 1
			if rand.Intn(2) == 1 {
				nearby.HasAnimat = true
				nearby.MoreHealth = rand.Intn(2) == 1
				nearby.OtherFace = Face(rand.Intn(8))
				nearby.OtherTargetFace = Face(rand.Intn(8))
			}
		}
		output := expected(input)
		data[i] = make([]float64, InputSize+OutputSize)
		input.Encode(data[i])
		output.Encode(data[i][InputSize:])
		fmt.Println(data[i])
	}
	return data
}
