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
		env.record.Delay[i] = RecordDelay
	}
	return env
}

var (
	GrassColor1 = color.RGBA{30, 60, 30, 255}
	GrassColor2 = color.RGBA{60, 120, 60, 255}
	GrassColor3 = color.RGBA{90, 180, 90, 255}
	GrassColor4 = color.RGBA{120, 240, 120, 255}
	GrassColor5 = color.RGBA{240, 240, 120, 255}
	BackGroundColor = color.RGBA{255, 255, 255, 255}
	AnimatColor = color.RGBA{255, 0, 0, 255}
	GridColor   = color.RGBA{0, 0, 0, 255}
	NormalColor = color.RGBA{0, 255, 0, 255}
)

func faceColor(f Face) color.Color {
	return AnimatColor
}

func headColor(tf Face) color.Color {
	if tf == 0 {
		return NormalColor
	}
	return AnimatColor
}

func bodyColor(hp int) color.Color {
	if hp > 0 {
		return NormalColor
	}
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
	img := image.NewPaletted(image.Rect(0, 0, 31*EnvSize-1, 31*EnvSize-1), []color.Color{ BackGroundColor, AnimatColor, GridColor, NormalColor,
		//GrassColor1,
		//GrassColor2,
		//GrassColor3,
		//GrassColor4,
		//GrassColor5,
		//AnimatColor,
	})
	for i := 1; i < EnvSize; i++ {
		for j := 0; j < 31*EnvSize-1; j++ {
			img.Set(31*i, j, GridColor)
			img.Set(j, 31*i, GridColor)
		}
	}
	for i := 0; i < EnvSize; i++ {
		for j := 0; j < EnvSize; j++ {
			cell := &env.Cell[i][j]
			if cell.Animat == nil {
				continue
				//img.Set(i, j, grassColor(cell.Food))
			}
			//else {
			//	img.Set(i, j, faceColor(cell.Animat.Face))
			//}
			a := cell.Animat
			switch a.Direction {
			case Down:
				for ii := 31*i+10; ii < 31*i+20; ii++ {
					for jj := 31*j; jj < 31*j+10; jj++ {
						img.Set(ii,jj,headColor(a.TargetFace))
					}
					for jj := 31*j+10; jj < 31*j+20; jj++ {
						img.Set(ii,jj,bodyColor(a.Health))
					}
				}
			case Right:
				for jj := 31*j+10; jj < 31*j+20; jj++ {
					for ii := 31*i; ii < 31*i+10; ii++ {
						img.Set(ii,jj,headColor(a.TargetFace))
					}
					for ii := 31*i+10; ii < 31*i+20; ii++ {
						img.Set(ii,jj,bodyColor(a.Health))
					} 
				}
			case Up:
				for ii := 31*i+10; ii < 31*i+20; ii++ {
					for jj := 31*j+10; jj < 31*j+20; jj++ {
						img.Set(ii,jj,bodyColor(a.Health))
					}
					for jj := 31*j+20; jj < 31*j+30; jj++ {
						img.Set(ii,jj,headColor(a.TargetFace))
					}
				}
			case Left:
				for jj := 31*j+10; jj < 31*j+20; jj++ {
					for ii := 31*i+10; ii < 31*i+20; ii++ {
						img.Set(ii,jj,bodyColor(a.Health))
					}
					for ii := 31*i+20; ii < 31*i+30; ii++ {
						img.Set(ii,jj,headColor(a.TargetFace))
					}
				}
			}
		}
	}
	env.record.Image[iter] = img
}
