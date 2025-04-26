package core

import (
	"fmt"
	"strings"
	"time"
)

type QuitCommand struct{ BaseCommand }

func (c *QuitCommand) Id() string {
	return "quit"
}
func (c *QuitCommand) Execute(ctx *AppContext, ui *UiContext, args []string) (State, error) {
	if len(args) > 1 && args[1] == "force" {
		return &ExitState{}, nil
	}
	state, err := ctx.GetCurrentState()
	if err != nil {
		return state, nil
	}
	if _, ok := state.(*ConfirmationDialogState); ok {
		return state, nil
	}
	return NewConfirmationDialog(&ExitState{}, "quit_confirm"), nil
}

type HelpCommand struct{ BaseCommand }

func (c *HelpCommand) Id() string {
	return "help"
}
func (c *HelpCommand) Execute(ctx *AppContext, ui *UiContext, args []string) (State, error) {
	state, err := ctx.GetCurrentState()
	if err != nil {
		return nil, err
	}
	desc := ui.GetLocalizedStateDescription(state)
	if desc == "" {
		ui.DisplayText(ui.GetLocalizedMsg(ui.AppLocalizer, "help_not_found") + "\r\n")
	} else {
		ui.DisplayText(fmt.Sprintf("%s\r\n", desc))
	}
	ui.DisplayText(ui.GetLocalizedMsg(ui.AppLocalizer, "available_commands") + "\r\n")
	for _, cmd := range ui.CommandRegistry.GetLocalCommands() {
		ui.DisplayText(fmt.Sprintf("%s: (%s).\r\n%s\r\n", ui.GetLocalizedCmdName(cmd), strings.Join(ui.GetLocalizedCmdAliases(cmd), ", "), ui.GetLocalizedCmdDescription(cmd)))
	}
	for _, cmd := range ui.CommandRegistry.GetGlobalCommands() {
		ui.DisplayText(fmt.Sprintf("%s: (%s).\r\n%s\r\n", ui.GetLocalizedCmdName(cmd), strings.Join(ui.GetLocalizedCmdAliases(cmd), ", "), ui.GetLocalizedCmdDescription(cmd)))
	}
	return state, nil
}

type BackCommand struct{ BaseCommand }

func (c *BackCommand) Id() string {
	return "back"
}
func (c *BackCommand) Execute(ctx *AppContext, ui *UiContext, args []string) (State, error) {
	state, err := ctx.GetPreviousState()
	if err != nil {
		return nil, err
	}
	return state, nil
}

type ExitCommand struct{ BaseCommand }

func (c *ExitCommand) Id() string {
	return "exit"
}
func (c *ExitCommand) Execute(ctx *AppContext, ui *UiContext, args []string) (State, error) {
	return ctx.Game.GetStartState(), nil
}

type VersionCommand struct{ BaseCommand }

func (c *VersionCommand) Id() string {
	return "version"
}
func (c *VersionCommand) Execute(ctx *AppContext, ui *UiContext, args []string) (State, error) {
	versionMsg := ui.GetLocalizedMsg(ui.AppLocalizer, "version_info")
	var displayTime string
	builtTime, err := time.Parse(time.RFC3339, BuildTime)
	if err != nil {
		displayTime = BuildTime
	} else {
		displayTime = builtTime.Format("02.01.2006 15:04:05")
	}
	versionMsg = fmt.Sprintf(versionMsg, Version, displayTime)
	ui.DisplayText(versionMsg + "\r\n")
	return ctx.GetCurrentState()
}

type ConfirmCommand struct{ BaseCommand }

func (c *ConfirmCommand) Id() string {
	return "confirm"
}
func (c *ConfirmCommand) Execute(ctx *AppContext, ui *UiContext, args []string) (State, error) {
	currentState, err := ctx.GetCurrentState()
	if err != nil {
		return currentState, err
	}
	confirmationState, ok := currentState.(*ConfirmationDialogState)
	if !ok {
		previousState, _ := ctx.GetPreviousState()
		return previousState, NewAppError(ErrInternal, "некорректный диалог подтверждения.", nil)
	}
	return confirmationState.nextState, nil
}

type CancelCommand struct{ BaseCommand }

func (c *CancelCommand) Id() string {
	return "cancel"
}
func (c *CancelCommand) Execute(ctx *AppContext, ui *UiContext, args []string) (State, error) {
	return ctx.GetPreviousState()
}
