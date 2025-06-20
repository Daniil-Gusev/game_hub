package guessnumber

import (
	"fmt"
	"game_hub/core"
)

type BaseGameState struct {
	core.BaseState
	game *Game
}

func (b *BaseGameState) Scope() core.Scope {
	return core.ScopeGame
}

func (b *BaseGameState) Init(ctx *core.AppContext, ui *core.UiContext) (core.State, error) {
	game, ok := ctx.Game.(*Game)
	if !ok {
		return &core.GameExitState{}, core.NewAppError(core.ErrInternal, "getting_gamedata_error", nil)
	}
	b.game = game
	return b, nil
}

type StartState struct{ BaseGameState }

func (s *StartState) Handle(ctx *core.AppContext, ui *core.UiContext, _ string) (core.State, error) {
	return NewMainMenu(ctx, ui, s.game), nil
}

func (s *StartState) RequiresInput() bool {
	return false
}

type MainMenuState struct{ BaseGameState }

func (m *MainMenuState) Id() string {
	return "main_menu"
}
func NewMainMenu(ctx *core.AppContext, ui *core.UiContext, game *Game) *core.MenuState {
	parentState := &MainMenuState{}
	options := []core.MenuOption{
		{Id: 0,
			Description: "exit_option",
			NextState:   func() core.State { return &core.GameExitState{} },
		},
		{Id: 1,
			Description: "start_game",
			Params: func() map[string]any {
				return map[string]any{"difficulty": ui.GetLocalizedMsg(ui.GameLocalizer, game.Difficulty.String())}
			},
			NextState: func() core.State { return &SelectMinNumberState{} },
		},
		{Id: 2,
			Description: "select_difficulty",
			NextState:   func() core.State { return &SelectDifficultyMenuState{} },
		},
	}
	return core.NewMenu(parentState, options, "")
}

// выбор минимального числа диапазона угадывания
type SelectMinNumberState struct{ BaseGameState }

func (s *SelectMinNumberState) Id() string {
	return "select_min_number"
}

func (s *SelectMinNumberState) Display(ctx *core.AppContext, ui *core.UiContext) {
	ui.DisplayText(ui.GetLocalizedStateMsg(s, "prompt") + "\r\n")
	ui.DisplayText(ui.GetLocalizedMsg(ui.GameLocalizer, "press_enter") + "\r\n")
	ui.DisplayText(fmt.Sprintf(ui.GetLocalizedMsg(ui.GameLocalizer, "current_value")+"\r\n", s.game.MinNumber))
}

func (s *SelectMinNumberState) Handle(ctx *core.AppContext, ui *core.UiContext, input string) (core.State, error) {
	num, err := ui.Validator.ParseOptionalIntInRange(input, s.game.MinNumber, s.game.MinRangeNumber, s.game.MaxRangeNumber)
	if err != nil {
		return s, err
	}
	s.game.MinNumber = num
	return &SelectMaxNumberState{}, nil
}

func (s *SelectMinNumberState) GetCommands() []core.Command {
	return []core.Command{
		&core.BackCommand{},
	}
}

// выбор максимального числа диапазона угадывания
type SelectMaxNumberState struct{ BaseGameState }

func (s *SelectMaxNumberState) Id() string {
	return "select_max_number"
}

func (s *SelectMaxNumberState) Display(ctx *core.AppContext, ui *core.UiContext) {
	ui.DisplayText(ui.GetLocalizedStateMsg(s, "prompt") + "\r\n")
	ui.DisplayText(ui.GetLocalizedMsg(ui.GameLocalizer, "press_enter") + "\r\n")
	ui.DisplayText(fmt.Sprintf(ui.GetLocalizedMsg(ui.GameLocalizer, "current_value")+"\r\n", s.game.MaxNumber))
}

func (s *SelectMaxNumberState) Handle(ctx *core.AppContext, ui *core.UiContext, input string) (core.State, error) {
	num, err := ui.Validator.ParseOptionalIntInRange(input, s.game.MaxNumber, s.game.MinNumber, s.game.MaxRangeNumber)
	if err != nil {
		return s, err
	}
	if (num - s.game.MinNumber) < s.game.MinRangeSize {
		ui.DisplayText(fmt.Sprintf(ui.GetLocalizedStateMsg(s, "range_too_small")+"\r\n", s.game.MinRangeSize))
		return s, nil
	}
	s.game.MaxNumber = num
	return &StartGameState{}, nil
}

func (s *SelectMaxNumberState) GetCommands() []core.Command {
	return []core.Command{
		&core.BackCommand{},
	}
}

// инициализация иначало игры
type StartGameState struct{ BaseGameState }

func (g *StartGameState) Id() string {
	return "start_game"
}

func (g *StartGameState) Init(ctx *core.AppContext, ui *core.UiContext) (core.State, error) {
	newState, initErr := g.BaseGameState.Init(ctx, ui)
	if initErr != nil {
		return newState, initErr
	}
	if err := g.game.Prepare(); err != nil {
		return NewMainMenu(ctx, ui, g.game), err
	}
	return g, nil
}

