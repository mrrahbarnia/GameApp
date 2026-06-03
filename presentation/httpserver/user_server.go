package httpserver

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/mrrahbarnia/GameApp/pkg/httpmsg"
	"github.com/mrrahbarnia/GameApp/presentation/dto"
	userservice "github.com/mrrahbarnia/GameApp/service/users"
)

func (s Server) userRegister(c *echo.Context) error {
	var req dto.RegisterRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid payload")
	}

	errFields, err := req.Validate()
	if err != nil {
		msg, code := httpmsg.Error(err)

		return c.JSON(code, map[string]any{"message": msg, "errors": errFields})
	}

	if resp, err := s.userSvc.Register(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	} else {
		return c.JSON(http.StatusCreated, resp)
	}
}

func (s Server) login(c *echo.Context) error {
	var req dto.LoginRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid payload")
	}

	if resp, err := s.userSvc.Login(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	} else {
		return c.JSON(http.StatusOK, resp)
	}
}

func (s Server) profile(c *echo.Context) error {
	authToken := c.Request().Header.Get("Authorization")

	claims, err := s.authSvc.ParseToken(authToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	resp, err := s.userSvc.Profile(userservice.ProfileRequest{UserID: claims.UserID})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}
