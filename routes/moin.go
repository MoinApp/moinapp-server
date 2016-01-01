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
		sendErrorCode(rw, err, http.StatusBadRequest)
		return
	}

	fmt.Printf("Moin request: %+v\n", body)

	currentUser := getUserFromRequest(req)
	targetUser := models.FindUserByName(body.Name)
	if !targetUser.IsResult() {
		sendErrorCode(rw, ErrUserNotFound, http.StatusBadRequest)
		return
	}

	push.SendMoinNotificationToUser(targetUser, currentUser)
	currentUser.AddRecentUser(targetUser)

	rw.WriteHeader(http.StatusOK)
}
