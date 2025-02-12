
// internal/handlers/print_jobs.go
package handlers

import (
    "encoding/json"
    "net/http"
    "print-automation/internal/models"
    "print-automation/internal/repository"
)

type PrintJobHandler struct {
    repo *repository.PrintJobRepository
}

func NewPrintJobHandler(repo *repository.PrintJobRepository) *PrintJobHandler {
    return &PrintJobHandler{repo: repo}
}

func (h *PrintJobHandler) Create(w http.ResponseWriter, r *http.Request) {
    var job models.PrintJob
    if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.repo.Create(&job); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(job)
}
