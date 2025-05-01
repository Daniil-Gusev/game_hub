package app

import (
	"fmt"
	"game_hub/core"
	"game_hub/utils"
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
func (s *GameSelectionMenuState) Init(ctx *core.AppContext, ui *core.UiContext) (core.State, error) {
	if err := ui.StateLocalizer.LoadTranslations(ctx.Config.Paths.AppStatesPath()); err != nil {
		return &core.ExitState{}, err
	}
	if err := ui.AppLocalizer.LoadTranslations(ctx.Config.Paths.AppTranslationsPath()); err != nil {
		return &core.ExitState{}, err
	}
	if err := ui.AppLocalizer.LoadOptionalTranslations(ctx.Config.Paths.GamesTranslationsPath()); err != nil {
		return &core.ExitState{}, err
	}
	return s, nil
}
func (s *GameSelectionMenuState) Display(_ *core.AppContext, ui *core.UiContext) {
	ui.DisplayText(ui.GetLocalizedStateMsg(s, "welcome") + "\r\n")
	ui.DisplayText(fmt.Sprintf("0. %s\r\n", ui.GetLocalizedStateMsg(s, "exit_option")))
	ui.DisplayText(ui.GetLocalizedStateMsg(s, "available_games") + "\r\n\r\n")
	for i, game := range s.AvailableGames {
		name := ui.GetOptionalLocalizedMsg(ui.AppLocalizer, game.GetId(), "name")
		desc := ui.GetOptionalLocalizedMsg(ui.AppLocalizer, game.GetId(), "description")
		author := ui.GetOptionalLocalizedMsg(ui.AppLocalizer, game.GetId(), "author")
		ui.DisplayText(fmt.Sprintf("%d. %s.\r\n%s\r\n%s: %s.\r\n\r\n", i+1, name, desc, utils.Capitalize(ui.GetLocalizedMsg(ui.AppLocalizer, "author")), author))
	}
	ui.DisplayText(ui.GetLocalizedStateMsg(s, "make_your_choice") + "\r\n")
}

func (s *GameSelectionMenuState) Handle(ctx *core.AppContext, ui *core.UiContext, input string) (core.State, error) {
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
