package main

import (
	"image"
	"image/color"
	"image/gif"
	"math/rand"
	"os"
)

type EnvCell struct {
	food   uint8
	growth uint8
	agent  *Agent
}

type Environment struct {
	cells [envSize][envSize]EnvCell

	// these 2 fields are for the purpose of recording and generating GIF in the end.
	paletted []color.Color
	record   gif.GIF
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
	env.record.Image[iter] = image.NewPaletted(image.Rect(0, 0, envSize, envSize), env.paletted)
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
			env.record.Image[iter].Set(i, j, env.paletted[env.cells[i][j].agent.appearance])

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
