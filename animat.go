package main

type Direction uint8

const (
	Up Direction = iota
	Left
	Down
	Right
)

type Face uint8 // 0-7

type Animat struct {
	Health     int
	Face       Face
	TargetFace Face
	Brain      Brain
	Direction  Direction
}

func (a *Animat) Do(x, y int, env *Environment) {
	output, _ := a.Brain.React(a.Observe(x, y, env))

	a.Move(output.Move, x, y, env)
	a.Eat(x, y, env)
	a.TargetFace = 0
	if output.Attack {
		a.Attack(x, y, env)
	}
}

func (a *Animat) Observe(x, y int, env *Environment) *BrainInput {
	input := new(BrainInput)
	input.MyFace = a.Face
	input.MyTargetFace = a.TargetFace

	for i, d := range []Direction{Up, Left, Down, Right} {
		x2, y2 := nextLoc(x, y, d)
		cell := &env.Cell[x2][y2]
		nearby := &input.Nearby[i]
		nearby.MoreFood = env.Cell[x][y].Food < cell.Food
		if cell.Animat == nil {
			continue
		}
		nearby.HasAnimat = true
		nearby.MoreHealth = a.Health < cell.Animat.Health
		nearby.OtherFace = cell.Animat.Face
		nearby.OtherTargetFace = cell.Animat.TargetFace
	}
	return input
}

func (a *Animat) Eat(x, y int, env *Environment) {
	a.Health += int(env.Cell[x][y].Food)
}

func relLoc(x, y, rx, ry int) (int, int) {
	return (x + rx + EnvSize) % EnvSize, (y + ry + EnvSize) % EnvSize
}

// calculate the coordinate of the cell in front of a given cell and direction.
// the world is wrapped around.
func nextLoc(x, y int, dir Direction) (int, int) {
	switch dir {
	case Up:
		return relLoc(x, y, 0, -1)
	case Left:
		return relLoc(x, y, -1, 0)
	case Down:
		return relLoc(x, y, 0, 1)
	case Right:
		return relLoc(x, y, 1, 0)
	default:
		return x, y
	}
}

func (a *Animat) Attack(x, y int, env *Environment) {
	x2, y2 := nextLoc(x, y, a.Direction)
	target := env.Cell[x2][y2].Animat

	if target == nil {
		return
	}
	a.TargetFace = target.Face
	target.Health -= AttackDamage

	env.Reputation[target][a.Face] -= 1.0
	for other := range env.Reputation {
		if other == target {
			continue
		}
		env.Reputation[other][a.Face] -= env.Reputation[other][target.Face] / ReputationFactor
	}
}

func (a *Animat) Move(move Move, x, y int, env *Environment) {
	switch move {
	case MoveUp:
		a.Direction = Up
	case MoveLeft:
		a.Direction = Left
	case MoveDown:
		a.Direction = Down
	case MoveRight:
		a.Direction = Right
	}
	if move != Stay {
		a.Health -= MoveCost
		x2, y2 := nextLoc(x, y, a.Direction)
		if env.Cell[x2][y2].Animat != nil { // cannot move forward if front cell is blocked.
			return
		}
		env.Cell[x2][y2].Animat = a
		env.Cell[x][y].Animat = nil
	}
}
