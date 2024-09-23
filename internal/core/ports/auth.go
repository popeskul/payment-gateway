package ports

import "github.com/popeskul/payment-gateway/internal/domain"

type AuthConfig interface {
	GetAccessTokenSecret() string
	GetRefreshTokenSecret() string
	GetAccessTokenTTL() int
	GetRefreshTokenTTL() int
}

type TokenStore interface {
	StoreRefreshToken(userID string, token string) error
	DeleteRefreshToken(userID string, token string) error
	IsRefreshTokenValid(userID string, token string) bool
}

type JWTManager interface {
	GenerateTokenPair(userID string) (accessToken, refreshToken string, err error)
	ValidateAccessToken(tokenString string) (*domain.Claims, error)
	RefreshTokens(userID string, refreshToken string) (newAccessToken, newRefreshToken string, err error)
	InvalidateRefreshToken(userID string, refreshToken string) error
}

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	ComparePasswordAndHash(password, hash string) error
}
