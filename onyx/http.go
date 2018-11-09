package onyx

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Create a new http request with the
// supplied parameters.
func NewRequest(method, path string) (*http.Request, error) {
	req, err := http.NewRequest(method, path, nil)

	return req, err
}

// Create a new http request, with a request body, with the
// supplied parameters.
func NewRequestWithBody(method, path, body string) (*http.Request, error) {
	data := strings.NewReader(body)
	req, err := http.NewRequest(method, path, data)

	return req, err
}

// Sets a header value for supplied http request.
func Header(req *http.Request, key, val string) {
	req.Header.Set(key, val)
}

// Check to see if handler returns code 200
// response
func Handle(req *http.Request, h http.HandlerFunc, t *testing.T) {

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(h)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
