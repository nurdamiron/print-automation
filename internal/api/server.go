// internal/api/server.go
package api

import (
    "print-automation/internal/config"
    "print-automation/internal/db"
    "print-automation/internal/handlers"
    "print-automation/internal/middleware"
    "print-automation/internal/service"
    "print-automation/internal/repository"
    "github.com/gorilla/mux"
)

type Server struct {
    Router      *mux.Router
    db          *db.Database
    config      *config.Config
    authHandler *handlers.AuthHandler
}

func NewServer(cfg *config.Config, db *db.Database) *Server {
    // Инициализируем необходимые репозитории
    userRepo := repository.NewUserRepository(db.DB)
    
    // Инициализируем сервисы
    authService := service.NewAuthService(userRepo, cfg.JWTSecret)
    
    // Инициализируем обработчики
    authHandler := handlers.NewAuthHandler(authService)

    s := &Server{
        Router:      mux.NewRouter(),
        db:          db,
        config:      cfg,
        authHandler: authHandler,
    }
    
    s.routes()
    return s
}

func (s *Server) routes() {
    // Auth routes
    s.Router.HandleFunc("/api/v1/auth/login", s.authHandler.Login).Methods("POST")
    s.Router.HandleFunc("/api/v1/auth/register", s.authHandler.Register).Methods("POST")

    // Protected routes
    api := s.Router.PathPrefix("/api/v1").Subrouter()
    api.Use(middleware.AuthMiddleware(s.config.JWTSecret))

    // Print Jobs
    api.HandleFunc("/print-jobs", s.CreatePrintJob).Methods("POST")
    api.HandleFunc("/print-jobs/{id}", s.GetPrintJob).Methods("GET")
    api.HandleFunc("/print-jobs/{id}/status", s.GetPrintJobStatus).Methods("GET")
    
    // Payments
    api.HandleFunc("/payments", s.CreatePayment).Methods("POST")
    api.HandleFunc("/payments/{id}/callback", s.PaymentCallback).Methods("POST")
    
    // Printers
    api.HandleFunc("/printers", s.ListPrinters).Methods("GET")
    api.HandleFunc("/printers/{id}/status", s.UpdatePrinterStatus).Methods("PUT")
}