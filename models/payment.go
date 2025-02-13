package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

type Payment struct {
    ID            string    `gorm:"type:varchar(36);primaryKey"`
    PrintJobID    string    `gorm:"type:varchar(36);not null"`
    Amount        float64   `gorm:"type:decimal(8,2);not null"`
    Status        string    `gorm:"type:varchar(50);not null;default:'created'"`
    PaymentMethod string    `gorm:"type:varchar(50)"`
    TransactionID string    `gorm:"type:varchar(100)"`
    CreatedAt     time.Time `gorm:"not null"`
    UpdatedAt     time.Time `gorm:"not null"`
}

func (pm *Payment) BeforeCreate(tx *gorm.DB) (err error) {
    pm.ID = uuid.New().String()
    pm.CreatedAt = time.Now()
    pm.UpdatedAt = time.Now()
    return
}

func (pm *Payment) BeforeUpdate(tx *gorm.DB) (err error) {
    pm.UpdatedAt = time.Now()
    return
}
