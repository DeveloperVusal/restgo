package response

import (
	"encoding/json"
	"net/http"
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
