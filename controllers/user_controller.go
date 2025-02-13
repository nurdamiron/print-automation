package controllers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "print-automation/config"
    "print-automation/models"
    "golang.org/x/crypto/bcrypt"
)

// Регистрация пользователя
func CreateUser(c *gin.Context) {
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Хеширование пароля
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось хешировать пароль"})
        return
    }
    user.PasswordHash = string(hashedPassword)

    if err := config.DB.Create(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, user)
}

// Получить пользователя по ID
func GetUserByID(c *gin.Context) {
    id := c.Param("id")
    var user models.User
    if err := config.DB.First(&user, "id = ?", id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
        return
    }
    c.JSON(http.StatusOK, user)
}

// Авторизация (упрощённый пример)
func LoginUser(c *gin.Context) {
    var credentials struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    if err := c.ShouldBindJSON(&credentials); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user models.User
    if err := config.DB.Where("email = ?", credentials.Email).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный email или пароль"})
        return
    }

    // Проверка хеша
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный email или пароль"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Успешная авторизация", "userId": user.ID})
}
