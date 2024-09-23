package services_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"

	"github.com/popeskul/payment-gateway/internal/core/domain/user"
	"github.com/popeskul/payment-gateway/internal/core/ports"
	"github.com/popeskul/payment-gateway/internal/core/services"
)

func TestUserService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := ports.NewMockUserRepository(ctrl)
	mockLogger := ports.NewMockLogger(ctrl)
	mockPasswordHasher := ports.NewMockPasswordHasher(ctrl)

	userService := services.NewUserService(mockRepo, mockLogger, mockPasswordHasher)

	tests := []struct {
		name          string
		req           *user.RegisterRequest
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Successful registration",
			req: &user.RegisterRequest{
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
			},
			setupMocks: func() {
				mockRepo.EXPECT().GetByEmail(gomock.Any(), "test@example.com").Return(nil, errors.New("not found"))
				mockPasswordHasher.EXPECT().HashPassword("password123").Return("hashedPassword", nil)
				mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "User already exists",
			req: &user.RegisterRequest{
				Email: "existing@example.com",
			},
			setupMocks: func() {
				mockRepo.EXPECT().GetByEmail(gomock.Any(), "existing@example.com").Return(&user.User{}, nil)
				mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any())
			},
			expectedError: errors.New("user with this email already exists"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			_, err := userService.Register(context.Background(), tt.req)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := ports.NewMockUserRepository(ctrl)
	mockLogger := ports.NewMockLogger(ctrl)
	mockPasswordHasher := ports.NewMockPasswordHasher(ctrl)

	userService := services.NewUserService(mockRepo, mockLogger, mockPasswordHasher)

	tests := []struct {
		name          string
		req           *user.LoginRequest
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Successful login",
			req: &user.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMocks: func() {
				mockRepo.EXPECT().GetByEmail(gomock.Any(), "test@example.com").Return(&user.User{PasswordHash: "hashedPassword"}, nil)
				mockPasswordHasher.EXPECT().ComparePasswordAndHash("password123", "hashedPassword").Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "User not found",
			req: &user.LoginRequest{
				Email: "nonexistent@example.com",
			},
			setupMocks: func() {
				mockRepo.EXPECT().GetByEmail(gomock.Any(), "nonexistent@example.com").Return(nil, errors.New("not found"))
				mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any())
			},
			expectedError: errors.New("invalid credentials"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			_, err := userService.Login(context.Background(), tt.req)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserService_UpdateProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := ports.NewMockUserRepository(ctrl)
	mockLogger := ports.NewMockLogger(ctrl)
	mockPasswordHasher := ports.NewMockPasswordHasher(ctrl)

	userService := services.NewUserService(mockRepo, mockLogger, mockPasswordHasher)

	tests := []struct {
		name          string
		id            string
		req           *user.UpdateProfileRequest
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Successful update",
			id:   "user123",
			req: &user.UpdateProfileRequest{
				FirstName: "John",
				LastName:  "Doe",
			},
			setupMocks: func() {
				mockRepo.EXPECT().GetByID(gomock.Any(), "user123").Return(&user.User{ID: "user123"}, nil)
				mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "User not found",
			id:   "nonexistent",
			req:  &user.UpdateProfileRequest{},
			setupMocks: func() {
				mockRepo.EXPECT().GetByID(gomock.Any(), "nonexistent").Return(nil, errors.New("not found"))
				mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any())
			},
			expectedError: errors.New("user not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			_, err := userService.UpdateProfile(context.Background(), tt.id, tt.req)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserService_ChangePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := ports.NewMockUserRepository(ctrl)
	mockLogger := ports.NewMockLogger(ctrl)
	mockPasswordHasher := ports.NewMockPasswordHasher(ctrl)

	userService := services.NewUserService(mockRepo, mockLogger, mockPasswordHasher)

	tests := []struct {
		name          string
		id            string
		req           *user.ChangePasswordRequest
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Successful password change",
			id:   "user123",
			req: &user.ChangePasswordRequest{
				OldPassword: "oldPass",
				NewPassword: "newPass",
			},
			setupMocks: func() {
				mockRepo.EXPECT().GetByID(gomock.Any(), "user123").Return(&user.User{ID: "user123", PasswordHash: "oldHash"}, nil)
				mockPasswordHasher.EXPECT().ComparePasswordAndHash("oldPass", "oldHash").Return(nil)
				mockPasswordHasher.EXPECT().HashPassword("newPass").Return("newHash", nil)
				mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Invalid old password",
			id:   "user123",
			req: &user.ChangePasswordRequest{
				OldPassword: "wrongPass",
				NewPassword: "newPass",
			},
			setupMocks: func() {
				mockRepo.EXPECT().GetByID(gomock.Any(), "user123").Return(&user.User{ID: "user123", PasswordHash: "oldHash"}, nil)
				mockPasswordHasher.EXPECT().ComparePasswordAndHash("wrongPass", "oldHash").Return(errors.New("invalid password"))
				mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any())
			},
			expectedError: errors.New("invalid old password"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := userService.ChangePassword(context.Background(), tt.id, tt.req)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
