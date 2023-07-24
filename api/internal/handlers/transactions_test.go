package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/LeonLow97/internal/transactions"
	"github.com/LeonLow97/internal/utils"

	"github.com/dgrijalva/jwt-go"
)

// Stub the Transaction Service
type MockService struct {
	ReturnServiceError    bool
	ReturnRepositoryError bool
}

func (s *MockService) GetTransactions(ctx context.Context, username string) (*transactions.Transactions, error) {
	return &transactions.Transactions{}, nil
}

func (s *MockService) CreateTransaction(ctx context.Context, senderName, beneficiaryName, beneficiaryNumber, amountTransferredCurrency, amountTransferred string) error {
	if s.ReturnServiceError {
		return utils.ServiceError{Message: "UNIT TESTING SERVICE ERROR. DO NOT BE ALARMED. THIS IS NOT AN ACTUAL ERROR."}
	}
	if s.ReturnRepositoryError {
		return utils.RepositoryError{Message: "UNIT TESTING REPOSITORY ERROR. DO NOT BE ALARMED. THIS IS NOT AN ACTUAL ERROR."}
	}
	return nil
}

func setupTest(t *testing.T, mockService *MockService, reqBody string) *httptest.ResponseRecorder {
	transactionHandler := transactionHandler{
		service: mockService,
	}

	req, _ := http.NewRequest(http.MethodPost, "/transaction", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(transactionHandler.CreateTransaction)

	expirationTime := time.Now().Add(60 * time.Minute)
	claims := &Claims{
		Username: "Alice",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtKey := []byte("secret_key")
	tokenString, _ := token.SignedString(jwtKey)
	bearerHeader := "Bearer " + tokenString
	req.Header.Set("Authorization", bearerHeader)

	handler.ServeHTTP(rr, req)

	return rr
}

func TestCreateTransaction(t *testing.T) {
	mockService := &MockService{
		ReturnServiceError:    false,
		ReturnRepositoryError: false,
	}

	reqBody := `{"beneficiary_name": "John Doe", "mobile_number": "123456789", "amount_transferred": "100.0", "amount_transferred_currency": "USD", "amount_received": 95.0, "amount_received_currency": "EUR"}`

	rr := setupTest(t, mockService, reqBody)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status %d, but got %d", http.StatusCreated, rr.Code)
	}
}

func TestCreateTransaction_ServiceError(t *testing.T) {
	mockService := &MockService{
		ReturnServiceError:    true,
		ReturnRepositoryError: false,
	}

	reqBody := `{"beneficiary_name": "John Doe", "mobile_number": "123456789", "amount_transferred": "100.0", "amount_transferred_currency": "USD", "amount_received": 95.0, "amount_received_currency": "EUR"}`

	rr := setupTest(t, mockService, reqBody)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, but got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestCreateTransaction_RepositoryError(t *testing.T) {
	mockService := &MockService{
		ReturnServiceError:    false,
		ReturnRepositoryError: true,
	}

	reqBody := `{"beneficiary_name": "John Doe", "mobile_number": "123456789", "amount_transferred": "100.0", "amount_transferred_currency": "USD", "amount_received": 95.0, "amount_received_currency": "EUR"}`

	rr := setupTest(t, mockService, reqBody)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, but got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestGetTransactions(t *testing.T) {

	mockService := &MockService{}

	transactionHandler := transactionHandler{
		service: mockService,
	}

	req, _ := http.NewRequest(http.MethodGet, "/transactions", nil)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(transactionHandler.GetTransactions)

	expirationTime := time.Now().Add(60 * time.Minute)
	claims := &Claims{
		Username: "Alice",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtKey := []byte("secret_key")
	tokenString, _ := token.SignedString(jwtKey)
	bearerHeader := "Bearer " + tokenString
	req.Header.Set("Authorization", bearerHeader)

	handler.ServeHTTP(rr, req)

	expectedStatusCode := http.StatusOK
	if rr.Code != expectedStatusCode {
		t.Errorf("%s: returned wrong status code; expected %d but got %d", "Status Code", expectedStatusCode, rr.Code)
	}

	expectedBody := `{"transactions":null}`
	if strings.TrimSpace(rr.Body.String()) != expectedBody {
		t.Errorf("Unexpected response body: expected %v, got %v", rr.Body.String(), expectedBody)
	}

}
