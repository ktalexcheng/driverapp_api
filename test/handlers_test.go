package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ktalexcheng/trailbrake_api/util"
)

func sendRequestToMockServer(t *testing.T, mg *util.MongoClient, method string, endpoint string, body io.Reader) *httptest.ResponseRecorder {
	// Start new test server
	testServer := NewTestServer(mg)
	defer testServer.Close()

	// Create new request
	req, err := http.NewRequest(method, testServer.URL+endpoint, body)
	if err != nil {
		t.Errorf("Failed to create request: %v", err)
	}

	// Create response recorder and send request
	rr := httptest.NewRecorder()
	testServer.Config.Handler.ServeHTTP(rr, req)

	return rr
}

func TestMain(m *testing.M) {
	m.Run()
}

func TestAuthTokenGet(t *testing.T) {
	// GET /auth/token should return 405
	mg := NewMockDB()
	rr := sendRequestToMockServer(t, mg, "GET", "/auth/token", nil)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Unexpected status code: got %v, expected %v", rr.Code, http.StatusMethodNotAllowed)
	}
}

func TestAuthTokenPost(t *testing.T) {
	// GET /auth/token should return 401 for invalid credentials
	mg := NewMockDB()
	rr := sendRequestToMockServer(t, mg, "POST", "/auth/token", nil)
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Unexpected status code: got %v, expected %v", rr.Code, http.StatusUnauthorized)
	}

	// GET /auth/token should return 200 for valid crednetials
	postBody, _ := json.Marshal(map[string]string{
		"email":    "test",
		"password": "1234",
	})
	postBodyBuffer := bytes.NewBuffer(postBody)

	rr = sendRequestToMockServer(t, mg, "POST", "/auth/token", postBodyBuffer)
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Unexpected status code: got %v, expected %v", rr.Code, http.StatusUnauthorized)
	}
}
