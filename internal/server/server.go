// internal/server/server.go
package server

import (
    "fmt"
    "net/http"
    "print-automation/internal/config"
    "print-automation/internal/handlers"
)

type Handlers struct {
    UserHandler     *handlers.UserHandler
    PrintJobHandler *handlers.PrintJobHandler
    PaymentHandler  *handlers.PaymentHandler
    PrinterHandler  *handlers.PrinterHandler
}

type Server struct {
    cfg      *config.Config
    handlers *Handlers
    router   *http.ServeMux
}

func NewServer(cfg *config.Config, handlers *Handlers) *Server {
    return &Server{
        cfg:      cfg,
        handlers: handlers,
        router:   http.NewServeMux(),
    }
}

// handleEndpoint обрабатывает эндпоинт с проверкой метода
func (s *Server) handleEndpoint(method string, path string, handler func(http.ResponseWriter, *http.Request)) {
    s.router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
        if r.Method != method && r.Method != "OPTIONS" {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }
        handler(w, r)
    })
}

// enableCORS добавляет CORS заголовки ко всем ответам
func (s *Server) enableCORS(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Устанавливаем CORS заголовки
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        // Обрабатываем preflight запросы
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        handler.ServeHTTP(w, r)
    })
}

// setupRoutes настраивает все маршруты приложения
func (s *Server) setupRoutes() {
    // Маршруты аутентификации и пользователей
    s.handleEndpoint("POST", "/api/users", s.handlers.UserHandler.Create)
    
    // Маршруты заданий печати
    s.handleEndpoint("POST", "/api/print-jobs", s.handlers.PrintJobHandler.Create)
    
    // Маршруты платежей
    s.handleEndpoint("POST", "/api/payments", s.handlers.PaymentHandler.ProcessPayment)
    s.handleEndpoint("GET", "/api/payments/status", s.handlers.PaymentHandler.GetPaymentStatus)
    
    // Маршруты для работы с принтерами
    s.handleEndpoint("POST", "/api/printers/connect", s.handlers.PrinterHandler.ConnectPrinter)
    s.handleEndpoint("POST", "/api/printers/disconnect", s.handlers.PrinterHandler.DisconnectPrinter)
    s.handleEndpoint("GET", "/api/printers/status", s.handlers.PrinterHandler.GetPrinterStatus)
    s.handleEndpoint("GET", "/api/printers/discover", s.handlers.PrinterHandler.DiscoverPrinters)
    s.handleEndpoint("POST", "/api/printers/print", s.handlers.PrinterHandler.PrintDocument)
    s.handleEndpoint("GET", "/api/printers/queue", s.handlers.PrinterHandler.GetPrinterQueue)

    // Health check
    s.handleEndpoint("GET", "/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })
}

// Run запускает HTTP сервер
func (s *Server) Run() error {
    // Настраиваем маршруты
    s.setupRoutes()
    
    // Создаем адрес для прослушивания
    addr := fmt.Sprintf("%s:%s", s.cfg.Server.Host, s.cfg.Server.Port)
    
    // Оборачиваем роутер в CORS middleware
    handler := s.enableCORS(s.router)
    
    // Запускаем сервер
    return http.ListenAndServe(addr, handler)
}