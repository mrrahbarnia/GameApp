package main

import (
	"time"

	"github.com/mrrahbarnia/GameApp/config"
	"github.com/mrrahbarnia/GameApp/delivery/httpserver"
	"github.com/mrrahbarnia/GameApp/infrastructure/bcrypt"
	"github.com/mrrahbarnia/GameApp/infrastructure/postgresql"
	authservice "github.com/mrrahbarnia/GameApp/service/auth"
	userservice "github.com/mrrahbarnia/GameApp/service/users"
)

const (
	JwtSignKey                 = "jwt_secret"
	AccessTokenSubject         = "ac"
	RefreshTokenSubject        = "rt"
	AccessTokenExpireDuration  = time.Hour * 24
	RefreshTokenExpireDuration = time.Hour * 24 * 7
)

func main() {
	cfg := config.Config{
		HTTPServer: config.HTTPServer{Port: 8090},
		Auth: authservice.Config{
			SignKey:               JwtSignKey,
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

	server := httpserver.New(cfg, userSvc, authSvc)
	server.Serve()

}

func setupServices(cfg config.Config) (authservice.Service, userservice.Service) {
	authSvc := authservice.New(cfg.Auth)
	bcrypt := bcrypt.New()
	PGRepo := postgresql.New(cfg.PostgreSQL)
	userSvc := userservice.New(PGRepo, bcrypt, authSvc)

	return authSvc, userSvc
}
