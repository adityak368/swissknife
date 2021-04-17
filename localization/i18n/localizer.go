package i18n

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/adityak368/swissknife/localization"
)

// i18nLocalizer is the i18n localizer that implements the localization.Localizer interface
type i18nLocalizer struct {
	locales map[string]localization.Translator
}

// AddLocale Adds a locale to the localizer
func (l *i18nLocalizer) AddLocale(localeKey string, translations map[string]string) {
	translator := &i18nTranslator{
		translations: translations,
	}
	l.locales[localeKey] = translator
}

// Translator Returns a translator for the locale. If no locale is added, it returns an empty translator
func (l *i18nLocalizer) Translator(localeKey string) localization.Translator {
	translator, ok := l.locales[localeKey]
	if !ok {
		return &i18nTranslator{
			translations: make(map[string]string),
		}
	}
	return translator
}

// LoadJSONLocalesFromFolder Parses and Loads all the locales in json format in a folder (filename is used as the key of the locale. Ex: en.json -> "en", de.json -> "de")
func (l *i18nLocalizer) LoadJSONLocalesFromFolder(localesPath string) error {
	err := filepath.Walk(localesPath, func(path string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !fileInfo.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			jsonData, err := ioutil.ReadAll(file)
			if err != nil {
				return err
			}
			var translations map[string]string
			json.Unmarshal(jsonData, &translations)
			fileName := fileInfo.Name()
			localeKey := strings.TrimSuffix(fileName, filepath.Ext(fileName))
			l.AddLocale(localeKey, translations)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

var i18n localization.Localizer

// Localizer returns a singleton i18n localizer if you need a plug and play option
// We dont need a thread safe singleton in most cases as we dont change any data once loaded. So it is not thread safe
func Localizer() localization.Localizer {
	if i18n == nil {
		i18n = &i18nLocalizer{
			locales: make(map[string]localization.Translator),
		}
	}
	return i18n
}
