// internal/database/db.go
package database

import (
    "database/sql"
    "fmt"
    "os"
    _ "github.com/lib/pq"
)

type Database struct {
    DB *sql.DB
}

// NewDatabase создает новое подключение к базе данных
func NewDatabase() (*Database, error) {
    // Получаем строку подключения из переменной окружения
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        return nil, fmt.Errorf("DATABASE_URL is not set")
    }

    // Открываем соединение
    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        return nil, fmt.Errorf("could not connect to database: %w", err)
    }

    // Проверяем соединение
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("could not ping database: %w", err)
    }

    // Настраиваем пул соединений
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)

    return &Database{DB: db}, nil
}

// Close закрывает соединение с базой данных
func (d *Database) Close() error {
    return d.DB.Close()
}

// TestConnection проверяет подключение к базе данных
func (d *Database) TestConnection() error {
    return d.DB.Ping()
}