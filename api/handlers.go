package api

import (
	"encoding/json"
	"net/http"
)

// Função para listar usuários cadastrados (apenas admins podem acessar)
func GetUsers(w http.ResponseWriter, r *http.Request) {
	claims, err := validateToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if !claims.Admin {
		http.Error(w, "Forbidden: Only admins can access this endpoint", http.StatusForbidden)
		return
	}

	var userList []map[string]string
	for username, password := range Users {
		userData := map[string]string{
			"username": username,
			"password": password,
		}
		userList = append(userList, userData)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userList)
}
