package game_template

import (
	"game_hub/core"
)

func (g *Game) CreateNew() core.GameInterface {
	return NewGame()
}

func (g *Game) GetId() string {
	return "gametemplate"
}

func (g *Game) GetStartState() core.State {
	return &StartState{}
}
