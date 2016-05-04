package routes

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/MoinApp/moinapp-server/models"
)

var (
	server           *httptest.Server
	defaultTransport http.RoundTripper
)

// Setup
// TestMain makes the setup for the tests to come.
func TestMain(m *testing.M) {
	models.InitDB(false)

	server = httptest.NewServer(CreateRouter(false))
	defer server.Close()

	defaultTransport = &http.Transport{}

	os.Exit(m.Run())
}
func path(path string) string {
	a := fmt.Sprintf("%v%v", server.URL, path)
	return a
}
func doRequest(req *http.Request) (*http.Response, error) {
	return defaultTransport.RoundTrip(req)
}

// Test Methods

func TestCreateServerOnDefaultPort(t *testing.T) {
	addr := StartListening(nil, nil)

	if result, port := testPortWithHostPort(defaultPort, addr.String()); !result {
		t.Errorf("Wrong default port. Expected: %v. Got: %v.", defaultPort, port)
	}
}
func TestCreateServerOnSpecificPort(t *testing.T) {
	const p uint16 = 8008
	os.Setenv("PORT", fmt.Sprintf("%v", p))
	addr := StartListening(nil, nil)

	if result, port := testPortWithHostPort(p, addr.String()); !result {
		t.Errorf("Wrong default port. Expected: %v. Got: %v.", p, port)
	}
}
func testPortWithHostPort(port uint16, hostPort string) (bool, string) {
	_, usedPortString, err := net.SplitHostPort(hostPort)
	if err != nil {
		return false, ""
	}

	usedPort, err := strconv.ParseUint(usedPortString, 10, 16)
	if err != nil {
		return false, usedPortString
	}

	return (uint16(usedPort) == port), usedPortString
}

func TestRootRedirects(t *testing.T) {
	req, _ := http.NewRequest("GET", path("/"), nil)
	res, err := doRequest(req)

	if err != nil {
		t.Error(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusFound {
		t.Errorf("Wrong status code. Expected: %v. Got: %v.", http.StatusFound, res.StatusCode)
	}
	if res.Header.Get("Location") != homeRedirectURL {
		t.Fatal("Wrong redirect.", "Expected:", homeRedirectURL, "Given:", res.Header.Get("Location"))
	}
}

func TestNewestApi(t *testing.T) {
	req, _ := http.NewRequest("GET", path("/v4"), nil)
	res, err := doRequest(req)

	if err != nil {
		t.Error(err)
	}
	res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Wrong status code. Expected: %v. Got: %v.", http.StatusNotFound, res.Status)
	}
}
