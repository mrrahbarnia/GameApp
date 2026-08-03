package middleware

import (
	mjwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
	"github.com/mrrahbarnia/GameApp/pkg/constant"
	authservice "github.com/mrrahbarnia/GameApp/service/auth"
)

func Auth(authSvc authservice.Service, authCfg authservice.Config) echo.MiddlewareFunc {
	return mjwt.WithConfig(mjwt.Config{
		ContextKey:    constant.ClaimContextKey,
		SigningKey:    []byte(authCfg.SignKey),
		SigningMethod: authCfg.SigningMethod,
		ParseTokenFunc: func(c *echo.Context, auth string) (interface{}, error) {

			claims, err := authSvc.ParseToken(auth)
			if err != nil {
				return nil, err
			}
			return claims, nil

		},
	})
}
