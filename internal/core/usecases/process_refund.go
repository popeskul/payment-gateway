package usecases

import (
	"context"

	"github.com/popeskul/payment-gateway/internal/core/ports"
)

type ProcessRefundUseCase struct {
	refundService ports.RefundService
}

func NewProcessRefundUseCase(refundService ports.RefundService) *ProcessRefundUseCase {
	return &ProcessRefundUseCase{refundService: refundService}
}

func (uc *ProcessRefundUseCase) Execute(ctx context.Context, refundID string) error {
	return uc.refundService.ProcessRefund(ctx, refundID)
}
