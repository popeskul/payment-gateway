package handlers

import (
	"encoding/json"
	"github.com/popeskul/payment-gateway/internal/infrastructure/metrics"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/popeskul/payment-gateway/internal/core/domain/payment"
)

func (h *Handler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var p payment.Payment
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		h.logger.Error("Failed to decode payment", "error", err)
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	merchantID := r.Context().Value("merchantID").(string)
	p.MerchantID = merchantID

	if err := h.services.Payments().CreatePayment(r.Context(), &p); err != nil {
		h.logger.Error("Failed to create payment", "error", err)
		http.Error(w, "Failed to create payment", http.StatusInternalServerError)
		metrics.PaymentTotal.WithLabelValues("failed").Inc()
		return
	}

	metrics.PaymentTotal.WithLabelValues("success").Inc()
	metrics.PaymentAmount.WithLabelValues(p.Currency).Observe(p.Amount)

	respondJSON(w, http.StatusCreated, p)
}

func (h *Handler) GetPayment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	p, err := h.services.Payments().GetPayment(r.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get payment", "error", err, "id", id)
		http.Error(w, "Payment not found", http.StatusNotFound)
		return
	}

	respondJSON(w, http.StatusOK, p)
}

func (h *Handler) UpdatePayment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var p payment.Payment
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		h.logger.Error("Failed to decode payment", "error", err)
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	p.ID = id
	if err := h.services.Payments().UpdatePayment(r.Context(), &p); err != nil {
		h.logger.Error("Failed to update payment", "error", err, "id", id)
		http.Error(w, "Failed to update payment", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, p)
}

func (h *Handler) ListPayments(w http.ResponseWriter, r *http.Request) {
	merchantID := r.Context().Value("merchantID").(string)
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit == 0 {
		limit = 10
	}

	payments, err := h.services.Payments().ListPayments(r.Context(), merchantID, limit, offset)
	if err != nil {
		h.logger.Error("Failed to list payments", "error", err)
		http.Error(w, "Failed to list payments", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, payments)
}

func (h *Handler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.services.Payments().ProcessPayment(r.Context(), id); err != nil {
		h.logger.Error("Failed to process payment", "error", err, "id", id)
		http.Error(w, "Failed to process payment", http.StatusInternalServerError)
		metrics.PaymentTotal.WithLabelValues("failed").Inc()
		return
	}

	metrics.PaymentTotal.WithLabelValues("processed").Inc()

	respondJSON(w, http.StatusOK, map[string]string{"message": "Payment processed successfully"})
}
