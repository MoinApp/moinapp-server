package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type user_signup_Request struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func serve_Users_SignUp(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var body user_signup_Request
	err := decoder.Decode(&body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Create user request: %+v.\n", body)
}
