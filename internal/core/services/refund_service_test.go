package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/popeskul/payment-gateway/internal/core/domain/payment"
	"github.com/popeskul/payment-gateway/internal/core/domain/refund"
	"github.com/popeskul/payment-gateway/internal/core/ports"
	"github.com/popeskul/payment-gateway/internal/core/services"
)

func TestRefundService_CreateRefund(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRefundRepo := ports.NewMockRefundRepository(ctrl)
	mockPaymentRepo := ports.NewMockPaymentRepository(ctrl)
	mockLogger := ports.NewMockLogger(ctrl)

	refundService := services.NewRefundService(mockRefundRepo, mockPaymentRepo, mockLogger)

	tests := []struct {
		name          string
		refund        *refund.Refund
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Successful refund creation",
			refund: &refund.Refund{
				PaymentID: "payment123",
				Amount:    50.0,
			},
			setupMocks: func() {
				mockPaymentRepo.EXPECT().GetByID(gomock.Any(), "payment123").Return(&payment.Payment{
					ID:     "payment123",
					Amount: 100.0,
					Status: payment.PaymentStatusCompleted,
				}, nil)
				mockRefundRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "Nil refund",
			refund: nil,
			setupMocks: func() {
				mockLogger.EXPECT().Error(gomock.Any())
			},
			expectedError: errors.New("refund cannot be nil"),
		},
		{
			name: "Payment not found",
			refund: &refund.Refund{
				PaymentID: "nonexistent",
				Amount:    50.0,
			},
			setupMocks: func() {
				mockPaymentRepo.EXPECT().GetByID(gomock.Any(), "nonexistent").Return(nil, errors.New("payment not found"))
				mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any())
			},
			expectedError: errors.New("failed to get payment: payment not found"),
		},
		{
			name: "Payment not completed",
			refund: &refund.Refund{
				PaymentID: "payment123",
				Amount:    50.0,
			},
			setupMocks: func() {
				mockPaymentRepo.EXPECT().GetByID(gomock.Any(), "payment123").Return(&payment.Payment{
					ID:     "payment123",
					Amount: 100.0,
					Status: payment.PaymentStatusPending,
				}, nil)
				mockLogger.EXPECT().Error(gomock.Any())
			},
			expectedError: errors.New("can only refund completed payments"),
		},
		{
			name: "Refund amount too high",
			refund: &refund.Refund{
				PaymentID: "payment123",
				Amount:    150.0,
			},
			setupMocks: func() {
				mockPaymentRepo.EXPECT().GetByID(gomock.Any(), "payment123").Return(&payment.Payment{
					ID:     "payment123",
					Amount: 100.0,
					Status: payment.PaymentStatusCompleted,
				}, nil)
				mockLogger.EXPECT().Error(gomock.Any())
			},
			expectedError: errors.New("refund amount cannot be greater than payment amount"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := refundService.CreateRefund(context.Background(), tt.refund)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRefundService_GetRefund(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRefundRepo := ports.NewMockRefundRepository(ctrl)
	mockPaymentRepo := ports.NewMockPaymentRepository(ctrl)
	mockLogger := ports.NewMockLogger(ctrl)

	refundService := services.NewRefundService(mockRefundRepo, mockPaymentRepo, mockLogger)

	tests := []struct {
		name           string
		refundID       string
		setupMocks     func()
		expectedRefund *refund.Refund
		expectedError  error
	}{
		{
			name:     "Successful refund retrieval",
			refundID: "refund123",
			setupMocks: func() {
				mockRefundRepo.EXPECT().GetByID(gomock.Any(), "refund123").Return(&refund.Refund{
					ID:     "refund123",
					Amount: 50.0,
				}, nil)
			},
			expectedRefund: &refund.Refund{
				ID:     "refund123",
				Amount: 50.0,
			},
			expectedError: nil,
		},
		{
			name:     "Refund not found",
			refundID: "nonexistent",
			setupMocks: func() {
				mockRefundRepo.EXPECT().GetByID(gomock.Any(), "nonexistent").Return(nil, errors.New("refund not found"))
			},
			expectedRefund: nil,
			expectedError:  errors.New("refund not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			refund, err := refundService.GetRefund(context.Background(), tt.refundID)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Nil(t, refund)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRefund, refund)
			}
		})
	}
}

func TestRefundService_UpdateRefund(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRefundRepo := ports.NewMockRefundRepository(ctrl)
	mockPaymentRepo := ports.NewMockPaymentRepository(ctrl)
	mockLogger := ports.NewMockLogger(ctrl)

	refundService := services.NewRefundService(mockRefundRepo, mockPaymentRepo, mockLogger)

	tests := []struct {
		name          string
		refund        *refund.Refund
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Successful refund update",
			refund: &refund.Refund{
				ID:     "refund123",
				Amount: 75.0,
			},
			setupMocks: func() {
				mockRefundRepo.EXPECT().GetByID(gomock.Any(), "refund123").Return(&refund.Refund{
					ID:        "refund123",
					Amount:    50.0,
					CreatedAt: time.Now(),
				}, nil)
				mockRefundRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "Nil refund",
			refund: nil,
			setupMocks: func() {
				mockLogger.EXPECT().Error(gomock.Any())
			},
			expectedError: errors.New("refund cannot be nil"),
		},
		{
			name: "Refund not found",
			refund: &refund.Refund{
				ID: "nonexistent",
			},
			setupMocks: func() {
				mockRefundRepo.EXPECT().GetByID(gomock.Any(), "nonexistent").Return(nil, errors.New("refund not found"))
				mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any())
			},
			expectedError: errors.New("refund with id nonexistent not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := refundService.UpdateRefund(context.Background(), tt.refund)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRefundService_ListRefunds(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRefundRepo := ports.NewMockRefundRepository(ctrl)
	mockPaymentRepo := ports.NewMockPaymentRepository(ctrl)
	mockLogger := ports.NewMockLogger(ctrl)

	refundService := services.NewRefundService(mockRefundRepo, mockPaymentRepo, mockLogger)

	tests := []struct {
		name           string
		paymentID      string
		limit          int
		offset         int
		setupMocks     func()
		expectedResult []*refund.Refund
		expectedError  error
	}{
		{
			name:      "Successful refunds list",
			paymentID: "payment123",
			limit:     10,
			offset:    0,
			setupMocks: func() {
				mockRefundRepo.EXPECT().List(gomock.Any(), "payment123", 10, 0).Return([]*refund.Refund{
					{ID: "refund1", Amount: 50.0},
					{ID: "refund2", Amount: 25.0},
				}, nil)
			},
			expectedResult: []*refund.Refund{
				{ID: "refund1", Amount: 50.0},
				{ID: "refund2", Amount: 25.0},
			},
			expectedError: nil,
		},
		{
			name:      "No refunds found",
			paymentID: "payment456",
			limit:     10,
			offset:    0,
			setupMocks: func() {
				mockRefundRepo.EXPECT().List(gomock.Any(), "payment456", 10, 0).Return([]*refund.Refund{}, nil)
			},
			expectedResult: []*refund.Refund{},
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			refunds, err := refundService.ListRefunds(context.Background(), tt.paymentID, tt.limit, tt.offset)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Nil(t, refunds)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, refunds)
			}
		})
	}
}

func TestRefundService_ProcessRefund(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRefundRepo := ports.NewMockRefundRepository(ctrl)
	mockPaymentRepo := ports.NewMockPaymentRepository(ctrl)
	mockLogger := ports.NewMockLogger(ctrl)

	refundService := services.NewRefundService(mockRefundRepo, mockPaymentRepo, mockLogger)

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
				mockRefundRepo.EXPECT().GetByID(gomock.Any(), "refund123").Return(&refund.Refund{
					ID:        "refund123",
					PaymentID: "payment123",
					Amount:    50.0,
					Status:    refund.RefundStatusPending,
				}, nil)
				mockPaymentRepo.EXPECT().GetByID(gomock.Any(), "payment123").Return(&payment.Payment{
					ID:     "payment123",
					Amount: 100.0,
				}, nil)
				mockRefundRepo.EXPECT().UpdateWithTransaction(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:     "Refund not found",
			refundID: "nonexistent",
			setupMocks: func() {
				mockRefundRepo.EXPECT().GetByID(gomock.Any(), "nonexistent").Return(nil, errors.New("refund not found"))
				mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any())
			},
			expectedError: errors.New("failed to get refund: refund not found"),
		},
		{
			name:     "Refund not in pending status",
			refundID: "refund456",
			setupMocks: func() {
				mockRefundRepo.EXPECT().GetByID(gomock.Any(), "refund456").Return(&refund.Refund{
					ID:     "refund456",
					Status: refund.RefundStatusCompleted,
				}, nil)
			},
			expectedError: errors.New("refund is not in pending status"),
		},
		{
			name:     "Payment not found",
			refundID: "refund789",
			setupMocks: func() {
				mockRefundRepo.EXPECT().GetByID(gomock.Any(), "refund789").Return(&refund.Refund{
					ID:        "refund789",
					PaymentID: "payment789",
					Status:    refund.RefundStatusPending,
				}, nil)
				mockPaymentRepo.EXPECT().GetByID(gomock.Any(), "payment789").Return(nil, errors.New("payment not found"))
				mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any())
			},
			expectedError: errors.New("failed to get payment: payment not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := refundService.ProcessRefund(context.Background(), tt.refundID)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
