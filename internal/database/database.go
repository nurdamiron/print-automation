// internal/database/database.go
package database

import (
    "database/sql"
    "fmt"
    "github.com/sirupsen/logrus"
    "print-automation/internal/config"
    "time"
)

type Database struct {
    *sql.DB
}

func NewDatabase(cfg *config.DatabaseConfig) (*Database, error) {
    logrus.Info("Initializing database connection")
    
    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
    )

    db, err := sql.Open("postgres", dsn)
    if err != nil {
        logrus.WithError(err).Error("Failed to open database connection")
        return nil, fmt.Errorf("error opening database: %w", err)
    }

    // Configure connection pool
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(5 * time.Minute)

    // Test connection
    if err := db.Ping(); err != nil {
        logrus.WithError(err).Error("Failed to ping database")
        return nil, fmt.Errorf("error connecting to the database: %w", err)
    }

    logrus.Info("Database connection established successfully")
    return &Database{db}, nil
}

func (d *Database) Close() error {
    logrus.Info("Closing database connection")
    return d.DB.Close()
}

// HealthCheck verifies database connectivity
func (d *Database) HealthCheck() error {
    start := time.Now()
    err := d.Ping()
    duration := time.Since(start)

    if err != nil {
        logrus.WithError(err).Error("Database health check failed")
        return err
    }

    logrus.WithField("duration_ms", duration.Milliseconds()).Info("Database health check successful")
    return nil
}