package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BernardN38/social-stream-backend/auth-service/internal/handler"
	"github.com/BernardN38/social-stream-backend/auth-service/internal/models"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/mock"
)

// Mock service for testing
type MockService struct {
	mock.Mock
}

func (m *MockService) RegisterUser(ctx context.Context, payload models.RegisterUserPayload) error {
	args := m.Called(ctx, payload)
	return args.Error(0)
}

func (m *MockService) LoginUser(ctx context.Context, payload models.LoginUserPayload) (int, error) {
	args := m.Called(ctx, payload)
	return args.Int(0), args.Error(1)
}
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

func TestLoginUser_Success(t *testing.T) {
	mockService := new(MockService)
	jwtAuth := jwtauth.New("HS512", []byte("qwertyuiopasdfghjklzxcvbnm123456qwertyuiopasdfghjklzxcvbnm123456"), nil)
	tm := handler.NewTokenManger(jwtAuth)
	h := &handler.Handler{Service: mockService, TokenManager: tm}

	// Create a sample login payload
	payload := models.LoginUserPayload{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Convert payload to JSON
	payloadBytes, _ := json.Marshal(payload)

	// Set up the mock to expect the LoginUser call
	mockService.On("LoginUser", mock.Anything, payload).Return(1, nil)

	// Create a new request
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(payloadBytes))
	w := httptest.NewRecorder()

	// Call the handler
	h.LoginUser(w, req)

	// Check the response
	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "1", string(w.Body.Bytes())) // Check if response body contains user ID

	// Check if the cookie is set
	cookies := res.Cookies()
	assert.Equal(t, len(cookies), 1)

	var jwtCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "jwt" {
			jwtCookie = cookie
			break
		}
	}
	jwtRegex := `^[A-Za-z0-9-_]+?\.[A-Za-z0-9-_]+?\.[A-Za-z0-9-_]+$`
	assert.MatchRegex(t, jwtCookie.Value, jwtRegex)
	mockService.AssertExpectations(t)
}

func TestLoginUser_BadRequest(t *testing.T) {
	mockService := new(MockService)
	h := &handler.Handler{Service: mockService}

	// Create an invalid login payload (missing password)
	payload := `{"email":"test@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte(payload)))
	w := httptest.NewRecorder()

	// Call the handler
	h.LoginUser(w, req)

	// Check the response
	res := w.Result()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	mockService.AssertExpectations(t)
}

func TestLoginUser_UserNotFound(t *testing.T) {
	mockService := new(MockService)
	h := &handler.Handler{Service: mockService}

	// Create a sample login payload
	payload := models.LoginUserPayload{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	// Convert payload to JSON
	payloadBytes, _ := json.Marshal(payload)

	// Set up the mock to return an error (user not found)
	mockService.On("LoginUser", mock.Anything, payload).Return(0, errors.New("user not found"))

	// Create a new request
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(payloadBytes))
	w := httptest.NewRecorder()

	// Call the handler
	h.LoginUser(w, req)

	// Check the response
	res := w.Result()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	mockService.AssertExpectations(t)
}
