package routes

import (
	"encoding/json"
	"fmt"
	"github.com/MoinApp/moinapp-server/models"
	"net/http"
)

type user_signup_Request struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type session_Response struct {
	// Token for a session for this user.
	SessionToken string `json:"session_token"`
}

func serve_Users_SignUp(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var body user_signup_Request
	err := decoder.Decode(&body)
	if err != nil {
		SendAPIError(err, rw)
		return
	}

	fmt.Printf("Create user request: %+v.\n", body)
	if !models.IsUsernameTaken(body.Name) {
		models.CreateUser(body.Name, body.Password, body.Email)

		data, _ := json.Marshal(session_Response{
			SessionToken: "null",
		})
		rw.Write(data)
	}
}

func serve_Users_Auth(rw http.ResponseWriter, req *http.Request) {

}
