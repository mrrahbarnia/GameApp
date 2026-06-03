package httpmsg

import (
	"net/http"

	"github.com/mrrahbarnia/GameApp/pkg/errmsg"
	"github.com/mrrahbarnia/GameApp/pkg/richerror"
)

func Error(err error) (message string, code int) {
	switch err.(type) {
	case richerror.RichError:
		re := err.(richerror.RichError)
		msg := re.Message()

		// we should not expose unexpected error messages
		code := mapKindToHTTPStatusCode(re.Kind())
		if code >= 500 {
			msg = errmsg.ErrorMsgSomethingWentWrong
		}

		return msg, code
	default:
		return err.Error(), http.StatusBadRequest
	}
}

func mapKindToHTTPStatusCode(kind richerror.Kind) int {
	switch kind {
	case richerror.KindInvalid:
		return http.StatusUnprocessableEntity
	case richerror.KindNotFound:
		return http.StatusNotFound
	case richerror.KindForbidden:
		return http.StatusForbidden
	case richerror.KindUnexpected:
		return http.StatusInternalServerError
	case richerror.KindConflict:
		return http.StatusConflict
	default:
		return http.StatusBadRequest
	}
}
