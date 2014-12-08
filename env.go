package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"math/rand"
	"os"

	"github.com/taylorchu/rbm"
)

type EnvCell struct {
	Food   uint8
	Animat *Animat
}

type Environment struct {
	Cell     [EnvSize][EnvSize]EnvCell
	Relation map[*Animat]map[Color]float64
	Good     map[*Animat][][]float64
	record   gif.GIF
}

func (env *Environment) WriteFile(filename string) (err error) {
	f, err := os.Create(filename)
	if err != nil {
		return
	}
	defer f.Close()
	err = gif.EncodeAll(f, &env.record)
	return
}

func (env *Environment) relCell(x, y, rx, ry int) *EnvCell {
	x2, y2 := relLoc(x, y, rx, ry)
	return &env.Cell[x2][y2]
}

func NewEnvironment() *Environment {
	env := new(Environment)
	env.Relation = make(map[*Animat]map[Color]float64)
	env.Good = make(map[*Animat][][]float64)

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
		brain := NewRBMBrain()
		brain.Load(*brainData)
		a := &Animat{
			Brain: brain,
			Color: Color((i % 3) + 1),
		}
		env.Relation[a] = make(map[Color]float64)
		// no position collision
		var x, y int
		for {
			x = rand.Intn(EnvSize)
			y = rand.Intn(EnvSize)
			if env.Cell[x][y].Animat == nil {
				break
			}
		}
		env.Cell[x][y].Animat = a
	}

	env.record.Image = make([]*image.Paletted, Iteration)
	env.record.Delay = make([]int, Iteration)
	for i := 0; i < Iteration; i++ {
		env.record.Delay[i] = RecordDelay
	}
	return env
}

var (
	bgColor1    = color.RGBA{0, 0, 0, 255}
	bgColor2    = color.RGBA{255, 255, 255, 255}
	bgColor3    = color.RGBA{128, 128, 128, 255}
	gridColor   = color.RGBA{0, 0, 0, 255}
	grassColor1 = color.RGBA{30, 60, 30, 255}
	grassColor2 = color.RGBA{60, 120, 60, 255}
	grassColor3 = color.RGBA{90, 180, 90, 255}
	grassColor4 = color.RGBA{120, 240, 120, 255}
	grassColor5 = color.RGBA{240, 240, 120, 255}
	attackColor = color.RGBA{255, 0, 0, 255}
	normalColor = color.RGBA{0, 255, 0, 255}
)

func grassColor(grass uint8) color.Color {
	switch grass {
	case 0:
		return grassColor5
	case 1:
		return grassColor4
	case 2:
		return grassColor3
	case 3:
		return grassColor2
	default:
		return grassColor1
	}
}

func skinColor(c Color) color.Color {
	switch c {
	case 1:
		return bgColor1
	case 2:
		return bgColor2
	default:
		return bgColor3
	}
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

		// observe, think, and do something
		input := a.Observe(l.X, l.Y, env)
		output, raw := a.Brain.React(input)
		a.Do(output, l.X, l.Y, env)

		// train
		if *output != *expected(input) {
			continue
		}
		fmt.Printf("animat %d\n", a.Color)
		env.Good[a] = append(env.Good[a], raw)

		if len(env.Good[a]) == rbm.BatchSize {
			a.Brain.Reward(env.Good[a], Reward)
			env.Good[a] = nil
		}
	}
	env.drawFrame(iter)
}

func (env *Environment) drawFrame(iter int) {
	img := image.NewPaletted(image.Rect(0, 0, (CellPixel+1)*EnvSize-1, (CellPixel+1)*EnvSize-1), []color.Color{
		bgColor1,
		bgColor2,
		bgColor3,
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
			a := cell.Animat

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
						img.Set(ii, jj, skinColor(a.Color))
					}
				}
			}
			switch a.Direction {
			case Up:
				for ii := a1 + 1; ii < a2-1; ii++ {
					for jj := b1 + 1; jj < b1+1+HeadPixel; jj++ {
						if a.TargetColor == 0 {
							img.Set(ii, jj, normalColor)
						} else {
							img.Set(ii, jj, attackColor)
						}
					}
				}
			case Left:
				for ii := a1 + 1; ii < a1+1+HeadPixel; ii++ {
					for jj := b1 + 1; jj < b2-1; jj++ {
						if a.TargetColor == 0 {
							img.Set(ii, jj, normalColor)
						} else {
							img.Set(ii, jj, attackColor)
						}
					}
				}
			case Down:
				for ii := a1 + 1; ii < a2-1; ii++ {
					for jj := b2 - 2; jj > b2-2-HeadPixel; jj-- {
						if a.TargetColor == 0 {
							img.Set(ii, jj, normalColor)
						} else {
							img.Set(ii, jj, attackColor)
						}
					}
				}
			case Right:
				for ii := a2 - 2; ii > a2-2-HeadPixel; ii-- {
					for jj := b1 + 1; jj < b2-1; jj++ {
						if a.TargetColor == 0 {
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
