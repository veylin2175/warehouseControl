package handlers

import (
	"WarehouseControl/internal/lib/api/response"
	"WarehouseControl/internal/models"
	"WarehouseControl/internal/storage/postgres"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userStorage postgres.UserStorageI
	jwtSecret   string
	log         *slog.Logger
}

func NewAuthHandler(userStorage postgres.UserStorageI, jwtSecret string, log *slog.Logger) *AuthHandler {
	return &AuthHandler{
		userStorage: userStorage,
		jwtSecret:   jwtSecret,
		log:         log,
	}
}

type loginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auth.Login"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", r.Header.Get("X-Request-ID")),
	)

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Warn("invalid request body", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.Error("invalid request body"))
		return
	}

	user, err := h.userStorage.GetUserByUsername(r.Context(), req.Username)
	if err != nil {
		log.Warn("user not found", slog.String("username", req.Username))
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response.Error("invalid credentials"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		log.Warn("invalid password", slog.String("username", req.Username))
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response.Error("invalid credentials"))
		return
	}

	// Создаем JWT токен
	claims := &models.JWTClaims{
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		log.Error("failed to sign token", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.Error("internal server error"))
		return
	}

	resp := loginResponse{Token: tokenString}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		response.Response
		Data loginResponse `json:"data,omitempty"`
	}{
		Response: response.Response{Status: response.StatusOK},
		Data:     resp,
	})
}
