package httpserver

import (
	"log/slog"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/mrrahbarnia/GameApp/config"
	authservice "github.com/mrrahbarnia/GameApp/service/auth"
	userservice "github.com/mrrahbarnia/GameApp/service/users"
)

type Server struct {
	config  config.Config
	userSvc userservice.Service
	authSvc authservice.Service
}

func New(config config.Config, userSvc userservice.Service, authSvc authservice.Service) Server {
	return Server{config: config, userSvc: userSvc, authSvc: authSvc}
}

func (s Server) Serve() {
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.GET("/health-check", s.healthCheck)
	e.POST("/users/register", s.userRegister)

	if err := e.Start(":8090"); err != nil {
		slog.Error("failed to start server", "error", err)
	}

}
