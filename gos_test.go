package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/cheikhshift/gos/gos-test"
)


var result int


func Test(t *testing.T) {
	// Testing web handler generated by Go Server (GopherSauce)
	// Renamed package name of application.go and bindata.go to gostest to import generated methods
	req, err := http.NewRequest("GET", "/index", nil)
	if err != nil {
		t.Fatal(err)
	}
	var expected = `Page data`

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handle := http.HandlerFunc(gostest.MakeHandler(gostest.Handler))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.

	handle.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	//No page set so in this test the request will return error 404
	if status := rr.Code; status != http.StatusOK  {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.

	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}


func GWeb(b *testing.B) {
	req, err := http.NewRequest("GET", "/index", nil)
	if err != nil {
		b.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handle := http.HandlerFunc(gostest.MakeHandler(gostest.Handler))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.

	handle.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		b.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.

}

func BenchmarkGWeb(b *testing.B) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.

	var r int

	for n := 0; n < b.N; n++ {
		r = 0

		GWeb(b)

	}
	result = r
}
