package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// SetupRoutes configura as rotas da API
func SetupRoutes() {
	r := mux.NewRouter()

	r.HandleFunc("/register", Register).Methods("POST")
	r.HandleFunc("/login", Login).Methods("POST")
	r.HandleFunc("/users", GetUsers).Methods("GET")
	r.HandleFunc("/buy", BuyStock).Methods("POST")
	r.HandleFunc("/order-status", ConsultaStatusOrdem).Methods("GET") // Adicionando nova rota

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
