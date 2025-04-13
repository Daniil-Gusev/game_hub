package game_template

import (
	"game_hub/core"
)

type BaseGameState struct {
	core.BaseState
	game *Game
}

func (b *BaseGameState) Scope() core.Scope {
	return core.ScopeGame
}

func (b *BaseGameState) Init(ctx *core.GameContext, ui *core.UiContext) (core.State, error) {
	game, ok := ctx.Game.(*Game)
	if !ok {
		return nil, core.NewAppError(core.ErrInternal, "getting_gamedata_error", nil)
	}
	b.game = game
	return b, nil
}

type StartState struct{ BaseGameState }

func (s *StartState) Handle(ctx *core.GameContext, ui *core.UiContext, _ string) (core.State, error) {
	return NewMainMenu(ctx, ui, s.game), nil
}

func (s *StartState) RequiresInput() bool {
	return false
}

type MainMenuState struct{ BaseGameState }

func (m *MainMenuState) Id() string {
	return "main_menu"
}

func NewMainMenu(ctx *core.GameContext, ui *core.UiContext, game *Game) *core.MenuState {
	parentState := &MainMenuState{}
	options := []core.MenuOption{
		{
			Id:          0,
			Description: "exit_option",
			NextState:   func() core.State { return &core.GameExitState{} },
		},
		{
			Id:          1,
			Description: "start_game",
			NextState:   func() core.State { return &GameState{} },
		},
	}
	return core.NewMenu(parentState, options, "")
}

type GameState struct{ BaseGameState }

func (g *GameState) Id() string {
	return "game"
}

func (g *GameState) Display(ctx *core.GameContext, ui *core.UiContext) {
	ui.Console.Write(ui.GetLocalizedStateMsg(g, "prompt") + "\r\n")
	ui.Console.Write("> ")
}

func (g *GameState) Handle(ctx *core.GameContext, ui *core.UiContext, input string) (core.State, error) {
	// Реализуйте логику обработки ввода
	return g, nil
}
