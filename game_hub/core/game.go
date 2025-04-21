package core

// GameInterface определяет общий интерфейс для игр
type GameInterface interface {
	CreateNew() GameInterface
	GetId() string
	GetStartState() State
}
