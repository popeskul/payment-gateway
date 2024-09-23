package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/popeskul/payment-gateway/internal/core/domain/payment"
	"github.com/popeskul/payment-gateway/internal/core/ports"
	"github.com/popeskul/payment-gateway/internal/core/services"
)

func TestPaymentService_CreatePayment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := ports.NewMockPaymentRepository(ctrl)
	mockAcquiringBank := ports.NewMockAcquiringBank(ctrl)
	mockLogger := ports.NewMockLogger(ctrl)

	paymentService := services.NewPaymentService(mockRepo, mockAcquiringBank, mockLogger)

	tests := []struct {
		name          string
		payment       *payment.Payment
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Successful payment creation",
			payment: &payment.Payment{
				MerchantID: "merchant123",
				Amount:     100.0,
				Currency:   "USD",
			},
			setupMocks: func() {
				mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:    "Nil payment",
			payment: nil,
			setupMocks: func() {
				mockLogger.EXPECT().Error("payment cannot be nil")
			},
			expectedError: errors.New("payment cannot be nil"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := paymentService.CreatePayment(context.Background(), tt.payment)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				if tt.payment != nil {
					assert.Equal(t, payment.PaymentStatusPending, tt.payment.Status)
					assert.NotZero(t, tt.payment.CreatedAt)
					assert.NotZero(t, tt.payment.UpdatedAt)
				}
			}
		})
	}
}

func TestPaymentService_GetPayment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := ports.NewMockPaymentRepository(ctrl)
	mockAcquiringBank := ports.NewMockAcquiringBank(ctrl)
	mockLogger := ports.NewMockLogger(ctrl)

	paymentService := services.NewPaymentService(mockRepo, mockAcquiringBank, mockLogger)

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
				mockRepo.EXPECT().GetByID(gomock.Any(), "payment123").Return(&payment.Payment{
					ID:     "payment123",
					Amount: 100.0,
				}, nil)
			},
			expectedPayment: &payment.Payment{
				ID:     "payment123",
				Amount: 100.0,
			},
			expectedError: nil,
		},
		{
			name:      "Payment not found",
			paymentID: "nonexistent",
			setupMocks: func() {
				mockRepo.EXPECT().GetByID(gomock.Any(), "nonexistent").Return(nil, errors.New("payment not found"))
			},
			expectedPayment: nil,
			expectedError:   errors.New("payment not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			payment, err := paymentService.GetPayment(context.Background(), tt.paymentID)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Nil(t, payment)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedPayment, payment)
			}
		})
	}
}

func TestPaymentService_UpdatePayment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := ports.NewMockPaymentRepository(ctrl)
	mockAcquiringBank := ports.NewMockAcquiringBank(ctrl)
	mockLogger := ports.NewMockLogger(ctrl)

	paymentService := services.NewPaymentService(mockRepo, mockAcquiringBank, mockLogger)

	tests := []struct {
		name          string
		payment       *payment.Payment
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Successful payment update",
			payment: &payment.Payment{
				ID:     "payment123",
				Amount: 150.0,
			},
			setupMocks: func() {
				mockRepo.EXPECT().GetByID(gomock.Any(), "payment123").Return(&payment.Payment{
					ID:        "payment123",
					Amount:    100.0,
					CreatedAt: time.Now(),
				}, nil)
				mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:    "Nil payment",
			payment: nil,
			setupMocks: func() {
				mockLogger.EXPECT().Error("payment cannot be nil")
			},
			expectedError: errors.New("payment cannot be nil"),
		},
		{
			name: "Payment not found",
			payment: &payment.Payment{
				ID: "nonexistent",
			},
			setupMocks: func() {
				mockRepo.EXPECT().GetByID(gomock.Any(), "nonexistent").Return(nil, errors.New("payment not found"))
				mockLogger.EXPECT().Error("payment not found", "id", "nonexistent")
			},
			expectedError: errors.New("payment with id nonexistent not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := paymentService.UpdatePayment(context.Background(), tt.payment)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				if tt.payment != nil {
					assert.NotZero(t, tt.payment.UpdatedAt)
				}
			}
		})
	}
}

