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

	if data, err := json.Marshal(response); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	} else {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(status)
		rw.Write([]byte(data))
	}
}
