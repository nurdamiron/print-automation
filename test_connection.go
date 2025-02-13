package main

import (
    "fmt"
    "log"
    "net"
    "os"
    "strconv"
    "time"

    "github.com/joho/godotenv"
)

func main() {
    // Загружаем .env
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatalf("Ошибка загрузки .env: %v", err)
    }

    // Читаем из окружения IP и порт
    printerIP := os.Getenv("PRINTER_IP")
    printerPortStr := os.Getenv("PRINTER_PORT")

    if printerIP == "" {
        log.Fatal("Не задана переменная окружения PRINTER_IP")
    }
    if printerPortStr == "" {
        log.Fatal("Не задана переменная окружения PRINTER_PORT")
    }

    printerPort, err := strconv.Atoi(printerPortStr)
    if err != nil {
        log.Fatalf("Некорректное значение PRINTER_PORT: %v", err)
    }

    // Формируем адрес вида "192.168.1.50:9100"
    addr := fmt.Sprintf("%s:%d", printerIP, printerPort)

    // Пытаемся подключиться к принтеру
    conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
    if err != nil {
        fmt.Println("Принтер недоступен:", err)
        return
    }
    defer conn.Close()

    fmt.Println("Успешное соединение с принтером на", addr)
}
