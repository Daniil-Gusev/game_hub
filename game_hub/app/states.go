package app

import (
	"fmt"
	"game_hub/core"
	"game_hub/utils"
)

type BaseAppState struct{ core.BaseState }

func (s *BaseAppState) Scope() core.Scope {
	return core.ScopeApp
}

type StartState struct{ BaseAppState }

func (s *StartState) Init(ctx *core.AppContext, ui *core.UiContext) (core.State, error) {
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

func (s *StartState) Handle(ctx *core.AppContext, ui *core.UiContext, _ string) (core.State, error) {
	return NewMainMenu(ctx, ui), nil
}

func (s *StartState) RequiresInput() bool {
	return false
}

type MainMenuState struct{ BaseAppState }

func (m *MainMenuState) Id() string {
	return "main_menu"
}

func NewMainMenu(ctx *core.AppContext, ui *core.UiContext) *core.MenuState {
	parentState := &MainMenuState{}
	options := []core.MenuOption{
		{Id: 0,
			Description: "exit_option",
			NextState:   func() core.State { return &core.ExitState{} },
		},
		{Id: 1,
			Description: "play_option",
			NextState: func() core.State {
				return NewGameSelectionMenu(ctx.AvailableGames)
			},
		},
		{Id: 2,
			Description: "change_language_option",
			NextState: func() core.State {
				return NewLanguageSelectionMenu(ui.LocalizationManager.AvailableLanguages())
			},
		},
	}
	return core.NewMenu(parentState, options, "")
}

type GameSelectionMenuState struct {
	BaseAppState
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

func (s *GameSelectionMenuState) Display(ctx *core.AppContext, ui *core.UiContext) {
	ui.DisplayText(ui.GetLocalizedStateMsg(s, "welcome") + "\r\n")
	ui.DisplayText(fmt.Sprintf("0. %s\r\n", ui.GetLocalizedStateMsg(s, "exit_option")))
	ui.DisplayText(ui.GetLocalizedStateMsg(s, "available_games") + "\r\n\r\n")
	for i, game := range s.AvailableGames {
		name := ui.GetOptionalLocalizedMsg(ui.AppLocalizer, game.GetId(), "name")
		desc := ui.GetOptionalLocalizedMsg(ui.AppLocalizer, game.GetId(), "description")
		author := ui.GetOptionalLocalizedMsg(ui.AppLocalizer, game.GetId(), "author")
		ui.DisplayText(fmt.Sprintf("%d. %s.\r\n%s\r\n%s: %s.\r\n\r\n", i+1, name, desc, utils.Capitalize(ui.GetLocalizedStateMsg(s, "author")), author))
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
		ui.DisplayText(ui.GetLocalizedStateMsg(s, "invalid_option") + "\r\n")
		return s, nil
	}
	if option == 0 {
		return NewMainMenu(ctx, ui), nil
	}
	selectedGame := s.AvailableGames[option-1]
	return &core.InitGameState{Game: selectedGame}, nil
}

type LanguageSelectionMenuState struct {
	BaseAppState
	availableLanguages []core.Language
}

func NewLanguageSelectionMenu(langs []core.Language) *LanguageSelectionMenuState {
	return &LanguageSelectionMenuState{
		availableLanguages: langs,
	}
}

func (m *LanguageSelectionMenuState) Id() string {
	return "language_selection_menu"
}

func (m *LanguageSelectionMenuState) Display(ctx *core.AppContext, ui *core.UiContext) {
	ui.DisplayText(ui.GetLocalizedStateMsg(m, "available_languages") + "\r\n\r\n")
	for _, lang := range m.availableLanguages {
		ui.DisplayText(fmt.Sprintf("%s: %s.\r\n", lang.Code, lang.Name))
	}
	ui.DisplayText("\r\n" + ui.GetLocalizedStateMsg(m, "prompt") + "\r\n")
}

func (s *LanguageSelectionMenuState) Handle(ctx *core.AppContext, ui *core.UiContext, input string) (core.State, error) {
	for _, lang := range s.availableLanguages {
		if lang.Code == input {
			if err := ui.LocalizationManager.SetCurrentLanguage(lang.Code); err != nil {
				return s, err
			}
			ui.CommandRegistry.UpdateAliases()
			ui.DisplayText(fmt.Sprintf(ui.GetLocalizedStateMsg(s, "selected")+"\r\n", lang.Name))
			return ctx.GetPreviousState()
		}
	}
	ui.DisplayText(ui.GetLocalizedStateMsg(s, "invalid_input") + "\r\n")
	return s, nil
}

func (s *LanguageSelectionMenuState) GetCommands() []core.Command {
	return []core.Command{
		&core.BackCommand{},
	}
}
