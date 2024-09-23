package services

import (
	"context"
	"errors"
	"time"

	"github.com/popeskul/payment-gateway/internal/core/domain/user"
	"github.com/popeskul/payment-gateway/internal/core/ports"
)

type userService struct {
	repo           ports.UserRepository
	logger         ports.Logger
	passwordHasher ports.PasswordHasher
}

func NewUserService(repo ports.UserRepository, logger ports.Logger, passwordHasher ports.PasswordHasher) ports.UserService {
	return &userService{
		repo:           repo,
		logger:         logger,
		passwordHasher: passwordHasher,
	}
}

func (s *userService) Register(ctx context.Context, req *user.RegisterRequest) (*user.User, error) {
	existingUser, err := s.repo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		s.logger.Error("User with this email already exists", "email", req.Email)
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := s.passwordHasher.HashPassword(req.Password)
	if err != nil {
		s.logger.Error("Failed to hash password", "error", err)
		return nil, errors.New("failed to create user")
	}

	newUser := &user.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = s.repo.Create(ctx, newUser)
	if err != nil {
		s.logger.Error("Failed to create user", "error", err)
		return nil, errors.New("failed to create user")
	}

	return newUser, nil
}

func (s *userService) Login(ctx context.Context, req *user.LoginRequest) (*user.User, error) {
	u, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Error("User not found", "email", req.Email)
		return nil, errors.New("invalid credentials")
	}

	err = s.passwordHasher.ComparePasswordAndHash(req.Password, u.PasswordHash)
	if err != nil {
		s.logger.Error("Invalid password", "email", req.Email)
		return nil, errors.New("invalid credentials")
	}

	return u, nil
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*user.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *userService) UpdateProfile(ctx context.Context, id string, req *user.UpdateProfileRequest) (*user.User, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("User not found", "id", id)
		return nil, errors.New("user not found")
	}

	u.FirstName = req.FirstName
	u.LastName = req.LastName
	u.UpdatedAt = time.Now()

	err = s.repo.Update(ctx, u)
	if err != nil {
		s.logger.Error("Failed to update user", "error", err)
		return nil, errors.New("failed to update user")
	}

	return u, nil
}

func (s *userService) ChangePassword(ctx context.Context, id string, req *user.ChangePasswordRequest) error {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("User not found", "id", id)
		return errors.New("user not found")
	}

	err = s.passwordHasher.ComparePasswordAndHash(req.OldPassword, u.PasswordHash)
	if err != nil {
		s.logger.Error("Invalid old password", "id", id)
		return errors.New("invalid old password")
	}

	hashedPassword, err := s.passwordHasher.HashPassword(req.NewPassword)
	if err != nil {
		s.logger.Error("Failed to hash new password", "error", err)
		return errors.New("failed to change password")
	}

	u.PasswordHash = hashedPassword
	u.UpdatedAt = time.Now()

	err = s.repo.Update(ctx, u)
	if err != nil {
		s.logger.Error("Failed to update user password", "error", err)
		return errors.New("failed to change password")
	}

	return nil
}
