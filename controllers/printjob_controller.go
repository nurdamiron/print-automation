package controllers

import (
    "net/http"
	"fmt"
    "os"
	"io"
    "path/filepath"
    "github.com/gin-gonic/gin"
    "print-automation/config"
    "print-automation/models"
	"print-automation/services"

)

// DownloadAndSendToPrinter скачивает файл по URL и отправляет на принтер
func DownloadAndSendToPrinter(jobID, fileURL, printerIP string, printerPort int) error {
    // Скачиваем файл во временную директорию
    tempFilePath := fmt.Sprintf("/tmp/%s.pdf", jobID)
    err := downloadFile(tempFilePath, fileURL)
    if err != nil {
        return fmt.Errorf("Ошибка скачивания файла: %v", err)
    }
    defer os.Remove(tempFilePath) // Почистим за собой после отправки

    // Отправляем на принтер
    err = services.SendToPrinterRaw(tempFilePath, printerIP, printerPort)
    if err != nil {
        return fmt.Errorf("Ошибка отправки на принтер: %v", err)
    }
    return nil
}


// downloadFile сохраняет файл по HTTP-ссылке localPath
func downloadFile(localPath, url string) error {
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    out, err := os.Create(localPath)
    if err != nil {
        return err
    }
    defer out.Close()

    _, err = io.Copy(out, resp.Body)
    return err
}


// Создать задание на печать
func CreatePrintJob(c *gin.Context) {
    var job models.PrintJob
    if err := c.ShouldBindJSON(&job); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := config.DB.Create(&job).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, job)
}


// Получить все задания
func GetAllPrintJobs(c *gin.Context) {
    var jobs []models.PrintJob
    if err := config.DB.Find(&jobs).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, jobs)
}

// Отправить готовое задание на принтер
func SendPrintJobHandler(c *gin.Context) {
    jobID := c.Param("id")

    // 1. Ищем задание
    var job models.PrintJob
    if err := config.DB.First(&job, "id = ?", jobID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Задание не найдено"})
        return
    }

    // 2. Ищем принтер
    var printer models.Printer
    if err := config.DB.First(&printer, "id = ?", job.PrinterID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Принтер не найден"})
        return
    }

    // 3. Проверяем путь к файлу
    if job.FileURL == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Нет пути к файлу (job.FileURL) для печати"})
        return
    }

    // Если fileURL — HTTP-ссылка, скачиваем и отправляем
    err := DownloadAndSendToPrinter(job.ID, job.FileURL, printer.IPAddress, printer.Port)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // 4. Обновляем статус задания в БД
    job.Status = "printing"
    if err := config.DB.Save(&job).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message":  "Задание отправлено на принтер",
        "job_id":   jobID,
        "filePath": filepath.Base(job.FileURL),
    })
}

// Обновить статус задания
func UpdatePrintJob(c *gin.Context) {
    id := c.Param("id")
    var job models.PrintJob
    if err := config.DB.First(&job, "id = ?", id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Задание не найдено"})
        return
    }

    var input struct {
        Status string  `json:"status"`
        Cost   float64 `json:"cost"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    job.Status = input.Status
    job.Cost = input.Cost
    if err := config.DB.Save(&job).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, job)
}