// internal/repository/print_job.go
package repository

import (
    "database/sql"
    "print-automation/internal/models"
)

type PrintJobRepository struct {
    db *sql.DB
}

func NewPrintJobRepository(db *sql.DB) *PrintJobRepository {
    return &PrintJobRepository{db: db}
}

func (r *PrintJobRepository) GetByUserID(userID string) ([]models.PrintJob, error) {
    query := `
        SELECT id, user_id, printer_id, file_name, file_url, status, copies, pages, cost, created_at, updated_at
        FROM print_jobs
        WHERE user_id = ?
    `
    rows, err := r.db.Query(query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var jobs []models.PrintJob
    for rows.Next() {
        var job models.PrintJob
        err := rows.Scan(
            &job.ID,
            &job.UserID,
            &job.PrinterID,
            &job.FileName,
            &job.FileURL,
            &job.Status,
            &job.Copies,
            &job.Pages,
            &job.Cost,
            &job.CreatedAt,
            &job.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        jobs = append(jobs, job)
    }
    return jobs, nil
}

func (r *PrintJobRepository) Update(job *models.PrintJob) error {
    query := `
        UPDATE print_jobs
        SET status = ?, copies = ?, pages = ?, cost = ?, updated_at = ?
        WHERE id = ?
    `
    _, err := r.db.Exec(query,
        job.Status,
        job.Copies,
        job.Pages,
        job.Cost,
        job.UpdatedAt,
        job.ID,
    )
    return err
}

func (r *PrintJobRepository) UpdateStatus(id string, status string) error {
    query := `
        UPDATE print_jobs
        SET status = ?, updated_at = CURRENT_TIMESTAMP
        WHERE id = ?
    `
    _, err := r.db.Exec(query, status, id)
    return err
}

func (r *PrintJobRepository) Exists(id string) (bool, error) {
    var exists bool
    query := "SELECT EXISTS(SELECT 1 FROM print_jobs WHERE id = ?)"
    err := r.db.QueryRow(query, id).Scan(&exists)
    return exists, err
}

func (r *PrintJobRepository) Create(job *models.PrintJob) error {
    query := `
        INSERT INTO print_jobs (id, user_id, printer_id, file_name, file_url, status, copies, pages, cost, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `
    _, err := r.db.Exec(query,
        job.ID,
        job.UserID,
        job.PrinterID,
        job.FileName,
        job.FileURL,
        job.Status,
        job.Copies,
        job.Pages,
        job.Cost,
        job.CreatedAt,
        job.UpdatedAt,
    )
    return err
}

func (r *PrintJobRepository) GetByID(id string) (*models.PrintJob, error) {
    job := &models.PrintJob{}
    query := `
        SELECT id, user_id, printer_id, file_name, file_url, status, copies, pages, cost, created_at, updated_at
        FROM print_jobs
        WHERE id = ?
    `
    err := r.db.QueryRow(query, id).Scan(
        &job.ID,
        &job.UserID,
        &job.PrinterID,
        &job.FileName,
        &job.FileURL,
        &job.Status,
        &job.Copies,
        &job.Pages,
        &job.Cost,
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
