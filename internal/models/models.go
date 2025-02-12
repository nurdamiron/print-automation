// internal/models/models.go
package models

import (
    "time"
)


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



