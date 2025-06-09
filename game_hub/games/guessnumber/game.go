package guessnumber

import (
	"game_hub/core"
	"math"
)

type Difficulty int

const (
	VeryEasy Difficulty = 1 + iota
	Easy
	Medium
	Hard
	VeryHard
)

func (d Difficulty) String() string {
	switch d {
	case VeryEasy:
		return "very_easy"
	case Easy:
		return "easy"
	case Medium:
		return "medium"
	case Hard:
		return "hard"
	case VeryHard:
		return "very_hard"
	default:
		return "unknown"
	}
}

func (d Difficulty) GetModifier() int {
	switch d {
	case VeryEasy:
		return 10
	case Easy:
		return 8
	case Medium:
		return 4
	case Hard:
		return 2
	case VeryHard:
		return 0
	default:
		return 0
	}
}

type Game struct {
	// минимальный диапазон угадывания чисел
	MinRangeSize int
	// минимально возможное число диапазона угадывания чисел
	MinRangeNumber int
	// максимально возможное число диапазона угадывания чисел
	MaxRangeNumber  int
	Difficulty      Difficulty
	MinNumber       int
	MaxNumber       int
	secretNumber    int
	attempts        int
	isWon           bool
	RandomGenerator *core.RandomGenerator
}

func NewGame() *Game {
	return &Game{
		MinRangeSize:    20,
		MinRangeNumber:  0,
		MaxRangeNumber:  math.MaxInt32,
		MinNumber:       1,
		MaxNumber:       100,
		Difficulty:      Medium,
		isWon:           false,
		RandomGenerator: core.NewRandomGenerator(),
	}
}

func (g *Game) Prepare() error {
	g.isWon = false
	secret, err := g.RandomGenerator.Generate(g.MinNumber, g.MaxNumber)
	if err != nil {
		return err
	}
	g.secretNumber = secret
	attempts, err := g.CalculateAttempts()
	if err != nil {
		return err
	}
	g.attempts = attempts
	return nil
}

// количество попыток на угадывание, как округленный до целого логарифм 2 плюс рандом от 0 до значения модификатора сложности
func (g *Game) CalculateAttempts() (int, error) {
	rangeSize := g.MaxNumber - g.MinNumber + 1
	if rangeSize < 1 || rangeSize > g.MaxRangeNumber {
		return 0, core.NewAppError(core.ErrInvalidRange, "Некорректный диапазон чисел для расчёта попыток.", nil)
	}
	minAttempts := int(math.Round(math.Log2(float64(rangeSize))))
	mod, err := g.RandomGenerator.Generate(0, g.Difficulty.GetModifier())
	if err != nil {
		return 0, err
	}
	return minAttempts + mod, nil
}

func (g *Game) GetAttempts() int {
	return g.attempts
}

func (g *Game) CheckWin() bool {
	return g.isWon
}

func (g *Game) CheckLoss() bool {
	return g.attempts < 1
}

func (g *Game) MakeGuess(guess int) {
	g.attempts--
	if guess == g.secretNumber {
		g.isWon = true
	}
}

func (g *Game) GetHint(guess int) string {
	if guess < g.secretNumber {
		return "hint_bigger"
	} else if guess > g.secretNumber {
		return "hint_smaller"
	}
	return "you_guessed"
}
