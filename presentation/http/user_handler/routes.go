package user_handler

import (
	"github.com/golang-jwt/jwt/v4"
	mjwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
	authservice "github.com/mrrahbarnia/GameApp/service/auth"
)

func (h Handler) SetUserRoutes(e *echo.Echo) {
	usersGroup := e.Group("/users")

	usersGroup.POST("/register", h.register)
	usersGroup.POST("/login", h.login)
	usersGroup.GET("/profile", h.profile, mjwt.WithConfig(mjwt.Config{
		ContextKey:    "user",
		SigningKey:    []byte("test"),
		SigningMethod: "HS256",
		ParseTokenFunc: func(c *echo.Context, auth string) (interface{}, error) {

			token, err := jwt.ParseWithClaims(auth, &authservice.Claims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte("test"), nil
			})

			if err != nil {
				return nil, err
			}

			if claims, ok := token.Claims.(*authservice.Claims); ok && token.Valid {
				return claims, nil
			} else {
				return nil, err
			}
		},
	}))
}
