package usecases_test

import (
	"context"
	"github.com/popeskul/payment-gateway/internal/core/domain/payment"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/popeskul/payment-gateway/internal/core/ports"
	"github.com/popeskul/payment-gateway/internal/core/usecases"
)

// TestProcessRefundUseCase_Execute остается без изменений

func TestNewUseCases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServices := ports.NewMockServices(ctrl)
	mockPaymentService := ports.NewMockPaymentService(ctrl)
	mockRefundService := ports.NewMockRefundService(ctrl)

	// Ожидаем, что Payments() будет вызван дважды
	mockServices.EXPECT().Payments().Return(mockPaymentService).Times(2)
	mockServices.EXPECT().Refunds().Return(mockRefundService)

	useCases := usecases.NewUseCases(mockServices)

	assert.NotNil(t, useCases)
	assert.NotNil(t, useCases.GetPaymentDetails)
	assert.NotNil(t, useCases.ProcessPayment)
	assert.NotNil(t, useCases.ProcessRefund)
}

// Дополнительные тесты для GetPaymentDetailsUseCase и ProcessPaymentUseCase

func TestGetPaymentDetailsUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPaymentService := ports.NewMockPaymentService(ctrl)
	useCase := usecases.NewGetPaymentDetailsUseCase(mockPaymentService)

	expectedPayment := &payment.Payment{ID: "payment123", Amount: 100}

	mockPaymentService.EXPECT().GetPayment(gomock.Any(), "payment123").Return(expectedPayment, nil)

	result, err := useCase.Execute(context.Background(), "payment123")

	assert.NoError(t, err)
	assert.Equal(t, expectedPayment, result)
}

func TestProcessPaymentUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPaymentService := ports.NewMockPaymentService(ctrl)
	useCase := usecases.NewProcessPaymentUseCase(mockPaymentService)

	mockPaymentService.EXPECT().ProcessPayment(gomock.Any(), "payment123").Return(nil)

	err := useCase.Execute(context.Background(), "payment123")

	assert.NoError(t, err)
}
