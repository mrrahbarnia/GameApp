package user_handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/mrrahbarnia/GameApp/presentation/dto"
	"github.com/mrrahbarnia/GameApp/presentation/http/dependencies"
)

func (h Handler) profile(c *echo.Context) error {
	userID := dependencies.GetUser(c)
	result, err := h.userSvc.Profile(userID)
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
