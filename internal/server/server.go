// internal/server/server.go
package server

import (
    "fmt"
    "net/http"
    "print-automation/internal/config"
)

type Server struct {
    config *config.Config
    router *http.ServeMux
}

func NewServer(cfg *config.Config) *Server {
    return &Server{
        config: cfg,
        router: http.NewServeMux(),
    }
}

func (s *Server) Run() error {
    // Регистрация маршрутов
    s.setupRoutes()

    addr := fmt.Sprintf("%s:%s", s.config.Server.Host, s.config.Server.Port)
    fmt.Printf("Server starting on %s\n", addr)
    
    return http.ListenAndServe(addr, s.router)
}

func (s *Server) setupRoutes() {
    s.router.HandleFunc("/health", s.handleHealth())
}

func (s *Server) handleHealth() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    }
}