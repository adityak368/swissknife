package localization

import (
	"strings"

	"github.com/adityak368/swissknife/localization/i18n"
	"github.com/labstack/echo/v4"
)

// EchoLocalizer returns a echo middleware for localization
func EchoLocalizer() echo.MiddlewareFunc {
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
