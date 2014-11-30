package main

import (
	//"fmt"
	"github.com/taylorchu/rbm"
	"math/rand"
)

const (
	Move = 3
	Eat  = 1 << (iota + 2)
	Attack
	Mate
)

type Brain interface {
	// input as described above; output consists of 5 effective bits - 2 for move, 1 for eat, 1 for attack, 1 for mate. (an agent can only attack/mate with the agent in front of it.)
	react([inputLen]uint8) uint8

	// a "train" method shall be added as training is not real-time (separated from react).

	// 2 brains reproduce an offspring (new brain). the second output is the offspring's appearance.
	reproduce(Brain) (Brain, uint8)
}

// NoBrain is a faked stuff for testing purpose.
type NoBrain struct{}

func (nb *NoBrain) react(input [inputLen]uint8) uint8 {
	return Eat | uint8(rand.Intn(4)) | Mate | Attack
}

func (nb *NoBrain) reproduce(mate Brain) (Brain, uint8) {
	return new(NoBrain), 0
}

type RBMBrain struct {
	m       *rbm.RBM
	history [][]uint8
}

const rbmSize = 8*inputLen + 5
const historySize = 10

func NewRBMBrain() *RBMBrain {
	return &RBMBrain{
		m: rbm.New(rbmSize),
	}
}

func (b *RBMBrain) reproduce(mate Brain) (Brain, uint8) {
	return b, 0
}

func expandBits(bs []uint8) (bits []uint8) {
	for _, b := range bs {
		var i uint8
		for i = 7; i < 8; i-- {
			//fmt.Println(b, i, b&(1<<i) != 0)
			if b&(1<<i) != 0 {
				bits = append(bits, 1)
			} else {
				bits = append(bits, 0)
			}
		}
	}
	return
}

func compress(bits []uint8) (b uint8) {
	for i := 0; i < len(bits); i++ {
		if bits[i] == 1 {
			b |= 1 << uint8(len(bits)-i-1)
		}
	}
	return
}

func (b *RBMBrain) react(input [inputLen]uint8) (output uint8) {
	rawInput := make([]uint8, rbmSize)
	copy(rawInput, expandBits(input[:]))
	rawOutput := b.m.Reconstruct(rawInput, 3)
	//b.history = append(b.history, rawOutput)
	//if len(b.history) > historySize {
	//	b.history = b.history[1:]
	//}
	return compress(rawOutput[rbmSize-5 : rbmSize])
}
