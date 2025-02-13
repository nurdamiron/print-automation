package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

type Printer struct {
    ID         string    `gorm:"type:varchar(36);primaryKey"`
    Name       string    `gorm:"type:varchar(255);not null"`
    IPAddress  string    `gorm:"type:varchar(64);not null"`
    Port       int       `gorm:"not null"`
    Protocol   string    `gorm:"type:varchar(20);not null"`
    IsOnline   bool      `gorm:"not null;default:false"`
    Status     string    `gorm:"type:varchar(50);not null;default:'UNKNOWN'"`
    CreatedAt  time.Time `gorm:"not null"`
    UpdatedAt  time.Time `gorm:"not null"`
}

// Хук GORM для генерации UUID и временных меток
func (p *Printer) BeforeCreate(tx *gorm.DB) (err error) {
    p.ID = uuid.New().String()
    p.CreatedAt = time.Now()
    p.UpdatedAt = time.Now()
    return
}

func (p *Printer) BeforeUpdate(tx *gorm.DB) (err error) {
    p.UpdatedAt = time.Now()
    return
}
