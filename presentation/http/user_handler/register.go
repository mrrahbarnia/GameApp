package user_handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/mrrahbarnia/GameApp/pkg/httpmsg"
	"github.com/mrrahbarnia/GameApp/presentation/dto"
)

func (h Handler) register(c *echo.Context) error {
	var req dto.RegisterRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid payload")
	}

	errFields, err := req.Validate()
	if err != nil {
		msg, code := httpmsg.Error(err)

		return c.JSON(code, map[string]any{"message": msg, "errors": errFields})
	}

	if result, err := h.userSvc.Register(req.Name, req.PhoneNumber, req.Password); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	} else {
		return c.JSON(http.StatusCreated,
			dto.RegisterResponse{
				UserID:      result.UserID,
				PhoneNumber: result.PhoneNumber,
				Name:        result.Name,
			},
		)
	}
}
