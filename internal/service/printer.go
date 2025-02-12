// internal/service/printer.go
package service

import (
    "fmt"
    "io"
    "sync"  // Используется для мьютексов
    "time"
    "github.com/sirupsen/logrus"
)

// PrintJob представляет задание печати
type PrintJob struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"created_at"`
    UserID    string    `json:"user_id"`
    PrinterID string    `json:"printer_id"`
    Copies    int       `json:"copies"`
    Pages     int       `json:"pages"`
}

// PrintQueue представляет очередь печати
type PrintQueue struct {
    Jobs []PrintJob `json:"jobs"`
}

// PrinterService управляет принтерами и заданиями печати
type PrinterService struct {
    connections map[string]*PrinterConnection
    printQueue  map[string][]PrintJob // ключ - printer ID
    mu          sync.RWMutex         // мьютекс для потокобезопасности
}

// NewPrinterService создает новый экземпляр сервиса принтеров
func NewPrinterService() *PrinterService {
    return &PrinterService{
        connections: make(map[string]*PrinterConnection),
        printQueue:  make(map[string][]PrintJob),
    }
}

// PrintDocument отправляет документ на печать
func (s *PrinterService) PrintDocument(printerID string, document io.Reader) error {
    logger := logrus.WithFields(logrus.Fields{
        "printer_id": printerID,
        "action":     "print_document",
    })

    logger.Info("Starting print document process")

    s.mu.RLock()
    connection, exists := s.connections[printerID]
    s.mu.RUnlock()

    if !exists {
        logger.Error("Printer not found")
        return fmt.Errorf("printer not found")
    }

    if !connection.IsConnected() {
        logger.Error("Printer not connected")
        return fmt.Errorf("printer not connected")
    }

    // Создаем новое задание печати
    job := PrintJob{
        ID:        fmt.Sprintf("job_%d", time.Now().UnixNano()),
        Name:      "Document",
        Status:    "pending",
        CreatedAt: time.Now(),
        PrinterID: printerID,
    }

    logger = logger.WithField("job_id", job.ID)
    logger.Info("Created new print job")

    // Добавляем задание в очередь
    s.mu.Lock()
    if _, exists := s.printQueue[printerID]; !exists {
        s.printQueue[printerID] = make([]PrintJob, 0)
    }
    s.printQueue[printerID] = append(s.printQueue[printerID], job)
    s.mu.Unlock()

    logger.Info("Added job to print queue")

    // Обрабатываем задание печати асинхронно
    go func() {
        logger.Info("Starting async print processing")
        
        // Здесь должна быть реальная логика печати
        if err := connection.Print(document); err != nil {
            logger.WithError(err).Error("Failed to print document")
            s.updateJobStatus(printerID, job.ID, "failed")
            return
        }
        
        s.updateJobStatus(printerID, job.ID, "completed")
        logger.Info("Print job completed successfully")
    }()

    return nil
}

// updateJobStatus обновляет статус задания печати
func (s *PrinterService) updateJobStatus(printerID, jobID, status string) {
    s.mu.Lock()
    defer s.mu.Unlock()

    for i, job := range s.printQueue[printerID] {
        if job.ID == jobID {
            s.printQueue[printerID][i].Status = status
            break
        }
    }
}

// ConnectPrinter подключает принтер
func (s *PrinterService) ConnectPrinter(id, name, ip string, port int) (*PrinterInfo, error) {
    logger := logrus.WithFields(logrus.Fields{
        "printer_id": id,
        "ip":        ip,
        "port":      port,
    })

    logger.Info("Attempting to connect printer")

    s.mu.Lock()
    defer s.mu.Unlock()

    if conn, exists := s.connections[id]; exists && conn.IsConnected() {
        logger.Info("Printer already connected")
        return &conn.info, nil
    }

    printerInfo := PrinterInfo{
        IP:       ip,
        Name:     name,
        Port:     port,
        Protocol: "RAW",
    }

    connection := NewPrinterConnection(printerInfo)
    
    if err := connection.Connect(); err != nil {
        logger.WithError(err).Error("Failed to connect to printer")
        return nil, fmt.Errorf("failed to connect to printer: %w", err)
    }

    status, err := connection.GetStatus()
    if err != nil {
        logger.WithError(err).Error("Failed to get printer status")
        connection.Disconnect()
        return nil, fmt.Errorf("failed to get printer status: %w", err)
    }

    connection.info.Status = status
    connection.info.IsOnline = true

    s.connections[id] = connection
    logger.Info("Printer connected successfully")

    return &connection.info, nil
}

// StartStatusMonitoring запускает мониторинг статуса принтера
func (s *PrinterService) StartStatusMonitoring(id string, interval time.Duration) error {
    logger := logrus.WithField("printer_id", id)
    logger.Info("Starting printer status monitoring")

    s.mu.RLock()
    connection, exists := s.connections[id]
    s.mu.RUnlock()

    if !exists {
        logger.Error("Printer not found")
        return fmt.Errorf("printer not found")
    }

    go func() {
        ticker := time.NewTicker(interval)
        defer ticker.Stop()

        for range ticker.C {
            if !connection.IsConnected() {
                logger.Info("Printer disconnected, stopping monitoring")
                return
            }

            status, err := connection.GetStatus()
            if err != nil {
                logger.WithError(err).Error("Failed to get printer status")
                continue
            }

            connection.info.Status = status
            logger.WithField("status", status).Debug("Updated printer status")
        }
    }()

    return nil
}

// GetPrinterStatus получает текущий статус принтера
func (s *PrinterService) GetPrinterStatus(id string) (*PrinterInfo, error) {
    logger := logrus.WithField("printer_id", id)
    logger.Info("Getting printer status")

    s.mu.RLock()
    connection, exists := s.connections[id]
    s.mu.RUnlock()

    if !exists {
        logger.Error("Printer not found")
        return nil, fmt.Errorf("printer not found")
    }

    if !connection.IsConnected() {
        logger.Error("Printer not connected")
        return nil, fmt.Errorf("printer not connected")
    }

    status, err := connection.GetStatus()
    if err != nil {
        logger.WithError(err).Error("Failed to get printer status")
        return nil, fmt.Errorf("failed to get printer status: %v", err)
    }

    connection.info.Status = status
    return &connection.info, nil
}

// DisconnectPrinter отключает принтер
func (s *PrinterService) DisconnectPrinter(id string) error {
    logger := logrus.WithField("printer_id", id)
    logger.Info("Disconnecting printer")

    s.mu.Lock()
    defer s.mu.Unlock()

    connection, exists := s.connections[id]
    if !exists {
        logger.Error("Printer not found")
        return fmt.Errorf("printer not found")
    }

    err := connection.Disconnect()
    if err != nil {
        logger.WithError(err).Error("Failed to disconnect printer")
        return fmt.Errorf("failed to disconnect printer: %v", err)
    }

    delete(s.connections, id)
    logger.Info("Printer disconnected successfully")
    return nil
}

// DiscoverPrinters ищет доступные принтеры в сети
func (s *PrinterService) DiscoverPrinters() ([]PrinterInfo, error) {
    logger := logrus.WithField("action", "discover_printers")
    logger.Info("Starting printer discovery")

    // Используем константу для маски подсети или получаем из конфигурации
    subnetMask := "255.255.255.0"
    
    printers, err := DiscoverPrinters(subnetMask)
    if err != nil {
        logger.WithError(err).Error("Printer discovery failed")
        return nil, fmt.Errorf("printer discovery failed: %w", err)
    }

    // Добавляем логирование для отладки
    logger.WithField("printer_count", len(printers)).Info("Discovered printers")
    for _, printer := range printers {
        logger.WithFields(logrus.Fields{
            "name": printer.Name,
            "ip":   printer.IP,
            "port": printer.Port,
        }).Debug("Found printer")
    }

    return printers, nil
}

func (s *PrinterService) CancelPrint(printerID string) error {
    s.mu.RLock()
    connection, exists := s.connections[printerID]
    s.mu.RUnlock()

    if !exists {
        return fmt.Errorf("printer not found")
    }

    if !connection.IsConnected() {
        return fmt.Errorf("printer not connected")
    }

    return connection.CancelPrint()
}

// GetPrintQueue возвращает текущую очередь печати для принтера
func (s *PrinterService) GetPrintQueue(printerID string) (*PrintQueue, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    _, exists := s.connections[printerID]
    if !exists {
        return nil, fmt.Errorf("printer not found")
    }

    jobs, exists := s.printQueue[printerID]
    if !exists {
        return &PrintQueue{Jobs: make([]PrintJob, 0)}, nil
    }

    return &PrintQueue{Jobs: jobs}, nil
}