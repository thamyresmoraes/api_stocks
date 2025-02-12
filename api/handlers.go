package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Função para listar usuários cadastrados (apenas admins podem acessar)
// @Summary Lista usuários cadastrados
// @Description Retorna a lista de usuários cadastrados com suas senhas (apenas administradores podem acessar)
// @Produce json
// @Success 200 {array} map[string]string "Lista de usuários cadastrados"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden: Only admins can access this endpoint"
// @Router /users [get]
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

// Simulação de mercado aberto
var marketOpen = true

// Simulação de saldo do usuário
var balance = map[string]float64{
	"test_user": 1000.0, // Usuário fictício com saldo inicial
}

// Simulação de pedidos já realizados
var previousOrders = make(map[string]BuyOrder)

// Simulação de status de ordens
var orderStatus = make(map[string]string)

// Função para realizar a compra de ações
// @Summary Compra ações
// @Description Permite ao usuário comprar ações, com validações de saldo e tipo de ordem
// @Accept  json
// @Produce  json
// @Param order body BuyOrder true "Dados da ordem de compra"
// @Success 200 {object} map[string]interface{} "Detalhes da compra"
// @Failure 400 {string} string "Erro ao processar a compra"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Market is closed"
// @Failure 400 {string} string "Insufficient funds"
// @Router /buy [post]
func BuyStock(w http.ResponseWriter, r *http.Request) {
	fmt.Println("🔹 Received BuyStock request")

	claims, err := validateToken(r)
	if err != nil {
		fmt.Println("❌ Unauthorized: Invalid Token")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var order BuyOrder
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		fmt.Println("❌ Invalid JSON request")
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	order.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	order.UserID = claims.Username

	fmt.Println("🔹 Order Details:", order)

	if !marketOpen {
		fmt.Println("❌ Market is closed")
		http.Error(w, "Market is closed", http.StatusBadRequest)
		return
	}

	userBalance, exists := balance[order.UserID]
	if !exists {
		fmt.Println("❌ User not found:", order.UserID)
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	if userBalance < order.Amount {
		fmt.Println("❌ Insufficient funds:", userBalance, "needed:", order.Amount)
		http.Error(w, "Insufficient funds", http.StatusBadRequest)
		return
	}

	orderKey := fmt.Sprintf("%s-%s-%s", order.UserID, order.Stock, order.OrderType)
	if _, exists := previousOrders[orderKey]; exists {
		fmt.Println("❌ Duplicate purchase request for:", orderKey)
		http.Error(w, "Duplicate purchase request", http.StatusBadRequest)
		return
	}

	balance[order.UserID] -= order.Amount
	previousOrders[orderKey] = order
	orderStatus[order.ID] = "Order sent"

	message := fmt.Sprintf("Successfully purchased %s for $%.2f", order.Stock, order.Amount)
	response := map[string]interface{}{
		"message":  message,
		"order_id": order.ID,
		"balance":  balance[order.UserID],
	}

	fmt.Println("✅ Order successful:", response)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ConsultaStatusOrdem retorna o status de uma ordem de compra enviada
// @Summary Consulta o status de uma ordem
// @Description Retorna o status atual de uma ordem de compra pelo ID da ordem
// @Accept json
// @Produce json
// @Param order_id path string true "ID da ordem"
// @Success 200 {object} map[string]string "Status da ordem"
// @Failure 400 {string} string "Ordem não encontrada"
// @Router /order-status/{order_id} [get]
func ConsultaStatusOrdem(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("order_id")
	if orderID == "" {
		http.Error(w, "Missing order ID", http.StatusBadRequest)
		return
	}

	status, exists := orderStatus[orderID]
	if !exists {
		http.Error(w, "Order not found", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"order_id": orderID, "status": status})
}

// MarketStock representa um ativo e seu preço atual
type MarketStock struct {
	Symbol string  `json:"symbol"`
	Name   string  `json:"name"`
	Price  float64 `json:"price"`
}

// Lista simulada de ações do mercado americano
var marketStocks = []MarketStock{
	{"AAPL", "Apple Inc.", 175.50},
	{"MSFT", "Microsoft Corporation", 320.75},
	{"GOOGL", "Alphabet Inc.", 140.30},
	{"AMZN", "Amazon.com Inc.", 135.20},
	{"TSLA", "Tesla Inc.", 215.90},
	{"NVDA", "NVIDIA Corporation", 470.65},
}

// GetMarketStocks retorna uma lista de ações do mercado americano e seus preços
// @Summary Retorna a lista de ativos do mercado americano e seus preços
// @Description Obtém uma lista de ações do mercado americano com seus preços de mercado
// @Accept json
// @Produce json
// @Success 200 {array} MarketStock "Lista de ações e preços"
// @Router /market-stocks [get]
func GetMarketStocks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(marketStocks)
}

// GetBalance retorna o saldo do usuário autenticado
// @Summary Retorna o saldo do usuário
// @Description Obtém o saldo disponível do usuário autenticado
// @Accept json
// @Produce json
// @Param Authorization header string true "Token JWT"
// @Success 200 {object} map[string]float64 "Saldo do usuário"
// @Failure 401 {string} string "Unauthorized"
// @Router /balance [get]
func GetBalance(w http.ResponseWriter, r *http.Request) {
	claims, err := validateToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userBalance, exists := balance[claims.Username]
	if !exists {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	fmt.Println("💰 Consulta de saldo:", claims.Username, "| Saldo atual:", userBalance)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"balance": userBalance})
}

// GetOrders retorna todas as ordens do usuário autenticado
// @Summary Lista todas as ordens do usuário
// @Description Retorna todas as ordens feitas pelo usuário autenticado
// @Produce json
// @Success 200 {array} BuyOrder "Lista de ordens"
// @Failure 401 {string} string "Unauthorized"
// @Router /orders [get]
func GetOrders(w http.ResponseWriter, r *http.Request) {
	claims, err := validateToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userOrders := []BuyOrder{}
	for _, order := range previousOrders {
		if order.UserID == claims.Username {
			userOrders = append(userOrders, order)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userOrders)
}
