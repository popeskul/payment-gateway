package handlers

import (
	"github.com/popeskul/payment-gateway/internal/core/ports"
)

type Handler struct {
	services   ports.Services
	logger     ports.Logger
	JWTManager ports.JWTManager
}

func NewHandler(services ports.Services, logger ports.Logger, jwtManager ports.JWTManager) *Handler {
	return &Handler{
		services:   services,
		logger:     logger,
		JWTManager: jwtManager,
	}
}
