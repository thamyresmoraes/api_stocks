package api_test

import (
	"net/http"
	"net/http/httptest"
	"stock_api/api"
	"testing"

	"github.com/stretchr/testify/assert"
)

func generateTestToken(username string, isAdmin bool) string {
	token, err := api.GenerateJWT(username, isAdmin) // ✅ Captura os dois retornos
	if err != nil {
		return "Erro ao gerar token"
	}
	return "Bearer " + token
}

func TestGetUsersUnauthorized(t *testing.T) {
	token := generateTestToken("test_user", false)

	req, _ := http.NewRequest("GET", "/users", nil)
	req.Header.Set("Authorization", token)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.GetUsers)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
}

func TestGetUsersAuthorized(t *testing.T) {
	token := generateTestToken("admin_user", true)

	req, _ := http.NewRequest("GET", "/users", nil)
	req.Header.Set("Authorization", token)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.GetUsers)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
