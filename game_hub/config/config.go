package config

// Config contains all application configuration settings.
type Config struct {
	Paths    *PathConfig
	Language *LanguageConfig
}

// NewConfig creates a new Config instance with initialized PathConfig and LanguageConfig.
func NewConfig(appName string) (*Config, error) {
	pathConfig, err := NewPathConfig(appName)
	if err != nil {
		return nil, err
	}
	languageConfig := NewLanguageConfig()
	return &Config{
		Paths:    pathConfig,
		Language: languageConfig,
	}, nil
}
