package core

import (
	"fmt"
)

type StateTranslation struct {
	Description map[string]string            `json:"description"`
	Messages    map[string]map[string]string `json:"messages"`
}

func (s StateTranslation) isLocalized(langs []string) error {
	for _, supportedLang := range langs {
		if len(s.Description) == 0 {
			break
		}
		if _, exists := s.Description[supportedLang]; exists {
			continue
		}
		return NewAppError(Err, "key_not_found", map[string]any{
			"key": fmt.Sprintf("description.%s", supportedLang),
		})
	}
	for msgKey, msgTrans := range s.Messages {
		for _, supportedLang := range langs {
			if _, exists := msgTrans[supportedLang]; exists {
				continue
			}
			return NewAppError(Err, "key_not_found", map[string]any{
				"key": fmt.Sprintf("Messages.%s.%s", msgKey, supportedLang),
			})
		}
	}
	return nil
}

type StateTranslations map[Scope]map[string]StateTranslation

type StateLocalizationData struct {
	Meta         LocalizationMetadata `json:"meta"`
	Translations StateTranslations    `json:"translations"`
}

type StateLocalizer struct {
	lm           *LocalizationManager
	Translations StateTranslations
}

func NewStateLocalizer(lm *LocalizationManager) *StateLocalizer {
	return &StateLocalizer{
		lm:           lm,
		Translations: make(StateTranslations),
	}
}

func (l *StateLocalizer) LoadTranslations(filePath string) error {
	var rawData StateLocalizationData
	if err := l.lm.loadLocalizationData(filePath, &rawData); err != nil {
		return err
	}

	supportedLanguages := rawData.Meta.SupportedLanguages
	for scope, states := range rawData.Translations {
		if !scope.IsValid() {
			return NewAppError(ErrLocalization, "invalid_state_localizations_scope", map[string]any{
				"file":  filePath,
				"scope": scope,
			})
		}
		if _, exists := l.Translations[scope]; !exists {
			l.Translations[scope] = make(map[string]StateTranslation)
		}
		for stateId, trans := range states {
			if trans.Description == nil {
				trans.Description = make(map[string]string)
			}
			if trans.Messages == nil {
				trans.Messages = make(map[string]map[string]string)
			}
			err := trans.isLocalized(supportedLanguages)
			if err != nil {
				locErr := NewAppError(ErrLocalization, "localization_file_translations_error", map[string]any{
					"file":  filePath,
					"path":  fmt.Sprintf("translations.%v.%s", scope, stateId),
					"error": err,
				})
				if l.lm.isCoreLocalization(filePath) {
					return locErr
				}
				l.lm.logError(locErr)
			}
			l.Translations[scope][stateId] = trans
		}
	}
	if l.lm.isCoreLocalization(filePath) {
		l.lm.updateAvailableLanguages(supportedLanguages)
	}
	return nil
}

func (l *StateLocalizer) GetDescription(scope Scope, stateId string) (string, error) {
	states, scopeExists := l.Translations[scope]
	if !scopeExists {
		return "", NewAppError(ErrLocalization, "scope_not_found", map[string]any{
			"scope": scope,
		})
	}
	trans, stateExists := states[stateId]
	if !stateExists {
		return "", NewAppError(ErrLocalization, "state_localization_not_found", map[string]any{
			"state": string(scope) + "." + stateId,
		})
	}
	return fetchTranslation(l.lm, trans.Description)
}

func (l *StateLocalizer) GetMessage(scope Scope, stateId, messageKey string) (string, error) {
	states, scopeExists := l.Translations[scope]
	if !scopeExists {
		return "", NewAppError(ErrLocalization, "scope_not_found", map[string]any{
			"scope": scope,
		})
	}
	trans, stateExists := states[stateId]
	if !stateExists {
		return "", NewAppError(ErrLocalization, "state_localization_not_found", map[string]any{
			"state": string(scope) + "." + stateId,
		})
	}
	message, messageExists := trans.Messages[messageKey]
	if !messageExists {
		return "", NewAppError(ErrLocalization, "key_not_found", map[string]any{
			"key": string(scope) + "." + stateId + "." + messageKey,
		})
	}
	return fetchTranslation(l.lm, message)
}
