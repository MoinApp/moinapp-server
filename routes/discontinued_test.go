package routes

import (
	"net/http"
	"testing"
)

func TestDiscontinuation(t *testing.T) {
	discontinuedRoutes := []string{
		"/api",
		"/api/v1",
		"/api/v2",
		"/api/v3",
	}

	for _, url := range discontinuedRoutes {
		req, _ := http.NewRequest("GET", server.URL+url, nil)
		res, err := doRequest(req)

		if err != nil {
			t.Error(err)
		}
		res.Body.Close()

		if res.StatusCode != http.StatusGone {
			t.Errorf("Wrong status code. Expected: %v. Got: %v.", http.StatusGone, res.Status)
		}
		if res.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Wrong content-type. Expected: %v. Got: %v.", "application/json", res.Header.Get("Content-Type"))
		}
	}
}
