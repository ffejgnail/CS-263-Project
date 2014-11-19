package main

import (
	"fmt"
	"math/rand"
	"os"
	"image"
	"image/color"
	"image/gif"
)

// size of the world. assume the world to be square instead of rectangular for simplicity.
const envSize int = 128

// max number of other agents that an agent can sense at a time.
const eyeRange int = 8

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

// this number is currently used for appearance generation purpose. it is expected to be remove once a real Brain is implemented.
var currentAgentNum uint8

const initAgentNum uint8 = 16
const numOfIterations int = 16
const initEnergy uint8 = 16
const initHealth uint8 = 64
const costOfMate uint8 = 128
const costOfMove uint8 = 128
const costOfMate uint8 = 128

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

// (x, y) coordinates corresponding to dir (direction) 0~3 are: (0, -1), (-1, 0), (0, 1), (1, 0)
// record of trajectory will be added.
type Agent struct {
	energy, health, appearance, dir uint8
	brain                           Brain
}

func (ag *Agent) eat(x int, y int, env *Environment) {
	if ag.energy > 255-env.cells[x][y].food {
		ag.energy = 255
		env.cells[x][y].food = 0		// all food in the cell is eaten; may waste a lot.
	}
	ag.energy = ag.energy + env.cells[x][y].food	// when food in cell is low (even empty), eating should be punished.
	env.cells[x][y].food = 0
	if ag.energy > 192 && ag.health < 240 {
		ag.health = ag.health + 16		// agent recovers health when its energy is high.
	}
}

// calculate the coordinate of the cell in front of a given cell and direction.
// the world is wrapped around.
func locPlusDir(x int, y int, dir uint8) (int, int) {
	switch dir {
	case 0:
		return x, (y + envSize - 1) % envSize
	case 1:
		return (x + envSize - 1) % envSize, y
	case 2:
		return x, (y + 1) % envSize
	case 3:
		return (x + 1) % envSize, y
	default:
		return x, y
	}
}

func (ag *Agent) attack(x int, y int, env *Environment) {
	x2, y2 := locPlusDir(x, y, ag.dir)

	// cannot attack agent with energy higher than yourself.
	if env.cells[x2][y2].agent == nil || env.cells[x2][y2].agent.energy > ag.energy {
		// maybe there should be some punishment for attacking nothing.
		return
	}

	// harm of an attack is your energy / 16
	harm := ag.energy >> 4

	// if harm is too big, avoid "under-flow"
	if harm > env.cells[x2][y2].agent.health {
		env.cells[x2][y2].agent.health = 0
		harm = 0
	}
	env.cells[x2][y2].agent.health = env.cells[x2][y2].agent.health - harm
}

func (ag *Agent) mate(x int, y int, env *Environment) {
	x2, y2 := locPlusDir(x, y, ag.dir)

	// in order to mate successfully, the agent must have energy, and be in the same direction with the agent in front of it.
	// current hard limit of total number of agents shall be replaced with some more clever manner of population control.
	if currentAgentNum == initAgentNum || ag.energy < costOfMate || env.cells[x2][y2].agent == nil || env.cells[x2][y2].agent.dir != ag.dir {
		return
	}
	ag.energy = ag.energy - costOfMate
	for {
		// offspring is created at some random (empty) location, with random direction.
		x3 := rand.Intn(envSize)
		y3 := rand.Intn(envSize)
		if env.cells[x3][y3].agent != nil {
			continue
		}
		env.cells[x3][y3].agent = new(Agent)
		env.cells[x3][y3].agent.brain, _ = ag.brain.reproduce(env.cells[x2][y2].agent.brain)
		env.cells[x3][y3].agent.energy = initEnergy
		env.cells[x3][y3].agent.health = initHealth

		// appearance shall be set by the second output of "reproduce".
		env.cells[x3][y3].agent.appearance = currentAgentNum
		currentAgentNum++

		env.cells[x3][y3].agent.dir = uint8(rand.Intn(4))
		break
	}
}

