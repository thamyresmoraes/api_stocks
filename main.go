package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"stock_api/docs"
	_ "stock_api/docs" // Importa os documentos gerados pelo Swagger

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Mercado aberto ou fechado
var marketOpen = true

// Saldo do usuário
var balance = 1000.0

// Para verificar compras duplicadas
var previousOrders = make(map[string]BuyOrder)

// Para armazenar o status das ordens
var orderStatus = make(map[string]string)

// Para armazenar os usuários cadastrados
var users = make(map[string]string) // Armazena login e senha (em um caso real, utilizaríamos um banco de dados)

// Secret Key para gerar o token JWT
var jwtKey = []byte("my_secret_key")

// Estrutura de dados para a compra de ações
type BuyOrder struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	Stock     string  `json:"stock"`
	Amount    float64 `json:"amount"`
	OrderType string  `json:"order_type"` // Limitada ou Mercado
}

// Estrutura para o usuário
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Estrutura para o login
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Estrutura para a resposta do token
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// Função para inicializar o Swagger
func init() {
	// Gera a documentação Swagger automaticamente
	docs.SwaggerInfo.Title = "Stock API"
	docs.SwaggerInfo.Description = "API para comprar ações"
	docs.SwaggerInfo.Version = "1.0"
}

func main() {
	r := mux.NewRouter()

	// Rota de status do mercado
	r.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Market is open"))
	}).Methods("GET")

	// Rota para compra de ações
	r.HandleFunc("/buy", buyStock).Methods("POST")

	// Rota para status da ordem
	r.HandleFunc("/order-status/{id}", getOrderStatus).Methods("GET")

	// Rota para acessar a documentação Swagger
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Rota para cadastrar usuário
	r.HandleFunc("/register", register).Methods("POST")

	// Rota para login (gera o token JWT)
	r.HandleFunc("/login", login).Methods("POST")

	// Iniciar o servidor
	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// Função de cadastro de usuário
// @Summary Cadastrar um novo usuário
// @Description Cria um novo usuário no sistema
// @Accept json
// @Produce json
// @Param user body User true "Dados do usuário"
// @Success 200 {string} string "User registered"
// @Failure 400 {string} string "Error"
// @Router /register [post]
func register(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Verifica se o usuário já existe
	if _, exists := users[user.Username]; exists {
		http.Error(w, "Username already taken", http.StatusBadRequest)
		return
	}

	// Armazena o usuário
	users[user.Username] = user.Password

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("User registered")
}

// Função de login e geração do token JWT
// @Summary Realiza o login e retorna o token JWT
// @Description Autentica o usuário e retorna o token JWT
// @Accept json
// @Produce json
// @Param credentials body Credentials true "Credenciais do usuário"
// @Success 200 {string} string "Token JWT"
// @Failure 400 {string} string "Error"
// @Router /login [post]
func login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Verifica as credenciais do usuário
	storedPassword, exists := users[creds.Username]
	if !exists || storedPassword != creds.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Cria o token JWT
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: creds.Username,
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

	// Retorna o token JWT
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

// Função para validar o token JWT
func validateToken(r *http.Request) (*Claims, error) {
	// Obtém o token do header da requisição
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return nil, fmt.Errorf("missing token")
	}

	// Analisa o token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Retorna as informações do token
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	return claims, nil
}

// Função para realizar a compra de ações
// @Summary Compra de ações
// @Description Permite ao usuário comprar ações, com validações de saldo e tipo de ordem
// @Accept json
// @Produce json
// @Param order body BuyOrder true "Compra de ações"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /buy [post]
func buyStock(w http.ResponseWriter, r *http.Request) {
	// Valida o token JWT
	_, err := validateToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Definindo o cabeçalho para permitir JSON
	w.Header().Set("Content-Type", "application/json")

	// Decodificando a requisição JSON
	var order BuyOrder
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Gera um ID único para a ordem
	order.ID = fmt.Sprintf("%d", time.Now().UnixNano())

	// Verifica se o mercado está aberto
	if !marketOpen {
		http.Error(w, "Market is closed", http.StatusBadRequest)
		return
	}

	// Verifica se o usuário tem saldo suficiente
	if balance < order.Amount {
		http.Error(w, "Insufficient funds", http.StatusBadRequest)
		return
	}

	// Verifica se a compra é duplicada (mesmo valor e tipo de ação)
	orderKey := fmt.Sprintf("%s-%s-%s", order.UserID, order.Stock, order.OrderType)
	if _, exists := previousOrders[orderKey]; exists {
		http.Error(w, "Duplicate purchase request", http.StatusBadRequest)
		return
	}

	// Processa a compra
	balance -= order.Amount
	previousOrders[orderKey] = order
	orderStatus[order.ID] = "Order sent"

	// Responde com sucesso
	var message string
	if order.OrderType == "limitada" {
		message = "Order sent"
	} else {
		message = fmt.Sprintf("Successfully purchased %s for $%.2f", order.Stock, order.Amount)
	}

	// Responde com o ID da ordem e a mensagem
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message":  message,
		"order_id": order.ID,
		"balance":  balance,
	}
	json.NewEncoder(w).Encode(response)
}

// Função para consultar o status da ordem
// @Summary Consultar o status da ordem
// @Description Permite ao usuário consultar o status de uma ordem pelo ID
// @Param id path string true "ID da ordem"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security bearer
// @Router /order-status/{id} [get]
func getOrderStatus(w http.ResponseWriter, r *http.Request) {
	// Valida o token JWT
	_, err := validateToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extrai o ID da ordem da URL
	params := mux.Vars(r)
	orderID := params["id"]

	// Verifica se a ordem existe
	status, exists := orderStatus[orderID]
	if !exists {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	// Responde com o status da ordem
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"order_id": orderID,
		"status":   status,
	}
	json.NewEncoder(w).Encode(response)
}
