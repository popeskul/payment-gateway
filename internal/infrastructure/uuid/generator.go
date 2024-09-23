package uuid

import (
	"github.com/google/uuid"
	"github.com/popeskul/payment-gateway/internal/core/ports"
)

type uuidGenerator struct{}

func NewUUIDGenerator() ports.UUIDGenerator {
	return &uuidGenerator{}
}

func (g *uuidGenerator) Generate() string {
	return uuid.New().String()
}
