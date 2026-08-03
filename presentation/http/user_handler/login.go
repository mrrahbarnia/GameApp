package user_handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/mrrahbarnia/GameApp/pkg/httpmsg"
	"github.com/mrrahbarnia/GameApp/presentation/dto"
)

// Login godoc
// @Summary      Login API
// @Description  Returns access and refresh tokens
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "Login credentials"
// @Success      200  {object}  dto.LoginResponse
// @Failure      404  {object}  map[string]string
// @Router       /users/login [post]
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

	if result, err := h.userSvc.Login(req.PhoneNumber, req.Password); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	} else {
		return c.JSON(
			http.StatusOK,
			dto.LoginResponse{AccessToken: result.AccessToken, RefreshToken: result.RefreshToken},
		)
	}
}
