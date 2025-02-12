// internal/service/payment.go
package service

import (
    "errors"
    "fmt"
    "github.com/google/uuid"
    "github.com/sirupsen/logrus"
    "print-automation/internal/models"
    "print-automation/internal/repository"
    "time"
)

type PaymentService struct {
    paymentRepo *repository.PaymentRepository
    printJobRepo *repository.PrintJobRepository
    demoMode bool
}

func (s *PaymentService) GetPaymentByJobID(printJobID string) (*models.Payment, error) {
    return s.paymentRepo.GetByPrintJobID(printJobID)
}

func NewPaymentService(paymentRepo *repository.PaymentRepository, printJobRepo *repository.PrintJobRepository, demoMode bool) *PaymentService {
    return &PaymentService{
        paymentRepo: paymentRepo,
        printJobRepo: printJobRepo,
        demoMode: demoMode,
    }
}

func (s *PaymentService) ProcessCallback(paymentID, status, transactionID string) error {
    logger := logrus.WithFields(logrus.Fields{
        "payment_id": paymentID,
        "status": status,
        "transaction_id": transactionID,
    })
    
    logger.Info("Processing payment callback")

    // Получаем платеж из БД
    payment, err := s.paymentRepo.GetByID(paymentID)
    if err != nil {
        logger.WithError(err).Error("Failed to get payment")
        return fmt.Errorf("failed to get payment: %w", err)
    }

    if payment == nil {
        logger.Error("Payment not found")
        return fmt.Errorf("payment not found")
    }

    // Обновляем статус платежа
    payment.Status = status
    payment.TransactionID = transactionID

    if err := s.paymentRepo.Update(payment); err != nil {
        logger.WithError(err).Error("Failed to update payment")
        return fmt.Errorf("failed to update payment: %w", err)
    }

    // Если платеж успешен, обновляем статус задания печати
    if status == "completed" {
        if err := s.printJobRepo.UpdateStatus(payment.PrintJobID, "ready_to_print"); err != nil {
            logger.WithError(err).Error("Failed to update print job status")
            return fmt.Errorf("failed to update print job status: %w", err)
        }
    }

    logger.Info("Payment callback processed successfully")
    return nil
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