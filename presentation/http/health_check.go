package httpserver

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func (s Server) healthCheck(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}
