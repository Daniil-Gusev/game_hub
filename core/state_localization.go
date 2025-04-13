package core

type StateTranslation struct {
	Description map[string]string            `json:"description"`
	Messages    map[string]map[string]string `json:"messages"`
}

type StateTranslations map[Scope]map[string]StateTranslation

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
	var rawData StateTranslations
	if err := l.lm.loadRawData(filePath, &rawData); err != nil {
		return err
	}

	for scope, states := range rawData {
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
			l.Translations[scope][stateId] = trans
		}
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
			"key": string(scope) + "." + stateId,
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
			"key": string(scope) + "." + stateId,
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
