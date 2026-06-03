package dto

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/mrrahbarnia/GameApp/pkg/errmsg"
	"github.com/mrrahbarnia/GameApp/pkg/richerror"
)

type RegisterRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type RegisterResponse struct {
	UserID      uint   `json:"user_id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}

func (req RegisterRequest) Validate() (map[string]string, error) {
	if err := validation.ValidateStruct(
		&req,
		validation.Field(&req.Name, validation.Required, validation.Length(3, 255)),
		// Todo - read password regex from .env
		validation.Field(&req.Password, validation.Required, validation.Match(regexp.MustCompile(`^.{8,}$`))),
		validation.Field(&req.PhoneNumber, validation.Required, validation.Match(regexp.MustCompile(phoneNumberRegex))),
	); err != nil {
		fieldErrs := make(map[string]string)

		errV, ok := err.(validation.Errors)
		if ok {
			for v, e := range errV {
				if e != nil {
					fieldErrs[v] = e.Error()
				}
			}
		}
		return fieldErrs,
			richerror.New("dto.register").
				WithErr(err).
				WithKind(richerror.KindInvalid).
				WithMessage(errmsg.ErrorMsgInvalidInput).
				WithMeta(map[string]any{"req": req})
	}

	return nil, nil
}
