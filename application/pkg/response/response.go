package response

import (
	"strings"

	"github.com/labstack/echo/v4"
)

func Error(c echo.Context, status int, err error) error {
	msg := err.Error()
	//  custom error message handler
	if strings.Contains(err.Error(), "duplicate key value") {
		msg = "resource already exists"
	}

	return c.JSON(status, echo.Map{
		"message": msg,
	})
}
