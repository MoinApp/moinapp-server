package v4

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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
	body.Name = strings.Trim(body.Name, " ")

	fmt.Printf("Moin request: %+v\n", body)

	if body.Name == "" {
		sendErrorCode(rw, ErrBadRequest, http.StatusBadRequest)
		return
	}

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
