// internal/repository/print_jobs.go
package repository

import (
    "database/sql"
    "github.com/google/uuid"
    "print-automation/internal/models"
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

// Добавляем метод для проверки существования задания
func (r *PrintJobRepository) Exists(id string) (bool, error) {
    var exists bool
    query := "SELECT EXISTS(SELECT 1 FROM print_jobs WHERE id = ?)"
    err := r.db.QueryRow(query, id).Scan(&exists)
    return exists, err
}

// Добавляем метод для получения задания по ID
func (r *PrintJobRepository) GetByID(id string) (*models.PrintJob, error) {
    job := &models.PrintJob{}
    query := `
        SELECT id, user_id, file_name, file_url, status, copies, created_at, updated_at
        FROM print_jobs
        WHERE id = ?
    `
    err := r.db.QueryRow(query, id).Scan(
        &job.ID,
        &job.UserID,
        &job.FileName,
        &job.FileURL,
        &job.Status,
        &job.Copies,
        &job.CreatedAt,
        &job.UpdatedAt,
    )
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return job, nil
}