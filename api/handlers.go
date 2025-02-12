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

	claims, err := validateToken(r)
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
	order.UserID = claims.Username // Associa a ordem ao usuário autenticado

	// Verifica se o mercado está aberto
	if !marketOpen {
		http.Error(w, "Market is closed", http.StatusBadRequest)
		return
	}

	// Verifica se o usuário tem saldo suficiente
	userBalance, exists := balance[order.UserID]
	if !exists {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	if userBalance < order.Amount {
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
	balance[order.UserID] -= order.Amount
	previousOrders[orderKey] = order
	orderStatus[order.ID] = "Order sent"

	// Responde com sucesso
	var message string
	if order.OrderType == "limit" {
		message = "Order placed successfully (Limit Order)"
	} else {
		message = fmt.Sprintf("Successfully purchased %s for $%.2f", order.Stock, order.Amount)
	}

	// Responde com o ID da ordem e a mensagem
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message":  message,
		"order_id": order.ID,
		"balance":  balance[order.UserID],
	}
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
