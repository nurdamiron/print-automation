// cmd/api/main.go
package main

import (
    "log"
    "print-automation/internal/config"
    "print-automation/internal/server"
)

func main() {
    // Загрузка конфигурации
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Инициализация и запуск сервера
    srv := server.NewServer(cfg)
    if err := srv.Run(); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}