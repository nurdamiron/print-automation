package main

import (
    "database/sql"
    "log"
    "fmt"
    "print-automation/internal/api"
    "print-automation/internal/config"
    "print-automation/internal/repository"
    "print-automation/internal/service"
    _ "github.com/go-sql-driver/mysql"  // Make sure this import exists
)

func main() {
    // Загружаем конфигурацию
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Error loading config: %v", err)
    }

    // Формируем DSN для MySQL
    dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true",
        cfg.Database.User,
        cfg.Database.Password,
        cfg.Database.Host,
        cfg.Database.DBName,
    )

    // Используем драйвер MySQL
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatalf("Failed to create database connection: %v", err)
    }
    defer db.Close()

    // Устанавливаем таймаут подключения
    db.SetConnMaxLifetime(0)
    db.SetMaxIdleConns(50)
    db.SetMaxOpenConns(50)

    // Проверяем соединение
    if err := db.Ping(); err != nil {
        log.Fatalf("Cannot connect to database: %v", err)
    }
    log.Println("✅ Connected to database successfully")

    // Остальной код остается тем же
    userRepo := repository.NewUserRepository(db)
    printJobRepo := repository.NewPrintJobRepository(db)
    paymentRepo := repository.NewPaymentRepository(db)

    printerService := service.NewPrinterService()
    printJobService := service.NewPrintJobService(printJobRepo, printerService, nil, nil)
    paymentService := service.NewPaymentService(paymentRepo, printJobRepo, true)
    authService := service.NewAuthService(userRepo, cfg.JWTSecret)

    server := api.NewServer(
        cfg,
        printerService,
        paymentService,
        authService,
        printJobService,
    )

    log.Printf("🚀 Starting server on %s:%s", cfg.Server.Host, cfg.Server.Port)

    if err := server.Run(); err != nil {
        log.Fatalf("❌ Server failed to start: %v", err)
    }
}