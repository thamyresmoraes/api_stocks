package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// Simulação de banco de dados
var Users = make(map[string]string)
var UserRoles = make(map[string]bool)

// Chave secreta para JWT
var jwtKey = []byte("my_secret_key")

// validateToken verifica e valida o token JWT
// @Summary Valida um token JWT
// @Description Verifica se um token JWT é válido
// @Param Authorization header string true "Token JWT"
// @Success 200 {object} TokenClaims "Claims do token"
// @Failure 401 {string} string "Unauthorized"
// @Router /validate-token [get]
func validateToken(r *http.Request) (*Claims, error) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return nil, fmt.Errorf("missing token")
	}

	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	return claims, nil
}

// Register cria um novo usuário
// @Summary Cria um novo usuário
// @Description Registra um novo usuário na plataforma
// @Accept json
// @Produce json
// @Param user body User true "Dados do usuário"
// @Success 200 {string} string "User registered"
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Could not hash password"
// @Router /register [post]
func Register(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if _, exists := Users[user.Username]; exists {
		http.Error(w, "Username already taken", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Could not hash password", http.StatusInternalServerError)
		return
	}

	Users[user.Username] = string(hashedPassword)
	UserRoles[user.Username] = user.Admin

	// 🔹 Adicionando saldo inicial ao usuário
	balance[user.Username] = 100000.0 // Pode ajustar o saldo inicial como preferir

	fmt.Println("✅ User registered:", user.Username, "Balance:", balance[user.Username])

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("User registered")
}

// Login autentica um usuário e gera um token JWT
// @Summary Realiza login do usuário
// @Description Autentica um usuário e retorna um token JWT
// @Accept json
// @Produce json
// @Param credentials body Credentials true "Credenciais do usuário"
// @Success 200 {object} map[string]string "Token JWT"
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Invalid credentials"
// @Failure 500 {string} string "Could not create token"
// @Router /login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	storedPassword, exists := Users[creds.Username]
	if !exists {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	isAdmin := UserRoles[creds.Username]

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: creds.Username,
		Admin:    isAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

// GenerateJWT gera um token JWT válido para um usuário
// @Summary Gera um token JWT para um usuário
// @Description Retorna um token JWT válido para um usuário específico
// @Param username query string true "Nome do usuário"
// @Param isAdmin query bool true "Indica se o usuário é admin"
// @Success 200 {string} string "Token JWT gerado"
// @Failure 500 {string} string "Erro ao gerar token"
// @Router /generate-token [get]
func GenerateJWT(username string, isAdmin bool) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		Admin:    isAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
