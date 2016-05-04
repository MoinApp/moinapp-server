package v4

// TODO JSON-Web-Tokens based authentication flow for all endpoints

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/MoinApp/moinapp-server/models"
	"github.com/MoinApp/moinapp-server/routes/v4/auth"
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

var (
	ErrInvalidCredentials = errors.New("Invalid credentials")
	ErrUsernameTaken      = errors.New("Username is taken already")
)

func serveSignUp(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var body signUpRequest
	err := decoder.Decode(&body)
	if err != nil {
		sendErrorCode(rw, err, http.StatusBadRequest)
		return
	}
	body.Email = strings.Trim(body.Email, " ")
	body.Name = strings.Trim(body.Name, " ")
	body.Password = strings.Trim(body.Password, " ")

	if body.Email == "" || body.Name == "" || body.Password == "" {
		sendErrorCode(rw, ErrBadRequest, http.StatusBadRequest)
		return
	}

	if models.IsUsernameTaken(body.Name) {
		sendErrorCode(rw, ErrUsernameTaken, http.StatusBadRequest)
		return
	}

	user := models.CreateUser(body.Name, body.Password, body.Email)
	tokenResponse, err := getSessionResponseTokenForUser(user)
	if err != nil {
		sendError(rw, err)
		return
	}
	rw.Write(tokenResponse)
}

func serveAuthentication(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var body authenticationRequest
	err := decoder.Decode(&body)
	if err != nil {
		sendErrorCode(rw, err, http.StatusBadRequest)
		return
	}
	body.Name = strings.Trim(body.Name, " ")
	body.Password = strings.Trim(body.Password, " ")

	if body.Name == "" || body.Password == "" {
		sendErrorCode(rw, ErrBadRequest, http.StatusBadRequest)
		return
	}

	user := models.FindUserWithCredentials(body.Name, body.Password)

	if !user.IsResult() {
		sendErrorCode(rw, ErrInvalidCredentials, http.StatusForbidden)
		return
	}

	tokenResponse, err := getSessionResponseTokenForUser(user)
	if err != nil {
		sendError(rw, err)
		return
	}
	rw.Write(tokenResponse)
}

func getSessionResponseTokenForUser(user *models.User) ([]byte, error) {
	token, err := auth.GenerateJWTToken(*user)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(sessionResponse{
		SessionToken: token,
	})
	if err != nil {
		return nil, err
	}

	return data, nil
}
