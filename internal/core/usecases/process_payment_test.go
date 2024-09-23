package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/popeskul/payment-gateway/internal/core/ports"
)

func TestProcessRefundUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRefundService := ports.NewMockRefundService(ctrl)
	useCase := NewProcessRefundUseCase(mockRefundService)

	tests := []struct {
		name          string
		refundID      string
		setupMocks    func()
		expectedError error
	}{
		{
			name:     "Successful refund processing",
			refundID: "refund123",
			setupMocks: func() {
				mockRefundService.EXPECT().ProcessRefund(gomock.Any(), "refund123").Return(nil)
			},
			expectedError: nil,
		},
		{
			name:     "Refund processing error",
			refundID: "refund456",
			setupMocks: func() {
				mockRefundService.EXPECT().ProcessRefund(gomock.Any(), "refund456").Return(errors.New("processing error"))
			},
			expectedError: errors.New("processing error"),
		},
		{
			name:     "Empty refund ID",
			refundID: "",
			setupMocks: func() {
				mockRefundService.EXPECT().ProcessRefund(gomock.Any(), "").Return(errors.New("invalid refund ID"))
			},
			expectedError: errors.New("invalid refund ID"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := useCase.Execute(context.Background(), tt.refundID)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Здесь можно добавить тесты для ProcessPaymentUseCase, если они еще не реализованы

func TestProcessPaymentUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPaymentService := ports.NewMockPaymentService(ctrl)
	useCase := NewProcessPaymentUseCase(mockPaymentService)

	tests := []struct {
		name          string
		paymentID     string
		setupMocks    func()
		expectedError error
	}{
		{
			name:      "Successful payment processing",
			paymentID: "payment123",
			setupMocks: func() {
				mockPaymentService.EXPECT().ProcessPayment(gomock.Any(), "payment123").Return(nil)
			},
			expectedError: nil,
		},
		{
			name:      "Payment processing error",
			paymentID: "payment456",
			setupMocks: func() {
				mockPaymentService.EXPECT().ProcessPayment(gomock.Any(), "payment456").Return(errors.New("processing error"))
			},
			expectedError: errors.New("processing error"),
		},
		{
			name:      "Empty payment ID",
			paymentID: "",
			setupMocks: func() {
				mockPaymentService.EXPECT().ProcessPayment(gomock.Any(), "").Return(errors.New("invalid payment ID"))
			},
			expectedError: errors.New("invalid payment ID"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := useCase.Execute(context.Background(), tt.paymentID)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
