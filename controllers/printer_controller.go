package controllers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "print-automation/config"
    "print-automation/models"
	"print-automation/services"

)

// Получить список принтеров
func GetAllPrinters(c *gin.Context) {
    var printers []models.Printer
    if err := config.DB.Find(&printers).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, printers)
}

// Создать (зарегистрировать) новый принтер
func CreatePrinter(c *gin.Context) {
    var printer models.Printer
    if err := c.ShouldBindJSON(&printer); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := config.DB.Create(&printer).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, printer)
}

func CheckPrinterConnectionHandler(c *gin.Context) {
    printerID := c.Param("id")

    // 1. Читаем из БД
    var printer models.Printer
    if err := config.DB.First(&printer, "id = ?", printerID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Принтер не найден"})
        return
    }

    // 2. Пытаемся открыть TCP-соединение 
    err := services.CheckPrinterConnection(printer.IPAddress, printer.Port)
    if err != nil {
        c.JSON(http.StatusOK, gin.H{"status": "offline", "message": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"status": "online", "message": "Успешное соединение"})
}

// Получить конкретный принтер
func GetPrinterByID(c *gin.Context) {
    id := c.Param("id")
    var printer models.Printer
    if err := config.DB.First(&printer, "id = ?", id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Принтер не найден"})
        return
    }
    c.JSON(http.StatusOK, printer)
}

// Обновить информацию о принтере
func UpdatePrinter(c *gin.Context) {
    id := c.Param("id")
    var printer models.Printer
    if err := config.DB.First(&printer, "id = ?", id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Принтер не найден"})
        return
    }

    var input models.Printer
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    printer.Name = input.Name
    printer.IPAddress = input.IPAddress
    printer.Port = input.Port
    printer.Protocol = input.Protocol
    printer.IsOnline = input.IsOnline
    printer.Status = input.Status

    if err := config.DB.Save(&printer).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, printer)
}

// Удалить принтер
func DeletePrinter(c *gin.Context) {
    id := c.Param("id")
    if err := config.DB.Delete(&models.Printer{}, "id = ?", id).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Принтер удалён"})
}
