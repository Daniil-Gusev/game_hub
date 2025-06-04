package core

import (
	"game_hub/config"
)

type AppContext struct {
	Config         *config.Config
	Game           GameInterface
	AvailableGames []GameInterface
	StateStack     *StateStack
	AppIsRunning   bool
	GoToMenu       bool
}

func (app *AppContext) GetCurrentState() (State, error) {
	if app.StateStack.IsEmpty() {
		return nil, NewAppError(ErrStateStack, "state_stack_empty", nil)
	}
	return app.StateStack.Peek(), nil
}
func (app *AppContext) GetPreviousState() (State, error) {
	if app.StateStack.IsEmpty() {
		return nil, NewAppError(ErrStateStack, "state_stack_empty", nil)
	}
	if len(app.StateStack.states) < 2 {
		return nil, NewAppError(ErrStateStack, "state_stack_insufficient", nil)
	}
	app.StateStack.Pop()
	return app.StateStack.Pop(), nil
}

func (app *AppContext) GoToState(nextState State, ui *UiContext) (State, error) {
	app.StateStack.Push(nextState)
	if newState, err := nextState.Init(app, ui); err != nil {
		if newState != nextState {
			app.StateStack.Pop()
			app.StateStack.Push(newState)
		}
		return newState, err
	}
	if err := ui.CommandRegistry.RegisterLocalCommands(nextState.GetCommands()); err != nil {
		return nextState, err
	}
	return nextState, nil
}
