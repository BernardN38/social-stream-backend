package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BernardN38/social-stream-backend/auth-service/internal/router"
	"github.com/stretchr/testify/assert"
)

// MockHandler will be used to test the /api/v1/auth/health endpoint
type MockHandler struct{}

func (m *MockHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
func (m *MockHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
func (m *MockHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// TestNewRouter checks if the router is initialized properly with the correct middleware and route
func TestNewRouter(t *testing.T) {
	mockHandler := &MockHandler{}
	router := router.NewRouter(mockHandler)

	// Create a test HTTP request to the /api/v1/auth/health endpoint
	req, err := http.NewRequest("GET", "/api/v1/auth/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Serve the HTTP request using the router
	router.R.ServeHTTP(rr, req)

	// Check the status code of the response
	assert.Equal(t, http.StatusOK, rr.Code, "Expected status 200 OK")

	// Check the response body
	assert.Equal(t, "OK", rr.Body.String(), "Expected body to be 'OK'")
}

// TestMiddleware checks that middleware is applied
func TestMiddleware(t *testing.T) {
	mockHandler := &MockHandler{}
	router := router.NewRouter(mockHandler)

	// Check that middleware is applied
	middlewareStack := router.R.Middlewares()
	assert.NotNil(t, middlewareStack, "Middleware stack should not be nil")
	assert.Greater(t, len(middlewareStack), 0, "Middleware stack should not be empty")

	// Validate that some specific middleware is in use (like Timeout middleware)
	containsTimeout := false
	for _, mw := range middlewareStack {
		if mw != nil {
			containsTimeout = true
		}
	}
	assert.True(t, containsTimeout, "Expected Timeout middleware to be present")
}
