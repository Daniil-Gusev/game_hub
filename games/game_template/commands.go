package gametemplate

import (
	"game_hub/core"
)

type CustomActionCommand struct{ core.GameCommand }

func (c *CustomActionCommand) Id() string {
	return "custom_action"
}

func (c *CustomActionCommand) Execute(ctx *core.GameContext, ui *core.UiContext, args []string) (core.State, error) {
	// Реализуйте логику команды
	return ctx.GetCurrentState()
}
