package services

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/popeskul/payment-gateway/internal/core/domain/merchant"
	"github.com/popeskul/payment-gateway/internal/core/ports"
)

type merchantService struct {
	repo   ports.MerchantRepository
	logger ports.Logger

	mu sync.RWMutex
}

func NewMerchantService(repo ports.MerchantRepository, logger ports.Logger) ports.MerchantService {
	return &merchantService{
		repo:   repo,
		logger: logger,
	}
}

func (s *merchantService) CreateMerchant(ctx context.Context, m *merchant.Merchant) error {
	if m == nil {
		s.logger.Error("merchant cannot be nil")
		return errors.New("merchant cannot be nil")
	}

	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()

	s.mu.Lock()
	defer s.mu.Unlock()

	existingMerchant, err := s.repo.GetByEmail(ctx, m.Email)
	if err == nil && existingMerchant != nil {
		s.logger.Error("merchant with this email already exists")
		return errors.New("merchant with this email already exists")
	}

	return s.repo.Create(ctx, m)
}

func (s *merchantService) GetMerchant(ctx context.Context, id string) (*merchant.Merchant, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.repo.GetByID(ctx, id)
}

func (s *merchantService) UpdateMerchant(ctx context.Context, m *merchant.Merchant) error {
	if m == nil {
		s.logger.Error("merchant cannot be nil")
		return errors.New("merchant cannot be nil")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	existing, err := s.repo.GetByID(ctx, m.ID)
	if err != nil {
		s.logger.Error("merchant not found")
		return fmt.Errorf("merchant with id %s not found", m.ID)
	}

	m.CreatedAt = existing.CreatedAt
	m.UpdatedAt = time.Now()

	return s.repo.Update(ctx, m)
}

func (s *merchantService) DeleteMerchant(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.repo.Delete(ctx, id)
}

func (s *merchantService) ListMerchants(ctx context.Context, limit, offset int) ([]*merchant.Merchant, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.repo.List(ctx, limit, offset)
}
