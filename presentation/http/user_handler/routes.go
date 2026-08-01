package user_handler

import "github.com/labstack/echo/v5"

func (h Handler) SetUserRoutes(e *echo.Echo) {
	usersGroup := e.Group("/users")

	usersGroup.POST("/register", h.register)
	usersGroup.POST("/login", h.login)
	usersGroup.GET("/profile", h.profile)
}
