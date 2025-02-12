// internal/handlers/payments.go
package handlers

import (
    "encoding/json"
    "net/http"
    "print-automation/internal/service"
)

type PaymentHandler struct {
    paymentService *service.PaymentService
}

type PaymentRequest struct {
    PrintJobID string  `json:"print_job_id"`
    Amount     float64 `json:"amount"`
}

func NewPaymentHandler(paymentService *service.PaymentService) *PaymentHandler {
    return &PaymentHandler{
        paymentService: paymentService,
    }
}

func (h *PaymentHandler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req PaymentRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    payment, err := h.paymentService.ProcessPayment(req.Amount, req.PrintJobID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(payment)
}

func (h *PaymentHandler) GetPaymentStatus(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    printJobID := r.URL.Query().Get("print_job_id")
    if printJobID == "" {
        http.Error(w, "print_job_id is required", http.StatusBadRequest)
        return
    }

    payment, err := h.paymentService.GetPaymentStatus(printJobID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if payment == nil {
        http.Error(w, "Payment not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(payment)
}