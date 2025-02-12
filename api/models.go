package api

import (
	"github.com/dgrijalva/jwt-go"
)

// Estruturas de dados
type BuyOrder struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	Stock     string  `json:"stock"`
	Amount    float64 `json:"amount"`
	OrderType string  `json:"order_type"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Admin    bool   `json:"admin,omitempty"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	Admin    bool   `json:"admin"`
	jwt.StandardClaims
}
