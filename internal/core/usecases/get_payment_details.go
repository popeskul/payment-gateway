package usecases

import (
	"context"

	"github.com/popeskul/payment-gateway/internal/core/domain/payment"
	"github.com/popeskul/payment-gateway/internal/core/ports"
)

type GetPaymentDetailsUseCase struct {
	paymentService ports.PaymentService
}

func NewGetPaymentDetailsUseCase(paymentService ports.PaymentService) *GetPaymentDetailsUseCase {
	return &GetPaymentDetailsUseCase{paymentService: paymentService}
}

func (uc *GetPaymentDetailsUseCase) Execute(ctx context.Context, paymentID string) (*payment.Payment, error) {
	return uc.paymentService.GetPayment(ctx, paymentID)
}
