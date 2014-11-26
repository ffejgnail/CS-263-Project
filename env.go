package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"math/rand"
)

type EnvCell struct {
	food   uint8
	growth uint8
	agent  *Agent
}

type Environment struct {
	cells [envSize][envSize]EnvCell

	record gif.GIF
}

func (env *Environment) WriteTo(w io.Writer) error { // generating GIF
	return gif.EncodeAll(w, &env.record)
}

func (env *Environment) setup() {
	rand.Seed(42)
	for i := 0; i < envSize; i++ {
		for j := 0; j < envSize; j++ {
			env.cells[i][j].growth = uint8(rand.Intn(2)) + 2
			//env.cells[i][j].food = env.cells[i][j].growth * 30
		}
	}

	for i := uint8(0); i < initAgentNum; i++ {
		x := rand.Intn(envSize)
		y := rand.Intn(envSize)
		agent := new(Agent)
		agent.brain = new(NoBrain)
		agent.energy = initEnergy
		agent.health = initHealth
		agent.appearance = i
		agent.dir = Direction(rand.Intn(4))
		env.cells[x][y].agent = agent
	}

	env.record.Image = make([]*image.Paletted, numOfIterations)
	env.record.Delay = make([]int, numOfIterations)
	for i := 0; i < numOfIterations; i++ {
		env.record.Delay[i] = 10
	}
}

var (
	colors = []color.Color{
		color.RGBA{30, 60, 30, 255},
		color.RGBA{60, 120, 60, 255},
		color.RGBA{90, 180, 90, 255},
		color.RGBA{120, 240, 120, 255},
		color.RGBA{240, 240, 120, 255},
		color.RGBA{255, 0, 0, 255},
	}
)

func appearanceColor(appear uint8) color.Color {
	return colors[5]
}

func grassColor(grass uint8) color.Color {
	if grass > 120 {
		return colors[0]
	}
	if grass > 90 {
		return colors[1]
	}
	if grass > 60 {
		return colors[2]
	}
	if grass > 30 {
		return colors[3]
	}
	return colors[4]
}

func (env *Environment) run(iter int) {
	moved := make(map[*Agent]bool)
	for i := 0; i < envSize; i++ {
		for j := 0; j < envSize; j++ {
			cell := &env.cells[i][j]
			if iter%grassGrowFreq == 0 {
				growth := uint8(rand.Intn(int(cell.growth)))
				if cell.food > 255-growth {
					cell.food = 255
				} else {
					cell.food += growth
				}
			}
			if cell.agent == nil {
				continue
			}
			if moved[cell.agent] {
				continue
			}
			moved[cell.agent] = true
			if cell.agent.health == 0 { // agent with no health die
				env.cells[i][j].agent = nil
				continue
			}

			var input [inputLen]uint8
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
			cnt := 0
LoopLabel:
			for dist := 1; dist<=eyeSight; dist++ {
				for k := 0; k < 2*dist; k++ {
					var x, y [4]int
					x[0] = (envSize+i-dist+k)%envSize
					y[0] = (envSize+j-dist)%envSize
					x[1] = (envSize+i+dist-k)%envSize
					y[1] = (j+dist)%envSize
					x[2] = (i+dist)%envSize
					y[2] = (envSize+j-dist+k)%envSize
					x[3] = (envSize+i-dist)%envSize
					y[3] = (envSize+j+dist-k)%envSize
					for q := 0; q < 4; q++ {
						cell2 := env.cells[x[q]][y[q]]
						if cell2.agent != nil {
							input[6+4*cnt] = cell2.agent.appearance
							input[6+4*cnt+1] =
								uint8((x[q]<<4)+y[q])
							input[6+4*cnt+2] =
								cell2.agent.energy
							input[6+4*cnt+3] =
								cell2.agent.health
							cnt++
							if cnt == eyeRange {
								break LoopLabel
							}
						}
					}
				}
			}

			output := cell.agent.brain.react(input)

			if output&Eat != 0 {
				cell.agent.eat(i, j, env)
			}
			if output&Attack != 0 {
				cell.agent.attack(i, j, env)
			}
			if output&Mate != 0 {
				cell.agent.mate(i, j, env)
			}
			cell.agent.move(output&Move, i, j, env)
		}
	}
	img := image.NewPaletted(image.Rect(0, 0, envSize, envSize), colors)
	for i := 0; i < envSize; i++ {
		for j := 0; j < envSize; j++ {
			cell := &env.cells[i][j]
			if cell.agent == nil {
				img.Set(i, j, grassColor(cell.food))
			} else {
				img.Set(i, j, appearanceColor(cell.agent.appearance))
			}
		}
	}
	env.record.Image[iter] = img
}
