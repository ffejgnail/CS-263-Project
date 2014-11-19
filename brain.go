package main

import (
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
	return uint8(rand.Intn(255))
}

func (nb *NoBrain) reproduce(mate Brain) (Brain, uint8) {
	return new(NoBrain), 0
}
