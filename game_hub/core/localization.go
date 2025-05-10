package core

import (
	"game_hub/config"
	"log"
	"sort"
)

type Language struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type LocalizationMetadata struct {
	SupportedLanguages []string `json:"supported_languages" validate:"required,min=1"`
}

func NewLocalizationMetadata() LocalizationMetadata {
	return LocalizationMetadata{
		SupportedLanguages: make([]string, 0, 10),
	}
}

type LocalizationData struct {
	Meta         LocalizationMetadata `json:"meta" validate:"required"`
	Translations map[string]any       `json:"translations" validate:"required"`
}

type LangDictData struct {
	Languages map[string]string `json:"languages" validate:"required"`
}

type LocalizationManager struct {
	cfg            *config.Config
	logger         Logger
	currentLang    string
	defaultLang    string
	availableLangs []Language
	langDict       map[string]string
}

func NewLocalizationManager(cfg *config.Config) (*LocalizationManager, error) {
	lm := &LocalizationManager{
		cfg:            cfg,
		availableLangs: make([]Language, 0, 10),
		langDict:       make(map[string]string, 10),
	}
	currentLang := cfg.Language.CurrentLanguage
	defaultLang := cfg.Language.DefaultLanguage
	dictFilePath := cfg.Paths.CoreLanguagesPath()
	var rawData LangDictData
	if err := LoadData(dictFilePath, &rawData); err != nil {
		return nil, NewAppError(ErrLocalization, "load_lang_dict_error", map[string]any{
			"file":  dictFilePath,
			"error": err,
		})
	}
	for key, value := range rawData.Languages {
		lm.langDict[key] = value
	}
	if !lm.isLanguageExists(currentLang) {
		return nil, NewAppError(ErrLocalization, "Current configuration Language is not supported.", map[string]any{"lang": defaultLang})
	}
	if !lm.isLanguageExists(defaultLang) {
		return nil, NewAppError(ErrLocalization, "Default configuration language is not supported.", map[string]any{"lang": defaultLang})
	}
	lm.currentLang = currentLang
	lm.defaultLang = defaultLang
	// не забыть перед использованием установить логгер
	return lm, nil
}

func (lm *LocalizationManager) SetLogger(logger Logger) {
	lm.logger = logger
}

func (lm *LocalizationManager) loadLocalizationData(filePath string, target any) error {
	data, err := ReadFile(filePath)
	if err != nil {
		return err
	}
	var rawData LocalizationData
	if err := DecodeData(data, &rawData); err != nil {
		return NewAppError(ErrLocalization, "file_parse_error", map[string]any{
			"file":  filePath,
			"error": err,
		})
	}
	validMeta, logErrors, err := lm.validateMetadata(rawData.Meta)
	if logErrors != nil {
		lm.logError(NewAppError(ErrLocalization, "invalid_localization_metadata", map[string]any{
			"file":  filePath,
			"error": logErrors,
		}))
	}
	if err != nil {
		return NewAppError(ErrLocalization, "invalid_localization_metadata", map[string]any{
			"file":  filePath,
			"error": err,
		})
	}
	rawData.Meta = validMeta
	updatedData, err := EncodeData(rawData)
	if err != nil {
		return NewAppError(ErrLocalization, "file_parse_error", map[string]any{
			"file":  filePath,
			"error": err,
		})
	}
	if err := DecodeData(updatedData, target); err != nil {
		return NewAppError(ErrLocalization, "file_parse_error", map[string]any{
			"file":  filePath,
			"error": err,
		})
	}
	return nil
}

func (lm *LocalizationManager) updateAvailableLanguages(supportedLangs []string) {
	langMap := make(map[string]void, len(lm.availableLangs))
	for _, lang := range lm.availableLangs {
		langMap[lang.Code] = void{}
	}

	for _, code := range supportedLangs {
		if _, exists := langMap[code]; !exists {
			name, ok := lm.langDict[code]
			if !ok {
				lm.log("Warning: Language code '%s' not found in dictionary.", code)
				continue
			}
			lm.availableLangs = append(lm.availableLangs, Language{Code: code, Name: name})
			langMap[code] = void{}
		}
	}

	sort.Slice(lm.availableLangs, func(i, j int) bool {
		return lm.availableLangs[i].Name < lm.availableLangs[j].Name
	})
}

func (lm *LocalizationManager) SetCurrentLanguage(lang string) error {
	if !lm.isLanguageSupported(lang) {
		return NewAppError(ErrLocalization, "lang_not_supported", map[string]any{"lang": lang})
	}
	lm.currentLang = lang
	return nil
}

func (lm *LocalizationManager) SetDefaultLanguage(lang string) error {
	if !lm.isLanguageSupported(lang) {
		return NewAppError(ErrLocalization, "lang_not_supported", map[string]any{"lang": lang})
	}
	lm.defaultLang = lang
	return nil
}

func (lm *LocalizationManager) CurrentLang() string {
	return lm.currentLang
}

func (lm *LocalizationManager) DefaultLang() string {
	return lm.defaultLang
}

func (lm *LocalizationManager) AvailableLanguages() []Language {
	return lm.availableLangs
}

func (lm *LocalizationManager) isLanguageSupported(code string) bool {
	for _, lang := range lm.availableLangs {
		if lang.Code == code {
			return true
		}
	}
	return false
}

func (lm *LocalizationManager) isLanguageExists(lang string) bool {
	_, exists := lm.langDict[lang]
	return exists
}
func (lm *LocalizationManager) getLanguageName(lang string) string {
	if !lm.isLanguageExists(lang) {
		return ""
	}
	return lm.langDict[lang]
}

func (lm *LocalizationManager) validateMetadata(meta LocalizationMetadata) (LocalizationMetadata, error, error) {
	validMeta := NewLocalizationMetadata()
	var logErrors error = nil
	if len(meta.SupportedLanguages) == 0 {
		return validMeta, nil, NewAppError(Err, "localization_metadata_missing_languages", nil)
	}
	validLangs := make([]string, 0, len(meta.SupportedLanguages))
	errs := make([]error, 0, 10)
	for _, lang := range meta.SupportedLanguages {
		if !lm.isLanguageExists(lang) {
			errs = append(errs, NewAppError(Err, "localization_metadata_invalid_language", map[string]any{"lang": lang}))
			continue
		}
		validLangs = append(validLangs, lang)
	}
	if len(errs) > 0 {
		logErrors = NewAppErrors(errs)
	}
	if len(validLangs) == 0 {
		return validMeta, logErrors, NewAppError(Err, "localization_metadata_missing_languages", nil)
	}
	isDefaultLang := false
	for _, lang := range validLangs {
		if lang == lm.defaultLang {
			isDefaultLang = true
		}
	}
	if len(errs) > 0 {
		logErrors = NewAppErrors(errs)
	}
	if !isDefaultLang {
		return validMeta, logErrors, NewAppError(Err, "localization_meta_missing_dedefaultlang", map[string]any{"lang": lm.getLanguageName(lm.defaultLang)})
	}
	validMeta.SupportedLanguages = validLangs
	return validMeta, logErrors, nil
}

func (lm *LocalizationManager) isCoreLocalization(filePath string) bool {
	return lm.cfg.Paths.IsCorePath(filePath)
}

func (lm *LocalizationManager) log(format string, v ...any) {
	if lm.logger == nil {
		log.Println("Logger is not set.")
		return
	}
	lm.logger.Printf(format, v...)
}
func (lm *LocalizationManager) logError(err error) {
	if lm.logger == nil {
		log.Println("Logger is not set.")
		return
	}
	lm.logger.Error(err)
}

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
