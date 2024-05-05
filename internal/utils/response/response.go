package response

import (
	"encoding/json"
	"log/slog"
	"net/http"

	aslog "apibgo/pkg/logger/feature/slog"
)

type Status string

const (
	StatusSuccess Status = "success"
	StatusError          = "error"
	StatusWarning        = "warning"
)

type Response struct {
	Code     byte
	Status   Status
	Message  string
	Result   interface{}
	Cookies  []*http.Cookie
	HttpCode int
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
