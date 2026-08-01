package user_handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
	userservice "github.com/mrrahbarnia/GameApp/service/users"
)

func (h Handler) profile(c *echo.Context) error {
	authToken := c.Request().Header.Get("Authorization")

	claims, err := h.authSvc.ParseToken(authToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	resp, err := h.userSvc.Profile(userservice.ProfileRequest{UserID: claims.UserID})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}
