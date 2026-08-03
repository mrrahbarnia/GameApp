package main

import (
	"fmt"
	"time"

	"github.com/mrrahbarnia/GameApp/config"
	"github.com/mrrahbarnia/GameApp/infrastructure/bcrypt"
	"github.com/mrrahbarnia/GameApp/infrastructure/postgresql"
	httpserver "github.com/mrrahbarnia/GameApp/presentation/http"

	authservice "github.com/mrrahbarnia/GameApp/service/auth"
	userservice "github.com/mrrahbarnia/GameApp/service/users"
)

const (
	JwtSignKey                 = "test"
	SigningMethod              = "HS256"
	AccessTokenSubject         = "ac"
	RefreshTokenSubject        = "rt"
	AccessTokenExpireDuration  = time.Hour * 24
	RefreshTokenExpireDuration = time.Hour * 24 * 7
)

// @title           GameApp API
// @version         1.0
// @description     This is a GameApp.
// @host            localhost:8080
// @BasePath        /api/v1
func main() {
	cfg := config.Config{
		HTTPServer: config.HTTPServer{Port: 8090},
		Auth: authservice.Config{
			SignKey:               JwtSignKey,
			SigningMethod:         SigningMethod,
			RefreshExpirationTime: RefreshTokenExpireDuration,
			AccessExpirationTime:  AccessTokenExpireDuration,
			RefreshSubject:        RefreshTokenSubject,
			AccessSubject:         AccessTokenSubject,
		},
		PostgreSQL: postgresql.Config{
			Username: "admin",
			Password: "123456",
			Host:     "localhost",
			Port:     5432,
			DBName:   "db",
		},
	}

	authSvc, userSvc := setupServices(cfg)

	server := httpserver.New(cfg, authSvc, userSvc)
	fmt.Println("start echo server")
	server.Serve()

}

func setupServices(cfg config.Config) (authservice.Service, userservice.Service) {
	authSvc := authservice.New(cfg.Auth)
	bcrypt := bcrypt.New()
	PGRepo := postgresql.New(cfg.PostgreSQL)
	userSvc := userservice.New(PGRepo, bcrypt, authSvc)

	return authSvc, userSvc
}
