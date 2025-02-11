// internal/database/database.go
package database

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    "print-automation/internal/config"
)

type Database struct {
    *sql.DB
}

func NewDatabase(cfg *config.DatabaseConfig) (*Database, error) {
    // Формируем строку подключения
    connStr := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
        cfg.Host,
        cfg.Port,
        cfg.User,
        cfg.Password,
        cfg.DBName,
    )

    // Открываем соединение
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, fmt.Errorf("error opening database: %w", err)
    }

    // Проверяем соединение
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("error connecting to the database: %w", err)
    }

    // Настраиваем пул соединений
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)

    return &Database{db}, nil
}

func (d *Database) Close() error {
    return d.DB.Close()
}

// Пример метода для проверки здоровья БД
func (d *Database) HealthCheck() error {
    return d.Ping()
}