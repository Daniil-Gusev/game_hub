package game_template

import (
	"game_hub/core"
)

type Game struct {
	RandomGenerator *core.RandomGenerator
}

func NewGame() *Game {
	return &Game{
		RandomGenerator: core.NewRandomGenerator(),
	}
}
