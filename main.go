package main

import (
    "fmt"
    "print-automation/config"
    "print-automation/routers"
)

func main() {
    // Инициализация БД
    config.InitDB()

    // Настройка роутера
    r := routers.SetupRouter()

    // Запуск сервера
    port := ":8080"
    fmt.Println("Сервер запущен на порту", port)
    if err := r.Run(port); err != nil {
        panic(err)
    }
}