func (g *StartGameState) Display(ctx *core.AppContext, ui *core.UiContext) {
	ui.DisplayText(fmt.Sprintf(ui.GetLocalizedStateMsg(g, "game_start")+"\r\n", g.game.MinNumber, g.game.MaxNumber, g.game.GetAttempts()))
}

func (g *StartGameState) Handle(_ *core.AppContext, _ *core.UiContext, _ string) (core.State, error) {
	return &GameState{}, nil
}

func (g *StartGameState) RequiresInput() bool {
	return false
}

// процесс угадывания числа
type GameState struct{ BaseGameState }

func (g *GameState) Id() string {
	return "game"
}

func (g *GameState) Display(ctx *core.AppContext, ui *core.UiContext) {
	ui.DisplayText(fmt.Sprintf(ui.GetLocalizedStateMsg(g, "attempts_left")+"\r\n", g.game.GetAttempts()))
}

func (g *GameState) Handle(ctx *core.AppContext, ui *core.UiContext, input string) (core.State, error) {
	num, err := ui.Validator.ParseIntInRange(input, g.game.MinNumber, g.game.MaxNumber)
	if err != nil {
		return g, err
	}
	return g.Guess(ctx, ui, num)
}

func (g *GameState) Guess(ctx *core.AppContext, ui *core.UiContext, guess int) (core.State, error) {
	g.game.MakeGuess(guess)
	if g.game.CheckWin() {
		return &EndGameState{}, nil
	}
	if g.game.CheckLoss() {
		return &EndGameState{}, nil
	}
	ui.DisplayText(fmt.Sprintf("%s\r\n", ui.GetLocalizedStateMsg(g, g.game.GetHint(guess))))
	return g, nil
}

func (g *GameState) GetCommands() []core.Command {
	return []core.Command{
		&core.ExitCommand{},
		&RestartCommand{},
	}
}

// конец игры
type EndGameState struct{ BaseGameState }

func (e *EndGameState) Id() string {
	return "end_game"
}

func (e *EndGameState) Display(ctx *core.AppContext, ui *core.UiContext) {
	if e.game.CheckWin() {
		ui.DisplayText(ui.GetLocalizedStateMsg(e, "win") + "\r\n")
	} else {
		ui.DisplayText(ui.GetLocalizedStateMsg(e, "loss") + "\r\n")
	}
}

func (e *EndGameState) Handle(ctx *core.AppContext, ui *core.UiContext, _ string) (core.State, error) {
	return NewEndMenu(ctx, ui, e.game), nil
}

func (e *EndGameState) RequiresInput() bool {
	return false
}

// меню после игры
type EndGameMenuState struct{ BaseGameState }

func (m *EndGameMenuState) Id() string {
	return "end_game_menu"
}

func NewEndMenu(ctx *core.AppContext, ui *core.UiContext, game *Game) *core.MenuState {
	parentState := &EndGameMenuState{}
	options := []core.MenuOption{
		{Id: 1,
			Description: "retry",
			NextState:   func() core.State { return &StartGameState{} },
		},
		{Id: 2,
			Description: "change_difficulty",
			NextState:   func() core.State { return &SelectDifficultyMenuState{} },
		},
		{Id: 3,
			Description: "main_menu",
			NextState:   func() core.State { return NewMainMenu(ctx, ui, game) },
		},
	}
	return core.NewMenu(parentState, options, "")
}

// меню выбора уровня сложности
type SelectDifficultyMenuState struct{ BaseGameState }

func (s *SelectDifficultyMenuState) Id() string {
	return "select_difficulty_menu"
}

func (s *SelectDifficultyMenuState) Display(ctx *core.AppContext, ui *core.UiContext) {
	ui.DisplayText(ui.GetLocalizedStateMsg(s, "prompt") + "\r\n")
	for d := VeryEasy; d <= VeryHard; d++ {
		ui.DisplayText(fmt.Sprintf("%d. %s.\r\n", d, ui.GetLocalizedMsg(ui.GameLocalizer, d.String())))
	}
}

func (s *SelectDifficultyMenuState) Handle(ctx *core.AppContext, ui *core.UiContext, input string) (core.State, error) {
	num, err := ui.Validator.ParseInt(input)
	if err != nil {
		ui.DisplayText(ui.GetLocalizedStateMsg(s, "invalid_input") + "\r\n")
		return s, nil
	}
	diff := Difficulty(num)
	if diff < VeryEasy || diff > VeryHard {
		ui.DisplayText(ui.GetLocalizedStateMsg(s, "invalid_option") + "\r\n")
		return s, nil
	}
	s.game.Difficulty = diff
	ui.DisplayText(fmt.Sprintf(ui.GetLocalizedStateMsg(s, "selected")+"\r\n", ui.GetLocalizedMsg(ui.GameLocalizer, diff.String())))
	return ctx.GetPreviousState()
}

func (s *SelectDifficultyMenuState) GetCommands() []core.Command {
	return []core.Command{
		&core.BackCommand{},
	}
}
