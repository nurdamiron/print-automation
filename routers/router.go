package routers

import (
    "github.com/gin-gonic/gin"
    "print-automation/controllers"
	"github.com/gin-contrib/cors"
	"time"


)

func SetupRouter() *gin.Engine {
    r := gin.Default()

	r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
        AllowHeaders:     []string{"Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge: 12 * time.Hour, 
    }))

    // Пример роутов для пользователей
    r.POST("/users", controllers.CreateUser)
    r.POST("/users/login", controllers.LoginUser)
    r.GET("/users/:id", controllers.GetUserByID)

    // Принтеры
    r.GET("/printers", controllers.GetAllPrinters)
    r.POST("/printers", controllers.CreatePrinter)
    r.GET("/printers/:id", controllers.GetPrinterByID)
    r.PUT("/printers/:id", controllers.UpdatePrinter)
    r.DELETE("/printers/:id", controllers.DeletePrinter)
	r.GET("/printers/:id/check", controllers.CheckPrinterConnectionHandler)

    // Задания на печать
    r.GET("/printjobs", controllers.GetAllPrintJobs)
    r.POST("/printjobs", controllers.CreatePrintJob)
    r.PUT("/printjobs/:id", controllers.UpdatePrintJob)
	r.POST("/printjobs/:id/send", controllers.SendPrintJobHandler)

    // Платежи
    r.GET("/payments", controllers.GetAllPayments)
    r.POST("/payments", controllers.CreatePayment)
    r.PUT("/payments/:id", controllers.UpdatePayment)

    return r
}
