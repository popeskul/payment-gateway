package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/popeskul/payment-gateway/internal/core/domain/refund"
	"github.com/popeskul/payment-gateway/internal/infrastructure/metrics"
)

func (h *Handler) CreateRefund(w http.ResponseWriter, r *http.Request) {
	var ref refund.Refund
	if err := json.NewDecoder(r.Body).Decode(&ref); err != nil {
		h.logger.Error("Failed to decode refund", "error", err)
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.services.Refunds().CreateRefund(r.Context(), &ref); err != nil {
		h.logger.Error("Failed to create refund", "error", err)
		http.Error(w, "Failed to create refund", http.StatusInternalServerError)
		metrics.RefundTotal.WithLabelValues("failed").Inc()
		return
	}

	metrics.RefundTotal.WithLabelValues("success").Inc()
	metrics.RefundAmount.WithLabelValues(ref.Currency).Observe(ref.Amount)

	respondJSON(w, http.StatusCreated, ref)
}

func (h *Handler) GetRefund(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ref, err := h.services.Refunds().GetRefund(r.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get refund", "error", err, "id", id)
		http.Error(w, "Refund not found", http.StatusNotFound)
		return
	}

	respondJSON(w, http.StatusOK, ref)
}

func (h *Handler) UpdateRefund(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var ref refund.Refund
	if err := json.NewDecoder(r.Body).Decode(&ref); err != nil {
		h.logger.Error("Failed to decode refund", "error", err)
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	ref.ID = id
	if err := h.services.Refunds().UpdateRefund(r.Context(), &ref); err != nil {
		h.logger.Error("Failed to update refund", "error", err, "id", id)
		http.Error(w, "Failed to update refund", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, ref)
}

func (h *Handler) ListRefunds(w http.ResponseWriter, r *http.Request) {
	paymentID := r.URL.Query().Get("payment_id")
	if paymentID == "" {
		h.logger.Error("Payment ID is required")
		http.Error(w, "Payment ID is required", http.StatusBadRequest)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit == 0 {
		limit = 10
	}

	refunds, err := h.services.Refunds().ListRefunds(r.Context(), paymentID, limit, offset)
	if err != nil {
		h.logger.Error("Failed to list refunds", "error", err)
		http.Error(w, "Failed to list refunds", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, refunds)
}

func (h *Handler) ProcessRefund(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.services.Refunds().ProcessRefund(r.Context(), id); err != nil {
		h.logger.Error("Failed to process refund", "error", err, "id", id)
		http.Error(w, "Failed to process refund", http.StatusInternalServerError)
		metrics.RefundTotal.WithLabelValues("failed").Inc()
		return
	}

	metrics.RefundTotal.WithLabelValues("processed").Inc()

	respondJSON(w, http.StatusOK, map[string]string{"message": "Refund processed successfully"})
}
