// internal/models/payment.go
package models

import "time"

type Payment struct {
    ID            string     `json:"id"`
    PrintJobID    string     `json:"print_job_id"`
    Amount        float64    `json:"amount"`
    Status        string     `json:"status"`
    PaymentMethod string     `json:"payment_method,omitempty"`
    TransactionID string     `json:"transaction_id,omitempty"`
    CreatedAt     time.Time  `json:"created_at"`
    UpdatedAt     time.Time  `json:"updated_at"`
}
