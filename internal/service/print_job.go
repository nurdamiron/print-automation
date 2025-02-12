// internal/service/print_job.go
package service

import (
    "fmt"
    "time"
    "io"
    "github.com/google/uuid"
    "github.com/sirupsen/logrus"
    "print-automation/internal/models"
    "print-automation/internal/repository"
)

// Константы статусов заданий печати
const (
    StatusPending    = "pending"     // Ожидает оплаты
    StatusPaid       = "paid"        // Оплачено
    StatusProcessing = "processing"  // В процессе печати
    StatusCompleted  = "completed"   // Завершено
    StatusFailed     = "failed"      // Ошибка печати
    StatusCancelled  = "cancelled"   // Отменено
)

// PrintJobService управляет заданиями печати
type PrintJobService struct {
    jobRepo         *repository.PrintJobRepository
    printerService  *PrinterService
    paymentService  *PaymentService
    logger          *logrus.Logger
}

// NewPrintJobService создает новый экземпляр сервиса
func NewPrintJobService(
    jobRepo *repository.PrintJobRepository,
    printerService *PrinterService,
    paymentService *PaymentService,
    logger *logrus.Logger,
) *PrintJobService {
    return &PrintJobService{
        jobRepo:         jobRepo,
        printerService:  printerService,
        paymentService:  paymentService,
        logger:         logger,
    }
}

