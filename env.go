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
	cells  [envSize][envSize]EnvCell
	friend map[*Agent]map[*Agent]float32
	score  map[*Agent][]float32

	record gif.GIF
}

func (env *Environment) WriteTo(w io.Writer) error { // generating GIF
	return gif.EncodeAll(w, &env.record)
}

func (env *Environment) setup() {
	env.friend = make(map[*Agent]map[*Agent]float32)
	env.score = make(map[*Agent][]float32)

	for i := 0; i < envSize; i++ {
		for j := 0; j < envSize; j++ {
			//env.cells[i][j].growth = uint8(rand.Intn(2)) + 2
			env.cells[i][j].food = uint8(rand.Intn(2))
		}
	}

	for i := uint8(0); i < initAgentNum; i++ {
		x := rand.Intn(envSize)
		y := rand.Intn(envSize)
		agent := new(Agent)
		agent.brain = NewRBMBrain()
		agent.energy = initEnergy
		agent.health = initHealth
		agent.appearance = i
		agent.dir = Direction(rand.Intn(4))
		env.cells[x][y].agent = agent
		env.friend[agent] = make(map[*Agent]float32)
		env.score[agent] = make([]float32, 2*trainScopeLen)
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
		color.RGBA{255, 255, 255, 255},
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
	if grass > 0 {
		return colors[4]
	}
	return colors[6]
}

func (env *Environment) run(iter int) {
	type AgentLocation struct {
		agent *Agent
		i     int
		j     int
	}
	var agents []*AgentLocation
	for i := 0; i < envSize; i++ {
		for j := 0; j < envSize; j++ {
			cell := &env.cells[i][j]
			if iter%grassGrowFreq == 0 && cell.growth > 0 {
				growth := uint8(rand.Intn(int(cell.growth)))
				cell.food = add(cell.food, growth)
			}
			if cell.agent == nil {
				continue
			}
			if cell.agent.health == 0 { // agent with no health die
				env.cells[i][j].agent = nil
				continue
			}
			agents = append(agents, &AgentLocation{
				agent: cell.agent,
				i:     i,
				j:     j,
			})
		}
	}
	for _, r := range rand.Perm(len(agents)) {
		al := agents[r]
		agent := al.agent
		// train
		agent.brain.train(env.score[agent][0])

		// observe, think, and do something
		agent.do(al.i, al.j, env)
		var fitness float32
		for k := range env.friend[agent] {
			fitness += float32(k.health) * env.friend[agent][k]
		}

		env.score[agent] = append(env.score[agent][1:], 0)
		for k := 0; k < trainScopeLen; k++ {
			env.score[agent][k] += fitness
		}
		for k := trainScopeLen; k < 2*trainScopeLen; k++ {
			env.score[agent][k] -= fitness
		}
	}
	env.drawFrame(iter)
}

func (env *Environment) drawFrame(iter int) {
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
