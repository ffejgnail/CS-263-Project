package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"math/rand"
	"os"
)

type EnvCell struct {
	Food   uint8 // 0, 1, 2, 3, 4
	Animat *Animat
}

type Environment struct {
	Cell       [EnvSize][EnvSize]EnvCell
	Reputation map[*Animat]map[Face]float64
	Score      map[*Animat][]float64

	record gif.GIF
}

func (env *Environment) WriteTo(w io.Writer) error { // generating GIF
	return gif.EncodeAll(w, &env.record)
}

func (env *Environment) relCell(x, y, rx, ry int) *EnvCell {
	x2, y2 := relLoc(x, y, rx, ry)
	return &env.Cell[x2][y2]
}

func NewEnvironment() *Environment {
	env := new(Environment)
	env.Reputation = make(map[*Animat]map[Face]float64)
	env.Score = make(map[*Animat][]float64)

	cx := EnvSize / 2
	cy := EnvSize / 2
	for i := -4; i <= 4; i++ {
		for j := -4; j <= 4; j++ {
			for k := 0; k <= 4; k++ {
				if i*i+j*j < k*k {
					env.relCell(cx, cy, i, j).Food = uint8(4 - k)
					break
				}
			}
		}
	}

	for i := uint8(0); i < InitAnimatNum; i++ {
		x := rand.Intn(EnvSize)
		y := rand.Intn(EnvSize)
		a := &Animat{
			Brain:     NewRBMBrain(),
			Face:      Face(i % 8),
			Direction: Direction(i % 4),
		}
		env.Cell[x][y].Animat = a
		env.Reputation[a] = make(map[Face]float64)
		env.Score[a] = make([]float64, 2*TrainScope)
		f, err := os.Open(*brainData)
		if err != nil {
			continue
		}
		a.Brain.(*RBMBrain).Decode(f)
		f.Close()
	}

	env.record.Image = make([]*image.Paletted, RecordIteration)
	env.record.Delay = make([]int, RecordIteration)
	for i := 0; i < RecordIteration; i++ {
		env.record.Delay[i] = 5
	}
	return env
}

var (
	GrassColor1 = color.RGBA{30, 60, 30, 255}
	GrassColor2 = color.RGBA{60, 120, 60, 255}
	GrassColor3 = color.RGBA{90, 180, 90, 255}
	GrassColor4 = color.RGBA{120, 240, 120, 255}
	GrassColor5 = color.RGBA{240, 240, 120, 255}
	AnimatColor = color.RGBA{255, 0, 0, 255}
)

func faceColor(f Face) color.Color {
	return AnimatColor
}

func grassColor(grass uint8) color.Color {
	if grass > 3 {
		return GrassColor1
	}
	if grass > 2 {
		return GrassColor2
	}
	if grass > 1 {
		return GrassColor3
	}
	if grass > 0 {
		return GrassColor4
	}
	return GrassColor5
}

func (env *Environment) Run(iter int) {
	type Location struct {
		Animat *Animat
		X      int
		Y      int
	}
	var list []*Location
	for i := 0; i < EnvSize; i++ {
		for j := 0; j < EnvSize; j++ {
			cell := &env.Cell[i][j]
			if cell.Animat == nil {
				continue
			}
			//if cell.Animat.Health == 0 { // no health die
			//f, err := os.Create(*brainData)
			//if err != nil {
			//	continue
			//}
			//cell.Animat.Brain.(*RBMBrain).Encode(f)
			//f.Close()
			//	env.Cell[i][j].Animat = nil
			//	continue
			//}
			list = append(list, &Location{
				Animat: cell.Animat,
				X:      i,
				Y:      j,
			})
		}
	}
	for _, r := range rand.Perm(len(list)) {
		l := list[r]
		a := l.Animat
		// train
		if Iteration-iter <= RecordIteration {
			fmt.Printf("animat %d score: %f; ",a.Face,env.Score[a][0])
		}
		a.Brain.Reward(env.Score[a][0])

		// observe, think, and do something
		a.Do(l.X, l.Y, env)
		fitness := float64(a.Health)
		//for f := range env.Reputation[a] {
		//	fitness += env.Reputation[a][f]
		//}

		if Iteration-iter <= RecordIteration {
			fmt.Printf("health: %f\n", fitness)
		}

		env.Score[a] = append(env.Score[a][1:], 0)
		for k := 0; k < TrainScope; k++ {
			env.Score[a][k] += fitness
		}
		for k := TrainScope; k < 2*TrainScope; k++ {
			env.Score[a][k] -= fitness
		}
	}
	if Iteration-iter <= RecordIteration {
		env.drawFrame(Iteration - iter - 1)
	}
}

func (env *Environment) drawFrame(iter int) {
	img := image.NewPaletted(image.Rect(0, 0, EnvSize, EnvSize), []color.Color{
		GrassColor1,
		GrassColor2,
		GrassColor3,
		GrassColor4,
		GrassColor5,
		AnimatColor,
	})
	for i := 0; i < EnvSize; i++ {
		for j := 0; j < EnvSize; j++ {
			cell := &env.Cell[i][j]
			if cell.Animat == nil {
				img.Set(i, j, grassColor(cell.Food))
			} else {
				img.Set(i, j, faceColor(cell.Animat.Face))
			}
		}
	}
	env.record.Image[iter] = img
}
