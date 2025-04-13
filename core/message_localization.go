package core

type MessageTranslations map[string]map[string]string

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
	var rawData MessageTranslations
	if err := l.lm.loadRawData(filePath, &rawData); err != nil {
		return err
	}
	l.CopyTranslations(l.Translations, rawData)
	return nil
}

func (l *MessageLocalizer) LoadOptionalTranslations(filePath string) error {
	var rawData map[string]MessageTranslations
	if err := l.lm.loadRawData(filePath, &rawData); err != nil {
		return err
	}
	for setName, set := range rawData {
		if _, exists := l.OptionalTranslations[setName]; !exists {
			l.OptionalTranslations[setName] = make(MessageTranslations)
		}
		l.CopyTranslations(l.OptionalTranslations[setName], set)
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

func (l *MessageLocalizer) CopyTranslations(dest, source map[string]map[string]string) {
	for msgKey, trans := range source {
		if _, exists := dest[msgKey]; !exists {
			dest[msgKey] = make(map[string]string)
		}
		for lang, msg := range trans {
			dest[msgKey][lang] = msg
		}
	}
}
