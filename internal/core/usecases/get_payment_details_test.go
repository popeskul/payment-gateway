package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/popeskul/payment-gateway/internal/core/domain/payment"
	"github.com/popeskul/payment-gateway/internal/core/ports"
)

func TestGetPaymentDetailsUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPaymentService := ports.NewMockPaymentService(ctrl)
	useCase := NewGetPaymentDetailsUseCase(mockPaymentService)

	tests := []struct {
		name            string
		paymentID       string
		setupMocks      func()
		expectedPayment *payment.Payment
		expectedError   error
	}{
		{
			name:      "Successful payment retrieval",
			paymentID: "payment123",
			setupMocks: func() {
				mockPaymentService.EXPECT().GetPayment(gomock.Any(), "payment123").Return(
					&payment.Payment{
						ID:     "payment123",
						Amount: 100.0,
						Status: payment.PaymentStatusCompleted,
					}, nil)
			},
			expectedPayment: &payment.Payment{
				ID:     "payment123",
				Amount: 100.0,
				Status: payment.PaymentStatusCompleted,
			},
			expectedError: nil,
		},
		{
			name:      "Payment not found",
			paymentID: "nonexistent",
			setupMocks: func() {
				mockPaymentService.EXPECT().GetPayment(gomock.Any(), "nonexistent").Return(nil, errors.New("payment not found"))
			},
			expectedPayment: nil,
			expectedError:   errors.New("payment not found"),
		},
		{
			name:      "Empty payment ID",
			paymentID: "",
			setupMocks: func() {
				mockPaymentService.EXPECT().GetPayment(gomock.Any(), "").Return(nil, errors.New("invalid payment ID"))
			},
			expectedPayment: nil,
			expectedError:   errors.New("invalid payment ID"),
		},
		{
			name:      "Database error",
			paymentID: "payment456",
			setupMocks: func() {
				mockPaymentService.EXPECT().GetPayment(gomock.Any(), "payment456").Return(nil, errors.New("database error"))
			},
			expectedPayment: nil,
			expectedError:   errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			payment, err := useCase.Execute(context.Background(), tt.paymentID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Nil(t, payment)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, payment)
				assert.Equal(t, tt.expectedPayment, payment)
			}
		})
	}
}

func TestNewGetPaymentDetailsUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPaymentService := ports.NewMockPaymentService(ctrl)
	useCase := NewGetPaymentDetailsUseCase(mockPaymentService)

	assert.NotNil(t, useCase)
	assert.Equal(t, mockPaymentService, useCase.paymentService)
}
