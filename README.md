# Game Hub

[English](#game-hub) | [Русский](README_ru.md)

Game Hub is a console-based Go application designed for playing a variety of mini-games, such as "Guess the Number" and "Rock, Paper, Scissors." The project emphasizes modularity, multilingual support, robust error handling, and cross-platform compatibility, making it easy to extend with new games and deploy on multiple operating systems.

## Features

- **Multilingual Support**: English and Russian localization for game interfaces and messages.
- **Modular Architecture**: Easily add new games using the provided `game_template`.
- **Finite State Machine**: Manages game states and commands for a seamless user experience.
- **Robust Error Handling**: Comprehensive input validation and localized error messages.
- **Cross-Platform Support**: Builds for Linux, Windows, and macOS (amd64 and arm64 architectures).
- **Installer and Portable Builds**: Options for default installation or portable usage.
- **Available Games**:
  - Guess the Number
  - Rock, Paper, Scissors
  - Game Template (for developers to create new games)

## Requirements

- **Go**: Version 1.24 or higher
- **Optional**: `7z` for faster archive compression during release builds (falls back to `zip` if unavailable)

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/Daniil-Gusev/game_hub.git
   cd game_hub
   ```

2. Verify Go installation:
   ```bash
   go version
   ```

3. Run the application:
   ```bash
   cp -r data game_hub/
   go run .
   ```

4. (Optional) Build a release:
   - For portable binaries:
     ```bash
     ./build_portable_release.sh vX.Y.Z
     ```
   - For installers:
     ```bash
     ./build_installable_release.sh vX.Y.Z
     ```
   Replace `vX.Y.Z` with the desired version (e.g., `v1.0.0`). Outputs are generated in the `release/` directory.

## Usage

1. Launch the application to access the main menu, where you can select a game or exit.
2. Enter the number corresponding to your choice (e.g., `1` for Guess the Number, `0` to exit).
3. Follow the in-game instructions, which are displayed in your configured language.
4. Use commands like `help`, `quit`, `back`, or game-specific commands (e.g., `restart`) for navigation.

## Adding a New Game

To create a new game, follow these steps:

1. Copy the `games/game_template` folder and rename it to a unique name (e.g., `mygame`).
2. Modify the files in the new folder:
   - **`interface.go`**: Update `GetId()` to return a unique identifier (e.g., `"mygame"`).
   - **`game.go`**: Implement the core game logic.
   - **`states.go` and `states.json`**: Define game states and their localized descriptions/messages.
   - **`translations.json`**: Add translations for game-specific messages.
   - **(Optional) `commands.go` and `commands.json`**: Define and localize custom commands.
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
4. Test the game by running the application and verifying it appears in the main menu.

## Project Structure

- game_hub
 - **`core/`**: Core application logic, including state management, command registry, localization, and error handling.
 - **`app/`**: Main menu and application-level states.
 - **`games/`**: Contains individual games.
  - **`game_template/`**: Template for creating new games.
   - **`guessnumber/`**: Implementation of the Guess the Number game.
  - **`rockpaperscissors/`**: Implementation of the Rock, Paper, Scissors game.
 - **`config/`**: Configuration and path management.
 - **`utils/`**: Shared utilities (e.g., text wrapping, parameter substitution).
 - **`main.go`**: Application entry point.
- **`installer/`**: Scripts and logic for generating installers and wrappers.
- **`build_portable_release.sh`**: Script for building portable binaries.
- **`build_installable_release.sh`**: Script for building installers.

## Localization

Game Hub supports English and Russian through JSON-based localization files located in `core/`, `app/`, and `games/` of the `data/` folder. To add a new language:

1. Add translations to the relevant JSON files (e.g., `core/translations.json`, `games/guessnumber/translations.json`) with a new language key (e.g., `"fr": "Bonjour"`).
2. Test the new language by setting it in the application configuration or passing it as a parameter.

## Building Releases

- **Portable Builds**: Use `build_portable_release.sh` to create standalone binaries with data files, archived as `.tar.gz` (Linux/macOS) or `.7z`/`.zip` (Windows).
- **Installers**: Use `build_installable_release.sh` to create installers with platform-specific wrappers (e.g., `.app` bundles for macOS, `.desktop` files for Linux, and shortcuts for Windows).

Both scripts support multiple platforms (Linux, Windows, macOS) and architectures (amd64, arm64). Ensure `7z` is installed for optimal compression on Windows.

## License

MIT License

## Authors

- Daniil Gusev

## Acknowledgments

- Built with [Go](https://golang.org/) for cross-platform compatibility.
- Inspired by classic console-based games and modular application design.