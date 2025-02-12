package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"stock_api/api"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// Teste de Registro de Usuário
func TestRegisterUser(t *testing.T) {
	reqBody, _ := json.Marshal(api.User{
		Username: "test_user",
		Password: "password123",
		Admin:    false,
	})

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.Register)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "User registered")
}

// Teste de Registro com Usuário Duplicado
func TestRegisterDuplicateUser(t *testing.T) {
	api.Users["existing_user"] = "hashed_password"

	reqBody, _ := json.Marshal(api.User{
		Username: "existing_user",
		Password: "password123",
	})

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.Register)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Username already taken")
}

// Teste de Login com Sucesso
func TestLoginSuccess(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	api.Users["test_user"] = string(hashedPassword)
	api.UserRoles["test_user"] = false

	reqBody, _ := json.Marshal(api.Credentials{
		Username: "test_user",
		Password: "password123",
	})

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.Login)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "token")
}

// Teste de Login com Senha Errada
func TestLoginWrongPassword(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	api.Users["test_user"] = string(hashedPassword)

	reqBody, _ := json.Marshal(api.Credentials{
		Username: "test_user",
		Password: "wrongpassword",
	})

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.Login)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid credentials")
}
