# Game Hub

[English](#game-hub) | [Русский](README_ru.md)

Game Hub is a console-based Go application that allows users to play various mini-games such as "Guess the Number" and "Rock, Paper, Scissors." The project supports localization (English and Russian), error handling, and a modular architecture for easily adding new games.

## Features

- Multilingual support (English, Russian).
- Modular structure for adding new games.
- State and command management via a finite state machine.
- Robust error handling and input validation.
- Available games:
  - Guess the Number
  - Rock, Paper, Scissors

## Requirements

- Go 1.24 or higher

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/Daniil-Gusev/game_hub.git
   cd game_hub
   ```

2. Ensure Go is installed:
   ```bash
   go version
   ```

3. Run the application:
   ```bash
   go run .
   ```

## Usage

1. Upon launch, you enter the main menu where you can select a game or exit.
2. Enter the number of the desired option to choose a game.
3. Follow the in-game instructions.
4. Use commands like `help`, `quit`, or `back` for navigation.

## Adding a New Game

1. Copy the `games/game_template` folder and rename it to a unique name for your game (e.g., `mygame`).
2. Edit the files in the new folder:
   - `interface.go`: Set a unique `GetId()`.
   - `game.go`: Implement the game logic.
   - `states.go` and `states.json`: Define game states and their localization.
   - `translations.json`: Add translations.
   - (Optional) `commands.go` and `commands.json`: Implement custom commands.
3. Register the game in `games/games.go` by updating the `AvailableGames()` function:
   ```go
   import "game_hub/games/mygame"
   ...
   return []core.GameInterface{
       guessnumber.NewGame(),
       rockpaperscissors.NewGame(),
       mygame.NewGame(),
   }
   ```
4. Verify that the game appears in the main menu.

## Project Structure

- `core/`: Core application logic (state management, commands, localization).
- `app/`: Main menu logic.
- `games/`: Folder containing games.
  - `game_template/`: Template for new games.
- `main.go`: Entry point.

## Localization

The project supports English and Russian languages. Localization is implemented via JSON files in the `core`, `app`, and `games` folders. To add a new language, include translations in the relevant JSON files with a new language key (e.g., `"fr": "Bonjour"`).

## License

MIT License

## Authors

- Daniil Gusev