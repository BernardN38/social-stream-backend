package application_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/BernardN38/social-stream-backend/auth-service/internal/application"
	"github.com/stretchr/testify/assert"
)

func setupEnv() {
	os.Setenv("port", "8080")
	os.Setenv("jwtSecret", "mysecret")
	os.Setenv("postgresDsn", "user=bernardn password=password host=postgres  port=5432 sslmode=disable dbname=auth_service")
	os.Setenv("dbName", "dbname")
}

// TestNewApp checks if the NewApp initializes properly
func TestNewApp(t *testing.T) {
	setupEnv()
	app := application.NewApp()

	// Assert the app router is not nil
	assert.NotNil(t, app.Router, "App router should not be nil")

	// Assert the router handler is not nil
	assert.NotNil(t, app.Router.R, "Router should have a valid http handler")
}

// TestRun checks if the app router responds as expected
func TestRun(t *testing.T) {
	setupEnv()
	// Initialize the application
	app := application.NewApp()

	// Create a test HTTP request
	req, err := http.NewRequest("GET", "/api/v1/auth/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Serve the request using the app's router
	app.Router.R.ServeHTTP(rr, req)

	// Assert the response code (assuming the root path is defined)
	// Change this based on the actual expected behavior of your app
	assert.Equal(t, http.StatusOK, rr.Code, "Expected response code 200")
}
