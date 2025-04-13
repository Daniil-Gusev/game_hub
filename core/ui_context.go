package core

import (
	"fmt"
	"strings"
)

type UiContext struct {
	Console             Console
	Msg                 string
	Validator           InputValidator
	ErrHandler          ErrorHandler
	LocalizationManager *LocalizationManager
	CommandRegistry     *CommandRegistry
	AppLocalizer        *MessageLocalizer
	GameLocalizer       *MessageLocalizer
	StateLocalizer      *StateLocalizer
}

func (ui *UiContext) DisplayError(err error) {
	msg := ui.ErrHandler.Handle(err)
	if msg != "" {
		ui.Console.Write(fmt.Sprintf("%s\r\n", msg))
	}
}
func (ui *UiContext) DisplayMessage() {
	if ui.Msg != "" {
		ui.Console.Write(fmt.Sprintf("%s", ui.Msg))
		ui.Msg = ""
	}
}
func (ui *UiContext) HandleInput(input string, ctx *GameContext) (State, error) {
	input = strings.TrimSpace(input)
	if cmd, args := ui.CommandRegistry.ParseInput(input); cmd != nil {
		return cmd.Execute(ctx, ui, args)
	}
	state, err := ctx.GetCurrentState()
	if err != nil {
		return nil, err
	}
	return state.Handle(ctx, ui, input)
}

func (ui *UiContext) GetLocalizedMsg(localizer *MessageLocalizer, key string) string {
	msg, err := localizer.Get(key)
	if err != nil {
		ui.DisplayError(err)
	}
	return msg
}
func (ui *UiContext) GetOptionalLocalizedMsg(localizer *MessageLocalizer, set string, key string) string {
	msg, err := localizer.GetOptional(set, key)
	if err != nil {
		ui.DisplayError(err)
	}
	return msg
}

func (ui *UiContext) GetLocalizedCmdName(cmd Command) string {
	name, err := ui.CommandRegistry.GetName(cmd)
	if err != nil {
		ui.DisplayError(err)
	}
	return name
}
func (ui *UiContext) GetLocalizedCmdDescription(cmd Command) string {
	desc, err := ui.CommandRegistry.GetDescription(cmd)
	if err != nil {
		ui.DisplayError(err)
	}
	return desc
}
func (ui *UiContext) GetLocalizedCmdAliases(cmd Command) []string {
	aliases, err := ui.CommandRegistry.GetAliases(cmd)
	if err != nil {
		ui.DisplayError(err)
	}
	return aliases
}

func (ui *UiContext) GetLocalizedStateDescription(state State) string {
	desc, err := ui.StateLocalizer.GetDescription(state.Scope(), state.Id())
	if err != nil {
		ui.DisplayError(err)
	}
	return desc
}
func (ui *UiContext) GetLocalizedStateMsg(state State, key string) string {
	msg, err := ui.StateLocalizer.GetMessage(state.Scope(), state.Id(), key)
	if err != nil {
		ui.DisplayError(err)
	}
	return msg
}
