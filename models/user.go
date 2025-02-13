package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

type User struct {
    ID           string    `gorm:"type:varchar(36);primaryKey"`
    Email        string    `gorm:"type:varchar(255);unique;not null"`
    PasswordHash string    `gorm:"type:varchar(255);not null"`
    CreatedAt    time.Time `gorm:"not null"`
    UpdatedAt    time.Time `gorm:"not null"`
}

// Генерация ID и временных меток
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
    u.ID = uuid.New().String()
    u.CreatedAt = time.Now()
    u.UpdatedAt = time.Now()
    return
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
    u.UpdatedAt = time.Now()
    return
}
