package routes

import (
	"encoding/json"
	"fmt"
	"github.com/MoinApp/moinapp-server/models"
	"net/http"
)

type signUpRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type authenticationRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type sessionResponse struct {
	// Token for a session for this user.
	SessionToken string `json:"session_token"`
}

func serveSignUp(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var body signUpRequest
	err := decoder.Decode(&body)
	if err != nil {
		SendAPIError(err, rw)
		return
	}

	fmt.Printf("Create user request: %+v.\n", body)
	if !models.IsUsernameTaken(body.Name) {
		models.CreateUser(body.Name, body.Password, body.Email)

		data, _ := json.Marshal(sessionResponse{
			SessionToken: "null",
		})
		rw.Write(data)
	}
}

func serveAuthentication(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var body authenticationRequest
	err := decoder.Decode(&body)
	if err != nil {
		SendAPIError(err, rw)
		return
	}

	fmt.Printf("Auth request: %+v\n", body)

	data, _ := json.Marshal(sessionResponse{
		SessionToken: "null",
	})
	rw.Write(data)
}
