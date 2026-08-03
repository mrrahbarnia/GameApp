package dependencies

import (
	"github.com/labstack/echo/v5"
	"github.com/mrrahbarnia/GameApp/pkg/constant"
	authservice "github.com/mrrahbarnia/GameApp/service/auth"
)

func GetUser(c *echo.Context) uint {
	claimsPtr, _ := c.Get(constant.ClaimContextKey).(*authservice.Claims)

	return *&claimsPtr.UserID
}
