// internal/repository/users.go
package repository

import (
    "database/sql"
    "github.com/google/uuid"
    "print-automation/internal/models"
    "time"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(email, passwordHash string) (*models.User, error) {
    user := &models.User{
        ID:           uuid.New().String(),
        Email:        email,
        PasswordHash: passwordHash,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }

    query := `
        INSERT INTO users (id, email, password_hash, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?)
    `
    _, err := r.db.Exec(query, 
        user.ID, 
        user.Email, 
        user.PasswordHash,
        user.CreatedAt,
        user.UpdatedAt,
    )
    if err != nil {
        return nil, err
    }

    return user, nil
}

