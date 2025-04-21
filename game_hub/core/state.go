package core

type State interface {
	Init(ctx *AppContext, ui *UiContext) (State, error)
	Handle(ctx *AppContext, ui *UiContext, input string) (State, error)
	Display(ctx *AppContext, ui *UiContext)
	GetCommands() []Command
	RequiresInput() bool
	Id() string
	Scope() Scope
}

type BaseState struct{}

func (b *BaseState) Id() string {
	return "unknown"
}
func (b *BaseState) GetCommands() []Command {
	return []Command{}
}
func (b *BaseState) RequiresInput() bool {
	return true
}
func (b *BaseState) Scope() Scope {
	return ScopeCore
}
func (b *BaseState) Init(ctx *AppContext, ui *UiContext) (State, error) {
	return b, nil
}
func (b *BaseState) Display(ctx *AppContext, ui *UiContext) {}
func (b *BaseState) Handle(ctx *AppContext, ui *UiContext, input string) (State, error) {
	return b, nil
}
