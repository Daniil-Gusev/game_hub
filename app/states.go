package app

import (
	"fmt"
	"game_hub/core"
)

type GameSelectionMenuState struct {
	core.BaseState
	AvailableGames []core.GameInterface
}

func NewGameSelectionMenu(availableGames []core.GameInterface) *GameSelectionMenuState {
	return &GameSelectionMenuState{
		AvailableGames: availableGames,
	}
}
func (s *GameSelectionMenuState) Id() string {
	return "game_selection_menu"
}
func (s *GameSelectionMenuState) Scope() core.Scope {
	return core.ScopeApp
}
func (s *GameSelectionMenuState) Init(ctx *core.GameContext, ui *core.UiContext) (core.State, error) {
	if err := ui.StateLocalizer.LoadTranslations("./app/states.json"); err != nil {
		return &core.ExitState{}, err
	}
	if err := ui.AppLocalizer.LoadTranslations("./app/translations.json"); err != nil {
		return &core.ExitState{}, err
	}
	if err := ui.AppLocalizer.LoadOptionalTranslations("./games/translations.json"); err != nil {
		return &core.ExitState{}, err
	}
	return s, nil
}
func (s *GameSelectionMenuState) Display(_ *core.GameContext, ui *core.UiContext) {
	ui.Console.Write(ui.GetLocalizedStateMsg(s, "welcome") + "\r\n")
	ui.Console.Write(fmt.Sprintf("0. %s\r\n", ui.GetLocalizedStateMsg(s, "exit_option")))
	ui.Console.Write(ui.GetLocalizedStateMsg(s, "available_games") + "\r\n\r\n")
	for i, game := range s.AvailableGames {
		name := ui.GetOptionalLocalizedMsg(ui.AppLocalizer, game.GetId(), "name")
		desc := ui.GetOptionalLocalizedMsg(ui.AppLocalizer, game.GetId(), "description")
		author := ui.GetOptionalLocalizedMsg(ui.AppLocalizer, game.GetId(), "author")
		ui.Console.Write(fmt.Sprintf("%d. %s.\r\n%s\r\n%s: %s.\r\n\r\n", i+1, name, desc, core.Capitalize(ui.GetLocalizedMsg(ui.AppLocalizer, "author")), author))
	}
	ui.Console.Write(ui.GetLocalizedStateMsg(s, "make_your_choice") + "\r\n")
	ui.Console.Write("> ")
}

func (s *GameSelectionMenuState) Handle(ctx *core.GameContext, ui *core.UiContext, input string) (core.State, error) {
	option, err := ui.Validator.ParseInt(input)
	if err != nil {
		return s, err
	}
	maxOption := len(s.AvailableGames)
	if option < 0 || option > maxOption {
		ui.Msg = ui.GetLocalizedStateMsg(s, "invalid_option") + "\r\n"
		return s, nil
	}
	if option == 0 {
		return &core.ExitState{}, nil
	}
	selectedGame := s.AvailableGames[option-1]
	return &core.InitGameState{Game: selectedGame}, nil
}
