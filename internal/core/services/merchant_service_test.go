package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/popeskul/payment-gateway/internal/core/domain/merchant"
	"github.com/popeskul/payment-gateway/internal/core/ports"
	"github.com/popeskul/payment-gateway/internal/core/services"
)

func TestMerchantService_CreateMerchant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := ports.NewMockMerchantRepository(ctrl)
	mockLogger := ports.NewMockLogger(ctrl)

	merchantService := services.NewMerchantService(mockRepo, mockLogger)

	tests := []struct {
		name          string
		merchant      *merchant.Merchant
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Successful merchant creation",
			merchant: &merchant.Merchant{
				Name:  "Test Merchant",
				Email: "test@example.com",
			},
			setupMocks: func() {
				mockRepo.EXPECT().GetByEmail(gomock.Any(), "test@example.com").Return(nil, errors.New("not found"))
				mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:     "Nil merchant",
			merchant: nil,
			setupMocks: func() {
				mockLogger.EXPECT().Error("merchant cannot be nil")
			},
			expectedError: errors.New("merchant cannot be nil"),
		},
		{
			name: "Merchant with existing email",
			merchant: &merchant.Merchant{
				Name:  "Test Merchant",
				Email: "existing@example.com",
			},
			setupMocks: func() {
				mockRepo.EXPECT().GetByEmail(gomock.Any(), "existing@example.com").Return(&merchant.Merchant{}, nil)
				mockLogger.EXPECT().Error("merchant with this email already exists")
			},
			expectedError: errors.New("merchant with this email already exists"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := merchantService.CreateMerchant(context.Background(), tt.merchant)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				if tt.merchant != nil {
					assert.NotZero(t, tt.merchant.CreatedAt)
					assert.NotZero(t, tt.merchant.UpdatedAt)
				}
			}
		})
	}
}

func TestMerchantService_GetMerchant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := ports.NewMockMerchantRepository(ctrl)
	mockLogger := ports.NewMockLogger(ctrl)

	merchantService := services.NewMerchantService(mockRepo, mockLogger)

	tests := []struct {
		name             string
		merchantID       string
		setupMocks       func()
		expectedMerchant *merchant.Merchant
		expectedError    error
	}{
		{
			name:       "Successful merchant retrieval",
			merchantID: "merchant123",
			setupMocks: func() {
				mockRepo.EXPECT().GetByID(gomock.Any(), "merchant123").Return(&merchant.Merchant{
					ID:   "merchant123",
					Name: "Test Merchant",
				}, nil)
			},
			expectedMerchant: &merchant.Merchant{
				ID:   "merchant123",
				Name: "Test Merchant",
			},
			expectedError: nil,
		},
		{
			name:       "Merchant not found",
			merchantID: "nonexistent",
			setupMocks: func() {
				mockRepo.EXPECT().GetByID(gomock.Any(), "nonexistent").Return(nil, errors.New("merchant not found"))
			},
			expectedMerchant: nil,
			expectedError:    errors.New("merchant not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			merchant, err := merchantService.GetMerchant(context.Background(), tt.merchantID)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Nil(t, merchant)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedMerchant, merchant)
			}
		})
	}
}

func TestMerchantService_UpdateMerchant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := ports.NewMockMerchantRepository(ctrl)
	mockLogger := ports.NewMockLogger(ctrl)

	merchantService := services.NewMerchantService(mockRepo, mockLogger)

	tests := []struct {
		name          string
		merchant      *merchant.Merchant
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Successful merchant update",
			merchant: &merchant.Merchant{
				ID:   "merchant123",
				Name: "Updated Merchant",
			},
			setupMocks: func() {
				mockRepo.EXPECT().GetByID(gomock.Any(), "merchant123").Return(&merchant.Merchant{
					ID:        "merchant123",
					Name:      "Original Merchant",
					CreatedAt: time.Now(),
				}, nil)
				mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:     "Nil merchant",
			merchant: nil,
			setupMocks: func() {
				mockLogger.EXPECT().Error("merchant cannot be nil")
			},
			expectedError: errors.New("merchant cannot be nil"),
		},
		{
			name: "Merchant not found",
			merchant: &merchant.Merchant{
				ID: "nonexistent",
			},
			setupMocks: func() {
				mockRepo.EXPECT().GetByID(gomock.Any(), "nonexistent").Return(nil, errors.New("merchant not found"))
				mockLogger.EXPECT().Error("merchant not found")
			},
			expectedError: errors.New("merchant with id nonexistent not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := merchantService.UpdateMerchant(context.Background(), tt.merchant)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				if tt.merchant != nil {
					assert.NotZero(t, tt.merchant.UpdatedAt)
				}
			}
		})
	}
}

func TestMerchantService_DeleteMerchant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := ports.NewMockMerchantRepository(ctrl)
	mockLogger := ports.NewMockLogger(ctrl)

	merchantService := services.NewMerchantService(mockRepo, mockLogger)

	tests := []struct {
		name          string
		merchantID    string
		setupMocks    func()
		expectedError error
	}{
		{
			name:       "Successful merchant deletion",
			merchantID: "merchant123",
			setupMocks: func() {
				mockRepo.EXPECT().Delete(gomock.Any(), "merchant123").Return(nil)
			},
			expectedError: nil,
		},
		{
			name:       "Merchant not found",
			merchantID: "nonexistent",
			setupMocks: func() {
				mockRepo.EXPECT().Delete(gomock.Any(), "nonexistent").Return(errors.New("merchant not found"))
			},
			expectedError: errors.New("merchant not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := merchantService.DeleteMerchant(context.Background(), tt.merchantID)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMerchantService_ListMerchants(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := ports.NewMockMerchantRepository(ctrl)
	mockLogger := ports.NewMockLogger(ctrl)

	merchantService := services.NewMerchantService(mockRepo, mockLogger)

	tests := []struct {
		name           string
		limit          int
		offset         int
		setupMocks     func()
		expectedResult []*merchant.Merchant
		expectedError  error
	}{
		{
			name:   "Successful merchants list",
			limit:  10,
			offset: 0,
			setupMocks: func() {
				mockRepo.EXPECT().List(gomock.Any(), 10, 0).Return([]*merchant.Merchant{
					{ID: "merchant1", Name: "Merchant 1"},
					{ID: "merchant2", Name: "Merchant 2"},
				}, nil)
			},
			expectedResult: []*merchant.Merchant{
				{ID: "merchant1", Name: "Merchant 1"},
				{ID: "merchant2", Name: "Merchant 2"},
			},
			expectedError: nil,
		},
		{
			name:   "No merchants found",
			limit:  10,
			offset: 0,
			setupMocks: func() {
				mockRepo.EXPECT().List(gomock.Any(), 10, 0).Return([]*merchant.Merchant{}, nil)
			},
			expectedResult: []*merchant.Merchant{},
			expectedError:  nil,
		},
		{
			name:   "Repository error",
			limit:  10,
			offset: 0,
			setupMocks: func() {
				mockRepo.EXPECT().List(gomock.Any(), 10, 0).Return(nil, errors.New("database error"))
			},
			expectedResult: nil,
			expectedError:  errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			merchants, err := merchantService.ListMerchants(context.Background(), tt.limit, tt.offset)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Nil(t, merchants)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, merchants)
			}
		})
	}
}
