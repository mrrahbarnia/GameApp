package config

import (
	"github.com/mrrahbarnia/GameApp/infrastructure/postgresql"
	authservice "github.com/mrrahbarnia/GameApp/service/auth"
)

type HTTPServer struct {
	Port int
}

type Config struct {
	HTTPServer HTTPServer
	Auth       authservice.Config
	PostgreSQL postgresql.Config
}
