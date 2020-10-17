package localization

// Localizer defines the interface for a localizer
type Localizer interface {
	AddLocale(localeKey string, translations map[string]string)
	Translator(localeKey string) Translator
	LoadJSONLocalesFromFolder(localesPath string) error
}
