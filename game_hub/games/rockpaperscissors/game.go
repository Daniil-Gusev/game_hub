package rockpaperscissors

import (
	"game_hub/core"
)

type Move int

const (
	Rock Move = iota
	Scissors
	Paper
)

func (m Move) String() string {
	switch m {
	case Rock:
		return "rock"
	case Scissors:
		return "scissors"
	case Paper:
		return "paper"
	default:
		return "unknown"
	}
}

type RoundResult int

const (
	Winning RoundResult = iota
	Draw
	Loss
)

type Game struct {
	MinRounds       int
	MaxRounds       int
	PlayerScore     int
	BotScore        int
	CurrentRound    int
	TotalRounds     int
	PlayerMove      Move
	BotMove         Move
	winTable        [3][3]RoundResult
	isWon           bool
	isLoss          bool
	RandomGenerator *core.RandomGenerator
}

func NewGame() *Game {
	return &Game{
		MinRounds:   3,
		MaxRounds:   100,
		PlayerScore: 0,
		BotScore:    0,
		TotalRounds: 3,
		winTable: [3][3]RoundResult{
			{Draw, Winning, Loss},
			{Loss, Draw, Winning},
			{Winning, Loss, Draw},
		},
		isWon:           false,
		isLoss:          false,
		RandomGenerator: core.NewRandomGenerator(),
	}
}
func (g *Game) Reset() {
	g.CurrentRound = 1
	g.PlayerScore, g.BotScore = 0, 0
	g.isWon, g.isLoss = false, false
}

func (g *Game) MakePlayerMove(playerMove Move) {
	g.PlayerMove = playerMove
}
func (g *Game) MakeBotMove() {
	var botMove Move
	rnd, _ := g.RandomGenerator.Generate(1, 3000)
	if rnd <= 1000 {
		botMove = Rock
	} else if rnd <= 2000 {
		botMove = Scissors
	} else {
		botMove = Paper
	}
	g.BotMove = botMove
}
func (g *Game) PlayRound() RoundResult {
	g.CurrentRound++
	result := g.winTable[int(g.PlayerMove)][int(g.BotMove)]
	switch result {
	case Winning:
		g.PlayerScore++
	case Loss:
		g.BotScore++
	case Draw:
		g.PlayerScore++
		g.BotScore++
	}
	if g.CurrentRound > g.TotalRounds {
		if g.PlayerScore > g.BotScore {
			g.isWon = true
		} else if g.PlayerScore < g.BotScore {
			g.isLoss = true
		}
		return result
	}
	return result
}
func (g *Game) CheckWin() bool {
	return g.isWon
}
func (g *Game) CheckLoss() bool {
	return g.isLoss
}
