package localization

import (
	"strings"

	"github.com/adityak368/swissknife/localization/i18n"
	"github.com/labstack/echo/v4"
)

// LocalizerConfig defines the localization config
type LocalizerConfig struct {
	InitializeFunc func() error
}

// DefaultLocalizerConfig defines the default localization config
var DefaultLocalizerConfig = LocalizerConfig{
	InitializeFunc: nil,
}

// EchoLocalizer returns a echo middleware for localization
func EchoLocalizer() echo.MiddlewareFunc {
	return EchoLocalizerWithConfig(DefaultLocalizerConfig)
}

// EchoLocalizerWithConfig returns a middleware for echo with config
func EchoLocalizerWithConfig(config LocalizerConfig) echo.MiddlewareFunc {
	if config.InitializeFunc != nil {
		config.InitializeFunc()
	}
	return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			language := strings.TrimSpace(c.Request().Header.Get("Accept-Language"))
			if language == "" {
				language = "en-US"
			}
			translator := i18n.Localizer().Translator(language)
			c.Set("translator", translator)
			return handlerFunc(c)
		}
	}
}
