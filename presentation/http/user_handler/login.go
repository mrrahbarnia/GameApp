package user_handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/mrrahbarnia/GameApp/pkg/httpmsg"
	"github.com/mrrahbarnia/GameApp/presentation/dto"
)

func (h Handler) login(c *echo.Context) error {
	var req dto.LoginRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid payload")
	}

	errFields, err := req.Validate()
	if err != nil {
		msg, code := httpmsg.Error(err)

		return c.JSON(code, map[string]any{"message": msg, "errors": errFields})
	}

	if resp, err := h.userSvc.Login(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	} else {
		return c.JSON(http.StatusOK, resp)
	}
}