func TestPaymentService_ListPayments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := ports.NewMockPaymentRepository(ctrl)
	mockAcquiringBank := ports.NewMockAcquiringBank(ctrl)
	mockLogger := ports.NewMockLogger(ctrl)

	paymentService := services.NewPaymentService(mockRepo, mockAcquiringBank, mockLogger)

	tests := []struct {
		name           string
		merchantID     string
		limit          int
		offset         int
		setupMocks     func()
		expectedResult []*payment.Payment
		expectedError  error
	}{
		{
			name:       "Successful payments list",
			merchantID: "merchant123",
			limit:      10,
			offset:     0,
			setupMocks: func() {
				mockRepo.EXPECT().List(gomock.Any(), "merchant123", 10, 0).Return([]*payment.Payment{
					{ID: "payment1", Amount: 100.0},
					{ID: "payment2", Amount: 200.0},
				}, nil)
			},
			expectedResult: []*payment.Payment{
				{ID: "payment1", Amount: 100.0},
				{ID: "payment2", Amount: 200.0},
			},
			expectedError: nil,
		},
		{
			name:       "No payments found",
			merchantID: "merchant456",
			limit:      10,
			offset:     0,
			setupMocks: func() {
				mockRepo.EXPECT().List(gomock.Any(), "merchant456", 10, 0).Return([]*payment.Payment{}, nil)
			},
			expectedResult: []*payment.Payment{},
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			payments, err := paymentService.ListPayments(context.Background(), tt.merchantID, tt.limit, tt.offset)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Nil(t, payments)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, payments)
			}
		})
	}
}

func TestPaymentService_ProcessPayment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := ports.NewMockPaymentRepository(ctrl)
	mockAcquiringBank := ports.NewMockAcquiringBank(ctrl)
	mockLogger := ports.NewMockLogger(ctrl)

	paymentService := services.NewPaymentService(mockRepo, mockAcquiringBank, mockLogger)

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
				mockRepo.EXPECT().GetByID(gomock.Any(), "payment123").Return(&payment.Payment{
					ID:     "payment123",
					Status: payment.PaymentStatusPending,
				}, nil)
				mockAcquiringBank.EXPECT().ProcessPayment(gomock.Any(), gomock.Any()).Return(nil)
				mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:      "Payment not found",
			paymentID: "nonexistent",
			setupMocks: func() {
				mockRepo.EXPECT().GetByID(gomock.Any(), "nonexistent").Return(nil, errors.New("payment not found"))
				mockLogger.EXPECT().Error("payment not found", "id", "nonexistent")
			},
			expectedError: errors.New("payment with id nonexistent not found"),
		},
		{
			name:      "Payment not in pending status",
			paymentID: "payment456",
			setupMocks: func() {
				mockRepo.EXPECT().GetByID(gomock.Any(), "payment456").Return(&payment.Payment{
					ID:     "payment456",
					Status: payment.PaymentStatusCompleted,
				}, nil)
				mockLogger.EXPECT().Error("payment is not in pending status", "id", "payment456")
			},
			expectedError: errors.New("payment is not in pending status"),
		},
		{
			name:      "Acquiring bank processing error",
			paymentID: "payment789",
			setupMocks: func() {
				mockRepo.EXPECT().GetByID(gomock.Any(), "payment789").Return(&payment.Payment{
					ID:     "payment789",
					Status: payment.PaymentStatusPending,
				}, nil)
				mockAcquiringBank.EXPECT().ProcessPayment(gomock.Any(), gomock.Any()).Return(errors.New("processing error"))
				mockLogger.EXPECT().Error("Failed to process payment", "error", errors.New("processing error"), "payment_id", "payment789")
			},
			expectedError: errors.New("processing error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := paymentService.ProcessPayment(context.Background(), tt.paymentID)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
