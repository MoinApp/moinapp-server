package routes

// TODO JSON-Web-Tokens based authentication flow for all endpoints

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MoinApp/moinapp-server/auth"
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
		user := models.CreateUser(body.Name, body.Password, body.Email)

		tokenResponse := getSessionResponseTokenForUser(user)
		rw.Write(tokenResponse)
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
	user := models.FindUserWithCredentials(body.Name, body.Password)

	if user == nil {
		// TODO: WTF
		panic(errors.New("Invalid credentials"))
	}

	tokenResponse := getSessionResponseTokenForUser(user)
	rw.Write(tokenResponse)
}

func getSessionResponseTokenForUser(user *models.User) []byte {
	token, err := auth.GenerateJWTToken(*user)
	if err != nil {
		// TODO: WTF
		panic(err)
	}

	data, err := json.Marshal(sessionResponse{
		SessionToken: token,
	})
	if err != nil {
		// TODO: WTF
		panic(err)
	}

	return data
}
