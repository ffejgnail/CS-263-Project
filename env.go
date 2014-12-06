package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"math/rand"

	"github.com/r9y9/nnet/rbm"
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
				if i*i+j*j <= k*k {
					env.relCell(cx, cy, i, j).Food = uint8(4 - k)
					break
				}
			}
		}
	}

	for i := uint8(0); i < InitAnimatNum; i++ {
		brain := NewRBMBrain2()
		if r, err := rbm.Load(*brainData); err == nil {
			brain.RBM = r
		}
		x := rand.Intn(EnvSize)
		y := rand.Intn(EnvSize)
		a := &Animat{
			Brain: brain,
			Face:  Face(i % 8),
		}
		env.Cell[x][y].Animat = a
		env.Reputation[a] = make(map[Face]float64)
		env.Score[a] = make([]float64, 2*TrainScope)
	}

	env.record.Image = make([]*image.Paletted, RecordIteration)
	env.record.Delay = make([]int, RecordIteration)
	for i := 0; i < RecordIteration; i++ {
		env.record.Delay[i] = RecordDelay
	}
	return env
}

var (
	backgroundColor = color.RGBA{255, 255, 255, 255}
	gridColor       = color.RGBA{0, 0, 0, 255}
	grassColor1     = color.RGBA{30, 60, 30, 255}
	grassColor2     = color.RGBA{60, 120, 60, 255}
	grassColor3     = color.RGBA{90, 180, 90, 255}
	grassColor4     = color.RGBA{120, 240, 120, 255}
	grassColor5     = color.RGBA{240, 240, 120, 255}
	attackColor     = color.RGBA{255, 0, 0, 255}
	normalColor     = color.RGBA{0, 255, 0, 255}
)

func grassColor(grass uint8) color.Color {
	if grass > 3 {
		return grassColor1
	}
	if grass > 2 {
		return grassColor2
	}
	if grass > 1 {
		return grassColor3
	}
	if grass > 0 {
		return grassColor4
	}
	return grassColor5
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
		fitness := float64(a.Health)
		//for f := range env.Reputation[a] {
		//	fitness += env.Reputation[a][f]
		//}

		env.Score[a] = append(env.Score[a][1:], 0)
		for k := 0; k < TrainScope; k++ {
			env.Score[a][k] += fitness
		}
		for k := TrainScope; k < 2*TrainScope; k++ {
			env.Score[a][k] -= fitness
		}

		// observe, think, and do something
		a.Do(l.X, l.Y, env)

		// train
		fmt.Printf("animat %d: %f\n", a.Face, env.Score[a][0])
		a.Brain.Reward(env.Score[a][0])
	}
	if Iteration-RecordIteration <= iter {
		env.drawFrame(iter - Iteration + RecordIteration)
	}
}

func (env *Environment) drawFrame(iter int) {
	img := image.NewPaletted(image.Rect(0, 0, (CellPixel+1)*EnvSize-1, (CellPixel+1)*EnvSize-1), []color.Color{
		backgroundColor,
		gridColor,
		grassColor1,
		grassColor2,
		grassColor3,
		grassColor4,
		grassColor5,
		attackColor,
		normalColor,
	})

	for i := 1; i < EnvSize; i++ {
		for j := 0; j < (CellPixel+1)*EnvSize-1; j++ {
			img.Set((CellPixel+1)*i, j, gridColor)
			img.Set(j, (CellPixel+1)*i, gridColor)
		}
	}
	for i := 0; i < EnvSize; i++ {
		for j := 0; j < EnvSize; j++ {
			cell := &env.Cell[i][j]

			for ii := i*(CellPixel+1) + 1; ii < (i+1)*(CellPixel+1); ii++ {
				for jj := j*(CellPixel+1) + 1; jj < (j+1)*(CellPixel+1); jj++ {
					img.Set(ii, jj, grassColor(cell.Food))
				}
			}

			if cell.Animat == nil {
				continue
			}
			a1 := i*(CellPixel+1) + MarginPixel + 1
			a2 := (i+1)*(CellPixel+1) - MarginPixel
			b1 := j*(CellPixel+1) + MarginPixel + 1
			b2 := (j+1)*(CellPixel+1) - MarginPixel
			for ii := a1; ii < a2; ii++ {
				for jj := b1; jj < b2; jj++ {
					if ii == a1 || ii == a2-1 ||
						jj == b1 || jj == b2-1 {
						img.Set(ii, jj, gridColor)
					} else {
						img.Set(ii, jj, backgroundColor)
					}
				}
			}
			a := cell.Animat
			switch a.Direction {
			case Up:
				for ii := a1 + 1; ii < a2-1; ii++ {
					for jj := b1 + 1; jj < b1+1+HeadPixel; jj++ {
						if a.TargetFace == 0 {
							img.Set(ii, jj, normalColor)
						} else {
							img.Set(ii, jj, attackColor)
						}
					}
				}
			case Left:
				for ii := a1 + 1; ii < a1+1+HeadPixel; ii++ {
					for jj := b1 + 1; jj < b2-1; jj++ {
						if a.TargetFace == 0 {
							img.Set(ii, jj, normalColor)
						} else {
							img.Set(ii, jj, attackColor)
						}
					}
				}
			case Down:
				for ii := a1 + 1; ii < a2-1; ii++ {
					for jj := b2 - 2; jj > b2-2-HeadPixel; jj-- {
						if a.TargetFace == 0 {
							img.Set(ii, jj, normalColor)
						} else {
							img.Set(ii, jj, attackColor)
						}
					}
				}
			case Right:
				for ii := a2 - 2; ii > a2-2-HeadPixel; ii-- {
					for jj := b1 + 1; jj < b2-1; jj++ {
						if a.TargetFace == 0 {
							img.Set(ii, jj, normalColor)
						} else {
							img.Set(ii, jj, attackColor)
						}
					}
				}
			}
		}
	}
	env.record.Image[iter] = img
}
