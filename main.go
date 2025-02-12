package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
)

// Estrutura do JSON de resposta
type Response struct {
	Message string `json:"message"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // Define que o retorno é JSON
	response := Response{Message: "API está rodando!"}
	json.NewEncoder(w).Encode(response) // Converte a struct para JSON e envia
}

func main() {
	http.HandleFunc("/", handler)

	// Força a API a rodar em IPv4
	listener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatal("Erro ao iniciar servidor:", err)
	}

	fmt.Println("✅ Server running on http://0.0.0.0:8080 (IPv4 only)")

	err = http.Serve(listener, nil)
	if err != nil {
		log.Fatal("Erro ao rodar servidor:", err)
	}
}
