package routes

import (
	"encoding/json"
	"log"
	"net/http"
)

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Fields  string `json:"fields"`
}

func sendError(rw http.ResponseWriter, baseError error) {
	sendErrorCode(rw, baseError, http.StatusInternalServerError)
}
func sendErrorCode(rw http.ResponseWriter, baseError error, errorCode int) {
	errorMessage := baseError.Error()

	response := errorResponse{
		Code:    errorCode,
		Message: errorMessage,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error serialisation failed: %v.", err)
		jsonResponse = []byte(errorMessage)
	}

	http.Error(rw, string(jsonResponse), errorCode)
}
