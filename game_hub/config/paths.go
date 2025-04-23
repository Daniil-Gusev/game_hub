package config

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

// PathConfig manages paths to configuration files.
type PathConfig struct {
	baseDir    string
	gamesDir   string
	isPortable bool
}

// NewPathConfig creates a new PathConfig instance and validates the configuration directory.
func NewPathConfig(appName string) (*PathConfig, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	exeDir := filepath.Dir(exePath)

	dataDir := filepath.Join(exeDir, "data")
	if _, err := os.Stat(dataDir); err == nil {
		// Portable mode: use game_hub_data as baseDir
		gamesDir := filepath.Join(dataDir, "games")
		if err := os.MkdirAll(gamesDir, 0755); err != nil {
			return nil, err
		}
		return &PathConfig{
			baseDir:    dataDir,
			gamesDir:   gamesDir,
			isPortable: true,
		}, nil
	}

	// Installed mode: use system configuration directory
	configDir, err := OsConfigDir(runtime.GOOS)
	if err != nil {
		return nil, err
	}
	baseDir := filepath.Join(configDir, appName)
	gamesDir := filepath.Join(baseDir, "games")

	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		return nil, errors.New("configuration directory does not exist: " + baseDir)
	}

	dir, err := os.Open(baseDir)
	if err != nil {
		return nil, err
	}
	defer dir.Close()
	_, err = dir.Readdirnames(1) // Check for at least one entry
	if err != nil {
		return nil, errors.New("configuration directory is empty: " + baseDir)
	}

	if err := os.MkdirAll(gamesDir, 0755); err != nil {
		return nil, err
	}

	return &PathConfig{
		baseDir:    baseDir,
		gamesDir:   gamesDir,
		isPortable: false,
	}, nil
}

func (pc *PathConfig) CoreTranslationsPath() string {
	return filepath.Join(pc.baseDir, "core", "translations.json")
}

// CoreStatesPath returns the path to states.json in core.
func (pc *PathConfig) CoreStatesPath() string {
	return filepath.Join(pc.baseDir, "core", "states.json")
}

// CoreGlobalCommandsPath returns the path to global_commands.json in core.
func (pc *PathConfig) CoreGlobalCommandsPath() string {
	return filepath.Join(pc.baseDir, "core", "global_commands.json")
}

// CoreLocalCommandsPath returns the path to local_commands.json in core.
func (pc *PathConfig) CoreLocalCommandsPath() string {
	return filepath.Join(pc.baseDir, "core", "local_commands.json")
}

// AppTranslationsPath returns the path to translations.json in app.
func (pc *PathConfig) AppTranslationsPath() string {
	return filepath.Join(pc.baseDir, "app", "translations.json")
}

// AppStatesPath returns the path to states.json in app.
func (pc *PathConfig) AppStatesPath() string {
	return filepath.Join(pc.baseDir, "app", "states.json")
}

// GamesTranslationsPath returns the path to translations.json for games.
func (pc *PathConfig) GamesTranslationsPath() string {
	return filepath.Join(pc.gamesDir, "translations.json")
}

// GameStatesPath returns the path to states.json for a specific game.
func (pc *PathConfig) GameStatesPath(gameID string) string {
	return filepath.Join(pc.gamesDir, gameID, "states.json")
}

// GameCommandsPath returns the path to commands.json for a specific game.
func (pc *PathConfig) GameCommandsPath(gameID string) string {
	return filepath.Join(pc.gamesDir, gameID, "commands.json")
}

// GameTranslationsPath returns the path to translations.json for a specific game.
func (pc *PathConfig) GameTranslationsPath(gameID string) string {
	return filepath.Join(pc.gamesDir, gameID, "translations.json")
}
func OsConfigDir(platform string) (string, error) {
	switch platform {
	case "linux":
		return "/usr/local/share", nil
	case "windows":
		return "C:\\ProgramData", nil
	case "darwin":
		return "/Library/Application Support", nil
	default:
		return "", errors.New("Platform " + platform + " is not supported.")
	}
}
