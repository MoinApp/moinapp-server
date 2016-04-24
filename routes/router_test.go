package routes

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/MoinApp/moinapp-server/models"
)

var (
	server           *httptest.Server
	defaultTransport http.RoundTripper = &http.Transport{}
)

func TestMain(m *testing.M) {
	log.Printf("Starting server...")
	models.InitDB(false)

	server = httptest.NewServer(CreateRouter(false))
	defer server.Close()

	log.Printf("Running on %q.", server.URL)

	os.Exit(m.Run())
}
func path(path string) string {
	a := fmt.Sprintf("%v%v", server.URL, path)
	fmt.Println(a)
	return a
}

func TestRootRedirectsToImage(t *testing.T) {
	req, _ := http.NewRequest("GET", server.URL+"/", nil)

	res, err := defaultTransport.RoundTrip(req)

	if err != nil {
		t.Error(err)
	}
	defer res.Body.Close()

	// this failes because the headers are all messed-up.
	// Why?
	if res.StatusCode != http.StatusFound {
		t.Errorf("Wrong status code. Expected: %v. Got: %v.", http.StatusFound, res.StatusCode)
	}
	if res.Header.Get("Location") != homeRedirectURL {
		t.Fatal("Wrong redirect.", "Expected:", homeRedirectURL, "Given:", res.Header.Get("Location"))
	}
}

func TestDiscontinuation(t *testing.T) {
	discontinuedRoutes := []string{
		"/api",
		"/api/v1",
		"/api/v2",
		"/api/v3",
	}

	for _, url := range discontinuedRoutes {
		req, _ := http.NewRequest("GET", server.URL+url, nil)

		res, err := defaultTransport.RoundTrip(req)

		if err != nil {
			t.Error(err)
		}
		res.Body.Close()

		if res.StatusCode != http.StatusGone {
			t.Errorf("Wrong status code. Expected: %v. Got: %v.", http.StatusGone, res.Status)
		}
	}
}

func TestNewestApi(t *testing.T) {
	req, _ := http.NewRequest("GET", server.URL+"/v4/", nil)

	res, err := defaultTransport.RoundTrip(req)

	if err != nil {
		t.Error(err)
	}
	res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Wrong status code. Expected: %v. Got: %v.", http.StatusNotFound, res.Status)
	}

}
