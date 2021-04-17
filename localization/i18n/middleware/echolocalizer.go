package localization

import (
	"strings"

	"github.com/adityak368/swissknife/localization/i18n"
	"github.com/labstack/echo/v4"
)

// EchoLocalizerConfig defines the localization config
type EchoLocalizerConfig struct {
	InitializeFunc func()
}

// DefaultLocalizerConfig defines the default localization config
var DefaultLocalizerConfig = EchoLocalizerConfig{
	InitializeFunc: nil,
}

// EchoLocalizer returns a echo middleware for localization
func EchoLocalizer() echo.MiddlewareFunc {
	return EchoLocalizerWithConfig(DefaultLocalizerConfig)
}

// EchoLocalizerWithConfig returns a middleware for echo with config
func EchoLocalizerWithConfig(config EchoLocalizerConfig) echo.MiddlewareFunc {
	if config.InitializeFunc != nil {
		config.InitializeFunc()
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			language := strings.TrimSpace(c.Request().Header.Get("Accept-Language"))
			if language == "" {
				language = "en-US"
			}
			translator := i18n.Localizer().Translator(language)
			c.Set("translator", translator)
			return next(c)
		}
	}
}
