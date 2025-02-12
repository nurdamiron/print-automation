// internal/handlers/printer.go
package handlers

import (
    "encoding/json"
    "net/http"
	"fmt"
    "print-automation/internal/service"
)

type PrinterHandler struct {
    printerService *service.PrinterService
}

func NewPrinterHandler(printerService *service.PrinterService) *PrinterHandler {
    return &PrinterHandler{
        printerService: printerService,
    }
}

// ConnectPrinter обрабатывает запрос на подключение к принтеру
func (h *PrinterHandler) ConnectPrinter(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req struct {
        ID   string `json:"id"`
        Name string `json:"name"`
        IP   string `json:"ip"`
        Port int    `json:"port"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    status, err := h.printerService.ConnectPrinter(req.ID, req.Name, req.IP, req.Port)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(status)
}

// GetPrinterStatus обрабатывает запрос на получение статуса принтера
func (h *PrinterHandler) GetPrinterStatus(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    printerID := r.URL.Query().Get("id")
    if printerID == "" {
        http.Error(w, "printer_id is required", http.StatusBadRequest)
        return
    }

    status, err := h.printerService.GetPrinterStatus(printerID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(status)
}

// DisconnectPrinter обрабатывает запрос на отключение принтера
func (h *PrinterHandler) DisconnectPrinter(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    printerID := r.URL.Query().Get("id")
    if printerID == "" {
        http.Error(w, "printer_id is required", http.StatusBadRequest)
        return
    }

    if err := h.printerService.DisconnectPrinter(printerID); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

// DiscoverPrinters обрабатывает запрос на поиск доступных принтеров
// DiscoverPrinters обрабатывает запрос на поиск доступных принтеров
func (h *PrinterHandler) DiscoverPrinters(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    printers, err := h.printerService.DiscoverPrinters()
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to discover printers: %v", err), http.StatusInternalServerError)
        return
    }

    // Передаем корректный тип данных из service
    response := struct {
        Message  string               `json:"message"`
        Printers []service.PrinterInfo `json:"printers"` // ✅ Указываем service.PrinterInfo
    }{
        Message:  fmt.Sprintf("Найдено принтеров: %d", len(printers)),
        Printers: printers,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}



// PrintDocument обрабатывает запрос на печать документа
func (h *PrinterHandler) PrintDocument(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    printerID := r.URL.Query().Get("id")
    if printerID == "" {
        http.Error(w, "printer_id is required", http.StatusBadRequest)
        return
    }

    // Максимальный размер файла - 32MB
    r.ParseMultipartForm(32 << 20)
    
    file, _, err := r.FormFile("document")
    if err != nil {
        http.Error(w, "Failed to get document file", http.StatusBadRequest)
        return
    }
    defer file.Close()

    if err := h.printerService.PrintDocument(printerID, file); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "status": "printing started",
    })
}

// GetPrinterQueue обрабатывает запрос на получение очереди печати
func (h *PrinterHandler) GetPrinterQueue(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    printerID := r.URL.Query().Get("id")
    if printerID == "" {
        http.Error(w, "printer_id is required", http.StatusBadRequest)
        return
    }

    queue, err := h.printerService.GetPrintQueue(printerID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(queue)
}