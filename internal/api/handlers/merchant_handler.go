package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/popeskul/payment-gateway/internal/core/domain/merchant"
)

func (h *Handler) CreateMerchant(w http.ResponseWriter, r *http.Request) {
	var m merchant.Merchant
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		h.logger.Error("Failed to decode merchant", "error", err)
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.services.Merchants().CreateMerchant(r.Context(), &m); err != nil {
		h.logger.Error("Failed to create merchant", "error", err)
		http.Error(w, "Failed to create merchant", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusCreated, m)
}

func (h *Handler) GetMerchant(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	m, err := h.services.Merchants().GetMerchant(r.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get merchant", "error", err, "id", id)
		http.Error(w, "Merchant not found", http.StatusNotFound)
		return
	}

	respondJSON(w, http.StatusOK, m)
}

func (h *Handler) UpdateMerchant(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var m merchant.Merchant
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		h.logger.Error("Failed to decode merchant", "error", err)
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	m.ID = id
	if err := h.services.Merchants().UpdateMerchant(r.Context(), &m); err != nil {
		h.logger.Error("Failed to update merchant", "error", err, "id", id)
		http.Error(w, "Failed to update merchant", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, m)
}

func (h *Handler) DeleteMerchant(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.services.Merchants().DeleteMerchant(r.Context(), id); err != nil {
		h.logger.Error("Failed to delete merchant", "error", err, "id", id)
		http.Error(w, "Failed to delete merchant", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Merchant deleted successfully"})
}

func (h *Handler) ListMerchants(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit == 0 {
		limit = 10
	}

	merchants, err := h.services.Merchants().ListMerchants(r.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Failed to list merchants", "error", err)
		http.Error(w, "Failed to list merchants", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, merchants)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
