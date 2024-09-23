package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/popeskul/payment-gateway/internal/core/domain/user"
	"github.com/popeskul/payment-gateway/internal/infrastructure/metrics"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var registerRequest user.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	newUser, err := h.services.Users().Register(r.Context(), &registerRequest)
	if err != nil {
		h.logger.Error("Failed to register user", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to register user")
		return
	}

	accessToken, refreshToken, err := h.JWTManager.GenerateTokenPair(newUser.ID)
	if err != nil {
		h.logger.Error("Failed to generate token pair", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to generate tokens")
		return
	}

	h.respondJSON(w, http.StatusCreated, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest user.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := h.services.Users().Login(r.Context(), &loginRequest)
	if err != nil {
		h.logger.Error("Failed to login", "error", err)
		h.respondError(w, http.StatusUnauthorized, "Invalid credentials")
		metrics.AuthenticationAttempts.WithLabelValues("failed").Inc()
		return
	}

	accessToken, refreshToken, err := h.JWTManager.GenerateTokenPair(user.ID)
	if err != nil {
		h.logger.Error("Failed to generate token pair", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to generate tokens")
		metrics.AuthenticationAttempts.WithLabelValues("failed").Inc()
		return
	}

	metrics.AuthenticationAttempts.WithLabelValues("success").Inc()

	h.respondJSON(w, http.StatusOK, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var refreshRequest struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&refreshRequest); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	claims, err := h.JWTManager.ValidateAccessToken(refreshRequest.RefreshToken)
	if err != nil {
		h.respondError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	newAccessToken, newRefreshToken, err := h.JWTManager.RefreshTokens(claims.UserID, refreshRequest.RefreshToken)
	if err != nil {
		h.logger.Error("Failed to refresh tokens", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to refresh tokens")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	refreshToken := r.Header.Get("X-Refresh-Token")

	if err := h.JWTManager.InvalidateRefreshToken(userID, refreshToken); err != nil {
		h.logger.Error("Failed to invalidate refresh token", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to logout")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Successfully logged out"})
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	user, err := h.services.Users().GetUserByID(r.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get user profile", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to get user profile")
		return
	}

	h.respondJSON(w, http.StatusOK, user)
}

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	var updateRequest user.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updatedUser, err := h.services.Users().UpdateProfile(r.Context(), userID, &updateRequest)
	if err != nil {
		h.logger.Error("Failed to update user profile", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to update user profile")
		return
	}

	h.respondJSON(w, http.StatusOK, updatedUser)
}

func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	var changePasswordRequest user.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&changePasswordRequest); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.services.Users().ChangePassword(r.Context(), userID, &changePasswordRequest); err != nil {
		h.logger.Error("Failed to change password", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to change password")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Password successfully changed"})
}

func (h *Handler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, map[string]string{"error": message})
}
