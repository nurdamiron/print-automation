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

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
    user := &models.User{}
    query := `
        SELECT id, email, password_hash, created_at, updated_at
        FROM users
        WHERE email = ?
    `
    err := r.db.QueryRow(query, email).Scan(
        &user.ID,
        &user.Email,
        &user.PasswordHash,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return user, nil
}

func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
    var exists bool
    query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)"
    err := r.db.QueryRow(query, email).Scan(&exists)
    return exists, err
}