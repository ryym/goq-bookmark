package handler

import (
	"net/http"

	"github.com/labstack/echo"
)

func Home() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "home", nil)
	}
}
