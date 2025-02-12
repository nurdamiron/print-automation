// internal/service/payment.go
package service

import (
    "errors"
    "fmt"
    "github.com/google/uuid"
    "print-automation/internal/models"
    "print-automation/internal/repository"
    "time"
)

type PaymentService struct {
    paymentRepo *repository.PaymentRepository
    printJobRepo *repository.PrintJobRepository
    demoMode bool
}

func NewPaymentService(paymentRepo *repository.PaymentRepository, printJobRepo *repository.PrintJobRepository, demoMode bool) *PaymentService {
    return &PaymentService{
        paymentRepo: paymentRepo,
        printJobRepo: printJobRepo,
        demoMode: demoMode,
    }
}

func (s *PaymentService) ProcessPayment(amount float64, printJobID string) (*models.Payment, error) {
    // Проверяем существование задания печати
    exists, err := s.printJobRepo.Exists(printJobID)
    if err != nil {
        return nil, fmt.Errorf("error checking print job: %w", err)
    }
    if !exists {
        return nil, errors.New("print job not found")
    }

    if !s.demoMode {
        return nil, errors.New("real payment processing not implemented")
    }
    
    payment := &models.Payment{
        ID:            uuid.New().String(),
        PrintJobID:    printJobID,
        Amount:        amount,
        Status:        "completed",
        PaymentMethod: "demo",
        TransactionID: "demo_" + uuid.New().String(),
        CreatedAt:     time.Now(),
        UpdatedAt:     time.Now(),
    }

    if err := s.paymentRepo.Create(payment); err != nil {
        return nil, fmt.Errorf("error creating payment: %w", err)
    }
    
    return payment, nil
}

func (s *PaymentService) GetPaymentStatus(printJobID string) (*models.Payment, error) {
    return s.paymentRepo.GetByPrintJobID(printJobID)
}