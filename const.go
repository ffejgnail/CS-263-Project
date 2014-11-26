package main

// size of the world. assume the world to be square instead of rectangular for simplicity.
const envSize int = 128

// max number of other agents that an agent can sense at a time.
const eyeRange int = 8

// eyesight of an agent is 2*eyeSight+1 by 2*eyeSight+1
const eyeSight int = 4

// input size to the RBM is 8 * inputLen bits.
// current design of input consists of:
//	1. amount of food in current cell
//	2. amount of food in front cell
//	3. self energy
//	4. self health
//	5. energy of agent in front cell (if applicable)
//	6. health of agent in front cell (if applicable)
//	7. events; for each event:
//		a. appearance
//		b. relative location
//		c. attack object's appearance
//		d. mate object's appearance
// agents' trajectory may be added as input to RBM in the future.
const inputLen int = 6 + 4*eyeRange

const (
	initAgentNum    uint8 = 32
	numOfIterations int   = 2000
	grassGrowFreq   int   = 8
	initEnergy      uint8 = 16
	initHealth      uint8 = 64
	costOfMate      uint8 = 5
	costOfMove      uint8 = 1
)
