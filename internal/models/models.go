// internal/models/models.go
package models

import (
    "time"
)

type User struct {
    ID           string    `json:"id"`
    Email        string    `json:"email"`
    PasswordHash string    `json:"-"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

type PrintJob struct {
    ID        string    `json:"id"`
    UserID    string    `json:"user_id"`
    FileName  string    `json:"file_name"`
    FileURL   string    `json:"file_url"`
    Status    string    `json:"status"`
    Copies    int       `json:"copies"`
    Pages     *int      `json:"pages,omitempty"`
    Cost      *float64  `json:"cost,omitempty"`
    PrinterID *string   `json:"printer_id,omitempty"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// internal/repository/print_jobs.go
package repository

import (
    "database/sql"
    "github.com/google/uuid"
    "your-project/internal/models"
)

type PrintJobRepository struct {
    db *sql.DB
}

func NewPrintJobRepository(db *sql.DB) *PrintJobRepository {
    return &PrintJobRepository{db: db}
}

func (r *PrintJobRepository) Create(job *models.PrintJob) error {
    job.ID = uuid.New().String()
    query := `
        INSERT INTO print_jobs (id, user_id, file_name, file_url, status, copies)
        VALUES (?, ?, ?, ?, ?, ?)
    `
    _, err := r.db.Exec(query, job.ID, job.UserID, job.FileName, job.FileURL, 
        "pending", job.Copies)
    return err
}

// internal/handlers/print_jobs.go
package handlers

import (
    "encoding/json"
    "net/http"
    "your-project/internal/models"
    "your-project/internal/repository"
)

type PrintJobHandler struct {
    repo *repository.PrintJobRepository
}

func NewPrintJobHandler(repo *repository.PrintJobRepository) *PrintJobHandler {
    return &PrintJobHandler{repo: repo}
}

func (h *PrintJobHandler) Create(w http.ResponseWriter, r *http.Request) {
    var job models.PrintJob
    if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.repo.Create(&job); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(job)
}

// internal/service/payment.go
package service

import (
    "errors"
    "your-project/internal/models"
)

type PaymentService struct {
    // В реальном проекте здесь будет интеграция с платежной системой
    demoMode bool
}

func NewPaymentService(demoMode bool) *PaymentService {
    return &PaymentService{demoMode: demoMode}
}

func (s *PaymentService) ProcessPayment(amount float64) (*models.Payment, error) {
    if !s.demoMode {
        return nil, errors.New("real payment processing not implemented")
    }
    
    // Демо-режим всегда возвращает успешный платеж
    return &models.Payment{
        Status: "completed",
        Amount: amount,
    }, nil
}

// cmd/api/main.go
package main

import (
    "database/sql"
    "log"
    "net/http"
    _ "github.com/go-sql-driver/mysql"
    "print-automation/internal/handlers"
    "print-automation/internal/repository"
    "print-automation/internal/service"
)

func main() {
    db, err := sql.Open("mysql", "user:pass@tcp(print.czwiyugwum02.eu-north-1.rds.amazonaws.com:3306)/print_service")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Инициализация репозиториев
    printJobRepo := repository.NewPrintJobRepository(db)
    
    // Инициализация сервисов
    paymentService := service.NewPaymentService(true) // demo mode
    
    // Инициализация обработчиков
    printJobHandler := handlers.NewPrintJobHandler(printJobRepo)
    
    // Маршрутизация
    http.HandleFunc("/api/print-jobs", printJobHandler.Create)
    
    log.Fatal(http.ListenAndServe(":8080", nil))
}