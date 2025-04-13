package core

import (
	"fmt"
	"os"
	"sort"
)

type InitGameState struct {
	BaseState
	Game GameInterface
}

func (g *InitGameState) Id() string {
	return "init_game"
}
func (g *InitGameState) Init(ctx *GameContext, ui *UiContext) (State, error) {
	basePath := fmt.Sprintf("./games/%s/", g.Game.GetId())
	statesPath := basePath + "states.json"
	if _, err := os.Stat(statesPath); err != nil {
		return &GameExitState{}, NewAppError(ErrLocalization, "file_open_error", map[string]any{
			"file":  statesPath,
			"error": fmt.Sprintf("%v", err),
		})
	}
	if err := ui.StateLocalizer.LoadTranslations(statesPath); err != nil {
		return &GameExitState{}, err
	}
	commandsPath := basePath + "commands.json"
	if _, err := os.Stat(commandsPath); err == nil {
		if err := ui.CommandRegistry.localLocalizer.LoadTranslations(commandsPath); err != nil {
			return &GameExitState{}, err
		}
	} else if !os.IsNotExist(err) {
		return &GameExitState{}, NewAppError(ErrLocalization, "file_open_error", map[string]any{
			"file":  commandsPath,
			"error": err,
		})
	}
	translationsPath := basePath + "translations.json"
	if _, err := os.Stat(translationsPath); err == nil {
		if err := ui.GameLocalizer.LoadTranslations(translationsPath); err != nil {
			return &GameExitState{}, err
		}
	} else if !os.IsNotExist(err) {
		return &GameExitState{}, NewAppError(ErrLocalization, "file_open_error", map[string]any{
			"file":  translationsPath,
			"error": err,
		})
	}
	return g, nil
}
func (g *InitGameState) Display(ctx *GameContext, ui *UiContext) {
	ui.Console.Write(fmt.Sprintf(ui.GetLocalizedStateMsg(g, "game_welcome"), ui.GetOptionalLocalizedMsg(ui.AppLocalizer, g.Game.GetId(), "name")) + "\r\n")
}
func (g *InitGameState) Handle(ctx *GameContext, ui *UiContext, input string) (State, error) {
	ctx.Game = g.Game.CreateNew()
	return g.Game.GetStartState(), nil
}
func (g *InitGameState) RequiresInput() bool {
	return false
}

// выход из приложения
type ExitState struct{ BaseState }

func (e *ExitState) Id() string {
	return "exit"
}
func (e *ExitState) Display(_ *GameContext, ui *UiContext) {
	ui.Console.Write(ui.GetLocalizedStateMsg(e, "exit") + "\r\n")
}
func (e *ExitState) Handle(ctx *GameContext, _ *UiContext, _ string) (State, error) {
	ctx.AppIsRunning = false
	return e, nil
}
func (e *ExitState) RequiresInput() bool {
	return false
}

// выход из игры
type GameExitState struct{ BaseState }

func (e *GameExitState) Id() string {
	return "game_exit"
}
func (e *GameExitState) RequiresInput() bool {
	return false
}
func (e *GameExitState) Display(_ *GameContext, ui *UiContext) {}

// служебное, возвращает в меню выбора игр в основном цикле
func (e *GameExitState) Handle(ctx *GameContext, ui *UiContext, _ string) (State, error) {
	ui.GameLocalizer = NewMessageLocalizer(ui.LocalizationManager)
	return e, nil
}

type ConfirmationDialogState struct {
	BaseState
	message   string
	nextState State
}

func (d *ConfirmationDialogState) Id() string {
	return "confirmation_dialog"
}
func (d *ConfirmationDialogState) Display(_ *GameContext, ui *UiContext) {
	ui.Console.Write(fmt.Sprintf("%s\r\n", ui.GetLocalizedMsg(ui.AppLocalizer, d.message)))
	ui.Console.Write("> ")
}
func (d *ConfirmationDialogState) Handle(ctx *GameContext, ui *UiContext, input string) (State, error) {
	ui.Msg = ui.GetLocalizedStateMsg(d, "confirmation_prompt")
	return d, nil
}
func (d *ConfirmationDialogState) GetCommands() []Command {
	return []Command{
		&ConfirmCommand{},
		&CancelCommand{},
	}
}
func NewConfirmationDialog(nextState State, message string) *ConfirmationDialogState {
	if message == "" {
		message = "confirmation_default"
	}
	return &ConfirmationDialogState{
		nextState: nextState,
		message:   message,
	}
}

type MenuOption struct {
	Id          int
	Description string
	Params      func() map[string]any
	NextState   func() State
}
type MenuState struct {
	BaseState
	ParentState State
	Options     []MenuOption
	OptionsMap  map[int]MenuOption
	Greeting    string
}

func (m *MenuState) Id() string {
	return "menu"
}
func (m *MenuState) Display(ctx *GameContext, ui *UiContext) {
	m.ShowGreeting(ctx, ui)
	for _, option := range m.Options {
		desc := ui.GetLocalizedStateMsg(m.ParentState, option.Description)
		if option.Params != nil {
			desc = SubstituteParams(desc, option.Params())
		}
		ui.Console.Write(fmt.Sprintf("%d. %s\r\n", option.Id, desc))
	}
	ui.Console.Write(ui.GetLocalizedStateMsg(m, "make_your_choice") + "\r\n")
	ui.Console.Write("> ")
}
func (m *MenuState) Handle(ctx *GameContext, ui *UiContext, input string) (State, error) {
	num, err := ui.Validator.ParseInt(input)
	if err != nil {
		return m, err
	}
	option, exists := m.OptionsMap[num]
	if !exists {
		ui.Msg = ui.GetLocalizedStateMsg(m, "invalid_option") + "\r\n"
		return m, nil
	}
	return option.NextState(), nil
}
func (m *MenuState) ShowGreeting(ctx *GameContext, ui *UiContext) {
	if m.Greeting != "" {
		ui.Console.Write(fmt.Sprintf("%s\r\n", m.Greeting))
		m.Greeting = ""
	}
}
func NewMenu(parentState State, options []MenuOption, greeting string) *MenuState {
	sort.Slice(options, func(i, j int) bool {
		return options[i].Id < options[j].Id
	})
	optionsMap := make(map[int]MenuOption, len(options))
	for _, option := range options {
		optionsMap[option.Id] = option
	}
	return &MenuState{
		ParentState: parentState,
		Options:     options,
		OptionsMap:  optionsMap,
		Greeting:    greeting,
	}
}
