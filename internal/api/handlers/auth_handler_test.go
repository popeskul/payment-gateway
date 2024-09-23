package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/popeskul/payment-gateway/internal/core/domain/user"
	"github.com/popeskul/payment-gateway/internal/core/ports"
)

func TestHandler_GetProfile(t *testing.T) {
	t.Run("Successful profile retrieval", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockServices := ports.NewMockServices(ctrl)
		mockUserService := ports.NewMockUserService(ctrl)
		mockLogger := ports.NewMockLogger(ctrl)
		mockJWTManager := ports.NewMockJWTManager(ctrl)

		h := NewHandler(mockServices, mockLogger, mockJWTManager)

		mockServices.EXPECT().Users().Return(mockUserService)

		expectedUser := &user.User{
			ID:        "user123",
			Email:     "test@example.com",
			FirstName: "John",
			LastName:  "Doe",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockUserService.EXPECT().GetUserByID(gomock.Any(), "user123").Return(expectedUser, nil)

		req, _ := http.NewRequest("GET", "/profile", nil)
		ctx := context.WithValue(req.Context(), "userID", "user123")
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		h.GetProfile(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var responseUser user.User
		err := json.Unmarshal(rr.Body.Bytes(), &responseUser)
		assert.NoError(t, err)

		assert.Equal(t, expectedUser.ID, responseUser.ID)
		assert.Equal(t, expectedUser.Email, responseUser.Email)
		assert.Equal(t, expectedUser.FirstName, responseUser.FirstName)
		assert.Equal(t, expectedUser.LastName, responseUser.LastName)
	})

	t.Run("Error retrieving profile", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockServices := ports.NewMockServices(ctrl)
		mockUserService := ports.NewMockUserService(ctrl)
		mockLogger := ports.NewMockLogger(ctrl)
		mockJWTManager := ports.NewMockJWTManager(ctrl)

		h := NewHandler(mockServices, mockLogger, mockJWTManager)

		mockServices.EXPECT().Users().Return(mockUserService)
		mockUserService.EXPECT().GetUserByID(gomock.Any(), "user123").Return(nil, errors.New("user not found"))
		mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any())

		req, _ := http.NewRequest("GET", "/profile", nil)
		ctx := context.WithValue(req.Context(), "userID", "user123")
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		h.GetProfile(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		var responseBody map[string]string
		json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.Equal(t, "Failed to get user profile", responseBody["error"])
	})
}
