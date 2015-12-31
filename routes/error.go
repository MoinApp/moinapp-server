package routes

import (
	"encoding/json"
	"net/http"
)

type apiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Fields  string `json:"fields"`
}

func newAPIError(baseError error) *apiError {
	return &apiError{
		Code:    -1,
		Message: baseError.Error(),
	}
}

func (e *apiError) Send(rw http.ResponseWriter) {
	if e.Code == -1 {
		e.Code = http.StatusBadRequest
	}

	jsonError, err := json.Marshal(e)
	if err != nil {
		jsonError = []byte(e.Message)
	}

	http.Error(rw, string(jsonError), e.Code)
}

func SendAPIError(baseError error, rw http.ResponseWriter) {
	newAPIError(baseError).Send(rw)
}
