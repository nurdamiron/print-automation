// cmd/api/main.go
package main

import (
    "database/sql"
    "log"
    "net/http"
    _ "github.com/go-sql-driver/mysql"
    "print-automation/internal/handlers"
    "print-automation/internal/repository"
    "print-automation/internal/service"
)

func main() {
    // Подключение к БД с правильными учетными данными
    db, err := sql.Open("mysql", "root:print0101@tcp(print.czwiyugwum02.eu-north-1.rds.amazonaws.com:3306)/root")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Проверка подключения к БД
    if err := db.Ping(); err != nil {
        log.Fatal("Cannot connect to database:", err)
    }
    log.Println("Connected to database successfully")

    // Инициализация репозиториев
    userRepo := repository.NewUserRepository(db)
    printJobRepo := repository.NewPrintJobRepository(db)
    paymentRepo := repository.NewPaymentRepository(db)
    
    // Инициализация сервисов
    paymentService := service.NewPaymentService(paymentRepo, printJobRepo, true)
    
    // Инициализация обработчиков
    userHandler := handlers.NewUserHandler(userRepo)
    printJobHandler := handlers.NewPrintJobHandler(printJobRepo)
    paymentHandler := handlers.NewPaymentHandler(paymentService)
    
    // Маршрутизация
    mux := http.NewServeMux()
    mux.HandleFunc("/api/users", userHandler.Create)
    mux.HandleFunc("/api/print-jobs", printJobHandler.Create)
    mux.HandleFunc("/api/payments", paymentHandler.ProcessPayment)
    mux.HandleFunc("/api/payments/status", paymentHandler.GetPaymentStatus)
    
    // Запуск сервера
    log.Println("Starting server on :8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}