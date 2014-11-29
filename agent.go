package main

import (
	"math/rand"
)

type Direction uint8

const (
	Up Direction = iota
	Left
	Down
	Right
)

// (x, y) coordinates corresponding to dir (direction) 0~3 are: (0, -1), (-1, 0), (0, 1), (1, 0)
// record of trajectory will be added.
type Agent struct {
	energy, health, appearance uint8
	dir                        Direction
	brain                      Brain
}

func (ag *Agent) do(i int, j int, env *Environment) {
	cell := &env.cells[i][j]
	output := cell.agent.brain.react(ag.observe(i, j, env))

	if output&Eat != 0 {
		cell.agent.eat(i, j, env)
	}
	if output&Attack != 0 {
		cell.agent.attack(i, j, env)
	}
	if false && output&Mate != 0 {
		cell.agent.mate(i, j, env)
	}
	cell.agent.move(output&Move, i, j, env)
}

func (ag *Agent) observe(i int, j int, env *Environment) (input [inputLen]uint8) {
	cell := &env.cells[i][j]
	input[0] = cell.food
	frontX, frontY := locPlusDir(i, j, cell.agent.dir)
	frontCell := env.cells[frontX][frontY]
	input[1] = frontCell.food
	input[2] = cell.agent.energy
	input[3] = cell.agent.health
	if frontCell.agent != nil {
		input[4] = frontCell.agent.energy
		input[5] = frontCell.agent.health
	}
	event := 0
observe:
	for dist := 1; dist <= sightRadius; dist++ {
		dir := Down //Down, Right, Up, Left
		x := i - dist
		y := j - dist
	radius:
		for {
			xx := (x + envSize) % envSize
			yy := (y + envSize) % envSize
			cell2 := env.cells[xx][yy]
			if cell2.agent != nil {
				input[6+4*event] = cell2.agent.appearance
				input[6+4*event+1] = (uint8(dir) + 3) % 4
				input[6+4*event+2] = cell2.agent.energy
				input[6+4*event+3] = cell2.agent.health
				event++
				if event == eventSize {
					break observe
				}
			}

			switch dir {
			case Down:
				if y < j+dist {
					y++
					if y == j+dist {
						dir = Right
					}
				}
			case Right:
				if x < i+dist {
					x++
					if x == i+dist {
						dir = Up
					}
				}
			case Up:
				if y > j-dist {
					y--
					if y == j-dist {
						dir = Left
					}
				}
			case Left:
				if x > i-dist {
					x--
					if x == i-dist {
						break radius
					}
				}
			}

		}
	}
	return
}

func (ag *Agent) eat(x int, y int, env *Environment) {
	ag.energy = add(ag.energy, env.cells[x][y].food)
	env.cells[x][y].food = 0
	if ag.energy > 192 && ag.health < 255 {
		ag.health++ // agent recovers health when its energy is high.
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
	agent := env.cells[x2][y2].agent

	// cannot attack agent with energy higher than yourself.
	if agent == nil || agent.energy > ag.energy {
		// maybe there should be some punishment for attacking nothing.
		return
	}

	// harm of an attack is your energy / 16
	harm := ag.energy >> 4
	agent.health = sub(agent.health, harm)
}

func (ag *Agent) mate(x int, y int, env *Environment) {
	x2, y2 := locPlusDir(x, y, ag.dir)
	agent := env.cells[x2][y2].agent

	// in order to mate successfully, the agent must have energy, and be in the same direction with the agent in front of it.
	// current hard limit of total number of agents shall be replaced with some more clever manner of population control.
	if agent == nil || agent.dir != ag.dir {
		return
	}
	ag.energy = sub(ag.energy, costOfMate)
	for {
		// offspring is created at some random (empty) location, with random direction.
		x3 := rand.Intn(envSize)
		y3 := rand.Intn(envSize)
		if env.cells[x3][y3].agent != nil {
			break
		}
		agent := new(Agent)
		// get brain from the mate
		agent.brain, _ = ag.brain.reproduce(env.cells[x2][y2].agent.brain)
		agent.energy = initEnergy
		agent.health = initHealth

		// appearance shall be set by the second output of "reproduce".
		agent.appearance = ag.appearance

		agent.dir = Direction(rand.Intn(4))
		env.cells[x3][y3].agent = agent
		break
	}
}

func (ag *Agent) move(op uint8, x int, y int, env *Environment) {
	switch op {
	case 0: // no op
	case 1: // turn left
		ag.dir = (ag.dir + 1) & 3
	case 2: // turn right
		ag.dir = (ag.dir + 3) & 3
	case 3:
		x2, y2 := locPlusDir(x, y, ag.dir)
		if env.cells[x2][y2].agent != nil { // cannot move forward if front cell is blocked.
			return
		}
		ag.energy = sub(ag.energy, costOfMove)
		env.cells[x2][y2].agent = ag
		env.cells[x][y].agent = nil
	}
}
