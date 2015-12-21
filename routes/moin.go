package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type moinRequest struct {
	Name string `json:"name"`
}

func serveMoin(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var body moinRequest
	err := decoder.Decode(&body)
	if err != nil {
		SendAPIError(err, rw)
		return
	}

	fmt.Printf("Moin request: %+v\n", body)
}
