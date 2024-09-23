package hasher

import (
	"github.com/popeskul/payment-gateway/internal/core/ports"
	"golang.org/x/crypto/bcrypt"
)

type bcryptPasswordHasher struct{}

func NewBcryptPasswordHasher() ports.PasswordHasher {
	return &bcryptPasswordHasher{}
}

func (h *bcryptPasswordHasher) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (h *bcryptPasswordHasher) ComparePasswordAndHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
