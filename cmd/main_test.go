// package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/assert"

// 	transaction "pismo-dev/api/handler/v1/transactions"
// 	"pismo-dev/internal/models"
// 	"pismo-dev/internal/repository"
// 	"pismo-dev/internal/service"
// )

// type inMemoryStorage struct {
// 	accounts     map[int]models.Account
// 	transactions []models.Transaction
// }

// func (s *inMemoryStorage) CreateAccount(account models.Account) error {
// 	// Implement logic to add account to in-memory storage (e.g., s.accounts[len(s.accounts)+1] = account)
// 	s.accounts[len(s.accounts)+1] = account
// 	return nil
// }

// func (s *inMemoryStorage) GetAccount(id int) (models.Account, error) {
// 	// Implement logic to retrieve account from in-memory storage (e.g., return s.accounts[id], nil)
// 	account, ok := s.accounts[id]
// 	if !ok {
// 		return models.Account{}, nil
// 	}
// 	return account, nil
// }

// func (s *inMemoryStorage) CreateTransaction(transaction models.Transaction) error {
// 	// Implement logic to add transaction to in-memory storage (e.g., s.transactions = append(s.transactions, transaction))
// 	s.transactions = append(s.transactions, transaction)
// 	return nil
// }

// func TestCreateAccount(t *testing.T) {
// 	// Setup gin router
// 	r := gin.Default()

// 	// Create in-memory storage
// 	repos := &inMemoryRepository{
// 		accounts:     nil,
// 		transactions: nil,
// 	}
// 	services := &service.Service{}

// 	// Create handlers with in-memory storage
// 	accountHandlers := transaction.InsertAccount(&service.Service{}, repository.New())
// 	transactionHandlers := transaction.InsertTransaction(&service.Service{})

// 	// Define routes for testing
// 	r.POST("/accounts", accountHandlers)
// 	r.POST("/transactions", transactionHandlers)

// 	// Test account creation
// 	requestBody, err := json.Marshal(models.Account{DocumentNumber: "12345678900"})
// 	if err != nil {
// 		t.Errorf("Error marshalling request body: %v", err)
// 		return
// 	}
// 	req, err := http.NewRequest("POST", "/accounts", bytes.NewReader(requestBody))
// 	if err != nil {
// 		t.Errorf("Error creating request: %v", err)
// 		return
// 	}
// 	w := httptest.NewRecorder()
// 	r.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)

// 	// Add more assertions to validate response data (optional)
// 	var response models.Account
// 	err = json.Unmarshal(w.Body.Bytes(), &response)
// 	if err != nil {
// 		t.Errorf("Error unmarshalling response: %v", err)
// 		return
// 	}
// 	assert.Equal(t, "12345678900", response.DocumentNumber)

// 	// Test transaction creation (similar structure)
// 	// ... (implement similar logic for testing transaction creation)
// }
