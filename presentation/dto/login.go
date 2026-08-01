package dto

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/mrrahbarnia/GameApp/pkg/errmsg"
	"github.com/mrrahbarnia/GameApp/pkg/richerror"
)

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (req LoginRequest) Validate() (map[string]string, error) {

	if err := validation.ValidateStruct(
		&req,
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
			richerror.New("dto.login").
				WithErr(err).
				WithKind(richerror.KindInvalid).
				WithMessage(errmsg.ErrorMsgInvalidInput).
				WithMeta(map[string]any{"req": req})
	}

	return nil, nil

}
