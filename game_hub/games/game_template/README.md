# Game Template

This template is designed for creating a new game in the Game Hub project.

## How to Use

1. Rename the `game_template` folder to a unique name for your game (e.g., `mygame`).
2. Update `GetId()` in `interface.go` to return a unique identifier for your game.
3. Implement the game logic in `game.go`.
4. Define game states in `states.go` and their localization in `states.json`.
5. Add translations to `translations.json`.
6. (Optional) Define commands in `commands.json` and implement them in `commands.go`.
7. Register the game in `games/games.go` by adding a call to `NewGame()` in the `AvailableGames()` function:
   ```go
   import "game_hub/games/mygame"
   ...
   return []core.GameInterface{
       guessnumber.NewGame(),
       rockpaperscissors.NewGame(),
       mygame.NewGame(),
   }
   ```

## File Structure

- `interface.go`: Implements the `core.GameInterface`.
- `game.go`: Contains the game logic.
- `states.go`: Defines game states.
- `states.json`: Localizes game states.
- `translations.json`: Contains translations that are common to states or commands.
- `commands.json`: (Optional) Localizes commands.
- `commands.go`: (Optional) Implements commands.