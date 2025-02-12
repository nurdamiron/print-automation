// internal/repository/payment.go
package repository

import (
    "database/sql"
    "print-automation/internal/models"
)

type PaymentRepository struct {
    db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
    return &PaymentRepository{db: db}
}

func (r *PaymentRepository) GetByID(id string) (*models.Payment, error) {
    payment := &models.Payment{}
    query := `
        SELECT id, print_job_id, amount, status, payment_method, transaction_id, created_at, updated_at
        FROM payments
        WHERE id = ?
    `
    err := r.db.QueryRow(query, id).Scan(
        &payment.ID,
        &payment.PrintJobID,
        &payment.Amount,
        &payment.Status,
        &payment.PaymentMethod,
        &payment.TransactionID,
        &payment.CreatedAt,
        &payment.UpdatedAt,
    )
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return payment, nil
}

func (r *PaymentRepository) Update(payment *models.Payment) error {
    query := `
        UPDATE payments
        SET status = ?, payment_method = ?, transaction_id = ?, updated_at = ?
        WHERE id = ?
    `
    _, err := r.db.Exec(query,
        payment.Status,
        payment.PaymentMethod,
        payment.TransactionID,
        payment.UpdatedAt,
        payment.ID,
    )
    return err
}

func (r *PaymentRepository) Create(payment *models.Payment) error {
    query := `
        INSERT INTO payments (id, print_job_id, amount, status, payment_method, transaction_id, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `
    _, err := r.db.Exec(query,
        payment.ID,
        payment.PrintJobID,
        payment.Amount,
        payment.Status,
        payment.PaymentMethod,
        payment.TransactionID,
        payment.CreatedAt,
        payment.UpdatedAt,
    )
    return err
}

func (r *PaymentRepository) GetByPrintJobID(printJobID string) (*models.Payment, error) {
    payment := &models.Payment{}
    query := `
        SELECT id, print_job_id, amount, status, payment_method, transaction_id, created_at, updated_at
        FROM payments
        WHERE print_job_id = ?
    `
    err := r.db.QueryRow(query, printJobID).Scan(
        &payment.ID,
        &payment.PrintJobID,
        &payment.Amount,
        &payment.Status,
        &payment.PaymentMethod,
        &payment.TransactionID,
        &payment.CreatedAt,
        &payment.UpdatedAt,
    )
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return payment, nil
}