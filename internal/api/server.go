// internal/api/server.go
package api

import (
    "github.com/gorilla/mux"
    "print-automation/internal/db"
    "print-automation/internal/config"
)

type Server struct {
    Router *mux.Router
    db     *db.Database
    config *config.Config
}

func NewServer(cfg *config.Config, db *db.Database) *Server {
    s := &Server{
        Router: mux.NewRouter(),
        db:     db,
        config: cfg,
    }
    s.routes()
    return s
}

func (s *Server) routes() {
    // Print Jobs
    s.Router.HandleFunc("/api/v1/print-jobs", s.CreatePrintJob).Methods("POST")
    s.Router.HandleFunc("/api/v1/print-jobs/{id}", s.GetPrintJob).Methods("GET")
    s.Router.HandleFunc("/api/v1/print-jobs/{id}/status", s.GetPrintJobStatus).Methods("GET")
    
    // Payments
    s.Router.HandleFunc("/api/v1/payments", s.CreatePayment).Methods("POST")
    s.Router.HandleFunc("/api/v1/payments/{id}/callback", s.PaymentCallback).Methods("POST")
    
    // Printers
    s.Router.HandleFunc("/api/v1/printers", s.ListPrinters).Methods("GET")
    s.Router.HandleFunc("/api/v1/printers/{id}/status", s.UpdatePrinterStatus).Methods("PUT")
}