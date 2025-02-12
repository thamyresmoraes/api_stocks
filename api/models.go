package api

import "github.com/dgrijalva/jwt-go"

// BuyOrder representa uma ordem de compra de ações
type BuyOrder struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	Stock     string  `json:"stock"`
	Amount    float64 `json:"amount"`
	OrderType string  `json:"order_type"`
}

// User representa um usuário do sistema
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Admin    bool   `json:"admin,omitempty"`
}

// Credentials representa as credenciais de login
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// TokenClaims representa os dados do token JWT sem a inclusão direta de StandardClaims
type TokenClaims struct {
	Username string `json:"username" example:"test_user"`
	Admin    bool   `json:"admin" example:"true"`
	Exp      int64  `json:"exp"`
}

// Claims representa as claims do token JWT
type Claims struct {
	Username string `json:"username"`
	Admin    bool   `json:"admin"`
	jwt.StandardClaims
}
