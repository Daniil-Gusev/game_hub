package rockpaperscissors

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
		return nil, core.NewAppError(core.ErrInternal, "getting_gamedata_error", nil)
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
			Params:      func() map[string]any { return map[string]any{"rounds": game.TotalRounds} },
			NextState: func() core.State {
				game.Reset()
				return &GameState{}
			},
		},
		{Id: 2,
			Description: "select_rounds",
			NextState:   func() core.State { return &SelectRoundsState{} },
		},
	}
	return core.NewMenu(parentState, options, "")
}

type GameState struct{ BaseGameState }

func (g *GameState) Id() string {
	return "game"
}
func (g *GameState) Display(ctx *core.AppContext, ui *core.UiContext) {
	ui.DisplayText(fmt.Sprintf(ui.GetLocalizedStateMsg(g, "score")+"\r\n", g.game.PlayerScore, g.game.BotScore))
	ui.DisplayText(fmt.Sprintf(ui.GetLocalizedStateMsg(g, "current_round")+"\r\n", g.game.CurrentRound, g.game.TotalRounds))
	ui.DisplayText(ui.GetLocalizedStateMsg(g, "prompt") + "\r\n")
	ui.DisplayText("> ")
}
func (g *GameState) Handle(ctx *core.AppContext, ui *core.UiContext, input string) (core.State, error) {
	option, err := ui.Validator.ParseInt(input)
	if err != nil {
		return g, err
	}
	var playerMove Move
	switch option {
	case 1:
		playerMove = Rock
	case 2:
		playerMove = Scissors
	case 3:
		playerMove = Paper
	default:
		ui.Msg = ui.GetLocalizedStateMsg(g, "invalid_option") + "\r\n"
		return g, nil
	}
	return g.Play(ctx, ui, playerMove)
}
func (g *GameState) Play(ctx *core.AppContext, ui *core.UiContext, playerMove Move) (core.State, error) {
	g.game.MakePlayerMove(playerMove)
	g.game.MakeBotMove()
	ui.DisplayText(fmt.Sprintf(ui.GetLocalizedStateMsg(g, "moves_info")+"\r\n", ui.GetLocalizedMsg(ui.GameLocalizer, g.game.PlayerMove.String()), ui.GetLocalizedMsg(ui.GameLocalizer, g.game.BotMove.String())))
	result := g.game.PlayRound()
	switch result {
	case Winning:
		ui.DisplayText(ui.GetLocalizedStateMsg(g, "round_win") + "\r\n")
	case Loss:
		ui.DisplayText(ui.GetLocalizedStateMsg(g, "round_loss") + "\r\n")
	case Draw:
		ui.DisplayText(ui.GetLocalizedStateMsg(g, "round_draw") + "\r\n")
	}
	if g.game.CurrentRound > g.game.TotalRounds {
		return &EndGameState{}, nil
	}
	return g, nil
}

type EndGameState struct{ BaseGameState }

func (e *EndGameState) Id() string {
	return "end_game"
}
func (e *EndGameState) Display(ctx *core.AppContext, ui *core.UiContext) {
	ui.DisplayText(fmt.Sprintf(ui.GetLocalizedStateMsg(e, "score")+"\r\n", e.game.PlayerScore, e.game.BotScore))
	if e.game.CheckWin() {
		ui.DisplayText(ui.GetLocalizedStateMsg(e, "win") + "\r\n")
	} else if e.game.CheckLoss() {
		ui.DisplayText(ui.GetLocalizedStateMsg(e, "loss") + "\r\n")
	} else {
		ui.DisplayText(ui.GetLocalizedStateMsg(e, "draw") + "\r\n")
	}
}
func (e *EndGameState) Handle(ctx *core.AppContext, ui *core.UiContext, _ string) (core.State, error) {
	return NewMainMenu(ctx, ui, e.game), nil
}
func (e *EndGameState) RequiresInput() bool {
	return false
}

type SelectRoundsState struct{ BaseGameState }

func (s *SelectRoundsState) Id() string {
	return "select_rounds"
}
func (s *SelectRoundsState) Display(ctx *core.AppContext, ui *core.UiContext) {
	ui.DisplayText(ui.GetLocalizedStateMsg(s, "prompt") + "\r\n")
	ui.DisplayText(fmt.Sprintf(ui.GetLocalizedStateMsg(s, "current_value")+"\r\n", s.game.TotalRounds))
	ui.DisplayText("> ")
}
func (s *SelectRoundsState) Handle(ctx *core.AppContext, ui *core.UiContext, input string) (core.State, error) {
	num, err := ui.Validator.ParseOptionalIntInRange(input, s.game.TotalRounds, s.game.MinRounds, s.game.MaxRounds)
	if err != nil {
		return s, err
	}
	s.game.TotalRounds = num
	ui.Msg = fmt.Sprintf(ui.GetLocalizedStateMsg(s, "selected")+"\r\n", num)
	return ctx.GetPreviousState()
}
