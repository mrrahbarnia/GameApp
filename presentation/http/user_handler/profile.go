package user_handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/mrrahbarnia/GameApp/presentation/dto"
)

func (h Handler) profile(c *echo.Context) error {
	authToken := c.Request().Header.Get("Authorization")

	claims, err := h.authSvc.ParseToken(authToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	result, err := h.userSvc.Profile(claims.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(
		http.StatusOK,
		dto.ProfileResponse{
			Name: result.Name,
		},
	)
}
