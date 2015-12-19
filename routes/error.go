package routes

import (
	"encoding/json"
	"net/http"
)

type APIError struct {
	Code    int32
	Message string
	Fields  string
}

func NewAPIError(baseError error) *APIError {
	return &APIError{
		Code:    -1,
		Message: baseError.Error(),
	}
}

func (e *APIError) Send(rw http.ResponseWriter) (int, error) {
	jsonError, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}

	return rw.Write(jsonError)
}

func SendAPIError(baseError error, rw http.ResponseWriter) {
	NewAPIError(baseError).Send(rw)
}
