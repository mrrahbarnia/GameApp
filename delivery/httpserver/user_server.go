package httpserver

import (
	"net/http"

	"github.com/labstack/echo/v5"
	userservice "github.com/mrrahbarnia/GameApp/service/users"
)

func (s Server) userRegister(c *echo.Context) error {
	// curl -X POST "http://localhost:8090/users/register" \
	// -H "Content-Type: application/json" \
	// -d '{"name": "testUser", "phone_number": "09131234567", "password": "12345678"}'
	var req userservice.RegisterRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid payload")
	}

	if resp, err := s.userSvc.Register(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	} else {
		return c.JSON(http.StatusCreated, resp)
	}
}

func (s Server) login(c *echo.Context) error {
	var req userservice.LoginRequest

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
