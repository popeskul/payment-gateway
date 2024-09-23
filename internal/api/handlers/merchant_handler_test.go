package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/popeskul/payment-gateway/internal/core/domain/merchant"
	"github.com/popeskul/payment-gateway/internal/core/ports"
)

func TestHandler_CreateMerchant(t *testing.T) {
	tests := []struct {
		name           string
		input          merchant.Merchant
		setupMocks     func(*ports.MockServices, *ports.MockMerchantService, *ports.MockLogger)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "Success",
			input: merchant.Merchant{
				Name: "Test Merchant",
			},
			setupMocks: func(ms *ports.MockServices, mms *ports.MockMerchantService, ml *ports.MockLogger) {
				ms.EXPECT().Merchants().Return(mms)
				mms.EXPECT().CreateMerchant(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: merchant.Merchant{
				Name: "Test Merchant",
			},
		},
		{
			name: "Service Error",
			input: merchant.Merchant{
				Name: "Test Merchant",
			},
			setupMocks: func(ms *ports.MockServices, mms *ports.MockMerchantService, ml *ports.MockLogger) {
				ms.EXPECT().Merchants().Return(mms)
				mms.EXPECT().CreateMerchant(gomock.Any(), gomock.Any()).Return(errors.New("service error"))
				ml.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any())
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Failed to create merchant\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockServices := ports.NewMockServices(ctrl)
			mockMerchantService := ports.NewMockMerchantService(ctrl)
			mockLogger := ports.NewMockLogger(ctrl)

			tt.setupMocks(mockServices, mockMerchantService, mockLogger)

			h := NewHandler(mockServices, mockLogger, nil)

			body, err := json.Marshal(tt.input)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", "/merchants", bytes.NewBuffer(body))
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			h.CreateMerchant(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			if tt.expectedStatus == http.StatusCreated {
				var response merchant.Merchant
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedBody, response)
			} else {
				assert.Equal(t, tt.expectedBody, rr.Body.String())
			}
		})
	}
}

func TestHandler_GetMerchant(t *testing.T) {
	tests := []struct {
		name           string
		merchantID     string
		setupMocks     func(*ports.MockServices, *ports.MockMerchantService, *ports.MockLogger)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:       "Success",
			merchantID: "123",
			setupMocks: func(ms *ports.MockServices, mms *ports.MockMerchantService, ml *ports.MockLogger) {
				ms.EXPECT().Merchants().Return(mms)
				mms.EXPECT().GetMerchant(gomock.Any(), "123").Return(&merchant.Merchant{ID: "123", Name: "Test Merchant"}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: merchant.Merchant{
				ID:   "123",
				Name: "Test Merchant",
			},
		},
		{
			name:       "Not Found",
			merchantID: "456",
			setupMocks: func(ms *ports.MockServices, mms *ports.MockMerchantService, ml *ports.MockLogger) {
				ms.EXPECT().Merchants().Return(mms)
				mms.EXPECT().GetMerchant(gomock.Any(), "456").Return(nil, errors.New("not found"))
				ml.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "Merchant not found\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockServices := ports.NewMockServices(ctrl)
			mockMerchantService := ports.NewMockMerchantService(ctrl)
			mockLogger := ports.NewMockLogger(ctrl)

			tt.setupMocks(mockServices, mockMerchantService, mockLogger)

			h := NewHandler(mockServices, mockLogger, nil)

			req, err := http.NewRequest("GET", "/merchants/"+tt.merchantID, nil)
			require.NoError(t, err)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.merchantID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()

			h.GetMerchant(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			if tt.expectedStatus == http.StatusOK {
				var response merchant.Merchant
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedBody, response)
			} else {
				assert.Equal(t, tt.expectedBody, rr.Body.String())
			}
		})
	}
}
