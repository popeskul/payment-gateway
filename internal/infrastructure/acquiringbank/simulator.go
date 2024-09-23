package acquiringbank

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/popeskul/payment-gateway/internal/core/domain/payment"
	"github.com/popeskul/payment-gateway/internal/core/domain/refund"
	"github.com/popeskul/payment-gateway/internal/core/ports"
)

type acquiringBankSimulator struct {
	logger          ports.Logger
	processingDelay time.Duration
	failureRate     float64
	randomGenerator *rand.Rand
}

func NewAcquiringBankSimulator(logger ports.Logger, processingDelay time.Duration, failureRate float64) ports.AcquiringBank {
	return &acquiringBankSimulator{
		logger:          logger,
		processingDelay: processingDelay,
		failureRate:     failureRate,
		randomGenerator: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *acquiringBankSimulator) ProcessPayment(ctx context.Context, p *payment.Payment) error {
	s.logger.Info("Processing payment", "payment_id", p.ID, "amount", p.Amount, "currency", p.Currency)

	// Simulate processing delay
	select {
	case <-time.After(s.processingDelay):
	case <-ctx.Done():
		return ctx.Err()
	}

	// Simulate random failures
	if s.randomGenerator.Float64() < s.failureRate {
		s.logger.Warn("Payment processing failed", "payment_id", p.ID)
		return fmt.Errorf("payment processing failed")
	}

	s.logger.Info("Payment processed successfully", "payment_id", p.ID)
	return nil
}

func (s *acquiringBankSimulator) ProcessRefund(ctx context.Context, r *refund.Refund) error {
	s.logger.Info("Processing refund", "refund_id", r.ID, "payment_id", r.PaymentID, "amount", r.Amount)

	// Simulate processing delay
	select {
	case <-time.After(s.processingDelay):
	case <-ctx.Done():
		return ctx.Err()
	}

	// Simulate random failures
	if s.randomGenerator.Float64() < s.failureRate {
		s.logger.Warn("Refund processing failed", "refund_id", r.ID)
		return fmt.Errorf("refund processing failed")
	}

	s.logger.Info("Refund processed successfully", "refund_id", r.ID)
	return nil
}

// SetProcessingDelay allows to configure processing delay
func (s *acquiringBankSimulator) SetProcessingDelay(delay time.Duration) {
	s.processingDelay = delay
}

// SetFailureRate allows to configure failure rate
func (s *acquiringBankSimulator) SetFailureRate(rate float64) {
	s.failureRate = rate
}
