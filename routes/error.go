package routes

import (
	"encoding/json"
	"log"
	"net/http"
)

type apiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Fields  string `json:"fields"`
}

func newAPIError(baseError error) *apiError {
	return newAPIErrorCode(baseError, -1)
}
func newAPIErrorCode(baseError error, code int) *apiError {
	return &apiError{
		Code:    code,
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

	if e.Code == http.StatusInternalServerError {
		log.Printf("Error code %v sent: %q.", e.Code, e.Message)
	}
}

func SendAPIError(baseError error, rw http.ResponseWriter) {
	newAPIError(baseError).Send(rw)
}
func SendAPIErrorCode(baseError error, code int, rw http.ResponseWriter) {
	newAPIErrorCode(baseError, code).Send(rw)
}