func (ag *Agent) move(op uint8, x int, y int, env *Environment) {
	if ag.energy < costOfMove || op != 3 {
		if op == 1 { // turn left
			ag.dir = (ag.dir + 1) & 3
		}
		if op == 2 { // turn right
			ag.dir = (ag.dir + 3) & 3
		}
		return
	}
	x2, y2 := locPlusDir(x, y, ag.dir)
	if env.cells[x2][y2].agent != nil { // cannot move forward if front cell is blocked.
		return
	}
	ag.energy = ag.energy - costOfMove
	env.cells[x2][y2].agent = ag
	env.cells[x][y].agent = nil
}

type EnvCell struct {
	food   uint8
	growth uint8
	agent  *Agent
}

type Environment struct {
	cells [envSize][envSize]EnvCell

	// these 2 fields are for the purpose of recording and generating GIF in the end.
	paletted []color.Color
	record gif.GIF
}

func (env *Environment) print() { // generating GIF
	outputImg, _ := os.Create("new.gif")
	defer outputImg.Close()
	gif.EncodeAll(outputImg, &env.record)
}

func (env *Environment) setup() {
	rand.Seed(42)
	for i := 0; i < envSize; i++ {
		for j := 0; j < envSize; j++ {
			env.cells[i][j].growth = 1
		}
		env.cells[i][i].growth = 10
	}

	for i := uint8(0); i < initAgentNum; i++ {
		x := rand.Intn(envSize)
		y := rand.Intn(envSize)
		env.cells[x][y].agent = new(Agent)
		env.cells[x][y].agent.brain = new(NoBrain)
		env.cells[x][y].agent.energy = initEnergy
		env.cells[x][y].agent.health = initHealth
		env.cells[x][y].agent.appearance = i
		env.cells[x][y].agent.dir = uint8(rand.Intn(4))
	}
	currentAgentNum = initAgentNum

	env.paletted = []color.Color{
		color.RGBA{255, 255, 255, 255},
		color.RGBA{192, 192, 192, 255},
		color.RGBA{128, 128, 128, 255},
		color.RGBA{0, 0, 0, 255},
		color.RGBA{255, 0, 0, 255},
		color.RGBA{128, 0, 0, 255},
		color.RGBA{255, 255, 0, 255},
		color.RGBA{128, 128, 0, 255},
		color.RGBA{0, 255, 0, 255},
		color.RGBA{0, 128, 0, 255},
		color.RGBA{0, 255, 255, 255},
		color.RGBA{0, 128, 128, 255},
		color.RGBA{0, 0, 255, 255},
		color.RGBA{0, 0, 128, 255},
		color.RGBA{255, 0, 255, 255},
		color.RGBA{128, 0, 128, 255}}
	env.record.Image = make([]*image.Paletted, numOfIterations)
	env.record.Delay = make([]int, numOfIterations)
	for i := 0; i < numOfIterations; i++ {
		env.record.Delay[i] = 50
	}
	env.record.LoopCount = 256
}

func (env *Environment) run(iter int) {
	env.record.Image[iter] = image.NewPaletted(image.Rect(0,0,envSize,envSize), env.paletted)
	for i := 0; i < envSize; i++ {
		for j := 0; j < envSize; j++ {
			if env.cells[i][j].food < 255-env.cells[i][j].growth {
				env.cells[i][j].food = env.cells[i][j].food + env.cells[i][j].growth
			}
			if env.cells[i][j].agent == nil || env.cells[i][j].agent.health == 0 { // agent with no health die
				env.cells[i][j].agent = nil
				continue
			}

			// recording. (currently "appearance" must < 16, o.w. array out of bound.)
			// TODO: determine color by agent's state (energy & health).
			//	i.e.: bad state (low health & energy) as red, good state (high health & energy) as green.
			env.record.Image[iter].Set(i,j,env.paletted[env.cells[i][j].agent.appearance])

			var input [inputLen]uint8 // TODO: write code to collect input for Brain.

			output := env.cells[i][j].agent.brain.react(input)

			if output&4 == 4 {
				env.cells[i][j].agent.eat(i, j, env)
			}
			if output&8 == 8 {
				env.cells[i][j].agent.attack(i, j, env)
			}
			if output&16 == 16 {
				env.cells[i][j].agent.mate(i, j, env)
			}
			env.cells[i][j].agent.move(output&3, i, j, env)
		}
	}
}

func main() {
	var env Environment
	env.setup()
	for i := 0; i < numOfIterations; i++ {
		env.run(i)
	}
	env.print()
	fmt.Println("end")
}
