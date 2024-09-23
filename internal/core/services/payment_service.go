package services

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/popeskul/payment-gateway/internal/core/domain/payment"
	"github.com/popeskul/payment-gateway/internal/core/ports"
)

type paymentService struct {
	repo          ports.PaymentRepository
	acquiringBank ports.AcquiringBank
	logger        ports.Logger

	mu sync.RWMutex
}

func NewPaymentService(repo ports.PaymentRepository, acquiringBank ports.AcquiringBank, logger ports.Logger) ports.PaymentService {
	return &paymentService{
		repo:          repo,
		acquiringBank: acquiringBank,
		logger:        logger,
	}
}

func (s *paymentService) CreatePayment(ctx context.Context, p *payment.Payment) error {
	if p == nil {
		s.logger.Error("payment cannot be nil")
		return errors.New("payment cannot be nil")
	}

	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	p.Status = payment.PaymentStatusPending

	s.mu.Lock()
	defer s.mu.Unlock()

	return s.repo.Create(ctx, p)
}

func (s *paymentService) GetPayment(ctx context.Context, id string) (*payment.Payment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.repo.GetByID(ctx, id)
}

func (s *paymentService) UpdatePayment(ctx context.Context, p *payment.Payment) error {
	if p == nil {
		s.logger.Error("payment cannot be nil")
		return errors.New("payment cannot be nil")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	existing, err := s.repo.GetByID(ctx, p.ID)
	if err != nil {
		s.logger.Error("payment not found", "id", p.ID)
		return fmt.Errorf("payment with id %s not found", p.ID)
	}

	p.CreatedAt = existing.CreatedAt
	p.UpdatedAt = time.Now()

	return s.repo.Update(ctx, p)
}

func (s *paymentService) ListPayments(ctx context.Context, merchantID string, limit, offset int) ([]*payment.Payment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.repo.List(ctx, merchantID, limit, offset)
}

func (s *paymentService) ProcessPayment(ctx context.Context, paymentID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	p, err := s.repo.GetByID(ctx, paymentID)
	if err != nil {
		s.logger.Error("payment not found", "id", paymentID)
		return fmt.Errorf("payment with id %s not found", paymentID)
	}

	if p.Status != payment.PaymentStatusPending {
		s.logger.Error("payment is not in pending status", "id", paymentID)
		return errors.New("payment is not in pending status")
	}

	err = s.acquiringBank.ProcessPayment(ctx, p)
	if err != nil {
		s.logger.Error("Failed to process payment", "error", err, "payment_id", paymentID)
		return err
	}

	p.Status = payment.PaymentStatusCompleted
	p.UpdatedAt = time.Now()

	return s.repo.Update(ctx, p)
}
