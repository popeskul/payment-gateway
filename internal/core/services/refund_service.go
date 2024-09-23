package services

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/popeskul/payment-gateway/internal/core/domain/payment"
	"github.com/popeskul/payment-gateway/internal/core/domain/refund"
	"github.com/popeskul/payment-gateway/internal/core/ports"
)

type refundService struct {
	refundRepo  ports.RefundRepository
	paymentRepo ports.PaymentRepository
	logger      ports.Logger

	mu sync.RWMutex
}

func NewRefundService(refundRepo ports.RefundRepository, paymentRepo ports.PaymentRepository, logger ports.Logger) ports.RefundService {
	return &refundService{
		refundRepo:  refundRepo,
		paymentRepo: paymentRepo,
		logger:      logger,
	}
}

func (s *refundService) CreateRefund(ctx context.Context, r *refund.Refund) error {
	if r == nil {
		s.logger.Error("refund cannot be nil")
		return errors.New("refund cannot be nil")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	p, err := s.paymentRepo.GetByID(ctx, r.PaymentID)
	if err != nil {
		s.logger.Error("failed to get payment", "error", err)
		return fmt.Errorf("failed to get payment: %w", err)
	}

	if p.Status != payment.PaymentStatusCompleted {
		s.logger.Error("can only refund completed payments")
		return errors.New("can only refund completed payments")
	}

	if r.Amount > p.Amount {
		s.logger.Error("refund amount cannot be greater than payment amount")
		return errors.New("refund amount cannot be greater than payment amount")
	}

	r.CreatedAt = time.Now()
	r.UpdatedAt = time.Now()
	r.Status = refund.RefundStatusPending

	return s.refundRepo.Create(ctx, r)
}

func (s *refundService) GetRefund(ctx context.Context, id string) (*refund.Refund, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.refundRepo.GetByID(ctx, id)
}

func (s *refundService) UpdateRefund(ctx context.Context, r *refund.Refund) error {
	if r == nil {
		s.logger.Error("refund cannot be nil")
		return errors.New("refund cannot be nil")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	existing, err := s.refundRepo.GetByID(ctx, r.ID)
	if err != nil {
		s.logger.Error("refund not found", "id", r.ID)
		return fmt.Errorf("refund with id %s not found", r.ID)
	}

	r.CreatedAt = existing.CreatedAt
	r.UpdatedAt = time.Now()

	return s.refundRepo.Update(ctx, r)
}

func (s *refundService) ListRefunds(ctx context.Context, paymentID string, limit, offset int) ([]*refund.Refund, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.refundRepo.List(ctx, paymentID, limit, offset)
}

func (s *refundService) ProcessRefund(ctx context.Context, refundID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	r, err := s.refundRepo.GetByID(ctx, refundID)
	if err != nil {
		s.logger.Error("failed to get refund", "error", err)
		return fmt.Errorf("failed to get refund: %w", err)
	}

	if r.Status != refund.RefundStatusPending {
		return errors.New("refund is not in pending status")
	}

	p, err := s.paymentRepo.GetByID(ctx, r.PaymentID)
	if err != nil {
		s.logger.Error("failed to get payment", "error", err)
		return fmt.Errorf("failed to get payment: %w", err)
	}

	r.Status = refund.RefundStatusCompleted
	r.UpdatedAt = time.Now()

	p.Amount -= r.Amount
	p.UpdatedAt = time.Now()

	return s.refundRepo.UpdateWithTransaction(ctx, r, p)
}
