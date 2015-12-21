package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type addPushTokenRequest struct {
	Token string `json:"token"`
	Type  string `json:"type"`
}

func serveAddPushToken(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var request addPushTokenRequest
	err := decoder.Decode(&request)
	if err != nil {
		SendAPIError(err, rw)
		return
	}

	fmt.Printf("Add Push Token request: %+v\n", request)
}
