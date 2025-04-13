package guessnumber

import (
	"game_hub/core"
)

type RestartCommand struct{ core.GameCommand }

func (c *RestartCommand) Id() string {
	return "restart"
}
func (c *RestartCommand) Execute(ctx *core.GameContext, ui *core.UiContext, args []string) (core.State, error) {
	return &StartGameState{}, nil
}
