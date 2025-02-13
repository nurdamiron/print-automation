package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

type PrintJob struct {
    ID        string    `gorm:"type:varchar(36);primaryKey"`
    UserID    string    `gorm:"type:varchar(36);not null"`
    PrinterID string    `gorm:"type:varchar(36);not null"`
    FileURL   string    `gorm:"type:varchar(255)"`
    Status    string    `gorm:"type:varchar(50);not null;default:'pending'"`
    Copies    int       `gorm:"not null;default:1"`
    Pages     int       `gorm:"not null;default:1"`
    Cost      float64   `gorm:"type:decimal(8,2)"`
    CreatedAt time.Time `gorm:"not null"`
    UpdatedAt time.Time `gorm:"not null"`
}

func (pj *PrintJob) BeforeCreate(tx *gorm.DB) (err error) {
    pj.ID = uuid.New().String()
    pj.CreatedAt = time.Now()
    pj.UpdatedAt = time.Now()
    return
}

func (pj *PrintJob) BeforeUpdate(tx *gorm.DB) (err error) {
    pj.UpdatedAt = time.Now()
    return
}