// Create создает новое задание печати
func (s *PrintJobService) Create(userID string, printerID string, document io.Reader, options PrintOptions) (*models.PrintJob, error) {
    log := s.logger.WithFields(logrus.Fields{
        "user_id": userID,
        "printer_id": printerID,
    })
    log.Info("Creating new print job")

    // Проверяем доступность принтера
    printer, err := s.printerService.GetPrinterStatus(printerID)
    if err != nil {
        log.WithError(err).Error("Failed to get printer status")
        return nil, fmt.Errorf("printer unavailable: %w", err)
    }

    if !printer.IsOnline {
        log.Error("Printer is offline")
        return nil, fmt.Errorf("printer is offline")
    }

    // Создаем новое задание
    job := &models.PrintJob{
        ID:        uuid.New().String(),
        UserID:    userID,
        PrinterID: printerID,
        Status:    StatusPending,
        Copies:    options.Copies,
        Pages:     &options.Pages,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    // Вычисляем стоимость
    cost, err := s.calculateCost(options)
    if err != nil {
        log.WithError(err).Error("Failed to calculate cost")
        return nil, fmt.Errorf("failed to calculate cost: %w", err)
    }
    job.Cost = &cost

    // Сохраняем документ
    fileURL, err := s.saveDocument(job.ID, document)
    if err != nil {
        log.WithError(err).Error("Failed to save document")
        return nil, fmt.Errorf("failed to save document: %w", err)
    }
    job.FileURL = fileURL

    // Сохраняем задание в БД
    if err := s.jobRepo.Create(job); err != nil {
        log.WithError(err).Error("Failed to create job in database")
        return nil, fmt.Errorf("failed to create job: %w", err)
    }

    log.WithField("job_id", job.ID).Info("Print job created successfully")
    return job, nil
}

// GetByID получает задание по ID
func (s *PrintJobService) GetByID(id string) (*models.PrintJob, error) {
    log := s.logger.WithField("job_id", id)
    log.Info("Getting print job")

    job, err := s.jobRepo.GetByID(id)
    if err != nil {
        log.WithError(err).Error("Failed to get print job")
        return nil, fmt.Errorf("failed to get job: %w", err)
    }

    if job == nil {
        return nil, fmt.Errorf("job not found")
    }

    return job, nil
}

// GetUserJobs получает все задания пользователя
func (s *PrintJobService) GetUserJobs(userID string) ([]models.PrintJob, error) {
    log := s.logger.WithField("user_id", userID)
    log.Info("Getting user print jobs")

    jobs, err := s.jobRepo.GetByUserID(userID)
    if err != nil {
        log.WithError(err).Error("Failed to get user jobs")
        return nil, fmt.Errorf("failed to get user jobs: %w", err)
    }

    return jobs, nil
}

// UpdateStatus обновляет статус задания
func (s *PrintJobService) UpdateStatus(id string, status string) error {
    log := s.logger.WithFields(logrus.Fields{
        "job_id": id,
        "status": status,
    })
    log.Info("Updating job status")

    job, err := s.jobRepo.GetByID(id)
    if err != nil {
        log.WithError(err).Error("Failed to get job")
        return fmt.Errorf("failed to get job: %w", err)
    }

    if job == nil {
        return fmt.Errorf("job not found")
    }

    prevStatus := job.Status
    job.Status = status
    job.UpdatedAt = time.Now()

    if err := s.jobRepo.Update(job); err != nil {
        log.WithError(err).Error("Failed to update job")
        return fmt.Errorf("failed to update job: %w", err)
    }

    log.WithFields(logrus.Fields{
        "previous_status": prevStatus,
        "new_status": status,
    }).Info("Job status updated successfully")

    return nil
}

// StartPrinting начинает процесс печати
func (s *PrintJobService) StartPrinting(id string) error {
    log := s.logger.WithField("job_id", id)
    log.Info("Starting print job")

    job, err := s.jobRepo.GetByID(id)
    if err != nil {
        log.WithError(err).Error("Failed to get job")
        return fmt.Errorf("failed to get job: %w", err)
    }

    if job == nil {
        return fmt.Errorf("job not found")
    }

    // Проверяем оплату
    payment, err := s.paymentService.GetPaymentByJobID(id)
    if err != nil {
        log.WithError(err).Error("Failed to get payment info")
        return fmt.Errorf("failed to get payment info: %w", err)
    }

    if payment == nil || payment.Status != "completed" {
        return fmt.Errorf("payment required")
    }

    // Проверяем принтер
    printer, err := s.printerService.GetPrinterStatus(job.PrinterID)
    if err != nil {
        log.WithError(err).Error("Failed to get printer status")
        return fmt.Errorf("failed to get printer status: %w", err)
    }

    if !printer.IsOnline {
        return fmt.Errorf("printer is offline")
    }

    // Обновляем статус и начинаем печать
    if err := s.UpdateStatus(id, StatusProcessing); err != nil {
        return err
    }

    // Асинхронная печать
    go func() {
        if err := s.processPrintJob(job); err != nil {
            log.WithError(err).Error("Print job failed")
            s.UpdateStatus(id, StatusFailed)
            return
        }
        s.UpdateStatus(id, StatusCompleted)
    }()

    return nil
}

// CancelJob отменяет задание
func (s *PrintJobService) CancelJob(id string) error {
    log := s.logger.WithField("job_id", id)
    log.Info("Cancelling print job")

    job, err := s.jobRepo.GetByID(id)
    if err != nil {
        log.WithError(err).Error("Failed to get job")
        return fmt.Errorf("failed to get job: %w", err)
    }

    if job == nil {
        return fmt.Errorf("job not found")
    }

    if job.Status == StatusCompleted {
        return fmt.Errorf("cannot cancel completed job")
    }

    if job.Status == StatusProcessing {
        // Пытаемся остановить печать
        if err := s.printerService.CancelPrint(job.PrinterID); err != nil {
            log.WithError(err).Error("Failed to cancel printing")
            return fmt.Errorf("failed to cancel printing: %w", err)
        }
    }

    return s.UpdateStatus(id, StatusCancelled)
}

// Вспомогательные структуры и методы

type PrintOptions struct {
    Copies    int
    Pages     int
    DoubleSided bool
    Color      bool
}

func (s *PrintJobService) calculateCost(options PrintOptions) (float64, error) {
    // Базовая стоимость страницы
    var baseCost float64 = 2.0
    if options.Color {
        baseCost = 5.0
    }

    // Скидка для двусторонней печати
    if options.DoubleSided {
        baseCost *= 0.8
    }

    totalPages := options.Pages * options.Copies
    totalCost := float64(totalPages) * baseCost

    return totalCost, nil
}

func (s *PrintJobService) processPrintJob(job *models.PrintJob) error {
    log := s.logger.WithFields(logrus.Fields{
        "job_id": job.ID,
        "printer_id": job.PrinterID,
    })
    log.Info("Processing print job")

    // Получаем документ
    document, err := s.getDocument(job.FileURL)
    if err != nil {
        return fmt.Errorf("failed to get document: %w", err)
    }
    defer document.Close()

    // Отправляем на печать
    if err := s.printerService.PrintDocument(job.PrinterID, document); err != nil {
        return fmt.Errorf("failed to print document: %w", err)
    }

    log.Info("Print job processed successfully")
    return nil
}

func (s *PrintJobService) saveDocument(jobID string, document io.Reader) (string, error) {
    // Здесь должна быть реализация сохранения документа
    // Например, в файловую систему или облачное хранилище
    fileURL := fmt.Sprintf("/documents/%s", jobID)
    return fileURL, nil
}

func (s *PrintJobService) getDocument(fileURL string) (io.ReadCloser, error) {
    // Здесь должна быть реализация получения документа
    // из файловой системы или облачного хранилища
    return nil, fmt.Errorf("not implemented")
}