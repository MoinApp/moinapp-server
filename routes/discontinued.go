package routes

import (
	"encoding/json"
	"errors"
	"net/http"
)

type discontinuationResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var (
	ErrAPIDiscontinued = errors.New("This API endpoint is discontinued.")
)

func discontinuationHandler(rw http.ResponseWriter, req *http.Request) {
	status := http.StatusGone

	response := discontinuationResponse{
		Code:    status,
		Message: ErrAPIDiscontinued.Error(),
	}

	data, err := json.Marshal(response)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	} else {
		http.Error(rw, string(data), status)
	}
}
