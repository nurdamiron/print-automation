// internal/handlers/print_jobs.go

package handlers

import (
    "encoding/json"
    "net/http"
    "fmt"
    "github.com/gorilla/mux"
    "github.com/sirupsen/logrus"
    "print-automation/internal/service"
)

type PrintJobHandler struct {
    printJobService *service.PrintJobService
}

func NewPrintJobHandler(printJobService *service.PrintJobService) *PrintJobHandler {
    return &PrintJobHandler{
        printJobService: printJobService,
    }
}

func (h *PrintJobHandler) Create(w http.ResponseWriter, r *http.Request) {
    logger := logrus.WithField("handler", "print_job_create")
    
    // Parse the multipart form data
    err := r.ParseMultipartForm(32 << 20) // 32MB max
    if err != nil {
        logger.WithError(err).Error("Failed to parse multipart form")
        http.Error(w, "Failed to parse form", http.StatusBadRequest)
        return
    }

    // Get the file
    file, _, err := r.FormFile("document")
    if err != nil {
        logger.WithError(err).Error("Failed to get document file")
        http.Error(w, "Failed to get document file", http.StatusBadRequest)
        return
    }
    defer file.Close()

    // Get print options from form
    copies := 1 // default value
    if copiesStr := r.FormValue("copies"); copiesStr != "" {
        fmt.Sscanf(copiesStr, "%d", &copies)
    }

    options := service.PrintOptions{
        Copies: copies,
        Pages: 1, // You might want to calculate this from the document
        DoubleSided: r.FormValue("double_sided") == "true",
        Color: r.FormValue("color") == "true",
    }

    // Get user ID from context
    userID := r.Context().Value("user_id").(string)
    printerID := r.FormValue("printer_id")

    // Create the print job
    job, err := h.printJobService.Create(userID, printerID, file, options)
    if err != nil {
        logger.WithError(err).Error("Failed to create print job")
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(job)
}

func (h *PrintJobHandler) Get(w http.ResponseWriter, r *http.Request) {
    logger := logrus.WithField("handler", "print_job_get")
    
    vars := mux.Vars(r)
    jobID := vars["id"]

    job, err := h.printJobService.GetByID(jobID)
    if err != nil {
        logger.WithError(err).Error("Failed to get print job")
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if job == nil {
        http.Error(w, "Print job not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(job)
}

func (h *PrintJobHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
    logger := logrus.WithField("handler", "print_job_status")
    
    vars := mux.Vars(r)
    jobID := vars["id"]

    job, err := h.printJobService.GetByID(jobID)
    if err != nil {
        logger.WithError(err).Error("Failed to get print job")
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if job == nil {
        http.Error(w, "Print job not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "status": job.Status,
    })
}