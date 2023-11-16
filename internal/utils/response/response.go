package response

import (
	"encoding/json"
	"net/http"
)

type Status string

const (
	Success Status = "success"
	Error          = "error"
	Warning        = "warning"
)

type Response struct {
	Code    byte
	Status  Status
	Message string
	Result  any
	Cookies []*http.Cookie
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
