package core

// GameInterface определяет общий интерфейс для игр
type GameInterface interface {
	CreateNew() GameInterface
	GetId() string
	GetStartState() State
}

type GameContext struct {
	Game         GameInterface
	StateStack   *StateStack
	AppIsRunning bool
}

func (gc *GameContext) GetCurrentState() (State, error) {
	if gc.StateStack.IsEmpty() {
		return nil, NewAppError(ErrStateStack, "state_stack_empty", nil)
	}
	return gc.StateStack.Peek(), nil
}
func (gc *GameContext) GetPreviousState() (State, error) {
	if gc.StateStack.IsEmpty() {
		return nil, NewAppError(ErrStateStack, "state_stack_empty", nil)
	}
	if len(gc.StateStack.states) < 2 {
		return nil, NewAppError(ErrStateStack, "state_stack_insufficient", nil)
	}
	gc.StateStack.Pop()
	return gc.StateStack.Pop(), nil
}
