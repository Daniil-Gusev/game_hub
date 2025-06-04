package core

type Scope string

const (
	ScopeCore Scope = "core"
	ScopeApp  Scope = "app"
	ScopeGame Scope = "game"
)

func (s Scope) IsValid() bool {
	switch s {
	case ScopeCore, ScopeApp, ScopeGame:
		return true
	default:
		return false
	}
}

type void struct{}
