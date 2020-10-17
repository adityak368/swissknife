package i18n

import "fmt"

// i18nTranslator is a 18n translator that implements the localization.Translator interface
type i18nTranslator struct {
	translations map[string]string
}

// Tr Translates a key. Format the string using additional parameters
func (t *i18nTranslator) Tr(key string, args ...interface{}) string {

	translatedString, ok := t.translations[key]
	if !ok {
		return key
	}
	if len(args) > 0 {
		return fmt.Sprintf(translatedString, args...)
	}

	return translatedString
}
