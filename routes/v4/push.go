package v4

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MoinApp/moinapp-server/models"
	"net/http"
)

type addPushTokenRequest struct {
	Token string `json:"token"`
	Type  string `json:"type"`
}

const (
	APNTokenType = "apns"
	GCMTokenType = "gcm"
)

var (
	ErrInvalidTokenType = errors.New("Invalid token type.")
)

func serveAddPushToken(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var request addPushTokenRequest
	err := decoder.Decode(&request)
	if err != nil {
		sendErrorCode(rw, err, http.StatusBadRequest)
		return
	}

	var tokenDBType models.TokenType
	switch request.Type {
	case APNTokenType:
		tokenDBType = models.APNToken
	case GCMTokenType:
		tokenDBType = models.GCMToken
	default:
		sendErrorCode(rw, ErrInvalidTokenType, http.StatusBadRequest)
		return
	}

	fmt.Printf("Add Push Token request: %+v\n", request)
	token := models.NewPushToken(tokenDBType, request.Token)
	currentUser := getUserFromRequest(req)

	currentUser.AddPushToken(token)

	rw.WriteHeader(http.StatusOK)
}
