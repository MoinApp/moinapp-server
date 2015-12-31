package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MoinApp/moinapp-server/push"
)

type moinRequest struct {
	Name string `json:"name"`
}

// TODO
func serveMoin(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var body moinRequest
	err := decoder.Decode(&body)
	if err != nil {
		SendAPIError(err, rw)
		return
	}

	fmt.Printf("Moin request: %+v\n", body)

	tokens := getUserFromRequest(req).GetPushTokens()
	for _, token := range tokens {
		push.SendPushNotification(token, "Moin")
	}
}
