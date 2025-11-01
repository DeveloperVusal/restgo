package response

import (
	"encoding/json"
	"log/slog"
	"net/http"

	aslog "apibgo/pkg/logger/feature/slog"
)

type Status string
type StatusBadError string

const (
	StatusSuccess Status = "success"
	StatusError   Status = "error"
	StatusWarning Status = "warning"
)

const (
	StatusError2 StatusBadError = "error"
)

type Response struct {
	Code     byte
	Status   Status
	Message  string
	Result   interface{}
	Cookies  []*http.Cookie
	HttpCode int
}

type DocSuccessResponse struct {
	Code    byte
	Status  Status
	Message string
	Result  interface{}
}

type DocErrorResponse struct {
	Code    byte
	Status  StatusBadError
	Message string
}

func (response *Response) CreateResponseData() []byte {
	if response.Result == nil {
		marshal, _ := json.Marshal(map[string]interface{}{
			"code":    response.Code,
			"status":  response.Status,
			"message": response.Message,
		})
		return marshal
	}

	marshal, _ := json.Marshal(map[string]interface{}{
		"code":    response.Code,
		"status":  response.Status,
		"message": response.Message,
		"result":  response.Result,
	})

	return marshal
}

func (response *Response) SetCookies(w *http.ResponseWriter, log *slog.Logger) {
	if len(response.Cookies) > 0 {
		for _, _cookie := range response.Cookies {
			if err := _cookie.Valid(); err != nil {
				log.Warn("cookie invalid", aslog.Err(err))
			} else {
				http.SetCookie(*w, _cookie)
			}
		}
	}
}
