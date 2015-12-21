package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type addPushTokenRequest struct {
	// Token brought to you by the Push Service Provider.
	Token string `json:"token"`
	// The Type of Pushservice (one of 'apns' and 'gcm')
	Type string `json:"type"`
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
