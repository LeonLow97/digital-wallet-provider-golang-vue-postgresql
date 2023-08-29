package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockService struct{}

func (s *MockService) Login(ctx context.Context, creds *Credentials) (*User, *Token, error) {
	return nil, nil, nil
}

func Test_Login(t *testing.T) {
	var mockService *MockService

	authHandler := authHandler{
		service: mockService,
	}

	// Create a mock request payload
	jsonPayload := []byte(`{"username": "testuser", "password": "testpass"}`)

	// generate request
	req, _ := http.NewRequest(http.MethodPost, "/login", strings.NewReader(string(jsonPayload)))
	req.Header.Set("Content-Type", "application/json")

	// generate response recorder
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(authHandler.Login)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, but got %d", http.StatusOK, rr.Code)
	}
}
