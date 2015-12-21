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

func (e *apiError) Send(rw http.ResponseWriter) (int, error) {
	jsonError, err := json.Marshal(e)
	if err != nil {
		var code int
		if e.Code != -1 {
			code = e.Code
		} else {
			code = http.StatusInternalServerError
		}
		http.Error(rw, e.Message, code)
		return 0, http.ErrNotSupported
	}

	return rw.Write(jsonError)
}

func SendAPIError(baseError error, rw http.ResponseWriter) {
	newAPIError(baseError).Send(rw)
}
