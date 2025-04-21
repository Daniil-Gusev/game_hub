package core

type StateStack struct {
	states []State
}

func NewStateStack() *StateStack {
	return &StateStack{states: make([]State, 0)}
}
func (s *StateStack) IsEmpty() bool {
	return len(s.states) == 0
}
func (s *StateStack) Push(state State) {
	s.states = append(s.states, state)
}
func (s *StateStack) Pop() State {
	if s.IsEmpty() {
		return nil
	}
	lastIndex := len(s.states) - 1
	state := s.states[lastIndex]
	s.states = s.states[:lastIndex]
	return state
}
func (s *StateStack) Peek() State {
	if s.IsEmpty() {
		return nil
	}
	lastIndex := len(s.states) - 1
	return s.states[lastIndex]
}
func (s *StateStack) Clear() {
	s.states = s.states[:0]
}
