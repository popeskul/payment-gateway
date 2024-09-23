package usecases

import (
	"context"

	"github.com/popeskul/payment-gateway/internal/core/ports"
)

type ProcessPaymentUseCase struct {
	paymentService ports.PaymentService
}

func NewProcessPaymentUseCase(paymentService ports.PaymentService) *ProcessPaymentUseCase {
	return &ProcessPaymentUseCase{paymentService: paymentService}
}

func (uc *ProcessPaymentUseCase) Execute(ctx context.Context, paymentID string) error {
	return uc.paymentService.ProcessPayment(ctx, paymentID)
}
