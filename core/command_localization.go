package core

import (
	"fmt"
)

type CommandTranslation struct {
	Name        map[string]string   `json:"name"`
	Description map[string]string   `json:"description"`
	Aliases     map[string][]string `json:"aliases"`
}

func NewCommandTranslation() CommandTranslation {
	return CommandTranslation{
		Name:        make(map[string]string),
		Description: make(map[string]string),
		Aliases:     make(map[string][]string),
	}
}

type CommandTranslations map[Scope]map[string]CommandTranslation

type CommandLocalizer struct {
	lm           *LocalizationManager
	Translations CommandTranslations
}

func NewCommandLocalizer(lm *LocalizationManager) *CommandLocalizer {
	return &CommandLocalizer{
		lm:           lm,
		Translations: make(CommandTranslations),
	}
}

func (l *CommandLocalizer) LoadTranslations(filePath string) error {
	var rawData CommandTranslations
	if err := l.lm.loadRawData(filePath, &rawData); err != nil {
		return err
	}
	for scope, cmds := range rawData {
		if !scope.IsValid() {
			return NewAppError(ErrLocalization, "invalid_command_localizations_scope", map[string]any{
				"file":  filePath,
				"scope": scope,
			})
		}
		if _, exists := l.Translations[scope]; !exists {
			l.Translations[scope] = make(map[string]CommandTranslation)
		}
		for cmdId, trans := range cmds {
			if len(trans.Name) == 0 {
				return NewAppError(ErrLocalization, "key_not_found", map[string]any{
					"key": fmt.Sprintf("%s.name", cmdId),
				})
				continue
			}
			if trans.Description == nil {
				trans.Description = make(map[string]string)
			}
			if trans.Aliases == nil {
				trans.Aliases = make(map[string][]string)
			}
			l.Translations[scope][cmdId] = trans
		}
	}
	return nil
}

func (l *CommandLocalizer) GetName(scope Scope, cmdId string) (string, error) {
	cmds, scopeExists := l.Translations[scope]
	if !scopeExists {
		return "", NewAppError(ErrLocalization, "scope_not_found", map[string]any{
			"scope": scope,
		})
	}
	trans, cmdExists := cmds[cmdId]
	if !cmdExists {
		return "", NewAppError(ErrLocalization, "command_localization_not_found", map[string]any{
			"key": fmt.Sprintf("%s.%s", scope, cmdId),
		})
	}
	return fetchTranslation(l.lm, trans.Name)
}

func (l *CommandLocalizer) GetDescription(scope Scope, cmdId string) (string, error) {
	cmds, scopeExists := l.Translations[scope]
	if !scopeExists {
		return "", NewAppError(ErrLocalization, "scope_not_found", map[string]any{
			"scope": scope,
		})
	}
	trans, cmdExists := cmds[cmdId]
	if !cmdExists {
		return "", NewAppError(ErrLocalization, "key_not_found", map[string]any{
			"key": fmt.Sprintf("%s.%s", scope, cmdId),
		})
	}
	return fetchTranslation(l.lm, trans.Description)
}

func (l *CommandLocalizer) GetAliases(scope Scope, cmdId string) ([]string, error) {
	cmds, scopeExists := l.Translations[scope]
	if !scopeExists {
		return nil, NewAppError(ErrLocalization, "scope_not_found", map[string]any{
			"scope": scope,
		})
	}

	trans, cmdExists := cmds[cmdId]
	if !cmdExists {
		return nil, NewAppError(ErrLocalization, "key_not_found", map[string]any{
			"key": fmt.Sprintf("%s.%s", scope, cmdId),
		})
	}

	aliases, err := fetchTranslation(l.lm, trans.Aliases)
	if err != nil {
		return []string{}, nil
	}
	return aliases, nil
}

func (l *CommandLocalizer) Exists(scope Scope, cmdId string) bool {
	cmds, scopeExists := l.Translations[scope]
	if !scopeExists {
		return false
	}
	_, exists := cmds[cmdId]
	return exists
}
