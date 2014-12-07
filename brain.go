package main

import (
	"fmt"
	"math/rand"
)

//   | 0 |
// 1 |   | 3
//   | 2 |

type Nearby struct {
	HasAnimat   bool
	Color       Color
	TargetColor Color
	MoreHealth  bool
	MoreFood    bool
	Nice        bool
}

type BrainInput struct {
	Color       Color
	TargetColor Color
	Nearby      [4]Nearby
}

const (
	InputSize  = 25 // 1 + 4 * 6
	OutputSize = 4
)

func (bi *BrainInput) Encode(b []float64) {
	if bi.Color == bi.TargetColor {
		b[0] = 1
	}
	for i := 0; i < 4; i++ {
		off := 1 + 6*i
		if bi.Nearby[i].MoreFood {
			b[off] = 1
		}
		if bi.Nearby[i].HasAnimat {
			b[off+1] = 1
			if bi.Nearby[i].Color == bi.Color {
				b[off+2] = 1
			}
			if bi.Nearby[i].TargetColor == bi.Color {
				b[off+3] = 1
			}
			if bi.Nearby[i].MoreHealth {
				b[off+4] = 1
			}
			if bi.Nearby[i].Nice {
				b[off+5] = 1
			}
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
	MoveRand
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
	Reward(data [][]float64, score int)
	Dump(filename string) error
	Load(filename string) error
}

func expected(input *BrainInput) (output *BrainOutput) {
	output = new(BrainOutput)
	output.Move = MoveRand
	for i := 0; i < 4; i++ {
		nearby := &input.Nearby[i]
		if nearby.MoreFood {
			output.Move = Move(i)
			if nearby.HasAnimat && (nearby.Color != input.Color || !nearby.Nice) {
				output.Attack = true
			}
			return
		}
	}
	return
}

func createTrainingData() [][]float64 {
	data := make([][]float64, PreTrain)
	for i := 0; i < PreTrain; i++ {
		input := new(BrainInput)
		input.Color = Color(rand.Intn(3) + 1)
		input.TargetColor = Color(rand.Intn(3) + 1)
		for i := 0; i < 4; i++ {
			nearby := &input.Nearby[i]
			nearby.MoreFood = rand.Intn(2) == 1
			if rand.Intn(2) == 1 {
				nearby.HasAnimat = true
				nearby.MoreHealth = rand.Intn(2) == 1
				nearby.Color = Color(rand.Intn(3) + 1)
				nearby.TargetColor = Color(rand.Intn(3) + 1)
				nearby.Nice = rand.Intn(2) == 1
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
