package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
    "print-automation/models"
	"github.com/joho/godotenv"
)

// Глобальная переменная для БД
var DB *gorm.DB

// InitDB инициализирует подключение к БД
func InitDB() {
	// Загружаем переменные окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Не удалось загрузить .env файл, используем системные переменные")
	}

	// Получаем переменные из окружения
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Формируем DSN строку
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbName)

	// Настройки логирования GORM
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Подключаемся к базе данных
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных: ", err)
	}

	fmt.Println("✅ Успешное подключение к базе данных!")

	// Автоматическая миграция
	migrate()
}

// migrate выполняет миграции таблиц
func migrate() {
    err := DB.AutoMigrate(
        &models.Printer{},
        &models.User{},
        &models.PrintJob{},
        &models.Payment{},
    )
	if err != nil {
		log.Fatal("Ошибка миграции: ", err)
	}
}
