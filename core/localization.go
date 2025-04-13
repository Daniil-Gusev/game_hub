package core

import (
	"encoding/json"
	"os"
)

type LocalizationManager struct {
	currentLang string
	defaultLang string
}

func NewLocalizationManager(lang string) *LocalizationManager {
	return &LocalizationManager{
		currentLang: lang,
		defaultLang: lang,
	}
}

func (lm *LocalizationManager) SetLanguage(lang string) {
	lm.currentLang = lang
}

func (lm *LocalizationManager) CurrentLang() string {
	return lm.currentLang
}

func (lm *LocalizationManager) DefaultLang() string {
	return lm.defaultLang
}

// loadRawData загружает переводы из файла в указанную структуру.
func (lm *LocalizationManager) loadRawData(filePath string, target any) error {
	file, err := os.Open(filePath)
	if err != nil {
		return NewAppError(ErrLocalization, "file_open_error", map[string]any{
			"file":  filePath,
			"error": err,
		})
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(target); err != nil {
		return NewAppError(ErrLocalization, "file_parse_error", map[string]any{
			"file":  filePath,
			"error": err,
		})
	}
	return nil
}

// fetchTranslation получает перевод для указанного языка с fallback на язык по умолчанию.
func fetchTranslation[T any](lm *LocalizationManager, translations map[string]T) (T, error) {
	if value, exists := translations[lm.currentLang]; exists {
		return value, nil
	}
	if value, exists := translations[lm.defaultLang]; exists {
		return value, nil
	}
	var zero T
	return zero, NewAppError(ErrLocalization, "lang_not_supported", map[string]any{
		"lang": lm.currentLang,
	})
}
