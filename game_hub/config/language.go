package config

type LanguageConfig struct {
	CurrentLanguage string
	DefaultLanguage string
}

func NewLanguageConfig() *LanguageConfig {
	return &LanguageConfig{
		CurrentLanguage: "en",
		DefaultLanguage: "en",
	}
}
