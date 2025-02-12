// internal/config/config.go
package config

import (
    "os"
    "github.com/joho/godotenv"
    "github.com/sirupsen/logrus"
    "fmt"
)

type Config struct {
    Server     ServerConfig
    Database   DatabaseConfig
    Log        LogConfig
    JWTSecret  string    // Новое поле
}

type ServerConfig struct {
    Host string
    Port string
}

type DatabaseConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
    SSLMode  string
    DSN      string    // Add this field
}

type LogConfig struct {
    Level string
    File  string
}

func LoadConfig() (*Config, error) {
    if err := godotenv.Load(); err != nil {
        logrus.Warn("No .env file found")
    }

    dbConfig := DatabaseConfig{
        Host:     getEnv("DB_HOST", "print.czwiyugwum02.eu-north-1.rds.amazonaws.com"),
        Port:     getEnv("DB_PORT", "3306"),  // Make sure we use MySQL port
        User:     getEnv("DB_USER", "root"),  // Default RDS admin username
        Password: getEnv("DB_PASSWORD", "print0101"),   // You'll need to set this
        DBName:   getEnv("DB_NAME", "root"),
        SSLMode:  getEnv("DB_SSL_MODE", "true"),
    }

    // Construct MySQL DSN string with proper SSL settings
    dbConfig.DSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&tls=true&timeout=5s",
        dbConfig.User,
        dbConfig.Password,
        dbConfig.Host,
        dbConfig.Port,
        dbConfig.DBName,
    )

    return &Config{
        Server: ServerConfig{
            Host: getEnv("SERVER_HOST", "localhost"),
            Port: getEnv("SERVER_PORT", "8080"),
        },
        Database: dbConfig,
        Log: LogConfig{
            Level: getEnv("LOG_LEVEL", "info"),
            File:  getEnv("LOG_FILE", "app.log"),
        },
        JWTSecret: getEnv("JWT_SECRET", "your-secret-key"),
    }, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}