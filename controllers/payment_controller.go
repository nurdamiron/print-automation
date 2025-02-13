package controllers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "print-automation/config"
    "print-automation/models"
)

// Создать платёжную запись
func CreatePayment(c *gin.Context) {
    var payment models.Payment
    if err := c.ShouldBindJSON(&payment); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := config.DB.Create(&payment).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, payment)
}

// Получить все платежи
func GetAllPayments(c *gin.Context) {
    var payments []models.Payment
    if err := config.DB.Find(&payments).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, payments)
}

// Обновить статус платежа
func UpdatePayment(c *gin.Context) {
    id := c.Param("id")
    var payment models.Payment
    if err := config.DB.First(&payment, "id = ?", id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Платёж не найден"})
        return
    }

    var input struct {
        Status        string `json:"status"`
        PaymentMethod string `json:"payment_method"`
        TransactionID string `json:"transaction_id"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    payment.Status = input.Status
    payment.PaymentMethod = input.PaymentMethod
    payment.TransactionID = input.TransactionID

    if err := config.DB.Save(&payment).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, payment)
}
