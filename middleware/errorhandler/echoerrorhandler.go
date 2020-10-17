package errorhandler

import (
	"net/http"

	"github.com/adityak368/swissknife/localization"
	"github.com/adityak368/swissknife/logger"
	"github.com/adityak368/swissknife/response"
	"github.com/labstack/echo/v4"
)

// EchoHTTPErrorHandlerMiddleware defines the error handler middleware for echo
func EchoHTTPErrorHandlerMiddleware(err error, c echo.Context) {

	switch e := err.(type) {
	case *response.Error:
		translator := c.Get("translator").(localization.Translator)
		c.JSON(e.Code, e.Message().AsTranslatedJSONMessage(translator))
	case *echo.HTTPError:
		c.JSON(e.Code, response.JSONMessage{Message: e.Error()})
	default:
		logger.Error.Println(e.Error())
		c.JSON(http.StatusInternalServerError, map[string]string{"message": http.StatusText(http.StatusInternalServerError)})
	}
}
