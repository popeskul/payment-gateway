package auth

import (
	"errors"
	"github.com/popeskul/payment-gateway/internal/domain"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/popeskul/payment-gateway/internal/core/ports"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
)

type JWTManager struct {
	config     ports.AuthConfig
	tokenStore ports.TokenStore
}

func NewJWTManager(config ports.AuthConfig, tokenStore ports.TokenStore) *JWTManager {
	return &JWTManager{
		config:     config,
		tokenStore: tokenStore,
	}
}

func (m *JWTManager) GenerateTokenPair(userID string) (accessToken, refreshToken string, err error) {
	accessToken, err = m.generateAccessToken(userID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = m.generateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	err = m.tokenStore.StoreRefreshToken(userID, refreshToken)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (m *JWTManager) ValidateAccessToken(tokenString string) (*domain.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.config.GetAccessTokenSecret()), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*domain.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

func (m *JWTManager) RefreshTokens(userID string, refreshToken string) (newAccessToken, newRefreshToken string, err error) {
	if !m.tokenStore.IsRefreshTokenValid(userID, refreshToken) {
		return "", "", ErrInvalidToken
	}

	err = m.tokenStore.DeleteRefreshToken(userID, refreshToken)
	if err != nil {
		return "", "", err
	}

	return m.GenerateTokenPair(userID)
}

func (m *JWTManager) InvalidateRefreshToken(userID string, refreshToken string) error {
	return m.tokenStore.DeleteRefreshToken(userID, refreshToken)
}

func (m *JWTManager) generateAccessToken(userID string) (string, error) {
	claims := &domain.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(m.config.GetAccessTokenTTL()) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.config.GetAccessTokenSecret()))
}

func (m *JWTManager) generateRefreshToken(userID string) (string, error) {
	claims := &domain.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(m.config.GetRefreshTokenTTL()) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.config.GetRefreshTokenSecret()))
}
