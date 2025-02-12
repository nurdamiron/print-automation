// internal/api/server.go
package api

import (
    "net/http"
    "github.com/gorilla/mux"
    "github.com/sirupsen/logrus"
    "print-automation/internal/config"
    "print-automation/internal/handlers"
    "print-automation/internal/middleware"
    "print-automation/internal/service"
)

type Server struct {
    Router         *mux.Router
    config         *config.Config
    authHandler    *handlers.AuthHandler
    printerHandler *handlers.PrinterHandler
    paymentHandler *handlers.PaymentHandler
    printJobHandler *handlers.PrintJobHandler
}

func NewServer(cfg *config.Config, 
    printerService *service.PrinterService,
    paymentService *service.PaymentService, 
    authService *service.AuthService,
    printJobService *service.PrintJobService) *Server {
    
    // Инициализируем обработчики
    authHandler := handlers.NewAuthHandler(authService)
    printerHandler := handlers.NewPrinterHandler(printerService)
    paymentHandler := handlers.NewPaymentHandler(paymentService)
    printJobHandler := handlers.NewPrintJobHandler(printJobService)

    s := &Server{
        Router:         mux.NewRouter(),
        config:         cfg,
        authHandler:    authHandler,
        printerHandler: printerHandler,
        paymentHandler: paymentHandler,
        printJobHandler: printJobHandler,
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
    api.HandleFunc("/print-jobs", s.printJobHandler.Create).Methods("POST")
    api.HandleFunc("/print-jobs/{id}", s.printJobHandler.Get).Methods("GET")
    api.HandleFunc("/print-jobs/{id}/status", s.printJobHandler.GetStatus).Methods("GET")
    
    // Payments
    api.HandleFunc("/payments", s.paymentHandler.ProcessPayment).Methods("POST")
    api.HandleFunc("/payments/{id}/callback", s.paymentHandler.PaymentCallback).Methods("POST")
    
    // Printers
    api.HandleFunc("/printers", s.printerHandler.DiscoverPrinters).Methods("GET")
    api.HandleFunc("/printers/{id}/status", s.printerHandler.GetPrinterStatus).Methods("GET")
    api.HandleFunc("/printers/{id}", s.printerHandler.ConnectPrinter).Methods("POST")
    api.HandleFunc("/printers/{id}", s.printerHandler.DisconnectPrinter).Methods("DELETE")
}

func (s *Server) Run() error {
    logger := logrus.WithFields(logrus.Fields{
        "host": s.config.Server.Host,
        "port": s.config.Server.Port,
    })
    
    addr := s.config.Server.Host + ":" + s.config.Server.Port
    logger.Info("Starting server")
    
    return http.ListenAndServe(addr, s.Router)
}

