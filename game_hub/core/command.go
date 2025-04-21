package core

type Command interface {
	Execute(ctx *AppContext, ui *UiContext, args []string) (State, error)
	Id() string
	Scope() Scope
}

type BaseCommand struct{}

func (c *BaseCommand) Execute(ctx *AppContext, ui *UiContext, args []string) (State, error) {
	ui.DisplayText(ui.GetLocalizedMsg(ui.AppLocalizer, "unknown_command_action"))
	return ctx.GetCurrentState()
}
func (c *BaseCommand) Id() string {
	return "unknown"
}
func (c *BaseCommand) Scope() Scope {
	return ScopeCore
}

type GameCommand struct{ BaseCommand }

func (c *GameCommand) Scope() Scope {
	return ScopeGame
}
func DefaultGlobalCommands() []Command {
	return []Command{
		&HelpCommand{},
		&QuitCommand{},
		&VersionCommand{},
	}
}
