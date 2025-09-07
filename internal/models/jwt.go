package models

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	Username string   `json:"username"`
	Role     UserRole `json:"role"`
	jwt.RegisteredClaims
}
