package core

import (
	"fmt"
)

type MessageTranslation map[string]string

type MessageTranslations map[string]MessageTranslation

func (m MessageTranslations) isLocalized(langs []string) error {
	for msgKey, msgTrans := range m {
		for _, supportedLang := range langs {
			if _, exists := msgTrans[supportedLang]; exists {
				continue
			}
			return NewAppError(Err, "key_not_found", map[string]any{
				"key": fmt.Sprintf("%s.%s", msgKey, supportedLang),
			})
		}
	}
	return nil
}

type OptionalMessageTranslations map[string]MessageTranslations

type MessageLocalizationData struct {
	Meta         LocalizationMetadata `json:"meta"`
	Translations MessageTranslations  `json:"translations"`
}

type OptionalMessageLocalizationData struct {
	Meta         LocalizationMetadata        `json:"meta"`
	Translations OptionalMessageTranslations `json:"translations"`
}

type MessageLocalizer struct {
	lm                   *LocalizationManager
	Translations         MessageTranslations
	OptionalTranslations map[string]MessageTranslations
}

func NewMessageLocalizer(lm *LocalizationManager) *MessageLocalizer {
	return &MessageLocalizer{
		lm:                   lm,
		Translations:         make(MessageTranslations),
		OptionalTranslations: make(map[string]MessageTranslations),
	}
}

func (l *MessageLocalizer) LoadTranslations(filePath string) error {
	var rawData MessageLocalizationData
	if err := l.lm.loadLocalizationData(filePath, &rawData); err != nil {
		return err
	}
	supportedLanguages := rawData.Meta.SupportedLanguages
	err := rawData.Translations.isLocalized(supportedLanguages)
	if err != nil {
		locErr := NewAppError(ErrLocalization, "localization_file_translations_error", map[string]any{
			"file":  filePath,
			"path":  "translations",
			"error": err,
		})
		if !l.lm.isCoreLocalization(filePath) {
			return locErr
		}
		l.lm.logError(locErr)
	}
	l.CopyTranslations(l.Translations, rawData.Translations)
	if l.lm.isCoreLocalization(filePath) {
		l.lm.updateAvailableLanguages(supportedLanguages)
	}
	return nil
}

func (l *MessageLocalizer) LoadOptionalTranslations(filePath string) error {
	var rawData OptionalMessageLocalizationData
	if err := l.lm.loadLocalizationData(filePath, &rawData); err != nil {
		return err
	}
	supportedLanguages := rawData.Meta.SupportedLanguages
	for setName, set := range rawData.Translations {
		if _, exists := l.OptionalTranslations[setName]; !exists {
			l.OptionalTranslations[setName] = make(MessageTranslations)
		}
		err := set.isLocalized(supportedLanguages)
		if err != nil {
			locErr := NewAppError(ErrLocalization, "localization_file_translations_error", map[string]any{
				"file":  filePath,
				"path":  fmt.Sprintf("translations.%s", setName),
				"error": err,
			})
			if !l.lm.isCoreLocalization(filePath) {
				return locErr
			}
			l.lm.logError(locErr)
		}
		l.CopyTranslations(l.OptionalTranslations[setName], set)
	}
	if l.lm.isCoreLocalization(filePath) {
		l.lm.updateAvailableLanguages(supportedLanguages)
	}
	return nil
}

func (l *MessageLocalizer) Get(key string) (string, error) {
	message, messageExists := l.Translations[key]
	if !messageExists {
		return "", NewAppError(ErrLocalization, "key_not_found", map[string]any{
			"key": key,
		})
	}
	return fetchTranslation(l.lm, message)
}

func (l *MessageLocalizer) GetOptional(messageSet, key string) (string, error) {
	set, exists := l.OptionalTranslations[messageSet]
	if !exists {
		return "", NewAppError(ErrLocalization, "set_not_found", map[string]any{
			"set": set,
		})
	}
	message, messageExists := set[key]
	if !messageExists {
		return "", NewAppError(ErrLocalization, "key_in_set_not_found", map[string]any{
			"key": key,
			"set": set,
		})
	}
	return fetchTranslation(l.lm, message)
}

func (l *MessageLocalizer) CopyTranslations(dest, source MessageTranslations) {
	for msgKey, trans := range source {
		if _, exists := dest[msgKey]; !exists {
			dest[msgKey] = make(map[string]string)
		}
		for lang, msg := range trans {
			dest[msgKey][lang] = msg
		}
	}
}
