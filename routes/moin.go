package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MoinApp/moinapp-server/models"
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
		SendAPIErrorCode(err, http.StatusInternalServerError, rw)
		return
	}

	fmt.Printf("Moin request: %+v\n", body)

	currentUser := getUserFromRequest(req)
	targetUser := models.FindUserByName(body.Name)
	if !targetUser.IsResult() {
		SendAPIError(ErrUserNotFound, rw)
		return
	}

	tokens := targetUser.GetPushTokens()
	if len(tokens) > 0 {
		message := fmt.Sprintf("Moin from %v.", currentUser.Name)
		push.SendPushNotificationToAll(tokens, message)
	}
}
