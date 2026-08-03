package user_handler

import (
	authservice "github.com/mrrahbarnia/GameApp/service/auth"
	userservice "github.com/mrrahbarnia/GameApp/service/users"
)

type Handler struct {
	userSvc    userservice.Service
	authConfig authservice.Config
	authSvc    authservice.Service
}

func New(userSvc userservice.Service, authCfg authservice.Config, authSvc authservice.Service) Handler {
	return Handler{
		userSvc:    userSvc,
		authSvc:    authSvc,
		authConfig: authCfg,
	}
}
