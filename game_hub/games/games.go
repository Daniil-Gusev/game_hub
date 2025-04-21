package games

import (
	"game_hub/core"
	"game_hub/games/guessnumber"
	"game_hub/games/rockpaperscissors"
)

func AvailableGames() []core.GameInterface {
	return []core.GameInterface{
		guessnumber.NewGame(),
		rockpaperscissors.NewGame(),
	}
}
