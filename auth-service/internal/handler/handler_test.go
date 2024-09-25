package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BernardN38/social-stream-backend/auth-service/internal/handler"
)

func TestCheckHealth(t *testing.T) {
	h := handler.Handler{}
	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Create an HTTP handler from our function
	handler := http.HandlerFunc(h.CheckHealth)

	// Serve the HTTP request to the handler
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect (200 OK)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body is what we expect
	expected := "auth service up and running"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
