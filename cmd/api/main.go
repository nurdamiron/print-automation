// cmd/api/main.go
package main

import (
    "database/sql"
    "log"
    "print-automation/internal/config"
    "print-automation/internal/handlers"  // imports all handlers
    "print-automation/internal/repository"
    "print-automation/internal/service"   // imports all services
    "print-automation/internal/server"

    _ "github.com/go-sql-driver/mysql"
)

func main() {
    // Загружаем конфигурацию
    cfg := &config.Config{
        Server: config.ServerConfig{
            Host: "localhost",
            Port: "8080",
        },
        Database: config.DatabaseConfig{
            DSN: "root:print0101@tcp(print.czwiyugwum02.eu-north-1.rds.amazonaws.com:3306)/root",
        },
    }

    // Инициализируем подключение к БД
    db, err := sql.Open("mysql", cfg.Database.DSN)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    // Настраиваем пул соединений
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)

    // Проверяем подключение
    if err := db.Ping(); err != nil {
        log.Fatalf("Cannot connect to database: %v", err)
    }
    log.Println("Connected to database successfully")

    // Инициализируем репозитории
    userRepo := repository.NewUserRepository(db)
    printJobRepo := repository.NewPrintJobRepository(db)
    paymentRepo := repository.NewPaymentRepository(db)

    // Инициализируем сервисы
    printerService := service.NewPrinterService()
    paymentService := service.NewPaymentService(paymentRepo, printJobRepo, true)

    // Инициализируем обработчики
    handlers := &server.Handlers{
        UserHandler:     handlers.NewUserHandler(userRepo),
        PrintJobHandler: handlers.NewPrintJobHandler(printJobRepo),
        PaymentHandler:  handlers.NewPaymentHandler(paymentService),
        PrinterHandler:  handlers.NewPrinterHandler(printerService),
    }

    // Создаем и настраиваем сервер
    srv := server.NewServer(cfg, handlers)

    // Запускаем сервер
    log.Printf("Starting server on %s:%s", cfg.Server.Host, cfg.Server.Port)
    if err := srv.Run(); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}