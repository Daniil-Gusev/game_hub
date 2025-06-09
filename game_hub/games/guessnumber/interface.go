package guessnumber

import (
	"game_hub/core"
)

func (g *Game) CreateNew() core.GameInterface {
	return NewGame()
}

func (g *Game) GetId() string {
	return "guessnumber"
}

func (g *Game) GetStartState() core.State {
	return &StartState{}
}
