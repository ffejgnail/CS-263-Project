package main

import (
	"math"
	"math/rand"
)

type Direction uint8

const (
	Up Direction = iota
	Down
	Left
	Right
)

// (x, y) coordinates corresponding to dir (direction) 0~3 are: (0, -1), (-1, 0), (0, 1), (1, 0)
// record of trajectory will be added.
type Agent struct {
	energy, health, appearance uint8
	dir                        Direction
	brain                      Brain
}

func (ag *Agent) eat(x int, y int, env *Environment) {
	ag.energy = (ag.energy + env.cells[x][y].food) % math.MaxUint8 // when food in cell is low (even empty), eating should be punished.
	env.cells[x][y].food = 0
	if ag.energy > 192 && ag.health < 240 {
		ag.health += 16 // agent recovers health when its energy is high.
	}
}

// calculate the coordinate of the cell in front of a given cell and direction.
// the world is wrapped around.
func locPlusDir(x int, y int, dir Direction) (int, int) {
	switch dir {
	case Up:
		return x, (y + envSize - 1) % envSize
	case Left:
		return (x + envSize - 1) % envSize, y
	case Down:
		return x, (y + 1) % envSize
	case Right:
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
	} else {
		env.cells[x2][y2].agent.health -= harm
	}
}

func (ag *Agent) mate(x int, y int, env *Environment) {
	x2, y2 := locPlusDir(x, y, ag.dir)

	// in order to mate successfully, the agent must have energy, and be in the same direction with the agent in front of it.
	// current hard limit of total number of agents shall be replaced with some more clever manner of population control.
	if currentAgentNum == initAgentNum || ag.energy < costOfMate || env.cells[x2][y2].agent == nil || env.cells[x2][y2].agent.dir != ag.dir {
		return
	}
	ag.energy -= costOfMate
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

		env.cells[x3][y3].agent.dir = Direction(rand.Intn(4))
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
