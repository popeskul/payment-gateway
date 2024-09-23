package usecases

import (
	"github.com/popeskul/payment-gateway/internal/core/ports"
)

type UseCases struct {
	GetPaymentDetails *GetPaymentDetailsUseCase
	ProcessPayment    *ProcessPaymentUseCase
	ProcessRefund     *ProcessRefundUseCase
}

func NewUseCases(services ports.Services) *UseCases {
	return &UseCases{
		GetPaymentDetails: NewGetPaymentDetailsUseCase(services.Payments()),
		ProcessPayment:    NewProcessPaymentUseCase(services.Payments()),
		ProcessRefund:     NewProcessRefundUseCase(services.Refunds()),
	}
}
