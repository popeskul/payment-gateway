package domain

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Claims struct {
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	jwt.RegisteredClaims
}
