package httpserver

import (
	"log/slog"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/mrrahbarnia/GameApp/config"
	"github.com/mrrahbarnia/GameApp/presentation/http/user_handler"
	authservice "github.com/mrrahbarnia/GameApp/service/auth"
	userservice "github.com/mrrahbarnia/GameApp/service/users"
)

type Server struct {
	config      config.Config
	userHandler user_handler.Handler
}

func New(config config.Config, authSvc authservice.Service, userSvc userservice.Service) Server {
	return Server{config: config, userHandler: user_handler.New(userSvc, authSvc)}
}

func (s Server) Serve() {
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.GET("/health-check", s.healthCheck)

	s.userHandler.SetUserRoutes(e)

	if err := e.Start(":8090"); err != nil {
		slog.Error("failed to start server", "error", err)
	}

}
