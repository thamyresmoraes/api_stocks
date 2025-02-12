package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"stock_api/api"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestIntegration(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/register", api.Register).Methods("POST")
	router.HandleFunc("/login", api.Login).Methods("POST")
	router.HandleFunc("/users", api.GetUsers).Methods("GET")

	// Criar usuário
	reqBody, _ := json.Marshal(api.User{
		Username: "integration_user",
		Password: "test1234",
		Admin:    true,
	})
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Login do usuário
	reqBody, _ = json.Marshal(api.Credentials{
		Username: "integration_user",
		Password: "test1234",
	})
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	var tokenResponse map[string]string
	json.Unmarshal(rr.Body.Bytes(), &tokenResponse)
	token := tokenResponse["token"]

	// Teste de acesso ao endpoint /users
	req, _ = http.NewRequest("GET", "/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
