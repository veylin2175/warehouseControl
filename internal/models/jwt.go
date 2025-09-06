package models

type JWTClaims struct {
	Username string   `json:"username"`
	Role     UserRole `json:"role"`
}
