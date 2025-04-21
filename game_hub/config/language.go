package config

// LanguageConfig manages language settings.
type LanguageConfig struct {
	CurrentLanguage string
}

// NewLanguageConfig creates a new LanguageConfig instance with a default language.
func NewLanguageConfig() *LanguageConfig {
	// Можно сделать настраиваемым, например, через переменную окружения
	return &LanguageConfig{
		CurrentLanguage: "ru", // По умолчанию русский
	}
}
