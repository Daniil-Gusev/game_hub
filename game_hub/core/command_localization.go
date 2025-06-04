package core

import (
	"fmt"
)

type CommandTranslation struct {
	Name        map[string]string   `json:"name"`
	Description map[string]string   `json:"description"`
	Aliases     map[string][]string `json:"aliases"`
}

func (c CommandTranslation) isLocalized(langs []string) error {
	for _, supportedLang := range langs {
		if _, exists := c.Name[supportedLang]; exists {
			continue
		}
		return NewAppError(Err, "key_not_found", map[string]any{
			"key": fmt.Sprintf("Name.%s", supportedLang),
		})
	}
	for _, supportedLang := range langs {
		if len(c.Description) == 0 {
			break
		}
		if _, exists := c.Description[supportedLang]; exists {
			continue
		}
		return NewAppError(Err, "key_not_found", map[string]any{
			"key": fmt.Sprintf("description.%s", supportedLang),
		})
	}
	for _, supportedLang := range langs {
		if len(c.Aliases) == 0 {
			break
		}
		if _, exists := c.Aliases[supportedLang]; exists {
			continue
		}
		return NewAppError(Err, "key_not_found", map[string]any{
			"key": fmt.Sprintf("aliases.%s", supportedLang),
		})
	}
	return nil
}

func NewCommandTranslation() CommandTranslation {
	return CommandTranslation{
		Name:        make(map[string]string),
		Description: make(map[string]string),
		Aliases:     make(map[string][]string),
	}
}

type CommandTranslations map[Scope]map[string]CommandTranslation

type CommandLocalizationData struct {
	Meta         LocalizationMetadata `json:"meta"`
	Translations CommandTranslations  `json:"translations"`
}

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
	var rawData CommandLocalizationData
	if err := l.lm.loadLocalizationData(filePath, &rawData); err != nil {
		return err
	}
	supportedLanguages := rawData.Meta.SupportedLanguages
	for scope, cmds := range rawData.Translations {
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
				l.lm.logError(NewAppError(ErrLocalization, "localization_file_translations_error", map[string]any{
					"file": filePath,
					"path": fmt.Sprintf("translations.%v.%s", scope, cmdId),
					"error": NewAppError(Err, "key_not_found", map[string]any{
						"key": fmt.Sprintf("%s.name", cmdId),
					}),
				}))
				continue
			}
			if trans.Description == nil {
				trans.Description = make(map[string]string)
			}
			if trans.Aliases == nil {
				trans.Aliases = make(map[string][]string)
			}
			err := trans.isLocalized(supportedLanguages)
			if err != nil {
				locErr := NewAppError(ErrLocalization, "localization_file_translations_error", map[string]any{
					"file":  filePath,
					"path":  fmt.Sprintf("translations.%v.%s", scope, cmdId),
					"error": err,
				})
				if l.lm.isCoreLocalization(filePath) {
					return locErr
				}
				l.lm.logError(locErr)
			}
			l.Translations[scope][cmdId] = trans
		}
	}
	if l.lm.isCoreLocalization(filePath) {
		l.lm.updateAvailableLanguages(supportedLanguages)
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
			"command": fmt.Sprintf("%s.%s", scope, cmdId),
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
		return "", NewAppError(ErrLocalization, "command_localization_not_found", map[string]any{
			"command": fmt.Sprintf("%s.%s", scope, cmdId),
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
		return nil, NewAppError(ErrLocalization, "command_localization_not_found", map[string]any{
			"command": fmt.Sprintf("%s.%s", scope, cmdId),
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
