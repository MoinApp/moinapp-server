package routes

import (
	"encoding/json"
	"fmt"
	"github.com/MoinApp/moinapp-server/models"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
	"strconv"
)

type userResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func newUserResponse(userModel *models.User) userResponse {
	modelID := strconv.Itoa(int(userModel.ID))

	return userResponse{
		ID:    modelID,
		Name:  userModel.Name,
		Email: userModel.Email,
	}
}

func serveGetUserProfile(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	username := vars["username"]

	user := models.FindUserByName(username)
	if !user.IsResult() {
		fmt.Printf("Requested user profile for \"%v\": No results found.", username)
		// TODO error message
		return
	}

	profile := newUserResponse(user)
	fmt.Printf("Requested user profile of \"%v\": %+v\n", username, profile)

	response, _ := json.Marshal(profile)
	rw.Write(response)
}

func serveSearchUser(rw http.ResponseWriter, req *http.Request) {
	uri, err := url.Parse(req.RequestURI)
	if err != nil {
		SendAPIError(err, rw)
		return
	}
	query := uri.Query()
	username := query.Get("username")

	users := models.FindUsersByName(username)
	profiles := make([]userResponse, len(users))

	for i, user := range users {
		profiles[i] = newUserResponse(user)
	}

	fmt.Printf("Searched for user \"%v\": %v result(s).\n", username, len(profiles))

	response, _ := json.Marshal(profiles)
	rw.Write(response)
}

func serveRecentUsers(rw http.ResponseWriter, req *http.Request) {
	fmt.Printf("Requested recents...\n")
}
